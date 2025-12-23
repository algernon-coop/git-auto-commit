param(
    [Parameter(Mandatory=$true)]
    [string]$Arch,
    
    [Parameter(Mandatory=$true)]
    [string]$Version,
    
    [Parameter(Mandatory=$true)]
    [string]$BinaryPath
)

$ErrorActionPreference = "Stop"

# Clean version (remove 'v' prefix if present)
$CleanVersion = $Version -replace '^v', ''

$OutputDir = "dist"
$MsiFile = "$OutputDir\git-auto-commit-windows-$Arch.msi"

Write-Host "Building MSI for $Arch architecture..."
Write-Host "Version: $CleanVersion"
Write-Host "Binary: $BinaryPath"

# Check if binary exists
if (-not (Test-Path $BinaryPath)) {
    Write-Error "Binary not found at $BinaryPath"
    exit 1
}

# Build MSI using go-msi
go-msi make `
    --msi "$MsiFile" `
    --version "$CleanVersion" `
    --arch "$Arch" `
    --path "wix.json" `
    --src "build\windows\templates"

Write-Host "MSI created: $MsiFile"
