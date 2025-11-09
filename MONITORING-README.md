# 📊 Ducla Agent Monitoring Test Suite

## 🎯 Overview

This test suite verifies Ducla Agent integration with Prometheus, Grafana, and AlertManager using Docker Compose. It provides comprehensive testing for monitoring, metrics collection, visualization, and alerting capabilities.

## 🚀 Quick Start

### 1. Start Monitoring Stack
```bash
./start-monitoring.sh
```

### 2. Run Quick Test
```bash
./quick-test.sh
```

### 3. Run Full Test Suite
```bash
./run-monitoring-tests.sh
```

### 4. Stop Monitoring Stack
```bash
./stop-monitoring.sh
```

## 📁 Files Overview

### Core Files
- **`docker-compose.monitoring.yml`** - Complete monitoring stack definition
- **`start-monitoring.sh`** - One-command startup script
- **`stop-monitoring.sh`** - Clean shutdown script

### Test Scripts
- **`quick-test.sh`** - Fast endpoint verification (2 minutes)
- **`run-monitoring-tests.sh`** - Comprehensive test suite (10 minutes)
- **`MONITORING-TEST-GUIDE.md`** - Detailed manual testing guide

### Configuration
- **`monitoring/prometheus.yml`** - Prometheus scraping configuration
- **`monitoring/grafana/`** - Grafana dashboards and datasources
- **`monitoring/rules/ducla-agent.yml`** - Alert rules
- **`monitoring/alertmanager.yml`** - AlertManager configuration

## 🧪 Test Types

### 1. Quick Test (`./quick-test.sh`)
**Duration**: ~2 minutes  
**Purpose**: Verify all services are running and accessible

```bash
🔍 Testing Ducla Agent Health...
✅ Ducla Agent Health is working (HTTP 200)
🔍 Testing Ducla Agent API...
✅ Ducla Agent API is working (HTTP 200)
🔍 Testing Prometheus...
✅ Prometheus is working (HTTP 200)
🔍 Testing Grafana...
✅ Grafana is working (HTTP 200)
```

### 2. Full Test Suite (`./run-monitoring-tests.sh`)
**Duration**: ~10 minutes  
**Purpose**: Comprehensive integration testing

**Test Cases**:
- ✅ Pre-test checks (Docker, disk space)
- ✅ Service startup and health
- ✅ Prometheus scraping integration
- ✅ Load testing and metrics collection
- ✅ Alert firing and resolution
- ✅ Performance benchmarks
- ✅ Container stability

**Sample Output**:
```bash
📊 TEST RESULTS SUMMARY
Total Tests: 18
Passed: 18
Failed: 0
Success Rate: 100%
```

### 3. Manual Testing (`MONITORING-TEST-GUIDE.md`)
**Duration**: ~30 minutes  
**Purpose**: Detailed step-by-step verification

## 🎛️ Services and Ports

| Service | Port | URL | Purpose |
|---------|------|-----|---------|
| Ducla Agent API | 8080 | http://localhost:8080 | Main API |
| Ducla Agent Health | 8081 | http://localhost:8081/health | Health checks |
| Ducla Agent Metrics | 9090 | http://localhost:9090/metrics | Prometheus metrics |
| Prometheus | 9091 | http://localhost:9091 | Metrics storage |
| Grafana | 3000 | http://localhost:3000 | Visualization |
| AlertManager | 9093 | http://localhost:9093 | Alert management |
| Node Exporter | 9100 | http://localhost:9100 | System metrics |
| cAdvisor | 8082 | http://localhost:8082 | Container metrics |

## 📊 Expected Metrics

### Ducla Agent Metrics
```prometheus
# Agent information
ducla_agent_info{version="1.0.0"} 1

# HTTP requests
ducla_http_requests_total{method="GET",endpoint="/api/v1/status"} 42

# Task execution
ducla_tasks_total 10
ducla_tasks_running 2
ducla_tasks_completed 8
ducla_tasks_queue_size 0

# Resource usage
process_resident_memory_bytes 52428800
process_cpu_seconds_total 1.23
```

### System Metrics (Node Exporter)
```prometheus
# CPU usage
node_cpu_seconds_total{mode="idle"} 12345

# Memory usage
node_memory_MemAvailable_bytes 8589934592

# Disk usage
node_filesystem_avail_bytes{mountpoint="/"} 42949672960
```

## 🚨 Alert Rules

### DuclaAgentDown
- **Condition**: `up{job="ducla-agent"} == 0`
- **Duration**: 1 minute
- **Severity**: Critical

### DuclaAgentHighMemory
- **Condition**: `process_resident_memory_bytes{job="ducla-agent"} > 100MB`
- **Duration**: 5 minutes
- **Severity**: Warning

### DuclaAgentHighTaskQueue
- **Condition**: `ducla_tasks_queue_size > 50`
- **Duration**: 2 minutes
- **Severity**: Warning

### DuclaAgentAPIErrors
- **Condition**: `rate(ducla_http_requests_total{status=~"5.."}[5m]) > 0.1`
- **Duration**: 2 minutes
- **Severity**: Warning

## 📈 Performance Benchmarks

### Resource Usage (Normal Operation)
| Service | CPU | Memory | Notes |
|---------|-----|--------|-------|
| Ducla Agent | 0.1-2% | 20-100MB | Scales with load |
| Prometheus | 0.5-5% | 100-500MB | Depends on retention |
| Grafana | 0.1-1% | 50-200MB | Stable |
| Node Exporter | 0.1-0.5% | 10-50MB | Very light |
| cAdvisor | 0.5-2% | 50-150MB | Container monitoring |
| AlertManager | 0.1-0.5% | 20-100MB | Minimal |

### Response Times
- **API Endpoint**: < 50ms
- **Metrics Endpoint**: < 100ms
- **Prometheus Query**: < 200ms
- **Grafana Dashboard**: < 3s load time

## 🐛 Troubleshooting

### Common Issues

#### Services Won't Start
```bash
# Check Docker resources
docker system df
docker system prune -f

# Check port conflicts
netstat -tlnp | grep -E "(3000|8080|9090|9091)"

# View logs
docker-compose -f docker-compose.monitoring.yml logs
```

#### No Metrics in Grafana
```bash
# Test Prometheus datasource
curl http://localhost:9091/api/v1/query?query=up

# Check Grafana datasource config
# Go to: Configuration > Data Sources > Test
```

#### High Resource Usage
```bash
# Monitor resources
docker stats

# Reduce Prometheus retention
# Edit monitoring/prometheus.yml:
# --storage.tsdb.retention.time=1d
```

#### Alerts Not Firing
```bash
# Check alert rules
curl http://localhost:9091/api/v1/rules

# Check AlertManager config
docker exec alertmanager cat /etc/alertmanager/alertmanager.yml
```

## 🔧 Customization

### Add Custom Metrics
Edit `monitoring/prometheus.yml` to add new scrape targets:
```yaml
scrape_configs:
  - job_name: 'my-app'
    static_configs:
      - targets: ['my-app:8080']
```

### Add Custom Dashboards
Place JSON files in `monitoring/grafana/dashboards/`:
```bash
cp my-dashboard.json monitoring/grafana/dashboards/
docker-compose -f docker-compose.monitoring.yml restart grafana
```

### Modify Alert Rules
Edit `monitoring/rules/ducla-agent.yml`:
```yaml
- alert: MyCustomAlert
  expr: my_metric > 100
  for: 5m
  labels:
    severity: warning
  annotations:
    summary: "Custom alert fired"
```

## 📚 Additional Resources

- **Prometheus Documentation**: https://prometheus.io/docs/
- **Grafana Documentation**: https://grafana.com/docs/
- **AlertManager Documentation**: https://prometheus.io/docs/alerting/latest/alertmanager/
- **Docker Compose Reference**: https://docs.docker.com/compose/

## 🎯 Success Criteria

### ✅ Test Suite Passes When:
- All containers start successfully
- All health checks return HTTP 200
- Prometheus scrapes Ducla Agent metrics
- Grafana displays data correctly
- Alerts fire and resolve properly
- Performance meets benchmarks
- No container restarts or errors

### 🎉 What This Proves:
- Ducla Agent monitoring integration works
- Metrics collection is reliable
- Visualization is functional
- Alerting system is operational
- Performance is acceptable
- System is production-ready

---

**🚀 Ready to test Ducla Agent monitoring? Start with `./start-monitoring.sh` and then run `./quick-test.sh`!**