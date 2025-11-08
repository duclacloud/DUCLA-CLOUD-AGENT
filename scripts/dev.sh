#!/bin/bash

# Development helper script

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

case "$1" in
    run)
        echo -e "${GREEN}Running agent in development mode...${NC}"
        go run ./cmd/agent --config configs/agent.yaml
        ;;
    
    build)
        echo -e "${GREEN}Building agent...${NC}"
        go build -o bin/ducla-agent ./cmd/agent
        ;;
    
    test)
        echo -e "${GREEN}Running tests...${NC}"
        go test -v -race -coverprofile=coverage.out ./...
        ;;
    
    coverage)
        echo -e "${GREEN}Generating coverage report...${NC}"
        go test -coverprofile=coverage.out ./...
        go tool cover -html=coverage.out -o coverage.html
        echo "Coverage report: coverage.html"
        ;;
    
    lint)
        echo -e "${GREEN}Running linter...${NC}"
        if command -v golangci-lint &> /dev/null; then
            golangci-lint run
        else
            echo -e "${YELLOW}golangci-lint not found, install it from https://golangci-lint.run/${NC}"
        fi
        ;;
    
    fmt)
        echo -e "${GREEN}Formatting code...${NC}"
        go fmt ./...
        ;;
    
    tidy)
        echo -e "${GREEN}Tidying dependencies...${NC}"
        go mod tidy
        ;;
    
    docker)
        echo -e "${GREEN}Building Docker image...${NC}"
        docker build -t ducla-cloud-agent:dev .
        ;;
    
    compose)
        echo -e "${GREEN}Starting with Docker Compose...${NC}"
        docker-compose -f docker-compose.dev.yml up
        ;;
    
    clean)
        echo -e "${GREEN}Cleaning build artifacts...${NC}"
        rm -rf bin/ dist/ coverage.out coverage.html
        go clean
        ;;
    
    *)
        echo "Usage: $0 {run|build|test|coverage|lint|fmt|tidy|docker|compose|clean}"
        echo ""
        echo "Commands:"
        echo "  run       - Run agent in development mode"
        echo "  build     - Build binary"
        echo "  test      - Run tests"
        echo "  coverage  - Generate coverage report"
        echo "  lint      - Run linter"
        echo "  fmt       - Format code"
        echo "  tidy      - Tidy dependencies"
        echo "  docker    - Build Docker image"
        echo "  compose   - Start with Docker Compose"
        echo "  clean     - Clean build artifacts"
        exit 1
        ;;
esac
