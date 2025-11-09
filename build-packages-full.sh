#!/bin/bash

# Full package builder for Ducla Cloud Agent with systemd service

set -e

VERSION="1.0.0"
PACKAGE_NAME="ducla-agent"

echo "🚀 Building full packages for Ducla Cloud Agent v${VERSION}..."

# Ensure we have the binary
if [ ! -f "dist/ducla-agent-v1.0.0" ]; then
    echo "📦 Building binary first..."
    ./build-v1.sh
fi

# Create package structure
echo "📁 Creating full package structure..."
mkdir -p pkg-full/{deb,rpm}/{usr/bin,usr/share/man/man1,etc/ducla,etc/systemd/system,var/log/ducla,opt/ducla}

# Copy binary
cp dist/ducla-agent-v1.0.0 pkg-full/deb/usr/bin/ducla-agent
cp dist/ducla-agent-v1.0.0 pkg-full/rpm/usr/bin/ducla-agent

# Copy man page
cp docs/ducla-agent.1 pkg-full/deb/usr/share/man/man1/
cp docs/ducla-agent.1 pkg-full/rpm/usr/share/man/man1/

# Create default config
cat > pkg-full/deb/etc/ducla/agent.yaml << 'EOF'
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

# Copy config to RPM
cp pkg-full/deb/etc/ducla/agent.yaml pkg-full/rpm/etc/ducla/

# Create systemd service
cat > pkg-full/deb/etc/systemd/system/ducla-agent.service << 'EOF'
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

# Copy service to RPM
cp pkg-full/deb/etc/systemd/system/ducla-agent.service pkg-full/rpm/etc/systemd/system/

# Create installation scripts
mkdir -p scripts-full

# Pre-install script
cat > scripts-full/pre-install.sh << 'EOF'
#!/bin/bash
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
cat > scripts-full/post-install.sh << 'EOF'
#!/bin/bash
# Set permissions
chown ducla:ducla /usr/bin/ducla-agent
chmod 755 /usr/bin/ducla-agent

# Set config permissions
chown root:ducla /etc/ducla/agent.yaml
chmod 640 /etc/ducla/agent.yaml

# Reload systemd and enable service
systemctl daemon-reload
systemctl enable ducla-agent

echo ""
echo "🎉 Ducla Cloud Agent installed successfully!"
echo ""
echo "📋 Next steps:"
echo "  1. Edit config: sudo nano /etc/ducla/agent.yaml"
echo "  2. Start service: sudo systemctl start ducla-agent"
echo "  3. Check status: sudo systemctl status ducla-agent"
echo "  4. View logs: sudo journalctl -u ducla-agent -f"
echo ""
echo "📊 API endpoints:"
echo "  Status: http://localhost:8080/api/v1/status"
echo "  Health: http://localhost:8081/health"
echo "  Metrics: http://localhost:9090/metrics"
EOF

# Pre-remove script
cat > scripts-full/pre-remove.sh << 'EOF'
#!/bin/bash
if systemctl is-active --quiet ducla-agent; then
    systemctl stop ducla-agent
fi

if systemctl is-enabled --quiet ducla-agent; then
    systemctl disable ducla-agent
fi
EOF

# Post-remove script
cat > scripts-full/post-remove.sh << 'EOF'
#!/bin/bash
systemctl daemon-reload

echo ""
echo "📦 Ducla Cloud Agent removed."
echo "💾 User data in /opt/ducla and /var/log/ducla preserved."
echo "🗑️  To completely remove: sudo userdel ducla && sudo rm -rf /opt/ducla /var/log/ducla"
EOF

# Make scripts executable
chmod +x scripts-full/*.sh

# Build DEB package with full features
echo "📦 Building full DEB package..."
fpm -s dir -t deb \
    --name "${PACKAGE_NAME}" \
    --version "${VERSION}" \
    --architecture "amd64" \
    --description "Ducla Cloud Agent - High-performance cloud agent for distributed task execution and system monitoring" \
    --maintainer "mandỵhades <mandỵhades@hotmail.com.vn>" \
    --url "https://github.com/duclacloud/DUCLA-CLOUD-AGENT" \
    --license "MIT" \
    --category "admin" \
    --depends "systemd" \
    --deb-systemd pkg-full/deb/etc/systemd/system/ducla-agent.service \
    --before-install scripts-full/pre-install.sh \
    --after-install scripts-full/post-install.sh \
    --before-remove scripts-full/pre-remove.sh \
    --after-remove scripts-full/post-remove.sh \
    --config-files /etc/ducla/agent.yaml \
    --directories /opt/ducla \
    --directories /var/log/ducla \
    --package pkg-full/deb/ \
    -C pkg-full/deb \
    .

# Build RPM package with full features
echo "📦 Building full RPM package..."
fpm -s dir -t rpm \
    --name "${PACKAGE_NAME}" \
    --version "${VERSION}" \
    --architecture "amd64" \
    --description "Ducla Cloud Agent - High-performance cloud agent for distributed task execution and system monitoring" \
    --maintainer "mandỵhades <mandỵhades@hotmail.com.vn>" \
    --url "https://github.com/duclacloud/DUCLA-CLOUD-AGENT" \
    --license "MIT" \
    --category "System Environment/Daemons" \
    --depends "systemd" \
    --before-install scripts-full/pre-install.sh \
    --after-install scripts-full/post-install.sh \
    --before-remove scripts-full/pre-remove.sh \
    --after-remove scripts-full/post-remove.sh \
    --config-files /etc/ducla/agent.yaml \
    --directories /opt/ducla \
    --directories /var/log/ducla \
    --package pkg-full/rpm/ \
    -C pkg-full/rpm \
    .

echo ""
echo "🎉 Full packages built successfully!"
echo ""
echo "📦 DEB Package (with systemd service):"
ls -la pkg-full/deb/*.deb
echo ""
echo "📦 RPM Package (with systemd service):"
ls -la pkg-full/rpm/*.rpm
echo ""
echo "📋 Installation commands:"
echo "  DEB: sudo dpkg -i pkg-full/deb/${PACKAGE_NAME}_${VERSION}_amd64.deb"
echo "  RPM: sudo rpm -ivh pkg-full/rpm/${PACKAGE_NAME}-${VERSION}-1.amd64.rpm"
echo ""
echo "🔧 After installation:"
echo "  sudo systemctl start ducla-agent"
echo "  sudo systemctl status ducla-agent"
echo "  curl http://localhost:8080/api/v1/status"

# Keep packages for testing
echo "📁 Packages saved in pkg-full/ directory"