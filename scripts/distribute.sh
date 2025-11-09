#!/bin/bash

# Distribution script for Ducla Cloud Agent packages

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

DIST_DIR="dist"
RELEASE_DIR="release"
VERSION=${VERSION:-$(git describe --tags --always --dirty 2>/dev/null | sed 's/^v//' || echo "1.0.0")}

echo -e "${GREEN}🚀 Distributing Ducla Cloud Agent v${VERSION}${NC}"
echo ""

# Create release directory
mkdir -p ${RELEASE_DIR}

# Copy packages to release directory
echo -e "${YELLOW}📦 Copying packages...${NC}"
cp ${DIST_DIR}/*.deb ${RELEASE_DIR}/ 2>/dev/null || true
cp ${DIST_DIR}/*.rpm ${RELEASE_DIR}/ 2>/dev/null || true
cp ${DIST_DIR}/ducla-agent-linux-* ${RELEASE_DIR}/ 2>/dev/null || true

# Create checksums
echo -e "${YELLOW}🔐 Generating checksums...${NC}"
cd ${RELEASE_DIR}
sha256sum * > SHA256SUMS
cd ..

# Create installation script
echo -e "${YELLOW}📝 Creating installation script...${NC}"
cat > ${RELEASE_DIR}/install.sh <<'EOF'
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
EOF

chmod +x ${RELEASE_DIR}/install.sh

# Create uninstall script
echo -e "${YELLOW}🗑️ Creating uninstall script...${NC}"
cat > ${RELEASE_DIR}/uninstall.sh <<'EOF'
#!/bin/bash

# Ducla Cloud Agent Uninstall Script

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

SERVICE_NAME="ducla-agent"

echo -e "${YELLOW}🗑️ Uninstalling Ducla Cloud Agent...${NC}"

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo -e "${RED}❌ This script must be run as root${NC}"
   exit 1
fi

# Stop and disable service
if systemctl is-active --quiet $SERVICE_NAME; then
    echo -e "${YELLOW}⏹️ Stopping service...${NC}"
    systemctl stop $SERVICE_NAME
fi

if systemctl is-enabled --quiet $SERVICE_NAME; then
    echo -e "${YELLOW}🔧 Disabling service...${NC}"
    systemctl disable $SERVICE_NAME
fi

# Detect OS and remove package
if [ -f /etc/os-release ]; then
    . /etc/os-release
    OS=$ID
    
    case $OS in
        ubuntu|debian)
            if dpkg -l | grep -q ducla-agent; then
                echo -e "${YELLOW}📦 Removing DEB package...${NC}"
                dpkg -r ducla-agent
            fi
            ;;
        centos|rhel|fedora)
            if rpm -q ducla-agent >/dev/null 2>&1; then
                echo -e "${YELLOW}📦 Removing RPM package...${NC}"
                rpm -e ducla-agent
            fi
            ;;
        *)
            echo -e "${YELLOW}🗑️ Removing manual installation...${NC}"
            rm -rf /opt/ducla
            rm -f /etc/systemd/system/$SERVICE_NAME.service
            rm -rf /etc/ducla
            systemctl daemon-reload
            ;;
    esac
fi

# Remove user (optional)
if id "ducla" &>/dev/null; then
    echo -e "${YELLOW}👤 Removing user 'ducla'...${NC}"
    userdel ducla 2>/dev/null || true
fi

# Remove logs (optional)
if [ -d "/var/log/ducla" ]; then
    echo -e "${YELLOW}📝 Removing logs...${NC}"
    rm -rf /var/log/ducla
fi

echo -e "${GREEN}✅ Ducla Cloud Agent uninstalled successfully!${NC}"
EOF

chmod +x ${RELEASE_DIR}/uninstall.sh

# Create README
echo -e "${YELLOW}📖 Creating README...${NC}"
cat > ${RELEASE_DIR}/README.md <<EOF
# Ducla Cloud Agent v${VERSION}

High-performance cloud agent for distributed task execution and system monitoring.

## 📦 Available Packages

- **Ubuntu/Debian**: \`ducla-agent_${VERSION}_amd64.deb\`
- **RedHat/CentOS/Fedora**: \`ducla-agent-${VERSION}.x86_64.rpm\`
- **Linux Binary**: \`ducla-agent-linux-amd64\`

## 🚀 Quick Installation

### Automatic Installation (Recommended)
\`\`\`bash
# Download and run installer
curl -sSL https://github.com/${GITHUB_REPO}/releases/latest/download/install.sh | sudo bash
\`\`\`

### Manual Installation

#### Ubuntu/Debian
\`\`\`bash
sudo dpkg -i ducla-agent_${VERSION}_amd64.deb
sudo systemctl enable ducla-agent
sudo systemctl start ducla-agent
\`\`\`

#### RedHat/CentOS/Fedora
\`\`\`bash
sudo rpm -ivh ducla-agent-${VERSION}.x86_64.rpm
sudo systemctl enable ducla-agent
sudo systemctl start ducla-agent
\`\`\`

#### Binary Installation
\`\`\`bash
# Run the install script in this directory
sudo ./install.sh
\`\`\`

## 🔧 Configuration

Edit \`/etc/ducla/agent.yaml\` to configure the agent:

\`\`\`yaml
agent:
  id: "\${HOSTNAME}"
  name: "\${HOSTNAME}"
  environment: "production"

api:
  http:
    enabled: true
    port: 8080

logging:
  level: "info"
  format: "json"
\`\`\`

## 📊 API Endpoints

- **Status**: http://localhost:8080/api/v1/status
- **Health**: http://localhost:8081/health  
- **Metrics**: http://localhost:9090/metrics

## 🔍 Service Management

\`\`\`bash
# Check status
sudo systemctl status ducla-agent

# View logs
sudo journalctl -u ducla-agent -f

# Restart service
sudo systemctl restart ducla-agent
\`\`\`

## 🗑️ Uninstallation

\`\`\`bash
# Run uninstall script
sudo ./uninstall.sh
\`\`\`

## 🔐 Security

- Service runs as non-root user \`ducla\`
- Configuration files have restricted permissions
- TLS support for secure communication

## 📋 System Requirements

- **OS**: Linux (Ubuntu 18.04+, CentOS 7+, RHEL 7+, Fedora 30+)
- **Architecture**: x86_64 (amd64)
- **RAM**: 512MB minimum
- **Disk**: 100MB free space

## 🆘 Support

- **Documentation**: https://github.com/${GITHUB_REPO}
- **Issues**: https://github.com/${GITHUB_REPO}/issues
- **Email**: support@ducla.cloud

---
**Version**: ${VERSION}  
**Build Date**: $(date '+%Y-%m-%d %H:%M:%S UTC')
EOF

# List release files
echo ""
echo -e "${GREEN}✅ Distribution complete!${NC}"
echo ""
echo -e "${BLUE}📁 Release files:${NC}"
ls -lh ${RELEASE_DIR}/

echo ""
echo -e "${BLUE}🚀 Next steps:${NC}"
echo "1. Test installation: sudo ./${RELEASE_DIR}/install.sh"
echo "2. Upload to GitHub releases"
echo "3. Update package repositories"
echo "4. Announce release"