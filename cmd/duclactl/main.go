package main

import (
	"fmt"
	"os"

	"github.com/ducla/cloud-agent/internal/cli"
	"github.com/spf13/cobra"
)

var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

func main() {
	rootCmd := cli.NewRootCommand(version, buildTime, gitCommit)
	
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
