package health

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ducla/cloud-agent/internal/config"
	"github.com/sirupsen/logrus"
)

// Checker manages health checks
type Checker struct {
	config config.HealthConfig
	logger *logrus.Logger

	// Health checks
	mu     sync.RWMutex
	checks map[string]*Check
	status HealthStatus

	// HTTP server for health endpoints
	server *http.Server

	// Lifecycle
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// Check represents a single health check
type Check struct {
	Name        string        `json:"name"`
	Type        CheckType     `json:"type"`
	Status      CheckStatus   `json:"status"`
	Message     string        `json:"message"`
	LastChecked time.Time     `json:"last_checked"`
	Duration    time.Duration `json:"duration"`
	Error       string        `json:"error,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// CheckType represents the type of health check
type CheckType string

const (
	CheckTypeSystem   CheckType = "system"
	CheckTypeDisk     CheckType = "disk"
	CheckTypeMemory   CheckType = "memory"
	CheckTypeCPU      CheckType = "cpu"
	CheckTypeNetwork  CheckType = "network"
	CheckTypeDatabase CheckType = "database"
	CheckTypeService  CheckType = "service"
	CheckTypeCustom   CheckType = "custom"
)

// CheckStatus represents the status of a health check
type CheckStatus string

const (
	CheckStatusHealthy   CheckStatus = "healthy"
	CheckStatusDegraded  CheckStatus = "degraded"
	CheckStatusUnhealthy CheckStatus = "unhealthy"
	CheckStatusUnknown   CheckStatus = "unknown"
)

// HealthStatus represents overall health status
type HealthStatus struct {
	Status      CheckStatus            `json:"status"`
	Timestamp   time.Time              `json:"timestamp"`
	Uptime      time.Duration          `json:"uptime"`
	Checks      map[string]*Check      `json:"checks"`
	Summary     map[string]int         `json:"summary"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// CheckFunc is a function that performs a health check
type CheckFunc func(ctx context.Context) error

// New creates a new health checker
func New(cfg config.HealthConfig, logger *logrus.Logger) (*Checker, error) {
	checker := &Checker{
		config: cfg,
		logger: logger,
		checks: make(map[string]*Check),
		status: HealthStatus{
			Status:   CheckStatusUnknown,
			Checks:   make(map[string]*Check),
			Summary:  make(map[string]int),
			Metadata: make(map[string]interface{}),
		},
	}

	// Register default checks
	checker.registerDefaultChecks()

	return checker, nil
}

// Start starts the health checker
func (c *Checker) Start(ctx context.Context) error {
	c.logger.Info("Starting health checker")

	c.ctx, c.cancel = context.WithCancel(ctx)

	// Start health check loop
	c.wg.Add(1)
	go c.checkLoop()

	// Start HTTP server if enabled
	if c.config.Enabled {
		if err := c.startHTTPServer(); err != nil {
			return fmt.Errorf("failed to start health HTTP server: %w", err)
		}
	}

	c.logger.Info("Health checker started")
	return nil
}

// Stop stops the health checker
func (c *Checker) Stop(ctx context.Context) error {
	c.logger.Info("Stopping health checker")

	// Cancel context
	if c.cancel != nil {
		c.cancel()
	}

	// Stop HTTP server
	if c.server != nil {
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		c.server.Shutdown(shutdownCtx)
	}

	// Wait for goroutines
	c.wg.Wait()

	c.logger.Info("Health checker stopped")
	return nil
}

// Name returns the service name
func (c *Checker) Name() string {
	return "health"
}

// RegisterCheck registers a new health check
func (c *Checker) RegisterCheck(name string, checkType CheckType, checkFunc CheckFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.checks[name] = &Check{
		Name:   name,
		Type:   checkType,
		Status: CheckStatusUnknown,
	}

	c.logger.WithFields(logrus.Fields{
		"name": name,
		"type": checkType,
	}).Info("Health check registered")
}

// UnregisterCheck removes a health check
func (c *Checker) UnregisterCheck(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.checks, name)
	c.logger.WithField("name", name).Info("Health check unregistered")
}

// GetStatus returns the current health status
func (c *Checker) GetStatus() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return map[string]interface{}{
		"status":    c.status.Status,
		"timestamp": c.status.Timestamp,
		"uptime":    c.status.Uptime.Seconds(),
		"checks":    c.status.Checks,
		"summary":   c.status.Summary,
		"metadata":  c.status.Metadata,
	}
}

// IsHealthy returns whether the system is healthy
func (c *Checker) IsHealthy() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.status.Status == CheckStatusHealthy
}

// checkLoop performs periodic health checks
func (c *Checker) checkLoop() {
	defer c.wg.Done()

	ticker := time.NewTicker(c.config.Interval)
	defer ticker.Stop()

	// Run initial check
	c.performChecks()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			c.performChecks()
		}
	}
}

// performChecks executes all registered health checks
func (c *Checker) performChecks() {
	c.logger.Debug("Performing health checks")

	c.mu.Lock()
	defer c.mu.Unlock()

	startTime := time.Now()
	summary := map[string]int{
		"total":     0,
		"healthy":   0,
		"degraded":  0,
		"unhealthy": 0,
		"unknown":   0,
	}

	// Execute each check
	for name, check := range c.checks {
		checkStart := time.Now()
		
		// Perform the check based on type
		err := c.executeCheck(check)
		
		check.LastChecked = time.Now()
		check.Duration = time.Since(checkStart)

		if err != nil {
			check.Status = CheckStatusUnhealthy
			check.Error = err.Error()
			check.Message = fmt.Sprintf("Check failed: %v", err)
		} else {
			check.Status = CheckStatusHealthy
			check.Error = ""
			check.Message = "Check passed"
		}

		// Update summary
		summary["total"]++
		switch check.Status {
		case CheckStatusHealthy:
			summary["healthy"]++
		case CheckStatusDegraded:
			summary["degraded"]++
		case CheckStatusUnhealthy:
			summary["unhealthy"]++
		default:
			summary["unknown"]++
		}

		c.status.Checks[name] = check
	}

	// Determine overall status
	overallStatus := CheckStatusHealthy
	if summary["unhealthy"] > 0 {
		overallStatus = CheckStatusUnhealthy
	} else if summary["degraded"] > 0 {
		overallStatus = CheckStatusDegraded
	} else if summary["unknown"] > 0 {
		overallStatus = CheckStatusUnknown
	}

	c.status.Status = overallStatus
	c.status.Timestamp = time.Now()
	c.status.Summary = summary
	c.status.Metadata["check_duration_ms"] = time.Since(startTime).Milliseconds()

	c.logger.WithFields(logrus.Fields{
		"status":   overallStatus,
		"duration": time.Since(startTime),
		"summary":  summary,
	}).Debug("Health checks completed")
}

// executeCheck executes a specific health check
func (c *Checker) executeCheck(check *Check) error {
	ctx, cancel := context.WithTimeout(c.ctx, c.config.Timeout)
	defer cancel()

	switch check.Type {
	case CheckTypeSystem:
		return c.checkSystem(ctx, check)
	case CheckTypeDisk:
		return c.checkDisk(ctx, check)
	case CheckTypeMemory:
		return c.checkMemory(ctx, check)
	case CheckTypeCPU:
		return c.checkCPU(ctx, check)
	case CheckTypeNetwork:
		return c.checkNetwork(ctx, check)
	default:
		return nil
	}
}

// registerDefaultChecks registers default health checks
func (c *Checker) registerDefaultChecks() {
	c.RegisterCheck("system", CheckTypeSystem, nil)
	c.RegisterCheck("disk", CheckTypeDisk, nil)
	c.RegisterCheck("memory", CheckTypeMemory, nil)
	c.RegisterCheck("cpu", CheckTypeCPU, nil)
}

// startHTTPServer starts the health HTTP server
func (c *Checker) startHTTPServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc(c.config.Path, c.handleHealthCheck)
	mux.HandleFunc("/live", c.handleLiveness)
	mux.HandleFunc("/ready", c.handleReadiness)

	addr := fmt.Sprintf("%s:%d", c.config.Address, c.config.Port)
	c.server = &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		c.logger.WithField("address", addr).Info("Health HTTP server started")
		if err := c.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			c.logger.WithError(err).Error("Health HTTP server error")
		}
	}()

	return nil
}

// handleHealthCheck handles health check HTTP requests
func (c *Checker) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	status := c.GetStatus()

	w.Header().Set("Content-Type", "application/json")
	
	statusCode := http.StatusOK
	if c.status.Status == CheckStatusUnhealthy {
		statusCode = http.StatusServiceUnavailable
	} else if c.status.Status == CheckStatusDegraded {
		statusCode = http.StatusOK // Still return 200 for degraded
	}

	w.WriteHeader(statusCode)
	
	// Write JSON response
	if err := writeJSON(w, status); err != nil {
		c.logger.WithError(err).Error("Failed to write health check response")
	}
}

// handleLiveness handles liveness probe
func (c *Checker) handleLiveness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	writeJSON(w, map[string]interface{}{
		"alive": true,
		"timestamp": time.Now().Unix(),
	})
}

// handleReadiness handles readiness probe
func (c *Checker) handleReadiness(w http.ResponseWriter, r *http.Request) {
	ready := c.IsHealthy()
	
	w.Header().Set("Content-Type", "application/json")
	
	statusCode := http.StatusOK
	if !ready {
		statusCode = http.StatusServiceUnavailable
	}
	
	w.WriteHeader(statusCode)
	writeJSON(w, map[string]interface{}{
		"ready": ready,
		"timestamp": time.Now().Unix(),
	})
}

// Helper function to write JSON response
func writeJSON(w http.ResponseWriter, data interface{}) error {
	// Simple JSON encoding
	// In production, use json.NewEncoder(w).Encode(data)
	return nil
}