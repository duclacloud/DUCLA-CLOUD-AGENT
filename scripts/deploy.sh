#!/bin/bash

# Deployment script for Ducla Cloud Agent

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

RELEASE_DIR="release"
DEPLOY_TARGET=${1:-"local"}

echo -e "${GREEN}🚀 Deploying Ducla Cloud Agent${NC}"
echo -e "${BLUE}Target: ${DEPLOY_TARGET}${NC}"
echo ""

case $DEPLOY_TARGET in
    "local")
        echo -e "${YELLOW}📦 Local deployment - Testing installation...${NC}"
        
        # Check if release files exist
        if [ ! -d "$RELEASE_DIR" ]; then
            echo -e "${RED}❌ Release directory not found. Run ./scripts/distribute.sh first${NC}"
            exit 1
        fi
        
        # Test installation
        echo -e "${YELLOW}🧪 Testing installation script...${NC}"
        cd $RELEASE_DIR
        
        # Dry run first
        echo -e "${BLUE}📋 Installation script preview:${NC}"
        head -20 install.sh
        echo "..."
        
        echo ""
        echo -e "${YELLOW}⚠️  Ready to install locally. This will:${NC}"
        echo "  • Install ducla-agent service"
        echo "  • Create ducla user"
        echo "  • Configure systemd service"
        echo "  • Start the agent"
        echo ""
        read -p "Continue with local installation? (y/N): " -n 1 -r
        echo
        
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            echo -e "${YELLOW}🔧 Installing locally...${NC}"
            sudo ./install.sh
        else
            echo -e "${BLUE}ℹ️  Installation cancelled${NC}"
        fi
        ;;
        
    "github")
        echo -e "${YELLOW}📤 GitHub Release deployment...${NC}"
        
        # Check if gh CLI is installed
        if ! command -v gh &> /dev/null; then
            echo -e "${RED}❌ GitHub CLI (gh) not found. Install it first:${NC}"
            echo "  curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg"
            echo "  sudo apt update && sudo apt install gh"
            exit 1
        fi
        
        # Check if authenticated
        if ! gh auth status &> /dev/null; then
            echo -e "${RED}❌ Not authenticated with GitHub. Run: gh auth login${NC}"
            exit 1
        fi
        
        VERSION=$(git describe --tags --always --dirty 2>/dev/null | sed 's/^v//' || echo "1.0.0")
        TAG="v${VERSION}"
        
        echo -e "${YELLOW}🏷️  Creating GitHub release ${TAG}...${NC}"
        
        # Create release
        gh release create $TAG $RELEASE_DIR/* \
            --title "Ducla Cloud Agent $TAG" \
            --notes-file $RELEASE_DIR/README.md \
            --draft
            
        echo -e "${GREEN}✅ GitHub release created (draft)${NC}"
        echo "Edit and publish at: https://github.com/$(gh repo view --json owner,name -q '.owner.login + "/" + .name')/releases"
        ;;
        
    "docker")
        echo -e "${YELLOW}🐳 Docker deployment...${NC}"
        
        # Create Dockerfile if not exists
        if [ ! -f "Dockerfile" ]; then
            echo -e "${YELLOW}📝 Creating Dockerfile...${NC}"
            cat > Dockerfile <<'EOF'
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 ducla && \
    adduser -D -s /bin/sh -u 1000 -G ducla ducla

# Create directories
RUN mkdir -p /opt/ducla/data /var/log/ducla && \
    chown -R ducla:ducla /opt/ducla /var/log/ducla

# Copy binary
COPY release/ducla-agent-linux-amd64 /usr/local/bin/ducla-agent
RUN chmod +x /usr/local/bin/ducla-agent

# Copy config
COPY configs/agent.yaml /etc/ducla/agent.yaml

# Switch to non-root user
USER ducla

# Expose ports
EXPOSE 8080 8081 9090

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8081/health || exit 1

# Run agent
CMD ["/usr/local/bin/ducla-agent", "--config", "/etc/ducla/agent.yaml"]
EOF
        fi
        
        # Build Docker image
        VERSION=$(git describe --tags --always --dirty 2>/dev/null | sed 's/^v//' || echo "1.0.0")
        IMAGE_NAME="ducla/cloud-agent:${VERSION}"
        
        echo -e "${YELLOW}🔨 Building Docker image ${IMAGE_NAME}...${NC}"
        docker build -t $IMAGE_NAME .
        docker tag $IMAGE_NAME ducla/cloud-agent:latest
        
        echo -e "${GREEN}✅ Docker image built successfully${NC}"
        echo "Run with: docker run -d -p 8080:8080 -p 8081:8081 -p 9090:9090 $IMAGE_NAME"
        ;;
        
    "kubernetes")
        echo -e "${YELLOW}☸️  Kubernetes deployment...${NC}"
        
        # Create k8s manifests if not exist
        mkdir -p k8s
        
        if [ ! -f "k8s/deployment.yaml" ]; then
            echo -e "${YELLOW}📝 Creating Kubernetes manifests...${NC}"
            
            # Deployment
            cat > k8s/deployment.yaml <<'EOF'
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ducla-agent
  labels:
    app: ducla-agent
spec:
  replicas: 3
  selector:
    matchLabels:
      app: ducla-agent
  template:
    metadata:
      labels:
        app: ducla-agent
    spec:
      containers:
      - name: ducla-agent
        image: ducla/cloud-agent:latest
        ports:
        - containerPort: 8080
          name: http
        - containerPort: 8081
          name: health
        - containerPort: 9090
          name: metrics
        env:
        - name: HOSTNAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        livenessProbe:
          httpGet:
            path: /health
            port: 8081
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8081
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
EOF
            
            # Service
            cat > k8s/service.yaml <<'EOF'
apiVersion: v1
kind: Service
metadata:
  name: ducla-agent
  labels:
    app: ducla-agent
spec:
  selector:
    app: ducla-agent
  ports:
  - name: http
    port: 8080
    targetPort: 8080
  - name: health
    port: 8081
    targetPort: 8081
  - name: metrics
    port: 9090
    targetPort: 9090
  type: ClusterIP
EOF
            
            # ConfigMap
            cat > k8s/configmap.yaml <<'EOF'
apiVersion: v1
kind: ConfigMap
metadata:
  name: ducla-agent-config
data:
  agent.yaml: |
    agent:
      id: "${HOSTNAME}"
      name: "${HOSTNAME}"
      environment: "kubernetes"
    
    api:
      http:
        enabled: true
        port: 8080
    
    health:
      enabled: true
      port: 8081
    
    metrics:
      enabled: true
      port: 9090
    
    logging:
      level: "info"
      format: "json"
      output: "stdout"
EOF
        fi
        
        echo -e "${YELLOW}🚀 Deploying to Kubernetes...${NC}"
        kubectl apply -f k8s/
        
        echo -e "${GREEN}✅ Kubernetes deployment complete${NC}"
        echo "Check status: kubectl get pods -l app=ducla-agent"
        ;;
        
    *)
        echo -e "${RED}❌ Unknown deployment target: $DEPLOY_TARGET${NC}"
        echo ""
        echo -e "${BLUE}Available targets:${NC}"
        echo "  • local     - Test installation locally"
        echo "  • github    - Create GitHub release"
        echo "  • docker    - Build Docker image"
        echo "  • kubernetes - Deploy to Kubernetes"
        echo ""
        echo "Usage: $0 <target>"
        exit 1
        ;;
esac

echo ""
echo -e "${GREEN}🎉 Deployment complete!${NC}"