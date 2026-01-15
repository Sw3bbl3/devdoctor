@echo off
REM DevDoctor Setup - Windows Launcher
REM Automatically detects and uses the best PowerShell version available

setlocal enabledelayedexpansion

echo.
echo ================================================================
echo   DevDoctor Setup
echo ================================================================
echo.

REM Check for PowerShell 7+ (pwsh)
where pwsh >nul 2>nul
if %errorlevel% equ 0 (
    echo [INFO] Using PowerShell 7+
    echo.
    pwsh -ExecutionPolicy Bypass -File "%~dp0setup.ps1" %*
    goto :end
)

REM Check for Windows PowerShell 5.x
where powershell >nul 2>nul
if %errorlevel% equ 0 (
    echo [INFO] Using Windows PowerShell
    echo.
    powershell -ExecutionPolicy Bypass -File "%~dp0setup.ps1" %*
    goto :end
)

REM No PowerShell found
echo [ERROR] PowerShell not found!
echo.
echo Please install PowerShell to run this setup script.
echo Download from: https://github.com/PowerShell/PowerShell
echo.
pause
exit /b 1

:end
if %errorlevel% neq 0 (
    echo.
    echo Setup encountered an error. Exit code: %errorlevel%
    pause
)
