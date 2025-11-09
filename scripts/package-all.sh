#!/bin/bash

# Build all packages

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

VERSION=${VERSION:-$(git describe --tags --always --dirty 2>/dev/null | sed 's/^v//' || echo "1.0.0")}

echo -e "${GREEN}Building all packages${NC}"
echo "Version: ${VERSION}"
echo ""

# Build DEB package
echo -e "${YELLOW}Building DEB package...${NC}"
./scripts/package-deb.sh

# Build RPM package
echo -e "${YELLOW}Building RPM package...${NC}"
./scripts/package-rpm.sh

echo ""
echo -e "${GREEN}All packages built successfully!${NC}"
echo ""
echo "Packages:"
ls -lh dist/
