package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/duclacloud/DUCLA-CLOUD-AGENT/internal/agent"
	"github.com/duclacloud/DUCLA-CLOUD-AGENT/internal/config"
	"github.com/sirupsen/logrus"
)

func main() {
	var (
		configFile = flag.String("config", "", "Configuration file path")
		showVer    = flag.Bool("version", false, "Show version information")
		debug      = flag.Bool("debug", false, "Enable debug logging")
		help       = flag.Bool("help", false, "Show help information")
	)
	flag.Parse()

	if *help {
		showCLIHelp()
		os.Exit(0)
	}

	if *showVer {
		PrintVersion()
		os.Exit(0)
	}

	// Check for CLI commands
	args := flag.Args()
	if len(args) > 0 {
		// Handle CLI commands
		handleCLICommands(args, *configFile, *debug)
		return
	}

	// Initialize logger
	log := logrus.New()
	if *debug {
		log.SetLevel(logrus.DebugLevel)
	}

	versionInfo := GetVersionInfo()
	log.WithFields(logrus.Fields{
		"version":    versionInfo.Version,
		"build_time": versionInfo.BuildTime,
		"git_commit": versionInfo.GitCommit,
	}).Info("Starting Ducla Cloud Agent")

	// Determine config file path
	if *configFile == "" {
		// Try current directory first, then system locations
		if _, err := os.Stat("agent.yaml"); err == nil {
			*configFile = "agent.yaml"
		} else if _, err := os.Stat("/etc/ducla/agent.yaml"); err == nil {
			*configFile = "/etc/ducla/agent.yaml"
		} else {
			*configFile = "agent.yaml" // Use default
		}
	}

	// Load configuration
	cfg, err := config.Load(*configFile)
	if err != nil {
		log.WithError(err).Fatal("Failed to load configuration")
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.WithError(err).Fatal("Invalid configuration")
	}

	log.WithField("config_file", *configFile).Info("Configuration loaded successfully")

	// Create agent instance
	agentInstance, err := agent.New(cfg, log)
	if err != nil {
		log.WithError(err).Fatal("Failed to create agent instance")
	}

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.WithField("signal", sig).Info("Received shutdown signal")
		cancel()
	}()

	// Start agent
	log.Info("Starting agent services...")
	if err := agentInstance.Start(ctx); err != nil {
		log.WithError(err).Fatal("Failed to start agent")
	}

	// Wait for shutdown signal
	<-ctx.Done()

	// Graceful shutdown
	log.Info("Shutting down agent...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := agentInstance.Stop(shutdownCtx); err != nil {
		log.WithError(err).Error("Error during shutdown")
		os.Exit(1)
	}

	log.Info("Agent stopped successfully")
}

// handleCLICommands processes CLI commands
func handleCLICommands(args []string, configFile string, debug bool) {
	// Initialize logger for CLI
	log := logrus.New()
	if debug {
		log.SetLevel(logrus.DebugLevel)
	}

	// Determine config file path
	if configFile == "" {
		// Try current directory first, then system locations
		if _, err := os.Stat("agent.yaml"); err == nil {
			configFile = "agent.yaml"
		} else if _, err := os.Stat("/etc/ducla/agent.yaml"); err == nil {
			configFile = "/etc/ducla/agent.yaml"
		} else {
			configFile = "agent.yaml" // Use default even if not exists
		}
	}

	// Load configuration for CLI commands
	cfg, err := config.Load(configFile)
	if err != nil {
		// For CLI commands, try to use default config if file doesn't exist
		cfg = &config.Config{
			Agent: config.AgentConfig{
				ID: "cli-agent",
				Name: "CLI Agent",
			},
			API: config.APIConfig{
				HTTP: config.HTTPConfig{
					Address: "127.0.0.1",
					Port: 8080,
				},
			},
			Health: config.HealthConfig{
				Address: "127.0.0.1",
				Port: 8081,
			},
			Metrics: config.MetricsConfig{
				Address: "127.0.0.1",
				Port: 9090,
			},
		}
		
		// Only show warning for config-dependent commands that need actual config
		if (args[0] == "show" && len(args) > 1 && (args[1] == "status" || args[1] == "health" || args[1] == "metrics")) ||
		   args[0] == "task" || args[0] == "file" || 
		   (args[0] == "config" && len(args) > 1 && args[1] == "test") {
			log.WithError(err).Warn("Could not load configuration, using defaults")
		}
	}

	// Handle the CLI command
	if err := handleCLICommand(args, cfg, log); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}