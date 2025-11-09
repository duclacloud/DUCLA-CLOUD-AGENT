# 🚀 DUCLA CLOUD AGENT

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![Release](https://img.shields.io/badge/Release-v1.0.0-green.svg)](https://github.com/duclacloud/DUCLA-CLOUD-AGENT/releases)

A high-performance, production-ready cloud agent for distributed task execution and system monitoring.

## Features

- 🚀 **High Performance**: Built with Go for optimal performance and low resource usage
- 🔒 **Secure**: JWT authentication, RBAC, audit logging, and TLS support
- 📊 **Monitoring**: Built-in metrics collection and health checks
- 🔌 **Extensible**: Plugin system for Docker, Kubernetes, and cloud providers
- 🌐 **Multi-Protocol**: HTTP REST and gRPC APIs
- 📦 **Easy Deployment**: Docker, Kubernetes, and systemd support

## Quick Start

### Using Docker

```bash
# Create environment file
cp .env.example .env

# Edit configuration
nano .env

# Start the agent
docker-compose up -d

# Check logs
docker-compose logs -f
```

### Using Binary

```bash
# Download latest release
curl -L https://github.com/duclacloud/DUCLA-CLOUD-AGENT/releases/latest/download/ducla-agent-linux-amd64.tar.gz | tar xz

# Run the agent
./ducla-agent --config configs/agent.yaml
```

### Using Installation Script

```bash
# Install (requires root)
curl -sSL https://raw.githubusercontent.com/your-org/ducla-cloud-agent/main/scripts/install.sh | sudo bash

# Configure
sudo nano /etc/ducla/agent.yaml

# Start service
sudo systemctl start ducla-agent
```

## Building from Source

### Prerequisites

- Go 1.21 or later
- Make (optional)

### Build

```bash
# Clone repository
git clone https://github.com/duclacloud/DUCLA-CLOUD-AGENT.git
cd ducla-cloud-agent

# Build for current platform
./scripts/build-cli.sh quick

# Or build for all platforms
./scripts/build-cli.sh

# Or use Make
make build
```

### Development

```bash
# Run in development mode
./scripts/dev.sh run

# Run tests
./scripts/dev.sh test

# Generate coverage report
./scripts/dev.sh coverage

# Run linter
./scripts/dev.sh lint

# Format code
./scripts/dev.sh fmt
```

## Configuration

The agent is configured via YAML file. See [configs/agent.yaml](configs/agent.yaml) for full configuration options.

Key configuration sections:

- **agent**: Agent identification and capabilities
- **master**: Master server connection settings
- **api**: HTTP and gRPC API configuration
- **security**: Authentication, RBAC, and audit settings
- **storage**: Data storage and cleanup
- **logging**: Log level and format
- **metrics**: Prometheus metrics configuration
- **health**: Health check endpoints
- **plugins**: Plugin system configuration
- **executor**: Task execution settings

## Deployment

### Docker Compose

See [docker-compose.yml](docker-compose.yml) for production deployment.

```bash
docker-compose up -d
```

### Kubernetes

See [k8s/](k8s/) directory for Kubernetes manifests.

```bash
# Deploy with kubectl
kubectl apply -k k8s/

# Or use Make
make k8s-deploy
```

For detailed deployment instructions, see [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md).

## API Endpoints

### HTTP API (Port 8080)

- `GET /api/v1/status` - Agent status
- `GET /api/v1/tasks` - List tasks
- `POST /api/v1/tasks` - Create task
- `GET /api/v1/tasks/:id` - Get task details
- `DELETE /api/v1/tasks/:id` - Cancel task

### Health Check (Port 8081)

- `GET /health` - Health status
- `GET /ready` - Readiness check

### Metrics (Port 9090)

- `GET /metrics` - Prometheus metrics

### gRPC API (Port 8443)

See [internal/api/grpc_types.go](internal/api/grpc_types.go) for gRPC service definitions.

## Monitoring

The agent exposes Prometheus metrics on `/metrics` endpoint:

- System metrics (CPU, memory, disk)
- Process metrics
- Agent-specific metrics
- Task execution metrics
- File operation metrics

## Security

- **Authentication**: JWT tokens
- **Authorization**: Role-based access control (RBAC)
- **Audit Logging**: All operations are logged
- **TLS**: Support for TLS encryption
- **Firewall**: IP-based access control
- **Rate Limiting**: Request rate limiting

## Plugins

The agent supports plugins for extended functionality:

- **Docker**: Container management
- **Kubernetes**: Pod and service management
- **AWS**: AWS service integration

## 🌐 REST API

Ducla Cloud Agent provides comprehensive REST APIs:

- **Main API (Port 8080)**: Task management, file operations, agent status
- **Health API (Port 8081)**: Health checks, readiness, liveness
- **Metrics API (Port 9090)**: Prometheus metrics, system uptime

### Quick API Examples
```bash
# Agent status
curl http://localhost:8080/api/v1/status

# Health check
curl http://localhost:8081/health

# Prometheus metrics
curl http://localhost:9090/metrics
```

**📖 Complete API Reference**: See [API-REFERENCE.md](API-REFERENCE.md) for full REST API documentation and output destinations.

## Architecture

```
ducla-cloud-agent/
├── cmd/agent/          # Main application entry point
├── internal/           # Internal packages
│   ├── agent/         # Core agent logic
│   ├── api/           # HTTP and gRPC APIs
│   ├── config/        # Configuration management
│   ├── executor/      # Task execution
│   ├── fileops/       # File operations
│   ├── health/        # Health checks
│   ├── metrics/       # Metrics collection
│   └── transport/     # Network transport
├── pkg/               # Public packages
├── configs/           # Configuration files
├── scripts/           # Build and deployment scripts
├── k8s/              # Kubernetes manifests
├── docs/             # Documentation
└── API-REFERENCE.md   # Complete REST API documentation
```

## 🤝 Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for details.

### Quick Start for Contributors

```bash
# Fork and clone
git clone https://github.com/YOUR_USERNAME/DUCLA-CLOUD-AGENT.git
cd DUCLA-CLOUD-AGENT

# Build and test
./build-v1.sh
./demo-auto.sh

# Create feature branch
git checkout -b feature/your-feature
```

## 📄 License

[MIT License](LICENSE) - see the [LICENSE](LICENSE) file for details.

## 👨‍💻 Author

**mandỵhades** - [mandỵhades@hotmail.com.vn](mailto:mandỵhades@hotmail.com.vn)

## 🆘 Support

- 📖 **Documentation**: [README-VI.md](README-VI.md) (Vietnamese)
- 🐛 **Issues**: [GitHub Issues](https://github.com/duclacloud/DUCLA-CLOUD-AGENT/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/duclacloud/DUCLA-CLOUD-AGENT/discussions)
- 📧 **Email**: [mandỵhades@hotmail.com.vn](mailto:mandỵhades@hotmail.com.vn)

## ⭐ Show Your Support

Give a ⭐️ if this project helped you!