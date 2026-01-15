#!/usr/bin/env pwsh
<#
.SYNOPSIS
    DevDoctor Setup - Automated development environment installer
.DESCRIPTION
    Installs Go and builds DevDoctor automatically. Supports Windows with multiple installation methods.
.PARAMETER SkipBuild
    Install Go but skip building DevDoctor
.PARAMETER Force
    Force reinstall even if Go is already installed
.PARAMETER Verbose
    Show detailed progress and debug information
.EXAMPLE
    .\setup.ps1
    .\setup.ps1 -Force
    .\setup.ps1 -SkipBuild
#>

param(
    [switch]$SkipBuild,
    [switch]$Force,
    [switch]$Verbose
)

$ErrorActionPreference = "Stop"
$Script:StartTime = Get-Date

# Configuration
$GO_VERSION = "1.25.6"
$GO_DOWNLOAD_BASE = "https://go.dev/dl"
$REQUIRED_GO_VERSION = "1.25"

# Colors
$Script:Colors = @{
    Header = "Cyan"
    Success = "Green"
    Warning = "Yellow"
    Error = "Red"
    Info = "Gray"
    Highlight = "White"
}

# Symbols
$Script:Symbols = @{
    Success = "[OK]"
    Error = "[ERROR]"
    Warning = "[WARN]"
    Info = "[INFO]"
    Arrow = "==>"
    Bullet = " *"
}

#region UI Functions

function Write-Header {
    param([string]$Text)
    Write-Host ""
    Write-Host "================================================================" -ForegroundColor $Colors.Header
    Write-Host "  $Text" -ForegroundColor $Colors.Header
    Write-Host "================================================================" -ForegroundColor $Colors.Header
    Write-Host ""
}

function Write-Step {
    param([string]$Text)
    Write-Host ""
    Write-Host "$($Symbols.Arrow) $Text" -ForegroundColor $Colors.Highlight
}

function Write-Success {
    param([string]$Text)
    Write-Host "$($Symbols.Success) $Text" -ForegroundColor $Colors.Success
}

function Write-Info {
    param([string]$Text, [switch]$NoNewline)
    if ($NoNewline) {
        Write-Host "$($Symbols.Info) $Text" -ForegroundColor $Colors.Info -NoNewline
    } else {
        Write-Host "$($Symbols.Info) $Text" -ForegroundColor $Colors.Info
    }
}

function Write-Warn {
    param([string]$Text)
    Write-Host "$($Symbols.Warning) $Text" -ForegroundColor $Colors.Warning
}

function Write-ErrorMsg {
    param([string]$Text)
    Write-Host "$($Symbols.Error) $Text" -ForegroundColor $Colors.Error
}

function Write-Detail {
    param([string]$Text)
    if ($Verbose) {
        Write-Host "    $Text" -ForegroundColor DarkGray
    }
}

function Show-Spinner {
    param([ScriptBlock]$Action, [string]$Message)
    
    $spinnerChars = @('|', '/', '-', '\')
    $job = Start-Job -ScriptBlock $Action
    $i = 0
    
    Write-Host "$($Symbols.Info) $Message " -ForegroundColor $Colors.Info -NoNewline
    
    while ($job.State -eq 'Running') {
        Write-Host "`b$($spinnerChars[$i % 4])" -NoNewline -ForegroundColor $Colors.Highlight
        Start-Sleep -Milliseconds 100
        $i++
    }
    
    $result = Receive-Job -Job $job -Wait
    Remove-Job -Job $job
    
    Write-Host "`b " -NoNewline
    return $result
}

#endregion

#region System Info

function Get-SystemInfo {
    $os = [System.Environment]::OSVersion
    $arch = if ([Environment]::Is64BitOperatingSystem) { "x64" } else { "x86" }
    $psVersion = $PSVersionTable.PSVersion.ToString()
    
    return @{
        OS = "Windows $($os.Version.Major).$($os.Version.Minor)"
        Architecture = $arch
        PowerShell = "PowerShell $psVersion"
    }
}

function Show-SystemInfo {
    Write-Step "System Information"
    $info = Get-SystemInfo
    Write-Info "OS: $($info.OS)"
    Write-Info "Architecture: $($info.Architecture)"
    Write-Info "Shell: $($info.PowerShell)"
}

#endregion

#region Go Installation

function Test-GoInstalled {
    try {
        $goVersion = & go version 2>&1
        if ($LASTEXITCODE -eq 0 -and $goVersion -match 'go version go([\d.]+)') {
            return @{
                Installed = $true
                Version = $matches[1]
                FullVersion = $goVersion
            }
        }
    } catch {
        Write-Detail "Go not found in PATH"
    }
    
    return @{ Installed = $false }
}

function Install-GoViaWinget {
    Write-Step "Installing Go via Windows Package Manager (winget)"
    
    try {
        Write-Host ""
        & winget install GoLang.Go --silent --accept-source-agreements --accept-package-agreements
        
        if ($LASTEXITCODE -eq 0 -or $LASTEXITCODE -eq -1978335189) {
            Start-Sleep -Seconds 2
            Update-EnvironmentPath
            Write-Success "Go installed successfully via winget"
            return $true
        }
    } catch {
        Write-Detail "Winget failed: $_"
    }
    
    return $false
}

function Install-GoViaChocolatey {
    Write-Step "Installing Go via Chocolatey"
    
    try {
        Write-Host ""
        & choco install golang -y
        
        if ($LASTEXITCODE -eq 0) {
            Start-Sleep -Seconds 2
            Update-EnvironmentPath
            Write-Success "Go installed successfully via Chocolatey"
            return $true
        }
    } catch {
        Write-Detail "Chocolatey failed: $_"
    }
    
    return $false
}

function Install-GoManual {
    Write-Step "Installing Go via Manual Download"
    
    $arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
    $goInstaller = "go$GO_VERSION.windows-$arch.msi"
    $downloadUrl = "$GO_DOWNLOAD_BASE/$goInstaller"
    $installerPath = Join-Path $env:TEMP $goInstaller
    
    try {
        Write-Info "Downloading: $goInstaller"
        Write-Detail "URL: $downloadUrl"
        Write-Detail "Destination: $installerPath"
        Write-Host ""
        
        # Download with progress
        $webClient = New-Object System.Net.WebClient
        
        Register-ObjectEvent -InputObject $webClient -EventName DownloadProgressChanged -Action {
            $percent = $Event.SourceEventArgs.ProgressPercentage
            $received = [math]::Round($Event.SourceEventArgs.BytesReceived / 1MB, 2)
            $total = [math]::Round($Event.SourceEventArgs.TotalBytesToReceive / 1MB, 2)
            
            Write-Progress -Activity "Downloading Go" `
                -Status "$percent% ($received MB / $total MB)" `
                -PercentComplete $percent
        } | Out-Null
        
        $downloadTask = $webClient.DownloadFileTaskAsync($downloadUrl, $installerPath)
        
        while (-not $downloadTask.IsCompleted) {
            Start-Sleep -Milliseconds 100
        }
        
        Write-Progress -Activity "Downloading Go" -Completed
        Get-EventSubscriber | Where-Object { $_.SourceObject -eq $webClient } | Unregister-Event
        $webClient.Dispose()
        
        if ($downloadTask.Exception) {
            throw $downloadTask.Exception
        }
        
        Write-Success "Download complete"
        
        # Install
        Write-Info "Installing Go (this may require administrator privileges)..."
        $process = Start-Process msiexec.exe -Wait -PassThru -ArgumentList "/i `"$installerPath`" /qn /norestart" -NoNewWindow
        
        if ($process.ExitCode -ne 0 -and $process.ExitCode -ne 3010) {
            throw "MSI installer failed with exit code: $($process.ExitCode)"
        }
        
        # Cleanup
        Remove-Item $installerPath -Force -ErrorAction SilentlyContinue
        
        Start-Sleep -Seconds 2
        Update-EnvironmentPath
        
        Write-Success "Go installed successfully"
        return $true
        
    } catch {
        Write-Detail "Manual installation failed: $_"
        if (Test-Path $installerPath) {
            Remove-Item $installerPath -Force -ErrorAction SilentlyContinue
        }
    }
    
    return $false
}

function Update-EnvironmentPath {
    $machinePath = [System.Environment]::GetEnvironmentVariable("Path", "Machine")
    $userPath = [System.Environment]::GetEnvironmentVariable("Path", "User")
    $env:Path = "$machinePath;$userPath"
    
    # Add common Go paths
    $goPaths = @(
        "C:\Program Files\Go\bin",
        "C:\Go\bin",
        "$env:USERPROFILE\go\bin"
    )
    
    foreach ($path in $goPaths) {
        if ((Test-Path $path) -and ($env:Path -notlike "*$path*")) {
            $env:Path = "$path;$env:Path"
            Write-Detail "Added to PATH: $path"
        }
    }
}

function Install-Go {
    Write-Step "Go Installation"
    
    # Check if already installed
    $goInfo = Test-GoInstalled
    
    if ($goInfo.Installed -and -not $Force) {
        Write-Success "Go is already installed"
        Write-Info "Version: $($goInfo.FullVersion)"
        
        if ($goInfo.Version -ge $REQUIRED_GO_VERSION) {
            Write-Success "Version meets requirements (>= $REQUIRED_GO_VERSION)"
            return $true
        } else {
            Write-Warn "Version $($goInfo.Version) is below recommended $REQUIRED_GO_VERSION"
            Write-Info "Consider upgrading Go"
        }
        return $true
    }
    
    if ($Force) {
        Write-Info "Force reinstall requested"
    }
    
    # Try installation methods in order
    $methods = @()
    
    if (Get-Command winget -ErrorAction SilentlyContinue) {
        $methods += @{ Name = "winget"; Function = { Install-GoViaWinget } }
    }
    
    if (Get-Command choco -ErrorAction SilentlyContinue) {
        $methods += @{ Name = "chocolatey"; Function = { Install-GoViaChocolatey } }
    }
    
    $methods += @{ Name = "manual"; Function = { Install-GoManual } }
    
    Write-Info "Available installation methods: $($methods.Name -join ', ')"
    
    foreach ($method in $methods) {
        Write-Detail "Trying: $($method.Name)"
        
        if (& $method.Function) {
            # Verify installation
            Start-Sleep -Seconds 1
            $goInfo = Test-GoInstalled
            
            if ($goInfo.Installed) {
                Write-Host ""
                Write-Success "Go installation verified"
                Write-Info "Version: $($goInfo.FullVersion)"
                return $true
            }
        }
    }
    
    throw "Failed to install Go using any method"
}

#endregion

#region Build

function Build-DevDoctor {
    Write-Step "Building DevDoctor"
    
    # Verify go is available
    $goInfo = Test-GoInstalled
    if (-not $goInfo.Installed) {
        throw "Go is not available. Cannot build DevDoctor."
    }
    
    try {
        # Download dependencies
        Write-Info "Downloading Go modules..."
        & go mod download
        
        if ($LASTEXITCODE -ne 0) {
            throw "Failed to download dependencies"
        }
        
        Write-Success "Dependencies downloaded"
        
        # Build
        Write-Info "Compiling DevDoctor..."
        $outputBinary = if ($IsWindows -or $env:OS -eq "Windows_NT") { "devdoctor.exe" } else { "devdoctor" }
        
        & go build -ldflags "-s -w" -o $outputBinary ./cmd/devdoctor
        
        if ($LASTEXITCODE -ne 0) {
            throw "Build failed"
        }
        
        # Verify binary
        if (-not (Test-Path $outputBinary)) {
            throw "Binary not found after build"
        }
        
        $binarySize = [math]::Round((Get-Item $outputBinary).Length / 1MB, 2)
        
        Write-Success "Build complete!"
        Write-Info "Binary: $outputBinary ($binarySize MB)"
        Write-Info "Location: $(Get-Location)\$outputBinary"
        
        return $true
        
    } catch {
        Write-ErrorMsg "Build failed: $_"
        throw
    }
}

# Install devdoctor globally so it's available as 'devdoctor'
function Ensure-UserPathContains {
    param([string]$PathToAdd)
    $userPath = [System.Environment]::GetEnvironmentVariable("Path", "User")
    if ($userPath -notlike "*$PathToAdd*") {
        $newPath = "$userPath;$PathToAdd"
        [System.Environment]::SetEnvironmentVariable("Path", $newPath, "User")
        Write-Info "Added to User PATH: $PathToAdd"
    } else {
        Write-Detail "Path already contains: $PathToAdd"
    }
}

function Install-DevDoctorGlobal {
    Write-Step "Installing DevDoctor globally"
    try {
        & go install ./cmd/devdoctor
        if ($LASTEXITCODE -ne 0) { throw "go install failed" }
        
        $gobin = Join-Path $env:USERPROFILE "go\bin"
        $exe = Join-Path $gobin "devdoctor.exe"
        
        if (-not (Test-Path $exe)) {
            # Fallback: copy local build
            $localExe = Join-Path (Get-Location) "devdoctor.exe"
            if (Test-Path $localExe) {
                New-Item -ItemType Directory -Path $gobin -Force | Out-Null
                Copy-Item $localExe $exe -Force
                Write-Info "Copied local binary to: $exe"
            } else {
                throw "Binary not found in $gobin or local build"
            }
        }
        
        Ensure-UserPathContains $gobin
        Update-EnvironmentPath
        
        Write-Success "DevDoctor installed to $gobin"
        Write-Info "You can now run: devdoctor"
        return $true
    } catch {
        Write-ErrorMsg "Global installation failed: $_"
        throw
    }
}

#endregion

#region Main

function Show-Banner {
    Write-Host ""
    Write-Host "  ____             ____             _             " -ForegroundColor Cyan
    Write-Host " |  _ \  _____   _|  _ \  ___   ___| |_ ___  _ __ " -ForegroundColor Cyan
    Write-Host " | | | |/ _ \ \ / / | | |/ _ \ / __| __/ _ \| '__|" -ForegroundColor Cyan
    Write-Host " | |_| |  __/\ V /| |_| | (_) | (__| || (_) | |   " -ForegroundColor Cyan
    Write-Host " |____/ \___| \_/ |____/ \___/ \___|\__\___/|_|   " -ForegroundColor Cyan
    Write-Host ""
    Write-Host "          Automated Development Environment Setup" -ForegroundColor Gray
    Write-Host ""
}

function Show-Summary {
    $elapsed = (Get-Date) - $Script:StartTime
    $duration = "{0:N1}s" -f $elapsed.TotalSeconds
    
    Write-Host ""
    Write-Host "================================================================" -ForegroundColor $Colors.Header
    Write-Host "  Setup Complete!" -ForegroundColor $Colors.Success
    Write-Host "================================================================" -ForegroundColor $Colors.Header
    Write-Host ""
    Write-Info "Time elapsed: $duration"
    Write-Host ""
    Write-Host "  Next steps:" -ForegroundColor $Colors.Highlight
    Write-Host "    1. Run DevDoctor:  " -NoNewline -ForegroundColor Gray
    Write-Host ".\devdoctor.exe" -ForegroundColor White
    Write-Host "    2. View help:      " -NoNewline -ForegroundColor Gray
    Write-Host ".\devdoctor.exe --help" -ForegroundColor White
    Write-Host "    3. Scan a project: " -NoNewline -ForegroundColor Gray
    Write-Host ".\devdoctor.exe -path C:\path\to\project" -ForegroundColor White
    Write-Host ""
}

function Main {
    try {
        Show-Banner
        
        if ($Verbose) {
            Write-Info "Verbose mode enabled"
        }
        
        Show-SystemInfo
        
        # Install Go
        Install-Go | Out-Null
        
        # Build DevDoctor
        if (-not $SkipBuild) {
            Write-Host ""
            Build-DevDoctor | Out-Null
            Write-Host ""
            Install-DevDoctorGlobal | Out-Null
        } else {
            Write-Info "Skipping build (--SkipBuild specified)"
        }
        
        Show-Summary
        
    } catch {
        Write-Host ""
        Write-Host "================================================================" -ForegroundColor $Colors.Error
        Write-Host "  Setup Failed" -ForegroundColor $Colors.Error
        Write-Host "================================================================" -ForegroundColor $Colors.Error
        Write-Host ""
        Write-ErrorMsg $_.Exception.Message
        Write-Host ""
        Write-Host "  Troubleshooting:" -ForegroundColor $Colors.Warning
        Write-Host "    * Try running as Administrator" -ForegroundColor Gray
        Write-Host "    * Check your internet connection" -ForegroundColor Gray
        Write-Host "    * Install Go manually: https://go.dev/dl/" -ForegroundColor Gray
        Write-Host "    * Run with -Verbose flag for details" -ForegroundColor Gray
        Write-Host ""
        exit 1
    }
}

# Run
Main

#endregion
