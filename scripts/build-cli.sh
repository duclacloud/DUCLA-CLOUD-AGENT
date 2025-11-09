#!/bin/bash

# Build script for Ducla Cloud Agent CLI
# Supports multiple platforms and architectures

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
APP_NAME="ducla-agent"
VERSION=${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev")}
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GO_VERSION=$(go version | awk '{print $3}')

# Build directory
BUILD_DIR="bin"
DIST_DIR="dist"

# Platforms to build for
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

# Ldflags for version information
LDFLAGS="-w -s \
    -X 'main.Version=${VERSION}' \
    -X 'main.BuildTime=${BUILD_TIME}' \
    -X 'main.GitCommit=${GIT_COMMIT}' \
    -X 'main.GoVersion=${GO_VERSION}'"

echo -e "${GREEN}Building Ducla Cloud Agent CLI${NC}"
echo "Version: ${VERSION}"
echo "Build Time: ${BUILD_TIME}"
echo "Git Commit: ${GIT_COMMIT}"
echo ""

# Clean previous builds
echo -e "${YELLOW}Cleaning previous builds...${NC}"
rm -rf ${BUILD_DIR} ${DIST_DIR}
mkdir -p ${BUILD_DIR} ${DIST_DIR}

# Build for current platform (quick build)
if [ "$1" == "quick" ]; then
    echo -e "${GREEN}Building for current platform...${NC}"
    go build -ldflags="${LDFLAGS}" -o ${BUILD_DIR}/${APP_NAME} ./cmd/agent
    echo -e "${GREEN}✓ Build complete: ${BUILD_DIR}/${APP_NAME}${NC}"
    exit 0
fi

# Build for all platforms
for platform in "${PLATFORMS[@]}"; do
    IFS='/' read -r -a platform_split <<< "$platform"
    GOOS="${platform_split[0]}"
    GOARCH="${platform_split[1]}"
    
    output_name="${APP_NAME}-${VERSION}-${GOOS}-${GOARCH}"
    
    if [ "$GOOS" == "windows" ]; then
        output_name="${output_name}.exe"
    fi
    
    echo -e "${YELLOW}Building for ${GOOS}/${GOARCH}...${NC}"
    
    GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 go build \
        -ldflags="${LDFLAGS}" \
        -o "${DIST_DIR}/${output_name}" \
        ./cmd/agent
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Built: ${output_name}${NC}"
        
        # Create archive
        cd ${DIST_DIR}
        if [ "$GOOS" == "windows" ]; then
            zip "${output_name%.exe}.zip" "${output_name}"
        else
            tar -czf "${output_name}.tar.gz" "${output_name}"
        fi
        cd ..
    else
        echo -e "${RED}✗ Failed to build for ${GOOS}/${GOARCH}${NC}"
    fi
done

# Generate checksums
echo -e "${YELLOW}Generating checksums...${NC}"
cd ${DIST_DIR}
sha256sum *.tar.gz *.zip 2>/dev/null > checksums.txt || shasum -a 256 *.tar.gz *.zip > checksums.txt
cd ..

echo ""
echo -e "${GREEN}Build complete!${NC}"
echo "Binaries are in: ${DIST_DIR}/"
ls -lh ${DIST_DIR}/
