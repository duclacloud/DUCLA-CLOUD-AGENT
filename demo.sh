#!/bin/bash

# Ducla Cloud Agent v1.0.0 Demo Script
# This script demonstrates the key features of the agent

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Demo configuration
AGENT_BINARY="./dist/ducla-agent"
CONFIG_FILE="demo-config.yaml"
DEMO_DIR="demo-workspace"

# Helper functions
print_header() {
    echo -e "\n${PURPLE}========================================${NC}"
    echo -e "${PURPLE}$1${NC}"
    echo -e "${PURPLE}========================================${NC}\n"
}

print_step() {
    echo -e "${CYAN}➤ $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

wait_for_input() {
    echo -e "\n${YELLOW}Press Enter to continue...${NC}"
    read
}

# Check if agent binary exists
check_binary() {
    if [ ! -f "$AGENT_BINARY" ]; then
        echo -e "${RED}❌ Agent binary not found at $AGENT_BINARY${NC}"
        echo -e "${YELLOW}Please run: ./build-v1.sh${NC}"
        exit 1
    fi
}

# Create demo configuration
create_demo_config() {
    print_step "Creating demo configuration..."
    
    cat > $CONFIG_FILE << EOF
# Demo Configuration for Ducla Cloud Agent v1.0.0
agent:
  id: "demo-agent"
  name: "Demo Agent"
  environment: "demo"
  region: "local"
  capabilities:
    - "file_operations"
    - "task_execution"
    - "system_monitoring"

# No master server for demo
master:
  url: ""
  token: ""
  max_reconnect_attempts: 0

api:
  http:
    enabled: true
    address: "127.0.0.1"
    port: 8080
    tls:
      enabled: false

security:
  auth:
    enabled: false
  rbac:
    enabled: false

storage:
  data_dir: "./$DEMO_DIR/data"
  temp_dir: "./$DEMO_DIR/tmp"
  cleanup:
    enabled: true
    interval: 30s
    max_age: 5m

logging:
  level: "info"
  format: "text"
  output: "stdout"

metrics:
  enabled: true
  address: "127.0.0.1"
  port: 9090

health:
  enabled: true
  address: "127.0.0.1"
  port: 8081

executor:
  workers: 3
  queue_size: 50
  task_timeout: 60s
EOF

    print_success "Demo configuration created: $CONFIG_FILE"
}

# Setup demo workspace
setup_workspace() {
    print_step "Setting up demo workspace..."
    
    mkdir -p $DEMO_DIR/{data,tmp,files}
    
    # Create some demo files
    echo "Hello from Ducla Agent!" > $DEMO_DIR/files/hello.txt
    echo "This is a test file" > $DEMO_DIR/files/test.txt
    echo "Demo data $(date)" > $DEMO_DIR/files/demo.log
    
    print_success "Demo workspace created: $DEMO_DIR/"
}

# Start agent in background
start_agent() {
    print_step "Starting Ducla Cloud Agent..."
    
    # Kill any existing agent process
    pkill -f "ducla-agent" 2>/dev/null || true
    sleep 1
    
    # Start agent in background
    $AGENT_BINARY -config $CONFIG_FILE > agent.log 2>&1 &
    AGENT_PID=$!
    
    # Wait for agent to start
    sleep 3
    
    # Check if agent is running
    if kill -0 $AGENT_PID 2>/dev/null; then
        print_success "Agent started successfully (PID: $AGENT_PID)"
    else
        echo -e "${RED}❌ Failed to start agent${NC}"
        cat agent.log
        exit 1
    fi
}

# Demo version information
demo_version() {
    print_header "🔍 VERSION INFORMATION"
    
    print_step "Checking agent version..."
    $AGENT_BINARY -version
    
    wait_for_input
}

# Demo health checks
demo_health() {
    print_header "🏥 HEALTH MONITORING"
    
    print_step "Checking agent health..."
    curl -s http://127.0.0.1:8081/health | jq '.' 2>/dev/null || curl -s http://127.0.0.1:8081/health
    
    echo -e "\n"
    print_step "Checking readiness..."
    curl -s http://127.0.0.1:8081/ready | jq '.' 2>/dev/null || curl -s http://127.0.0.1:8081/ready
    
    wait_for_input
}

# Demo metrics
demo_metrics() {
    print_header "📊 METRICS COLLECTION"
    
    print_step "Fetching Prometheus metrics..."
    echo -e "${BLUE}Sample metrics:${NC}"
    curl -s http://127.0.0.1:9090/metrics | head -20
    echo "..."
    
    echo -e "\n"
    print_info "Full metrics available at: http://127.0.0.1:9090/metrics"
    
    wait_for_input
}

# Demo agent status
demo_status() {
    print_header "📋 AGENT STATUS"
    
    print_step "Getting agent status..."
    curl -s http://127.0.0.1:8080/api/v1/status | jq '.' 2>/dev/null || curl -s http://127.0.0.1:8080/api/v1/status
    
    wait_for_input
}

# Demo task execution
demo_tasks() {
    print_header "⚡ TASK EXECUTION"
    
    print_step "Creating a simple task..."
    TASK_RESPONSE=$(curl -s -X POST http://127.0.0.1:8080/api/v1/tasks \
        -H "Content-Type: application/json" \
        -d '{
            "type": "shell",
            "name": "demo-task",
            "command": "echo",
            "args": ["Hello from Ducla Agent Task!"]
        }')
    
    echo "$TASK_RESPONSE" | jq '.' 2>/dev/null || echo "$TASK_RESPONSE"
    
    # Extract task ID
    TASK_ID=$(echo "$TASK_RESPONSE" | jq -r '.data.task_id' 2>/dev/null || echo "unknown")
    
    if [ "$TASK_ID" != "unknown" ] && [ "$TASK_ID" != "null" ]; then
        echo -e "\n"
        print_step "Getting task details..."
        sleep 1
        curl -s http://127.0.0.1:8080/api/v1/tasks/$TASK_ID | jq '.' 2>/dev/null || curl -s http://127.0.0.1:8080/api/v1/tasks/$TASK_ID
    fi
    
    echo -e "\n"
    print_step "Listing all tasks..."
    curl -s http://127.0.0.1:8080/api/v1/tasks | jq '.' 2>/dev/null || curl -s http://127.0.0.1:8080/api/v1/tasks
    
    wait_for_input
}

# Demo file operations
demo_files() {
    print_header "📁 FILE OPERATIONS"
    
    print_step "Listing files in demo directory..."
    curl -s -X POST http://127.0.0.1:8080/api/v1/files \
        -H "Content-Type: application/json" \
        -d "{
            \"type\": \"list\",
            \"source_path\": \"$(pwd)/$DEMO_DIR/files\"
        }" | jq '.' 2>/dev/null || curl -s -X POST http://127.0.0.1:8080/api/v1/files \
        -H "Content-Type: application/json" \
        -d "{
            \"type\": \"list\",
            \"source_path\": \"$(pwd)/$DEMO_DIR/files\"
        }"
    
    echo -e "\n"
    print_step "Copying a file..."
    curl -s -X POST http://127.0.0.1:8080/api/v1/files \
        -H "Content-Type: application/json" \
        -d "{
            \"type\": \"copy\",
            \"source_path\": \"$(pwd)/$DEMO_DIR/files/hello.txt\",
            \"dest_path\": \"$(pwd)/$DEMO_DIR/files/hello-copy.txt\"
        }" | jq '.' 2>/dev/null || curl -s -X POST http://127.0.0.1:8080/api/v1/files \
        -H "Content-Type: application/json" \
        -d "{
            \"type\": \"copy\",
            \"source_path\": \"$(pwd)/$DEMO_DIR/files/hello.txt\",
            \"dest_path\": \"$(pwd)/$DEMO_DIR/files/hello-copy.txt\"
        }"
    
    echo -e "\n"
    print_step "Verifying file was copied..."
    ls -la $DEMO_DIR/files/
    
    wait_for_input
}

# Demo API endpoints
demo_api() {
    print_header "🌐 API ENDPOINTS SUMMARY"
    
    print_info "Available API endpoints:"
    echo -e "${BLUE}HTTP REST API (Port 8080):${NC}"
    echo "  GET  http://127.0.0.1:8080/api/v1/status"
    echo "  GET  http://127.0.0.1:8080/api/v1/tasks"
    echo "  POST http://127.0.0.1:8080/api/v1/tasks"
    echo "  GET  http://127.0.0.1:8080/api/v1/tasks/{id}"
    echo "  POST http://127.0.0.1:8080/api/v1/files"
    
    echo -e "\n${BLUE}Health Check API (Port 8081):${NC}"
    echo "  GET  http://127.0.0.1:8081/health"
    echo "  GET  http://127.0.0.1:8081/ready"
    
    echo -e "\n${BLUE}Metrics API (Port 9090):${NC}"
    echo "  GET  http://127.0.0.1:9090/metrics"
    
    wait_for_input
}

# Stop agent
stop_agent() {
    print_step "Stopping agent..."
    
    if [ ! -z "$AGENT_PID" ] && kill -0 $AGENT_PID 2>/dev/null; then
        kill $AGENT_PID
        sleep 2
        print_success "Agent stopped"
    else
        print_warning "Agent was not running"
    fi
}

# Cleanup
cleanup() {
    print_step "Cleaning up demo files..."
    
    rm -f $CONFIG_FILE agent.log
    rm -rf $DEMO_DIR
    
    print_success "Cleanup completed"
}

# Main demo function
main() {
    print_header "🚀 DUCLA CLOUD AGENT v1.0.0 DEMO"
    
    print_info "This demo will showcase the key features of Ducla Cloud Agent:"
    echo "  • Version information"
    echo "  • Health monitoring"
    echo "  • Metrics collection"
    echo "  • Agent status"
    echo "  • Task execution"
    echo "  • File operations"
    echo "  • API endpoints"
    
    wait_for_input
    
    # Setup
    check_binary
    create_demo_config
    setup_workspace
    start_agent
    
    # Demo features
    demo_version
    demo_health
    demo_metrics
    demo_status
    demo_tasks
    demo_files
    demo_api
    
    # Cleanup
    stop_agent
    cleanup
    
    print_header "🎉 DEMO COMPLETED"
    print_success "Thank you for trying Ducla Cloud Agent v1.0.0!"
    print_info "For more information, see README-VI.md"
}

# Handle Ctrl+C
trap 'echo -e "\n${YELLOW}Demo interrupted${NC}"; stop_agent; cleanup; exit 1' INT

# Run demo
main