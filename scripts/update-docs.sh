#!/bin/bash

# Script to update documentation to GitHub

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}📚 Updating Documentation to GitHub${NC}"
echo ""

# Check if we're in a git repository
if [ ! -d ".git" ]; then
    echo -e "${RED}❌ Not in a git repository${NC}"
    exit 1
fi

# Check if git is configured
if ! git config user.name >/dev/null 2>&1; then
    echo -e "${RED}❌ Git user.name not configured${NC}"
    echo "Run: git config --global user.name 'Your Name'"
    exit 1
fi

if ! git config user.email >/dev/null 2>&1; then
    echo -e "${RED}❌ Git user.email not configured${NC}"
    echo "Run: git config --global user.email 'your.email@example.com'"
    exit 1
fi

# Check for uncommitted changes
if ! git diff-index --quiet HEAD --; then
    echo -e "${YELLOW}⚠️  You have uncommitted changes${NC}"
    echo ""
    git status --porcelain
    echo ""
    read -p "Do you want to commit these changes? (y/N): " -n 1 -r
    echo
    
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${YELLOW}📝 Committing changes...${NC}"
        
        # Add all documentation files
        git add README.md
        git add README-VI.md
        git add USER-GUIDE.md
        git add DEPLOYMENT-SUMMARY.md
        git add API-REFERENCE.md
        git add dong-file-rpm-deb-exe.md
        git add CONTRIBUTING.md
        git add LICENSE
        git add docs/
        git add scripts/
        git add configs/
        git add release/
        git add dist/
        
        # Commit with descriptive message
        COMMIT_MSG="📚 Update documentation and release files

- Add comprehensive USER-GUIDE.md
- Update DEPLOYMENT-SUMMARY.md with latest info
- Add distribution and deployment scripts
- Include DEB/RPM packages and Docker image
- Update API documentation
- Add installation and uninstall scripts

Version: $(git describe --tags --always --dirty 2>/dev/null || echo '1.0.0')
Date: $(date '+%Y-%m-%d %H:%M:%S UTC')"

        git commit -m "$COMMIT_MSG"
        echo -e "${GREEN}✅ Changes committed${NC}"
    else
        echo -e "${BLUE}ℹ️  Skipping commit${NC}"
    fi
fi

# Check if we have a remote
if ! git remote get-url origin >/dev/null 2>&1; then
    echo -e "${RED}❌ No 'origin' remote configured${NC}"
    echo "Add remote: git remote add origin <repository-url>"
    exit 1
fi

# Get current branch
CURRENT_BRANCH=$(git branch --show-current)
echo -e "${BLUE}📍 Current branch: ${CURRENT_BRANCH}${NC}"

# Push to GitHub
echo -e "${YELLOW}🚀 Pushing to GitHub...${NC}"
git push origin $CURRENT_BRANCH

echo -e "${GREEN}✅ Documentation updated on GitHub!${NC}"
echo ""

# Show repository info
REPO_URL=$(git remote get-url origin | sed 's/\.git$//')
if [[ $REPO_URL == git@* ]]; then
    REPO_URL=$(echo $REPO_URL | sed 's/git@github.com:/https:\/\/github.com\//')
fi

echo -e "${BLUE}🔗 Repository: ${REPO_URL}${NC}"
echo -e "${BLUE}📚 Documentation files updated:${NC}"
echo "  • README.md - Main project documentation"
echo "  • README-VI.md - Vietnamese documentation"
echo "  • USER-GUIDE.md - Comprehensive user guide"
echo "  • DEPLOYMENT-SUMMARY.md - Deployment information"
echo "  • API-REFERENCE.md - API documentation"
echo "  • CONTRIBUTING.md - Contribution guidelines"

echo ""
echo -e "${BLUE}📦 Release files available:${NC}"
echo "  • DEB package: dist/ducla-agent_*.deb"
echo "  • RPM package: dist/ducla-agent-*.rpm"
echo "  • Docker image: ducla/cloud-agent:latest"
echo "  • Installation scripts: release/install.sh"

echo ""
echo -e "${BLUE}🎯 Next steps:${NC}"
echo "1. Create GitHub release: ./scripts/deploy.sh github"
echo "2. Update package repositories"
echo "3. Announce release to users"
echo "4. Monitor deployment metrics"

echo ""
echo -e "${GREEN}🎉 Documentation update complete!${NC}"