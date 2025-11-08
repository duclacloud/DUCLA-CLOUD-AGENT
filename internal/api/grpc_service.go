package api

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AgentService implements the gRPC agent service
type AgentService struct {
	UnimplementedAgentAPIServer
	agent  AgentInterface
	logger *logrus.Logger
}

// NewAgentService creates a new agent service
func NewAgentService(agent AgentInterface, logger *logrus.Logger) *AgentService {
	return &AgentService{
		agent:  agent,
		logger: logger,
	}
}

// GetInfo returns agent information
func (s *AgentService) GetInfo(ctx context.Context, req *InfoRequest) (*InfoResponse, error) {
	s.logger.Debug("GetInfo called")

	cfg := s.agent.GetConfig()

	return &InfoResponse{
		AgentId:      cfg.Agent.ID,
		AgentName:    cfg.Agent.Name,
		Environment:  cfg.Agent.Environment,
		Region:       cfg.Agent.Region,
		Zone:         cfg.Agent.Zone,
		Version:      "1.0.0",
		Tags:         cfg.Agent.Tags,
		Capabilities: cfg.Agent.Capabilities,
	}, nil
}

// GetStatus returns agent status
func (s *AgentService) GetStatus(ctx context.Context, req *StatusRequest) (*StatusResponse, error) {
	s.logger.Debug("GetStatus called")

	stats := s.agent.GetExecutor().GetStats()

	response := &StatusResponse{
		Running: s.agent.IsRunning(),
		Tasks: &TaskStats{
			Total:     int32(stats["total_tasks"].(int)),
			Running:   int32(stats["running_tasks"].(int)),
			Completed: int32(stats["completed_tasks"].(int)),
		},
	}

	// Add metrics if requested
	if req.IncludeMetrics && s.agent.GetMetrics() != nil {
		metrics := s.agent.GetMetrics().GetMetrics()
		// Convert metrics to protobuf format
		response.Metrics = convertMetrics(metrics)
	}

	return response, nil
}

// SubmitTask submits a task for execution
func (s *AgentService) SubmitTask(ctx context.Context, req *TaskRequest) (*TaskResponse, error) {
	s.logger.WithFields(logrus.Fields{
		"task_type": req.Type,
		"command":   req.Command,
	}).Info("SubmitTask called")

	// Convert request to task data
	taskData := map[string]interface{}{
		"type":        req.Type,
		"name":        req.Name,
		"command":     req.Command,
		"args":        req.Args,
		"env":         req.Env,
		"working_dir": req.WorkingDir,
		"timeout":     req.Timeout,
		"metadata":    req.Metadata,
	}

	// Submit task
	taskID, err := s.agent.GetExecutor().SubmitTask(taskData)
	if err != nil {
		s.logger.WithError(err).Error("Failed to submit task")
		return nil, status.Errorf(codes.Internal, "failed to submit task: %v", err)
	}

	return &TaskResponse{
		TaskId:  taskID,
		Status:  "queued",
		Message: "Task submitted successfully",
	}, nil
}

// GetTask retrieves task information
func (s *AgentService) GetTask(ctx context.Context, req *TaskDetailRequest) (*TaskDetailResponse, error) {
	s.logger.WithField("task_id", req.TaskId).Debug("GetTask called")

	task, err := s.agent.GetExecutor().GetTask(req.TaskId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "task not found: %v", err)
	}

	// Convert task to response
	// This would need proper type conversion
	return &TaskDetailResponse{
		TaskId: req.TaskId,
		Status: "running", // Would come from actual task
	}, nil
}

// CancelTask cancels a running task
func (s *AgentService) CancelTask(ctx context.Context, req *TaskDetailRequest) (*TaskResponse, error) {
	s.logger.WithField("task_id", req.TaskId).Info("CancelTask called")

	if err := s.agent.GetExecutor().CancelTask(req.TaskId); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to cancel task: %v", err)
	}

	return &TaskResponse{
		TaskId:  req.TaskId,
		Status:  "cancelled",
		Message: "Task cancelled successfully",
	}, nil
}

// ListTasks lists all tasks
func (s *AgentService) ListTasks(ctx context.Context, req *ListTasksRequest) (*ListTasksResponse, error) {
	s.logger.Debug("ListTasks called")

	var tasks []interface{}
	if req.Filter == "running" {
		tasks = s.agent.GetExecutor().ListRunningTasks()
	} else {
		tasks = s.agent.GetExecutor().ListTasks()
	}

	// Convert tasks to response format
	taskList := make([]*TaskSummary, 0, len(tasks))
	for _, task := range tasks {
		// This would need proper type conversion
		taskList = append(taskList, &TaskSummary{
			TaskId: "task-id", // Would come from actual task
			Status: "running",
		})
	}

	return &ListTasksResponse{
		Tasks: taskList,
		Total: int32(len(taskList)),
	}, nil
}

// ExecuteFileOperation executes a file operation
func (s *AgentService) ExecuteFileOperation(ctx context.Context, req *FileOperationRequest) (*FileOperationResponse, error) {
	s.logger.WithFields(logrus.Fields{
		"operation": req.Operation,
		"source":    req.SourcePath,
		"dest":      req.DestPath,
	}).Info("ExecuteFileOperation called")

	// Convert request to operation data
	opData := map[string]interface{}{
		"type":        req.Operation,
		"source_path": req.SourcePath,
		"dest_path":   req.DestPath,
		"recursive":   req.Recursive,
		"overwrite":   req.Overwrite,
		"metadata":    req.Metadata,
	}

	// Execute operation
	result, err := s.agent.GetFileOps().ExecuteOperation(ctx, opData)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to execute file operation: %v", err)
	}

	return &FileOperationResponse{
		Success: true,
		Message: "Operation completed successfully",
		Result:  convertToStringMap(result),
	}, nil
}

// GetTransferStatus retrieves transfer status
func (s *AgentService) GetTransferStatus(ctx context.Context, req *TransferStatusRequest) (*TransferStatusResponse, error) {
	s.logger.WithField("transfer_id", req.TransferId).Debug("GetTransferStatus called")

	transfer, err := s.agent.GetFileOps().GetTransfer(req.TransferId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "transfer not found: %v", err)
	}

	// Convert transfer to response
	// This would need proper type conversion
	return &TransferStatusResponse{
		TransferId: req.TransferId,
		Status:     "running",
		Progress:   0.0,
	}, nil
}

// CancelTransfer cancels a file transfer
func (s *AgentService) CancelTransfer(ctx context.Context, req *TransferStatusRequest) (*FileOperationResponse, error) {
	s.logger.WithField("transfer_id", req.TransferId).Info("CancelTransfer called")

	if err := s.agent.GetFileOps().CancelTransfer(req.TransferId); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to cancel transfer: %v", err)
	}

	return &FileOperationResponse{
		Success: true,
		Message: "Transfer cancelled successfully",
	}, nil
}

// HealthCheck performs a health check
func (s *AgentService) HealthCheck(ctx context.Context, req *HealthCheckRequest) (*HealthCheckResponse, error) {
	s.logger.Debug("HealthCheck called")

	healthy := s.agent.IsRunning()

	response := &HealthCheckResponse{
		Healthy: healthy,
		Status:  "healthy",
	}

	if !healthy {
		response.Status = "unhealthy"
	}

	// Add detailed checks if requested
	if req.Detailed && s.agent.GetHealth() != nil {
		healthStatus := s.agent.GetHealth().GetStatus()
		response.Checks = convertToStringMap(healthStatus)
	}

	return response, nil
}

// StreamLogs streams agent logs (bidirectional streaming)
func (s *AgentService) StreamLogs(stream AgentAPI_StreamLogsServer) error {
	s.logger.Info("StreamLogs called")

	// This would implement actual log streaming
	// For now, just return not implemented
	return status.Errorf(codes.Unimplemented, "log streaming not yet implemented")
}

// StreamMetrics streams agent metrics (server streaming)
func (s *AgentService) StreamMetrics(req *MetricsRequest, stream AgentAPI_StreamMetricsServer) error {
	s.logger.Info("StreamMetrics called")

	// This would implement actual metrics streaming
	// For now, just return not implemented
	return status.Errorf(codes.Unimplemented, "metrics streaming not yet implemented")
}

// Helper functions

func convertMetrics(metrics map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for key, value := range metrics {
		result[key] = fmt.Sprintf("%v", value)
	}
	return result
}

func convertToStringMap(data map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for key, value := range data {
		result[key] = fmt.Sprintf("%v", value)
	}
	return result
}