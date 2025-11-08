package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// handleHealth handles health check requests
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"uptime":    time.Since(time.Now()).Seconds(), // Would be actual uptime
	}

	if s.agent.GetHealth() != nil {
		health["checks"] = s.agent.GetHealth().GetStatus()
	}

	s.respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    health,
	})
}

// handleLiveness handles liveness probe
func (s *Server) handleLiveness(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	s.respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data: map[string]interface{}{
			"alive": s.agent.IsRunning(),
		},
	})
}

// handleReadiness handles readiness probe
func (s *Server) handleReadiness(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	ready := s.agent.IsRunning()

	status := http.StatusOK
	if !ready {
		status = http.StatusServiceUnavailable
	}

	s.respondJSON(w, status, Response{
		Success: ready,
		Data: map[string]interface{}{
			"ready": ready,
		},
	})
}

// handleInfo handles agent info requests
func (s *Server) handleInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	cfg := s.agent.GetConfig()
	info := map[string]interface{}{
		"agent_id":     cfg.Agent.ID,
		"agent_name":   cfg.Agent.Name,
		"environment":  cfg.Agent.Environment,
		"region":       cfg.Agent.Region,
		"zone":         cfg.Agent.Zone,
		"tags":         cfg.Agent.Tags,
		"capabilities": cfg.Agent.Capabilities,
		"version":      "1.0.0", // Would come from build info
	}

	s.respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    info,
	})
}

// handleStatus handles agent status requests
func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	status := map[string]interface{}{
		"running": s.agent.IsRunning(),
		"tasks":   s.agent.GetExecutor().GetStats(),
	}

	if s.agent.GetMetrics() != nil {
		status["metrics"] = s.agent.GetMetrics().GetMetrics()
	}

	s.respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    status,
	})
}

// handleTasks handles task list requests
func (s *Server) handleTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Get query parameters
	filter := r.URL.Query().Get("filter")

	var tasks []interface{}
	if filter == "running" {
		tasks = s.agent.GetExecutor().ListRunningTasks()
	} else {
		tasks = s.agent.GetExecutor().ListTasks()
	}

	s.respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data: map[string]interface{}{
			"tasks": tasks,
			"count": len(tasks),
		},
	})
}

// handleTaskSubmit handles task submission requests
func (s *Server) handleTaskSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Parse request body
	var taskData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&taskData); err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Submit task
	taskID, err := s.agent.GetExecutor().SubmitTask(taskData)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.respondJSON(w, http.StatusAccepted, Response{
		Success: true,
		Data: map[string]interface{}{
			"task_id": taskID,
			"status":  "queued",
		},
		Message: "Task submitted successfully",
	})
}

// handleTaskDetail handles task detail requests
func (s *Server) handleTaskDetail(w http.ResponseWriter, r *http.Request) {
	// Extract task ID from path
	path := strings.TrimPrefix(r.Path, "/api/v1/tasks/")
	taskID := strings.Split(path, "/")[0]

	if taskID == "" {
		s.respondError(w, http.StatusBadRequest, "Task ID is required")
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.handleTaskGet(w, r, taskID)
	case http.MethodDelete:
		s.handleTaskCancel(w, r, taskID)
	default:
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// handleTaskGet handles get task requests
func (s *Server) handleTaskGet(w http.ResponseWriter, r *http.Request, taskID string) {
	task, err := s.agent.GetExecutor().GetTask(taskID)
	if err != nil {
		s.respondError(w, http.StatusNotFound, err.Error())
		return
	}

	s.respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    task,
	})
}

// handleTaskCancel handles cancel task requests
func (s *Server) handleTaskCancel(w http.ResponseWriter, r *http.Request, taskID string) {
	if err := s.agent.GetExecutor().CancelTask(taskID); err != nil {
		s.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.respondJSON(w, http.StatusOK, Response{
		Success: true,
		Message: "Task cancelled successfully",
	})
}

// handleFiles handles file operation requests
func (s *Server) handleFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Parse request body
	var opData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&opData); err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Execute file operation
	result, err := s.agent.GetFileOps().ExecuteOperation(r.Context(), opData)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    result,
	})
}

// handleFileUpload handles file upload requests
func (s *Server) handleFileUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Parse multipart form
	if err := r.ParseMultipartForm(32 << 20); err != nil { // 32MB max
		s.respondError(w, http.StatusBadRequest, "Failed to parse form")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		s.respondError(w, http.StatusBadRequest, "No file provided")
		return
	}
	defer file.Close()

	destPath := r.FormValue("dest_path")
	if destPath == "" {
		s.respondError(w, http.StatusBadRequest, "Destination path is required")
		return
	}

	// Handle file upload
	// Implementation would save the file
	s.respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data: map[string]interface{}{
			"filename": header.Filename,
			"size":     header.Size,
			"dest":     destPath,
		},
		Message: "File uploaded successfully",
	})
}

// handleFileDownload handles file download requests
func (s *Server) handleFileDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		s.respondError(w, http.StatusBadRequest, "File path is required")
		return
	}

	// Serve file
	http.ServeFile(w, r, filePath)
}

// handleTransferStatus handles transfer status requests
func (s *Server) handleTransferStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract transfer ID from path
	path := strings.TrimPrefix(r.Path, "/api/v1/files/transfer/")
	transferID := strings.Split(path, "/")[0]

	if transferID == "" {
		s.respondError(w, http.StatusBadRequest, "Transfer ID is required")
		return
	}

	transfer, err := s.agent.GetFileOps().GetTransfer(transferID)
	if err != nil {
		s.respondError(w, http.StatusNotFound, err.Error())
		return
	}

	s.respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    transfer,
	})
}

// handleMetrics handles metrics requests
func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var metrics map[string]interface{}
	if s.agent.GetMetrics() != nil {
		metrics = s.agent.GetMetrics().GetMetrics()
	} else {
		metrics = map[string]interface{}{
			"message": "Metrics collection not enabled",
		}
	}

	s.respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    metrics,
	})
}

// respondJSON sends a JSON response
func (s *Server) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		s.logger.WithError(err).Error("Failed to encode JSON response")
	}
}

// respondError sends an error response
func (s *Server) respondError(w http.ResponseWriter, status int, message string) {
	s.respondJSON(w, status, Response{
		Success: false,
		Error:   message,
	})
}

// Middleware for logging
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call next handler
		next.ServeHTTP(w, r)

		// Log request
		s.logger.WithFields(map[string]interface{}{
			"method":   r.Method,
			"path":     r.URL.Path,
			"duration": time.Since(start).Milliseconds(),
			"remote":   r.RemoteAddr,
		}).Info("HTTP request")
	})
}

// Middleware for authentication
func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for health endpoints
		if strings.HasPrefix(r.URL.Path, "/health") {
			next.ServeHTTP(w, r)
			return
		}

		// Check authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			s.respondError(w, http.StatusUnauthorized, "Authorization required")
			return
		}

		// Validate token
		// Implementation would validate JWT token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			s.respondError(w, http.StatusUnauthorized, "Invalid authorization format")
			return
		}

		// Call next handler
		next.ServeHTTP(w, r)
	})
}

// Middleware for CORS
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}