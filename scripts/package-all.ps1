# Build all Windows packages

param(
    [string]$Version = "1.0.0"
)

$ErrorActionPreference = "Stop"

Write-Host "Building all Windows packages" -ForegroundColor Green
Write-Host "Version: $Version"
Write-Host ""

# Build portable package
Write-Host "Building portable package..." -ForegroundColor Yellow
& .\scripts\build-windows-portable.ps1 -Version $Version

# Build installer if Inno Setup is available
Write-Host ""
Write-Host "Building installer..." -ForegroundColor Yellow
& .\scripts\package-windows.ps1 -Version $Version

Write-Host ""
Write-Host "All Windows packages built successfully!" -ForegroundColor Green
Write-Host ""
Write-Host "Packages:"
Get-ChildItem dist\ | Format-Table Name, Length, LastWriteTime
