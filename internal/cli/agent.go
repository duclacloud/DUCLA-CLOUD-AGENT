package cli

import (
	"fmt"

	"github.com/ducla/cloud-agent/internal/cli/client"
	"github.com/spf13/cobra"
)

func NewAgentCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Manage agent operations",
		Long:  "Commands for managing and monitoring the Ducla Cloud Agent",
	}

	cmd.AddCommand(newAgentInfoCommand())
	cmd.AddCommand(newAgentStatusCommand())
	cmd.AddCommand(newAgentListCommand())

	return cmd
}

func newAgentInfoCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Get agent information",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := client.NewClient(globalFlags.AgentURL, globalFlags.Token)
			
			info, err := c.GetAgentInfo(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to get agent info: %w", err)
			}

			return printOutput(info, globalFlags.Output)
		},
	}
}

func newAgentStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Get agent status",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := client.NewClient(globalFlags.AgentURL, globalFlags.Token)
			
			status, err := c.GetAgentStatus(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to get agent status: %w", err)
			}

			return printOutput(status, globalFlags.Output)
		},
	}
}

func newAgentListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all agents (requires master server)",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := client.NewClient(globalFlags.AgentURL, globalFlags.Token)
			
			agents, err := c.ListAgents(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to list agents: %w", err)
			}

			return printOutput(agents, globalFlags.Output)
		},
	}
}


