# Ducla Cloud Agent v1.0.0-1-g23271c6-dirty

High-performance cloud agent for distributed task execution and system monitoring.

## 📦 Available Packages

- **Ubuntu/Debian**: `ducla-agent_1.0.0-1-g23271c6-dirty_amd64.deb`
- **RedHat/CentOS/Fedora**: `ducla-agent-1.0.0-1-g23271c6-dirty.x86_64.rpm`
- **Linux Binary**: `ducla-agent-linux-amd64`

## 🚀 Quick Installation

### Automatic Installation (Recommended)
```bash
# Download and run installer
curl -sSL https://github.com//releases/latest/download/install.sh | sudo bash
```

### Manual Installation

#### Ubuntu/Debian
```bash
sudo dpkg -i ducla-agent_1.0.0-1-g23271c6-dirty_amd64.deb
sudo systemctl enable ducla-agent
sudo systemctl start ducla-agent
```

#### RedHat/CentOS/Fedora
```bash
sudo rpm -ivh ducla-agent-1.0.0-1-g23271c6-dirty.x86_64.rpm
sudo systemctl enable ducla-agent
sudo systemctl start ducla-agent
```

#### Binary Installation
```bash
# Run the install script in this directory
sudo ./install.sh
```

## 🔧 Configuration

Edit `/etc/ducla/agent.yaml` to configure the agent:

```yaml
agent:
  id: "${HOSTNAME}"
  name: "${HOSTNAME}"
  environment: "production"

api:
  http:
    enabled: true
    port: 8080

logging:
  level: "info"
  format: "json"
```

## 📊 API Endpoints

- **Status**: http://localhost:8080/api/v1/status
- **Health**: http://localhost:8081/health  
- **Metrics**: http://localhost:9090/metrics

## 🔍 Service Management

```bash
# Check status
sudo systemctl status ducla-agent

# View logs
sudo journalctl -u ducla-agent -f

# Restart service
sudo systemctl restart ducla-agent
```

## 🗑️ Uninstallation

```bash
# Run uninstall script
sudo ./uninstall.sh
```

## 🔐 Security

- Service runs as non-root user `ducla`
- Configuration files have restricted permissions
- TLS support for secure communication

## 📋 System Requirements

- **OS**: Linux (Ubuntu 18.04+, CentOS 7+, RHEL 7+, Fedora 30+)
- **Architecture**: x86_64 (amd64)
- **RAM**: 512MB minimum
- **Disk**: 100MB free space

## 🆘 Support

- **Documentation**: https://github.com/
- **Issues**: https://github.com//issues
- **Email**: support@ducla.cloud

---
**Version**: 1.0.0-1-g23271c6-dirty  
**Build Date**: 2025-11-09 15:57:44 UTC
