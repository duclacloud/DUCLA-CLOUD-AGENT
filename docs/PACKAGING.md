# Packaging Guide

This guide explains how to build distribution packages for different platforms.

## Prerequisites

### For DEB packages (Ubuntu/Debian)
```bash
sudo apt-get install dpkg-dev
```

### For RPM packages (RedHat/CentOS/Fedora)
```bash
# RHEL/CentOS
sudo yum install rpm-build

# Fedora
sudo dnf install rpm-build
```

### For Windows packages
- **Portable ZIP**: No additional tools needed
- **Installer**: [Inno Setup 6](https://jrsoftware.org/isdl.php) (optional)

## Building Packages

### Ubuntu/Debian (DEB)

```bash
# Build DEB package
./scripts/package-deb.sh

# Or with specific version
VERSION=1.0.0 ./scripts/package-deb.sh

# Install
sudo dpkg -i dist/ducla-agent_1.0.0_amd64.deb
sudo apt-get install -f  # Install dependencies if needed

# Uninstall
sudo apt-get remove ducla-agent

# Purge (remove config files)
sudo apt-get purge ducla-agent
```

### RedHat/CentOS/Fedora (RPM)

```bash
# Build RPM package
./scripts/package-rpm.sh

# Or with specific version
VERSION=1.0.0 ./scripts/package-rpm.sh

# Install
sudo rpm -ivh dist/ducla-agent-1.0.0-1.x86_64.rpm

# Or with yum/dnf
sudo yum install dist/ducla-agent-1.0.0-1.x86_64.rpm

# Uninstall
sudo rpm -e ducla-agent
```

### Windows

#### Portable Package (ZIP)

```powershell
# Build portable package
.\scripts\build-windows-portable.ps1 -Version "1.0.0"

# Extract and use
Expand-Archive dist\ducla-agent-1.0.0-windows-amd64-portable.zip -DestinationPath C:\ducla-agent
cd C:\ducla-agent
.\start.bat
```

#### Windows Installer (EXE)

```powershell
# Build installer (requires Inno Setup)
.\scripts\package-windows.ps1 -Version "1.0.0"

# Run installer
.\dist\ducla-agent-1.0.0-windows-amd64-setup.exe
```

### Build All Packages

```bash
# Linux (DEB + RPM)
VERSION=1.0.0 ./scripts/package-all.sh

# Windows (Portable + Installer)
.\scripts\package-all.ps1 -Version "1.0.0"
```

## Package Contents

### DEB/RPM Packages Include:
- Binary: `/usr/bin/ducla-agent`
- Configuration: `/etc/ducla/agent.yaml`
- Systemd service: `/lib/systemd/system/ducla-agent.service`
- Data directory: `/opt/ducla/data`
- Log directory: `/var/log/ducla`
- Service user: `ducla`

### Windows Portable Package Includes:
- Binary: `ducla-agent.exe`
- Configuration: `config/agent.yaml`
- Batch files: `start.bat`, `install-service.bat`, `uninstall-service.bat`
- Documentation: `README.txt`

### Windows Installer Includes:
- All portable package contents
- Automatic service installation
- Start menu shortcuts
- Uninstaller

## Post-Installation

### Linux (DEB/RPM)

1. Edit configuration:
```bash
sudo nano /etc/ducla/agent.yaml
```

2. Set environment variables (optional):
```bash
sudo nano /etc/default/ducla-agent
```

3. Start service:
```bash
sudo systemctl start ducla-agent
sudo systemctl status ducla-agent
```

4. Enable auto-start:
```bash
sudo systemctl enable ducla-agent
```

### Windows

#### Portable Package

1. Edit `config\agent.yaml`
2. Run `start.bat` for console mode
3. Or run `install-service.bat` as Administrator for service mode

#### Installer

1. Edit `C:\ProgramData\Ducla\agent.yaml`
2. Service is automatically installed and started
3. Manage via Services app or `sc` command

## Service Management

### Linux (systemd)

```bash
# Start
sudo systemctl start ducla-agent

# Stop
sudo systemctl stop ducla-agent

# Restart
sudo systemctl restart ducla-agent

# Status
sudo systemctl status ducla-agent

# Logs
sudo journalctl -u ducla-agent -f
```

### Windows

```powershell
# Start
sc start DuclaAgent

# Stop
sc stop DuclaAgent

# Query status
sc query DuclaAgent

# View logs
Get-Content C:\ProgramData\Ducla\logs\agent.log -Tail 50 -Wait
```

## Troubleshooting

### DEB Package Issues

```bash
# Check package info
dpkg -l | grep ducla-agent

# List package files
dpkg -L ducla-agent

# Verify package
dpkg -V ducla-agent

# Reconfigure
sudo dpkg-reconfigure ducla-agent
```

### RPM Package Issues

```bash
# Check package info
rpm -qi ducla-agent

# List package files
rpm -ql ducla-agent

# Verify package
rpm -V ducla-agent
```

### Windows Service Issues

```powershell
# Check service status
Get-Service DuclaAgent

# View service details
sc qc DuclaAgent

# Check event logs
Get-EventLog -LogName Application -Source DuclaAgent -Newest 10
```

## Building for Distribution

### Create Release Packages

```bash
# Set version
export VERSION=1.0.0

# Build all Linux packages
./scripts/package-all.sh

# Generate checksums
cd dist
sha256sum *.deb *.rpm > checksums.txt
```

### Signing Packages

#### DEB Packages
```bash
dpkg-sig --sign builder ducla-agent_1.0.0_amd64.deb
```

#### RPM Packages
```bash
rpm --addsign ducla-agent-1.0.0-1.x86_64.rpm
```

## Repository Setup

### APT Repository (Debian/Ubuntu)

```bash
# Create repository structure
mkdir -p repo/deb/pool/main
cp dist/*.deb repo/deb/pool/main/

# Generate Packages file
cd repo/deb
dpkg-scanpackages pool/main /dev/null | gzip -9c > pool/main/Packages.gz
```

### YUM Repository (RedHat/CentOS)

```bash
# Create repository structure
mkdir -p repo/rpm
cp dist/*.rpm repo/rpm/

# Create repository metadata
createrepo repo/rpm
```

## Continuous Integration

See `.github/workflows/release.yml` for automated package building on release.

## Support

For issues with packaging:
- Check logs in `/var/log/ducla/` (Linux) or `C:\ProgramData\Ducla\logs\` (Windows)
- Review systemd journal: `journalctl -u ducla-agent`
- Open an issue: https://github.com/your-org/ducla-cloud-agent/issues
