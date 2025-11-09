# 📊 Hướng Dẫn Test Ducla Agent với Prometheus & Grafana

## 🎯 Mục Tiêu Test

Test này sẽ verify:
- ✅ Ducla Agent metrics collection
- ✅ Prometheus scraping và storage
- ✅ Grafana visualization
- ✅ AlertManager notifications
- ✅ API endpoints integration
- ✅ Real-time monitoring capabilities

## 📋 Prerequisites

### System Requirements
- **OS**: Linux (Ubuntu 18.04+, CentOS 7+)
- **Docker**: Version 20.10+
- **Docker Compose**: Version 1.29+
- **RAM**: Minimum 2GB, Recommended 4GB+
- **Disk**: 5GB free space
- **Network**: Internet access

### Pre-test Checklist
```bash
# Check Docker
docker --version
docker info

# Check Docker Compose
docker-compose --version

# Check available resources
free -h
df -h
```

## 🚀 Test Setup

### 1. Prepare Environment
```bash
# Clone repository (if not already done)
git clone https://github.com/duclacloud/DUCLA-CLOUD-AGENT.git
cd DUCLA-CLOUD-AGENT

# Ensure scripts are executable
chmod +x start-monitoring.sh stop-monitoring.sh

# Create necessary directories
mkdir -p data logs
sudo chown -R $USER:$USER data logs
```

### 2. Start Monitoring Stack
```bash
# Start all services
./start-monitoring.sh

# Expected output:
# 🚀 Starting Ducla Agent Monitoring Stack
# ========================================
# 📁 Creating directories...
# 🔨 Building Ducla Agent Docker image...
# 🐳 Starting containers...
# ⏳ Waiting for services to start...
# 🔍 Checking service health...
# ✅ Ducla Agent: http://localhost:8080 (API), http://localhost:8081 (Health)
# ✅ Prometheus: http://localhost:9091
# ✅ Grafana: http://localhost:3000 (admin/admin123)
```

### 3. Verify Services
```bash
# Check container status
docker-compose -f docker-compose.monitoring.yml ps

# Expected output: All services should be "Up"
#     Name                   Command               State                    Ports
# -------------------------------------------------------------------------------------
# alertmanager    /bin/alertmanager --config ...   Up      0.0.0.0:9093->9093/tcp
# cadvisor        /usr/bin/cadvisor -logtostderr   Up      0.0.0.0:8082->8080/tcp
# ducla-agent     /usr/bin/ducla-agent -config ... Up      0.0.0.0:8080->8080/tcp, ...
# grafana         /run.sh                          Up      0.0.0.0:3000->3000/tcp
# node-exporter   /bin/node_exporter --path.pr ... Up      0.0.0.0:9100->9100/tcp
# prometheus      /bin/prometheus --config.fil ... Up      0.0.0.0:9091->9090/tcp
```

## 🧪 Test Cases

### Test Case 1: Ducla Agent Health Check

#### Objective
Verify Ducla Agent is running và responding to health checks.

#### Steps
```bash
# 1. Test health endpoint
curl -i http://localhost:8081/health

# Expected: HTTP/1.1 200 OK
# Content-Type: application/json

# 2. Test API status
curl -s http://localhost:8080/api/v1/status | jq .

# Expected JSON response:
# {
#   "success": true,
#   "data": {
#     "running": true,
#     "tasks": {
#       "worker_count": 5,
#       "queue_size": 0,
#       "running_tasks": 0,
#       "completed_tasks": 0,
#       "total_tasks": 0
#     },
#     "metrics": {}
#   }
# }

# 3. Test metrics endpoint
curl -s http://localhost:9090/metrics | head -20

# Expected: Prometheus metrics format
# # HELP ducla_agent_info Agent information
# # TYPE ducla_agent_info gauge
# ducla_agent_info{version="1.0.0",build_time="2025-11-09"} 1
```

#### Success Criteria
- ✅ Health endpoint returns HTTP 200
- ✅ API status shows `"running": true`
- ✅ Metrics endpoint returns Prometheus format
- ✅ No error messages in logs

### Test Case 2: Prometheus Scraping

#### Objective
Verify Prometheus is successfully scraping Ducla Agent metrics.

#### Steps
```bash
# 1. Access Prometheus UI
open http://localhost:9091
# Or: firefox http://localhost:9091

# 2. Check targets status
# Navigate to: Status > Targets
# Look for: ducla-agent job

# 3. Query metrics via API
curl -s "http://localhost:9091/api/v1/query?query=up{job=\"ducla-agent\"}" | jq .

# Expected response:
# {
#   "status": "success",
#   "data": {
#     "resultType": "vector",
#     "result": [
#       {
#         "metric": {
#           "__name__": "up",
#           "instance": "ducla-agent:9090",
#           "job": "ducla-agent"
#         },
#         "value": [1699534800, "1"]
#       }
#     ]
#   }
# }

# 4. Test specific Ducla metrics
curl -s "http://localhost:9091/api/v1/query?query=ducla_tasks_total" | jq .
curl -s "http://localhost:9091/api/v1/query?query=ducla_http_requests_total" | jq .
```

#### Success Criteria
- ✅ Ducla Agent target shows as "UP" in Prometheus
- ✅ `up{job="ducla-agent"}` returns value "1"
- ✅ Ducla-specific metrics are available
- ✅ Scrape duration < 1 second

### Test Case 3: Grafana Dashboard

#### Objective
Verify Grafana can visualize Ducla Agent metrics.

#### Steps
```bash
# 1. Access Grafana
open http://localhost:3000
# Login: admin / admin123

# 2. Check datasource
# Navigate to: Configuration > Data Sources
# Verify: Prometheus datasource is configured and working

# 3. Import dashboard
# Navigate to: Dashboards
# Look for: "Ducla Cloud Agent Dashboard"

# 4. Test dashboard panels
# Verify panels show data:
# - Agent Status (should show "1" - UP)
# - API Requests Rate
# - Task Execution
# - Memory Usage
# - CPU Usage

# 5. Test queries manually
# Navigate to: Explore
# Try queries:
# - up{job="ducla-agent"}
# - rate(ducla_http_requests_total[5m])
# - process_resident_memory_bytes{job="ducla-agent"}
```

#### Success Criteria
- ✅ Grafana login successful
- ✅ Prometheus datasource connected
- ✅ Dashboard loads without errors
- ✅ All panels show data
- ✅ Real-time updates working

### Test Case 4: Load Testing

#### Objective
Generate load on Ducla Agent và verify metrics collection.

#### Steps
```bash
# 1. Generate API requests
for i in {1..100}; do
  curl -s http://localhost:8080/api/v1/status > /dev/null
  echo "Request $i completed"
  sleep 0.1
done

# 2. Create some tasks (if task API exists)
for i in {1..10}; do
  curl -X POST http://localhost:8080/api/v1/tasks \
    -H "Content-Type: application/json" \
    -d '{"task": "test-task-'$i'", "action": "echo", "params": {"message": "Hello World"}}'
  sleep 1
done

# 3. Monitor metrics during load
watch -n 2 'curl -s http://localhost:9090/metrics | grep -E "(ducla_http_requests_total|ducla_tasks_)"'

# 4. Check Grafana dashboard
# Refresh dashboard and observe:
# - API request rate increase
# - Task execution metrics
# - Memory/CPU usage changes
```

#### Success Criteria
- ✅ API request metrics increase during load
- ✅ Task metrics update correctly
- ✅ Resource usage metrics reflect load
- ✅ No errors or timeouts during load

### Test Case 5: Alert Testing

#### Objective
Test AlertManager integration và alert rules.

#### Steps
```bash
# 1. Check AlertManager
open http://localhost:9093

# 2. Trigger alert by stopping Ducla Agent
docker-compose -f docker-compose.monitoring.yml stop ducla-agent

# 3. Wait for alert (should trigger after 1 minute)
# Check AlertManager UI for "DuclaAgentDown" alert

# 4. Check Prometheus alerts
open http://localhost:9091/alerts

# 5. Restart agent and verify alert resolves
docker-compose -f docker-compose.monitoring.yml start ducla-agent

# 6. Test webhook (if configured)
curl -X POST http://localhost:8080/api/v1/alerts \
  -H "Content-Type: application/json" \
  -d '{
    "alerts": [
      {
        "status": "firing",
        "labels": {
          "alertname": "TestAlert",
          "severity": "warning"
        },
        "annotations": {
          "summary": "Test alert from monitoring test"
        }
      }
    ]
  }'
```

#### Success Criteria
- ✅ Alert fires when agent is down
- ✅ Alert appears in AlertManager
- ✅ Alert resolves when agent restarts
- ✅ Webhook endpoint receives alerts

### Test Case 6: Data Persistence

#### Objective
Verify data persistence across restarts.

#### Steps
```bash
# 1. Generate some metrics data
for i in {1..50}; do
  curl -s http://localhost:8080/api/v1/status > /dev/null
  sleep 2
done

# 2. Stop monitoring stack
./stop-monitoring.sh

# 3. Restart monitoring stack
./start-monitoring.sh

# 4. Check if historical data is preserved
# In Grafana, set time range to "Last 1 hour"
# Verify data from before restart is still visible

# 5. Query historical data via Prometheus API
curl -s "http://localhost:9091/api/v1/query_range?query=up{job=\"ducla-agent\"}&start=$(date -d '1 hour ago' +%s)&end=$(date +%s)&step=60" | jq .
```

#### Success Criteria
- ✅ Historical data preserved after restart
- ✅ Grafana shows data from before restart
- ✅ Prometheus query_range returns historical data
- ✅ No data gaps during restart

## 📊 Performance Benchmarks

### Expected Metrics
```bash
# Resource Usage (normal operation)
docker stats --no-stream

# Expected ranges:
# ducla-agent:    CPU: 0.1-2%,   Memory: 20-100MB
# prometheus:     CPU: 0.5-5%,   Memory: 100-500MB
# grafana:        CPU: 0.1-1%,   Memory: 50-200MB
# node-exporter:  CPU: 0.1-0.5%, Memory: 10-50MB
# cadvisor:       CPU: 0.5-2%,   Memory: 50-150MB
# alertmanager:   CPU: 0.1-0.5%, Memory: 20-100MB
```

### Performance Tests
```bash
# 1. Metrics collection latency
time curl -s http://localhost:9090/metrics > /dev/null
# Expected: < 100ms

# 2. API response time
time curl -s http://localhost:8080/api/v1/status > /dev/null
# Expected: < 50ms

# 3. Prometheus query performance
time curl -s "http://localhost:9091/api/v1/query?query=up" > /dev/null
# Expected: < 200ms

# 4. Grafana dashboard load time
# Manual test: Dashboard should load in < 3 seconds
```

## 🐛 Troubleshooting

### Common Issues

#### Issue 1: Containers won't start
```bash
# Check Docker resources
docker system df
docker system prune -f

# Check port conflicts
netstat -tlnp | grep -E "(3000|8080|8081|9090|9091|9093)"

# Check logs
docker-compose -f docker-compose.monitoring.yml logs ducla-agent
```

#### Issue 2: Prometheus can't scrape Ducla Agent
```bash
# Check network connectivity
docker exec prometheus ping ducla-agent

# Check Ducla Agent metrics endpoint
docker exec ducla-agent curl localhost:9090/metrics

# Check Prometheus config
docker exec prometheus cat /etc/prometheus/prometheus.yml
```

#### Issue 3: Grafana shows no data
```bash
# Test Prometheus datasource
curl -s http://localhost:9091/api/v1/query?query=up

# Check Grafana logs
docker-compose -f docker-compose.monitoring.yml logs grafana

# Verify datasource configuration
# In Grafana: Configuration > Data Sources > Test
```

#### Issue 4: High resource usage
```bash
# Check container resources
docker stats

# Reduce Prometheus retention
# Edit monitoring/prometheus.yml:
# --storage.tsdb.retention.time=1d

# Restart with resource limits
docker-compose -f docker-compose.monitoring.yml up -d --scale cadvisor=0
```

## 📝 Test Report Template

### Test Execution Report

**Date**: ___________  
**Tester**: ___________  
**Environment**: ___________  

#### Test Results Summary
| Test Case | Status | Duration | Notes |
|-----------|--------|----------|-------|
| Health Check | ✅/❌ | ___s | |
| Prometheus Scraping | ✅/❌ | ___s | |
| Grafana Dashboard | ✅/❌ | ___s | |
| Load Testing | ✅/❌ | ___s | |
| Alert Testing | ✅/❌ | ___s | |
| Data Persistence | ✅/❌ | ___s | |

#### Performance Metrics
- **Ducla Agent Memory**: ___MB
- **Total Stack Memory**: ___MB
- **API Response Time**: ___ms
- **Metrics Collection Time**: ___ms
- **Dashboard Load Time**: ___s

#### Issues Found
1. ___________
2. ___________
3. ___________

#### Recommendations
1. ___________
2. ___________
3. ___________

## 🧹 Cleanup

### Stop and Clean
```bash
# Stop monitoring stack
./stop-monitoring.sh

# Remove all data (optional)
docker-compose -f docker-compose.monitoring.yml down -v

# Clean Docker resources
docker system prune -f
docker volume prune -f

# Remove created directories
rm -rf data logs
```

## 🎯 Success Criteria Summary

### ✅ Test Passes If:
- All containers start successfully
- Ducla Agent responds to health checks
- Prometheus scrapes metrics without errors
- Grafana displays data correctly
- Alerts fire and resolve properly
- Performance meets benchmarks
- Data persists across restarts

### ❌ Test Fails If:
- Any container fails to start
- Health checks return errors
- Prometheus shows scrape errors
- Grafana shows no data
- Alerts don't fire when expected
- Performance below benchmarks
- Data loss during restarts

## 📚 Additional Resources

- **Prometheus Documentation**: https://prometheus.io/docs/
- **Grafana Documentation**: https://grafana.com/docs/
- **Docker Compose Reference**: https://docs.docker.com/compose/
- **Ducla Agent API Reference**: [API-REFERENCE.md](API-REFERENCE.md)

---

**🎉 Happy Testing! This comprehensive test ensures Ducla Agent monitoring integration works perfectly.**