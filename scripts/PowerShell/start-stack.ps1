#
# Copyright (c) 2025. Jakub Marszałek / hesoyamcc
# All Rights Reserved.
#

#
# Configuration / Constants
#
# Path related constants
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
$RootDir = Resolve-Path (Join-Path $ScriptDir "..\..")
$TestScript = Join-Path $PSScriptRoot "test-stack.ps1"
$UiFile = Join-Path $ScriptDir "_ui.ps1"
$EnvFile = Join-Path $RootDir ".env.docker"

# Waiting time before running healthcheck (in seconds)
$SleepSeconds = 10

# Import UI file
if (Test-Path $UiFile) {
    . $UiFile
} else {
    Write-Host "UI helper file not found: $UiFile" -ForegroundColor Red
    exit 1
}

#
# Stack start sequence
#
writeHeader -Text "Starting docker stack..."

# Step 1: Start whole docker stack
Write-Host "1. Starting whole stack..." -ForegroundColor Yellow
docker-compose --env-file $EnvFile up -d
Write-Host ""

# Step 2: Wait untill everything start properly
Write-Host "2. Wating for services to start ($SleepSeconds seconds)..." -ForegroundColor Yellow
Start-Sleep -Seconds $SleepSeconds
Write-Host ""

# Step 3: Test whole stack's health
Write-Host "3. Running stack's healthcheck..." -ForegroundColor Yellow
& $TestScript
Write-Host ""

# Summary
writeHeader -Text "Start completed"
