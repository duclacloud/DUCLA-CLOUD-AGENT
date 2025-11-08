# Network Configuration Guide

## IP Address Ranges Explained

### Private IP Ranges (RFC 1918)
These are **private network addresses** used in internal networks, not "fake" IPs:

| Range | CIDR | Description | Total IPs |
|-------|------|-------------|-----------|
| `10.0.0.0/8` | Class A Private | 10.0.0.0 - 10.255.255.255 | 16,777,216 |
| `172.16.0.0/12` | Class B Private | 172.16.0.0 - 172.31.255.255 | 1,048,576 |
| `192.168.0.0/16` | Class C Private | 192.168.0.0 - 192.168.255.255 | 65,536 |

### Special Use IP Ranges

| Range | CIDR | Description |
|-------|------|-------------|
| `127.0.0.0/8` | Loopback | Local machine (localhost) |
| `169.254.0.0/16` | Link-Local | Auto-assigned when DHCP fails |
| `224.0.0.0/4` | Multicast | Multicast addresses |
| `240.0.0.0/4` | Reserved | Reserved for future use |

## Firewall Configuration

### Default Configuration (Development)
```yaml
security:
  firewall:
    enabled: false  # Disabled by default for development
    allowed_ips:
      - "10.0.0.0/8"      # Private networks
      - "172.16.0.0/12"   # Private networks  
      - "192.168.0.0/16"  # Private networks
      - "127.0.0.0/8"     # Loopback
```

### Production Configuration
```yaml
security:
  firewall:
    enabled: true
    allowed_ips:
      # Only allow specific networks/IPs
      - "10.10.0.0/16"        # Your internal network
      - "203.0.113.100/32"    # Master server public IP
      - "198.51.100.0/24"     # Management network
    blocked_ips:
      - "0.0.0.0/8"           # Invalid source
      - "224.0.0.0/4"         # Multicast
      - "240.0.0.0/4"         # Reserved
    rate_limiting: true
    requests_per_min: 100     # Stricter rate limiting
```

### Cloud Environment Examples

#### AWS VPC
```yaml
security:
  firewall:
    enabled: true
    allowed_ips:
      - "10.0.0.0/16"         # Your VPC CIDR
      - "10.1.0.0/16"         # Peered VPC
      # AWS service endpoints
      - "52.95.0.0/16"        # S3 endpoints (example)
```

#### Azure Virtual Network
```yaml
security:
  firewall:
    enabled: true
    allowed_ips:
      - "10.0.0.0/16"         # Your VNet CIDR
      - "10.1.0.0/16"         # Peered VNet
      # Azure service IPs
      - "20.0.0.0/8"          # Azure public IPs (example)
```

#### Google Cloud VPC
```yaml
security:
  firewall:
    enabled: true
    allowed_ips:
      - "10.128.0.0/16"       # Your VPC subnet
      - "10.129.0.0/16"       # Another subnet
      # GCP service ranges
      - "35.199.192.0/19"     # Google services (example)
```

## Network Security Best Practices

### 1. Principle of Least Privilege
```yaml
# ❌ Too permissive
allowed_ips:
  - "0.0.0.0/0"  # Allows everything

# ✅ Specific networks only
allowed_ips:
  - "10.0.1.0/24"    # Specific subnet
  - "192.168.1.100/32"  # Specific IP
```

### 2. Environment-Specific Configuration
```yaml
# Development
security:
  firewall:
    enabled: false  # More permissive for development

# Staging  
security:
  firewall:
    enabled: true
    allowed_ips:
      - "10.0.0.0/16"  # Internal network

# Production
security:
  firewall:
    enabled: true
    allowed_ips:
      - "10.0.1.0/24"  # Very specific range
```

### 3. Rate Limiting
```yaml
security:
  firewall:
    rate_limiting: true
    requests_per_min: 1000  # Adjust based on expected load
```

## Common Network Scenarios

### 1. Corporate Network
```yaml
security:
  firewall:
    enabled: true
    allowed_ips:
      - "192.168.0.0/16"      # Corporate LAN
      - "10.0.0.0/8"          # Corporate WAN
      - "203.0.113.0/24"      # Office public IP range
```

### 2. Data Center
```yaml
security:
  firewall:
    enabled: true
    allowed_ips:
      - "10.10.0.0/16"        # Management network
      - "10.20.0.0/16"        # Application network
      - "10.30.0.0/16"        # Storage network
```

### 3. Kubernetes Cluster
```yaml
security:
  firewall:
    enabled: true
    allowed_ips:
      - "10.244.0.0/16"       # Pod network (Flannel)
      - "10.96.0.0/12"        # Service network
      - "172.17.0.0/16"       # Docker bridge
```

### 4. Multi-Cloud Setup
```yaml
security:
  firewall:
    enabled: true
    allowed_ips:
      # AWS VPC
      - "10.0.0.0/16"
      # Azure VNet  
      - "10.1.0.0/16"
      # GCP VPC
      - "10.2.0.0/16"
      # VPN connections
      - "192.168.100.0/24"
```

## Troubleshooting Network Issues

### 1. Connection Refused
```bash
# Check if IP is allowed
curl -v http://agent-ip:8080/health

# Check firewall rules
iptables -L | grep ducla
```

### 2. Rate Limited
```bash
# Check rate limiting logs
journalctl -u ducla-agent | grep "rate limit"

# Adjust rate limits in config
requests_per_min: 2000  # Increase limit
```

### 3. DNS Resolution
```yaml
master:
  url: "https://master.internal.company.com:8443"  # Use FQDN
  # or
  url: "https://10.0.1.100:8443"  # Use IP directly
```

## Security Considerations

### 1. Never Allow 0.0.0.0/0 in Production
```yaml
# ❌ NEVER do this in production
allowed_ips:
  - "0.0.0.0/0"  # Allows entire internet

# ✅ Be specific
allowed_ips:
  - "10.0.1.0/24"  # Only your subnet
```

### 2. Use TLS for All Communications
```yaml
master:
  url: "https://master.example.com:8443"  # HTTPS
  # or
  url: "wss://master.example.com:8443"    # WSS
  # or  
  url: "grpcs://master.example.com:8443"  # gRPC with TLS
```

### 3. Regular Security Audits
```bash
# Check current connections
netstat -an | grep :8080

# Check firewall rules
iptables -L -n

# Monitor access logs
tail -f /var/log/ducla/audit.log
```

## Environment Variables Override

You can override network settings using environment variables:

```bash
# Allow specific IP
export DUCLA_ALLOWED_IPS="10.0.1.100/32,192.168.1.0/24"

# Enable firewall
export DUCLA_FIREWALL_ENABLED=true

# Set rate limit
export DUCLA_RATE_LIMIT=500
```