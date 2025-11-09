package metrics

import (
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// getCPUMetrics returns CPU usage and load average
func getCPUMetrics() (usage float64, loadAvg float64) {
	// Get load average
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		data, err := os.ReadFile("/proc/loadavg")
		if err == nil {
			fields := strings.Fields(string(data))
			if len(fields) > 0 {
				loadAvg, _ = strconv.ParseFloat(fields[0], 64)
			}
		} else {
			// Fallback to uptime command
			cmd := exec.Command("uptime")
			output, err := cmd.Output()
			if err == nil {
				parts := strings.Split(string(output), "load average:")
				if len(parts) == 2 {
					loads := strings.Split(strings.TrimSpace(parts[1]), ",")
					if len(loads) > 0 {
						loadAvg, _ = strconv.ParseFloat(strings.TrimSpace(loads[0]), 64)
					}
				}
			}
		}
	}

	// Calculate CPU usage from load average
	numCPU := runtime.NumCPU()
	if numCPU > 0 {
		usage = (loadAvg / float64(numCPU)) * 100
	}

	return usage, loadAvg
}

// getMemoryMetrics returns total and used memory
func getMemoryMetrics() (total, used uint64) {
	switch runtime.GOOS {
	case "linux":
		return getLinuxMemoryMetrics()
	case "darwin":
		return getDarwinMemoryMetrics()
	case "windows":
		return getWindowsMemoryMetrics()
	default:
		return 0, 0
	}
}

// getLinuxMemoryMetrics gets memory metrics on Linux
func getLinuxMemoryMetrics() (total, used uint64) {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0, 0
	}

	var available uint64
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		value, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			continue
		}

		// Convert from KB to bytes
		value *= 1024

		switch fields[0] {
		case "MemTotal:":
			total = value
		case "MemAvailable:":
			available = value
		}
	}

	if total > available {
		used = total - available
	}

	return total, used
}

// getDarwinMemoryMetrics gets memory metrics on macOS
func getDarwinMemoryMetrics() (total, used uint64) {
	// Get total memory
	cmd := exec.Command("sysctl", "-n", "hw.memsize")
	output, err := cmd.Output()
	if err != nil {
		return 0, 0
	}

	total, _ = strconv.ParseUint(strings.TrimSpace(string(output)), 10, 64)

	// Get memory usage from vm_stat
	cmd = exec.Command("vm_stat")
	output, err = cmd.Output()
	if err != nil {
		return total, 0
	}

	lines := strings.Split(string(output), "\n")
	var activePages, inactivePages, wiredPages uint64

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		value, _ := strconv.ParseUint(strings.TrimSuffix(fields[2], "."), 10, 64)

		if strings.Contains(line, "Pages active:") {
			activePages = value
		} else if strings.Contains(line, "Pages inactive:") {
			inactivePages = value
		} else if strings.Contains(line, "Pages wired down:") {
			wiredPages = value
		}
	}

	// Page size is typically 4096 bytes
	pageSize := uint64(4096)
	used = (activePages + inactivePages + wiredPages) * pageSize

	return total, used
}

// getWindowsMemoryMetrics gets memory metrics on Windows
func getWindowsMemoryMetrics() (total, used uint64) {
	// Placeholder for Windows implementation
	// Would use Windows API calls
	return 8 * 1024 * 1024 * 1024, 4 * 1024 * 1024 * 1024
}

// getDiskMetrics returns total and used disk space
func getDiskMetrics() (total, used uint64) {
	var stat syscall.Statfs_t

	// Check root filesystem
	if err := syscall.Statfs("/", &stat); err != nil {
		return 0, 0
	}

	// Calculate disk usage
	total = stat.Blocks * uint64(stat.Bsize)
	free := stat.Bfree * uint64(stat.Bsize)
	used = total - free

	return total, used
}

// TasksCollector collects task execution metrics
type TasksCollector struct {
	executor interface {
		GetStats() map[string]interface{}
	}
}

// NewTasksCollector creates a new tasks collector
func NewTasksCollector(executor interface{ GetStats() map[string]interface{} }) *TasksCollector {
	return &TasksCollector{executor: executor}
}

// Name returns the collector name
func (c *TasksCollector) Name() string {
	return "tasks"
}

// Collect collects task metrics
func (c *TasksCollector) Collect() (map[string]*Metric, error) {
	metrics := make(map[string]*Metric)

	if c.executor == nil {
		return metrics, nil
	}

	stats := c.executor.GetStats()

	// Total tasks
	if total, ok := stats["total_tasks"].(int); ok {
		metrics["agent_tasks_total"] = &Metric{
			Name:      "agent_tasks_total",
			Type:      MetricTypeGauge,
			Value:     float64(total),
			Help:      "Total number of tasks",
			Timestamp: time.Now(),
		}
	}

	// Running tasks
	if running, ok := stats["running_tasks"].(int); ok {
		metrics["agent_tasks_running"] = &Metric{
			Name:      "agent_tasks_running",
			Type:      MetricTypeGauge,
			Value:     float64(running),
			Help:      "Number of running tasks",
			Timestamp: time.Now(),
		}
	}

	// Completed tasks
	if completed, ok := stats["completed_tasks"].(int); ok {
		metrics["agent_tasks_completed_total"] = &Metric{
			Name:      "agent_tasks_completed_total",
			Type:      MetricTypeCounter,
			Value:     float64(completed),
			Help:      "Total number of completed tasks",
			Timestamp: time.Now(),
		}
	}

	// Queue size
	if queueSize, ok := stats["queue_size"].(int); ok {
		metrics["agent_tasks_queue_size"] = &Metric{
			Name:      "agent_tasks_queue_size",
			Type:      MetricTypeGauge,
			Value:     float64(queueSize),
			Help:      "Number of tasks in queue",
			Timestamp: time.Now(),
		}
	}

	return metrics, nil
}

// FilesCollector collects file operation metrics
type FilesCollector struct {
	fileops interface {
		GetStats() map[string]interface{}
	}
}

// NewFilesCollector creates a new files collector
func NewFilesCollector(fileops interface{ GetStats() map[string]interface{} }) *FilesCollector {
	return &FilesCollector{fileops: fileops}
}

// Name returns the collector name
func (c *FilesCollector) Name() string {
	return "files"
}

// Collect collects file operation metrics
func (c *FilesCollector) Collect() (map[string]*Metric, error) {
	metrics := make(map[string]*Metric)

	if c.fileops == nil {
		return metrics, nil
	}

	// This would collect file operation statistics
	// Placeholder implementation
	metrics["agent_file_operations_total"] = &Metric{
		Name:      "agent_file_operations_total",
		Type:      MetricTypeCounter,
		Value:     0,
		Help:      "Total number of file operations",
		Timestamp: time.Now(),
	}

	metrics["agent_file_transfers_active"] = &Metric{
		Name:      "agent_file_transfers_active",
		Type:      MetricTypeGauge,
		Value:     0,
		Help:      "Number of active file transfers",
		Timestamp: time.Now(),
	}

	metrics["agent_file_bytes_transferred_total"] = &Metric{
		Name:      "agent_file_bytes_transferred_total",
		Type:      MetricTypeCounter,
		Value:     0,
		Help:      "Total bytes transferred",
		Timestamp: time.Now(),
	}

	return metrics, nil
}

// NetworkCollector collects network metrics
type NetworkCollector struct{}

// NewNetworkCollector creates a new network collector
func NewNetworkCollector() *NetworkCollector {
	return &NetworkCollector{}
}

// Name returns the collector name
func (c *NetworkCollector) Name() string {
	return "network"
}

// Collect collects network metrics
func (c *NetworkCollector) Collect() (map[string]*Metric, error) {
	metrics := make(map[string]*Metric)

	// Network metrics would be collected here
	// This is a placeholder implementation
	metrics["agent_network_bytes_sent_total"] = &Metric{
		Name:      "agent_network_bytes_sent_total",
		Type:      MetricTypeCounter,
		Value:     0,
		Help:      "Total bytes sent over network",
		Timestamp: time.Now(),
	}

	metrics["agent_network_bytes_received_total"] = &Metric{
		Name:      "agent_network_bytes_received_total",
		Type:      MetricTypeCounter,
		Value:     0,
		Help:      "Total bytes received over network",
		Timestamp: time.Now(),
	}

	metrics["agent_network_connections_active"] = &Metric{
		Name:      "agent_network_connections_active",
		Type:      MetricTypeGauge,
		Value:     0,
		Help:      "Number of active network connections",
		Timestamp: time.Now(),
	}

	return metrics, nil
}