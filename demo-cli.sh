#!/bin/bash

# Demo CLI commands for Ducla Cloud Agent

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

AGENT_BINARY="./dist/ducla-agent"

print_header() {
    echo -e "\n${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}\n"
}

print_step() {
    echo -e "${YELLOW}➤ $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

# Check if binary exists
if [ ! -f "$AGENT_BINARY" ]; then
    echo "❌ Agent binary not found. Please run: ./build-v1.sh"
    exit 1
fi

print_header "🚀 DUCLA CLOUD AGENT CLI DEMO"

print_step "Testing help command..."
$AGENT_BINARY --help | head -10
echo ""

print_step "Testing version command..."
$AGENT_BINARY show version
echo ""

print_step "Testing config display..."
$AGENT_BINARY show config
echo ""

print_step "Testing config validation..."
$AGENT_BINARY config validate
echo ""

print_step "Testing with custom config..."
$AGENT_BINARY -config agent.yaml show config | head -5
echo ""

print_step "Testing man page..."
echo "📚 Man page available:"
echo "  man ducla-agent"
echo ""
echo "Sample from man page:"
man ducla-agent | head -5
echo ""

print_header "🎉 CLI DEMO COMPLETED"

print_success "CLI commands are working!"
echo ""
echo "📋 Available commands:"
echo "  ducla-agent --help           # Show help"
echo "  ducla-agent show version     # Show version"
echo "  ducla-agent show config      # Show config"
echo "  ducla-agent config validate  # Validate config"
echo "  man ducla-agent              # Show manual"
echo ""
echo "🚀 To test with running agent:"
echo "  1. Start agent: ducla-agent"
echo "  2. In another terminal:"
echo "     ducla-agent show status"
echo "     ducla-agent show health"
echo "     ducla-agent show tasks"
echo "     ducla-agent file list /tmp"