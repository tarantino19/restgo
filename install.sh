#!/bin/bash

# REST API Summarizer Installation Script
# Usage: curl -sSL https://raw.githubusercontent.com/tarantino19/restgo/main/install.sh | bash

set -e

echo "🚀 Installing REST API Summarizer..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.24+ first."
    echo "   Visit: https://golang.org/dl/"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | cut -d' ' -f3 | sed 's/go//')
REQUIRED_VERSION="1.24"

if ! printf '%s\n%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V -C; then
    echo "❌ Go version $GO_VERSION is too old. Please upgrade to Go $REQUIRED_VERSION or later."
    exit 1
fi

echo "✅ Go version $GO_VERSION detected"

# Install the tool
echo "📦 Installing restapisummarizer..."
go install github.com/tarantino19/restgo@latest

# Check if installation was successful
if command -v restapisummarizer &> /dev/null; then
    echo "✅ Installation successful!"
    echo ""
    echo "🎉 REST API Summarizer is now installed!"
    echo ""
    echo "Next steps:"
    echo "1. Get a Gemini API key: https://aistudio.google.com/app/apikey"
    echo "2. Set your API key: restapisummarizer config set api-key YOUR_KEY"
    echo "3. Start analyzing: restapisummarizer sum"
    echo ""
    echo "Run 'restapisummarizer --help' for more information."
else
    echo "❌ Installation failed. Make sure \$GOPATH/bin is in your \$PATH"
    echo "   Add this to your shell profile:"
    echo "   export PATH=\$PATH:\$(go env GOPATH)/bin"
fi 