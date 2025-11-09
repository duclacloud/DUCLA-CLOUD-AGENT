# 📖 Hướng Dẫn Sử Dụng Ducla Cloud Agent

## 📋 Mục Lục

0. [Giới Thiệu Ducla Cloud Agent](#0-giới-thiệu-ducla-cloud-agent)
1. [Cài Đặt Ducla Agent](#1-cài-đặt-ducla-agent)
2. [Các Lệnh Cơ Bản](#2-các-lệnh-cơ-bản-check-ducla-agent)
3. [Chỉnh Sửa File Config](#3-chỉnh-sửa-tham-số-file-config)
4. [Output và Kết Nối Hệ Thống](#4-các-lệnh-output-và-kết-nối-hệ-thống)
5. [Tích Hợp AI Systems và Workflow Automation](#5-tích-hợp-ai-systems-và-workflow-automation)

---

## 0. Giới Thiệu Ducla Cloud Agent

### 🚀 Ducla Cloud Agent là gì?

**Ducla Cloud Agent** là một công cụ monitoring và automation mạnh mẽ được thiết kế để quản lý và giám sát các hệ thống cloud infrastructure. Agent hoạt động như một cầu nối thông minh giữa các services, applications và các hệ thống monitoring/logging khác nhau.

### 🎯 Mục Đích và Ứng Dụng

#### **Infrastructure Monitoring**

- Giám sát real-time các metrics hệ thống (CPU, Memory, Disk, Network)
- Thu thập và phân tích logs từ applications và services
- Monitoring health status của các microservices
- Alerting và notification khi có sự cố

#### **Task Automation**

- Thực thi các automation tasks theo schedule
- Batch processing và data pipeline management
- Deployment automation và CI/CD integration
- System maintenance và cleanup tasks

#### **Data Integration**

- Kết nối và đồng bộ dữ liệu giữa các hệ thống khác nhau
- Export metrics đến Prometheus, Grafana, InfluxDB
- Stream logs đến Elasticsearch, Kafka, AWS CloudWatch
- Backup và archiving dữ liệu đến AWS S3, Google Cloud Storage

#### **AI-Powered Intelligence**

- Tích hợp với ChatGPT, Claude, Amazon Q để phân tích hệ thống thông minh
- AI-powered troubleshooting và root cause analysis
- Automated insights và performance optimization recommendations
- Natural language interface để query và điều khiển hệ thống
- Multi-AI comparison và consensus-based decision making

#### **Workflow Automation với N8N**

- Visual workflow automation với drag-and-drop interface
- Automated incident response workflows
- Multi-step data processing pipelines
- Integration với 200+ services và APIs
- Event-driven automation và real-time triggers
- AI-enhanced workflows với ChatGPT, Claude nodes

### 🏗️ Kiến Trúc và Thành Phần

#### **Core Components**

```
┌─────────────────────────────────────────────────────────┐
│                 Ducla Cloud Agent                       │
├─────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  │   API       │  │  Executor   │  │   Metrics   │     │
│  │  Server     │  │   Engine    │  │ Collector   │     │
│  │ (REST/gRPC) │  │ (Workers)   │  │(Prometheus) │     │
│  └─────────────┘  └─────────────┘  └─────────────┘     │
├─────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  │   Config    │  │   Storage   │  │   Output    │     │
│  │  Manager    │  │   Engine    │  │  Handlers   │     │
│  │   (YAML)    │  │ (Local/S3)  │  │(Multi-dest) │     │
│  └─────────────┘  └─────────────┘  └─────────────┘     │
└─────────────────────────────────────────────────────────┘
```

#### **Key Features**

🔧 **Multi-Protocol Support**

- REST API (HTTP/HTTPS) cho web integration
- gRPC cho high-performance communication
- WebSocket cho real-time streaming
- CLI interface cho system administration

📊 **Rich Metrics Collection**

- System metrics (CPU, Memory, Disk I/O, Network)
- Application metrics (custom counters, gauges, histograms)
- Business metrics (KPIs, SLAs, performance indicators)
- Infrastructure metrics (containers, VMs, cloud resources)

🔄 **Flexible Task Execution**

- Multi-threaded worker pool với configurable size
- Task queuing và priority management
- Retry logic với exponential backoff
- Timeout handling và resource cleanup

🤖 **AI-Powered Intelligence**

- **ChatGPT Integration**: Performance analysis, troubleshooting insights
- **Claude Integration**: Security analysis, compliance checking
- **Amazon Q Integration**: AWS-specific recommendations và cost optimization
- **Natural Language Interface**: Chat với hệ thống bằng tiếng tự nhiên
- **Multi-AI Consensus**: Kết hợp insights từ multiple AI systems

🔄 **N8N Workflow Automation**

- **Visual Workflows**: Drag-and-drop workflow designer
- **200+ Integrations**: Slack, GitHub, AWS, Google Cloud, databases
- **Event-Driven**: Webhooks, schedules, file watchers, API triggers
- **AI-Enhanced Nodes**: ChatGPT, Claude, OpenAI nodes trong workflows
- **Error Handling**: Retry logic, fallbacks, notifications

🌐 **Multi-Destination Output**

- **Monitoring**: Prometheus, Grafana, InfluxDB, Datadog
- **Logging**: Elasticsearch, Splunk, Fluentd, Logstash
- **Streaming**: Kafka, RabbitMQ, AWS Kinesis, Google Pub/Sub
- **Storage**: AWS S3, Google Cloud Storage, Azure Blob
- **Alerting**: PagerDuty, Slack, Email, Webhook
- **AI Systems**: ChatGPT, Claude, Amazon Q, custom AI endpoints

### 🎨 Use Cases và Scenarios

#### **1. Microservices Monitoring**

```yaml
# Scenario: E-commerce platform với 50+ microservices
services:
  - user-service (port 8001)
  - product-service (port 8002)
  - order-service (port 8003)
  - payment-service (port 8004)

monitoring:
  - Health checks mỗi 30 giây
  - Performance metrics collection
  - Error rate và latency tracking
  - Auto-scaling triggers
```

#### **2. Log Aggregation và Analysis**

```yaml
# Scenario: Multi-region application logs
sources:
  - Application logs (JSON format)
  - Nginx access logs
  - Database slow query logs
  - System audit logs

destinations:
  - Elasticsearch cluster (search & analytics)
  - AWS S3 (long-term storage)
  - Slack (critical alerts)
```

#### **3. DevOps Automation**

```yaml
# Scenario: CI/CD pipeline integration
tasks:
  - Code deployment verification
  - Database migration execution
  - Cache warming và preloading
  - Performance testing automation
  - Rollback procedures
```

#### **4. Cloud Cost Optimization**

```yaml
# Scenario: AWS resource monitoring
metrics:
  - EC2 instance utilization
  - RDS connection counts
  - S3 storage usage
  - Lambda execution costs

actions:
  - Auto-shutdown unused instances
  - Resize over/under-provisioned resources
  - Archive old data to cheaper storage
```

### 🔒 Security và Compliance

#### **Security Features**

- **Authentication**: API keys, JWT tokens, OAuth2
- **Authorization**: Role-based access control (RBAC)
- **Encryption**: TLS/SSL for data in transit, AES-256 for data at rest
- **Audit Logging**: Complete audit trail của tất cả operations
- **Network Security**: IP whitelisting, VPN support, private networking

#### **Compliance Support**

- **GDPR**: Data privacy và right to be forgotten
- **SOC 2**: Security controls và monitoring
- **HIPAA**: Healthcare data protection
- **PCI DSS**: Payment card industry standards

### 📈 Performance và Scalability

#### **Performance Characteristics**

- **Throughput**: 10,000+ events/second per agent
- **Latency**: < 10ms average processing time
- **Memory**: 50MB baseline, scales with workload
- **CPU**: Multi-core utilization với worker pools
- **Storage**: Configurable local buffering và compression

#### **Scalability Options**

- **Horizontal**: Multiple agents với load balancing
- **Vertical**: Configurable worker counts và resource limits
- **Cloud-native**: Kubernetes deployment với auto-scaling
- **Edge deployment**: Lightweight mode cho IoT và edge devices

### 🛠️ Technology Stack

#### **Core Technologies**

- **Language**: Go (Golang) - performance và concurrency
- **HTTP Framework**: Gin/Echo - fast HTTP routing
- **gRPC**: Protocol Buffers - efficient binary communication
- **Database**: BadgerDB/BoltDB - embedded key-value storage
- **Metrics**: Prometheus client libraries
- **Configuration**: YAML/JSON với validation

#### **Integration Libraries**

- **AWS SDK**: S3, CloudWatch, Kinesis, SQS
- **Kubernetes Client**: Pod management, service discovery
- **Docker API**: Container monitoring và management
- **Database Drivers**: PostgreSQL, MySQL, MongoDB, Redis

### 🌟 Competitive Advantages

#### **So với các giải pháp khác**

**vs. Prometheus Node Exporter**

- ✅ Multi-destination output (không chỉ Prometheus)
- ✅ Task execution capabilities
- ✅ Built-in API server
- ✅ Advanced configuration management
- ✅ **AI-powered analysis** và insights

**vs. Fluentd/Fluent Bit**

- ✅ Metrics collection (không chỉ logs)
- ✅ Task automation features
- ✅ REST API interface
- ✅ Better performance với Go
- ✅ **ChatGPT/Claude integration** cho log analysis

**vs. Telegraf**

- ✅ Task execution engine
- ✅ Advanced retry logic
- ✅ Built-in health checks
- ✅ More flexible configuration
- ✅ **N8N workflow automation** integration

**vs. Datadog/New Relic Agents**

- ✅ Open source và self-hosted
- ✅ **Multi-AI support** (ChatGPT + Claude + Amazon Q)
- ✅ **N8N visual workflows** cho automation
- ✅ No vendor lock-in
- ✅ **Natural language interface** cho system queries

**vs. Custom Solutions**

- ✅ Production-ready với comprehensive testing
- ✅ Professional documentation và support
- ✅ Security best practices built-in
- ✅ Regular updates và maintenance
- ✅ **AI-first architecture** với built-in intelligence

### 🎯 Target Users

#### **DevOps Engineers**

- Infrastructure monitoring và automation
- CI/CD pipeline integration
- Incident response và troubleshooting

#### **Site Reliability Engineers (SRE)**

- Service level monitoring
- Capacity planning và performance optimization
- Disaster recovery và business continuity

#### **Platform Engineers**

- Multi-tenant infrastructure management
- Developer productivity tools
- Internal platform services

#### **System Administrators**

- Server monitoring và maintenance
- Log management và analysis
- Security compliance và auditing

---

## 1. Cài Đặt Ducla Agent

### 🐧 Ubuntu/Debian Systems

#### Cài đặt từ DEB package

```bash
# Tải package
wget https://github.com/duclacloud/DUCLA-CLOUD-AGENT/releases/download/v1.0.0/ducla-agent_1.0.0_amd64.deb

# Cài đặt
sudo dpkg -i ducla-agent_1.0.0_amd64.deb

# Khởi động service
sudo systemctl enable ducla-agent
sudo systemctl start ducla-agent
```

### 🎩 RHEL/CentOS/Fedora Systems

#### Cài đặt từ RPM package

```bash
# Tải package
wget https://github.com/duclacloud/DUCLA-CLOUD-AGENT/releases/download/v1.0.0/ducla-agent-1.0.0-1.x86_64.rpm

# Cài đặt
sudo rpm -ivh ducla-agent-1.0.0-1.x86_64.rpm

# Khởi động service
sudo systemctl enable ducla-agent
sudo systemctl start ducla-agent
```

### 🔧 Manual Installation (Binary)

#### Cài đặt từ binary

```bash
# Tải binary
wget https://github.com/duclacloud/DUCLA-CLOUD-AGENT/releases/download/v1.0.0/ducla-agent-linux-amd64.tar.gz

# Giải nén
tar -xzf ducla-agent-linux-amd64.tar.gz

# Copy binary
sudo cp ducla-agent /usr/local/bin/
sudo chmod +x /usr/local/bin/ducla-agent

# Tạo user và directories
sudo useradd -r -s /bin/false ducla
sudo mkdir -p /etc/ducla /opt/ducla /var/log/ducla
sudo chown ducla:ducla /opt/ducla /var/log/ducla

# Tạo config file
sudo cp agent.yaml /etc/ducla/
sudo chown root:ducla /etc/ducla/agent.yaml
sudo chmod 640 /etc/ducla/agent.yaml
```

### ✅ Xác Nhận Cài Đặt

```bash
# Kiểm tra version
ducla-agent show version

# Kiểm tra service status
sudo systemctl status ducla-agent

# Kiểm tra API endpoints
curl http://localhost:8080/api/v1/status
curl -I http://localhost:8081/health
```

---

## 2. Các Lệnh Cơ Bản Check Ducla Agent

### 📊 Hiển Thị Thông Tin

#### Xem version và build info

```bash
ducla-agent show version
```

**Output:**

```
Ducla Cloud Agent
  Version:    1.0.0
  Build Time: 2025-11-09_05:21:13_UTC
  Git Commit: b83e291
  Go Version: go1.21.5
  OS/Arch:    linux/amd64
```

#### Xem cấu hình hiện tại

```bash
ducla-agent show config
```

#### Xem trạng thái agent

```bash
ducla-agent show status
```

### 🔍 Kiểm Tra Hệ Thống

#### Kiểm tra service systemd

```bash
# Trạng thái service
sudo systemctl status ducla-agent

# Xem logs
sudo journalctl -u ducla-agent -f

# Xem logs với filter
sudo journalctl -u ducla-agent --since "1 hour ago"
```

#### Kiểm tra API endpoints

```bash
# API chính (port 8080)
curl http://localhost:8080/api/v1/status | jq .

# Health check (port 8081)
curl -I http://localhost:8081/health

# Metrics Prometheus (port 9090)
curl http://localhost:9090/metrics | head -20
```

#### Kiểm tra processes và resources

```bash
# Xem process
ps aux | grep ducla-agent

# Xem memory usage
sudo systemctl show ducla-agent --property=MemoryCurrent

# Xem network connections
sudo netstat -tlnp | grep ducla-agent
```

### 🛠️ Quản Lý Service

#### Điều khiển service

```bash
# Khởi động
sudo systemctl start ducla-agent

# Dừng
sudo systemctl stop ducla-agent

# Restart
sudo systemctl restart ducla-agent

# Reload config
sudo systemctl reload ducla-agent

# Enable auto-start
sudo systemctl enable ducla-agent

# Disable auto-start
sudo systemctl disable ducla-agent
```

### 📋 Validation và Testing

#### Validate config file

```bash
ducla-agent config validate
ducla-agent config validate -config /path/to/custom.yaml
```

#### Test connectivity

```bash
ducla-agent config test
```

---

## 3. Chỉnh Sửa Tham Số File Config

### 📁 Vị Trí Config File

**Default locations:**

- `/etc/ducla/agent.yaml` (system-wide)
- `./agent.yaml` (current directory)
- `~/.ducla/agent.yaml` (user home)

### ⚙️ Cấu Trúc Config File

#### File config mẫu (`/etc/ducla/agent.yaml`):

```yaml
# Ducla Cloud Agent Configuration
agent:
  id: "ducla-agent"
  name: "Ducla Cloud Agent"
  environment: "production" # development, staging, production
  region: "us-east-1"

# API Configuration
api:
  http:
    enabled: true
    address: "127.0.0.1"
    port: 8080
    timeout: "30s"
  grpc:
    enabled: false
    address: "0.0.0.0"
    port: 8443
    tls_enabled: false

# Health Check Configuration
health:
  enabled: true
  address: "127.0.0.1"
  port: 8081

# Metrics Configuration
metrics:
  enabled: true
  address: "127.0.0.1"
  port: 9090
  path: "/metrics"

# Storage Configuration
storage:
  data_dir: "/opt/ducla/data"
  temp_dir: "/tmp/ducla"
  max_size: "10GB"

# Executor Configuration
executor:
  workers: 5
  queue_size: 100
  task_timeout: "5m"
  retry_attempts: 3

# Logging Configuration
logging:
  level: "info" # debug, info, warn, error
  format: "json" # json, text
  output: "stdout" # stdout, file, syslog
  file_path: "/var/log/ducla/agent.log"
  max_size: "100MB"
  max_backups: 5
  max_age: 30

# Security Configuration
security:
  api_key: ""
  tls:
    enabled: false
    cert_file: ""
    key_file: ""
    ca_file: ""
```

### 🔧 Các Tham Số Quan Trọng

#### Agent Settings

```yaml
agent:
  id: "my-agent-001" # Unique agent identifier
  name: "Production Agent" # Human readable name
  environment: "production" # Environment tag
  region: "ap-southeast-1" # AWS region or location
  tags: # Custom tags
    - "web-server"
    - "monitoring"
```

#### Performance Tuning

```yaml
executor:
  workers: 10 # Số worker threads (default: 5)
  queue_size: 500 # Kích thước queue (default: 100)
  task_timeout: "10m" # Timeout cho tasks (default: 5m)
  retry_attempts: 5 # Số lần retry (default: 3)
  batch_size: 50 # Batch processing size
```

#### Resource Limits

```yaml
storage:
  data_dir: "/opt/ducla/data"
  temp_dir: "/tmp/ducla"
  max_size: "50GB" # Giới hạn storage
  cleanup_interval: "1h" # Tần suất cleanup

memory:
  max_usage: "2GB" # Giới hạn memory
  gc_percent: 100 # Go GC tuning
```

### 🔄 Reload Config

#### Sau khi chỉnh sửa config:

```bash
# Validate config trước
ducla-agent config validate

# Reload service
sudo systemctl reload ducla-agent

# Hoặc restart nếu cần
sudo systemctl restart ducla-agent

# Kiểm tra config mới
ducla-agent show config
```

---

## 4. Các Lệnh Output và Kết Nối Hệ Thống

### 📤 Output Destinations

Ducla Agent hỗ trợ gửi dữ liệu đến nhiều hệ thống khác nhau:

#### Cấu hình trong `agent.yaml`:

```yaml
outputs:
  # Prometheus Metrics
  prometheus:
    enabled: true
    endpoint: "http://prometheus:9090/api/v1/write"
    interval: "30s"

  # Elasticsearch
  elasticsearch:
    enabled: false
    hosts: ["http://elasticsearch:9200"]
    index: "ducla-agent"

  # AWS S3
  s3:
    enabled: false
    bucket: "my-ducla-logs"
    region: "us-east-1"
    prefix: "agent-logs/"

  # Kafka
  kafka:
    enabled: false
    brokers: ["kafka:9092"]
    topic: "ducla-events"

  # InfluxDB
  influxdb:
    enabled: false
    url: "http://influxdb:8086"
    database: "ducla"

  # Grafana
  grafana:
    enabled: false
    url: "http://grafana:3000"
    api_key: "your-api-key"
```

### 🔗 Kết Nối AWS S3

#### 1. Cấu hình S3 Output

```yaml
outputs:
  s3:
    enabled: true
    bucket: "my-company-ducla-logs"
    region: "ap-southeast-1"
    prefix: "production/agent-logs/"
    access_key_id: "AKIA..." # Hoặc dùng IAM role
    secret_access_key: "..." # Hoặc dùng IAM role
    session_token: "" # Nếu dùng temporary credentials

    # Advanced settings
    upload_interval: "5m" # Tần suất upload
    batch_size: 1000 # Số records per batch
    compression: "gzip" # Nén dữ liệu
    encryption: "AES256" # Mã hóa S3

    # File naming
    filename_template: "ducla-{date}-{time}-{hostname}.json.gz"
    date_format: "2006-01-02"
    time_format: "15-04-05"
```

#### 2. AWS Credentials Setup

```bash
# Option 1: AWS CLI
aws configure
AWS Access Key ID: AKIA...
AWS Secret Access Key: ...
Default region name: ap-southeast-1
Default output format: json

# Option 2: Environment variables
export AWS_ACCESS_KEY_ID="AKIA..."
export AWS_SECRET_ACCESS_KEY="..."
export AWS_DEFAULT_REGION="ap-southeast-1"

# Option 3: IAM Role (recommended for EC2)
# Attach IAM role với S3 permissions đến EC2 instance
```

#### 3. S3 Bucket Policy

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::ACCOUNT:user/ducla-agent"
      },
      "Action": ["s3:PutObject", "s3:PutObjectAcl", "s3:GetObject"],
      "Resource": "arn:aws:s3:::my-company-ducla-logs/*"
    }
  ]
}
```

#### 4. Test S3 Connection

```bash
# Test S3 connectivity
ducla-agent output test s3

# Upload test file
ducla-agent output upload s3 --file test.json

# List S3 objects
aws s3 ls s3://my-company-ducla-logs/production/agent-logs/
```

### 📊 Kết Nối Prometheus

#### 1. Cấu hình Prometheus Output

```yaml
outputs:
  prometheus:
    enabled: true

    # Push Gateway mode
    push_gateway:
      enabled: true
      url: "http://prometheus-pushgateway:9091"
      job_name: "ducla-agent"
      instance: "agent-001"
      push_interval: "30s"

    # Remote Write mode
    remote_write:
      enabled: false
      url: "http://prometheus:9090/api/v1/write"
      timeout: "30s"

    # Metrics configuration
    metrics:
      prefix: "ducla_"
      labels:
        environment: "production"
        region: "ap-southeast-1"

    # Custom metrics
    custom_metrics:
      - name: "task_duration_seconds"
        type: "histogram"
        help: "Task execution duration"
        buckets: [0.1, 0.5, 1.0, 2.5, 5.0, 10.0]

      - name: "task_total"
        type: "counter"
        help: "Total number of tasks"
        labels: ["status", "type"]
```

#### 2. Prometheus Server Configuration

```yaml
# prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: "ducla-agent"
    static_configs:
      - targets: ["ducla-agent:9090"]
    scrape_interval: 30s
    metrics_path: /metrics

  - job_name: "ducla-pushgateway"
    static_configs:
      - targets: ["pushgateway:9091"]
```

#### 3. Test Prometheus Connection

```bash
# Test metrics endpoint
curl http://localhost:9090/metrics

# Test push gateway
ducla-agent output test prometheus

# Push custom metric
ducla-agent metrics push --name "test_metric" --value 42 --labels "env=test"
```

### 🔍 Kết Nối Elasticsearch

#### 1. Cấu hình Elasticsearch Output

```yaml
outputs:
  elasticsearch:
    enabled: true
    hosts:
      - "http://elasticsearch-01:9200"
      - "http://elasticsearch-02:9200"

    # Index configuration
    index: "ducla-agent-{+yyyy.MM.dd}"
    doc_type: "_doc"

    # Authentication
    username: "ducla_user"
    password: "secure_password"

    # SSL/TLS
    ssl:
      enabled: true
      ca_file: "/etc/ssl/certs/elasticsearch-ca.pem"
      cert_file: "/etc/ssl/certs/client.pem"
      key_file: "/etc/ssl/private/client-key.pem"

    # Performance
    bulk_size: 1000
    flush_interval: "30s"
    timeout: "60s"

    # Template
    template:
      enabled: true
      name: "ducla-agent"
      pattern: "ducla-agent-*"
      settings:
        number_of_shards: 1
        number_of_replicas: 1
```

#### 2. Elasticsearch Index Template

```json
{
  "index_patterns": ["ducla-agent-*"],
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 1,
    "index.refresh_interval": "30s"
  },
  "mappings": {
    "properties": {
      "@timestamp": { "type": "date" },
      "agent_id": { "type": "keyword" },
      "level": { "type": "keyword" },
      "message": { "type": "text" },
      "task_id": { "type": "keyword" },
      "duration": { "type": "long" },
      "status": { "type": "keyword" }
    }
  }
}
```

#### 3. Test Elasticsearch Connection

```bash
# Test connection
ducla-agent output test elasticsearch

# Send test document
ducla-agent output send elasticsearch --data '{"test": "message"}'

# Query Elasticsearch
curl -X GET "elasticsearch:9200/ducla-agent-*/_search?pretty"
```

### 📨 Kết Nối Kafka

#### 1. Cấu hình Kafka Output

```yaml
outputs:
  kafka:
    enabled: true
    brokers:
      - "kafka-01:9092"
      - "kafka-02:9092"
      - "kafka-03:9092"

    # Topic configuration
    topic: "ducla-events"
    partition_key: "agent_id"

    # Security
    security:
      protocol: "SASL_SSL"
      mechanism: "PLAIN"
      username: "ducla_producer"
      password: "secure_password"

    # SSL
    ssl:
      ca_file: "/etc/ssl/certs/kafka-ca.pem"
      cert_file: "/etc/ssl/certs/kafka-client.pem"
      key_file: "/etc/ssl/private/kafka-client-key.pem"

    # Producer settings
    producer:
      acks: "all"
      retries: 3
      batch_size: 16384
      linger_ms: 5
      buffer_memory: 33554432
      compression_type: "gzip"
```

#### 2. Test Kafka Connection

```bash
# Test connection
ducla-agent output test kafka

# Send test message
ducla-agent output send kafka --message '{"event": "test"}'

# Consume messages (for testing)
kafka-console-consumer --bootstrap-server kafka:9092 --topic ducla-events --from-beginning
```

### 🔧 CLI Commands cho Output Management

#### Quản lý outputs

```bash
# List all configured outputs
ducla-agent output list

# Test specific output
ducla-agent output test s3
ducla-agent output test prometheus
ducla-agent output test elasticsearch

# Enable/disable outputs
ducla-agent output enable s3
ducla-agent output disable elasticsearch

# Send test data
ducla-agent output send s3 --file test.json
ducla-agent output send prometheus --metric test_metric=42

# Check output status
ducla-agent output status
ducla-agent output status s3

# View output statistics
ducla-agent output stats
ducla-agent output stats --output prometheus --since "1h"
```

#### Monitoring outputs

```bash
# Real-time output monitoring
ducla-agent output monitor

# Output logs
ducla-agent output logs s3
ducla-agent output logs --follow

# Output metrics
ducla-agent output metrics
ducla-agent output metrics prometheus
```

---

## 5. Tích Hợp AI Systems và Workflow Automation

### 🤖 AI Systems Integration

Ducla Cloud Agent có thể hoạt động như "chân tay" cho các hệ thống AI, cung cấp dữ liệu thực tế và thực thi các actions được AI đề xuất.

#### 🧠 Amazon Q Integration

##### Cấu hình Amazon Q Connector

```yaml
ai_integrations:
  amazon_q:
    enabled: true
    region: "us-east-1"
    application_id: "your-q-app-id"

    # Authentication
    credentials:
      access_key_id: "AKIA..."
      secret_access_key: "..."

    # Data streaming to Q
    data_sources:
      - type: "metrics"
        format: "json"
        interval: "5m"

      - type: "logs"
        format: "structured"
        level: "info"

      - type: "events"
        format: "cloudtrail"

    # Q Query Interface
    query_endpoint: "https://your-q-app.us-east-1.amazonaws.com/api/v1/query"

    # Auto-response to Q queries
    auto_response:
      enabled: true
      allowed_actions:
        - "get_system_status"
        - "get_metrics"
        - "list_services"
        - "check_health"
```

##### Amazon Q Use Cases

```bash
# Q có thể hỏi agent về system status
Q: "What's the current CPU usage on production servers?"
Agent Response: {
  "cpu_usage": "45%",
  "memory_usage": "67%",
  "disk_usage": "23%",
  "timestamp": "2025-11-09T12:00:00Z"
}

# Q có thể yêu cầu thực hiện actions
Q: "Can you restart the web service on server-01?"
Agent Action: ducla-agent task execute restart-service --target server-01 --service web

# Q có thể phân tích logs và đưa ra insights
Q: "Analyze error patterns in the last 24 hours"
Agent: Streams structured logs to Q for analysis
```

#### 🤖 ChatGPT/OpenAI Integration

##### Cấu hình OpenAI Connector

```yaml
ai_integrations:
  openai:
    enabled: true
    api_key: "sk-..."
    model: "gpt-4"

    # Context sharing
    context_sharing:
      system_metrics: true
      application_logs: true
      infrastructure_state: true

    # Automated insights
    insights:
      enabled: true
      schedule: "0 */6 * * *" # Every 6 hours
      topics:
        - "performance_analysis"
        - "error_pattern_detection"
        - "capacity_planning"
        - "security_anomalies"

    # Action execution
    action_approval:
      required: true
      timeout: "5m"
      allowed_actions:
        - "scale_services"
        - "restart_services"
        - "update_configs"
```

##### ChatGPT Integration Examples

```python
# Python script để ChatGPT tương tác với Ducla Agent
import requests
import openai

def get_system_status():
    """ChatGPT có thể gọi function này để lấy system status"""
    response = requests.get("http://ducla-agent:8080/api/v1/status")
    return response.json()

def analyze_with_gpt(system_data):
    """Gửi system data cho ChatGPT phân tích"""
    prompt = f"""
    Analyze this system data and provide insights:
    {system_data}

    Please identify:
    1. Performance bottlenecks
    2. Potential issues
    3. Optimization recommendations
    4. Action items
    """

    response = openai.ChatCompletion.create(
        model="gpt-4",
        messages=[{"role": "user", "content": prompt}]
    )

    return response.choices[0].message.content

# Workflow example
system_data = get_system_status()
insights = analyze_with_gpt(system_data)
print(f"AI Insights: {insights}")
```

#### 🧠 Claude/Anthropic Integration

##### Cấu hình Claude Connector

```yaml
ai_integrations:
  claude:
    enabled: true
    api_key: "sk-ant-..."
    model: "claude-3-opus"

    # Specialized for infrastructure analysis
    specialization:
      - "infrastructure_optimization"
      - "security_analysis"
      - "cost_optimization"
      - "compliance_checking"

    # Structured data format for Claude
    data_format:
      metrics: "prometheus_format"
      logs: "structured_json"
      events: "timeline_format"
```

### 🔄 N8N Workflow Integration

#### N8N Workflow Examples

##### 1. Automated Incident Response Workflow

```json
{
  "name": "Ducla Agent - Incident Response",
  "nodes": [
    {
      "name": "Ducla Webhook",
      "type": "n8n-nodes-base.webhook",
      "parameters": {
        "path": "ducla-alert",
        "httpMethod": "POST"
      }
    },
    {
      "name": "Parse Alert",
      "type": "n8n-nodes-base.function",
      "parameters": {
        "functionCode": "// Parse incoming alert from Ducla Agent\nconst alert = items[0].json;\nreturn [{\n  json: {\n    severity: alert.severity,\n    service: alert.service,\n    message: alert.message,\n    timestamp: alert.timestamp\n  }\n}];"
      }
    },
    {
      "name": "Check Severity",
      "type": "n8n-nodes-base.if",
      "parameters": {
        "conditions": {
          "string": [
            {
              "value1": "={{$json.severity}}",
              "operation": "equal",
              "value2": "critical"
            }
          ]
        }
      }
    },
    {
      "name": "Slack Alert",
      "type": "n8n-nodes-base.slack",
      "parameters": {
        "channel": "#alerts",
        "text": "🚨 Critical Alert: {{$json.message}}"
      }
    },
    {
      "name": "Auto Remediation",
      "type": "n8n-nodes-base.httpRequest",
      "parameters": {
        "method": "POST",
        "url": "http://ducla-agent:8080/api/v1/tasks",
        "body": {
          "task": "auto_remediate",
          "service": "={{$json.service}}",
          "action": "restart"
        }
      }
    }
  ]
}
```

##### 2. Performance Monitoring Workflow

```json
{
  "name": "Ducla Agent - Performance Monitor",
  "nodes": [
    {
      "name": "Schedule",
      "type": "n8n-nodes-base.cron",
      "parameters": {
        "triggerTimes": {
          "item": [
            {
              "mode": "everyMinute",
              "minute": 5
            }
          ]
        }
      }
    },
    {
      "name": "Get Metrics",
      "type": "n8n-nodes-base.httpRequest",
      "parameters": {
        "method": "GET",
        "url": "http://ducla-agent:8080/api/v1/metrics"
      }
    },
    {
      "name": "Analyze with ChatGPT",
      "type": "n8n-nodes-base.openAi",
      "parameters": {
        "operation": "chat",
        "model": "gpt-4",
        "messages": {
          "messageValues": [
            {
              "role": "user",
              "content": "Analyze these metrics and suggest optimizations: {{$json}}"
            }
          ]
        }
      }
    },
    {
      "name": "Store Insights",
      "type": "n8n-nodes-base.httpRequest",
      "parameters": {
        "method": "POST",
        "url": "http://ducla-agent:8080/api/v1/insights",
        "body": {
          "source": "chatgpt",
          "insights": "={{$json.choices[0].message.content}}",
          "timestamp": "={{new Date().toISOString()}}"
        }
      }
    }
  ]
}
```

##### 3. Multi-AI Analysis Workflow

```json
{
  "name": "Ducla Agent - Multi-AI Analysis",
  "nodes": [
    {
      "name": "Data Collection",
      "type": "n8n-nodes-base.httpRequest",
      "parameters": {
        "method": "GET",
        "url": "http://ducla-agent:8080/api/v1/comprehensive-data"
      }
    },
    {
      "name": "Split Analysis",
      "type": "n8n-nodes-base.splitInBatches",
      "parameters": {
        "batchSize": 1
      }
    },
    {
      "name": "ChatGPT Analysis",
      "type": "n8n-nodes-base.openAi",
      "parameters": {
        "operation": "chat",
        "model": "gpt-4",
        "messages": {
          "messageValues": [
            {
              "role": "system",
              "content": "You are a DevOps expert. Analyze system performance."
            },
            {
              "role": "user",
              "content": "{{$json.data}}"
            }
          ]
        }
      }
    },
    {
      "name": "Claude Analysis",
      "type": "n8n-nodes-base.httpRequest",
      "parameters": {
        "method": "POST",
        "url": "https://api.anthropic.com/v1/messages",
        "headers": {
          "x-api-key": "sk-ant-...",
          "content-type": "application/json"
        },
        "body": {
          "model": "claude-3-opus-20240229",
          "max_tokens": 1000,
          "messages": [
            {
              "role": "user",
              "content": "Focus on security and compliance analysis: {{$json.data}}"
            }
          ]
        }
      }
    },
    {
      "name": "Amazon Q Query",
      "type": "n8n-nodes-base.aws",
      "parameters": {
        "service": "qbusiness",
        "operation": "chatSync",
        "applicationId": "your-q-app-id",
        "userMessage": "Analyze this infrastructure data: {{$json.data}}"
      }
    },
    {
      "name": "Combine Insights",
      "type": "n8n-nodes-base.function",
      "parameters": {
        "functionCode": "// Combine insights from multiple AI systems\nconst chatgpt = items[0].json.choices[0].message.content;\nconst claude = items[1].json.content[0].text;\nconst amazonq = items[2].json.systemMessage;\n\nreturn [{\n  json: {\n    combined_analysis: {\n      chatgpt_insights: chatgpt,\n      claude_security: claude,\n      amazonq_recommendations: amazonq,\n      timestamp: new Date().toISOString()\n    }\n  }\n}];"
      }
    },
    {
      "name": "Execute Actions",
      "type": "n8n-nodes-base.httpRequest",
      "parameters": {
        "method": "POST",
        "url": "http://ducla-agent:8080/api/v1/ai-actions",
        "body": "={{$json.combined_analysis}}"
      }
    }
  ]
}
```

### 🔧 AI-Powered CLI Commands

#### Intelligent System Analysis

```bash
# AI-powered system analysis
ducla-agent ai analyze --provider chatgpt --focus performance
ducla-agent ai analyze --provider claude --focus security
ducla-agent ai analyze --provider amazonq --focus cost-optimization

# Multi-AI comparison
ducla-agent ai compare --providers "chatgpt,claude,amazonq" --topic "infrastructure-health"

# AI-suggested actions
ducla-agent ai suggest --context current-metrics --provider chatgpt
ducla-agent ai execute --suggestion-id "ai-001" --approve
```

#### Conversational Interface

```bash
# Natural language queries
ducla-agent chat "What's causing high CPU usage?"
ducla-agent chat "How can I optimize memory usage?"
ducla-agent chat "Show me security vulnerabilities"

# AI-powered troubleshooting
ducla-agent troubleshoot --ai-assistant chatgpt --issue "slow-response-time"
ducla-agent troubleshoot --ai-assistant claude --issue "memory-leak"
```

### 📊 AI Integration Examples

#### Real-world Scenario 1: E-commerce Platform

```yaml
# AI-powered e-commerce monitoring
scenario: "Black Friday Traffic Spike"

ai_workflow: 1. Ducla Agent collects real-time metrics
  2. ChatGPT analyzes traffic patterns
  3. Claude evaluates security risks
  4. Amazon Q suggests scaling strategies
  5. N8N orchestrates auto-scaling actions
  6. Slack notifications with AI insights

metrics_collected:
  - "requests_per_second: 15000"
  - "response_time_p95: 250ms"
  - "error_rate: 0.02%"
  - "cpu_usage: 78%"

ai_insights:
  chatgpt: "Traffic spike detected. Recommend horizontal scaling of web tier."
  claude: "No security anomalies. DDoS protection holding steady."
  amazonq: "Cost-optimal scaling: Add 3 instances for 2 hours."

actions_taken:
  - Auto-scale web servers (3 → 6 instances)
  - Enable CDN burst mode
  - Increase database connection pool
  - Alert on-call team
```

#### Real-world Scenario 2: DevOps Pipeline

```yaml
# AI-enhanced CI/CD pipeline
scenario: "Deployment Quality Gate"

ai_workflow: 1. Ducla Agent monitors deployment metrics
  2. AI systems analyze deployment health
  3. N8N workflow decides rollback/proceed
  4. Automated notifications with AI reasoning

deployment_metrics:
  - "deployment_time: 8m 32s"
  - "test_pass_rate: 98.5%"
  - "performance_regression: -2.1%"
  - "error_spike: false"

ai_analysis:
  chatgpt: "Deployment successful. Minor performance regression within acceptable range."
  claude: "Security scan passed. No new vulnerabilities introduced."
  amazonq: "Resource utilization optimal. No scaling needed."

decision: "PROCEED - All AI systems recommend deployment continuation"
```

### 🎯 Best Practices cho AI Integration

#### 1. Data Privacy và Security

```yaml
ai_security:
  data_anonymization: true
  pii_filtering: enabled
  encryption_in_transit: true
  audit_logging: comprehensive

  # Chỉ gửi metadata, không gửi sensitive data
  allowed_data_types:
    - "system_metrics"
    - "performance_counters"
    - "error_patterns"
    - "resource_utilization"

  blocked_data_types:
    - "user_data"
    - "payment_info"
    - "personal_identifiers"
    - "api_keys"
```

#### 2. Cost Optimization

```yaml
ai_cost_control:
  # Rate limiting
  rate_limits:
    chatgpt: "100 requests/hour"
    claude: "50 requests/hour"
    amazonq: "200 requests/hour"

  # Smart caching
  cache_duration: "15m"
  cache_similar_queries: true

  # Batch processing
  batch_size: 10
  batch_interval: "5m"
```

#### 3. Reliability và Fallbacks

```yaml
ai_reliability:
  # Fallback chain
  primary_ai: "chatgpt"
  fallback_ai: "claude"
  final_fallback: "rule_based_system"

  # Timeout handling
  timeout: "30s"
  retry_attempts: 3

  # Quality checks
  confidence_threshold: 0.8
  human_review_required: true # for critical actions
```

---

## 🔍 Troubleshooting

### 🚨 Common Issues

#### Service không start

```bash
# Check service status
sudo systemctl status ducla-agent

# Check logs
sudo journalctl -u ducla-agent --no-pager

# Check config
ducla-agent config validate

# Check permissions
ls -la /etc/ducla/agent.yaml
ls -la /opt/ducla/
```

#### API endpoints không response

```bash
# Check if ports are listening
sudo netstat -tlnp | grep ducla-agent

# Check firewall
sudo ufw status
sudo iptables -L

# Test locally
curl -v http://127.0.0.1:8080/api/v1/status
```

#### Output connection issues

```bash
# Test specific output
ducla-agent output test s3
ducla-agent output test prometheus

# Check network connectivity
ping prometheus-server
telnet elasticsearch 9200

# Check credentials
aws sts get-caller-identity
```

### 📞 Support

- **Documentation**: https://github.com/duclacloud/DUCLA-CLOUD-AGENT
- **Issues**: https://github.com/duclacloud/DUCLA-CLOUD-AGENT/issues
- **Email**: mandỵhades@hotmail.com.vn

---

**🎉 Chúc bạn sử dụng Ducla Cloud Agent thành công!**
