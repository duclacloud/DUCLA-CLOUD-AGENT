#!/bin/bash

# Simple build script for Ducla Cloud Agent

set -e

# Version information
VERSION="1.0.0"
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S_UTC')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build flags
LDFLAGS="-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.gitCommit=${GIT_COMMIT}"

echo "🚀 Building Ducla Cloud Agent v${VERSION}..."
echo "📅 Build Time: ${BUILD_TIME}"
echo "🔗 Git Commit: ${GIT_COMMIT}"
echo ""

# Build
go build -ldflags "${LDFLAGS}" -o ducla-agent ./cmd/agent

echo "✅ Build completed successfully!"
echo "📦 Binary: ./ducla-agent"
echo ""

# Test version
echo "🔍 Testing version info..."
./ducla-agent -version