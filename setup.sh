#!/usr/bin/env bash
#
# DevDoctor Setup - Automated development environment installer
# Installs Go and builds DevDoctor automatically for Linux/macOS
#

set -e

# Configuration
GO_VERSION="1.25.6"
GO_DOWNLOAD_BASE="https://go.dev/dl"
REQUIRED_GO_VERSION="1.25"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
GRAY='\033[0;90m'
WHITE='\033[1;37m'
NC='\033[0m' # No Color

# Flags
SKIP_BUILD=false
FORCE=false
VERBOSE=false
START_TIME=$(date +%s)

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --skip-build)
            SKIP_BUILD=true
            shift
            ;;
        --force)
            FORCE=true
            shift
            ;;
        --verbose|-v)
            VERBOSE=true
            shift
            ;;
        --help|-h)
            echo "DevDoctor Setup - Automated installer"
            echo ""
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --skip-build    Install Go but skip building DevDoctor"
            echo "  --force         Force reinstall even if Go exists"
            echo "  --verbose, -v   Show detailed progress"
            echo "  --help, -h      Show this help message"
            echo ""
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

#region UI Functions

print_header() {
    echo ""
    echo -e "${CYAN}================================================================${NC}"
    echo -e "${CYAN}  $1${NC}"
    echo -e "${CYAN}================================================================${NC}"
    echo ""
}

print_step() {
    echo ""
    echo -e "${WHITE}==> $1${NC}"
}

print_success() {
    echo -e "${GREEN}[OK] $1${NC}"
}

print_info() {
    echo -e "${GRAY}[INFO] $1${NC}"
}

print_warn() {
    echo -e "${YELLOW}[WARN] $1${NC}"
}

print_error() {
    echo -e "${RED}[ERROR] $1${NC}"
}

print_detail() {
    if [ "$VERBOSE" = true ]; then
        echo -e "    ${GRAY}$1${NC}"
    fi
}

#endregion

#region System Info

detect_platform() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)
    
    case "$os" in
        linux)
            OS="linux"
            ;;
        darwin)
            OS="darwin"
            ;;
        *)
            print_error "Unsupported operating system: $os"
            exit 1
            ;;
    esac
    
    case "$arch" in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        armv7l|armv6l)
            ARCH="armv6l"
            ;;
        *)
            print_error "Unsupported architecture: $arch"
            exit 1
            ;;
    esac
}

show_system_info() {
    print_step "System Information"
    print_info "OS: $(uname -s) $(uname -r)"
    print_info "Architecture: $ARCH"
    print_info "Shell: $SHELL"
}

#endregion

#region Go Installation

check_go() {
    if command -v go &> /dev/null; then
        local version=$(go version | grep -oP 'go\K[0-9.]+')
        echo "$version"
        return 0
    fi
    return 1
}

install_go_homebrew() {
    print_step "Installing Go via Homebrew"
    
    if ! command -v brew &> /dev/null; then
        print_detail "Homebrew not found"
        return 1
    fi
    
    echo ""
    brew install go
    
    if [ $? -eq 0 ]; then
        print_success "Go installed successfully via Homebrew"
        return 0
    fi
    
    return 1
}

install_go_manual() {
    print_step "Installing Go via Manual Download"
    
    detect_platform
    
    local go_archive="go${GO_VERSION}.${OS}-${ARCH}.tar.gz"
    local download_url="${GO_DOWNLOAD_BASE}/${go_archive}"
    local tmp_dir=$(mktemp -d)
    
    print_info "Downloading: $go_archive"
    print_detail "URL: $download_url"
    print_detail "Temp: $tmp_dir"
    echo ""
    
    # Download
    if command -v curl &> /dev/null; then
        curl -fL --progress-bar "$download_url" -o "$tmp_dir/$go_archive"
    elif command -v wget &> /dev/null; then
        wget --show-progress -q "$download_url" -O "$tmp_dir/$go_archive"
    else
        print_error "Neither curl nor wget found"
        return 1
    fi
    
    if [ $? -ne 0 ]; then
        rm -rf "$tmp_dir"
        return 1
    fi
    
    print_success "Download complete"
    
    # Determine installation directory
    if [ "$EUID" -eq 0 ]; then
        INSTALL_DIR="/usr/local"
    else
        INSTALL_DIR="$HOME/.local"
        mkdir -p "$INSTALL_DIR"
    fi
    
    print_info "Installing to: $INSTALL_DIR"
    
    # Remove existing installation
    if [ -d "$INSTALL_DIR/go" ]; then
        print_detail "Removing existing Go installation..."
        rm -rf "$INSTALL_DIR/go"
    fi
    
    # Extract
    tar -C "$INSTALL_DIR" -xzf "$tmp_dir/$go_archive"
    
    # Cleanup
    rm -rf "$tmp_dir"
    
    # Update PATH
    export PATH="$INSTALL_DIR/go/bin:$PATH"
    
    # Add to shell profile
    local shell_profile=""
    if [ -n "$BASH_VERSION" ]; then
        shell_profile="$HOME/.bashrc"
    elif [ -n "$ZSH_VERSION" ]; then
        shell_profile="$HOME/.zshrc"
    fi
    
    if [ -n "$shell_profile" ] && [ -f "$shell_profile" ]; then
        if ! grep -q "# DevDoctor Go setup" "$shell_profile"; then
            {
                echo ""
                echo "# DevDoctor Go setup"
                echo "export PATH=\"$INSTALL_DIR/go/bin:\$PATH\""
                echo "export PATH=\"\$HOME/go/bin:\$PATH\""
            } >> "$shell_profile"
            print_detail "Added Go to $shell_profile"
        fi
    fi
    
    print_success "Go installed successfully"
    return 0
}

install_go() {
    print_step "Go Installation"
    
    # Check if already installed
    if go_version=$(check_go); then
        if [ "$FORCE" = false ]; then
            print_success "Go is already installed"
            print_info "Version: go$go_version"
            
            if [ "$(printf '%s\n' "$REQUIRED_GO_VERSION" "$go_version" | sort -V | head -n1)" = "$REQUIRED_GO_VERSION" ]; then
                print_success "Version meets requirements (>= $REQUIRED_GO_VERSION)"
            else
                print_warn "Version $go_version is below recommended $REQUIRED_GO_VERSION"
            fi
            return 0
        else
            print_info "Force reinstall requested"
        fi
    fi
    
    # Try installation methods
    local methods=()
    
    if command -v brew &> /dev/null; then
        methods+=("homebrew")
    fi
    
    methods+=("manual")
    
    print_info "Available installation methods: ${methods[*]}"
    
    for method in "${methods[@]}"; do
        print_detail "Trying: $method"
        
        case $method in
            homebrew)
                if install_go_homebrew; then
                    sleep 1
                    if go_version=$(check_go); then
                        echo ""
                        print_success "Go installation verified"
                        print_info "Version: go$go_version"
                        return 0
                    fi
                fi
                ;;
            manual)
                if install_go_manual; then
                    sleep 1
                    if go_version=$(check_go); then
                        echo ""
                        print_success "Go installation verified"
                        print_info "Version: go$go_version"
                        return 0
                    fi
                fi
                ;;
        esac
    done
    
    print_error "Failed to install Go using any method"
    return 1
}

#endregion

#region Build

build_devdoctor() {
    print_step "Building DevDoctor"
    
    # Verify go is available
    if ! command -v go &> /dev/null; then
        print_error "Go is not available. Cannot build DevDoctor."
        return 1
    fi
    
    # Download dependencies
    print_info "Downloading Go modules..."
    go mod download
    
    if [ $? -ne 0 ]; then
        print_error "Failed to download dependencies"
        return 1
    fi
    
    print_success "Dependencies downloaded"
    
    # Build
    print_info "Compiling DevDoctor..."
    go build -ldflags "-s -w" -o devdoctor ./cmd/devdoctor
    
    if [ $? -ne 0 ]; then
        print_error "Build failed"
        return 1
    fi
    
    # Make executable
    chmod +x devdoctor
    
    # Verify binary
    if [ ! -f "devdoctor" ]; then
        print_error "Binary not found after build"
        return 1
    fi
    
    local binary_size=$(du -h devdoctor | cut -f1)
    
    print_success "Build complete!"
    print_info "Binary: devdoctor ($binary_size)"
    print_info "Location: $(pwd)/devdoctor"
    
    return 0
}

# Install devdoctor globally so it's available as 'devdoctor'
install_devdoctor_global() {
    print_step "Installing DevDoctor globally"
    
    # Use go install to put binary in GOPATH/bin
    go install ./cmd/devdoctor
    if [ $? -ne 0 ]; then
        print_error "go install failed"
        return 1
    fi
    
    local gobin="$HOME/go/bin"
    local exe="$gobin/devdoctor"
    
    if [ ! -f "$exe" ]; then
        # Fallback: copy local binary
        mkdir -p "$gobin"
        cp ./devdoctor "$exe"
        print_info "Copied local binary to: $exe"
    fi
    
    # Ensure PATH contains GOPATH/bin persistently
    ensure_path_contains() {
        local dir="$1"
        case "$SHELL" in
            */bash)
                local profile="$HOME/.bashrc"
                ;;
            */zsh)
                local profile="$HOME/.zshrc"
                ;;
            *)
                local profile="$HOME/.profile"
                ;;
        esac
        
        if ! echo "$PATH" | grep -q "$dir"; then
            echo "" >> "$profile"
            echo "# DevDoctor global install" >> "$profile"
            echo "export PATH=\"$dir:\$PATH\"" >> "$profile"
            print_info "Added to PATH in $profile: $dir"
        fi
    }
    
    ensure_path_contains "$gobin"
    
    print_success "DevDoctor installed to $gobin"
    print_info "You can now run: devdoctor"
    return 0
}

#endregion

#region Main

show_banner() {
    echo ""
    echo -e "${CYAN}  ____             ____             _             ${NC}"
    echo -e "${CYAN} |  _ \  _____   _|  _ \  ___   ___| |_ ___  _ __ ${NC}"
    echo -e "${CYAN} | | | |/ _ \ \ / / | | |/ _ \ / __| __/ _ \| '__|${NC}"
    echo -e "${CYAN} | |_| |  __/\ V /| |_| | (_) | (__| || (_) | |   ${NC}"
    echo -e "${CYAN} |____/ \___| \_/ |____/ \___/ \___|\__\___/|_|   ${NC}"
    echo ""
    echo -e "${GRAY}          Automated Development Environment Setup${NC}"
    echo ""
}

show_summary() {
    local end_time=$(date +%s)
    local duration=$((end_time - START_TIME))
    
    echo ""
    echo -e "${CYAN}================================================================${NC}"
    echo -e "${GREEN}  Setup Complete!${NC}"
    echo -e "${CYAN}================================================================${NC}"
    echo ""
    print_info "Time elapsed: ${duration}s"
    echo ""
    echo -e "${WHITE}  Next steps:${NC}"
    echo -e "${GRAY}    1. Run DevDoctor:  ${WHITE}./devdoctor${NC}"
    echo -e "${GRAY}    2. View help:      ${WHITE}./devdoctor --help${NC}"
    echo -e "${GRAY}    3. Scan a project: ${WHITE}./devdoctor -path /path/to/project${NC}"
    echo ""
}

main() {
    show_banner
    
    if [ "$VERBOSE" = true ]; then
        print_info "Verbose mode enabled"
    fi
    
    show_system_info
    
    # Install Go
    if ! install_go; then
        echo ""
        print_header "Setup Failed"
        print_error "Go installation failed"
        echo ""
        echo -e "${YELLOW}  Troubleshooting:${NC}"
        echo -e "${GRAY}    * Check your internet connection${NC}"
        echo -e "${GRAY}    * Install Go manually: https://go.dev/dl/${NC}"
        echo -e "${GRAY}    * Run with --verbose flag for details${NC}"
        echo ""
        exit 1
    fi
    
    # Build DevDoctor
    if [ "$SKIP_BUILD" = false ]; then
        echo ""
        if ! build_devdoctor; then
            echo ""
            print_header "Setup Failed"
            print_error "Build failed"
            exit 1
        fi
        echo ""
        install_devdoctor_global || true
    else
        print_info "Skipping build (--skip-build specified)"
    fi
    
    show_summary
}

# Run
main

#endregion
