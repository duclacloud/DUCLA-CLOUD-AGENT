# Build portable Windows package (ZIP)

param(
    [string]$Version = "1.0.0"
)

$ErrorActionPreference = "Stop"

$APP_NAME = "ducla-agent"
$BUILD_DIR = "build\windows-portable"
$DIST_DIR = "dist"

Write-Host "Building Windows portable package" -ForegroundColor Green
Write-Host "Version: $Version"
Write-Host ""

# Clean and create build directory
Write-Host "Creating build directory..." -ForegroundColor Yellow
if (Test-Path $BUILD_DIR) { Remove-Item -Recurse -Force $BUILD_DIR }
New-Item -ItemType Directory -Force -Path "$BUILD_DIR\$APP_NAME" | Out-Null

# Build binary
Write-Host "Building binary..." -ForegroundColor Yellow
$env:GOOS = "windows"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"

$LDFLAGS = "-w -s -X 'main.Version=$Version'"
go build -ldflags="$LDFLAGS" -o "$BUILD_DIR\$APP_NAME\$APP_NAME.exe" .\cmd\agent

if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed!" -ForegroundColor Red
    exit 1
}

# Copy configuration
New-Item -ItemType Directory -Force -Path "$BUILD_DIR\$APP_NAME\config" | Out-Null
Copy-Item "configs\agent.yaml" "$BUILD_DIR\$APP_NAME\config\agent.yaml"

# Create batch files
@"
@echo off
echo Starting Ducla Cloud Agent...
"%~dp0$APP_NAME.exe" --config "%~dp0config\agent.yaml"
pause
"@ | Out-File -FilePath "$BUILD_DIR\$APP_NAME\start.bat" -Encoding ASCII

@"
@echo off
echo Installing Ducla Cloud Agent as Windows Service...
"%~dp0$APP_NAME.exe" install
echo.
echo Service installed successfully!
echo Use 'sc start DuclaAgent' to start the service
pause
"@ | Out-File -FilePath "$BUILD_DIR\$APP_NAME\install-service.bat" -Encoding ASCII

@"
@echo off
echo Uninstalling Ducla Cloud Agent service...
"%~dp0$APP_NAME.exe" uninstall
echo.
echo Service uninstalled successfully!
pause
"@ | Out-File -FilePath "$BUILD_DIR\$APP_NAME\uninstall-service.bat" -Encoding ASCII

# Create README
@"
Ducla Cloud Agent - Windows Portable
=====================================

Version: $Version

Quick Start:
1. Edit config\agent.yaml with your settings
2. Double-click start.bat to run the agent
3. Or install as Windows service using install-service.bat

Files:
- $APP_NAME.exe          - Main executable
- config\agent.yaml      - Configuration file
- start.bat              - Start agent in console mode
- install-service.bat    - Install as Windows service
- uninstall-service.bat  - Uninstall Windows service

Configuration:
Edit config\agent.yaml and set:
- master.url: Your master server URL
- master.token: Your authentication token

Or set environment variables:
- DUCLA_MASTER_URL
- DUCLA_AGENT_TOKEN
- DUCLA_JWT_SECRET

Running as Service:
1. Run install-service.bat as Administrator
2. Start service: sc start DuclaAgent
3. Stop service: sc stop DuclaAgent
4. Check status: sc query DuclaAgent

Data Directories:
- Data: %PROGRAMDATA%\Ducla\data
- Logs: %PROGRAMDATA%\Ducla\logs

Support:
- Documentation: https://github.com/your-org/ducla-cloud-agent
- Issues: https://github.com/your-org/ducla-cloud-agent/issues
"@ | Out-File -FilePath "$BUILD_DIR\$APP_NAME\README.txt" -Encoding UTF8

# Copy license
if (Test-Path "LICENSE") {
    Copy-Item "LICENSE" "$BUILD_DIR\$APP_NAME\LICENSE.txt"
}

# Create ZIP archive
Write-Host "Creating ZIP archive..." -ForegroundColor Yellow
$ZipFile = "$DIST_DIR\$APP_NAME-$Version-windows-amd64-portable.zip"
New-Item -ItemType Directory -Force -Path $DIST_DIR | Out-Null

if (Test-Path $ZipFile) { Remove-Item $ZipFile }
Compress-Archive -Path "$BUILD_DIR\$APP_NAME\*" -DestinationPath $ZipFile

Write-Host ""
Write-Host "Windows portable package built successfully!" -ForegroundColor Green
Get-ChildItem $ZipFile | Format-Table Name, Length, LastWriteTime

Write-Host ""
Write-Host "To use:" -ForegroundColor Yellow
Write-Host "1. Extract the ZIP file"
Write-Host "2. Edit config\agent.yaml"
Write-Host "3. Run start.bat or install-service.bat"
