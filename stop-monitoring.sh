#!/bin/bash

# Stop Ducla Agent Monitoring Stack

set -e

echo "🛑 Stopping Ducla Agent Monitoring Stack"
echo "========================================"

# Stop containers
echo "🐳 Stopping containers..."
docker-compose -f docker-compose.monitoring.yml down

# Optional: Remove volumes (uncomment if you want to clean data)
# echo "🗑️  Removing volumes..."
# docker-compose -f docker-compose.monitoring.yml down -v

echo ""
echo "✅ Monitoring Stack Stopped Successfully!"
echo ""
echo "🔧 To start again:"
echo "  ./start-monitoring.sh"
echo ""
echo "🗑️  To remove all data:"
echo "  docker-compose -f docker-compose.monitoring.yml down -v"