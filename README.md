# REST API Summarizer

A powerful CLI tool that automatically discovers REST API endpoints in your codebase and generates AI-powered summaries using Google's Gemini API.

## Features

- 🔍 **Recursive Directory Scanning**: Automatically scans all subdirectories in your project
- 🤖 **AI-Powered Summaries**: Uses Google's Gemini API to generate concise endpoint descriptions
- 🌐 **Multi-Language Support**: Works with multiple programming languages and frameworks
- 📊 **Beautiful Output**: Displays results in a clean, formatted table
- ⚡ **Fast & Efficient**: Optimized scanning with intelligent directory filtering

## Supported Frameworks

- **JavaScript/TypeScript**: Express.js, Fastify, Koa
- **Python**: Flask, FastAPI, Django
- **Go**: Gin, Echo, Gorilla Mux
- **Java**: Spring Boot, JAX-RS
- **Ruby**: Ruby on Rails
- **C#**: ASP.NET Core
- **PHP**: Laravel, Symfony (coming soon)

## Installation

### 🚀 Quick Install (Recommended)

**One-liner installation** - automatically detects your OS and downloads the right binary:

```bash
curl -sSL https://raw.githubusercontent.com/tarantino19/restgo/main/scripts/install.sh | bash
```

### 📦 Other Installation Methods

<details>
<summary><strong>Pre-built Binaries</strong> (No Go required)</summary>

Download from [GitHub Releases](https://github.com/tarantino19/restgo/releases):

**Linux (x64):**

```bash
curl -L -o restapisummarizer https://github.com/tarantino19/restgo/releases/latest/download/restapisummarizer-linux-amd64
chmod +x restapisummarizer
sudo mv restapisummarizer /usr/local/bin/
```

**macOS (Intel):**

```bash
curl -L -o restapisummarizer https://github.com/tarantino19/restgo/releases/latest/download/restapisummarizer-darwin-amd64
chmod +x restapisummarizer
sudo mv restapisummarizer /usr/local/bin/
```

**macOS (Apple Silicon):**

```bash
curl -L -o restapisummarizer https://github.com/tarantino19/restgo/releases/latest/download/restapisummarizer-darwin-arm64
chmod +x restapisummarizer
sudo mv restapisummarizer /usr/local/bin/
```

**Windows:**
Download `restapisummarizer-windows-amd64.exe` from the releases page.

</details>

<details>
<summary><strong>Go Install</strong> (Requires Go 1.24+)</summary>

```bash
go install github.com/tarantino19/restgo@latest
```

Make sure `$GOPATH/bin` is in your `$PATH`.

</details>

<details>
<summary><strong>Build from Source</strong></summary>

```bash
git clone https://github.com/tarantino19/restgo.git
cd restgo
go build -o restapisummarizer
sudo mv restapisummarizer /usr/local/bin/
```

</details>

<details>
<summary><strong>Package Managers</strong></summary>

**Homebrew (macOS/Linux):**

```bash
# Coming soon
brew install tarantino19/tap/restapisummarizer
```

**Chocolatey (Windows):**

```bash
# Coming soon
choco install restapisummarizer
```

**Snap (Linux):**

```bash
# Coming soon
sudo snap install restapisummarizer
```

</details>

### Requirements

- Google Gemini API key ([Get one free](https://aistudio.google.com/app/apikey))
- No other dependencies required for pre-built binaries!

## Quick Start

### 1. Get a Gemini API Key

1. Visit [Google AI Studio](https://makersuite.google.com/app/apikey)
2. Create a new API key
3. Copy the key for the next step

### 2. Set Your API Key

```bash
restapisummarizer config set api-key YOUR_API_KEY_HERE
```

### 3. Analyze Your Project

```bash
# Analyze current directory and all subdirectories
restapisummarizer sum

# Analyze a specific directory
restapisummarizer sum /path/to/your/project
```

## Command Reference

### Main Commands

| Command   | Description                                    | Usage                                   |
| --------- | ---------------------------------------------- | --------------------------------------- |
| `sum`     | Analyze REST API endpoints in a directory tree | `restapisummarizer sum [directory]`     |
| `config`  | Manage configuration settings                  | `restapisummarizer config <subcommand>` |
| `help`    | Show help information                          | `restapisummarizer help`                |
| `version` | Show version information                       | `restapisummarizer version`             |

### Config Subcommands

| Command              | Description                   | Usage                                           |
| -------------------- | ----------------------------- | ----------------------------------------------- |
| `config set api-key` | Set the Gemini API key        | `restapisummarizer config set api-key YOUR_KEY` |
| `config get api-key` | View current API key (masked) | `restapisummarizer config get api-key`          |

## Usage Examples

### Basic Usage

```bash
# Analyze current directory (recursively scans all subdirectories)
restapisummarizer sum

# Analyze a specific project
restapisummarizer sum ~/projects/my-api

# Set API key
restapisummarizer config set api-key AIzaSyC...your-key-here

# Check current API key
restapisummarizer config get api-key
# Output: Current API key: AIza...here
```

### Example Output

```
🔍 Scanning directory tree: /Users/example/my-api
This will recursively scan all subdirectories...

  📂 Entering: /Users/example/my-api/src
  📂 Entering: /Users/example/my-api/src/controllers
  ✓ Found 5 endpoints in src/controllers/users.js
  ✓ Found 3 endpoints in src/controllers/auth.js
  📂 Entering: /Users/example/my-api/src/routes
  ✓ Found 2 endpoints in src/routes/products.js

✓ Scan complete! Analyzed 12 files, found 10 endpoints

🤖 Initializing Gemini AI...

📝 Generating endpoint summaries...

🔍 REST API Endpoints Summary
Found 10 endpoints

┌────────┬─────────────────────┬──────────────────────────┬────────────────────────────────────────────┐
│ Method │ Path                │ File                     │ Summary                                    │
├────────┼─────────────────────┼──────────────────────────┼────────────────────────────────────────────┤
│ GET    │ /api/users          │ src/controllers/users.js:15 │ Retrieves a list of all users           │
│ POST   │ /api/users          │ src/controllers/users.js:28 │ Creates a new user account              │
│ GET    │ /api/users/:id      │ src/controllers/users.js:45 │ Fetches a specific user by ID           │
│ PUT    │ /api/users/:id      │ src/controllers/users.js:62 │ Updates user information                │
│ DELETE │ /api/users/:id      │ src/controllers/users.js:79 │ Deletes a user account                  │
│ POST   │ /api/auth/login     │ src/controllers/auth.js:12  │ Authenticates user credentials          │
│ POST   │ /api/auth/logout    │ src/controllers/auth.js:34  │ Ends user session                       │
│ POST   │ /api/auth/refresh   │ src/controllers/auth.js:56  │ Refreshes authentication token          │
│ GET    │ /api/products       │ src/routes/products.js:8    │ Lists all available products            │
│ GET    │ /api/products/:id   │ src/routes/products.js:22   │ Gets product details by ID              │
└────────┴─────────────────────┴──────────────────────────┴────────────────────────────────────────────┘

📁 Endpoints by File:
  src/controllers/users.js (5 endpoints)
    • GET /api/users
    • POST /api/users
    • GET /api/users/:id
    • PUT /api/users/:id
    • DELETE /api/users/:id

  src/controllers/auth.js (3 endpoints)
    • POST /api/auth/login
    • POST /api/auth/logout
    • POST /api/auth/refresh

  src/routes/products.js (2 endpoints)
    • GET /api/products
    • GET /api/products/:id

✅ Analysis completed in 23s
```

## Configuration

### Config File Location

- **macOS/Linux**: `~/.restapisummarizer/config.yaml`
- **Windows**: `%USERPROFILE%\.restapisummarizer\config.yaml`

### Environment Variables

You can also set the API key using environment variables:

```bash
# Option 1
export GEMINI_API_KEY=YOUR_API_KEY_HERE

# Option 2 (with prefix)
export RESTAPI_GEMINI_API_KEY=YOUR_API_KEY_HERE
```

## How It Works

1. **Recursive Scanning**: The tool starts at the specified directory and recursively walks through all subdirectories
2. **File Detection**: It identifies source code files based on their extensions (.js, .ts, .py, .go, etc.)
3. **Pattern Matching**: Uses regex patterns specific to each framework to find REST API endpoint definitions
4. **Context Extraction**: Captures surrounding code for better AI analysis
5. **AI Summary**: Sends the code context to Gemini API to generate human-readable summaries
6. **Display Results**: Formats everything in a beautiful table with color-coded HTTP methods

## Ignored Directories

The tool automatically skips these directories to improve performance:

- `node_modules`, `vendor`, `.git`, `dist`, `build`, `target`
- `__pycache__`, `.venv`, `venv`, `env`
- `.idea`, `.vscode`, `coverage`
- `test`, `tests`, `spec`, `specs`
- `.next`, `out`, `tmp`, `temp`, `cache`, `.cache`, `logs`

## Tips & Best Practices

1. **Large Codebases**: The tool handles large projects efficiently by skipping non-source directories
2. **API Key Security**: Your API key is stored locally and never transmitted except to Google's API
3. **Rate Limiting**: Free tier allows 15 requests/minute - the tool automatically handles this
4. **Accuracy**: Follow standard REST API patterns in your framework for best results

## Troubleshooting

### "No endpoints found"

- Ensure you're in a directory with source code
- Check that your endpoints follow standard patterns for your framework
- Try running with `-v` flag for verbose output (coming soon)

### "API key not set"

```bash
# Set your API key
restapisummarizer config set api-key YOUR_KEY

# Or use environment variable
export GEMINI_API_KEY=YOUR_KEY
```

### "Rate limit exceeded"

- The tool automatically handles rate limiting
- For very large codebases, the analysis may take longer
- Consider upgrading your Gemini API plan for higher limits

### "Permission denied"

```bash
# Make the binary executable
chmod +x restapisummarizer

# Or install globally
sudo mv restapisummarizer /usr/local/bin/
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Adding New Framework Support

To add support for a new framework:

1. Add patterns to `internal/analyzer/patterns.go`
2. Test with sample projects
3. Submit a PR with examples

## License

MIT License - see LICENSE file for details

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI framework
- Uses [Google Gemini](https://ai.google.dev/) for AI summaries
- Table formatting by [tablewriter](https://github.com/olekukonko/tablewriter)
- Colors by [fatih/color](https://github.com/fatih/color)
