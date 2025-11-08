// Code generated for gRPC types. This is a simplified version.
// In production, these would be generated from .proto files using protoc.

package api

import (
	"context"

	"google.golang.org/grpc"
)

// AgentAPI service definition
type AgentAPIServer interface {
	GetInfo(context.Context, *InfoRequest) (*InfoResponse, error)
	GetStatus(context.Context, *StatusRequest) (*StatusResponse, error)
	SubmitTask(context.Context, *TaskRequest) (*TaskResponse, error)
	GetTask(context.Context, *TaskDetailRequest) (*TaskDetailResponse, error)
	CancelTask(context.Context, *TaskDetailRequest) (*TaskResponse, error)
	ListTasks(context.Context, *ListTasksRequest) (*ListTasksResponse, error)
	ExecuteFileOperation(context.Context, *FileOperationRequest) (*FileOperationResponse, error)
	GetTransferStatus(context.Context, *TransferStatusRequest) (*TransferStatusResponse, error)
	CancelTransfer(context.Context, *TransferStatusRequest) (*FileOperationResponse, error)
	HealthCheck(context.Context, *HealthCheckRequest) (*HealthCheckResponse, error)
	StreamLogs(AgentAPI_StreamLogsServer) error
	StreamMetrics(*MetricsRequest, AgentAPI_StreamMetricsServer) error
}

// UnimplementedAgentAPIServer provides default implementations
type UnimplementedAgentAPIServer struct{}

func (UnimplementedAgentAPIServer) GetInfo(context.Context, *InfoRequest) (*InfoResponse, error) {
	return nil, nil
}
func (UnimplementedAgentAPIServer) GetStatus(context.Context, *StatusRequest) (*StatusResponse, error) {
	return nil, nil
}
func (UnimplementedAgentAPIServer) SubmitTask(context.Context, *TaskRequest) (*TaskResponse, error) {
	return nil, nil
}
func (UnimplementedAgentAPIServer) GetTask(context.Context, *TaskDetailRequest) (*TaskDetailResponse, error) {
	return nil, nil
}
func (UnimplementedAgentAPIServer) CancelTask(context.Context, *TaskDetailRequest) (*TaskResponse, error) {
	return nil, nil
}
func (UnimplementedAgentAPIServer) ListTasks(context.Context, *ListTasksRequest) (*ListTasksResponse, error) {
	return nil, nil
}
func (UnimplementedAgentAPIServer) ExecuteFileOperation(context.Context, *FileOperationRequest) (*FileOperationResponse, error) {
	return nil, nil
}
func (UnimplementedAgentAPIServer) GetTransferStatus(context.Context, *TransferStatusRequest) (*TransferStatusResponse, error) {
	return nil, nil
}
func (UnimplementedAgentAPIServer) CancelTransfer(context.Context, *TransferStatusRequest) (*FileOperationResponse, error) {
	return nil, nil
}
func (UnimplementedAgentAPIServer) HealthCheck(context.Context, *HealthCheckRequest) (*HealthCheckResponse, error) {
	return nil, nil
}
func (UnimplementedAgentAPIServer) StreamLogs(AgentAPI_StreamLogsServer) error {
	return nil
}
func (UnimplementedAgentAPIServer) StreamMetrics(*MetricsRequest, AgentAPI_StreamMetricsServer) error {
	return nil
}

// Stream interfaces
type AgentAPI_StreamLogsServer interface {
	Send(*LogEntry) error
	Recv() (*LogRequest, error)
	grpc.ServerStream
}

type AgentAPI_StreamMetricsServer interface {
	Send(*MetricsResponse) error
	grpc.ServerStream
}

// Request/Response types

type InfoRequest struct{}

type InfoResponse struct {
	AgentId      string            `json:"agent_id"`
	AgentName    string            `json:"agent_name"`
	Environment  string            `json:"environment"`
	Region       string            `json:"region"`
	Zone         string            `json:"zone"`
	Version      string            `json:"version"`
	Tags         map[string]string `json:"tags"`
	Capabilities []string          `json:"capabilities"`
}

type StatusRequest struct {
	IncludeMetrics bool `json:"include_metrics"`
}

type StatusResponse struct {
	Running bool                `json:"running"`
	Tasks   *TaskStats          `json:"tasks"`
	Metrics map[string]string   `json:"metrics,omitempty"`
}

type TaskStats struct {
	Total     int32 `json:"total"`
	Running   int32 `json:"running"`
	Completed int32 `json:"completed"`
}

type TaskRequest struct {
	Type       string            `json:"type"`
	Name       string            `json:"name"`
	Command    string            `json:"command"`
	Args       []string          `json:"args"`
	Env        map[string]string `json:"env"`
	WorkingDir string            `json:"working_dir"`
	Timeout    int32             `json:"timeout"`
	Metadata   map[string]string `json:"metadata"`
}

type TaskResponse struct {
	TaskId  string `json:"task_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type TaskDetailRequest struct {
	TaskId string `json:"task_id"`
}

type TaskDetailResponse struct {
	TaskId     string            `json:"task_id"`
	Type       string            `json:"type"`
	Status     string            `json:"status"`
	ExitCode   int32             `json:"exit_code"`
	Output     string            `json:"output"`
	Error      string            `json:"error"`
	StartedAt  int64             `json:"started_at"`
	FinishedAt int64             `json:"finished_at"`
	Metadata   map[string]string `json:"metadata"`
}

type ListTasksRequest struct {
	Filter string `json:"filter"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

type ListTasksResponse struct {
	Tasks []*TaskSummary `json:"tasks"`
	Total int32          `json:"total"`
}

type TaskSummary struct {
	TaskId    string `json:"task_id"`
	Type      string `json:"type"`
	Status    string `json:"status"`
	StartedAt int64  `json:"started_at"`
}

type FileOperationRequest struct {
	Operation  string            `json:"operation"`
	SourcePath string            `json:"source_path"`
	DestPath   string            `json:"dest_path"`
	Recursive  bool              `json:"recursive"`
	Overwrite  bool              `json:"overwrite"`
	Metadata   map[string]string `json:"metadata"`
}

type FileOperationResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Result  map[string]string `json:"result"`
}

type TransferStatusRequest struct {
	TransferId string `json:"transfer_id"`
}

type TransferStatusResponse struct {
	TransferId  string  `json:"transfer_id"`
	Status      string  `json:"status"`
	Progress    float32 `json:"progress"`
	Transferred int64   `json:"transferred"`
	Total       int64   `json:"total"`
	Error       string  `json:"error,omitempty"`
}

type HealthCheckRequest struct {
	Detailed bool `json:"detailed"`
}

type HealthCheckResponse struct {
	Healthy bool              `json:"healthy"`
	Status  string            `json:"status"`
	Checks  map[string]string `json:"checks,omitempty"`
}

type LogRequest struct {
	Level  string `json:"level"`
	Follow bool   `json:"follow"`
}

type LogEntry struct {
	Timestamp int64  `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Fields    map[string]string `json:"fields"`
}

type MetricsRequest struct {
	Interval int32 `json:"interval"`
}

type MetricsResponse struct {
	Timestamp int64             `json:"timestamp"`
	Metrics   map[string]string `json:"metrics"`
}

// RegisterAgentAPIServer registers the service with gRPC server
func RegisterAgentAPIServer(s *grpc.Server, srv AgentAPIServer) {
	// In production, this would be generated by protoc
	// For now, we'll leave it as a placeholder
}