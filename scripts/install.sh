#!/bin/bash

# Installation script for Ducla Cloud Agent
# Supports Linux and macOS

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Configuration
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/ducla"
DATA_DIR="/opt/ducla/data"
LOG_DIR="/var/log/ducla"
SERVICE_USER="ducla"
SERVICE_GROUP="ducla"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo -e "${RED}Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

echo -e "${GREEN}Ducla Cloud Agent Installer${NC}"
echo "OS: $OS"
echo "Architecture: $ARCH"
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo -e "${RED}Please run as root or with sudo${NC}"
    exit 1
fi

# Download latest release
echo -e "${YELLOW}Downloading latest release...${NC}"
LATEST_VERSION=$(curl -s https://api.github.com/repos/your-org/ducla-cloud-agent/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
    echo -e "${RED}Failed to get latest version${NC}"
    exit 1
fi

DOWNLOAD_URL="https://github.com/your-org/ducla-cloud-agent/releases/download/${LATEST_VERSION}/ducla-agent-${LATEST_VERSION}-${OS}-${ARCH}.tar.gz"

echo "Downloading version: $LATEST_VERSION"
curl -L -o /tmp/ducla-agent.tar.gz "$DOWNLOAD_URL"

# Extract binary
echo -e "${YELLOW}Installing binary...${NC}"
tar -xzf /tmp/ducla-agent.tar.gz -C /tmp
mv /tmp/ducla-agent-* "$INSTALL_DIR/ducla-agent"
chmod +x "$INSTALL_DIR/ducla-agent"

# Create user and group
echo -e "${YELLOW}Creating service user...${NC}"
if ! id "$SERVICE_USER" &>/dev/null; then
    useradd -r -s /bin/false -d /opt/ducla "$SERVICE_USER"
fi

# Create directories
echo -e "${YELLOW}Creating directories...${NC}"
mkdir -p "$CONFIG_DIR" "$DATA_DIR" "$LOG_DIR"
chown -R "$SERVICE_USER:$SERVICE_GROUP" "$DATA_DIR" "$LOG_DIR"

# Install configuration
if [ ! -f "$CONFIG_DIR/agent.yaml" ]; then
    echo -e "${YELLOW}Installing default configuration...${NC}"
    cat > "$CONFIG_DIR/agent.yaml" <<'EOF'
# Ducla Cloud Agent Configuration
agent:
  id: "${HOSTNAME}"
  name: "${HOSTNAME}"
  environment: "production"

master:
  url: "${DUCLA_MASTER_URL}"
  token: "${DUCLA_AGENT_TOKEN}"

api:
  http:
    enabled: true
    port: 8080
  grpc:
    enabled: true
    port: 8443

logging:
  level: "info"
  format: "json"
  output: "stdout"
EOF
    chmod 644 "$CONFIG_DIR/agent.yaml"
fi

# Install systemd service
if command -v systemctl &> /dev/null; then
    echo -e "${YELLOW}Installing systemd service...${NC}"
    cat > /etc/systemd/system/ducla-agent.service <<EOF
[Unit]
Description=Ducla Cloud Agent
After=network.target
Documentation=https://github.com/your-org/ducla-cloud-agent

[Service]
Type=simple
User=$SERVICE_USER
Group=$SERVICE_GROUP
ExecStart=$INSTALL_DIR/ducla-agent --config $CONFIG_DIR/agent.yaml
Restart=on-failure
RestartSec=10s

# Security
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=$DATA_DIR $LOG_DIR

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload
    systemctl enable ducla-agent
    
    echo -e "${GREEN}✓ Systemd service installed${NC}"
fi

# Cleanup
rm -f /tmp/ducla-agent.tar.gz

echo ""
echo -e "${GREEN}Installation complete!${NC}"
echo ""
echo "Next steps:"
echo "1. Edit configuration: $CONFIG_DIR/agent.yaml"
echo "2. Set environment variables:"
echo "   export DUCLA_MASTER_URL=https://master.ducla.cloud"
echo "   export DUCLA_AGENT_TOKEN=your-token"
echo "3. Start the service:"
echo "   systemctl start ducla-agent"
echo "4. Check status:"
echo "   systemctl status ducla-agent"
