#
# Copyright (c) 2025. Jakub Marszałek / hesoyamcc
# All Rights Reserved.
#

param(
    [switch]$Remove
)

# ===================================================
# Shut down whole stack
#
# Use flag:
#     -Remove    - To also remove volumes
#
# NOTE: By default none of the volumes will be removed
# ====================================================

#
# Configuration / Constants
#
# Path related constants
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
$RootDir = Resolve-Path (Join-Path $ScriptDir "..\..")
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
# Shutdown sequence
#
writeHeader -Text "Stack shutting down..."

# Step 1: Shutdown whole docker stack
if ($Remove) {
    Write-Host "1. Whole stack shutting down (removing volumes)..." -ForegroundColor Yellow
    docker-compose --env-file $EnvFile down -v
    Write-Host ""
} else {
    Write-Host "1. Whole stack shutting down (keeping volumes)..." -ForegroundColor Yellow
    docker-compose --env-file $EnvFile down
    Write-Host ""
}

# Summary
writeHeader -Text "Shutdown completed"
