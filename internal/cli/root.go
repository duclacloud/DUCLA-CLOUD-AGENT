package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

type GlobalFlags struct {
	AgentURL   string
	Token      string
	ConfigFile string
	Output     string
	Verbose    bool
}

var globalFlags GlobalFlags

func NewRootCommand(version, buildTime, gitCommit string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "duclactl",
		Short: "Ducla Cloud Agent CLI",
		Long: `duclactl is a command-line interface for managing and interacting with Ducla Cloud Agents.
		
It provides commands to monitor agent status, execute tasks, manage files, and more.`,
		Version: fmt.Sprintf("%s (built: %s, commit: %s)", version, buildTime, gitCommit),
	}

	// Global flags
	rootCmd.PersistentFlags().StringVar(&globalFlags.AgentURL, "agent-url", "http://localhost:8080", "Agent API URL")
	rootCmd.PersistentFlags().StringVar(&globalFlags.Token, "token", "", "Authentication token")
	rootCmd.PersistentFlags().StringVarP(&globalFlags.ConfigFile, "config", "c", "", "Config file path")
	rootCmd.PersistentFlags().StringVarP(&globalFlags.Output, "output", "o", "table", "Output format (table, json, yaml)")
	rootCmd.PersistentFlags().BoolVarP(&globalFlags.Verbose, "verbose", "v", false, "Verbose output")

	// Add subcommands
	rootCmd.AddCommand(NewAgentCommand())
	rootCmd.AddCommand(NewTaskCommand())
	rootCmd.AddCommand(NewFileCommand())
	rootCmd.AddCommand(NewHealthCommand())
	rootCmd.AddCommand(NewMetricsCommand())
	rootCmd.AddCommand(NewConfigCommand())

	return rootCmd
}
