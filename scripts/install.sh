#!/bin/bash

# REST API Summarizer Installation Script
# Usage: curl -sSL https://raw.githubusercontent.com/tarantino19/restgo/main/scripts/install.sh | bash

set -e

REPO="tarantino19/restgo"
BINARY_NAME="restapisummarizer"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Installing REST API Summarizer...${NC}"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) echo -e "${RED}❌ Unsupported architecture: $ARCH${NC}"; exit 1 ;;
esac

case $OS in
    linux) OS="linux" ;;
    darwin) OS="darwin" ;;
    *) echo -e "${RED}❌ Unsupported OS: $OS${NC}"; exit 1 ;;
esac

echo -e "${GREEN}✅ Detected: $OS/$ARCH${NC}"

# Get latest release
echo -e "${BLUE}📦 Fetching latest release...${NC}"
LATEST_RELEASE=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_RELEASE" ]; then
    echo -e "${RED}❌ Could not fetch latest release. Falling back to Go install...${NC}"
    
    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        echo -e "${RED}❌ Go is not installed. Please install Go 1.24+ first.${NC}"
        echo -e "${YELLOW}   Visit: https://golang.org/dl/${NC}"
        exit 1
    fi
    
    echo -e "${BLUE}📦 Installing via Go...${NC}"
    go install github.com/$REPO@latest
    
    if command -v $BINARY_NAME &> /dev/null; then
        echo -e "${GREEN}✅ Installation successful via Go!${NC}"
    else
        echo -e "${RED}❌ Installation failed. Make sure \$GOPATH/bin is in your \$PATH${NC}"
        exit 1
    fi
else
    # Download binary
    BINARY_URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/${BINARY_NAME}-${OS}-${ARCH}"
    TEMP_FILE="/tmp/$BINARY_NAME"
    
    echo -e "${BLUE}📥 Downloading $BINARY_NAME $LATEST_RELEASE for $OS/$ARCH...${NC}"
    
    if curl -L -o "$TEMP_FILE" "$BINARY_URL"; then
        chmod +x "$TEMP_FILE"
        
        # Install to /usr/local/bin or ~/bin
        if [ -w "/usr/local/bin" ]; then
            INSTALL_DIR="/usr/local/bin"
        elif [ -d "$HOME/bin" ]; then
            INSTALL_DIR="$HOME/bin"
        else
            mkdir -p "$HOME/bin"
            INSTALL_DIR="$HOME/bin"
        fi
        
        echo -e "${BLUE}📦 Installing to $INSTALL_DIR...${NC}"
        
        if [ "$INSTALL_DIR" = "/usr/local/bin" ]; then
            if sudo mv "$TEMP_FILE" "$INSTALL_DIR/$BINARY_NAME"; then
                echo -e "${GREEN}✅ Installation successful!${NC}"
            else
                echo -e "${RED}❌ Failed to install to $INSTALL_DIR${NC}"
                exit 1
            fi
        else
            mv "$TEMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
            echo -e "${GREEN}✅ Installation successful!${NC}"
            
            # Check if ~/bin is in PATH
            if [[ ":$PATH:" != *":$HOME/bin:"* ]]; then
                echo -e "${YELLOW}⚠️  Add $HOME/bin to your PATH:${NC}"
                echo -e "${YELLOW}   export PATH=\$PATH:$HOME/bin${NC}"
            fi
        fi
    else
        echo -e "${RED}❌ Failed to download binary. Trying Go install...${NC}"
        
        if ! command -v go &> /dev/null; then
            echo -e "${RED}❌ Go is not installed. Please install Go 1.24+ first.${NC}"
            echo -e "${YELLOW}   Visit: https://golang.org/dl/${NC}"
            exit 1
        fi
        
        go install github.com/$REPO@latest
    fi
fi

# Verify installation
if command -v $BINARY_NAME &> /dev/null; then
    echo ""
    echo -e "${GREEN}🎉 REST API Summarizer is now installed!${NC}"
    echo ""
    echo -e "${BLUE}Next steps:${NC}"
    echo -e "${YELLOW}1. Get a Gemini API key: https://aistudio.google.com/app/apikey${NC}"
    echo -e "${YELLOW}2. Set your API key: $BINARY_NAME config set api-key YOUR_KEY${NC}"
    echo -e "${YELLOW}3. Start analyzing: $BINARY_NAME sum${NC}"
    echo ""
    echo -e "${BLUE}Run '$BINARY_NAME --help' for more information.${NC}"
    echo ""
    echo -e "${GREEN}Version: $($BINARY_NAME --version 2>/dev/null || echo 'unknown')${NC}"
else
    echo -e "${RED}❌ Installation verification failed${NC}"
    exit 1
fi 