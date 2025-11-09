Name:           ducla-agent
Version:        1.0.0_1_g23271c6_dirty
Release:        1%{?dist}
Summary:        Ducla Cloud Agent for distributed task execution

License:        MIT
URL:            https://github.com/your-org/ducla-cloud-agent
Source0:        %{name}-%{version}.tar.gz

BuildArch:      x86_64
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
install -m 0755 %{_builddir}/ducla-agent %{buildroot}%{_bindir}/ducla-agent

# Install configuration
mkdir -p %{buildroot}%{_sysconfdir}/ducla
cat > %{buildroot}%{_sysconfdir}/ducla/agent.yaml <<'EOFCONFIG'
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
EOFCONFIG

# Install systemd service
mkdir -p %{buildroot}/usr/lib/systemd/system
cat > %{buildroot}/usr/lib/systemd/system/ducla-agent.service <<'EOFSERVICE'
[Unit]
Description=Ducla Cloud Agent
After=network.target
Documentation=https://github.com/your-org/ducla-cloud-agent

[Service]
Type=simple
User=ducla
Group=ducla
ExecStart=/usr/bin/ducla-agent --config /etc/ducla/agent.yaml
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
getent passwd ducla >/dev/null ||     useradd -r -g ducla -d /opt/ducla -s /sbin/nologin     -c "Ducla Cloud Agent" ducla
exit 0

%post
%systemd_post ducla-agent.service

%preun
%systemd_preun ducla-agent.service

%postun
%systemd_postun_with_restart ducla-agent.service

%files
/usr/bin/ducla-agent
%config(noreplace) /etc/ducla/agent.yaml
/usr/lib/systemd/system/ducla-agent.service
%dir %attr(0755,ducla,ducla) /opt/ducla
%dir %attr(0755,ducla,ducla) /opt/ducla/data
%dir %attr(0755,ducla,ducla) /var/log/ducla

%changelog
* Sun Nov 09 2025 Builder <builder@ducla.cloud> - 1.0.0_1_g23271c6_dirty-1
- Release 1.0.0_1_g23271c6_dirty
