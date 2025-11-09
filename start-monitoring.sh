#!/bin/bash

# Start Ducla Agent Monitoring Stack
# This script starts Prometheus, Grafana, and Ducla Agent for testing

set -e

echo "🚀 Starting Ducla Agent Monitoring Stack"
echo "========================================"

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    echo "❌ docker-compose not found. Please install docker-compose."
    exit 1
fi

# Create necessary directories
echo "📁 Creating directories..."
mkdir -p data logs
sudo chown -R 1000:1000 data logs

# Build Ducla Agent image if needed
if ! docker images | grep -q ducla-agent; then
    echo "🔨 Building Ducla Agent Docker image..."
    docker build -t ducla-agent .
fi

# Start the monitoring stack
echo "🐳 Starting containers..."
docker-compose -f docker-compose.monitoring.yml up -d

# Wait for services to be ready
echo "⏳ Waiting for services to start..."
sleep 30

# Check service health
echo "🔍 Checking service health..."

# Check Ducla Agent
if curl -s http://localhost:8081/health > /dev/null; then
    echo "✅ Ducla Agent: http://localhost:8080 (API), http://localhost:8081 (Health)"
else
    echo "❌ Ducla Agent health check failed"
fi

# Check Prometheus
if curl -s http://localhost:9091/-/healthy > /dev/null; then
    echo "✅ Prometheus: http://localhost:9091"
else
    echo "❌ Prometheus health check failed"
fi

# Check Grafana
if curl -s http://localhost:3000/api/health > /dev/null; then
    echo "✅ Grafana: http://localhost:3000 (admin/admin123)"
else
    echo "❌ Grafana health check failed"
fi

echo ""
echo "🎉 Monitoring Stack Started Successfully!"
echo "========================================"
echo ""
echo "📊 Access URLs:"
echo "  • Ducla Agent API:    http://localhost:8080"
echo "  • Ducla Agent Health: http://localhost:8081/health"
echo "  • Ducla Agent Metrics: http://localhost:9090/metrics"
echo "  • Prometheus:         http://localhost:9091"
echo "  • Grafana:            http://localhost:3000 (admin/admin123)"
echo "  • AlertManager:       http://localhost:9093"
echo "  • Node Exporter:      http://localhost:9100"
echo "  • cAdvisor:           http://localhost:8082"
echo ""
echo "🔧 Useful Commands:"
echo "  • View logs:          docker-compose -f docker-compose.monitoring.yml logs -f"
echo "  • Stop stack:         docker-compose -f docker-compose.monitoring.yml down"
echo "  • Restart service:    docker-compose -f docker-compose.monitoring.yml restart ducla-agent"
echo ""
echo "📈 Test Metrics:"
echo "  curl http://localhost:9090/metrics"
echo "  curl http://localhost:8080/api/v1/status"
echo ""
echo "🎯 Next Steps:"
echo "  1. Open Grafana: http://localhost:3000"
echo "  2. Login with admin/admin123"
echo "  3. Check 'Ducla Cloud Agent Dashboard'"
echo "  4. Test API endpoints and watch metrics"