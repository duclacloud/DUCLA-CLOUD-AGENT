#!/bin/bash

# Build DEB and RPM packages for Ducla Cloud Agent

set -e

# Version information
VERSION="1.0.0"
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S_UTC')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Package information
PACKAGE_NAME="ducla-agent"
DESCRIPTION="Ducla Cloud Agent - High-performance cloud agent for distributed task execution"
MAINTAINER="mandỵhades <mandỵhades@hotmail.com.vn>"
HOMEPAGE="https://github.com/duclacloud/DUCLA-CLOUD-AGENT"
ARCH="amd64"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_header() {
    echo -e "\n${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}\n"
}

print_step() {
    echo -e "${YELLOW}➤ $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check dependencies
check_dependencies() {
    print_step "Checking dependencies..."
    
    # Check for fpm
    if ! command -v fpm &> /dev/null; then
        print_error "fpm is required but not installed"
        echo "Install with: gem install fpm"
        echo "Or on Ubuntu: sudo apt install ruby-dev build-essential && gem install fpm"
        exit 1
    fi
    
    # Check for rpmbuild (for RPM)
    if ! command -v rpmbuild &> /dev/null; then
        print_step "Installing rpmbuild..."
        sudo apt update
        sudo apt install -y rpm
    fi
    
    print_success "Dependencies checked"
}

# Build binary
build_binary() {
    print_step "Building binary..."
    
    # Build flags
    LDFLAGS="-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.gitCommit=${GIT_COMMIT}"
    
    # Build for Linux AMD64
    GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o dist/${PACKAGE_NAME} ./cmd/agent
    
    # Make executable
    chmod +x dist/${PACKAGE_NAME}
    
    print_success "Binary built: dist/${PACKAGE_NAME}"
}

# Create package structure
create_package_structure() {
    print_step "Creating package structure..."
    
    # Create directories
    mkdir -p packaging/{deb,rpm}/{usr/bin,etc/ducla,etc/systemd/system,var/log/ducla,opt/ducla}
    
    # Copy binary
    cp dist/${PACKAGE_NAME} packaging/deb/usr/bin/
    cp dist/${PACKAGE_NAME} packaging/rpm/usr/bin/
    
    # Create default config
    cat > packaging/deb/etc/ducla/agent.yaml << 'EOF'
# Ducla Cloud Agent Configuration
agent:
  id: "${HOSTNAME}"
  name: "${HOSTNAME}"
  environment: "production"
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
  data_dir: "/opt/ducla/data"
  temp_dir: "/tmp/ducla"
  cleanup:
    enabled: true
    interval: 1h
    max_age: 24h

logging:
  level: "info"
  format: "json"
  output: "/var/log/ducla/agent.log"

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
  task_timeout: 300s
EOF

    # Copy config to RPM structure
    cp packaging/deb/etc/ducla/agent.yaml packaging/rpm/etc/ducla/

    # Create systemd service file
    cat > packaging/deb/etc/systemd/system/ducla-agent.service << 'EOF'
[Unit]
Description=Ducla Cloud Agent
Documentation=https://github.com/duclacloud/DUCLA-CLOUD-AGENT
After=network.target

[Service]
Type=simple
User=ducla
Group=ducla
WorkingDirectory=/opt/ducla
ExecStart=/usr/bin/ducla-agent -config /etc/ducla/agent.yaml
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=ducla-agent

# Security settings
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/ducla /var/log/ducla /tmp/ducla

[Install]
WantedBy=multi-user.target
EOF

    # Copy service file to RPM structure
    cp packaging/deb/etc/systemd/system/ducla-agent.service packaging/rpm/etc/systemd/system/

    print_success "Package structure created"
}

# Build DEB package
build_deb() {
    print_step "Building DEB package..."
    
    fpm -s dir -t deb \
        --name "${PACKAGE_NAME}" \
        --version "${VERSION}" \
        --architecture "${ARCH}" \
        --description "${DESCRIPTION}" \
        --maintainer "${MAINTAINER}" \
        --url "${HOMEPAGE}" \
        --license "MIT" \
        --category "admin" \
        --depends "systemd" \
        --deb-systemd packaging/deb/etc/systemd/system/ducla-agent.service \
        --before-install scripts/pre-install.sh \
        --after-install scripts/post-install.sh \
        --before-remove scripts/pre-remove.sh \
        --after-remove scripts/post-remove.sh \
        --config-files /etc/ducla/agent.yaml \
        --directories /opt/ducla \
        --directories /var/log/ducla \
        --package dist/ \
        -C packaging/deb \
        .
    
    print_success "DEB package created: dist/${PACKAGE_NAME}_${VERSION}_${ARCH}.deb"
}

# Build RPM package
build_rpm() {
    print_step "Building RPM package..."
    
    fpm -s dir -t rpm \
        --name "${PACKAGE_NAME}" \
        --version "${VERSION}" \
        --architecture "${ARCH}" \
        --description "${DESCRIPTION}" \
        --maintainer "${MAINTAINER}" \
        --url "${HOMEPAGE}" \
        --license "MIT" \
        --category "System Environment/Daemons" \
        --depends "systemd" \
        --before-install scripts/pre-install.sh \
        --after-install scripts/post-install.sh \
        --before-remove scripts/pre-remove.sh \
        --after-remove scripts/post-remove.sh \
        --config-files /etc/ducla/agent.yaml \
        --directories /opt/ducla \
        --directories /var/log/ducla \
        --package dist/ \
        -C packaging/rpm \
        .
    
    print_success "RPM package created: dist/${PACKAGE_NAME}-${VERSION}-1.${ARCH}.rpm"
}

# Create installation scripts
create_install_scripts() {
    print_step "Creating installation scripts..."
    
    mkdir -p scripts
    
    # Pre-install script
    cat > scripts/pre-install.sh << 'EOF'
#!/bin/bash
# Pre-installation script for Ducla Cloud Agent

# Create ducla user and group
if ! getent group ducla >/dev/null; then
    groupadd --system ducla
fi

if ! getent passwd ducla >/dev/null; then
    useradd --system --gid ducla --home-dir /opt/ducla --shell /bin/false ducla
fi

# Create directories
mkdir -p /opt/ducla /var/log/ducla /tmp/ducla
chown ducla:ducla /opt/ducla /var/log/ducla /tmp/ducla
chmod 755 /opt/ducla /var/log/ducla /tmp/ducla
EOF

    # Post-install script
    cat > scripts/post-install.sh << 'EOF'
#!/bin/bash
# Post-installation script for Ducla Cloud Agent

# Set permissions
chown ducla:ducla /usr/bin/ducla-agent
chmod 755 /usr/bin/ducla-agent

# Set config permissions
chown root:ducla /etc/ducla/agent.yaml
chmod 640 /etc/ducla/agent.yaml

# Reload systemd and enable service
systemctl daemon-reload
systemctl enable ducla-agent

echo "Ducla Cloud Agent installed successfully!"
echo "Configuration: /etc/ducla/agent.yaml"
echo "Start service: sudo systemctl start ducla-agent"
echo "View logs: sudo journalctl -u ducla-agent -f"
EOF

    # Pre-remove script
    cat > scripts/pre-remove.sh << 'EOF'
#!/bin/bash
# Pre-removal script for Ducla Cloud Agent

# Stop and disable service
if systemctl is-active --quiet ducla-agent; then
    systemctl stop ducla-agent
fi

if systemctl is-enabled --quiet ducla-agent; then
    systemctl disable ducla-agent
fi
EOF

    # Post-remove script
    cat > scripts/post-remove.sh << 'EOF'
#!/bin/bash
# Post-removal script for Ducla Cloud Agent

# Reload systemd
systemctl daemon-reload

# Note: We don't remove user/group or data directories
# in case user wants to reinstall later
echo "Ducla Cloud Agent removed."
echo "User data in /opt/ducla and /var/log/ducla preserved."
echo "To completely remove: sudo userdel ducla && sudo rm -rf /opt/ducla /var/log/ducla"
EOF

    # Make scripts executable
    chmod +x scripts/*.sh
    
    print_success "Installation scripts created"
}

# Cleanup
cleanup() {
    print_step "Cleaning up temporary files..."
    rm -rf packaging
    print_success "Cleanup completed"
}

# Main function
main() {
    print_header "🚀 BUILDING DUCLA CLOUD AGENT PACKAGES"
    
    echo "Building packages for Ducla Cloud Agent v${VERSION}"
    echo "Target: ${ARCH} architecture"
    echo "Packages: DEB and RPM"
    echo ""
    
    # Create dist directory
    mkdir -p dist
    
    # Build process
    check_dependencies
    build_binary
    create_install_scripts
    create_package_structure
    build_deb
    build_rpm
    cleanup
    
    print_header "🎉 PACKAGE BUILD COMPLETED"
    
    echo "Packages created:"
    ls -la dist/*.deb dist/*.rpm 2>/dev/null || echo "No packages found"
    
    echo ""
    echo "Installation commands:"
    echo "  DEB: sudo dpkg -i dist/${PACKAGE_NAME}_${VERSION}_${ARCH}.deb"
    echo "  RPM: sudo rpm -ivh dist/${PACKAGE_NAME}-${VERSION}-1.${ARCH}.rpm"
    echo ""
    echo "After installation:"
    echo "  sudo systemctl start ducla-agent"
    echo "  sudo systemctl status ducla-agent"
}

# Handle Ctrl+C
trap 'echo -e "\n${YELLOW}Build interrupted${NC}"; cleanup; exit 1' INT

# Run main function
main "$@"