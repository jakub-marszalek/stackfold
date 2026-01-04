#
# Copyright (c) 2025. Jakub Marszałek / hesoyamcc
# All Rights Reserved.
#

#
# Default configuration
#

# UI Constants
$global:UIConfig = @{
    LineChar = "-"
    EdgeChar = "+"
    LineWidth = 80
    ShortLineWidth = 55
    TitleColor = "Cyan"
    SeparatorLineColor = "DarkGray"
}

# Info tags used in scripts
$START_TAG="[START]"
$TEST_TAG="[TEST]"
$BUILD_TAG="[BUILD]"
$SHUTDOWN_TAG="[SHUTDOWN]"
$RESET_TAG="[RESET]"

#
# UI functions
#

# Print Stackr logo
function showLogo {
    Write-Host "    _____ __             __   ____      __    __ " -ForegroundColor Cyan
    Write-Host "   / ___// /_____ ______/ /__/ __/___  / /___/ / " -ForegroundColor Cyan
    Write-Host "   \__ \/ __/ __ '/ ___/ //_/ /_/ __ \/ / __  /  " -ForegroundColor Cyan
    Write-Host "  ___/ / /_/ /_/ / /__/ ,< / __/ /_/ / / /_/ /   " -ForegroundColor Cyan
    Write-Host " /____/\__/\__,_/\___/_/|_/_/  \____/_/\__,_/    " -ForegroundColor Cyan
    Write-Host ""
}

# Draw simple line
function drawLine {
    param (
        [int]$Width = $global:UIConfig.LineWidth,
        [string]$EdgeChar = $global:UIConfig.EdgeChar,
        [string]$LineChar = $global:UIConfig.LineChar,
        [string]$Color = $global:UIConfig.TitleColor
    )

    $InnerWidth = [Math]::Max($Width - 2, 1)
    $Line = "$EdgeChar" + ($LineChar * $InnerWidth) + "$EdgeChar"

    Write-Host $Line -ForegroundColor $Color
}

# Draw short line
function drawShortLine {
    param (
        [int]$Width = $global:UIConfig.ShortLineWidth,
        [string]$LineChar = $global:UIConfig.LineChar,
        [string]$Color = $global:UIConfig.SeparatorLineColor
    )

    $Line = ($LineChar * $Width)
    Write-Host $Line -ForegroundColor $Color
}

# Write centered header text
function writeHeader {
    param (
        [string]$Text,
        [int]$Width = $global:UIConfig.LineWidth,
        [string]$Color = $global:UIConfig.TitleColor
    )

    DrawLine
        $Padding = [Math]::Max(0, [Math]::Floor(($Width - $Text.Length) / 2))
        $Line = (" " * $Padding) + $Text
        Write-Host $Line -ForegroundColor $Color
    DrawLine
}
