package cli

import (
	"fmt"

	"github.com/duclacloud/DUCLA-CLOUD-AGENT/internal/cli/client"
	"github.com/spf13/cobra"
)

func NewHealthCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "health",
		Short: "Check agent health",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := client.NewClient(globalFlags.AgentURL, globalFlags.Token)
			
			health, err := c.GetHealth(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to get health status: %w", err)
			}

			if health.Status == "healthy" {
				fmt.Println("✓ Agent is healthy")
			} else {
				fmt.Printf("✗ Agent is unhealthy: %s\n", health.Status)
			}

			return printOutput(health, globalFlags.Output)
		},
	}

	return cmd
}
