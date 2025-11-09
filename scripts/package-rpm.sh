#!/bin/bash

# RPM packaging script for RedHat/CentOS/Fedora

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Configuration
APP_NAME="ducla-agent"
VERSION=${VERSION:-$(git describe --tags --always --dirty 2>/dev/null | sed 's/^v//' | sed 's/-/_/g' || echo "1.0.0")}
RELEASE=${RELEASE:-1}
ARCH=${ARCH:-x86_64}

RPM_BUILD_DIR="build/rpm"
SPEC_FILE="${RPM_BUILD_DIR}/SPECS/${APP_NAME}.spec"

echo -e "${GREEN}Building RPM package${NC}"
echo "Version: ${VERSION}"
echo "Release: ${RELEASE}"
echo "Architecture: ${ARCH}"
echo ""

# Check for rpmbuild
if ! command -v rpmbuild &> /dev/null; then
    echo -e "${RED}rpmbuild not found. Install it with:${NC}"
    echo "  RHEL/CentOS: sudo yum install rpm-build"
    echo "  Fedora: sudo dnf install rpm-build"
    exit 1
fi

# Create RPM build directory structure
echo -e "${YELLOW}Creating RPM build directory...${NC}"
mkdir -p ${RPM_BUILD_DIR}/{BUILD,RPMS,SOURCES,SPECS,SRPMS}

# Build binary
echo -e "${YELLOW}Building binary...${NC}"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X 'main.Version=${VERSION}'" \
    -o ${RPM_BUILD_DIR}/BUILD/${APP_NAME} \
    ./cmd/agent

# Create spec file
echo -e "${YELLOW}Creating RPM spec file...${NC}"
cat > ${SPEC_FILE} <<EOF
Name:           ${APP_NAME}
Version:        ${VERSION}
Release:        ${RELEASE}%{?dist}
Summary:        Ducla Cloud Agent for distributed task execution

License:        MIT
URL:            https://github.com/your-org/ducla-cloud-agent
Source0:        %{name}-%{version}.tar.gz

BuildArch:      ${ARCH}
Requires:       systemd

%description
Ducla Cloud Agent is a high-performance agent for distributed task execution,
system monitoring, and cloud resource management.

%prep
# No prep needed for pre-built binary

%build
# No build needed for pre-built binary

%install
rm -rf %{buildroot}

# Install binary
mkdir -p %{buildroot}%{_bindir}
install -m 0755 %{_builddir}/${APP_NAME} %{buildroot}%{_bindir}/${APP_NAME}

# Install configuration
mkdir -p %{buildroot}%{_sysconfdir}/ducla
cat > %{buildroot}%{_sysconfdir}/ducla/agent.yaml <<'EOFCONFIG'
agent:
  id: "\${HOSTNAME}"
  name: "\${HOSTNAME}"
  environment: "production"

master:
  url: "\${DUCLA_MASTER_URL}"
  token: "\${DUCLA_AGENT_TOKEN}"

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
EOFCONFIG

# Install systemd service
mkdir -p %{buildroot}/usr/lib/systemd/system
cat > %{buildroot}/usr/lib/systemd/system/${APP_NAME}.service <<'EOFSERVICE'
[Unit]
Description=Ducla Cloud Agent
After=network.target
Documentation=https://github.com/your-org/ducla-cloud-agent

[Service]
Type=simple
User=ducla
Group=ducla
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
EOFSERVICE

# Create directories
mkdir -p %{buildroot}/opt/ducla/data
mkdir -p %{buildroot}/var/log/ducla

%pre
# Create user and group
getent group ducla >/dev/null || groupadd -r ducla
getent passwd ducla >/dev/null || \
    useradd -r -g ducla -d /opt/ducla -s /sbin/nologin \
    -c "Ducla Cloud Agent" ducla
exit 0

%post
%systemd_post ${APP_NAME}.service

%preun
%systemd_preun ${APP_NAME}.service

%postun
%systemd_postun_with_restart ${APP_NAME}.service

%files
/usr/bin/${APP_NAME}
%config(noreplace) /etc/ducla/agent.yaml
/usr/lib/systemd/system/${APP_NAME}.service
%dir %attr(0755,ducla,ducla) /opt/ducla
%dir %attr(0755,ducla,ducla) /opt/ducla/data
%dir %attr(0755,ducla,ducla) /var/log/ducla

%changelog
* $(date '+%a %b %d %Y') Builder <builder@ducla.cloud> - ${VERSION}-${RELEASE}
- Release ${VERSION}
EOF

# Build RPM
echo -e "${YELLOW}Building RPM package...${NC}"
rpmbuild --define "_topdir $(pwd)/${RPM_BUILD_DIR}" \
         --define "_builddir $(pwd)/${RPM_BUILD_DIR}/BUILD" \
         -bb ${SPEC_FILE}

# Move RPM to dist directory
mkdir -p dist
mv ${RPM_BUILD_DIR}/RPMS/${ARCH}/*.rpm dist/

echo ""
echo -e "${GREEN}RPM package built successfully!${NC}"
ls -lh dist/*.rpm
