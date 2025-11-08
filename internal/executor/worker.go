package executor

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// Worker represents a task worker
type Worker struct {
	id         int
	taskQueue  <-chan *Task
	resultChan chan<- *TaskResult
	logger     *logrus.Logger
}

// NewWorker creates a new worker
func NewWorker(id int, taskQueue <-chan *Task, resultChan chan<- *TaskResult, logger *logrus.Logger) *Worker {
	return &Worker{
		id:         id,
		taskQueue:  taskQueue,
		resultChan: resultChan,
		logger:     logger,
	}
}

// Start starts the worker
func (w *Worker) Start(ctx context.Context) {
	w.logger.WithField("worker_id", w.id).Info("Worker started")

	for {
		select {
		case <-ctx.Done():
			w.logger.WithField("worker_id", w.id).Info("Worker stopped")
			return
		case task, ok := <-w.taskQueue:
			if !ok {
				w.logger.WithField("worker_id", w.id).Info("Task queue closed, worker stopping")
				return
			}

			w.executeTask(ctx, task)
		}
	}
}

// executeTask executes a single task
func (w *Worker) executeTask(ctx context.Context, task *Task) {
	w.logger.WithFields(logrus.Fields{
		"worker_id": w.id,
		"task_id":   task.ID,
		"task_type": task.Type,
	}).Info("Executing task")

	// Update task status
	task.Status = TaskStatusRunning
	task.StartedAt = time.Now()

	// Create result
	result := &TaskResult{
		TaskID:    task.ID,
		StartedAt: task.StartedAt,
		Metadata:  make(map[string]interface{}),
	}

	// Execute based on task type
	var err error
	switch task.Type {
	case TaskTypeCommand:
		err = w.executeCommand(task.ctx, task, result)
	case TaskTypeScript:
		err = w.executeScript(task.ctx, task, result)
	case TaskTypeFile:
		err = w.executeFileOperation(task.ctx, task, result)
	case TaskTypeHTTP:
		err = w.executeHTTPRequest(task.ctx, task, result)
	case TaskTypeDocker:
		err = w.executeDockerTask(task.ctx, task, result)
	case TaskTypeKubernetes:
		err = w.executeKubernetesTask(task.ctx, task, result)
	case TaskTypeCustom:
		err = w.executeCustomTask(task.ctx, task, result)
	default:
		err = fmt.Errorf("unsupported task type: %s", task.Type)
	}

	// Update result
	result.FinishedAt = time.Now()
	result.Duration = result.FinishedAt.Sub(result.StartedAt)

	if err != nil {
		result.Status = TaskStatusFailed
		result.Error = err.Error()
		w.logger.WithError(err).WithFields(logrus.Fields{
			"worker_id": w.id,
			"task_id":   task.ID,
		}).Error("Task execution failed")
	} else {
		result.Status = TaskStatusCompleted
		w.logger.WithFields(logrus.Fields{
			"worker_id": w.id,
			"task_id":   task.ID,
			"duration":  result.Duration,
		}).Info("Task execution completed")
	}

	// Send result
	select {
	case w.resultChan <- result:
	case <-ctx.Done():
		w.logger.WithField("task_id", task.ID).Warn("Context cancelled while sending result")
	case <-time.After(5 * time.Second):
		w.logger.WithField("task_id", task.ID).Error("Timeout sending result")
	}
}

// executeCommand executes a command task
func (w *Worker) executeCommand(ctx context.Context, task *Task, result *TaskResult) error {
	executor := NewCommandExecutor(w.logger)
	return executor.Execute(ctx, task, result)
}

// executeScript executes a script task
func (w *Worker) executeScript(ctx context.Context, task *Task, result *TaskResult) error {
	executor := NewScriptExecutor(w.logger)
	return executor.Execute(ctx, task, result)
}

// executeFileOperation executes a file operation task
func (w *Worker) executeFileOperation(ctx context.Context, task *Task, result *TaskResult) error {
	executor := NewFileExecutor(w.logger)
	return executor.Execute(ctx, task, result)
}

// executeHTTPRequest executes an HTTP request task
func (w *Worker) executeHTTPRequest(ctx context.Context, task *Task, result *TaskResult) error {
	executor := NewHTTPExecutor(w.logger)
	return executor.Execute(ctx, task, result)
}

// executeDockerTask executes a Docker task
func (w *Worker) executeDockerTask(ctx context.Context, task *Task, result *TaskResult) error {
	executor := NewDockerExecutor(w.logger)
	return executor.Execute(ctx, task, result)
}

// executeKubernetesTask executes a Kubernetes task
func (w *Worker) executeKubernetesTask(ctx context.Context, task *Task, result *TaskResult) error {
	executor := NewKubernetesExecutor(w.logger)
	return executor.Execute(ctx, task, result)
}

// executeCustomTask executes a custom task
func (w *Worker) executeCustomTask(ctx context.Context, task *Task, result *TaskResult) error {
	executor := NewCustomExecutor(w.logger)
	return executor.Execute(ctx, task, result)
}

// TaskExecutor interface for different task types
type TaskExecutor interface {
	Execute(ctx context.Context, task *Task, result *TaskResult) error
}