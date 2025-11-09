package agent

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/duclacloud/DUCLA-CLOUD-AGENT/internal/api"
	"github.com/duclacloud/DUCLA-CLOUD-AGENT/internal/config"
	"github.com/duclacloud/DUCLA-CLOUD-AGENT/internal/executor"
	"github.com/duclacloud/DUCLA-CLOUD-AGENT/internal/fileops"
	"github.com/duclacloud/DUCLA-CLOUD-AGENT/internal/health"
	"github.com/duclacloud/DUCLA-CLOUD-AGENT/internal/metrics"
	"github.com/duclacloud/DUCLA-CLOUD-AGENT/internal/transport"
	"github.com/sirupsen/logrus"
)

// Agent represents the main agent instance
type Agent struct {
	config    *config.Config
	logger    *logrus.Logger
	transport transport.Transport
	api       *api.Server
	executor  *executor.Executor
	fileops   *fileops.Manager
	health    *health.Checker
	metrics   *metrics.Collector
	
	// Internal state
	mu       sync.RWMutex
	running  bool
	services []Service
}

// Service represents a service that can be started and stopped
type Service interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Name() string
}

// New creates a new agent instance
func New(cfg *config.Config, logger *logrus.Logger) (*Agent, error) {
	agent := &Agent{
		config:   cfg,
		logger:   logger,
		services: make([]Service, 0),
	}

	// Initialize transport layer (only if master URL is provided)
	if cfg.Master.URL != "" {
		transportClient, err := transport.New(cfg.Master, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to create transport: %w", err)
		}
		agent.transport = transportClient
	}

	// Initialize executor
	executorInstance, err := executor.New(cfg.Executor, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create executor: %w", err)
	}
	agent.executor = executorInstance
	agent.services = append(agent.services, executorInstance)

	// Initialize file operations manager
	fileopsManager, err := fileops.New(cfg.Storage, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create fileops manager: %w", err)
	}
	agent.fileops = fileopsManager

	// Initialize health checker
	if cfg.Health.Enabled {
		healthChecker, err := health.New(cfg.Health, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to create health checker: %w", err)
		}
		agent.health = healthChecker
		agent.services = append(agent.services, healthChecker)
	}

	// Initialize metrics collector
	if cfg.Metrics.Enabled {
		metricsCollector, err := metrics.New(cfg.Metrics, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to create metrics collector: %w", err)
		}
		agent.metrics = metricsCollector
		agent.services = append(agent.services, metricsCollector)
	}

	// Initialize API server
	apiServer, err := api.New(cfg.API, agent, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create API server: %w", err)
	}
	agent.api = apiServer
	agent.services = append(agent.services, apiServer)

	return agent, nil
}

// Start starts the agent and all its services
func (a *Agent) Start(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.running {
		return fmt.Errorf("agent is already running")
	}

	a.logger.Info("Starting Ducla Cloud Agent")

	// Start all services
	for _, service := range a.services {
		a.logger.WithField("service", service.Name()).Info("Starting service")
		if err := service.Start(ctx); err != nil {
			a.logger.WithError(err).WithField("service", service.Name()).Error("Failed to start service")
			// Stop already started services
			a.stopServices(ctx)
			return fmt.Errorf("failed to start service %s: %w", service.Name(), err)
		}
		a.logger.WithField("service", service.Name()).Info("Service started successfully")
	}

	// Connect to master server (if transport is configured)
	if a.transport != nil {
		if err := a.transport.Connect(ctx); err != nil {
			a.logger.WithError(err).Error("Failed to connect to master server")
			a.stopServices(ctx)
			return fmt.Errorf("failed to connect to master: %w", err)
		}
	} else {
		a.logger.Info("Running in standalone mode (no master server)")
	}

	// Start heartbeat
	go a.heartbeatLoop(ctx)

	// Start message handling
	go a.messageLoop(ctx)

	a.running = true
	a.logger.Info("Ducla Cloud Agent started successfully")

	return nil
}

// Stop stops the agent and all its services
func (a *Agent) Stop(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.running {
		return nil
	}

	a.logger.Info("Stopping Ducla Cloud Agent")

	// Disconnect from master (if transport is configured)
	if a.transport != nil {
		if err := a.transport.Disconnect(); err != nil {
			a.logger.WithError(err).Error("Error disconnecting from master")
		}
	}

	// Stop all services
	a.stopServices(ctx)

	a.running = false
	a.logger.Info("Ducla Cloud Agent stopped")

	return nil
}

// stopServices stops all services in reverse order
func (a *Agent) stopServices(ctx context.Context) {
	for i := len(a.services) - 1; i >= 0; i-- {
		service := a.services[i]
		a.logger.WithField("service", service.Name()).Info("Stopping service")
		if err := service.Stop(ctx); err != nil {
			a.logger.WithError(err).WithField("service", service.Name()).Error("Error stopping service")
		} else {
			a.logger.WithField("service", service.Name()).Info("Service stopped")
		}
	}
}

// heartbeatLoop sends periodic heartbeats to the master server
func (a *Agent) heartbeatLoop(ctx context.Context) {
	ticker := time.NewTicker(a.config.Master.HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := a.sendHeartbeat(); err != nil {
				a.logger.WithError(err).Error("Failed to send heartbeat")
			}
		}
	}
}

// messageLoop handles incoming messages from the master server
func (a *Agent) messageLoop(ctx context.Context) {
	// Skip message loop if no transport (standalone mode)
	if a.transport == nil {
		<-ctx.Done()
		return
	}
	
	for {
		select {
		case <-ctx.Done():
			return
		default:
			message, err := a.transport.ReceiveMessage(ctx)
			if err != nil {
				a.logger.WithError(err).Error("Failed to receive message")
				time.Sleep(time.Second) // Avoid tight loop on errors
				continue
			}

			go a.handleMessage(ctx, message)
		}
	}
}

// sendHeartbeat sends a heartbeat message to the master server
func (a *Agent) sendHeartbeat() error {
	// Skip heartbeat if no transport (standalone mode)
	if a.transport == nil {
		return nil
	}
	
	heartbeat := &transport.Message{
		Type: transport.MessageTypeHeartbeat,
		Data: map[string]interface{}{
			"agent_id":   a.config.Agent.ID,
			"timestamp":  time.Now().Unix(),
			"status":     "healthy",
			"version":    "1.0.0", // TODO: Get from build info
		},
	}

	return a.transport.SendMessage(heartbeat)
}

// handleMessage processes incoming messages from the master server
func (a *Agent) handleMessage(ctx context.Context, message *transport.Message) {
	a.logger.WithFields(logrus.Fields{
		"type":       message.Type,
		"message_id": message.ID,
	}).Debug("Handling message")

	switch message.Type {
	case transport.MessageTypeTask:
		a.handleTaskMessage(ctx, message)
	case transport.MessageTypeFileOperation:
		a.handleFileOperationMessage(ctx, message)
	case transport.MessageTypeHealthCheck:
		a.handleHealthCheckMessage(ctx, message)
	case transport.MessageTypeConfig:
		a.handleConfigMessage(ctx, message)
	default:
		a.logger.WithField("type", message.Type).Warn("Unknown message type")
	}
}

// handleTaskMessage handles task execution messages
func (a *Agent) handleTaskMessage(ctx context.Context, message *transport.Message) {
	task, err := executor.ParseTask(message.Data)
	if err != nil {
		a.logger.WithError(err).Error("Failed to parse task")
		a.sendErrorResponse(message, err)
		return
	}

	// Execute task
	result, err := a.executor.ExecuteTask(ctx, task)
	if err != nil {
		a.logger.WithError(err).Error("Failed to execute task")
		a.sendErrorResponse(message, err)
		return
	}

	// Send result back to master
	response := &transport.Message{
		Type:      transport.MessageTypeTaskResult,
		ID:        message.ID,
		ReplyTo:   message.ID,
		Data:      map[string]interface{}{
			"task_id": result.TaskID,
			"status":  result.Status,
			"output":  result.Output,
			"error":   result.Error,
		},
	}

	if err := a.transport.SendMessage(response); err != nil {
		a.logger.WithError(err).Error("Failed to send task result")
	}
}

// handleFileOperationMessage handles file operation messages
func (a *Agent) handleFileOperationMessage(ctx context.Context, message *transport.Message) {
	operation, err := fileops.ParseOperation(message.Data)
	if err != nil {
		a.logger.WithError(err).Error("Failed to parse file operation")
		a.sendErrorResponse(message, err)
		return
	}

	// Execute file operation
	result, err := a.fileops.ExecuteOperation(ctx, operation)
	if err != nil {
		a.logger.WithError(err).Error("Failed to execute file operation")
		a.sendErrorResponse(message, err)
		return
	}

	// Send result back to master
	response := &transport.Message{
		Type:    transport.MessageTypeFileOperationResult,
		ID:      message.ID,
		ReplyTo: message.ID,
		Data:    result,
	}

	if err := a.transport.SendMessage(response); err != nil {
		a.logger.WithError(err).Error("Failed to send file operation result")
	}
}

// handleHealthCheckMessage handles health check messages
func (a *Agent) handleHealthCheckMessage(ctx context.Context, message *transport.Message) {
	var healthStatus map[string]interface{}

	if a.health != nil {
		healthStatus = a.health.GetStatus()
	} else {
		healthStatus = map[string]interface{}{
			"status": "healthy",
			"checks": map[string]interface{}{},
		}
	}

	response := &transport.Message{
		Type:    transport.MessageTypeHealthCheckResult,
		ID:      message.ID,
		ReplyTo: message.ID,
		Data:    healthStatus,
	}

	if err := a.transport.SendMessage(response); err != nil {
		a.logger.WithError(err).Error("Failed to send health check result")
	}
}

// handleConfigMessage handles configuration update messages
func (a *Agent) handleConfigMessage(ctx context.Context, message *transport.Message) {
	a.logger.Info("Received configuration update message")
	// TODO: Implement dynamic configuration updates
	
	response := &transport.Message{
		Type:    transport.MessageTypeConfigResult,
		ID:      message.ID,
		ReplyTo: message.ID,
		Data: map[string]interface{}{
			"status": "acknowledged",
		},
	}

	if err := a.transport.SendMessage(response); err != nil {
		a.logger.WithError(err).Error("Failed to send config update acknowledgment")
	}
}

// sendErrorResponse sends an error response for a message
func (a *Agent) sendErrorResponse(originalMessage *transport.Message, err error) {
	response := &transport.Message{
		Type:    transport.MessageTypeError,
		ID:      originalMessage.ID,
		ReplyTo: originalMessage.ID,
		Data: map[string]interface{}{
			"error": err.Error(),
		},
	}

	if sendErr := a.transport.SendMessage(response); sendErr != nil {
		a.logger.WithError(sendErr).Error("Failed to send error response")
	}
}

// GetConfig returns the agent configuration
func (a *Agent) GetConfig() *config.Config {
	return a.config
}

// GetExecutor returns the task executor
func (a *Agent) GetExecutor() api.ExecutorInterface {
	return a.executor
}

// GetFileOps returns the file operations manager
func (a *Agent) GetFileOps() api.FileOpsInterface {
	return a.fileops
}

// GetHealth returns the health checker
func (a *Agent) GetHealth() api.HealthInterface {
	return a.health
}

// GetMetrics returns the metrics collector
func (a *Agent) GetMetrics() api.MetricsInterface {
	return a.metrics
}

// IsRunning returns whether the agent is currently running
func (a *Agent) IsRunning() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.running
}