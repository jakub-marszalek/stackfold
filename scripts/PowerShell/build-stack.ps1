#
# Copyright (c) 2025. Jakub Marszałek / hesoyamcc
# All Rights Reserved.
#

param(
    [switch]$BuildAndStart
)

# ===================================================
# Build Whole Stack
#
# Use flag:
#     -BuildAndStart - To build images and start them
#
# NOTE: By default only build will be executed
# ===================================================

#
# Configuration / Constants
#
# Path related constants
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
$RootDir = Resolve-Path (Join-Path $ScriptDir "..\..")
$StartScript = Join-Path $PSScriptRoot "start-stack.ps1"
$UiFile = Join-Path $ScriptDir "_ui.ps1"
$EnvFile = Join-Path $RootDir ".env.docker"

# Import UI file
if (Test-Path $UiFile) {
    . $UiFile
} else {
    Write-Host "UI helper file not found: $UiFile" -ForegroundColor Red
    exit 1
}

#
# Build sequence
#
writeHeader -Text "Building docker stack..."

if ($BuildAndStart) {
    Write-Host "Option: Build and Start" -ForegroundColor Cyan
}
Write-Host ""

# Validate .env.docker path
if (-Not (Test-Path $EnvFile)) {
    Write-Host ".env.docker not found in: $EnvFile" -ForegroundColor Red
    exit 1
}

# Step 1: Build whole docker stack
Write-Host "1. Building docker stack..." -ForegroundColor Yellow
docker-compose --env-file $EnvFile build --no-cache
Write-Host ""

# Step 2: (If build and start enabled)
if ($BuildAndStart) {
    Write-Host "2. Starting docker stack..." -ForegroundColor Yellow
    & $StartScript
    Write-Host ""
}

# Summary
writeHeader -Text "Build completed"
