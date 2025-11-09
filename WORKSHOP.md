# 🚀 Ducla Cloud Agent v1.0.0 Workshop

## Giới Thiệu Workshop

Chào mừng bạn đến với workshop **Ducla Cloud Agent v1.0.0**! Workshop này sẽ hướng dẫn bạn từng bước để khám phá và trải nghiệm tất cả các tính năng chính của agent.

## 📋 Mục Tiêu Workshop

Sau khi hoàn thành workshop này, bạn sẽ:

- ✅ Hiểu được kiến trúc và tính năng của Ducla Cloud Agent
- ✅ Biết cách build và deploy agent
- ✅ Thành thạo việc sử dụng các API endpoints
- ✅ Có thể giám sát và troubleshoot agent
- ✅ Sẵn sàng triển khai trong môi trường thực tế

## 🎯 Đối Tượng

- **Developers**: Muốn tích hợp agent vào hệ thống
- **DevOps Engineers**: Cần deploy và quản lý agent
- **System Administrators**: Giám sát và vận hành agent
- **Technical Leaders**: Đánh giá và quyết định sử dụng

## ⏱️ Thời Gian

- **Tổng thời gian**: 60-90 phút
- **Chuẩn bị**: 10 phút
- **Demo tự động**: 30 phút
- **Thực hành**: 30-40 phút
- **Q&A**: 10 phút

## 🛠️ Yêu Cầu Chuẩn Bị

### Hệ Thống
- **OS**: Linux (Ubuntu 20.04+ khuyến nghị)
- **RAM**: Tối thiểu 2GB
- **Disk**: 500MB trống
- **Network**: Kết nối internet

### Phần Mềm
- **Go**: 1.21+ (sẽ hướng dẫn cài đặt)
- **curl**: Để test API
- **jq**: Để format JSON (tùy chọn)
- **git**: Để clone repository

### Kiến Thức
- Cơ bản về Linux command line
- Hiểu biết về REST API
- Kinh nghiệm với YAML configuration

## 📚 Cấu Trúc Workshop

### Phase 1: Chuẩn Bị (10 phút)
1. **Environment Setup**
   - Cài đặt Go
   - Clone repository
   - Kiểm tra dependencies

2. **Build Agent**
   - Build từ source code
   - Tạo version 1.0.0
   - Verify binary

### Phase 2: Demo Tự Động (30 phút)
3. **Automated Demo**
   - Chạy script demo tự động
   - Khám phá tất cả tính năng
   - Hiểu workflow cơ bản

### Phase 3: Thực Hành (30-40 phút)
4. **Manual Testing**
   - Cấu hình agent
   - Test từng API endpoint
   - Monitoring và troubleshooting

5. **Advanced Features**
   - Security configuration
   - Performance tuning
   - Production deployment

### Phase 4: Wrap-up (10 phút)
6. **Q&A và Next Steps**
   - Thảo luận use cases
   - Best practices
   - Roadmap

---

## 🚀 Bắt Đầu Workshop

### Step 1: Environment Setup

#### 1.1 Cài Đặt Go (nếu chưa có)

```bash
# Kiểm tra Go version hiện tại
go version

# Nếu chưa có hoặc version < 1.21, cài đặt mới
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz

# Thêm vào PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify
go version
```

#### 1.2 Clone Repository

```bash
# Clone repository (hoặc sử dụng source code có sẵn)
git clone <repository-url>
cd ducla-cloud-agent

# Hoặc nếu đã có source
cd ducla-cloud-agent
```

#### 1.3 Install Dependencies

```bash
# Cài đặt jq để format JSON
sudo apt update
sudo apt install -y jq curl

# Download Go dependencies
go mod tidy
```

### Step 2: Build Agent

#### 2.1 Build Version 1.0.0

```bash
# Build agent với version info
./build-v1.sh

# Verify build
ls -la dist/
./dist/ducla-agent -version
```

**Expected Output:**
```
Ducla Cloud Agent
  Version:    1.0.0
  Build Time: 2025-11-09_XX:XX:XX_UTC
  Git Commit: xxxxxxx
  Go Version: go1.21.5
  OS/Arch:    linux/amd64
```

### Step 3: Automated Demo

#### 3.1 Chạy Demo Script

```bash
# Chạy demo tự động
./demo.sh
```

**Demo sẽ showcase:**
- ✅ Version information
- ✅ Health monitoring
- ✅ Metrics collection
- ✅ Agent status
- ✅ Task execution
- ✅ File operations
- ✅ API endpoints overview

#### 3.2 Quan Sát và Ghi Chú

Trong quá trình demo, hãy chú ý:
- **Startup time**: Agent khởi động nhanh như thế nào?
- **API responses**: Format và structure của responses
- **Health checks**: Các metrics được monitor
- **Task execution**: Workflow của task processing
- **File operations**: Các thao tác file được hỗ trợ

### Step 4: Manual Testing

#### 4.1 Cấu Hình Agent

Tạo file `workshop-config.yaml`:

```yaml
agent:
  id: "workshop-agent"
  name: "Workshop Demo Agent"
  environment: "workshop"
  region: "local"
  capabilities:
    - "file_operations"
    - "task_execution"
    - "system_monitoring"

master:
  url: ""
  token: ""
  max_reconnect_attempts: 0

api:
  http:
    enabled: true
    address: "0.0.0.0"
    port: 8080
    tls:
      enabled: false

security:
  auth:
    enabled: false

storage:
  data_dir: "./workshop-data"
  temp_dir: "./workshop-tmp"

logging:
  level: "debug"
  format: "json"

metrics:
  enabled: true
  address: "0.0.0.0"
  port: 9090

health:
  enabled: true
  address: "0.0.0.0"
  port: 8081

executor:
  workers: 5
  queue_size: 100
  task_timeout: 120s
```

#### 4.2 Start Agent

```bash
# Start agent
./dist/ducla-agent -config workshop-config.yaml -debug
```

#### 4.3 Test API Endpoints

**Terminal mới để test APIs:**

```bash
# 1. Agent Status
curl -s http://localhost:8080/api/v1/status | jq

# 2. Health Check
curl -s http://localhost:8081/health | jq

# 3. Metrics
curl -s http://localhost:9090/metrics | head -20

# 4. Create Task
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "type": "shell",
    "name": "workshop-task",
    "command": "ls",
    "args": ["-la", "/tmp"]
  }' | jq

# 5. List Tasks
curl -s http://localhost:8080/api/v1/tasks | jq

# 6. File Operations
mkdir -p workshop-files
echo "Workshop demo file" > workshop-files/demo.txt

curl -X POST http://localhost:8080/api/v1/files \
  -H "Content-Type: application/json" \
  -d "{
    \"type\": \"list\",
    \"source_path\": \"$(pwd)/workshop-files\"
  }" | jq
```

### Step 5: Advanced Features

#### 5.1 Performance Testing

```bash
# Tạo nhiều tasks đồng thời
for i in {1..10}; do
  curl -X POST http://localhost:8080/api/v1/tasks \
    -H "Content-Type: application/json" \
    -d "{
      \"type\": \"shell\",
      \"name\": \"perf-test-$i\",
      \"command\": \"sleep\",
      \"args\": [\"2\"]
    }" &
done

wait

# Kiểm tra task status
curl -s http://localhost:8080/api/v1/tasks | jq '.data.tasks | length'
```

#### 5.2 Monitoring

```bash
# Monitor health
watch -n 2 'curl -s http://localhost:8081/health | jq .data.summary'

# Monitor metrics
curl -s http://localhost:9090/metrics | grep ducla_
```

#### 5.3 Error Handling

```bash
# Test invalid task
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "type": "shell",
    "name": "error-task",
    "command": "invalid-command"
  }' | jq

# Test invalid file operation
curl -X POST http://localhost:8080/api/v1/files \
  -H "Content-Type: application/json" \
  -d '{
    "type": "copy",
    "source_path": "/nonexistent/file.txt",
    "dest_path": "/tmp/copy.txt"
  }' | jq
```

## 📊 Workshop Checklist

### ✅ Chuẩn Bị
- [ ] Go 1.21+ installed
- [ ] Repository cloned
- [ ] Dependencies installed
- [ ] Agent built successfully

### ✅ Demo
- [ ] Automated demo completed
- [ ] All features demonstrated
- [ ] API endpoints tested
- [ ] Performance observed

### ✅ Thực Hành
- [ ] Manual configuration
- [ ] API testing completed
- [ ] Error handling tested
- [ ] Performance testing done

### ✅ Advanced
- [ ] Security features explored
- [ ] Monitoring setup
- [ ] Production considerations discussed

## 🎯 Key Takeaways

### Tính Năng Chính
1. **High Performance**: Go-based, low resource usage
2. **Complete API**: REST + gRPC support
3. **Monitoring**: Built-in health checks and metrics
4. **Flexibility**: Configurable and extensible
5. **Production Ready**: Security, logging, deployment options

### Use Cases
- **CI/CD Pipelines**: Task execution in build/deploy
- **System Automation**: File operations and system tasks
- **Monitoring**: Health checks and metrics collection
- **Microservices**: Service-to-service communication
- **Edge Computing**: Lightweight agent deployment

### Best Practices
- **Configuration**: Use environment variables for secrets
- **Security**: Enable TLS and authentication in production
- **Monitoring**: Set up proper logging and metrics
- **Deployment**: Use containers or systemd services
- **Scaling**: Configure worker pools based on workload

## 🚀 Next Steps

### Immediate Actions
1. **Evaluate**: Assess fit for your use case
2. **Prototype**: Build a small proof of concept
3. **Test**: Performance and load testing
4. **Security**: Review security requirements
5. **Deploy**: Plan production deployment

### Long-term Planning
1. **Integration**: Plan integration with existing systems
2. **Monitoring**: Set up comprehensive monitoring
3. **Scaling**: Plan for horizontal scaling
4. **Maintenance**: Establish update and maintenance procedures
5. **Training**: Train team on operation and troubleshooting

## 📞 Support & Resources

### Documentation
- **README-VI.md**: Vietnamese documentation
- **CHANGELOG.md**: Version history and features
- **API Documentation**: Complete API reference

### Community
- **GitHub Issues**: Bug reports and feature requests
- **Discussions**: Community discussions and Q&A
- **Discord/Slack**: Real-time community support

### Professional Support
- **Consulting**: Architecture and implementation guidance
- **Training**: Team training and workshops
- **Support**: Production support and maintenance

---

## 🎉 Kết Thúc Workshop

Cảm ơn bạn đã tham gia workshop **Ducla Cloud Agent v1.0.0**!

### Feedback
Vui lòng chia sẻ feedback về workshop:
- Nội dung có hữu ích không?
- Thời gian có phù hợp không?
- Có phần nào cần cải thiện?
- Bạn có sẵn sàng sử dụng trong dự án thực tế?

### Stay Connected
- ⭐ Star repository nếu bạn thấy hữu ích
- 🐛 Report bugs nếu phát hiện
- 💡 Suggest features cho version tiếp theo
- 📢 Share với team và community

**Happy Coding! 🚀**