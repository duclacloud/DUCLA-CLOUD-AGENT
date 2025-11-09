#!/bin/bash

# Automated Monitoring Test Suite for Ducla Agent
# This script runs all test cases automatically and generates a report

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test results
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0
TEST_RESULTS=()

# Logging
LOG_FILE="monitoring-test-$(date +%Y%m%d_%H%M%S).log"

log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" | tee -a "$LOG_FILE"
}

print_header() {
    echo -e "\n${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}\n"
    log "TEST: $1"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
    log "SUCCESS: $1"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
    log "ERROR: $1"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
    log "WARNING: $1"
}

run_test() {
    local test_name="$1"
    local test_command="$2"
    local expected_result="$3"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    echo -e "\n${YELLOW}Running: $test_name${NC}"
    
    if eval "$test_command"; then
        print_success "$test_name"
        TEST_RESULTS+=("✅ $test_name")
        PASSED_TESTS=$((PASSED_TESTS + 1))
        return 0
    else
        print_error "$test_name"
        TEST_RESULTS+=("❌ $test_name")
        FAILED_TESTS=$((FAILED_TESTS + 1))
        return 1
    fi
}

wait_for_service() {
    local service_url="$1"
    local service_name="$2"
    local max_attempts=30
    local attempt=1
    
    echo "Waiting for $service_name to be ready..."
    while [ $attempt -le $max_attempts ]; do
        if curl -s "$service_url" > /dev/null 2>&1; then
            print_success "$service_name is ready"
            return 0
        fi
        echo "Attempt $attempt/$max_attempts - $service_name not ready yet..."
        sleep 2
        attempt=$((attempt + 1))
    done
    
    print_error "$service_name failed to start within $((max_attempts * 2)) seconds"
    return 1
}

# Main test execution
main() {
    print_header "🧪 DUCLA AGENT MONITORING TEST SUITE"
    
    log "Starting monitoring test suite"
    log "Test log file: $LOG_FILE"
    
    # Pre-test checks
    print_header "📋 PRE-TEST CHECKS"
    
    run_test "Docker is running" \
        "docker info > /dev/null 2>&1" \
        "Docker should be running"
    
    run_test "Docker Compose is available" \
        "command -v docker-compose > /dev/null 2>&1" \
        "Docker Compose should be installed"
    
    run_test "Sufficient disk space (>2GB)" \
        "[ $(df . | tail -1 | awk '{print $4}') -gt 2000000 ]" \
        "At least 2GB free space required"
    
    # Start monitoring stack
    print_header "🚀 STARTING MONITORING STACK"
    
    echo "Starting monitoring stack..."
    if ./start-monitoring.sh > /dev/null 2>&1; then
        print_success "Monitoring stack started"
    else
        print_error "Failed to start monitoring stack"
        exit 1
    fi
    
    # Wait for services
    print_header "⏳ WAITING FOR SERVICES"
    
    wait_for_service "http://localhost:8081/health" "Ducla Agent Health"
    wait_for_service "http://localhost:8080/api/v1/status" "Ducla Agent API"
    wait_for_service "http://localhost:9091/-/healthy" "Prometheus"
    wait_for_service "http://localhost:3000/api/health" "Grafana"
    
    # Test Case 1: Health Checks
    print_header "🏥 TEST CASE 1: HEALTH CHECKS"
    
    run_test "Ducla Agent health endpoint" \
        "curl -s -o /dev/null -w '%{http_code}' http://localhost:8081/health | grep -q '200'" \
        "Health endpoint should return HTTP 200"
    
    run_test "Ducla Agent API status" \
        "curl -s http://localhost:8080/api/v1/status | jq -r '.success' | grep -q 'true'" \
        "API should return success: true"
    
    run_test "Ducla Agent metrics endpoint" \
        "curl -s http://localhost:9090/metrics | grep -q 'ducla_'" \
        "Metrics endpoint should return Ducla metrics"
    
    # Test Case 2: Prometheus Integration
    print_header "📊 TEST CASE 2: PROMETHEUS INTEGRATION"
    
    run_test "Prometheus targets Ducla Agent" \
        "curl -s http://localhost:9091/api/v1/targets | jq -r '.data.activeTargets[] | select(.job==\"ducla-agent\") | .health' | grep -q 'up'" \
        "Prometheus should show Ducla Agent target as up"
    
    run_test "Prometheus can query Ducla metrics" \
        "curl -s 'http://localhost:9091/api/v1/query?query=up{job=\"ducla-agent\"}' | jq -r '.data.result[0].value[1]' | grep -q '1'" \
        "Prometheus should return up=1 for Ducla Agent"
    
    run_test "Ducla-specific metrics available" \
        "curl -s 'http://localhost:9091/api/v1/query?query=ducla_agent_info' | jq -r '.data.result | length' | grep -v '^0$'" \
        "Ducla-specific metrics should be available"
    
    # Test Case 3: Load Testing
    print_header "🔥 TEST CASE 3: LOAD TESTING"
    
    echo "Generating load on Ducla Agent..."
    for i in {1..20}; do
        curl -s http://localhost:8080/api/v1/status > /dev/null &
    done
    wait
    
    sleep 5  # Wait for metrics to update
    
    run_test "API request metrics increase" \
        "curl -s 'http://localhost:9091/api/v1/query?query=ducla_http_requests_total' | jq -r '.data.result | length' | grep -v '^0$'" \
        "HTTP request metrics should be collected"
    
    run_test "Memory usage metrics available" \
        "curl -s 'http://localhost:9091/api/v1/query?query=process_resident_memory_bytes{job=\"ducla-agent\"}' | jq -r '.data.result[0].value[1]' | grep -E '^[0-9]+$'" \
        "Memory usage metrics should be numeric"
    
    # Test Case 4: Alert Testing
    print_header "🚨 TEST CASE 4: ALERT TESTING"
    
    echo "Testing alerts by stopping Ducla Agent..."
    docker-compose -f docker-compose.monitoring.yml stop ducla-agent > /dev/null 2>&1
    
    sleep 70  # Wait for alert to fire (1 minute + buffer)
    
    run_test "Alert fires when agent is down" \
        "curl -s http://localhost:9091/api/v1/alerts | jq -r '.data.alerts[] | select(.labels.alertname==\"DuclaAgentDown\") | .state' | grep -q 'firing'" \
        "DuclaAgentDown alert should be firing"
    
    echo "Restarting Ducla Agent..."
    docker-compose -f docker-compose.monitoring.yml start ducla-agent > /dev/null 2>&1
    
    wait_for_service "http://localhost:8081/health" "Ducla Agent Health"
    
    # Test Case 5: Performance Benchmarks
    print_header "⚡ TEST CASE 5: PERFORMANCE BENCHMARKS"
    
    run_test "API response time < 100ms" \
        "timeout 5s bash -c 'time curl -s http://localhost:8080/api/v1/status > /dev/null' 2>&1 | grep real | awk '{print $2}' | sed 's/[^0-9.]//g' | awk '{print ($1 < 0.1)}' | grep -q '1'" \
        "API should respond within 100ms"
    
    run_test "Metrics endpoint response time < 200ms" \
        "timeout 5s bash -c 'time curl -s http://localhost:9090/metrics > /dev/null' 2>&1 | grep real | awk '{print $2}' | sed 's/[^0-9.]//g' | awk '{print ($1 < 0.2)}' | grep -q '1'" \
        "Metrics endpoint should respond within 200ms"
    
    run_test "Ducla Agent memory usage < 200MB" \
        "docker stats --no-stream --format 'table {{.MemUsage}}' ducla-agent | tail -1 | sed 's/MiB.*//' | awk '{print ($1 < 200)}' | grep -q '1'" \
        "Ducla Agent should use less than 200MB memory"
    
    # Test Case 6: Container Health
    print_header "🐳 TEST CASE 6: CONTAINER HEALTH"
    
    run_test "All containers are running" \
        "[ $(docker-compose -f docker-compose.monitoring.yml ps | grep 'Up' | wc -l) -eq 6 ]" \
        "All 6 containers should be running"
    
    run_test "No container restarts" \
        "docker-compose -f docker-compose.monitoring.yml ps | grep -v 'Up.*seconds' | grep -q 'Up'" \
        "Containers should be stable (not recently restarted)"
    
    # Generate final report
    print_header "📊 TEST RESULTS SUMMARY"
    
    echo -e "\n${BLUE}Test Execution Summary:${NC}"
    echo "Total Tests: $TOTAL_TESTS"
    echo -e "Passed: ${GREEN}$PASSED_TESTS${NC}"
    echo -e "Failed: ${RED}$FAILED_TESTS${NC}"
    echo -e "Success Rate: $(( PASSED_TESTS * 100 / TOTAL_TESTS ))%"
    
    echo -e "\n${BLUE}Detailed Results:${NC}"
    for result in "${TEST_RESULTS[@]}"; do
        echo "$result"
    done
    
    # Resource usage summary
    echo -e "\n${BLUE}Resource Usage Summary:${NC}"
    docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.MemPerc}}"
    
    # Log summary
    log "Test completed. Total: $TOTAL_TESTS, Passed: $PASSED_TESTS, Failed: $FAILED_TESTS"
    
    echo -e "\n${BLUE}Full test log saved to: $LOG_FILE${NC}"
    
    # Cleanup prompt
    echo -e "\n${YELLOW}Test completed. Do you want to stop the monitoring stack? (y/N)${NC}"
    read -r response
    if [[ "$response" =~ ^[Yy]$ ]]; then
        echo "Stopping monitoring stack..."
        ./stop-monitoring.sh
        print_success "Monitoring stack stopped"
    else
        echo -e "\n${BLUE}Monitoring stack is still running:${NC}"
        echo "  • Grafana: http://localhost:3000 (admin/admin123)"
        echo "  • Prometheus: http://localhost:9091"
        echo "  • Ducla Agent: http://localhost:8080"
        echo ""
        echo "To stop later, run: ./stop-monitoring.sh"
    fi
    
    # Exit with appropriate code
    if [ $FAILED_TESTS -eq 0 ]; then
        echo -e "\n${GREEN}🎉 All tests passed! Ducla Agent monitoring integration is working perfectly.${NC}"
        exit 0
    else
        echo -e "\n${RED}❌ Some tests failed. Please check the logs and fix issues.${NC}"
        exit 1
    fi
}

# Trap to ensure cleanup on script exit
trap 'echo -e "\n${YELLOW}Test interrupted. You may need to run ./stop-monitoring.sh to clean up.${NC}"' INT TERM

# Run main function
main "$@"