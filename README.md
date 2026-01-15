# DevDoctor 🩺

DevDoctor is a tiny, cross-platform CLI that diagnoses why a freshly cloned project does not run locally.

## Features

- 🔍 **Auto-detection** - Automatically detects project types (Node.js, Python, Go, Java, Ruby, Rust, .NET, Docker)
- 🚫 **Read-only** - Never modifies your system or installs anything automatically
- 📊 **Clear reporting** - Provides actionable suggestions to fix issues
- 🌍 **Cross-platform** - Works on Windows, macOS, and Linux
- 🚀 **Single binary** - No runtime dependencies required

## Installation

### Quick Setup (Recommended)

The easiest way to get started is using our automated setup script that installs Go and builds DevDoctor:

**Windows (PowerShell):**
```powershell
# Clone and setup
git clone https://github.com/Sw3bbl3/devdoctor.git
cd devdoctor

# Run setup (PowerShell)
.\setup.ps1

# Or use the batch launcher
setup.bat

# Advanced options
.\setup.ps1 -Verbose      # Show detailed progress
.\setup.ps1 -Force        # Force reinstall Go
.\setup.ps1 -SkipBuild    # Only install Go
```

**Windows (CMD):**
```cmd
git clone https://github.com/Sw3bbl3/devdoctor.git
cd devdoctor
setup.bat
```

**Linux/macOS:**
```bash
git clone https://github.com/Sw3bbl3/devdoctor.git
cd devdoctor
chmod +x setup.sh

# Run setup
./setup.sh

# Advanced options
./setup.sh --verbose      # Show detailed progress
./setup.sh --force        # Force reinstall Go
./setup.sh --skip-build   # Only install Go
```

The setup script will automatically:
- ✅ Detect your operating system and architecture
- ✅ Check for existing Go installation
- ✅ Install Go using the best method available:
  - Windows: winget → Chocolatey → MSI installer
  - macOS: Homebrew → Manual download
  - Linux: Manual download and extraction
- ✅ Download and cache project dependencies
- ✅ Build optimized DevDoctor binary
- ✅ Verify installation and provide usage instructions

**Features:**
- 🎨 Beautiful terminal UI with progress indicators
- 🔄 Automatic fallback to alternative installation methods
- 📊 Real-time download progress with percentage and speed
- ✔️ Version checking and requirement validation
- 🔍 Verbose mode for troubleshooting
- ⚡ Production-ready and developer-friendly

### Download Pre-built Binary

Download the latest release from the [releases page](https://github.com/Sw3bbl3/devdoctor/releases).

### Manual Build from Source

If you already have Go installed:

```bash
git clone https://github.com/Sw3bbl3/devdoctor.git
cd devdoctor
go build -o devdoctor ./cmd/devdoctor
```

## Usage

Run DevDoctor in your project directory:

```bash
devdoctor
```

Or specify a path:

```bash
devdoctor -path /path/to/project
```

### Version & Updates

```bash
devdoctor -version
devdoctor -check-update
devdoctor -update
```

## Example Output

```
╔═══════════════════════════════════════════════════════════════╗
║                         DEVDOCTOR                             ║
║              Project Diagnostic Report                        ║
╚═══════════════════════════════════════════════════════════════╝

Scanning: /home/user/my-project

📋 Detected Project Types:
  ✓ Node.js
    - package.json

⚠️  WARNINGS (Issues that may cause problems):
─────────────────────────────────────────────────────────────────

[Node.js] Dependencies not installed (node_modules directory not found)
   💡 Run 'npm install' or 'yarn install' to install dependencies

═════════════════════════════════════════════════════════════════
Summary: 0 error(s), 1 warning(s), 0 info

⚠️  Consider addressing the warnings to ensure smooth operation.
═════════════════════════════════════════════════════════════════
```

## Supported Project Types

- **Node.js** - Detects `package.json`, checks for `node_modules`, verifies Node version requirements
- **Python** - Detects `requirements.txt`, `setup.py`, `pyproject.toml`, checks for virtual environments
- **Go** - Detects `go.mod`, checks for `go.sum` and vendor directory
- **Java** - Detects `pom.xml` (Maven) or `build.gradle` (Gradle), checks for build artifacts
- **Ruby** - Detects `Gemfile`, checks for `Gemfile.lock`
- **Rust** - Detects `Cargo.toml`, checks for build artifacts
- **.NET** - Detects `.csproj`, `.sln` files, checks for build artifacts
- **Docker** - Detects `Dockerfile`, `docker-compose.yml`, checks Docker daemon status

## What DevDoctor Checks

### For All Projects
- ✅ Required development tools are installed (e.g., `node`, `python`, `go`)
- ✅ Tools are accessible in PATH

### Project-Specific Checks
- ✅ Dependencies are installed
- ✅ Build artifacts exist
- ✅ Configuration files are present
- ✅ Version requirements (where specified)
- ✅ Environment files (`.env`) when examples exist

## What DevDoctor Does NOT Do

- ❌ Install tools or dependencies automatically
- ❌ Modify your project files
- ❌ Make system-level changes
- ❌ Use AI or cloud services
- ❌ Require internet connectivity (except for Docker daemon check)

## Exit Codes

- `0` - No issues found or no supported project detected
- `1` - Issues detected that may prevent the project from running

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) file for details.
