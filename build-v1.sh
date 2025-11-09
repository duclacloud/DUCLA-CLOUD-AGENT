#!/bin/bash

# Build Ducla Cloud Agent v1.0.0 for current platform

set -e

# Version information
VERSION="1.0.0"
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S_UTC')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build flags
LDFLAGS="-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.gitCommit=${GIT_COMMIT}"

echo "🚀 Building Ducla Cloud Agent v${VERSION} for $(uname -s)/$(uname -m)..."
echo "📅 Build Time: ${BUILD_TIME}"
echo "🔗 Git Commit: ${GIT_COMMIT}"
echo ""

# Create dist directory
mkdir -p dist

# Build for current platform
go build -ldflags "${LDFLAGS}" -o dist/ducla-agent-v${VERSION} ./cmd/agent

# Also create a symlink for convenience
ln -sf ducla-agent-v${VERSION} dist/ducla-agent

echo "✅ Build completed successfully!"
echo "📦 Binary: ./dist/ducla-agent-v${VERSION}"
echo "🔗 Symlink: ./dist/ducla-agent"
echo ""

# Test version
echo "🔍 Testing version info..."
./dist/ducla-agent -version

echo ""
echo "🎉 Ducla Cloud Agent v${VERSION} is ready!"
echo "📁 Files created:"
ls -la dist/