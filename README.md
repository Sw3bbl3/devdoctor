# DevDoctor 🩺

DevDoctor is a tiny, cross-platform CLI that diagnoses why a freshly cloned project does not run locally.

## Features

- 🔍 **Auto-detection** - Automatically detects project types (Node.js, Python, Go, Java, Ruby, Rust, .NET, Docker)
- 🚫 **Read-only** - Never modifies your system or installs anything automatically
- 📊 **Clear reporting** - Provides actionable suggestions to fix issues
- 🌍 **Cross-platform** - Works on Windows, macOS, and Linux
- 🚀 **Single binary** - No runtime dependencies required

## Installation

### Download Pre-built Binary

Download the latest release from the [releases page](https://github.com/Sw3bbl3/devdoctor/releases).

### Build from Source

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
