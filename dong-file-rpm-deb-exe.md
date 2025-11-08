Packages được tạo:
1. Ubuntu/Debian (DEB)
Script: scripts/package-deb.sh
Output: ducla-agent_1.0.0_amd64.deb
Bao gồm: binary, systemd service, config, user/group setup
2. RedHat/CentOS/Fedora (RPM)
Script: scripts/package-rpm.sh
Output: ducla-agent-1.0.0-1.x86_64.rpm
Bao gồm: binary, systemd service, config, user/group setup
3. Windows Portable (ZIP)
Script: scripts/build-windows-portable.ps1
Output: ducla-agent-1.0.0-windows-amd64-portable.zip
Bao gồm: exe, config, batch files (start, install-service, uninstall-service)
4. Windows Installer (EXE)
Script: scripts/package-windows.ps1
Output: ducla-agent-1.0.0-windows-amd64-setup.exe
Yêu cầu: Inno Setup (tự động fallback về portable nếu không có)
Cách sử dụng:
# Ubuntu/Debian
./scripts/package-deb.sh
sudo dpkg -i dist/ducla-agent_1.0.0_amd64.deb

# RedHat/CentOS
./scripts/package-rpm.sh
sudo rpm -ivh dist/ducla-agent-1.0.0-1.x86_64.rpm

# Build tất cả Linux packages
./scripts/package-all.sh
# Windows Portable
.\scripts\build-windows-portable.ps1 -Version "1.0.0"

# Windows Installer
.\scripts\package-windows.ps1 -Version "1.0.0"

# Build tất cả Windows packages
.\scripts\package-all.ps1 -Version "1.0.0"
Tính năng chính:
✅ DEB/RPM: Systemd service, user/group tự động, pre/post install scripts ✅ Windows: Service installation, batch helpers, portable mode ✅ CI/CD: GitHub Actions workflow cho automated packaging ✅ Documentation: Chi tiết trong docs/PACKAGING.md ✅ Windows Service: Support install/uninstall service trong cmd/agent/service_windows.go

Tất cả packages đều production-ready với proper service management, security, và configuration handling!