#
# Copyright (c) 2025. Jakub Marszałek / hesoyamcc
# All Rights Reserved.
#

#
# Configuration / Constants
#
# Path related constants
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
$ShutdownScript = Join-Path $PSScriptRoot "shutdown-stack.ps1"
$BuildScript = Join-Path $PSScriptRoot "build-stack.ps1"
$UiFile = Join-Path $ScriptDir "_ui.ps1"

# Import UI file
if (Test-Path $UiFile) {
    . $UiFile
} else {
    Write-Host "UI helper file not found: $UiFile" -ForegroundColor Red
    exit 1
}

#
# Reset sequence
#
writeHeader -Text "Reseting docker stack..."

# Step 1: Shutdown stack and remove volumes
Write-Host "1. Whole stack shutting down..." -ForegroundColor Yellow
& $ShutdownScript -Remove
Write-Host ""

# Step 2: Clean build with startup
Write-Host "2. Building and starting stack..." -ForegroundColor Yellow
& $BuildScript -BuildAndStart
Write-Host ""

# Summary
writeHeader -Text "Reset completed"
