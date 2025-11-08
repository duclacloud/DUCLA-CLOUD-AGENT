# Ducla Cloud Agent - Hướng Dẫn Thiết Lập Đầy Đủ

## Tổng Quan Dự Án

Ducla Cloud Agent là một agent hiệu suất cao, production-ready cho distributed task execution và system monitoring, được xây dựng bằng Go.

## Cấu Trúc Dự Án

```
ducla-cloud-agent/
├── cmd/agent/                  # Main application entry point
│   ├── main.go                # Entry point với version support
│   ├── version.go             # Version information
│   └── service_windows.go     # Windows service support
├── internal/                   # Internal packages
│   ├── agent/                 # Core agent logic
│   ├── api/                   # HTTP và gRPC APIs
│   │   ├── server.go
│   │   ├── http_handlers.go
│   │   ├── grpc_service.go
│   │   └── grpc_types.go
│   ├── config/                # Configuration management
│   ├── executor/              # Task execution
│   ├── fileops/               # File operations
│   │   └── fileops.go
│   ├── health/                # Health checks
│   │   ├── health.go
│   │   └── checks.go
│   ├── metrics/               # Metrics collection
│   │   ├── metrics.go
│   │   └── collectors.go
│   └── transport/             # Network transport
├── pkg/                       # Public packages
│   └── utils/
├── configs/                   # Configuration files
│   └── agent.yaml            # Main configuration
├── scripts/                   # Build và deployment scripts
│   ├── build-cli.sh          # Linux/macOS build script
│   ├── build-cli.ps1         # Windows build script
│   ├── dev.sh                # Development helper
│   ├── install.sh            # Linux installation
│   ├── uninstall.sh          # Linux uninstallation
│   ├── release.sh            # Release automation
│   ├── package-deb.sh        # Ubuntu/Debian packaging
│   ├── package-rpm.sh        # RedHat/CentOS packaging
│   ├── package-windows.ps1   # Windows installer
│   ├── build-windows-portable.ps1  # Windows portable
│   ├── package-all.sh        # Build all Linux packages
│   └── package-all.ps1       # Build all Windows packages
├── k8s/                      # Kubernetes manifests
│   ├── namespace.yaml
│   ├── serviceaccount.yaml
│   ├── configmap.yaml
│   ├── secret.yaml
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── hpa.yaml
│   ├── pdb.yaml
│   ├── ingress.yaml
│   ├── servicemonitor.yaml
│   └── kustomization.yaml
├── docs/                     # Documentation
│   ├── DEPLOYMENT.md
│   ├── PACKAGING.md
│   └── NETWORK.md
├── .github/workflows/        # CI/CD workflows
│   ├── build.yml
│   ├── release.yml
│   └── package.yml
├── Dockerfile                # Multi-stage Docker build
├── docker-compose.yml        # Production Docker Compose
├── docker-compose.dev.yml    # Development Docker Compose
├── .dockerignore
├── Makefile                  # Build automation
├── .golangci.yml            # Linter configuration
├── .env.example             # Environment variables template
├── go.mod                   # Go dependencies
└── README.md                # Project documentation
```

## Các Tính Năng Chính

### 1. Core Features
- ✅ High Performance: Built với Go
- ✅ Secure: JWT authentication, RBAC, audit logging, TLS
- ✅ Monitoring: Prometheus metrics, health checks
- ✅ Extensible: Plugin system (Docker, Kubernetes, AWS)
- ✅ Multi-Protocol: HTTP REST và gRPC APIs
- ✅ Production-Ready: Resource limits, proper logging

### 2. API Endpoints

#### HTTP API (Port 8080)
- `GET /api/v1/status` - Agent status
- `GET /api/v1/tasks` - List tasks
- `POST /api/v1/tasks` - Create task
- `GET /api/v1/tasks/:id` - Get task details
- `DELETE /api/v1/tasks/:id` - Cancel task

#### Health Check (Port 8081)
- `GET /health` - Health status
- `GET /ready` - Readiness check

#### Metrics (Port 9090)
- `GET /metrics` - Prometheus metrics

#### gRPC API (Port 8443)
- Task execution services
- File operations
- System monitoring

### 3. Configuration

File cấu hình chính: `configs/agent.yaml`

```yaml
agent:
  id: "${HOSTNAME}"
  name: "${HOSTNAME}"
  environment: "production"
  capabilities:
    - "file_operations"
    - "task_execution"
    - "system_monitoring"
    - "container_management"

master:
  url: "${DUCLA_MASTER_URL}"
  token: "${DUCLA_AGENT_TOKEN}"
  heartbeat_interval: 30s

api:
  http:
    enabled: true
    port: 8080
  grpc:
    enabled: true
    port: 8443

security:
  jwt:
    secret: "${DUCLA_JWT_SECRET}"
  rbac:
    enabled: true
  audit:
    enabled: true

metrics:
  enabled: true
  port: 9090
  collectors:
    - "system"
    - "process"
    - "agent"
```

## Hướng Dẫn Build và Deploy

### 1. Build Binary

#### Linux/macOS
```bash
# Quick build cho platform hiện tại
./scripts/build-cli.sh quick

# Build cho tất cả platforms
./scripts/build-cli.sh

# Hoặc dùng Make
make build
```

#### Windows
```powershell
# Quick build
.\scripts\build-cli.ps1 -Mode quick

# Build cho tất cả platforms
.\scripts\build-cli.ps1
```

### 2. Development

```bash
# Run trong development mode
./scripts/dev.sh run

# Run tests
./scripts/dev.sh test

# Generate coverage report
./scripts/dev.sh coverage

# Run linter
./scripts/dev.sh lint

# Format code
./scripts/dev.sh fmt

# Tidy dependencies
./scripts/dev.sh tidy
```

### 3. Docker Deployment

#### Docker Compose (Production)
```bash
# Copy environment file
cp .env.example .env

# Edit configuration
nano .env

# Start services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

#### Docker Compose (Development)
```bash
docker-compose -f docker-compose.dev.yml up
```

#### Build Docker Image
```bash
# Build
docker build -t ducla-cloud-agent:latest .

# Run
docker run -d \
  --name ducla-agent \
  -p 8080:8080 -p 8081:8081 -p 8443:8443 -p 9090:9090 \
  -e DUCLA_MASTER_URL=https://master.ducla.cloud \
  -e DUCLA_AGENT_TOKEN=your-token \
  ducla-cloud-agent:latest
```

### 4. Kubernetes Deployment

```bash
# Deploy tất cả resources
kubectl apply -k k8s/

# Check status
kubectl get all -n ducla-system

# View logs
kubectl logs -n ducla-system -l app=ducla-agent -f

# Scale deployment
kubectl scale deployment ducla-agent -n ducla-system --replicas=5

# Delete deployment
kubectl delete -k k8s/
```

#### Hoặc dùng Make
```bash
make k8s-deploy
make k8s-status
make k8s-logs
make k8s-delete
```

## Packaging cho Distribution

### 1. Ubuntu/Debian (DEB Package)

```bash
# Build DEB package
./scripts/package-deb.sh

# Hoặc với version cụ thể
VERSION=1.0.0 ./scripts/package-deb.sh

# Install
sudo dpkg -i dist/ducla-agent_1.0.0_amd64.deb
sudo apt-get install -f

# Uninstall
sudo apt-get remove ducla-agent

# Purge (xóa cả config)
sudo apt-get purge ducla-agent
```

**Package bao gồm:**
- Binary: `/usr/bin/ducla-agent`
- Config: `/etc/ducla/agent.yaml`
- Systemd service: `/lib/systemd/system/ducla-agent.service`
- Data: `/opt/ducla/data`
- Logs: `/var/log/ducla`
- User: `ducla`

### 2. RedHat/CentOS/Fedora (RPM Package)

```bash
# Build RPM package
./scripts/package-rpm.sh

# Hoặc với version cụ thể
VERSION=1.0.0 ./scripts/package-rpm.sh

# Install
sudo rpm -ivh dist/ducla-agent-1.0.0-1.x86_64.rpm

# Hoặc với yum/dnf
sudo yum install dist/ducla-agent-1.0.0-1.x86_64.rpm

# Uninstall
sudo rpm -e ducla-agent
```

### 3. Windows Portable Package (ZIP)

```powershell
# Build portable package
.\scripts\build-windows-portable.ps1 -Version "1.0.0"

# Extract và sử dụng
Expand-Archive dist\ducla-agent-1.0.0-windows-amd64-portable.zip -DestinationPath C:\ducla-agent
cd C:\ducla-agent

# Run trong console mode
.\start.bat

# Hoặc install as service
.\install-service.bat
```

**Package bao gồm:**
- `ducla-agent.exe` - Binary
- `config\agent.yaml` - Configuration
- `start.bat` - Start trong console mode
- `install-service.bat` - Install Windows service
- `uninstall-service.bat` - Uninstall service
- `README.txt` - Documentation

### 4. Windows Installer (EXE)

```powershell
# Build installer (yêu cầu Inno Setup)
.\scripts\package-windows.ps1 -Version "1.0.0"

# Run installer
.\dist\ducla-agent-1.0.0-windows-amd64-setup.exe
```

**Installer features:**
- Automatic service installation
- Start menu shortcuts
- Uninstaller
- Configuration wizard

### 5. Build All Packages

```bash
# Linux (DEB + RPM)
VERSION=1.0.0 ./scripts/package-all.sh

# Windows (Portable + Installer)
.\scripts\package-all.ps1 -Version "1.0.0"
```

## Service Management

### Linux (systemd)

```bash
# Start service
sudo systemctl start ducla-agent

# Stop service
sudo systemctl stop ducla-agent

# Restart service
sudo systemctl restart ducla-agent

# Check status
sudo systemctl status ducla-agent

# Enable auto-start
sudo systemctl enable ducla-agent

# Disable auto-start
sudo systemctl disable ducla-agent

# View logs
sudo journalctl -u ducla-agent -f
```

### Windows

```powershell
# Start service
sc start DuclaAgent

# Stop service
sc stop DuclaAgent

# Query status
sc query DuclaAgent

# View service details
sc qc DuclaAgent

# View logs
Get-Content C:\ProgramData\Ducla\logs\agent.log -Tail 50 -Wait
```

## Configuration

### Environment Variables

Tạo file `.env` từ template:
```bash
cp .env.example .env
```

Các biến quan trọng:
```bash
# Master Server
DUCLA_MASTER_URL=https://master.ducla.cloud
DUCLA_AGENT_TOKEN=your-agent-token-here

# Security
DUCLA_JWT_SECRET=your-jwt-secret-here

# Agent Settings
DUCLA_ENVIRONMENT=production
DUCLA_REGION=us-east-1
DUCLA_ZONE=us-east-1a

# Logging
DUCLA_LOG_LEVEL=info
```

### Linux Configuration Files

```bash
# Main config
/etc/ducla/agent.yaml

# Environment variables (systemd)
/etc/default/ducla-agent

# Data directory
/opt/ducla/data

# Logs
/var/log/ducla
```

### Windows Configuration Files

```
# Main config
C:\ProgramData\Ducla\agent.yaml

# Data directory
C:\ProgramData\Ducla\data

# Logs
C:\ProgramData\Ducla\logs
```

## Monitoring và Health Checks

### Health Check Endpoints

```bash
# Health check
curl http://localhost:8081/health

# Readiness check
curl http://localhost:8081/ready
```

### Prometheus Metrics

```bash
# View metrics
curl http://localhost:9090/metrics
```

**Available metrics:**
- System metrics (CPU, memory, disk)
- Process metrics
- Agent-specific metrics
- Task execution metrics
- File operation metrics

### Kubernetes Monitoring

```bash
# Check pod health
kubectl get pods -n ducla-system

# View pod logs
kubectl logs -n ducla-system -l app=ducla-agent -f

# Check metrics
kubectl top pods -n ducla-system

# View events
kubectl get events -n ducla-system
```

## CI/CD Workflows

### GitHub Actions

**Build Workflow** (`.github/workflows/build.yml`):
- Runs on push to main/develop
- Runs tests with coverage
- Runs linter
- Builds for multiple platforms
- Uploads artifacts

**Release Workflow** (`.github/workflows/release.yml`):
- Triggers on version tags (v*)
- Builds binaries for all platforms
- Creates GitHub release
- Builds and pushes Docker images

**Package Workflow** (`.github/workflows/package.yml`):
- Builds DEB packages
- Builds RPM packages
- Builds Windows packages
- Creates release with all packages

### Creating a Release

```bash
# Create and push tag
./scripts/release.sh v1.0.0

# Or manually
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

## Troubleshooting

### Check Logs

```bash
# Linux (systemd)
sudo journalctl -u ducla-agent -f

# Linux (file)
sudo tail -f /var/log/ducla/agent.log

# Windows
Get-Content C:\ProgramData\Ducla\logs\agent.log -Tail 50 -Wait

# Docker
docker-compose logs -f ducla-agent

# Kubernetes
kubectl logs -n ducla-system -l app=ducla-agent -f
```

### Common Issues

#### Service won't start
```bash
# Check configuration
ducla-agent --config /etc/ducla/agent.yaml --debug

# Check permissions
ls -la /opt/ducla /var/log/ducla

# Check systemd status
systemctl status ducla-agent
```

#### Connection issues
```bash
# Test master server connection
curl -v $DUCLA_MASTER_URL

# Check network
netstat -tulpn | grep ducla-agent

# Check firewall
sudo ufw status
```

#### High resource usage
```bash
# Check metrics
curl http://localhost:9090/metrics

# Monitor resources
top -p $(pgrep ducla-agent)

# Check logs for errors
journalctl -u ducla-agent --since "1 hour ago"
```

## Security Best Practices

1. **Use strong secrets**: Generate secure JWT secrets
2. **Enable TLS**: Configure TLS for gRPC and HTTP
3. **Firewall rules**: Restrict access to necessary IPs
4. **Regular updates**: Keep agent updated
5. **Monitor logs**: Enable audit logging
6. **RBAC**: Configure role-based access control
7. **Secrets management**: Use external secrets (Vault, etc.)

## Performance Tuning

### Resource Limits

**Docker:**
```yaml
deploy:
  resources:
    limits:
      cpus: '2'
      memory: 1G
    reservations:
      cpus: '0.5'
      memory: 256M
```

**Kubernetes:**
```yaml
resources:
  requests:
    cpu: 100m
    memory: 128Mi
  limits:
    cpu: 1000m
    memory: 512Mi
```

### Configuration Tuning

```yaml
executor:
  max_concurrent_tasks: 10
  task_timeout: 30m
  worker_pool_size: 5
  queue_size: 100

storage:
  cleanup:
    enabled: true
    interval: 1h
    max_age: 24h
    max_disk_usage: 80
```

## Support và Documentation

- **Documentation**: `docs/` directory
- **GitHub Issues**: https://github.com/your-org/ducla-cloud-agent/issues
- **GitHub Discussions**: https://github.com/your-org/ducla-cloud-agent/discussions

## License

MIT License - See LICENSE file for details

---

## Quick Reference Commands

### Development
```bash
./scripts/dev.sh run          # Run agent
./scripts/dev.sh test         # Run tests
./scripts/dev.sh lint         # Run linter
./scripts/dev.sh fmt          # Format code
```

### Build
```bash
./scripts/build-cli.sh quick  # Quick build
./scripts/build-cli.sh        # Multi-platform build
make build                    # Using Make
```

### Package
```bash
./scripts/package-deb.sh      # DEB package
./scripts/package-rpm.sh      # RPM package
./scripts/package-all.sh      # All Linux packages
.\scripts\package-all.ps1     # All Windows packages
```

### Docker
```bash
docker-compose up -d          # Start
docker-compose logs -f        # Logs
docker-compose down           # Stop
```

### Kubernetes
```bash
kubectl apply -k k8s/         # Deploy
kubectl get all -n ducla-system  # Status
kubectl logs -n ducla-system -l app=ducla-agent -f  # Logs
kubectl delete -k k8s/        # Delete
```

### Service
```bash
# Linux
sudo systemctl start ducla-agent
sudo systemctl status ducla-agent
sudo journalctl -u ducla-agent -f

# Windows
sc start DuclaAgent
sc query DuclaAgent
```

---

**Tài liệu này được tạo tự động từ quá trình setup dự án Ducla Cloud Agent**
