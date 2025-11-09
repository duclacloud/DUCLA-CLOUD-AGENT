# ☸️ K3s Monitoring Stack Alternative

## 🎯 When to Use K3s Instead of Docker

### Use K3s if:
- ✅ **Production-like testing**: Need to test Kubernetes-specific features
- ✅ **Helm charts**: Want to use Prometheus Operator, Grafana Operator
- ✅ **Service mesh**: Testing with Istio, Linkerd
- ✅ **Auto-scaling**: Testing HPA, VPA
- ✅ **Cloud-native**: Preparing for EKS, GKE, AKS deployment

### Use Docker if:
- ✅ **Quick testing**: Just want to test Ducla Agent integration
- ✅ **Development**: Local development và debugging
- ✅ **Resource limited**: Laptop với limited RAM/CPU
- ✅ **Simple setup**: Don't want Kubernetes complexity

## 🚀 K3s Setup (Alternative)

### 1. Install K3s
```bash
# Install K3s
curl -sfL https://get.k3s.io | sh -

# Check status
sudo k3s kubectl get nodes

# Copy kubeconfig
mkdir -p ~/.kube
sudo cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
sudo chown $USER:$USER ~/.kube/config
```

### 2. Install Prometheus Stack
```bash
# Add Helm repo
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

# Install kube-prometheus-stack
helm install monitoring prometheus-community/kube-prometheus-stack \
  --namespace monitoring \
  --create-namespace \
  --set prometheus.service.type=NodePort \
  --set prometheus.service.nodePort=30090 \
  --set grafana.service.type=NodePort \
  --set grafana.service.nodePort=30300 \
  --set alertmanager.service.type=NodePort \
  --set alertmanager.service.nodePort=30093
```

### 3. Deploy Ducla Agent
```yaml
# ducla-agent-k3s.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ducla-agent
  namespace: monitoring
spec:
  replicas: 1
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
        image: ducla-agent:latest
        ports:
        - containerPort: 8080
        - containerPort: 8081
        - containerPort: 9090
        env:
        - name: DUCLA_CONFIG_PATH
          value: "/etc/ducla/agent.yaml"
        volumeMounts:
        - name: config
          mountPath: /etc/ducla
      volumes:
      - name: config
        configMap:
          name: ducla-agent-config
---
apiVersion: v1
kind: Service
metadata:
  name: ducla-agent
  namespace: monitoring
  labels:
    app: ducla-agent
spec:
  selector:
    app: ducla-agent
  ports:
  - name: api
    port: 8080
    targetPort: 8080
    nodePort: 30080
  - name: health
    port: 8081
    targetPort: 8081
    nodePort: 30081
  - name: metrics
    port: 9090
    targetPort: 9090
    nodePort: 30090
  type: NodePort
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: ducla-agent-config
  namespace: monitoring
data:
  agent.yaml: |
    agent:
      id: "ducla-agent-k3s"
      name: "Ducla Cloud Agent K3s"
      environment: "testing"
    api:
      http:
        enabled: true
        address: "0.0.0.0"
        port: 8080
    health:
      enabled: true
      address: "0.0.0.0"
      port: 8081
    metrics:
      enabled: true
      address: "0.0.0.0"
      port: 9090
```

### 4. ServiceMonitor for Prometheus
```yaml
# ducla-agent-servicemonitor.yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: ducla-agent
  namespace: monitoring
  labels:
    app: ducla-agent
spec:
  selector:
    matchLabels:
      app: ducla-agent
  endpoints:
  - port: metrics
    interval: 30s
    path: /metrics
```

### 5. Deploy Commands
```bash
# Apply manifests
kubectl apply -f ducla-agent-k3s.yaml
kubectl apply -f ducla-agent-servicemonitor.yaml

# Check status
kubectl get pods -n monitoring
kubectl get svc -n monitoring

# Access services
echo "Ducla Agent API: http://localhost:30080"
echo "Ducla Agent Health: http://localhost:30081"
echo "Prometheus: http://localhost:30090"
echo "Grafana: http://localhost:30300"
```

## 📊 Comparison Summary

| Feature | Docker Compose | K3s |
|---------|----------------|-----|
| **Setup Time** | 5 minutes | 15-30 minutes |
| **Resource Usage** | ~500MB RAM | ~1GB+ RAM |
| **Learning Curve** | Easy | Medium |
| **Production-like** | No | Yes |
| **Debugging** | Easy | Medium |
| **Scaling** | Manual | Auto |
| **Service Discovery** | Manual | Automatic |
| **Monitoring Stack** | Custom | Prometheus Operator |

## 🎯 Final Recommendation

### For Ducla Agent Development/Testing: **Docker Compose**
- Faster iteration
- Easier debugging
- Resource efficient
- Focus on integration testing

### For Production Preparation: **K3s**
- Production-like environment
- Kubernetes-native monitoring
- Auto-scaling testing
- Cloud deployment preparation

## 🚀 Quick Start Commands

### Docker Compose (Recommended)
```bash
./start-monitoring.sh
# Access: http://localhost:3000 (Grafana)
```

### K3s (Advanced)
```bash
# Install K3s
curl -sfL https://get.k3s.io | sh -

# Install monitoring stack
helm install monitoring prometheus-community/kube-prometheus-stack -n monitoring --create-namespace

# Deploy Ducla Agent
kubectl apply -f ducla-agent-k3s.yaml
```

**🎉 Choose based on your needs: Docker for quick testing, K3s for production-like environment!**