#!/bin/bash

# DEB packaging script for Ubuntu/Debian

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Configuration
APP_NAME="ducla-agent"
VERSION=${VERSION:-$(git describe --tags --always --dirty 2>/dev/null | sed 's/^v//' || echo "1.0.0")}
ARCH=${ARCH:-amd64}
MAINTAINER="Ducla Team <team@ducla.cloud>"

DEB_BUILD_DIR="build/deb"
PACKAGE_NAME="${APP_NAME}_${VERSION}_${ARCH}"
PACKAGE_DIR="${DEB_BUILD_DIR}/${PACKAGE_NAME}"

echo -e "${GREEN}Building DEB package${NC}"
echo "Version: ${VERSION}"
echo "Architecture: ${ARCH}"
echo ""

# Check for dpkg-deb
if ! command -v dpkg-deb &> /dev/null; then
    echo -e "${RED}dpkg-deb not found. Install it with:${NC}"
    echo "  sudo apt-get install dpkg-dev"
    exit 1
fi

# Clean and create build directory
echo -e "${YELLOW}Creating DEB build directory...${NC}"
rm -rf ${DEB_BUILD_DIR}
mkdir -p ${PACKAGE_DIR}/{DEBIAN,usr/bin,etc/ducla,lib/systemd/system,opt/ducla/data,var/log/ducla}

# Build binary
echo -e "${YELLOW}Building binary...${NC}"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X 'main.Version=${VERSION}'" \
    -o ${PACKAGE_DIR}/usr/bin/${APP_NAME} \
    ./cmd/agent

chmod 755 ${PACKAGE_DIR}/usr/bin/${APP_NAME}

# Create control file
echo -e "${YELLOW}Creating control file...${NC}"
cat > ${PACKAGE_DIR}/DEBIAN/control <<EOF
Package: ${APP_NAME}
Version: ${VERSION}
Section: utils
Priority: optional
Architecture: ${ARCH}
Maintainer: ${MAINTAINER}
Description: Ducla Cloud Agent for distributed task execution
 Ducla Cloud Agent is a high-performance agent for distributed task execution,
 system monitoring, and cloud resource management.
 .
 Features:
  - High performance with low resource usage
  - Secure with JWT authentication and RBAC
  - Built-in metrics and health checks
  - Plugin system for extensibility
  - Multi-protocol support (HTTP REST and gRPC)
Homepage: https://github.com/your-org/ducla-cloud-agent
Depends: systemd
EOF

# Create preinst script
cat > ${PACKAGE_DIR}/DEBIAN/preinst <<'EOF'
#!/bin/bash
set -e

# Create user and group
if ! getent group ducla >/dev/null; then
    addgroup --system ducla
fi

if ! getent passwd ducla >/dev/null; then
    adduser --system --ingroup ducla --home /opt/ducla \
            --no-create-home --shell /usr/sbin/nologin \
            --gecos "Ducla Cloud Agent" ducla
fi

exit 0
EOF
chmod 755 ${PACKAGE_DIR}/DEBIAN/preinst

# Create postinst script
cat > ${PACKAGE_DIR}/DEBIAN/postinst <<'EOF'
#!/bin/bash
set -e

# Set ownership
chown -R ducla:ducla /opt/ducla /var/log/ducla

# Reload systemd
systemctl daemon-reload

# Enable service
systemctl enable ducla-agent.service || true

echo ""
echo "Ducla Cloud Agent installed successfully!"
echo ""
echo "Next steps:"
echo "1. Edit configuration: /etc/ducla/agent.yaml"
echo "2. Set environment variables in /etc/default/ducla-agent"
echo "3. Start the service: systemctl start ducla-agent"
echo "4. Check status: systemctl status ducla-agent"
echo ""

exit 0
EOF
chmod 755 ${PACKAGE_DIR}/DEBIAN/postinst

# Create prerm script
cat > ${PACKAGE_DIR}/DEBIAN/prerm <<'EOF'
#!/bin/bash
set -e

# Stop service
systemctl stop ducla-agent.service || true

exit 0
EOF
chmod 755 ${PACKAGE_DIR}/DEBIAN/prerm

# Create postrm script
cat > ${PACKAGE_DIR}/DEBIAN/postrm <<'EOF'
#!/bin/bash
set -e

if [ "$1" = "purge" ]; then
    # Remove user and group
    if getent passwd ducla >/dev/null; then
        deluser ducla || true
    fi
    
    if getent group ducla >/dev/null; then
        delgroup ducla || true
    fi
    
    # Remove data directories
    rm -rf /opt/ducla /var/log/ducla /etc/ducla
fi

# Reload systemd
systemctl daemon-reload || true

exit 0
EOF
chmod 755 ${PACKAGE_DIR}/DEBIAN/postrm

# Create configuration file
cat > ${PACKAGE_DIR}/etc/ducla/agent.yaml <<'EOF'
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

# Create systemd service file
cat > ${PACKAGE_DIR}/lib/systemd/system/${APP_NAME}.service <<EOF
[Unit]
Description=Ducla Cloud Agent
After=network.target
Documentation=https://github.com/your-org/ducla-cloud-agent

[Service]
Type=simple
User=ducla
Group=ducla
EnvironmentFile=-/etc/default/${APP_NAME}
ExecStart=/usr/bin/${APP_NAME} --config /etc/ducla/agent.yaml
Restart=on-failure
RestartSec=10s

# Security
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/ducla /var/log/ducla

[Install]
WantedBy=multi-user.target
EOF

# Create environment file template
mkdir -p ${PACKAGE_DIR}/etc/default
cat > ${PACKAGE_DIR}/etc/default/${APP_NAME} <<'EOF'
# Ducla Cloud Agent environment variables
# Uncomment and set these values

# DUCLA_MASTER_URL=https://master.ducla.cloud
# DUCLA_AGENT_TOKEN=your-token-here
# DUCLA_JWT_SECRET=your-secret-here
# DUCLA_LOG_LEVEL=info
EOF

# Create conffiles
echo "/etc/ducla/agent.yaml" > ${PACKAGE_DIR}/DEBIAN/conffiles
echo "/etc/default/${APP_NAME}" >> ${PACKAGE_DIR}/DEBIAN/conffiles

# Build DEB package
echo -e "${YELLOW}Building DEB package...${NC}"
dpkg-deb --build ${PACKAGE_DIR}

# Move DEB to dist directory
mkdir -p dist
mv ${DEB_BUILD_DIR}/${PACKAGE_NAME}.deb dist/

echo ""
echo -e "${GREEN}DEB package built successfully!${NC}"
ls -lh dist/*.deb

echo ""
echo "To install:"
echo "  sudo dpkg -i dist/${PACKAGE_NAME}.deb"
echo "  sudo apt-get install -f  # Install dependencies if needed"
