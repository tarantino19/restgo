# Distribution Guide

This document outlines all the ways users can install and use REST API Summarizer.

## 🚀 Distribution Methods

### 1. GitHub Releases (Automated)

**How it works:**

- GitHub Actions automatically builds binaries for all platforms when you create a tag
- Users can download pre-built binaries without needing Go
- Supports Linux, macOS (Intel + Apple Silicon), and Windows

**To create a release:**

```bash
git tag v1.0.0
git push origin v1.0.0
```

This triggers the GitHub Action that:

1. Builds binaries for all platforms
2. Creates a GitHub release
3. Uploads all binaries and checksums

### 2. One-liner Installation Script

**User experience:**

```bash
curl -sSL https://raw.githubusercontent.com/tarantino19/restgo/main/scripts/install.sh | bash
```

**Features:**

- Auto-detects OS and architecture
- Downloads appropriate binary from GitHub releases
- Falls back to `go install` if binary download fails
- Handles installation permissions
- Provides helpful next steps

### 3. Go Install (For Go Users)

```bash
go install github.com/tarantino19/restgo@latest
```

**Pros:**

- Always gets latest version
- Compiles for exact user environment
- Familiar to Go developers

**Cons:**

- Requires Go installation
- Slower than downloading pre-built binary

### 4. Manual Binary Download

Users can manually download from:

- GitHub Releases page
- Direct URLs (e.g., `https://github.com/tarantino19/restgo/releases/latest/download/restapisummarizer-linux-amd64`)

### 5. Package Managers (Future)

**Homebrew (macOS/Linux):**

```bash
brew install tarantino19/tap/restapisummarizer
```

**Chocolatey (Windows):**

```bash
choco install restapisummarizer
```

**Snap (Linux):**

```bash
sudo snap install restapisummarizer
```

**APT/YUM packages:**

- Future consideration for major Linux distributions

## 📊 Distribution Analytics

### Current Status:

- ✅ GitHub Releases (automated)
- ✅ One-liner install script
- ✅ Go install
- ✅ Manual downloads
- 🔄 Package managers (in progress)

### Recommended Priority:

1. **GitHub Releases** - Primary distribution method
2. **Install Script** - Best user experience
3. **Go Install** - For Go developers
4. **Homebrew** - Popular on macOS
5. **Chocolatey** - Windows users
6. **Snap** - Modern Linux distributions

## 🛠️ Maintenance

### For Releases:

1. Update version in relevant files
2. Create and push git tag
3. GitHub Actions handles the rest
4. Update package manager formulas if needed

### For Install Script:

- Test on different OS/architectures
- Update if GitHub API changes
- Ensure fallback mechanisms work

### For Package Managers:

- Submit to official repositories
- Maintain formula/spec files
- Update when new versions are released

## 📈 Usage Statistics

Track installation methods through:

- GitHub release download counts
- Package manager statistics
- Install script analytics (if implemented)

## 🔧 Development Tools

Use the Makefile for local development:

```bash
make build        # Build for current platform
make build-all    # Build for all platforms
make install      # Install locally
make release      # Create release archives
```

## 🎯 User Onboarding

After installation, users need to:

1. Get a Gemini API key
2. Configure the key: `restapisummarizer config set api-key YOUR_KEY`
3. Test: `restapisummarizer sum`

The install script provides these instructions automatically.
