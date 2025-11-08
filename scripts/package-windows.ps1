# Windows installer packaging script using Inno Setup

param(
    [string]$Version = "1.0.0"
)

$ErrorActionPreference = "Stop"

# Configuration
$APP_NAME = "DuclaAgent"
$APP_DISPLAY_NAME = "Ducla Cloud Agent"
$PUBLISHER = "Ducla Cloud"
$APP_URL = "https://github.com/your-org/ducla-cloud-agent"
$BUILD_DIR = "build\windows"
$DIST_DIR = "dist"

Write-Host "Building Windows installer" -ForegroundColor Green
Write-Host "Version: $Version"
Write-Host ""

# Check for Inno Setup
$InnoSetupPath = "C:\Program Files (x86)\Inno Setup 6\ISCC.exe"
if (-not (Test-Path $InnoSetupPath)) {
    Write-Host "Inno Setup not found at: $InnoSetupPath" -ForegroundColor Red
    Write-Host "Download from: https://jrsoftware.org/isdl.php" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Alternatively, building portable ZIP package..." -ForegroundColor Yellow
    
    # Build portable package
    & .\scripts\build-windows-portable.ps1 -Version $Version
    exit 0
}

# Clean and create build directory
Write-Host "Creating build directory..." -ForegroundColor Yellow
if (Test-Path $BUILD_DIR) { Remove-Item -Recurse -Force $BUILD_DIR }
New-Item -ItemType Directory -Force -Path $BUILD_DIR | Out-Null
New-Item -ItemType Directory -Force -Path "$BUILD_DIR\bin" | Out-Null
New-Item -ItemType Directory -Force -Path "$BUILD_DIR\config" | Out-Null

# Build binary
Write-Host "Building binary..." -ForegroundColor Yellow
$env:GOOS = "windows"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"

$LDFLAGS = "-w -s -X 'main.Version=$Version'"
go build -ldflags="$LDFLAGS" -o "$BUILD_DIR\bin\$APP_NAME.exe" .\cmd\agent

if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed!" -ForegroundColor Red
    exit 1
}

# Copy configuration
Copy-Item "configs\agent.yaml" "$BUILD_DIR\config\agent.yaml"

# Create README
@"
Ducla Cloud Agent - Windows Installation
=========================================

Installation Directory: C:\Program Files\Ducla Cloud Agent\

Configuration:
- Config file: C:\ProgramData\Ducla\agent.yaml
- Data directory: C:\ProgramData\Ducla\data
- Log directory: C:\ProgramData\Ducla\logs

Service Management:
- The agent is installed as a Windows service
- Service name: DuclaAgent
- Start service: sc start DuclaAgent
- Stop service: sc stop DuclaAgent
- Service status: sc query DuclaAgent

Configuration:
1. Edit C:\ProgramData\Ducla\agent.yaml
2. Set environment variables (optional):
   - DUCLA_MASTER_URL
   - DUCLA_AGENT_TOKEN
   - DUCLA_JWT_SECRET
3. Restart the service

Uninstallation:
- Use "Add or Remove Programs" in Windows Settings
- Or run: "C:\Program Files\Ducla Cloud Agent\uninstall.exe"

Support:
- Documentation: https://github.com/your-org/ducla-cloud-agent
- Issues: https://github.com/your-org/ducla-cloud-agent/issues
"@ | Out-File -FilePath "$BUILD_DIR\README.txt" -Encoding UTF8

# Create Inno Setup script
Write-Host "Creating Inno Setup script..." -ForegroundColor Yellow
@"
#define MyAppName "$APP_DISPLAY_NAME"
#define MyAppVersion "$Version"
#define MyAppPublisher "$PUBLISHER"
#define MyAppURL "$APP_URL"
#define MyAppExeName "$APP_NAME.exe"
#define MyAppServiceName "$APP_NAME"

[Setup]
AppId={{A1B2C3D4-E5F6-7890-ABCD-EF1234567890}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
AppSupportURL={#MyAppURL}
AppUpdatesURL={#MyAppURL}
DefaultDirName={autopf}\Ducla Cloud Agent
DefaultGroupName=Ducla Cloud Agent
DisableProgramGroupPage=yes
LicenseFile=LICENSE
OutputDir=$DIST_DIR
OutputBaseFilename=ducla-agent-{#MyAppVersion}-windows-amd64-setup
Compression=lzma
SolidCompression=yes
WizardStyle=modern
PrivilegesRequired=admin
ArchitecturesAllowed=x64
ArchitecturesInstallIn64BitMode=x64

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Files]
Source: "$BUILD_DIR\bin\{#MyAppExeName}"; DestDir: "{app}"; Flags: ignoreversion
Source: "$BUILD_DIR\config\agent.yaml"; DestDir: "{commonappdata}\Ducla"; Flags: onlyifdoesntexist uninsneveruninstall
Source: "$BUILD_DIR\README.txt"; DestDir: "{app}"; Flags: ignoreversion
Source: "LICENSE"; DestDir: "{app}"; Flags: ignoreversion

[Dirs]
Name: "{commonappdata}\Ducla"
Name: "{commonappdata}\Ducla\data"
Name: "{commonappdata}\Ducla\logs"

[Icons]
Name: "{group}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"
Name: "{group}\Configuration"; Filename: "{commonappdata}\Ducla\agent.yaml"
Name: "{group}\Logs"; Filename: "{commonappdata}\Ducla\logs"
Name: "{group}\Uninstall {#MyAppName}"; Filename: "{uninstallexe}"

[Run]
Filename: "{app}\{#MyAppExeName}"; Parameters: "install"; StatusMsg: "Installing service..."; Flags: runhidden
Filename: "sc"; Parameters: "start {#MyAppServiceName}"; StatusMsg: "Starting service..."; Flags: runhidden

[UninstallRun]
Filename: "sc"; Parameters: "stop {#MyAppServiceName}"; Flags: runhidden
Filename: "{app}\{#MyAppExeName}"; Parameters: "uninstall"; Flags: runhidden

[Code]
function InitializeSetup(): Boolean;
begin
  Result := True;
  if not IsAdminLoggedOn then
  begin
    MsgBox('Administrator privileges are required to install this application.', mbError, MB_OK);
    Result := False;
  end;
end;

procedure CurStepChanged(CurStep: TSetupStep);
begin
  if CurStep = ssPostInstall then
  begin
    // Create environment variables or registry entries if needed
  end;
end;
"@ | Out-File -FilePath "$BUILD_DIR\installer.iss" -Encoding UTF8

# Build installer
Write-Host "Building installer..." -ForegroundColor Yellow
& $InnoSetupPath "$BUILD_DIR\installer.iss"

if ($LASTEXITCODE -eq 0) {
    Write-Host ""
    Write-Host "Windows installer built successfully!" -ForegroundColor Green
    Get-ChildItem "$DIST_DIR\ducla-agent-*-setup.exe" | Format-Table Name, Length, LastWriteTime
} else {
    Write-Host "Installer build failed!" -ForegroundColor Red
    exit 1
}
