#
# Copyright (c) 2025. Jakub Marszałek / hesoyamcc
# All Rights Reserved.
#

# Import UI and config files
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
$UiFile = Join-Path $ScriptDir "_ui.ps1"
$ConfigFile = Join-Path $ScriptDir "_list-config.ps1"

if (Test-Path $UiFile) {
    . $UiFile
} else {
    Write-Host "UI helper file not found: $UiFile" -ForegroundColor Red
    exit 1
}

if (Test-Path $ConfigFile) {
    . $ConfigFile
} else {
    Write-Host "UI helper file not found: $ConfigFile" -ForegroundColor Red
    exit 1
}

# Display list of PowerShell scripts
$scriptInfo = @(
    @{
        Script = "test-stack.ps1"
        Description = "Performs a full Health check of all services in the stack."
        Flags = "None"
    },
    @{
        Script = "start-stack.ps1"
        Description = "Starts the stack based on already built images and runs the Health check."
        Flags = "None"
    },
    @{
        Script = "build-stack.ps1"
        Description = "Builds Docker images for the entire stack."
        Flags = "-BuildAndStart (build and start stack)"
    },
    @{
        Script = "shutdown-stack.ps1"
        Description = "Stops all containers in the stack."
        Flags = "-Remove (also removes volumes)"
    },
    @{
        Script = "reset-stack.ps1"
        Description = "Performs a full stack reset: shutdown + build + start."
        Flags = "None"
    }
)

#
# List sequence
#
DrawLine
ShowLogo

# Variables loaded from _list-config.ps1
DrawLine
    Write-Host "Project     : $Project" -ForegroundColor White
    Write-Host "Author      : $Author" -ForegroundColor White
    Write-Host "Environment : $Environment" -ForegroundColor White
    Write-Host "Version     : $Version" -ForegroundColor White
    Write-Host ""
    Write-Host "$Description" -ForegroundColor White
DrawLine

# Display table / list
foreach ($s in $scriptInfo) {
    Write-Host "Script      : $($s.Script)" -ForegroundColor Yellow
    Write-Host "Description : $($s.Description)" -ForegroundColor White
    Write-Host "Flags       : $($s.Flags)" -ForegroundColor Green
    DrawShortLine
}

# Display footer
Write-Host ""
Write-Host "Usage: .\scripts\PowerShell\<script_name>.ps1 [flags]" -ForegroundColor Cyan
Write-Host "Example: .\scripts\PowerShell\build-stack.ps1 -BuildAndStart" -ForegroundColor Cyan
DrawLine
