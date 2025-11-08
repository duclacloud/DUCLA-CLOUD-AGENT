#!/bin/bash

# Release script for creating new versions

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Check if version is provided
if [ -z "$1" ]; then
    echo -e "${RED}Usage: $0 <version>${NC}"
    echo "Example: $0 v1.0.0"
    exit 1
fi

VERSION=$1

# Validate version format
if [[ ! $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo -e "${RED}Invalid version format. Use: vX.Y.Z${NC}"
    exit 1
fi

echo -e "${GREEN}Creating release $VERSION${NC}"
echo ""

# Check if git is clean
if [ -n "$(git status --porcelain)" ]; then
    echo -e "${RED}Git working directory is not clean${NC}"
    git status --short
    exit 1
fi

# Check if on main/master branch
BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [ "$BRANCH" != "main" ] && [ "$BRANCH" != "master" ]; then
    echo -e "${YELLOW}Warning: Not on main/master branch (current: $BRANCH)${NC}"
    read -p "Continue? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Run tests
echo -e "${YELLOW}Running tests...${NC}"
go test ./...

# Update version in files if needed
# echo $VERSION > VERSION

# Build binaries
echo -e "${YELLOW}Building binaries...${NC}"
VERSION=$VERSION ./scripts/build-cli.sh

# Create git tag
echo -e "${YELLOW}Creating git tag...${NC}"
git tag -a "$VERSION" -m "Release $VERSION"

# Push tag
echo -e "${YELLOW}Pushing tag to remote...${NC}"
git push origin "$VERSION"

echo ""
echo -e "${GREEN}Release $VERSION created successfully!${NC}"
echo ""
echo "Next steps:"
echo "1. Create GitHub release with binaries from dist/"
echo "2. Update documentation"
echo "3. Announce the release"
