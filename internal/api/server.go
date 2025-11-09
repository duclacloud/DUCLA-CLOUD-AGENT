package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/duclacloud/DUCLA-CLOUD-AGENT/internal/config"
	"github.com/duclacloud/DUCLA-CLOUD-AGENT/internal/executor"
	"github.com/duclacloud/DUCLA-CLOUD-AGENT/internal/fileops"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// Server represents the API server
type Server struct {
	config config.APIConfig
	logger *logrus.Logger
	agent  AgentInterface

	// HTTP server
	httpServer *http.Server
	httpMux    *http.ServeMux

	// gRPC server
	grpcServer *grpc.Server
	grpcLis    net.Listener

	// Lifecycle
	mu      sync.Mutex
	running bool
}

// AgentInterface defines the interface for agent operations
type AgentInterface interface {
	GetConfig() *config.Config
	GetExecutor() ExecutorInterface
	GetFileOps() FileOpsInterface
	GetHealth() HealthInterface
	GetMetrics() MetricsInterface
	IsRunning() bool
}

// ExecutorInterface defines the interface for task executor
type ExecutorInterface interface {
	SubmitTask(task *executor.Task) (string, error)
	GetTask(taskID string) (*executor.Task, error)
	GetTaskResult(taskID string) (*executor.TaskResult, error)
	CancelTask(taskID string) error
	ListTasks() []*executor.Task
	ListRunningTasks() []*executor.Task
	GetStats() map[string]interface{}
}

// FileOpsInterface defines the interface for file operations
type FileOpsInterface interface {
	ExecuteOperation(ctx context.Context, op *fileops.Operation) (map[string]interface{}, error)
	GetTransfer(transferID string) (*fileops.Transfer, error)
	CancelTransfer(transferID string) error
	CalculateChecksum(path string, algorithm string) (string, error)
}

// HealthInterface defines the interface for health checks
type HealthInterface interface {
	GetStatus() map[string]interface{}
}

// MetricsInterface defines the interface for metrics
type MetricsInterface interface {
	GetMetrics() map[string]interface{}
}

// New creates a new API server
func New(cfg config.APIConfig, agent AgentInterface, logger *logrus.Logger) (*Server, error) {
	server := &Server{
		config: cfg,
		logger: logger,
		agent:  agent,
	}

	return server, nil
}

// Start starts the API server
func (s *Server) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return fmt.Errorf("API server is already running")
	}

	s.logger.Info("Starting API server")

	// Start HTTP server if enabled
	if s.config.HTTP.Enabled {
		if err := s.startHTTPServer(); err != nil {
			return fmt.Errorf("failed to start HTTP server: %w", err)
		}
	}

	// Start gRPC server if enabled
	if s.config.GRPC.Enabled {
		if err := s.startGRPCServer(); err != nil {
			// Stop HTTP server if gRPC fails
			if s.httpServer != nil {
				s.httpServer.Close()
			}
			return fmt.Errorf("failed to start gRPC server: %w", err)
		}
	}

	s.running = true
	s.logger.Info("API server started successfully")

	return nil
}

// Stop stops the API server
func (s *Server) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return nil
	}

	s.logger.Info("Stopping API server")

	// Stop HTTP server
	if s.httpServer != nil {
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
			s.logger.WithError(err).Error("Error shutting down HTTP server")
		}
	}

	// Stop gRPC server
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}

	s.running = false
	s.logger.Info("API server stopped")

	return nil
}

// Name returns the service name
func (s *Server) Name() string {
	return "api"
}

// startHTTPServer starts the HTTP server
func (s *Server) startHTTPServer() error {
	s.logger.WithFields(logrus.Fields{
		"address": s.config.HTTP.Address,
		"port":    s.config.HTTP.Port,
	}).Info("Starting HTTP server")

	// Create HTTP mux
	s.httpMux = http.NewServeMux()

	// Register HTTP handlers
	s.registerHTTPHandlers()

	// Create HTTP server
	addr := fmt.Sprintf("%s:%d", s.config.HTTP.Address, s.config.HTTP.Port)
	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      s.httpMux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		var err error
		if s.config.HTTP.TLS.Enabled {
			err = s.httpServer.ListenAndServeTLS(
				s.config.HTTP.TLS.CertFile,
				s.config.HTTP.TLS.KeyFile,
			)
		} else {
			err = s.httpServer.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			s.logger.WithError(err).Error("HTTP server error")
		}
	}()

	s.logger.WithField("address", addr).Info("HTTP server started")
	return nil
}

// startGRPCServer starts the gRPC server
func (s *Server) startGRPCServer() error {
	s.logger.WithFields(logrus.Fields{
		"address": s.config.GRPC.Address,
		"port":    s.config.GRPC.Port,
	}).Info("Starting gRPC server")

	// Create listener
	addr := fmt.Sprintf("%s:%d", s.config.GRPC.Address, s.config.GRPC.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	s.grpcLis = lis

	// Create gRPC server with options
	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(10 * 1024 * 1024), // 10MB
		grpc.MaxSendMsgSize(10 * 1024 * 1024), // 10MB
	}

	// Add TLS if enabled
	if s.config.GRPC.TLS.Enabled {
		// TLS credentials would be added here
		s.logger.Info("gRPC TLS enabled")
	}

	s.grpcServer = grpc.NewServer(opts...)

	// Register gRPC services
	s.registerGRPCServices()

	// Start server in goroutine
	go func() {
		if err := s.grpcServer.Serve(lis); err != nil {
			s.logger.WithError(err).Error("gRPC server error")
		}
	}()

	s.logger.WithField("address", addr).Info("gRPC server started")
	return nil
}

// registerHTTPHandlers registers HTTP handlers
func (s *Server) registerHTTPHandlers() {
	// Health endpoints
	s.httpMux.HandleFunc("/health", s.handleHealth)
	s.httpMux.HandleFunc("/health/live", s.handleLiveness)
	s.httpMux.HandleFunc("/health/ready", s.handleReadiness)

	// Agent info
	s.httpMux.HandleFunc("/api/v1/info", s.handleInfo)
	s.httpMux.HandleFunc("/api/v1/status", s.handleStatus)

	// Task endpoints
	s.httpMux.HandleFunc("/api/v1/tasks", s.handleTasks)
	s.httpMux.HandleFunc("/api/v1/tasks/submit", s.handleTaskSubmit)
	s.httpMux.HandleFunc("/api/v1/tasks/", s.handleTaskDetail)

	// File operation endpoints
	s.httpMux.HandleFunc("/api/v1/files", s.handleFiles)
	s.httpMux.HandleFunc("/api/v1/files/upload", s.handleFileUpload)
	s.httpMux.HandleFunc("/api/v1/files/download", s.handleFileDownload)
	s.httpMux.HandleFunc("/api/v1/files/transfer/", s.handleTransferStatus)

	// Metrics endpoint
	s.httpMux.HandleFunc("/api/v1/metrics", s.handleMetrics)

	s.logger.Info("HTTP handlers registered")
}

// registerGRPCServices registers gRPC services
func (s *Server) registerGRPCServices() {
	// Register agent service
	agentService := NewAgentService(s.agent, s.logger)
	RegisterAgentAPIServer(s.grpcServer, agentService)

	s.logger.Info("gRPC services registered")
}

// IsRunning returns whether the server is running
func (s *Server) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}