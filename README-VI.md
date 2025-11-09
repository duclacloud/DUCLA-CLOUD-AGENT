# Ducla Cloud Agent - Hướng Dẫn Tiếng Việt

## Giới Thiệu

**Ducla Cloud Agent** là một agent cloud hiệu suất cao được viết bằng Go, được thiết kế để thực thi các tác vụ phân tán và giám sát hệ thống trong môi trường production.

## Tính Năng Chính

- 🚀 **Hiệu suất cao**: Được xây dựng bằng Go để tối ưu hiệu suất và sử dụng ít tài nguyên
- 🔒 **Bảo mật**: Hỗ trợ JWT authentication, RBAC, audit logging, và TLS
- 📊 **Giám sát**: Thu thập metrics tích hợp và health checks
- 🔌 **Mở rộng**: Hệ thống plugin cho Docker, Kubernetes, và cloud providers
- 🌐 **Đa giao thức**: HTTP REST và gRPC APIs
- 📦 **Triển khai dễ dàng**: Hỗ trợ Docker, Kubernetes, và systemd

## Yêu Cầu Hệ Thống

- **Hệ điều hành**: Linux (Ubuntu, CentOS, RHEL, etc.)
- **Go**: Phiên bản 1.21 trở lên
- **RAM**: Tối thiểu 512MB
- **Disk**: Tối thiểu 100MB trống
- **Network**: Kết nối internet để tải dependencies

## Cài Đặt

### 1. Cài Đặt Go (nếu chưa có)

```bash
# Tải Go 1.21.5
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz

# Cài đặt
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz

# Thêm vào PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Kiểm tra
go version
```

### 2. Build Dự Án

```bash
# Clone repository (nếu từ git)
git clone <repository-url>
cd ducla-cloud-agent

# Hoặc nếu đã có source code
cd ducla-cloud-agent

# Tải dependencies
go mod tidy

# Build binary
go build -o ducla-agent ./cmd/agent

# Kiểm tra build thành công
ls -la ducla-agent
```

## Cấu Hình

### File Cấu Hình Mặc Định

Agent sẽ tự động tìm config file theo thứ tự:
1. `agent.yaml` (current directory)
2. `/etc/ducla/agent.yaml` (system-wide)
3. File được chỉ định với `-config`

### Tạo File Cấu Hình

Tạo file `agent.yaml`:

```yaml
# Cấu hình cơ bản
agent:
  id: "my-agent"
  name: "My Ducla Agent"
  environment: "production"
  region: "vietnam"
  capabilities:
    - "file_operations"
    - "task_execution"
    - "system_monitoring"

# Kết nối master server (tùy chọn)
master:
  url: "ws://your-master-server:9000"
  token: "your-auth-token"
  connect_timeout: 30s
  heartbeat_interval: 30s

# API server
api:
  http:
    enabled: true
    address: "0.0.0.0"
    port: 8080
    tls:
      enabled: false

# Bảo mật
security:
  auth:
    enabled: false
  rbac:
    enabled: false

# Lưu trữ
storage:
  data_dir: "./data"
  temp_dir: "./tmp"
  cleanup:
    enabled: true
    interval: 1h
    max_age: 24h

# Logging
logging:
  level: "info"
  format: "text"
  output: "stdout"

# Metrics
metrics:
  enabled: true
  address: "0.0.0.0"
  port: 9090

# Health check
health:
  enabled: true
  address: "0.0.0.0"
  port: 8081

# Task executor
executor:
  workers: 4
  queue_size: 100
  task_timeout: 300s
```

## Chạy Agent

### Chạy Cơ Bản

```bash
# Chạy với config mặc định
./ducla-agent

# Chạy với config tùy chỉnh
./ducla-agent -config agent.yaml

# Chạy với debug mode
./ducla-agent -config agent.yaml -debug

# Xem version
./ducla-agent -version

# Xem help
./ducla-agent --help
```

### CLI Commands

```bash
# Xem trạng thái agent
./ducla-agent show status

# Xem sức khỏe hệ thống
./ducla-agent show health

# Liệt kê tất cả tasks
./ducla-agent show tasks

# Liệt kê tasks đang chạy
./ducla-agent show tasks running

# Tạo task mới
./ducla-agent task create "echo 'Hello World'"

# Hủy task
./ducla-agent task cancel <task-id>

# Liệt kê files
./ducla-agent file list /tmp

# Copy file
./ducla-agent file copy /tmp/source.txt /tmp/dest.txt

# Move file
./ducla-agent file move /tmp/source.txt /tmp/dest.txt

# Delete file
./ducla-agent file delete /tmp/file.txt

# Validate config
./ducla-agent config validate

# Test config và connectivity
./ducla-agent config test

# Xem manual đầy đủ
man ducla-agent
```

### Chạy Như Service (systemd)

Tạo file `/etc/systemd/system/ducla-agent.service`:

```ini
[Unit]
Description=Ducla Cloud Agent
After=network.target

[Service]
Type=simple
User=ducla
Group=ducla
WorkingDirectory=/opt/ducla
ExecStart=/opt/ducla/ducla-agent -config /etc/ducla/agent.yaml
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

Khởi động service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable ducla-agent
sudo systemctl start ducla-agent
sudo systemctl status ducla-agent
```

## API Endpoints

### 🚀 Main HTTP API (Port 8080)

#### Agent Status & Info
```bash
# Agent status
curl http://localhost:8080/api/v1/status

# Configuration info  
curl http://localhost:8080/api/v1/config

# Version info
curl http://localhost:8080/api/v1/version
```

#### Quản Lý Tasks
```bash
# Liệt kê tất cả tasks
curl http://localhost:8080/api/v1/tasks

# Liệt kê tasks đang chạy
curl http://localhost:8080/api/v1/tasks?filter=running

# Tạo task mới
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "type": "shell",
    "name": "test-task",
    "command": "echo",
    "args": ["Hello World"]
  }'

#### File Operations
```bash
# Liệt kê files
curl -X POST http://localhost:8080/api/v1/files \
  -H "Content-Type: application/json" \
  -d '{
    "type": "list",
    "source_path": "/tmp"
  }'

# Copy file
curl -X POST http://localhost:8080/api/v1/files \
  -H "Content-Type: application/json" \
  -d '{
    "type": "copy",
    "source_path": "/tmp/source.txt",
    "dest_path": "/tmp/dest.txt"
  }'

# Move file
curl -X POST http://localhost:8080/api/v1/files \
  -H "Content-Type: application/json" \
  -d '{
    "type": "move",
    "source_path": "/tmp/old.txt",
    "dest_path": "/tmp/new.txt"
  }'

# Delete file
curl -X POST http://localhost:8080/api/v1/files \
  -H "Content-Type: application/json" \
  -d '{
    "type": "delete",
    "source_path": "/tmp/unwanted.txt"
  }'
```

### 🏥 Health Check API (Port 8081)
```bash
# Kiểm tra sức khỏe tổng thể
curl http://localhost:8081/health

# Kiểm tra readiness
curl http://localhost:8081/ready

# Kiểm tra liveness
curl http://localhost:8081/live
```

### 📊 Metrics API (Port 9090)
```bash
# Xem Prometheus metrics
curl http://localhost:9090/metrics

# System uptime
curl http://localhost:9090/api/v1/uptime

# Custom metrics
curl http://localhost:9090/api/v1/metrics
```

### 📖 Chi Tiết API
Xem [API-REFERENCE.md](API-REFERENCE.md) để biết đầy đủ REST API commands, examples và output destinations.

## 📤 Output Destinations

Ducla Cloud Agent hỗ trợ gửi logs và metrics đến nhiều destinations:

### ☁️ Cloud Services
```bash
# AWS CloudWatch
cloudwatch_logs-o cloudwatch_logs -p log_group_name=ducla-agent -p region=ap-southeast-1

# AWS S3
s3-o s3 -p bucket=ducla-logs -p region=ap-southeast-1 -p total_file_size=5M

# Google Stackdriver
stackdriver-o stackdriver -p resource=k8s_container -p project_id=my-project
```

### 📊 Monitoring & Analytics
```bash
# Prometheus
prometheus-o prometheus -p host=10.0.0.5 -p port=9090 -p format=metrics

# Elasticsearch
es-o es -p host=10.0.0.10 -p port=9200 -p index=ducla-logs

# Datadog
datadog-o datadog -p apikey=YOUR_API_KEY -p tags=env:prod,service:ducla-agent

# Grafana Loki
loki-o loki -p host=http://10.0.0.5:3100 -p labels=job:ducla-agent
```

### 🌐 Network Protocols
```bash
# HTTP/HTTPS
http-o http -p host=10.0.0.5 -p port=8080 -p format=json

# TCP/UDP
tcp-o tcp -p host=10.0.0.10 -p port=9000
udp-o udp -p host=10.0.0.10 -p port=514

# Apache Kafka
kafka-o kafka -p brokers=localhost:9092 -p topics=logs
```

### 💬 Notifications
```bash
# Slack
slack-o slack -p webhook_url=https://hooks.slack.com/services/YOUR/WEBHOOK

# Email
email-o email -p smtp_host=smtp.gmail.com -p to=admin@example.com

# Discord
discord-o discord -p webhook_url=https://discord.com/api/webhooks/YOUR/WEBHOOK
```

**📖 Chi tiết đầy đủ**: Xem [API-REFERENCE.md](API-REFERENCE.md) để biết tất cả output destinations và cấu hình.

## Giám Sát

### Health Checks

Agent tự động kiểm tra:

- **System**: CPU, Memory, Disk usage
- **Network**: Kết nối mạng
- **Services**: Trạng thái các service nội bộ

### Metrics

Agent thu thập metrics cho:

- System metrics (CPU, memory, disk)
- Process metrics
- Task execution metrics
- File operation metrics
- Custom application metrics

### Logging

Logs được xuất ra theo format JSON hoặc text:

```bash
# Xem logs realtime
tail -f /var/log/ducla-agent.log

# Hoặc nếu chạy với systemd
journalctl -u ducla-agent -f
```

## Troubleshooting

### Lỗi Thường Gặp

#### 1. Không thể bind port

```
Error: bind: address already in use
```

**Giải pháp**: Thay đổi port trong config hoặc kill process đang sử dụng port

#### 2. Permission denied

```
Error: permission denied
```

**Giải pháp**: Chạy với quyền phù hợp hoặc thay đổi user/group

#### 3. Config file not found

```
Error: config file not found
```

**Giải pháp**: Kiểm tra đường dẫn config file

#### 4. Master connection failed

```
Error: Failed to connect to master server
```

**Giải pháp**: Kiểm tra URL master server và network connectivity

### Debug Mode

Chạy với debug để xem thông tin chi tiết:

```bash
./ducla-agent -config config.yaml -debug
```

### Kiểm Tra Logs

```bash
# Xem logs với journalctl
sudo journalctl -u ducla-agent -n 100

# Xem logs file
tail -f /var/log/ducla-agent.log
```

## Development

### Build từ Source

```bash
# Clone repository
git clone <repo-url>
cd ducla-cloud-agent

# Install dependencies
go mod tidy

# Run tests
go test ./...

# Build
go build -o ducla-agent ./cmd/agent

# Run development server
go run ./cmd/agent -config config.yaml -debug
```

### Cấu Trúc Dự Án

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
└── docs/             # Documentation
```

## Deployment

### Docker

```bash
# Build Docker image
docker build -t ducla-agent .

# Run container
docker run -d \
  --name ducla-agent \
  -p 8080:8080 \
  -p 8081:8081 \
  -p 9090:9090 \
  -v $(pwd)/config.yaml:/etc/ducla/agent.yaml \
  ducla-agent
```

### Kubernetes

```bash
# Deploy với kubectl
kubectl apply -f k8s/

# Hoặc với Helm
helm install ducla-agent ./helm-chart
```

## Bảo Mật

### Khuyến Nghị

1. **Sử dụng TLS**: Bật TLS cho tất cả API endpoints
2. **Authentication**: Bật JWT authentication
3. **RBAC**: Cấu hình role-based access control
4. **Firewall**: Chỉ mở các port cần thiết
5. **Audit Logging**: Bật audit logging cho tracking
6. **Regular Updates**: Cập nhật thường xuyên

### Cấu Hình TLS

```yaml
api:
  http:
    tls:
      enabled: true
      cert_file: "/etc/ducla/tls/server.crt"
      key_file: "/etc/ducla/tls/server.key"
      ca_file: "/etc/ducla/tls/ca.crt"
```

## Hỗ Trợ

- **Documentation**: [docs/](docs/)
- **Issues**: Tạo issue trên GitHub repository
- **Community**: Tham gia Discord/Slack community

## License

MIT License - xem file [LICENSE](LICENSE) để biết thêm chi tiết.
