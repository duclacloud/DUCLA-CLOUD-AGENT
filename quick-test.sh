#!/bin/bash

# Quick Test Script for Ducla Agent Monitoring
# This is a simplified version for quick verification

set -e

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_status() {
    echo -e "${BLUE}🔍 $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

test_endpoint() {
    local url="$1"
    local name="$2"
    local expected_code="${3:-200}"
    
    print_status "Testing $name..."
    
    if response=$(curl -s -w "%{http_code}" "$url" 2>/dev/null); then
        http_code="${response: -3}"
        if [ "$http_code" = "$expected_code" ]; then
            print_success "$name is working (HTTP $http_code)"
            return 0
        else
            print_error "$name returned HTTP $http_code (expected $expected_code)"
            return 1
        fi
    else
        print_error "$name is not accessible"
        return 1
    fi
}

main() {
    echo -e "${BLUE}"
    echo "🚀 Ducla Agent Quick Test"
    echo "========================"
    echo -e "${NC}"
    
    # Test all endpoints
    test_endpoint "http://localhost:8081/health" "Ducla Agent Health"
    test_endpoint "http://localhost:8080/api/v1/status" "Ducla Agent API"
    test_endpoint "http://localhost:9090/metrics" "Ducla Agent Metrics"
    test_endpoint "http://localhost:9091/-/healthy" "Prometheus"
    test_endpoint "http://localhost:3000/api/health" "Grafana"
    test_endpoint "http://localhost:9093/-/healthy" "AlertManager"
    
    echo ""
    print_status "Testing Prometheus integration..."
    
    # Test if Prometheus can see Ducla Agent
    if curl -s "http://localhost:9091/api/v1/query?query=up{job=\"ducla-agent\"}" | grep -q '"value":\[.*,"1"\]'; then
        print_success "Prometheus is scraping Ducla Agent successfully"
    else
        print_error "Prometheus cannot scrape Ducla Agent"
    fi
    
    echo ""
    print_status "Checking container status..."
    
    # Check if all containers are running
    running_containers=$(docker-compose -f docker-compose.monitoring.yml ps | grep "Up" | wc -l)
    total_containers=6
    
    if [ "$running_containers" -eq "$total_containers" ]; then
        print_success "All $total_containers containers are running"
    else
        print_error "Only $running_containers/$total_containers containers are running"
        echo "Run: docker-compose -f docker-compose.monitoring.yml ps"
    fi
    
    echo ""
    echo -e "${BLUE}📊 Quick Access URLs:${NC}"
    echo "  • Ducla Agent API:    http://localhost:8080"
    echo "  • Ducla Agent Health: http://localhost:8081/health"
    echo "  • Prometheus:         http://localhost:9091"
    echo "  • Grafana:            http://localhost:3000 (admin/admin123)"
    echo "  • AlertManager:       http://localhost:9093"
    
    echo ""
    echo -e "${BLUE}🔧 Quick Commands:${NC}"
    echo "  • View logs:          docker-compose -f docker-compose.monitoring.yml logs -f"
    echo "  • Restart agent:      docker-compose -f docker-compose.monitoring.yml restart ducla-agent"
    echo "  • Stop stack:         ./stop-monitoring.sh"
    echo "  • Full test:          ./run-monitoring-tests.sh"
    
    echo ""
    print_success "Quick test completed! 🎉"
}

main "$@"