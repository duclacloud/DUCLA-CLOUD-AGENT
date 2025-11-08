package cli

import (
	"fmt"

	"github.com/ducla/cloud-agent/internal/cli/client"
	"github.com/spf13/cobra"
)

func NewMetricsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metrics",
		Short: "Get agent metrics",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := client.NewClient(globalFlags.AgentURL, globalFlags.Token)
			
			metrics, err := c.GetMetrics(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to get metrics: %w", err)
			}

			return printOutput(metrics, globalFlags.Output)
		},
	}

	return cmd
}
