# 🌐 Ducla Cloud Agent - REST API Reference

## 📋 Overview

Ducla Cloud Agent cung cấp comprehensive REST API để quản lý và monitor agent từ xa. API được chia thành 3 endpoints chính với các ports khác nhau.

## 🔌 API Endpoints

### 🚀 Main HTTP API (Port 8080)
```bash
# Base URL
http://localhost:8080
```

### 🏥 Health Check API (Port 8081)
```bash
# Base URL  
http://localhost:8081
```

### 📊 Metrics API (Port 9090)
```bash
# Base URL
http://localhost:9090
```

---

## 📡 Main HTTP API Commands (Port 8080)

### Agent Status & Information

#### Get Agent Status
```bash
curl http://localhost:8080/api/v1/status
```
**Response:**
```json
{
  "success": true,
  "data": {
    "running": true,
    "tasks": {
      "total_tasks": 0,
      "running_tasks": 0,
      "completed_tasks": 0,
      "queue_size": 0,
      "worker_count": 5
    },
    "metrics": {}
  }
}
```

#### Get Agent Configuration
```bash
curl http://localhost:8080/api/v1/config
```

#### Get Agent Version
```bash
curl http://localhost:8080/api/v1/version
```

### Task Management

#### List All Tasks
```bash
curl http://localhost:8080/api/v1/tasks
```

#### List Running Tasks Only
```bash
curl http://localhost:8080/api/v1/tasks?filter=running
```

#### Create New Task
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "type": "shell",
    "name": "example-task",
    "command": "echo",
    "args": ["Hello World"]
  }'
```

#### Get Task Details
```bash
curl http://localhost:8080/api/v1/tasks/{task-id}
```

#### Cancel Task
```bash
curl -X DELETE http://localhost:8080/api/v1/tasks/{task-id}
```

### File Operations

#### List Files
```bash
curl -X POST http://localhost:8080/api/v1/files \
  -H "Content-Type: application/json" \
  -d '{
    "type": "list",
    "source_path": "/tmp"
  }'
```

#### Copy File
```bash
curl -X POST http://localhost:8080/api/v1/files \
  -H "Content-Type: application/json" \
  -d '{
    "type": "copy",
    "source_path": "/tmp/source.txt",
    "dest_path": "/tmp/dest.txt"
  }'
```

#### Move File
```bash
curl -X POST http://localhost:8080/api/v1/files \
  -H "Content-Type: application/json" \
  -d '{
    "type": "move",
    "source_path": "/tmp/old.txt",
    "dest_path": "/tmp/new.txt"
  }'
```

#### Delete File
```bash
curl -X POST http://localhost:8080/api/v1/files \
  -H "Content-Type: application/json" \
  -d '{
    "type": "delete",
    "source_path": "/tmp/unwanted.txt"
  }'
```

#### Get File Stats
```bash
curl -X POST http://localhost:8080/api/v1/files \
  -H "Content-Type: application/json" \
  -d '{
    "type": "stat",
    "source_path": "/tmp/file.txt"
  }'
```

---

## 🏥 Health Check API Commands (Port 8081)

### System Health
```bash
curl http://localhost:8081/health
```
**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2025-11-09T12:00:00Z",
  "checks": {
    "system": "healthy",
    "disk": "healthy", 
    "memory": "healthy",
    "cpu": "healthy"
  },
  "summary": {
    "total": 4,
    "healthy": 4,
    "unhealthy": 0,
    "degraded": 0,
    "unknown": 0
  }
}
```

### Readiness Check
```bash
curl http://localhost:8081/ready
```

### Liveness Check
```bash
curl http://localhost:8081/live
```

---

## 📊 Metrics API Commands (Port 9090)

### Prometheus Metrics
```bash
curl http://localhost:9090/metrics
```
**Sample Output:**
```
# HELP ducla_agent_uptime_seconds Agent uptime in seconds
# TYPE ducla_agent_uptime_seconds counter
ducla_agent_uptime_seconds 3600

# HELP ducla_tasks_total Total number of tasks processed
# TYPE ducla_tasks_total counter
ducla_tasks_total{status="completed"} 10
ducla_tasks_total{status="failed"} 2

# HELP ducla_system_cpu_usage CPU usage percentage
# TYPE ducla_system_cpu_usage gauge
ducla_system_cpu_usage 15.5

# HELP ducla_system_memory_usage Memory usage percentage
# TYPE ducla_system_memory_usage gauge
ducla_system_memory_usage 45.2
```

### Custom Metrics Endpoint
```bash
curl http://localhost:9090/api/v1/metrics
```

### System Uptime
```bash
curl http://localhost:9090/api/v1/uptime
```

---

## 🔧 Advanced API Usage

### Authentication (if enabled)
```bash
# With JWT token
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:8080/api/v1/status
```

### Batch Operations
```bash
# Multiple file operations
curl -X POST http://localhost:8080/api/v1/files/batch \
  -H "Content-Type: application/json" \
  -d '{
    "operations": [
      {
        "type": "copy",
        "source_path": "/tmp/file1.txt",
        "dest_path": "/backup/file1.txt"
      },
      {
        "type": "copy", 
        "source_path": "/tmp/file2.txt",
        "dest_path": "/backup/file2.txt"
      }
    ]
  }'
```

### Streaming Endpoints
```bash
# Stream task logs
curl http://localhost:8080/api/v1/tasks/{task-id}/logs/stream

# Stream metrics
curl http://localhost:9090/api/v1/metrics/stream
```

---

## 📤 Output Destinations

Ducla Cloud Agent hỗ trợ gửi logs và metrics đến nhiều destinations khác nhau:

### 🔥 Prometheus
```bash
# Cấu hình Prometheus scraping
prometheus-o prometheus -p host=10.0.0.5 -p port=9090 -p format=metrics
```

### 🌐 HTTP/HTTPS
```bash
# Gửi qua HTTP POST
http-o http -p host=10.0.0.5 -p port=8080 -p format=json

# Gửi qua HTTPS với authentication
https-o https -p host=api.example.com -p port=443 -p format=json -p auth=bearer:TOKEN
```

### 🔌 TCP/UDP
```bash
# Gửi log qua TCP
tcp-o tcp -p host=10.0.0.10 -p port=9000

# Gửi log qua UDP
udp-o udp -p host=10.0.0.10 -p port=514
```

### 🔍 Elasticsearch
```bash
# Gửi sang Elasticsearch
es-o es -p host=10.0.0.10 -p port=9200 -p index=ducla-logs

# Với authentication
es-o es -p host=elastic.example.com -p port=9200 -p index=logs -p user=admin -p password=secret
```

### 📨 Apache Kafka
```bash
# Gửi sang Kafka
kafka-o kafka -p brokers=localhost:9092 -p topics=logs

# Multiple brokers
kafka-o kafka -p brokers=broker1:9092,broker2:9092 -p topics=ducla-logs,system-logs
```

### ☁️ AWS CloudWatch Logs
```bash
# Gửi AWS CloudWatch
cloudwatch_logs-o cloudwatch_logs -p log_group_name=ducla-agent -p region=ap-southeast-1

# Với custom stream
cloudwatch_logs-o cloudwatch_logs -p log_group_name=my-app -p log_stream_name=agent-001 -p region=us-east-1
```

### 🪣 AWS S3
```bash
# Gửi sang S3
s3-o s3 -p bucket=ducla-logs -p region=ap-southeast-1 -p total_file_size=5M

# Với custom prefix
s3-o s3 -p bucket=mybucket -p region=us-west-2 -p prefix=logs/agent/ -p total_file_size=10M
```

### 🐕 Datadog
```bash
# Gửi log sang Datadog
datadog-o datadog -p apikey=YOUR_API_KEY

# Với custom tags
datadog-o datadog -p apikey=YOUR_API_KEY -p tags=env:prod,service:ducla-agent
```

### 📊 Grafana Loki
```bash
# Gửi Promtail/Loki
loki-o loki -p host=http://10.0.0.5:3100

# Với labels
loki-o loki -p host=http://loki.example.com:3100 -p labels=job:ducla-agent,env:production
```

### 🏗️ Google Cloud Stackdriver
```bash
# Gửi Google Stackdriver
stackdriver-o stackdriver -p resource=k8s_container

# Với custom resource
stackdriver-o stackdriver -p resource=gce_instance -p project_id=my-project
```

### 🗃️ InfluxDB
```bash
# Gửi InfluxDB
influxdb-o influxdb -p host=10.0.0.5 -p port=8086 -p database=ducla

# InfluxDB v2
influxdb2-o influxdb2 -p host=influx.example.com -p org=myorg -p bucket=logs -p token=YOUR_TOKEN
```

### 📧 Email Notifications
```bash
# Gửi email alerts
email-o email -p smtp_host=smtp.gmail.com -p smtp_port=587 -p to=admin@example.com

# Với authentication
email-o email -p smtp_host=mail.example.com -p user=alerts@example.com -p password=secret
```

### 💬 Slack/Discord
```bash
# Gửi Slack
slack-o slack -p webhook_url=https://hooks.slack.com/services/YOUR/WEBHOOK/URL

# Gửi Discord
discord-o discord -p webhook_url=https://discord.com/api/webhooks/YOUR/WEBHOOK
```

### 🚫 Debug/Null Output
```bash
# Bỏ qua output (debug flow)
null-o null

# Debug output (console)
stdout-o stdout -p format=json
```

---

## 🔧 Configuration Examples

### Basic HTTP Output
```yaml
outputs:
  - name: http_endpoint
    type: http
    config:
      host: "10.0.0.5"
      port: 8080
      format: "json"
      endpoint: "/api/v1/logs"
```

### Multiple Outputs
```yaml
outputs:
  - name: prometheus
    type: prometheus
    config:
      host: "monitoring.example.com"
      port: 9090
      
  - name: elasticsearch
    type: elasticsearch
    config:
      host: "elastic.example.com"
      port: 9200
      index: "ducla-logs"
      
  - name: slack_alerts
    type: slack
    config:
      webhook_url: "https://hooks.slack.com/services/YOUR/WEBHOOK"
      channel: "#alerts"
```

---

## 🚀 Quick Start Examples

### Monitor Agent Status
```bash
#!/bin/bash
# Monitor script
while true; do
  echo "=== Agent Status ==="
  curl -s http://localhost:8080/api/v1/status | jq '.data'
  
  echo "=== Health Check ==="
  curl -s http://localhost:8081/health | jq '.status'
  
  sleep 30
done
```

### Automated Task Creation
```bash
#!/bin/bash
# Create daily backup task
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "type": "shell",
    "name": "daily-backup",
    "command": "tar",
    "args": ["-czf", "/backup/daily-$(date +%Y%m%d).tar.gz", "/data"]
  }'
```

### Metrics Collection
```bash
#!/bin/bash
# Collect and forward metrics
curl -s http://localhost:9090/metrics | \
  curl -X POST http://prometheus.example.com:9091/api/v1/import/prometheus \
    -H "Content-Type: text/plain" \
    --data-binary @-
```

---

## 📚 API Response Formats

### Success Response
```json
{
  "success": true,
  "data": {
    // Response data
  },
  "message": "Operation completed successfully"
}
```

### Error Response
```json
{
  "success": false,
  "error": "Error message",
  "code": "ERROR_CODE",
  "details": {
    // Additional error details
  }
}
```

### Pagination
```json
{
  "success": true,
  "data": {
    "items": [...],
    "pagination": {
      "page": 1,
      "per_page": 20,
      "total": 100,
      "total_pages": 5
    }
  }
}
```

---

## 🔐 Security Considerations

### API Authentication
- JWT tokens for API access
- API key authentication
- IP-based access control
- Rate limiting

### TLS Configuration
```yaml
api:
  http:
    tls:
      enabled: true
      cert_file: "/etc/ducla/tls/server.crt"
      key_file: "/etc/ducla/tls/server.key"
```

### Firewall Rules
```bash
# Allow API access
sudo ufw allow 8080/tcp  # HTTP API
sudo ufw allow 8081/tcp  # Health checks
sudo ufw allow 9090/tcp  # Metrics
```

---

## 🎯 Best Practices

1. **Use Health Checks**: Monitor `/health` endpoint regularly
2. **Implement Retries**: Handle temporary failures gracefully
3. **Monitor Metrics**: Set up Prometheus scraping
4. **Secure APIs**: Use TLS and authentication in production
5. **Log Everything**: Enable comprehensive logging
6. **Test Outputs**: Verify all output destinations work correctly

---

## 📞 Support

- **Documentation**: [README-VI.md](README-VI.md)
- **Issues**: [GitHub Issues](https://github.com/duclacloud/DUCLA-CLOUD-AGENT/issues)
- **API Questions**: Create issue with `api` label