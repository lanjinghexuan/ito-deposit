# PowerShell script to start ito-deposit service
Write-Host "Starting ito-deposit service..." -ForegroundColor Green
Write-Host "Current directory: $(Get-Location)" -ForegroundColor Yellow

# Change to project root directory
$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $scriptPath
Write-Host "Changed to project root: $(Get-Location)" -ForegroundColor Yellow

# Check if configs directory exists
if (Test-Path "./configs") {
    Write-Host "Config directory found: ./configs" -ForegroundColor Green
    if (Test-Path "./configs/config.yaml") {
        Write-Host "Config file found: ./configs/config.yaml" -ForegroundColor Green
    } else {
        Write-Host "Config file NOT found: ./configs/config.yaml" -ForegroundColor Red
    }
} else {
    Write-Host "Config directory NOT found: ./configs" -ForegroundColor Red
}

# Start the service
Write-Host "Starting kratos run..." -ForegroundColor Green
kratos run