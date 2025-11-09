#!/bin/bash

# Build script for Ducla Cloud Agent with version info

set -e

# Version information
VERSION=${VERSION:-"1.0.0"}
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S_UTC')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GO_VERSION=$(go version | awk '{print $3}')

# Build flags
LDFLAGS="-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.gitCommit=${GIT_COMMIT} -X main.goVersion=${GO_VERSION}"

echo "Building Ducla Cloud Agent v${VERSION}..."
echo "Build Time: ${BUILD_TIME}"
echo "Git Commit: ${GIT_COMMIT}"
echo "Go Version: ${GO_VERSION}"
echo ""

# Create output directory
mkdir -p dist

# Build for current platform
echo "Building for current platform..."
go build -ldflags "${LDFLAGS}" -o dist/ducla-agent ./cmd/agent

# Build for multiple platforms
echo "Building for multiple platforms..."

# Linux AMD64
GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o dist/ducla-agent-linux-amd64 ./cmd/agent

# Linux ARM64
GOOS=linux GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o dist/ducla-agent-linux-arm64 ./cmd/agent

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o dist/ducla-agent-windows-amd64.exe ./cmd/agent

# macOS AMD64
GOOS=darwin GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o dist/ducla-agent-darwin-amd64 ./cmd/agent

# macOS ARM64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o dist/ducla-agent-darwin-arm64 ./cmd/agent

echo ""
echo "Build completed successfully!"
echo "Binaries created in dist/ directory:"
ls -la dist/

echo ""
echo "Testing version info..."
./dist/ducla-agent -version