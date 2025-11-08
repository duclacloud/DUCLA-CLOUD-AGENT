package executor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ducla/cloud-agent/internal/config"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Executor manages task execution
type Executor struct {
	config config.ExecutorConfig
	logger *logrus.Logger

	// Task management
	mu            sync.RWMutex
	tasks         map[string]*Task
	runningTasks  map[string]*Task
	completedTasks map[string]*Task

	// Worker pool
	workers    []*Worker
	taskQueue  chan *Task
	resultChan chan *TaskResult

	// Lifecycle
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// Task represents a task to be executed
type Task struct {
	ID          string                 `json:"id"`
	Type        TaskType               `json:"type"`
	Name        string                 `json:"name"`
	Command     string                 `json:"command"`
	Args        []string               `json:"args"`
	Env         map[string]string      `json:"env"`
	WorkingDir  string                 `json:"working_dir"`
	Timeout     time.Duration          `json:"timeout"`
	Priority    int                    `json:"priority"`
	Metadata    map[string]interface{} `json:"metadata"`
	
	// Execution state
	Status      TaskStatus    `json:"status"`
	StartedAt   time.Time     `json:"started_at,omitempty"`
	FinishedAt  time.Time     `json:"finished_at,omitempty"`
	Result      *TaskResult   `json:"result,omitempty"`
	Error       error         `json:"error,omitempty"`
	
	// Cancellation
	ctx         context.Context
	cancel      context.CancelFunc
}

// TaskType represents the type of task
type TaskType string

const (
	TaskTypeCommand    TaskType = "command"
	TaskTypeScript     TaskType = "script"
	TaskTypeFile       TaskType = "file"
	TaskTypeHTTP       TaskType = "http"
	TaskTypeDocker     TaskType = "docker"
	TaskTypeKubernetes TaskType = "kubernetes"
	TaskTypeCustom     TaskType = "custom"
)

// TaskStatus represents the status of a task
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusQueued    TaskStatus = "queued"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
	TaskStatusCancelled TaskStatus = "cancelled"
	TaskStatusTimeout   TaskStatus = "timeout"
)

// TaskResult represents the result of task execution
type TaskResult struct {
	TaskID       string                 `json:"task_id"`
	Status       TaskStatus             `json:"status"`
	ExitCode     int                    `json:"exit_code"`
	Output       string                 `json:"output"`
	Error        string                 `json:"error"`
	StartedAt    time.Time              `json:"started_at"`
	FinishedAt   time.Time              `json:"finished_at"`
	Duration     time.Duration          `json:"duration"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// New creates a new executor instance
func New(cfg config.ExecutorConfig, logger *logrus.Logger) (*Executor, error) {
	executor := &Executor{
		config:         cfg,
		logger:         logger,
		tasks:          make(map[string]*Task),
		runningTasks:   make(map[string]*Task),
		completedTasks: make(map[string]*Task),
		taskQueue:      make(chan *Task, cfg.QueueSize),
		resultChan:     make(chan *TaskResult, cfg.QueueSize),
		workers:        make([]*Worker, cfg.WorkerPoolSize),
	}

	return executor, nil
}

// Start starts the executor and worker pool
func (e *Executor) Start(ctx context.Context) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.logger.Info("Starting task executor")

	e.ctx, e.cancel = context.WithCancel(ctx)

	// Start workers
	for i := 0; i < e.config.WorkerPoolSize; i++ {
		worker := NewWorker(i, e.taskQueue, e.resultChan, e.logger)
		e.workers[i] = worker

		e.wg.Add(1)
		go func(w *Worker) {
			defer e.wg.Done()
			w.Start(e.ctx)
		}(worker)
	}

	// Start result handler
	e.wg.Add(1)
	go e.handleResults()

	e.logger.WithField("workers", e.config.WorkerPoolSize).Info("Task executor started")
	return nil
}

// Stop stops the executor and all workers
func (e *Executor) Stop(ctx context.Context) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.logger.Info("Stopping task executor")

	// Cancel all running tasks
	for _, task := range e.runningTasks {
		if task.cancel != nil {
			task.cancel()
		}
	}

	// Cancel context
	if e.cancel != nil {
		e.cancel()
	}

	// Close channels
	close(e.taskQueue)

	// Wait for workers to finish
	e.wg.Wait()

	close(e.resultChan)

	e.logger.Info("Task executor stopped")
	return nil
}

// Name returns the service name
func (e *Executor) Name() string {
	return "executor"
}

// ExecuteTask executes a task
func (e *Executor) ExecuteTask(ctx context.Context, task *Task) (*TaskResult, error) {
	// Validate task
	if err := e.validateTask(task); err != nil {
		return nil, fmt.Errorf("invalid task: %w", err)
	}

	// Set task ID if not provided
	if task.ID == "" {
		task.ID = uuid.New().String()
	}

	// Set default timeout
	if task.Timeout == 0 {
		task.Timeout = e.config.TaskTimeout
	}

	// Create task context with timeout
	taskCtx, cancel := context.WithTimeout(ctx, task.Timeout)
	task.ctx = taskCtx
	task.cancel = cancel

	// Set initial status
	task.Status = TaskStatusQueued

	// Store task
	e.mu.Lock()
	e.tasks[task.ID] = task
	e.mu.Unlock()

	e.logger.WithFields(logrus.Fields{
		"task_id": task.ID,
		"type":    task.Type,
		"name":    task.Name,
	}).Info("Task queued for execution")

	// Queue task
	select {
	case e.taskQueue <- task:
		// Task queued successfully
	case <-ctx.Done():
		cancel()
		return nil, ctx.Err()
	case <-time.After(5 * time.Second):
		cancel()
		return nil, fmt.Errorf("task queue is full")
	}

	// Wait for result
	select {
	case <-taskCtx.Done():
		if taskCtx.Err() == context.DeadlineExceeded {
			task.Status = TaskStatusTimeout
			return &TaskResult{
				TaskID:     task.ID,
				Status:     TaskStatusTimeout,
				Error:      "task execution timeout",
				FinishedAt: time.Now(),
			}, nil
		}
		return nil, taskCtx.Err()
	case <-ctx.Done():
		cancel()
		return nil, ctx.Err()
	}
}

// SubmitTask submits a task for asynchronous execution
func (e *Executor) SubmitTask(task *Task) (string, error) {
	// Validate task
	if err := e.validateTask(task); err != nil {
		return "", fmt.Errorf("invalid task: %w", err)
	}

	// Set task ID if not provided
	if task.ID == "" {
		task.ID = uuid.New().String()
	}

	// Set default timeout
	if task.Timeout == 0 {
		task.Timeout = e.config.TaskTimeout
	}

	// Create task context with timeout
	taskCtx, cancel := context.WithTimeout(e.ctx, task.Timeout)
	task.ctx = taskCtx
	task.cancel = cancel

	// Set initial status
	task.Status = TaskStatusQueued

	// Store task
	e.mu.Lock()
	e.tasks[task.ID] = task
	e.mu.Unlock()

	e.logger.WithFields(logrus.Fields{
		"task_id": task.ID,
		"type":    task.Type,
		"name":    task.Name,
	}).Info("Task submitted for execution")

	// Queue task
	select {
	case e.taskQueue <- task:
		return task.ID, nil
	case <-time.After(5 * time.Second):
		cancel()
		return "", fmt.Errorf("task queue is full")
	}
}

// GetTask retrieves a task by ID
func (e *Executor) GetTask(taskID string) (*Task, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	task, exists := e.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}

	return task, nil
}

// GetTaskResult retrieves the result of a task
func (e *Executor) GetTaskResult(taskID string) (*TaskResult, error) {
	task, err := e.GetTask(taskID)
	if err != nil {
		return nil, err
	}

	if task.Result == nil {
		return nil, fmt.Errorf("task result not available yet")
	}

	return task.Result, nil
}

// CancelTask cancels a running task
func (e *Executor) CancelTask(taskID string) error {
	e.mu.RLock()
	task, exists := e.tasks[taskID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("task not found: %s", taskID)
	}

	if task.Status != TaskStatusRunning && task.Status != TaskStatusQueued {
		return fmt.Errorf("task cannot be cancelled (status: %s)", task.Status)
	}

	if task.cancel != nil {
		task.cancel()
	}

	task.Status = TaskStatusCancelled

	e.logger.WithField("task_id", taskID).Info("Task cancelled")
	return nil
}

// ListTasks returns all tasks
func (e *Executor) ListTasks() []*Task {
	e.mu.RLock()
	defer e.mu.RUnlock()

	tasks := make([]*Task, 0, len(e.tasks))
	for _, task := range e.tasks {
		tasks = append(tasks, task)
	}

	return tasks
}

// ListRunningTasks returns all running tasks
func (e *Executor) ListRunningTasks() []*Task {
	e.mu.RLock()
	defer e.mu.RUnlock()

	tasks := make([]*Task, 0, len(e.runningTasks))
	for _, task := range e.runningTasks {
		tasks = append(tasks, task)
	}

	return tasks
}

// GetStats returns executor statistics
func (e *Executor) GetStats() map[string]interface{} {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return map[string]interface{}{
		"total_tasks":     len(e.tasks),
		"running_tasks":   len(e.runningTasks),
		"completed_tasks": len(e.completedTasks),
		"queue_size":      len(e.taskQueue),
		"worker_count":    len(e.workers),
	}
}

// handleResults processes task results
func (e *Executor) handleResults() {
	defer e.wg.Done()

	for {
		select {
		case <-e.ctx.Done():
			return
		case result, ok := <-e.resultChan:
			if !ok {
				return
			}

			e.processResult(result)
		}
	}
}

// processResult processes a single task result
func (e *Executor) processResult(result *TaskResult) {
	e.mu.Lock()
	defer e.mu.Unlock()

	task, exists := e.tasks[result.TaskID]
	if !exists {
		e.logger.WithField("task_id", result.TaskID).Warn("Received result for unknown task")
		return
	}

	// Update task with result
	task.Result = result
	task.Status = result.Status
	task.FinishedAt = result.FinishedAt

	// Move from running to completed
	delete(e.runningTasks, task.ID)
	e.completedTasks[task.ID] = task

	// Cancel task context
	if task.cancel != nil {
		task.cancel()
	}

	e.logger.WithFields(logrus.Fields{
		"task_id":  task.ID,
		"status":   result.Status,
		"duration": result.Duration,
	}).Info("Task completed")
}

// validateTask validates a task
func (e *Executor) validateTask(task *Task) error {
	if task == nil {
		return fmt.Errorf("task is nil")
	}

	if task.Type == "" {
		return fmt.Errorf("task type is required")
	}

	if task.Command == "" && task.Type == TaskTypeCommand {
		return fmt.Errorf("command is required for command tasks")
	}

	return nil
}

// ParseTask parses task data from a map
func ParseTask(data map[string]interface{}) (*Task, error) {
	task := &Task{
		Env:      make(map[string]string),
		Metadata: make(map[string]interface{}),
	}

	// Parse required fields
	if id, ok := data["id"].(string); ok {
		task.ID = id
	}

	if taskType, ok := data["type"].(string); ok {
		task.Type = TaskType(taskType)
	}

	if name, ok := data["name"].(string); ok {
		task.Name = name
	}

	if command, ok := data["command"].(string); ok {
		task.Command = command
	}

	// Parse args
	if args, ok := data["args"].([]interface{}); ok {
		task.Args = make([]string, len(args))
		for i, arg := range args {
			if argStr, ok := arg.(string); ok {
				task.Args[i] = argStr
			}
		}
	}

	// Parse env
	if env, ok := data["env"].(map[string]interface{}); ok {
		for key, value := range env {
			if valueStr, ok := value.(string); ok {
				task.Env[key] = valueStr
			}
		}
	}

	// Parse working directory
	if workingDir, ok := data["working_dir"].(string); ok {
		task.WorkingDir = workingDir
	}

	// Parse timeout
	if timeout, ok := data["timeout"].(float64); ok {
		task.Timeout = time.Duration(timeout) * time.Second
	}

	// Parse priority
	if priority, ok := data["priority"].(float64); ok {
		task.Priority = int(priority)
	}

	// Parse metadata
	if metadata, ok := data["metadata"].(map[string]interface{}); ok {
		task.Metadata = metadata
	}

	return task, nil
}