package main

import (
	"fmt"
	"runtime"
)

// Version information - set by build flags
var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
	goVersion = runtime.Version()
)

// VersionInfo holds version information
type VersionInfo struct {
	Version   string `json:"version"`
	BuildTime string `json:"build_time"`
	GitCommit string `json:"git_commit"`
	GoVersion string `json:"go_version"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
}

// GetVersionInfo returns version information
func GetVersionInfo() VersionInfo {
	return VersionInfo{
		Version:   version,
		BuildTime: buildTime,
		GitCommit: gitCommit,
		GoVersion: goVersion,
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
}

// PrintVersion prints version information to stdout
func PrintVersion() {
	info := GetVersionInfo()
	fmt.Printf("Ducla Cloud Agent\n")
	fmt.Printf("  Version:    %s\n", info.Version)
	fmt.Printf("  Build Time: %s\n", info.BuildTime)
	fmt.Printf("  Git Commit: %s\n", info.GitCommit)
	fmt.Printf("  Go Version: %s\n", info.GoVersion)
	fmt.Printf("  OS/Arch:    %s/%s\n", info.OS, info.Arch)
}