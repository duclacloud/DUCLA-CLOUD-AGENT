package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/duclacloud/DUCLA-CLOUD-AGENT/internal/config"
	"github.com/sirupsen/logrus"
)

// CLI commands handler
func handleCLICommand(args []string, cfg *config.Config, logger *logrus.Logger) error {
	if len(args) < 1 {
		return fmt.Errorf("no command specified")
	}

	command := args[0]
	subArgs := args[1:]

	switch command {
	case "show":
		return handleShowCommand(subArgs, cfg, logger)
	case "task":
		return handleTaskCommand(subArgs, cfg, logger)
	case "file":
		return handleFileCommand(subArgs, cfg, logger)
	case "config":
		return handleConfigCommand(subArgs, cfg, logger)
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

// Handle 'show' commands
func handleShowCommand(args []string, cfg *config.Config, logger *logrus.Logger) error {
	if len(args) < 1 {
		return fmt.Errorf("show command requires a subcommand")
	}

	subcommand := args[0]
	apiURL := fmt.Sprintf("http://%s:%d", cfg.API.HTTP.Address, cfg.API.HTTP.Port)

	switch subcommand {
	case "status":
		return showStatus(apiURL)
	case "health":
		return showHealth(cfg)
	case "metrics":
		return showMetrics(cfg)
	case "config":
		return showConfig(cfg)
	case "tasks":
		if len(args) > 1 && args[1] == "running" {
			return showTasks(apiURL, "running")
		}
		return showTasks(apiURL, "")
	case "version":
		PrintVersion()
		return nil
	default:
		return fmt.Errorf("unknown show subcommand: %s", subcommand)
	}
}

// Handle 'task' commands
func handleTaskCommand(args []string, cfg *config.Config, logger *logrus.Logger) error {
	if len(args) < 1 {
		return fmt.Errorf("task command requires a subcommand")
	}

	subcommand := args[0]
	apiURL := fmt.Sprintf("http://%s:%d", cfg.API.HTTP.Address, cfg.API.HTTP.Port)

	switch subcommand {
	case "create":
		if len(args) < 2 {
			return fmt.Errorf("task create requires a command")
		}
		command := strings.Join(args[1:], " ")
		return createTask(apiURL, command)
	case "cancel":
		if len(args) < 2 {
			return fmt.Errorf("task cancel requires a task ID")
		}
		return cancelTask(apiURL, args[1])
	case "logs":
		if len(args) < 2 {
			return fmt.Errorf("task logs requires a task ID")
		}
		return showTaskLogs(apiURL, args[1])
	default:
		return fmt.Errorf("unknown task subcommand: %s", subcommand)
	}
}

// Handle 'file' commands
func handleFileCommand(args []string, cfg *config.Config, logger *logrus.Logger) error {
	if len(args) < 1 {
		return fmt.Errorf("file command requires a subcommand")
	}

	subcommand := args[0]
	apiURL := fmt.Sprintf("http://%s:%d", cfg.API.HTTP.Address, cfg.API.HTTP.Port)

	switch subcommand {
	case "list":
		if len(args) < 2 {
			return fmt.Errorf("file list requires a path")
		}
		return listFiles(apiURL, args[1])
	case "copy":
		if len(args) < 3 {
			return fmt.Errorf("file copy requires source and destination")
		}
		return copyFile(apiURL, args[1], args[2])
	case "move":
		if len(args) < 3 {
			return fmt.Errorf("file move requires source and destination")
		}
		return moveFile(apiURL, args[1], args[2])
	case "delete":
		if len(args) < 2 {
			return fmt.Errorf("file delete requires a path")
		}
		return deleteFile(apiURL, args[1])
	default:
		return fmt.Errorf("unknown file subcommand: %s", subcommand)
	}
}

// Handle 'config' commands
func handleConfigCommand(args []string, cfg *config.Config, logger *logrus.Logger) error {
	if len(args) < 1 {
		return fmt.Errorf("config command requires a subcommand")
	}

	subcommand := args[0]

	switch subcommand {
	case "validate":
		return validateConfig(cfg)
	case "test":
		return testConfig(cfg)
	default:
		return fmt.Errorf("unknown config subcommand: %s", subcommand)
	}
}

// Show status
func showStatus(apiURL string) error {
	resp, err := http.Get(apiURL + "/api/v1/status")
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Println("🔍 Agent Status:")
	printJSON(result)
	return nil
}

// Show health
func showHealth(cfg *config.Config) error {
	healthURL := fmt.Sprintf("http://%s:%d/health", cfg.Health.Address, cfg.Health.Port)
	resp, err := http.Get(healthURL)
	if err != nil {
		return fmt.Errorf("failed to get health: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Println("🏥 System Health:")
	printJSON(result)
	return nil
}

// Show metrics
func showMetrics(cfg *config.Config) error {
	metricsURL := fmt.Sprintf("http://%s:%d/metrics", cfg.Metrics.Address, cfg.Metrics.Port)
	resp, err := http.Get(metricsURL)
	if err != nil {
		return fmt.Errorf("failed to get metrics: %w", err)
	}
	defer resp.Body.Close()

	fmt.Println("📊 Metrics:")
	// For metrics, we'll show first 20 lines
	buf := make([]byte, 2048)
	n, _ := resp.Body.Read(buf)
	lines := strings.Split(string(buf[:n]), "\n")
	for i, line := range lines {
		if i >= 20 {
			fmt.Println("... (truncated)")
			break
		}
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		fmt.Println(line)
	}
	return nil
}

// Show config
func showConfig(cfg *config.Config) error {
	fmt.Println("⚙️  Configuration:")
	
	fmt.Printf("Agent ID: %s\n", cfg.Agent.ID)
	fmt.Printf("Agent Name: %s\n", cfg.Agent.Name)
	fmt.Printf("Environment: %s\n", cfg.Agent.Environment)
	fmt.Printf("Region: %s\n", cfg.Agent.Region)
	
	fmt.Printf("\nAPI Configuration:\n")
	fmt.Printf("  HTTP: %s:%d (enabled: %t)\n", cfg.API.HTTP.Address, cfg.API.HTTP.Port, cfg.API.HTTP.Enabled)
	fmt.Printf("  gRPC: %s:%d (enabled: %t)\n", cfg.API.GRPC.Address, cfg.API.GRPC.Port, cfg.API.GRPC.Enabled)
	
	fmt.Printf("\nStorage:\n")
	fmt.Printf("  Data Dir: %s\n", cfg.Storage.DataDir)
	fmt.Printf("  Temp Dir: %s\n", cfg.Storage.TempDir)
	
	fmt.Printf("\nExecutor:\n")
	fmt.Printf("  Workers: %d\n", cfg.Executor.WorkerPoolSize)
	fmt.Printf("  Queue Size: %d\n", cfg.Executor.QueueSize)
	fmt.Printf("  Task Timeout: %s\n", cfg.Executor.TaskTimeout)
	
	return nil
}

// Show tasks
func showTasks(apiURL, filter string) error {
	url := apiURL + "/api/v1/tasks"
	if filter != "" {
		url += "?filter=" + filter
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get tasks: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if filter == "running" {
		fmt.Println("⚡ Running Tasks:")
	} else {
		fmt.Println("📋 All Tasks:")
	}
	printJSON(result)
	return nil
}

// Create task
func createTask(apiURL, command string) error {
	// Parse command into parts
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}

	taskData := map[string]interface{}{
		"type":    "shell",
		"name":    "cli-task",
		"command": parts[0],
		"args":    parts[1:],
	}

	jsonData, err := json.Marshal(taskData)
	if err != nil {
		return fmt.Errorf("failed to marshal task data: %w", err)
	}

	resp, err := http.Post(apiURL+"/api/v1/tasks", "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Println("✅ Task Created:")
	printJSON(result)
	return nil
}

// Cancel task
func cancelTask(apiURL, taskID string) error {
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", apiURL+"/api/v1/tasks/"+taskID, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to cancel task: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Printf("✅ Task %s cancelled successfully\n", taskID)
	} else {
		fmt.Printf("❌ Failed to cancel task %s (status: %d)\n", taskID, resp.StatusCode)
	}
	return nil
}

// Show task logs
func showTaskLogs(apiURL, taskID string) error {
	resp, err := http.Get(apiURL + "/api/v1/tasks/" + taskID)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Printf("📄 Task %s Logs:\n", taskID)
	printJSON(result)
	return nil
}

// List files
func listFiles(apiURL, path string) error {
	operation := map[string]interface{}{
		"type":        "list",
		"source_path": path,
	}

	jsonData, err := json.Marshal(operation)
	if err != nil {
		return fmt.Errorf("failed to marshal operation: %w", err)
	}

	resp, err := http.Post(apiURL+"/api/v1/files", "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return fmt.Errorf("failed to list files: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Printf("📁 Files in %s:\n", path)
	printJSON(result)
	return nil
}

// Copy file
func copyFile(apiURL, source, dest string) error {
	operation := map[string]interface{}{
		"type":        "copy",
		"source_path": source,
		"dest_path":   dest,
	}

	jsonData, err := json.Marshal(operation)
	if err != nil {
		return fmt.Errorf("failed to marshal operation: %w", err)
	}

	resp, err := http.Post(apiURL+"/api/v1/files", "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Printf("📋 Copy %s → %s:\n", source, dest)
	printJSON(result)
	return nil
}

// Move file
func moveFile(apiURL, source, dest string) error {
	operation := map[string]interface{}{
		"type":        "move",
		"source_path": source,
		"dest_path":   dest,
	}

	jsonData, err := json.Marshal(operation)
	if err != nil {
		return fmt.Errorf("failed to marshal operation: %w", err)
	}

	resp, err := http.Post(apiURL+"/api/v1/files", "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return fmt.Errorf("failed to move file: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Printf("🔄 Move %s → %s:\n", source, dest)
	printJSON(result)
	return nil
}

// Delete file
func deleteFile(apiURL, path string) error {
	operation := map[string]interface{}{
		"type":        "delete",
		"source_path": path,
	}

	jsonData, err := json.Marshal(operation)
	if err != nil {
		return fmt.Errorf("failed to marshal operation: %w", err)
	}

	resp, err := http.Post(apiURL+"/api/v1/files", "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Printf("🗑️  Delete %s:\n", path)
	printJSON(result)
	return nil
}

// Validate config
func validateConfig(cfg *config.Config) error {
	fmt.Println("✅ Configuration validation:")
	
	if err := cfg.Validate(); err != nil {
		fmt.Printf("❌ Configuration is invalid: %s\n", err)
		return err
	}
	
	fmt.Println("✅ Configuration is valid")
	return nil
}

// Test config
func testConfig(cfg *config.Config) error {
	fmt.Println("🔍 Testing configuration...")
	
	// Test API connectivity
	apiURL := fmt.Sprintf("http://%s:%d", cfg.API.HTTP.Address, cfg.API.HTTP.Port)
	client := &http.Client{Timeout: 5 * time.Second}
	
	fmt.Printf("Testing API connectivity to %s...\n", apiURL)
	resp, err := client.Get(apiURL + "/api/v1/status")
	if err != nil {
		fmt.Printf("❌ API test failed: %s\n", err)
		return err
	}
	resp.Body.Close()
	
	fmt.Println("✅ API connectivity test passed")
	
	// Test health endpoint
	healthURL := fmt.Sprintf("http://%s:%d/health", cfg.Health.Address, cfg.Health.Port)
	fmt.Printf("Testing health endpoint at %s...\n", healthURL)
	resp, err = client.Get(healthURL)
	if err != nil {
		fmt.Printf("❌ Health test failed: %s\n", err)
		return err
	}
	resp.Body.Close()
	
	fmt.Println("✅ Health endpoint test passed")
	
	// Test metrics endpoint
	metricsURL := fmt.Sprintf("http://%s:%d/metrics", cfg.Metrics.Address, cfg.Metrics.Port)
	fmt.Printf("Testing metrics endpoint at %s...\n", metricsURL)
	resp, err = client.Get(metricsURL)
	if err != nil {
		fmt.Printf("❌ Metrics test failed: %s\n", err)
		return err
	}
	resp.Body.Close()
	
	fmt.Println("✅ Metrics endpoint test passed")
	fmt.Println("🎉 All configuration tests passed!")
	
	return nil
}

// Helper function to print JSON nicely
func printJSON(data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Error formatting JSON: %s\n", err)
		return
	}
	fmt.Println(string(jsonData))
}

// Show help for CLI commands
func showCLIHelp() {
	fmt.Println("Ducla Cloud Agent CLI Commands:")
	fmt.Println("")
	fmt.Println("USAGE:")
	fmt.Println("  ducla-agent [OPTIONS] [COMMAND]")
	fmt.Println("")
	fmt.Println("OPTIONS:")
	fmt.Println("  -c, --config FILE    Configuration file path")
	fmt.Println("  -d, --debug          Enable debug logging")
	fmt.Println("  -v, --version        Show version information")
	fmt.Println("  -h, --help           Show this help message")
	fmt.Println("")
	fmt.Println("COMMANDS:")
	fmt.Println("  start                Start the agent daemon (default)")
	fmt.Println("")
	fmt.Println("  show status          Display agent runtime status")
	fmt.Println("  show health          Display system health information")
	fmt.Println("  show metrics         Display current metrics")
	fmt.Println("  show config          Display current configuration")
	fmt.Println("  show tasks           List all tasks")
	fmt.Println("  show tasks running   List only running tasks")
	fmt.Println("  show version         Display version information")
	fmt.Println("")
	fmt.Println("  task create COMMAND  Create and execute a new task")
	fmt.Println("  task cancel TASK_ID  Cancel a running task")
	fmt.Println("  task logs TASK_ID    Show logs for a specific task")
	fmt.Println("")
	fmt.Println("  file list PATH       List files in directory")
	fmt.Println("  file copy SRC DEST   Copy file from source to destination")
	fmt.Println("  file move SRC DEST   Move file from source to destination")
	fmt.Println("  file delete PATH     Delete file or directory")
	fmt.Println("")
	fmt.Println("  config validate      Validate configuration file")
	fmt.Println("  config test          Test configuration and connectivity")
	fmt.Println("")
	fmt.Println("EXAMPLES:")
	fmt.Println("  ducla-agent                              # Start agent")
	fmt.Println("  ducla-agent show status                  # Show status")
	fmt.Println("  ducla-agent task create \"echo hello\"     # Create task")
	fmt.Println("  ducla-agent file list /tmp               # List files")
	fmt.Println("  ducla-agent config validate              # Validate config")
	fmt.Println("")
	fmt.Println("For more information, see: man ducla-agent")
}