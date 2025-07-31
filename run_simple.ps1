# Simple run script for ito-deposit
Write-Host "Simple run script for ito-deposit" -ForegroundColor Green
Write-Host "Current directory: $(Get-Location)" -ForegroundColor Yellow

# 确保在项目根目录
if (-not (Test-Path "configs/config.yaml")) {
    Write-Host "ERROR: Must run from project root directory" -ForegroundColor Red
    Write-Host "Current directory: $(Get-Location)" -ForegroundColor Red
    Write-Host "Expected to find: configs/config.yaml" -ForegroundColor Red
    Read-Host "Press Enter to exit"
    exit 1
}

Write-Host "Config file found: configs/config.yaml" -ForegroundColor Green

# 编译程序
Write-Host "Building..." -ForegroundColor Green
go build -o bin/ito-deposit.exe ./cmd/ito-deposit
if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed" -ForegroundColor Red
    Read-Host "Press Enter to exit"
    exit 1
}

# 运行程序
Write-Host "Running..." -ForegroundColor Green
& "./bin/ito-deposit.exe"