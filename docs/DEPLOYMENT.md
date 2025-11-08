# Deployment Guide

## Docker Deployment

### Prerequisites
- Docker 20.10+
- Docker Compose 2.0+

### Quick Start with Docker Compose

1. Copy environment file:
```bash
cp .env.example .env
```

2. Edit `.env` with your configuration:
```bash
DUCLA_MASTER_URL=https://master.ducla.cloud
DUCLA_AGENT_TOKEN=your-token-here
DUCLA_JWT_SECRET=your-secret-here
```

3. Start the agent:
```bash
docker-compose up -d
```

4. Check logs:
```bash
docker-compose logs -f ducla-agent
```

5. Check health:
```bash
curl http://localhost:8081/health
```

### Development Mode

For development with hot reload:
```bash
docker-compose -f docker-compose.dev.yml up
```

### Build and Run Manually

Build the image:
```bash
docker build -t ducla-cloud-agent:latest .
```

Run the container:
```bash
docker run -d \
  --name ducla-agent \
  -p 8080:8080 \
  -p 8081:8081 \
  -p 8443:8443 \
  -p 9090:9090 \
  -e DUCLA_MASTER_URL=https://master.ducla.cloud \
  -e DUCLA_AGENT_TOKEN=your-token \
  -e DUCLA_JWT_SECRET=your-secret \
  ducla-cloud-agent:latest
```

## Kubernetes Deployment

### Prerequisites
- Kubernetes 1.24+
- kubectl configured
- Kustomize (optional, included in kubectl 1.14+)

### Quick Start

1. Create namespace:
```bash
kubectl apply -f k8s/namespace.yaml
```

2. Update secrets in `k8s/secret.yaml`:
```yaml
stringData:
  DUCLA_AGENT_TOKEN: "your-actual-token"
  DUCLA_JWT_SECRET: "your-actual-secret"
  DUCLA_MASTER_URL: "https://master.ducla.cloud"
```

3. Deploy all resources:
```bash
kubectl apply -k k8s/
```

4. Check deployment status:
```bash
kubectl get all -n ducla-system
```

5. View logs:
```bash
kubectl logs -n ducla-system -l app=ducla-agent -f
```

### Using Kustomize

Deploy with custom configuration:
```bash
kubectl kustomize k8s/ | kubectl apply -f -
```

### Scaling

Manual scaling:
```bash
kubectl scale deployment ducla-agent -n ducla-system --replicas=5
```

Auto-scaling is configured via HPA (Horizontal Pod Autoscaler) in `k8s/hpa.yaml`.

### Monitoring

The agent exposes Prometheus metrics on port 9090. If you have Prometheus Operator installed:
```bash
kubectl apply -f k8s/servicemonitor.yaml
```

### Ingress

To expose the agent externally, update `k8s/ingress.yaml` with your domain and apply:
```bash
kubectl apply -f k8s/ingress.yaml
```

### Updating Configuration

Update ConfigMap:
```bash
kubectl edit configmap ducla-agent-config -n ducla-system
```

Restart pods to apply changes:
```bash
kubectl rollout restart deployment ducla-agent -n ducla-system
```

### Cleanup

Remove all resources:
```bash
kubectl delete -k k8s/
```

## Using Makefile

The project includes a Makefile for common operations:

```bash
# Build binary
make build

# Build Docker image
make docker-build

# Push to registry
make docker-push REGISTRY=your-registry.io

# Deploy to Kubernetes
make k8s-deploy

# View logs
make k8s-logs

# Check status
make k8s-status

# Clean up
make k8s-delete
```

## Production Considerations

### Security
- Use proper secrets management (e.g., Sealed Secrets, External Secrets Operator)
- Enable TLS for gRPC and HTTP endpoints
- Configure network policies
- Use read-only root filesystem
- Run as non-root user (already configured)

### High Availability
- Deploy at least 3 replicas
- Configure pod anti-affinity (already configured)
- Use PodDisruptionBudget (already configured)
- Enable HPA for auto-scaling

### Monitoring
- Enable Prometheus metrics collection
- Set up alerting rules
- Monitor health endpoints
- Track resource usage

### Storage
- Use persistent volumes for data that needs to survive pod restarts
- Configure backup strategies
- Monitor disk usage

### Networking
- Configure proper ingress/load balancer
- Set up service mesh if needed
- Configure network policies for security
- Use proper DNS configuration
