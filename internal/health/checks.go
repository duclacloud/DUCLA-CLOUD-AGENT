package health

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// checkSystem performs system health check
func (c *Checker) checkSystem(ctx context.Context, check *Check) error {
	// Check if system is responsive
	startTime := time.Now()
	
	// Basic system checks
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("failed to get hostname: %w", err)
	}

	check.Metadata = map[string]interface{}{
		"hostname":   hostname,
		"os":         runtime.GOOS,
		"arch":       runtime.GOARCH,
		"go_version": runtime.Version(),
		"num_cpu":    runtime.NumCPU(),
		"num_goroutine": runtime.NumGoroutine(),
	}

	// Check system responsiveness
	duration := time.Since(startTime)
	if duration > 1*time.Second {
		check.Status = CheckStatusDegraded
		check.Message = fmt.Sprintf("System response slow: %v", duration)
	}

	return nil
}

// checkDisk performs disk health check
func (c *Checker) checkDisk(ctx context.Context, check *Check) error {
	var stat syscall.Statfs_t
	
	// Check root filesystem
	if err := syscall.Statfs("/", &stat); err != nil {
		return fmt.Errorf("failed to get disk stats: %w", err)
	}

	// Calculate disk usage
	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bfree * uint64(stat.Bsize)
	used := total - free
	usedPercent := float64(used) / float64(total) * 100

	check.Metadata = map[string]interface{}{
		"total_bytes":   total,
		"used_bytes":    used,
		"free_bytes":    free,
		"used_percent":  usedPercent,
		"total_gb":      float64(total) / 1024 / 1024 / 1024,
		"free_gb":       float64(free) / 1024 / 1024 / 1024,
	}

	// Check thresholds
	if usedPercent > 95 {
		return fmt.Errorf("disk usage critical: %.2f%%", usedPercent)
	} else if usedPercent > 85 {
		check.Status = CheckStatusDegraded
		check.Message = fmt.Sprintf("Disk usage high: %.2f%%", usedPercent)
	}

	return nil
}

// checkMemory performs memory health check
func (c *Checker) checkMemory(ctx context.Context, check *Check) error {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Get system memory info
	totalMem, freeMem, err := getSystemMemory()
	if err != nil {
		return fmt.Errorf("failed to get system memory: %w", err)
	}

	usedMem := totalMem - freeMem
	usedPercent := float64(usedMem) / float64(totalMem) * 100

	check.Metadata = map[string]interface{}{
		"total_bytes":      totalMem,
		"used_bytes":       usedMem,
		"free_bytes":       freeMem,
		"used_percent":     usedPercent,
		"alloc_bytes":      m.Alloc,
		"sys_bytes":        m.Sys,
		"num_gc":           m.NumGC,
		"heap_alloc_bytes": m.HeapAlloc,
		"heap_sys_bytes":   m.HeapSys,
	}

	// Check thresholds
	if usedPercent > 95 {
		return fmt.Errorf("memory usage critical: %.2f%%", usedPercent)
	} else if usedPercent > 85 {
		check.Status = CheckStatusDegraded
		check.Message = fmt.Sprintf("Memory usage high: %.2f%%", usedPercent)
	}

	return nil
}

// checkCPU performs CPU health check
func (c *Checker) checkCPU(ctx context.Context, check *Check) error {
	// Get CPU usage
	cpuUsage, loadAvg, err := getCPUInfo()
	if err != nil {
		return fmt.Errorf("failed to get CPU info: %w", err)
	}

	numCPU := runtime.NumCPU()
	
	check.Metadata = map[string]interface{}{
		"num_cpu":      numCPU,
		"cpu_usage":    cpuUsage,
		"load_average": loadAvg,
		"load_per_cpu": loadAvg / float64(numCPU),
	}

	// Check thresholds
	loadPerCPU := loadAvg / float64(numCPU)
	if loadPerCPU > 2.0 {
		return fmt.Errorf("CPU load critical: %.2f per CPU", loadPerCPU)
	} else if loadPerCPU > 1.5 {
		check.Status = CheckStatusDegraded
		check.Message = fmt.Sprintf("CPU load high: %.2f per CPU", loadPerCPU)
	}

	return nil
}

// checkNetwork performs network health check
func (c *Checker) checkNetwork(ctx context.Context, check *Check) error {
	// Check network connectivity
	// Try to resolve a hostname
	cmd := exec.CommandContext(ctx, "ping", "-c", "1", "-W", "1", "8.8.8.8")
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		check.Metadata = map[string]interface{}{
			"error": string(output),
		}
		return fmt.Errorf("network connectivity check failed: %w", err)
	}

	check.Metadata = map[string]interface{}{
		"connectivity": "ok",
	}

	return nil
}

// getSystemMemory returns total and free system memory
func getSystemMemory() (total, free uint64, err error) {
	switch runtime.GOOS {
	case "linux":
		return getLinuxMemory()
	case "darwin":
		return getDarwinMemory()
	case "windows":
		return getWindowsMemory()
	default:
		return 0, 0, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

// getLinuxMemory gets memory info on Linux
func getLinuxMemory() (total, free uint64, err error) {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0, 0, err
	}

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
			free = value
		}
	}

	return total, free, nil
}

// getDarwinMemory gets memory info on macOS
func getDarwinMemory() (total, free uint64, err error) {
	// Get total memory
	cmd := exec.Command("sysctl", "-n", "hw.memsize")
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}

	total, err = strconv.ParseUint(strings.TrimSpace(string(output)), 10, 64)
	if err != nil {
		return 0, 0, err
	}

	// Get free memory (simplified)
	cmd = exec.Command("vm_stat")
	output, err = cmd.Output()
	if err != nil {
		return 0, 0, err
	}

	// Parse vm_stat output
	lines := strings.Split(string(output), "\n")
	var freePages, inactivePages uint64
	
	for _, line := range lines {
		if strings.Contains(line, "Pages free:") {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				freePages, _ = strconv.ParseUint(strings.TrimSuffix(fields[2], "."), 10, 64)
			}
		} else if strings.Contains(line, "Pages inactive:") {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				inactivePages, _ = strconv.ParseUint(strings.TrimSuffix(fields[2], "."), 10, 64)
			}
		}
	}

	// Page size is typically 4096 bytes
	pageSize := uint64(4096)
	free = (freePages + inactivePages) * pageSize

	return total, free, nil
}

// getWindowsMemory gets memory info on Windows
func getWindowsMemory() (total, free uint64, err error) {
	// This would use Windows API calls
	// For now, return placeholder values
	return 8 * 1024 * 1024 * 1024, 4 * 1024 * 1024 * 1024, nil
}

// getCPUInfo returns CPU usage and load average
func getCPUInfo() (usage float64, loadAvg float64, err error) {
	switch runtime.GOOS {
	case "linux", "darwin":
		return getUnixCPUInfo()
	case "windows":
		return getWindowsCPUInfo()
	default:
		return 0, 0, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

// getUnixCPUInfo gets CPU info on Unix-like systems
func getUnixCPUInfo() (usage float64, loadAvg float64, err error) {
	// Get load average
	data, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		// Try uptime command as fallback
		cmd := exec.Command("uptime")
		output, err := cmd.Output()
		if err != nil {
			return 0, 0, err
		}

		// Parse uptime output
		parts := strings.Split(string(output), "load average:")
		if len(parts) == 2 {
			loads := strings.Split(strings.TrimSpace(parts[1]), ",")
			if len(loads) > 0 {
				loadAvg, _ = strconv.ParseFloat(strings.TrimSpace(loads[0]), 64)
			}
		}
	} else {
		fields := strings.Fields(string(data))
		if len(fields) > 0 {
			loadAvg, _ = strconv.ParseFloat(fields[0], 64)
		}
	}

	// CPU usage would require sampling over time
	// For now, use load average as proxy
	usage = loadAvg / float64(runtime.NumCPU()) * 100

	return usage, loadAvg, nil
}

// getWindowsCPUInfo gets CPU info on Windows
func getWindowsCPUInfo() (usage float64, loadAvg float64, err error) {
	// This would use Windows API calls
	// For now, return placeholder values
	return 25.0, 1.0, nil
}

// CheckDatabase performs database health check
func CheckDatabase(ctx context.Context, dsn string) error {
	// This would check database connectivity
	// Implementation depends on database type
	return nil
}

// CheckService performs service health check
func CheckService(ctx context.Context, url string) error {
	// This would check if a service is reachable
	// Implementation would use HTTP client
	return nil
}

// CheckCustom performs custom health check
func CheckCustom(ctx context.Context, checkFunc CheckFunc) error {
	if checkFunc == nil {
		return nil
	}
	return checkFunc(ctx)
}