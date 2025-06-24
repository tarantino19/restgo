# REST API Summarizer Installation Script for Windows
# Usage: iwr -useb https://raw.githubusercontent.com/tarantino19/restgo/main/scripts/install.ps1 | iex

$ErrorActionPreference = "Stop"

$REPO = "tarantino19/restgo"
$BINARY_NAME = "restapisummarizer"

Write-Host "🚀 Installing REST API Summarizer..." -ForegroundColor Blue

# Detect architecture
$ARCH = $env:PROCESSOR_ARCHITECTURE
switch ($ARCH) {
    "AMD64" { $ARCH = "amd64" }
    "ARM64" { $ARCH = "arm64" }
    default { 
        Write-Host "❌ Unsupported architecture: $ARCH" -ForegroundColor Red
        exit 1 
    }
}

Write-Host "✅ Detected: windows/$ARCH" -ForegroundColor Green

# Get latest release
Write-Host "📦 Fetching latest release..." -ForegroundColor Blue
try {
    $LATEST_RELEASE = (Invoke-RestMethod -Uri "https://api.github.com/repos/$REPO/releases/latest").tag_name
} catch {
    Write-Host "❌ Could not fetch latest release. Trying Go install..." -ForegroundColor Red
    
    # Check if Go is installed
    if (!(Get-Command go -ErrorAction SilentlyContinue)) {
        Write-Host "❌ Go is not installed. Please install Go 1.24+ first." -ForegroundColor Red
        Write-Host "   Visit: https://golang.org/dl/" -ForegroundColor Yellow
        exit 1
    }
    
    Write-Host "📦 Installing via Go..." -ForegroundColor Blue
    go install github.com/$REPO@latest
    
    if (Get-Command $BINARY_NAME -ErrorAction SilentlyContinue) {
        Write-Host "✅ Installation successful via Go!" -ForegroundColor Green
    } else {
        Write-Host "❌ Installation failed. Make sure `$GOPATH/bin is in your PATH" -ForegroundColor Red
        exit 1
    }
    exit 0
}

# Download binary
$BINARY_URL = "https://github.com/$REPO/releases/download/$LATEST_RELEASE/$BINARY_NAME-windows-$ARCH.exe"
$TEMP_FILE = "$env:TEMP\$BINARY_NAME.exe"

Write-Host "📥 Downloading $BINARY_NAME $LATEST_RELEASE for windows/$ARCH..." -ForegroundColor Blue

try {
    Invoke-WebRequest -Uri $BINARY_URL -OutFile $TEMP_FILE
    
    # Determine install directory
    $INSTALL_DIR = ""
    if (Test-Path "C:\Program Files\") {
        $INSTALL_DIR = "C:\Program Files\$BINARY_NAME"
    } else {
        $INSTALL_DIR = "$env:USERPROFILE\bin"
    }
    
    # Create directory if it doesn't exist
    if (!(Test-Path $INSTALL_DIR)) {
        New-Item -ItemType Directory -Path $INSTALL_DIR -Force | Out-Null
    }
    
    Write-Host "📦 Installing to $INSTALL_DIR..." -ForegroundColor Blue
    
    # Copy binary
    Copy-Item $TEMP_FILE "$INSTALL_DIR\$BINARY_NAME.exe" -Force
    
    # Add to PATH if not already there
    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    if ($currentPath -notlike "*$INSTALL_DIR*") {
        Write-Host "⚠️  Adding $INSTALL_DIR to your PATH..." -ForegroundColor Yellow
        [Environment]::SetEnvironmentVariable("PATH", "$currentPath;$INSTALL_DIR", "User")
        Write-Host "   Please restart your terminal or run: refreshenv" -ForegroundColor Yellow
    }
    
    Write-Host "✅ Installation successful!" -ForegroundColor Green
    
} catch {
    Write-Host "❌ Failed to download binary. Trying Go install..." -ForegroundColor Red
    
    if (!(Get-Command go -ErrorAction SilentlyContinue)) {
        Write-Host "❌ Go is not installed. Please install Go 1.24+ first." -ForegroundColor Red
        Write-Host "   Visit: https://golang.org/dl/" -ForegroundColor Yellow
        exit 1
    }
    
    go install "github.com/$REPO@latest"
}

# Clean up temp file
if (Test-Path $TEMP_FILE) {
    Remove-Item $TEMP_FILE -Force
}

# Verify installation
if (Get-Command $BINARY_NAME -ErrorAction SilentlyContinue) {
    Write-Host ""
    Write-Host "🎉 REST API Summarizer is now installed!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Next steps:" -ForegroundColor Blue
    Write-Host "1. Get a Gemini API key: https://aistudio.google.com/app/apikey" -ForegroundColor Yellow
    Write-Host "2. Set your API key: $BINARY_NAME config set api-key YOUR_KEY" -ForegroundColor Yellow
    Write-Host "3. Start analyzing: $BINARY_NAME sum" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Run '$BINARY_NAME --help' for more information." -ForegroundColor Blue
    Write-Host ""
    try {
        $version = & $BINARY_NAME --version 2>$null
        Write-Host "Version: $version" -ForegroundColor Green
    } catch {
        Write-Host "Version: unknown" -ForegroundColor Green
    }
} else {
    Write-Host "❌ Installation verification failed" -ForegroundColor Red
    Write-Host "   You may need to restart your terminal or run: refreshenv" -ForegroundColor Yellow
    exit 1
} 