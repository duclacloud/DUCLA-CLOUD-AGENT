package executor

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

// CommandExecutor executes command tasks
type CommandExecutor struct {
	logger *logrus.Logger
}

// NewCommandExecutor creates a new command executor
func NewCommandExecutor(logger *logrus.Logger) *CommandExecutor {
	return &CommandExecutor{
		logger: logger,
	}
}

// Execute executes a command task
func (e *CommandExecutor) Execute(ctx context.Context, task *Task, result *TaskResult) error {
	e.logger.WithFields(logrus.Fields{
		"task_id": task.ID,
		"command": task.Command,
		"args":    task.Args,
	}).Debug("Executing command")

	// Create command
	cmd := exec.CommandContext(ctx, task.Command, task.Args...)

	// Set working directory
	if task.WorkingDir != "" {
		cmd.Dir = task.WorkingDir
	}

	// Set environment variables
	if len(task.Env) > 0 {
		env := make([]string, 0, len(task.Env))
		for key, value := range task.Env {
			env = append(env, fmt.Sprintf("%s=%s", key, value))
		}
		cmd.Env = append(cmd.Env, env...)
	}

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Execute command
	startTime := time.Now()
	err := cmd.Run()
	duration := time.Since(startTime)

	// Get exit code
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				exitCode = status.ExitStatus()
			}
		}
	}

	// Update result
	result.ExitCode = exitCode
	result.Output = stdout.String()
	if stderr.Len() > 0 {
		result.Error = stderr.String()
	}
	result.Metadata["duration_ms"] = duration.Milliseconds()
	result.Metadata["command"] = task.Command
	result.Metadata["args"] = task.Args

	e.logger.WithFields(logrus.Fields{
		"task_id":   task.ID,
		"exit_code": exitCode,
		"duration":  duration,
	}).Debug("Command execution completed")

	if err != nil && exitCode != 0 {
		return fmt.Errorf("command failed with exit code %d: %s", exitCode, stderr.String())
	}

	return nil
}

// ScriptExecutor executes script tasks
type ScriptExecutor struct {
	logger *logrus.Logger
}

// NewScriptExecutor creates a new script executor
func NewScriptExecutor(logger *logrus.Logger) *ScriptExecutor {
	return &ScriptExecutor{
		logger: logger,
	}
}

// Execute executes a script task
func (e *ScriptExecutor) Execute(ctx context.Context, task *Task, result *TaskResult) error {
	e.logger.WithFields(logrus.Fields{
		"task_id": task.ID,
		"script":  task.Command,
	}).Debug("Executing script")

	// Determine script interpreter
	interpreter := "/bin/sh"
	if scriptType, ok := task.Metadata["script_type"].(string); ok {
		switch scriptType {
		case "bash":
			interpreter = "/bin/bash"
		case "python":
			interpreter = "/usr/bin/python3"
		case "ruby":
			interpreter = "/usr/bin/ruby"
		case "perl":
			interpreter = "/usr/bin/perl"
		}
	}

	// Create command
	cmd := exec.CommandContext(ctx, interpreter, "-c", task.Command)

	// Set working directory
	if task.WorkingDir != "" {
		cmd.Dir = task.WorkingDir
	}

	// Set environment variables
	if len(task.Env) > 0 {
		env := make([]string, 0, len(task.Env))
		for key, value := range task.Env {
			env = append(env, fmt.Sprintf("%s=%s", key, value))
		}
		cmd.Env = append(cmd.Env, env...)
	}

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Execute script
	startTime := time.Now()
	err := cmd.Run()
	duration := time.Since(startTime)

	// Get exit code
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				exitCode = status.ExitStatus()
			}
		}
	}

	// Update result
	result.ExitCode = exitCode
	result.Output = stdout.String()
	if stderr.Len() > 0 {
		result.Error = stderr.String()
	}
	result.Metadata["duration_ms"] = duration.Milliseconds()
	result.Metadata["interpreter"] = interpreter

	e.logger.WithFields(logrus.Fields{
		"task_id":   task.ID,
		"exit_code": exitCode,
		"duration":  duration,
	}).Debug("Script execution completed")

	if err != nil && exitCode != 0 {
		return fmt.Errorf("script failed with exit code %d: %s", exitCode, stderr.String())
	}

	return nil
}

// FileExecutor executes file operation tasks
type FileExecutor struct {
	logger *logrus.Logger
}

// NewFileExecutor creates a new file executor
func NewFileExecutor(logger *logrus.Logger) *FileExecutor {
	return &FileExecutor{
		logger: logger,
	}
}

// Execute executes a file operation task
func (e *FileExecutor) Execute(ctx context.Context, task *Task, result *TaskResult) error {
	e.logger.WithFields(logrus.Fields{
		"task_id":   task.ID,
		"operation": task.Command,
	}).Debug("Executing file operation")

	// Parse operation
	operation := task.Command
	args := task.Args

	switch operation {
	case "copy":
		if len(args) < 2 {
			return fmt.Errorf("copy operation requires source and destination")
		}
		return e.copyFile(ctx, args[0], args[1], result)
	case "move":
		if len(args) < 2 {
			return fmt.Errorf("move operation requires source and destination")
		}
		return e.moveFile(ctx, args[0], args[1], result)
	case "delete":
		if len(args) < 1 {
			return fmt.Errorf("delete operation requires file path")
		}
		return e.deleteFile(ctx, args[0], result)
	case "chmod":
		if len(args) < 2 {
			return fmt.Errorf("chmod operation requires file path and mode")
		}
		return e.chmodFile(ctx, args[0], args[1], result)
	case "chown":
		if len(args) < 2 {
			return fmt.Errorf("chown operation requires file path and owner")
		}
		return e.chownFile(ctx, args[0], args[1], result)
	default:
		return fmt.Errorf("unsupported file operation: %s", operation)
	}
}

func (e *FileExecutor) copyFile(ctx context.Context, src, dst string, result *TaskResult) error {
	cmd := exec.CommandContext(ctx, "cp", "-r", src, dst)
	output, err := cmd.CombinedOutput()
	result.Output = string(output)
	return err
}

func (e *FileExecutor) moveFile(ctx context.Context, src, dst string, result *TaskResult) error {
	cmd := exec.CommandContext(ctx, "mv", src, dst)
	output, err := cmd.CombinedOutput()
	result.Output = string(output)
	return err
}

func (e *FileExecutor) deleteFile(ctx context.Context, path string, result *TaskResult) error {
	cmd := exec.CommandContext(ctx, "rm", "-rf", path)
	output, err := cmd.CombinedOutput()
	result.Output = string(output)
	return err
}

func (e *FileExecutor) chmodFile(ctx context.Context, path, mode string, result *TaskResult) error {
	cmd := exec.CommandContext(ctx, "chmod", mode, path)
	output, err := cmd.CombinedOutput()
	result.Output = string(output)
	return err
}

func (e *FileExecutor) chownFile(ctx context.Context, path, owner string, result *TaskResult) error {
	cmd := exec.CommandContext(ctx, "chown", owner, path)
	output, err := cmd.CombinedOutput()
	result.Output = string(output)
	return err
}

// HTTPExecutor executes HTTP request tasks
type HTTPExecutor struct {
	logger *logrus.Logger
}

// NewHTTPExecutor creates a new HTTP executor
func NewHTTPExecutor(logger *logrus.Logger) *HTTPExecutor {
	return &HTTPExecutor{
		logger: logger,
	}
}

// Execute executes an HTTP request task
func (e *HTTPExecutor) Execute(ctx context.Context, task *Task, result *TaskResult) error {
	e.logger.WithFields(logrus.Fields{
		"task_id": task.ID,
		"method":  task.Command,
		"url":     task.Args[0],
	}).Debug("Executing HTTP request")

	// Build curl command
	method := strings.ToUpper(task.Command)
	if method == "" {
		method = "GET"
	}

	args := []string{"-X", method}

	// Add headers
	if headers, ok := task.Metadata["headers"].(map[string]interface{}); ok {
		for key, value := range headers {
			args = append(args, "-H", fmt.Sprintf("%s: %v", key, value))
		}
	}

	// Add body
	if body, ok := task.Metadata["body"].(string); ok && body != "" {
		args = append(args, "-d", body)
	}

	// Add URL
	if len(task.Args) > 0 {
		args = append(args, task.Args[0])
	} else {
		return fmt.Errorf("URL is required for HTTP request")
	}

	// Execute curl command
	cmd := exec.CommandContext(ctx, "curl", args...)
	output, err := cmd.CombinedOutput()

	result.Output = string(output)
	result.Metadata["method"] = method
	result.Metadata["url"] = task.Args[0]

	if err != nil {
		result.Error = err.Error()
		return err
	}

	return nil
}

// DockerExecutor executes Docker tasks
type DockerExecutor struct {
	logger *logrus.Logger
}

// NewDockerExecutor creates a new Docker executor
func NewDockerExecutor(logger *logrus.Logger) *DockerExecutor {
	return &DockerExecutor{
		logger: logger,
	}
}

// Execute executes a Docker task
func (e *DockerExecutor) Execute(ctx context.Context, task *Task, result *TaskResult) error {
	e.logger.WithFields(logrus.Fields{
		"task_id":   task.ID,
		"operation": task.Command,
	}).Debug("Executing Docker operation")

	// Build docker command
	args := []string{task.Command}
	args = append(args, task.Args...)

	cmd := exec.CommandContext(ctx, "docker", args...)
	output, err := cmd.CombinedOutput()

	result.Output = string(output)
	result.Metadata["operation"] = task.Command

	if err != nil {
		result.Error = err.Error()
		return err
	}

	return nil
}

// KubernetesExecutor executes Kubernetes tasks
type KubernetesExecutor struct {
	logger *logrus.Logger
}

// NewKubernetesExecutor creates a new Kubernetes executor
func NewKubernetesExecutor(logger *logrus.Logger) *KubernetesExecutor {
	return &KubernetesExecutor{
		logger: logger,
	}
}

// Execute executes a Kubernetes task
func (e *KubernetesExecutor) Execute(ctx context.Context, task *Task, result *TaskResult) error {
	e.logger.WithFields(logrus.Fields{
		"task_id":   task.ID,
		"operation": task.Command,
	}).Debug("Executing Kubernetes operation")

	// Build kubectl command
	args := []string{task.Command}
	args = append(args, task.Args...)

	cmd := exec.CommandContext(ctx, "kubectl", args...)
	output, err := cmd.CombinedOutput()

	result.Output = string(output)
	result.Metadata["operation"] = task.Command

	if err != nil {
		result.Error = err.Error()
		return err
	}

	return nil
}

// CustomExecutor executes custom tasks
type CustomExecutor struct {
	logger *logrus.Logger
}

// NewCustomExecutor creates a new custom executor
func NewCustomExecutor(logger *logrus.Logger) *CustomExecutor {
	return &CustomExecutor{
		logger: logger,
	}
}

// Execute executes a custom task
func (e *CustomExecutor) Execute(ctx context.Context, task *Task, result *TaskResult) error {
	e.logger.WithFields(logrus.Fields{
		"task_id": task.ID,
		"custom":  task.Command,
	}).Debug("Executing custom task")

	// Custom task execution logic
	// This can be extended with plugin system
	result.Output = "Custom task execution not implemented"
	result.Metadata["custom_type"] = task.Command

	return fmt.Errorf("custom task execution not implemented")
}