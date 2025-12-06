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
$UiFile = Join-Path $ScriptDir "_ui.ps1"
$EnvFile = Join-Path $RootDir ".env.docker"

# Stack related constants
$PHP_CONTAINER = "php"
$PHP_HEALTH_CHECK_FILE = "/var/www/html/docker/Health/Stack/stack_health_check.php"
$NGINX_URL = "http://localhost:8080"
$VITE_CONTAINER = "node"
$VITE_LOG_SEARCH = "ready in"

# Ports definition (name:port)
$ports = @(
    @{name = "Nginx"; port = 8080},
    @{name = "PostgreSQL"; port = 5432},
    @{name = "Redis"; port = 6379},
    @{name = "Vite"; port = 5173}
)

# Import UI file
if (Test-Path $UiFile) {
    . $UiFile
} else {
    Write-Host "UI helper file not found: $UiFile" -ForegroundColor Red
    exit 1
}

#
# Stack tests sequence
#
writeHeader -Text "Testing docker stack..."

# Test 1: Docker Compose Status
Write-Host "$TEST_TAG Testing Docker Compose Status..." -ForegroundColor Yellow
    $containers = docker-compose --env-file $EnvFile ps --format "table {{.Names}}\t{{.Status}}"
    Write-Host $containers
Write-Host ""

# Test 2: Nginx
Write-Host "$TEST_TAG Testing Nginx server (HTTP)..." -ForegroundColor Yellow
    try {
        $response = Invoke-WebRequest -Uri $NGINX_URL -UseBasicParsing -ErrorAction Stop
        Write-Host "Nginx: OK (Status $($response.StatusCode))" -ForegroundColor Green
    } catch {
        Write-Host "Nginx: ERROR - $($_.Exception.Message)" -ForegroundColor Red
    }
Write-Host ""

# Test 3: PHP Health Check
Write-Host "$TEST_TAG Testing PHP Stack..." -ForegroundColor Yellow
    $phpResult = docker-compose --env-file $EnvFile exec -T $PHP_CONTAINER php $PHP_HEALTH_CHECK_FILE
    $phpResult | ConvertFrom-Json | ConvertTo-Json -Depth 10 | Write-Host
Write-Host ""

# Test 4: Vite Dev Server
Write-Host "$TEST_TAG Testing Vite Dev Server..." -ForegroundColor Yellow
    $viteStatus = docker-compose --env-file $EnvFile logs --tail=5 $VITE_CONTAINER | Select-String $VITE_LOG_SEARCH
    if ($viteStatus) {
        Write-Host "Vite: Running" -ForegroundColor Green
        Write-Host $viteStatus
    } else {
        Write-Host "WARNING! Vite: Check logs with: docker-compose --env-file $EnvFile logs node" -ForegroundColor Yellow
    }
Write-Host ""

# Test 5: Ports
Write-Host "$TEST_TAG Testing Open Ports..." -ForegroundColor Yellow
foreach ($port in $ports) {
    $connection = New-Object System.Net.Sockets.TcpClient
    try {
        $connection.Connect("127.0.0.1", $port.port)
        Write-Host "$($port.name): Listening on port $($port.port) - OK" -ForegroundColor Green
    } catch {
        Write-Host "$($port.name): NOT listening on port $($port.port) - ERROR" -ForegroundColor Red
    } finally {
        $connection.Close()
    }
}

# Summary
writeHeader -Text "Test completed"
