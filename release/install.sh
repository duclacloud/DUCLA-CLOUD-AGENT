#!/bin/bash

# Ducla Cloud Agent Installation Script

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

GITHUB_REPO="duclacloud/DUCLA-CLOUD-AGENT"
INSTALL_DIR="/opt/ducla"
CONFIG_DIR="/etc/ducla"
SERVICE_NAME="ducla-agent"

echo -e "${BLUE}🚀 Ducla Cloud Agent Installer${NC}"
echo ""

# Detect OS
if [ -f /etc/os-release ]; then
    . /etc/os-release
    OS=$ID
    VERSION_ID=$VERSION_ID
else
    echo -e "${RED}❌ Cannot detect OS${NC}"
    exit 1
fi

echo -e "${YELLOW}📋 Detected OS: $OS $VERSION_ID${NC}"

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo -e "${RED}❌ This script must be run as root${NC}"
   exit 1
fi

# Install based on OS
case $OS in
    ubuntu|debian)
        echo -e "${YELLOW}📦 Installing DEB package...${NC}"
        if [ -f "ducla-agent_*_amd64.deb" ]; then
            dpkg -i ducla-agent_*_amd64.deb
            apt-get install -f -y
        else
            echo -e "${RED}❌ DEB package not found${NC}"
            exit 1
        fi
        ;;
    centos|rhel|fedora)
        echo -e "${YELLOW}📦 Installing RPM package...${NC}"
        if [ -f "ducla-agent-*.x86_64.rpm" ]; then
            rpm -ivh ducla-agent-*.x86_64.rpm
        else
            echo -e "${RED}❌ RPM package not found${NC}"
            exit 1
        fi
        ;;
    *)
        echo -e "${YELLOW}📦 Installing binary...${NC}"
        if [ -f "ducla-agent-linux-amd64" ]; then
            # Create directories
            mkdir -p $INSTALL_DIR/bin
            mkdir -p $CONFIG_DIR
            mkdir -p /var/log/ducla
            
            # Install binary
            cp ducla-agent-linux-amd64 $INSTALL_DIR/bin/ducla-agent
            chmod +x $INSTALL_DIR/bin/ducla-agent
            
            # Create user
            useradd -r -s /bin/false -d $INSTALL_DIR ducla 2>/dev/null || true
            
            # Set permissions
            chown -R ducla:ducla $INSTALL_DIR
            chown -R ducla:ducla /var/log/ducla
            
            # Create config
            cat > $CONFIG_DIR/agent.yaml <<'EOFCONFIG'
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
  output: "/var/log/ducla/agent.log"
EOFCONFIG
            
            # Create systemd service
            cat > /etc/systemd/system/$SERVICE_NAME.service <<'EOFSERVICE'
[Unit]
Description=Ducla Cloud Agent
After=network.target

[Service]
Type=simple
User=ducla
Group=ducla
ExecStart=/opt/ducla/bin/ducla-agent --config /etc/ducla/agent.yaml
Restart=on-failure
RestartSec=10s

[Install]
WantedBy=multi-user.target
EOFSERVICE
            
            systemctl daemon-reload
        else
            echo -e "${RED}❌ Binary not found${NC}"
            exit 1
        fi
        ;;
esac

# Enable and start service
echo -e "${YELLOW}🔧 Configuring service...${NC}"
systemctl enable $SERVICE_NAME
systemctl start $SERVICE_NAME

# Check status
sleep 2
if systemctl is-active --quiet $SERVICE_NAME; then
    echo -e "${GREEN}✅ Ducla Cloud Agent installed and started successfully!${NC}"
    echo ""
    echo -e "${BLUE}📊 Service Status:${NC}"
    systemctl status $SERVICE_NAME --no-pager -l
    echo ""
    echo -e "${BLUE}🔗 API Endpoints:${NC}"
    echo "  • Status: http://localhost:8080/api/v1/status"
    echo "  • Health: http://localhost:8081/health"
    echo "  • Metrics: http://localhost:9090/metrics"
else
    echo -e "${RED}❌ Failed to start Ducla Cloud Agent${NC}"
    echo "Check logs: journalctl -u $SERVICE_NAME -f"
    exit 1
fi
