package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Config struct {
	AgentURL string `yaml:"agent_url"`
	Token    string `yaml:"token"`
	Output   string `yaml:"output"`
}

func NewConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage CLI configuration",
	}

	cmd.AddCommand(newConfigInitCommand())
	cmd.AddCommand(newConfigShowCommand())
	cmd.AddCommand(newConfigSetCommand())

	return cmd
}

func newConfigInitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize configuration file",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath := getConfigPath()
			
			if _, err := os.Stat(configPath); err == nil {
				return fmt.Errorf("config file already exists at %s", configPath)
			}

			cfg := Config{
				AgentURL: "http://localhost:8080",
				Token:    "",
				Output:   "table",
			}

			data, err := yaml.Marshal(cfg)
			if err != nil {
				return fmt.Errorf("failed to marshal config: %w", err)
			}

			if err := os.WriteFile(configPath, data, 0600); err != nil {
				return fmt.Errorf("failed to write config: %w", err)
			}

			fmt.Printf("Configuration initialized at %s\n", configPath)
			return nil
		},
	}
}

func newConfigShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return err
			}

			data, err := yaml.Marshal(cfg)
			if err != nil {
				return fmt.Errorf("failed to marshal config: %w", err)
			}

			fmt.Print(string(data))
			return nil
		},
	}
}

func newConfigSetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "set [key] [value]",
		Short: "Set configuration value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				cfg = &Config{}
			}

			key, value := args[0], args[1]
			
			switch key {
			case "agent_url":
				cfg.AgentURL = value
			case "token":
				cfg.Token = value
			case "output":
				cfg.Output = value
			default:
				return fmt.Errorf("unknown config key: %s", key)
			}

			data, err := yaml.Marshal(cfg)
			if err != nil {
				return fmt.Errorf("failed to marshal config: %w", err)
			}

			configPath := getConfigPath()
			if err := os.WriteFile(configPath, data, 0600); err != nil {
				return fmt.Errorf("failed to write config: %w", err)
			}

			fmt.Printf("Configuration updated: %s = %s\n", key, value)
			return nil
		},
	}
}

func getConfigPath() string {
	if globalFlags.ConfigFile != "" {
		return globalFlags.ConfigFile
	}
	
	home, err := os.UserHomeDir()
	if err != nil {
		return ".duclactl.yaml"
	}
	
	return home + "/.duclactl.yaml"
}

func loadConfig() (*Config, error) {
	configPath := getConfigPath()
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
