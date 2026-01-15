# DevDoctor Setup Guide

## Quick Start

### Windows
```powershell
.\setup.ps1
```

### Linux/macOS
```bash
chmod +x setup.sh
./setup.sh
```

## Advanced Usage

### PowerShell Options

```powershell
# Show detailed progress and debug information
.\setup.ps1 -Verbose

# Force reinstall Go even if already installed
.\setup.ps1 -Force

# Only install Go, skip building DevDoctor
.\setup.ps1 -SkipBuild

# Combine options
.\setup.ps1 -Verbose -Force
```

### Bash Options

```bash
# Show detailed progress
./setup.sh --verbose

# Force reinstall Go
./setup.sh --force

# Only install Go
./setup.sh --skip-build

# Show help
./setup.sh --help
```

## Installation Methods

### Windows
The setup script tries these methods in order:
1. **winget** (Windows Package Manager) - Fastest, modern Windows 10/11
2. **Chocolatey** - Popular package manager
3. **MSI Installer** - Direct download with progress bar

### macOS
1. **Homebrew** - Standard package manager for macOS
2. **Manual Download** - Direct download and extraction

### Linux
1. **Manual Download** - Downloads and extracts to `~/.local/go` or `/usr/local/go`

## Troubleshooting

### "Go not found after installation"
**Solution:** Restart your terminal or run:
```powershell
# Windows
$env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path", "User")

# Linux/macOS
source ~/.bashrc  # or ~/.zshrc
```

### "Permission denied"
**Windows:** Run PowerShell as Administrator
**Linux/macOS:** Use `sudo` or install to user directory

### "Download is slow"
**Solution:** The script will show progress. If using MSI on Windows, this is normal (50MB download).
Use `winget` or `choco` for faster installation.

### "Script execution is disabled"
**Windows PowerShell:**
```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### "Build failed"
**Solution:** Run with verbose mode to see details:
```powershell
.\setup.ps1 -Verbose
```

Common causes:
- Network issues during dependency download
- Insufficient disk space
- Antivirus blocking compilation

## Manual Installation

If the automated setup fails, you can install manually:

1. **Install Go** from https://go.dev/dl/
2. **Clone the repository:**
   ```bash
   git clone https://github.com/Sw3bbl3/devdoctor.git
   cd devdoctor
   ```
3. **Build:**
   ```bash
   go build -o devdoctor ./cmd/devdoctor
   ```

## Features

✅ **Automatic Platform Detection** - Detects OS and architecture
✅ **Multiple Install Methods** - Tries best method first, falls back automatically
✅ **Progress Indicators** - Real-time download and build progress
✅ **Version Checking** - Ensures Go version meets requirements
✅ **Smart PATH Management** - Automatically configures environment
✅ **Error Recovery** - Clear error messages and troubleshooting tips
✅ **Production Ready** - Used by developers worldwide

## System Requirements

- **Windows:** Windows 10/11 with PowerShell 5.1+
- **macOS:** macOS 10.15+ (Catalina or later)
- **Linux:** Any modern distribution
- **Disk Space:** ~100MB for Go + ~10MB for DevDoctor
- **Internet:** Required for downloading Go and dependencies

## Support

- **Issues:** https://github.com/Sw3bbl3/devdoctor/issues
- **Discussions:** https://github.com/Sw3bbl3/devdoctor/discussions
- **Documentation:** https://github.com/Sw3bbl3/devdoctor

## License

MIT License - See LICENSE file for details
