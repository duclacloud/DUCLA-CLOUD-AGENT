package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ducla/cloud-agent/internal/agent"
	"github.com/ducla/cloud-agent/internal/config"
	"github.com/ducla/cloud-agent/pkg/utils/logger"
	"github.com/sirupsen/logrus"
)

func main() {
	var (
		configFile = flag.String("config", "/etc/ducla/agent.yaml", "Configuration file path")
		showVer    = flag.Bool("version", false, "Show version information")
		debug      = flag.Bool("debug", false, "Enable debug logging")
	)
	flag.Parse()

	if *showVer {
		PrintVersion()
		os.Exit(0)
	}

	// Initialize logger
	log := logger.New()
	if *debug {
		log.SetLevel(logrus.DebugLevel)
	}

	versionInfo := GetVersionInfo()
	log.WithFields(logrus.Fields{
		"version":    versionInfo.Version,
		"build_time": versionInfo.BuildTime,
		"git_commit": versionInfo.GitCommit,
	}).Info("Starting Ducla Cloud Agent")

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