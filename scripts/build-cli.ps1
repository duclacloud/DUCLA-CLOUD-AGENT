# Build script for Ducla Cloud Agent CLI (PowerShell)
# Supports multiple platforms and architectures

param(
    [string]$Mode = "all"
)

$ErrorActionPreference = "Stop"

# Configuration
$APP_NAME = "ducla-agent"
$VERSION = if ($env:VERSION) { $env:VERSION } else { 
    try { git describe --tags --always --dirty 2>$null } catch { "dev" }
}
$BUILD_TIME = (Get-Date).ToUniversalTime().ToString("yyyy-MM-dd_HH:mm:ss")
$GIT_COMMIT = try { git rev-parse --short HEAD 2>$null } catch { "unknown" }
$GO_VERSION = (go version).Split()[2]

# Build directories
$BUILD_DIR = "bin"
$DIST_DIR = "dist"

# Platforms to build for
$PLATFORMS = @(
    @{OS="linux"; ARCH="amd64"},
    @{OS="linux"; ARCH="arm64"},
    @{OS="darwin"; ARCH="amd64"},
    @{OS="darwin"; ARCH="arm64"},
    @{OS="windows"; ARCH="amd64"}
)

# Ldflags for version information
$LDFLAGS = "-w -s " +
    "-X 'main.Version=$VERSION' " +
    "-X 'main.BuildTime=$BUILD_TIME' " +
    "-X 'main.GitCommit=$GIT_COMMIT' " +
    "-X 'main.GoVersion=$GO_VERSION'"

Write-Host "Building Ducla Cloud Agent CLI" -ForegroundColor Green
Write-Host "Version: $VERSION"
Write-Host "Build Time: $BUILD_TIME"
Write-Host "Git Commit: $GIT_COMMIT"
Write-Host ""

# Clean previous builds
Write-Host "Cleaning previous builds..." -ForegroundColor Yellow
if (Test-Path $BUILD_DIR) { Remove-Item -Recurse -Force $BUILD_DIR }
if (Test-Path $DIST_DIR) { Remove-Item -Recurse -Force $DIST_DIR }
New-Item -ItemType Directory -Force -Path $BUILD_DIR | Out-Null
New-Item -ItemType Directory -Force -Path $DIST_DIR | Out-Null

# Quick build for current platform
if ($Mode -eq "quick") {
    Write-Host "Building for current platform..." -ForegroundColor Green
    go build -ldflags="$LDFLAGS" -o "$BUILD_DIR\$APP_NAME.exe" .\cmd\agent
    Write-Host "✓ Build complete: $BUILD_DIR\$APP_NAME.exe" -ForegroundColor Green
    exit 0
}

# Build for all platforms
foreach ($platform in $PLATFORMS) {
    $GOOS = $platform.OS
    $GOARCH = $platform.ARCH
    
    $output_name = "$APP_NAME-$VERSION-$GOOS-$GOARCH"
    
    if ($GOOS -eq "windows") {
        $output_name = "$output_name.exe"
    }
    
    Write-Host "Building for $GOOS/$GOARCH..." -ForegroundColor Yellow
    
    $env:GOOS = $GOOS
    $env:GOARCH = $GOARCH
    $env:CGO_ENABLED = "0"
    
    go build -ldflags="$LDFLAGS" -o "$DIST_DIR\$output_name" .\cmd\agent
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ Built: $output_name" -ForegroundColor Green
        
        # Create archive
        Push-Location $DIST_DIR
        if ($GOOS -eq "windows") {
            Compress-Archive -Path $output_name -DestinationPath "$($output_name -replace '\.exe$','').zip" -Force
        } else {
            tar -czf "$output_name.tar.gz" $output_name
        }
        Pop-Location
    } else {
        Write-Host "✗ Failed to build for $GOOS/$GOARCH" -ForegroundColor Red
    }
}

# Generate checksums
Write-Host "Generating checksums..." -ForegroundColor Yellow
Push-Location $DIST_DIR
Get-ChildItem -Filter *.tar.gz,*.zip | Get-FileHash -Algorithm SHA256 | 
    Select-Object @{Name="Hash";Expression={$_.Hash}}, @{Name="File";Expression={$_.Path | Split-Path -Leaf}} |
    ForEach-Object { "$($_.Hash)  $($_.File)" } | Out-File -FilePath checksums.txt -Encoding ASCII
Pop-Location

Write-Host ""
Write-Host "Build complete!" -ForegroundColor Green
Write-Host "Binaries are in: $DIST_DIR\"
Get-ChildItem $DIST_DIR | Format-Table Name, Length, LastWriteTime
