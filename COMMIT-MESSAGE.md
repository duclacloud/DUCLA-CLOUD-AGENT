feat: initial release of Ducla Cloud Agent v1.0.0

Complete cloud agent with comprehensive feature set:

🖥️ CLI Interface:
- Complete CLI commands (show, task, file, config)
- Professional man page (man ducla-agent)
- Auto config detection (agent.yaml)
- No warning messages for clean UX

🌐 REST API:
- Main HTTP API (Port 8080): Task management, file operations
- Health Check API (Port 8081): System health, readiness, liveness  
- Metrics API (Port 9090): Prometheus metrics, system uptime
- Complete API documentation (API-REFERENCE.md)

📤 Output Destinations:
- Cloud Services: AWS CloudWatch, S3, Google Stackdriver
- Monitoring: Prometheus, Elasticsearch, Datadog, Grafana Loki
- Network: HTTP/HTTPS, TCP/UDP, Apache Kafka
- Notifications: Slack, Email, Discord
- 15+ supported destinations with examples

📚 Documentation:
- Professional README with badges
- Complete Vietnamese documentation (README-VI.md)
- Comprehensive API reference (API-REFERENCE.md)
- Contributing guidelines (CONTRIBUTING.md)
- Workshop materials (WORKSHOP.md)
- CLI features guide (CLI-FEATURES.md)

📦 Package Management:
- DEB packages for Ubuntu/Debian
- RPM packages for CentOS/RHEL
- Systemd service integration
- Proper maintainer metadata

🛠️ Development Tools:
- Automated build scripts with version info
- Demo scripts for feature showcase
- Package creation scripts
- Git configuration (.gitignore, LICENSE)

🔒 Production Ready:
- TLS support for all APIs
- JWT authentication framework
- RBAC authorization system
- Comprehensive logging and monitoring
- Security best practices

Technical Details:
- Language: Go 1.21+
- Binary Size: ~17MB optimized
- Memory Usage: ~50MB baseline
- Platforms: Linux (amd64, arm64)
- License: MIT

Repository: github.com/duclacloud/DUCLA-CLOUD-AGENT
Maintainer: mandỵhades <mandỵhades@hotmail.com.vn>

This release represents a complete, production-ready cloud agent
suitable for enterprise deployment and open source community use.