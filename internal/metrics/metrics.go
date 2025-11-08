package metrics

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/ducla/cloud-agent/internal/config"
	"github.com/sirupsen/logrus"
)

// Collector manages metrics collection
type Collector struct {
	config config.MetricsConfig
	logger *logrus.Logger

	// Metrics storage
	mu      sync.RWMutex
	metrics map[string]*Metric

	// Collectors
	collectors map[string]MetricCollector

	// HTTP server for Prometheus endpoint
	server *http.Server

	// Lifecycle
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// Metric represents a single metric
type Metric struct {
	Name        string            `json:"name"`
	Type        MetricType        `json:"type"`
	Value       float64           `json:"value"`
	Labels      map[string]string `json:"labels"`
	Help        string            `json:"help"`
	Timestamp   time.Time         `json:"timestamp"`
}

// MetricType represents the type of metric
type MetricType string

const (
	MetricTypeCounter   MetricType = "counter"
	MetricTypeGauge     MetricType = "gauge"
	MetricTypeHistogram MetricType = "histogram"
	MetricTypeSummary   MetricType = "summary"
)

// MetricCollector interface for different metric collectors
type MetricCollector interface {
	Name() string
	Collect() (map[string]*Metric, error)
}

// New creates a new metrics collector
func New(cfg config.MetricsConfig, logger *logrus.Logger) (*Collector, error) {
	collector := &Collector{
		config:     cfg,
		logger:     logger,
		metrics:    make(map[string]*Metric),
		collectors: make(map[string]MetricCollector),
	}

	// Register default collectors
	if err := collector.registerDefaultCollectors(); err != nil {
		return nil, fmt.Errorf("failed to register collectors: %w", err)
	}

	return collector, nil
}

// Start starts the metrics collector
func (c *Collector) Start(ctx context.Context) error {
	c.logger.Info("Starting metrics collector")

	c.ctx, c.cancel = context.WithCancel(ctx)

	// Start collection loop
	c.wg.Add(1)
	go c.collectionLoop()

	// Start HTTP server for Prometheus endpoint
	if c.config.Enabled {
		if err := c.startHTTPServer(); err != nil {
			return fmt.Errorf("failed to start metrics HTTP server: %w", err)
		}
	}

	c.logger.Info("Metrics collector started")
	return nil
}

// Stop stops the metrics collector
func (c *Collector) Stop(ctx context.Context) error {
	c.logger.Info("Stopping metrics collector")

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

	c.logger.Info("Metrics collector stopped")
	return nil
}

// Name returns the service name
func (c *Collector) Name() string {
	return "metrics"
}

// RegisterCollector registers a new metric collector
func (c *Collector) RegisterCollector(collector MetricCollector) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.collectors[collector.Name()] = collector
	c.logger.WithField("name", collector.Name()).Info("Metric collector registered")
}

// UnregisterCollector removes a metric collector
func (c *Collector) UnregisterCollector(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.collectors, name)
	c.logger.WithField("name", name).Info("Metric collector unregistered")
}

// RecordMetric records a single metric
func (c *Collector) RecordMetric(name string, metricType MetricType, value float64, labels map[string]string, help string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.metrics[name] = &Metric{
		Name:      name,
		Type:      metricType,
		Value:     value,
		Labels:    labels,
		Help:      help,
		Timestamp: time.Now(),
	}
}

// IncrementCounter increments a counter metric
func (c *Collector) IncrementCounter(name string, labels map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.metricKey(name, labels)
	if metric, exists := c.metrics[key]; exists {
		metric.Value++
		metric.Timestamp = time.Now()
	} else {
		c.metrics[key] = &Metric{
			Name:      name,
			Type:      MetricTypeCounter,
			Value:     1,
			Labels:    labels,
			Timestamp: time.Now(),
		}
	}
}

// SetGauge sets a gauge metric value
func (c *Collector) SetGauge(name string, value float64, labels map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.metricKey(name, labels)
	c.metrics[key] = &Metric{
		Name:      name,
		Type:      MetricTypeGauge,
		Value:     value,
		Labels:    labels,
		Timestamp: time.Now(),
	}
}

// GetMetrics returns all current metrics
func (c *Collector) GetMetrics() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[string]interface{})
	for name, metric := range c.metrics {
		result[name] = map[string]interface{}{
			"type":      metric.Type,
			"value":     metric.Value,
			"labels":    metric.Labels,
			"timestamp": metric.Timestamp.Unix(),
		}
	}

	return result
}

// collectionLoop performs periodic metrics collection
func (c *Collector) collectionLoop() {
	defer c.wg.Done()

	ticker := time.NewTicker(c.config.Interval)
	defer ticker.Stop()

	// Run initial collection
	c.collectMetrics()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			c.collectMetrics()
		}
	}
}

// collectMetrics collects metrics from all registered collectors
func (c *Collector) collectMetrics() {
	c.logger.Debug("Collecting metrics")

	c.mu.Lock()
	collectors := make([]MetricCollector, 0, len(c.collectors))
	for _, collector := range c.collectors {
		collectors = append(collectors, collector)
	}
	c.mu.Unlock()

	// Collect from each collector
	for _, collector := range collectors {
		metrics, err := collector.Collect()
		if err != nil {
			c.logger.WithError(err).WithField("collector", collector.Name()).Error("Failed to collect metrics")
			continue
		}

		// Store collected metrics
		c.mu.Lock()
		for name, metric := range metrics {
			c.metrics[name] = metric
		}
		c.mu.Unlock()
	}

	c.logger.WithField("count", len(c.metrics)).Debug("Metrics collection completed")
}

// registerDefaultCollectors registers default metric collectors
func (c *Collector) registerDefaultCollectors() error {
	// Register collectors based on configuration
	for _, collectorName := range c.config.Collectors {
		switch collectorName {
		case "system":
			c.RegisterCollector(NewSystemCollector(c.logger))
		case "process":
			c.RegisterCollector(NewProcessCollector(c.logger))
		case "agent":
			c.RegisterCollector(NewAgentCollector(c.logger))
		case "tasks":
			// Tasks collector would be registered by executor
		case "files":
			// Files collector would be registered by fileops
		default:
			c.logger.WithField("collector", collectorName).Warn("Unknown collector type")
		}
	}

	return nil
}

// startHTTPServer starts the metrics HTTP server
func (c *Collector) startHTTPServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc(c.config.Path, c.handleMetrics)

	addr := fmt.Sprintf("%s:%d", c.config.Address, c.config.Port)
	c.server = &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		c.logger.WithField("address", addr).Info("Metrics HTTP server started")
		if err := c.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			c.logger.WithError(err).Error("Metrics HTTP server error")
		}
	}()

	return nil
}

// handleMetrics handles Prometheus metrics endpoint
func (c *Collector) handleMetrics(w http.ResponseWriter, r *http.Request) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	w.Header().Set("Content-Type", "text/plain; version=0.0.4")

	// Write metrics in Prometheus format
	for _, metric := range c.metrics {
		c.writePrometheusMetric(w, metric)
	}
}

// writePrometheusMetric writes a metric in Prometheus format
func (c *Collector) writePrometheusMetric(w http.ResponseWriter, metric *Metric) {
	// Write HELP line
	if metric.Help != "" {
		fmt.Fprintf(w, "# HELP %s %s\n", metric.Name, metric.Help)
	}

	// Write TYPE line
	fmt.Fprintf(w, "# TYPE %s %s\n", metric.Name, metric.Type)

	// Write metric line
	if len(metric.Labels) > 0 {
		labels := ""
		first := true
		for key, value := range metric.Labels {
			if !first {
				labels += ","
			}
			labels += fmt.Sprintf("%s=\"%s\"", key, value)
			first = false
		}
		fmt.Fprintf(w, "%s{%s} %f %d\n", metric.Name, labels, metric.Value, metric.Timestamp.UnixMilli())
	} else {
		fmt.Fprintf(w, "%s %f %d\n", metric.Name, metric.Value, metric.Timestamp.UnixMilli())
	}
}

// metricKey generates a unique key for a metric with labels
func (c *Collector) metricKey(name string, labels map[string]string) string {
	if len(labels) == 0 {
		return name
	}

	key := name
	for k, v := range labels {
		key += fmt.Sprintf("_%s_%s", k, v)
	}
	return key
}

// SystemCollector collects system metrics
type SystemCollector struct {
	logger *logrus.Logger
}

// NewSystemCollector creates a new system collector
func NewSystemCollector(logger *logrus.Logger) *SystemCollector {
	return &SystemCollector{logger: logger}
}

// Name returns the collector name
func (c *SystemCollector) Name() string {
	return "system"
}

// Collect collects system metrics
func (c *SystemCollector) Collect() (map[string]*Metric, error) {
	metrics := make(map[string]*Metric)

	// CPU metrics
	cpuUsage, loadAvg := getCPUMetrics()
	metrics["system_cpu_usage_percent"] = &Metric{
		Name:      "system_cpu_usage_percent",
		Type:      MetricTypeGauge,
		Value:     cpuUsage,
		Help:      "System CPU usage percentage",
		Timestamp: time.Now(),
	}
	metrics["system_load_average"] = &Metric{
		Name:      "system_load_average",
		Type:      MetricTypeGauge,
		Value:     loadAvg,
		Help:      "System load average (1 minute)",
		Timestamp: time.Now(),
	}

	// Memory metrics
	totalMem, usedMem := getMemoryMetrics()
	metrics["system_memory_total_bytes"] = &Metric{
		Name:      "system_memory_total_bytes",
		Type:      MetricTypeGauge,
		Value:     float64(totalMem),
		Help:      "Total system memory in bytes",
		Timestamp: time.Now(),
	}
	metrics["system_memory_used_bytes"] = &Metric{
		Name:      "system_memory_used_bytes",
		Type:      MetricTypeGauge,
		Value:     float64(usedMem),
		Help:      "Used system memory in bytes",
		Timestamp: time.Now(),
	}

	// Disk metrics
	totalDisk, usedDisk := getDiskMetrics()
	metrics["system_disk_total_bytes"] = &Metric{
		Name:      "system_disk_total_bytes",
		Type:      MetricTypeGauge,
		Value:     float64(totalDisk),
		Help:      "Total disk space in bytes",
		Timestamp: time.Now(),
	}
	metrics["system_disk_used_bytes"] = &Metric{
		Name:      "system_disk_used_bytes",
		Type:      MetricTypeGauge,
		Value:     float64(usedDisk),
		Help:      "Used disk space in bytes",
		Timestamp: time.Now(),
	}

	return metrics, nil
}

// ProcessCollector collects process metrics
type ProcessCollector struct {
	logger *logrus.Logger
}

// NewProcessCollector creates a new process collector
func NewProcessCollector(logger *logrus.Logger) *ProcessCollector {
	return &ProcessCollector{logger: logger}
}

// Name returns the collector name
func (c *ProcessCollector) Name() string {
	return "process"
}

// Collect collects process metrics
func (c *ProcessCollector) Collect() (map[string]*Metric, error) {
	metrics := make(map[string]*Metric)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Process memory metrics
	metrics["process_memory_alloc_bytes"] = &Metric{
		Name:      "process_memory_alloc_bytes",
		Type:      MetricTypeGauge,
		Value:     float64(m.Alloc),
		Help:      "Bytes of allocated heap objects",
		Timestamp: time.Now(),
	}
	metrics["process_memory_sys_bytes"] = &Metric{
		Name:      "process_memory_sys_bytes",
		Type:      MetricTypeGauge,
		Value:     float64(m.Sys),
		Help:      "Total bytes of memory obtained from OS",
		Timestamp: time.Now(),
	}
	metrics["process_memory_heap_alloc_bytes"] = &Metric{
		Name:      "process_memory_heap_alloc_bytes",
		Type:      MetricTypeGauge,
		Value:     float64(m.HeapAlloc),
		Help:      "Bytes of allocated heap objects",
		Timestamp: time.Now(),
	}

	// GC metrics
	metrics["process_gc_count_total"] = &Metric{
		Name:      "process_gc_count_total",
		Type:      MetricTypeCounter,
		Value:     float64(m.NumGC),
		Help:      "Total number of GC cycles",
		Timestamp: time.Now(),
	}
	metrics["process_gc_pause_seconds"] = &Metric{
		Name:      "process_gc_pause_seconds",
		Type:      MetricTypeGauge,
		Value:     float64(m.PauseTotalNs) / 1e9,
		Help:      "Total GC pause time in seconds",
		Timestamp: time.Now(),
	}

	// Goroutine metrics
	metrics["process_goroutines"] = &Metric{
		Name:      "process_goroutines",
		Type:      MetricTypeGauge,
		Value:     float64(runtime.NumGoroutine()),
		Help:      "Number of goroutines",
		Timestamp: time.Now(),
	}

	// CPU metrics
	metrics["process_cpu_cores"] = &Metric{
		Name:      "process_cpu_cores",
		Type:      MetricTypeGauge,
		Value:     float64(runtime.NumCPU()),
		Help:      "Number of CPU cores",
		Timestamp: time.Now(),
	}

	return metrics, nil
}

// AgentCollector collects agent-specific metrics
type AgentCollector struct {
	logger *logrus.Logger
}

// NewAgentCollector creates a new agent collector
func NewAgentCollector(logger *logrus.Logger) *AgentCollector {
	return &AgentCollector{logger: logger}
}

// Name returns the collector name
func (c *AgentCollector) Name() string {
	return "agent"
}

// Collect collects agent metrics
func (c *AgentCollector) Collect() (map[string]*Metric, error) {
	metrics := make(map[string]*Metric)

	// Agent uptime
	// This would be calculated from agent start time
	metrics["agent_uptime_seconds"] = &Metric{
		Name:      "agent_uptime_seconds",
		Type:      MetricTypeGauge,
		Value:     time.Since(time.Now()).Seconds(), // Placeholder
		Help:      "Agent uptime in seconds",
		Timestamp: time.Now(),
	}

	// Agent info
	metrics["agent_info"] = &Metric{
		Name: "agent_info",
		Type: MetricTypeGauge,
		Value: 1,
		Labels: map[string]string{
			"version": "1.0.0",
			"os":      runtime.GOOS,
			"arch":    runtime.GOARCH,
		},
		Help:      "Agent information",
		Timestamp: time.Now(),
	}

	return metrics, nil
}