#!/bin/bash

# Install man page for Ducla Cloud Agent

set -e

echo "📚 Installing man page for Ducla Cloud Agent..."

# Check if man page exists
if [ ! -f "docs/ducla-agent.1" ]; then
    echo "❌ Man page not found: docs/ducla-agent.1"
    exit 1
fi

# Create man directory if it doesn't exist
sudo mkdir -p /usr/local/man/man1

# Copy man page
sudo cp docs/ducla-agent.1 /usr/local/man/man1/

# Update man database
sudo mandb

echo "✅ Man page installed successfully!"
echo ""
echo "📖 You can now use:"
echo "  man ducla-agent"
echo ""
echo "🔍 Test it:"
man ducla-agent