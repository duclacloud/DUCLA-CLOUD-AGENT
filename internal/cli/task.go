package cli

import (
	"fmt"

	"github.com/ducla/cloud-agent/internal/cli/client"
	"github.com/spf13/cobra"
)

func NewTaskCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "Manage tasks",
		Long:  "Commands for managing and monitoring tasks on the agent",
	}

	cmd.AddCommand(newTaskListCommand())
	cmd.AddCommand(newTaskGetCommand())
	cmd.AddCommand(newTaskExecuteCommand())
	cmd.AddCommand(newTaskCancelCommand())

	return cmd
}

func newTaskListCommand() *cobra.Command {
	var status string
	
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := client.NewClient(globalFlags.AgentURL, globalFlags.Token)
			
			tasks, err := c.ListTasks(cmd.Context(), status)
			if err != nil {
				return fmt.Errorf("failed to list tasks: %w", err)
			}

			return printOutput(tasks, globalFlags.Output)
		},
	}
	
	cmd.Flags().StringVar(&status, "status", "", "Filter by status (pending, running, completed, failed)")
	
	return cmd
}

func newTaskGetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "get [task-id]",
		Short: "Get task details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := client.NewClient(globalFlags.AgentURL, globalFlags.Token)
			
			task, err := c.GetTask(cmd.Context(), args[0])
			if err != nil {
				return fmt.Errorf("failed to get task: %w", err)
			}

			return printOutput(task, globalFlags.Output)
		},
	}
}

func newTaskExecuteCommand() *cobra.Command {
	var (
		taskType string
		payload  string
		timeout  string
	)
	
	cmd := &cobra.Command{
		Use:   "execute",
		Short: "Execute a new task",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := client.NewClient(globalFlags.AgentURL, globalFlags.Token)
			
			result, err := c.ExecuteTask(cmd.Context(), taskType, payload, timeout)
			if err != nil {
				return fmt.Errorf("failed to execute task: %w", err)
			}

			fmt.Printf("Task created: %s\n", result.TaskID)
			return printOutput(result, globalFlags.Output)
		},
	}
	
	cmd.Flags().StringVar(&taskType, "type", "", "Task type (required)")
	cmd.Flags().StringVar(&payload, "payload", "{}", "Task payload (JSON)")
	cmd.Flags().StringVar(&timeout, "timeout", "30m", "Task timeout")
	cmd.MarkFlagRequired("type")
	
	return cmd
}

func newTaskCancelCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "cancel [task-id]",
		Short: "Cancel a running task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := client.NewClient(globalFlags.AgentURL, globalFlags.Token)
			
			if err := c.CancelTask(cmd.Context(), args[0]); err != nil {
				return fmt.Errorf("failed to cancel task: %w", err)
			}

			fmt.Printf("Task %s cancelled successfully\n", args[0])
			return nil
		},
	}
}
