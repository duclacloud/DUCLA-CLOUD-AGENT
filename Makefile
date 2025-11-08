.PHONY: help build docker-build docker-push k8s-deploy k8s-delete clean

# Variables
IMAGE_NAME ?= ducla-cloud-agent
IMAGE_TAG ?= latest
REGISTRY ?= docker.io/your-registry
NAMESPACE ?= ducla-system

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the Go binary
	go build -o bin/ducla-agent ./cmd/agent

test: ## Run tests
	go test -v ./...

docker-build: ## Build Docker image
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

docker-push: ## Push Docker image to registry
	docker tag $(IMAGE_NAME):$(IMAGE_TAG) $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)
	docker push $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)

docker-run: ## Run Docker container locally
	docker run -d --name ducla-agent \
		-p 8080:8080 -p 8081:8081 -p 8443:8443 -p 9090:9090 \
		--env-file .env \
		$(IMAGE_NAME):$(IMAGE_TAG)

docker-stop: ## Stop Docker container
	docker stop ducla-agent && docker rm ducla-agent

compose-up: ## Start services with Docker Compose
	docker-compose up -d

compose-down: ## Stop services with Docker Compose
	docker-compose down

compose-logs: ## View Docker Compose logs
	docker-compose logs -f

k8s-deploy: ## Deploy to Kubernetes
	kubectl apply -k k8s/

k8s-delete: ## Delete from Kubernetes
	kubectl delete -k k8s/

k8s-logs: ## View Kubernetes logs
	kubectl logs -n $(NAMESPACE) -l app=ducla-agent -f

k8s-status: ## Check Kubernetes deployment status
	kubectl get all -n $(NAMESPACE)

clean: ## Clean build artifacts
	rm -rf bin/
	go clean

.DEFAULT_GOAL := help
