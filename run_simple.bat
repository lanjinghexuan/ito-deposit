@echo off
echo Simple run script for ito-deposit
echo Current directory: %CD%

REM 确保在项目根目录
if not exist "configs\config.yaml" (
    echo ERROR: Must run from project root directory
    echo Current directory: %CD%
    echo Expected to find: configs\config.yaml
    pause
    exit /b 1
)

echo Config file found: configs\config.yaml

REM 编译程序
echo Building...
go build -o bin\ito-deposit.exe .\cmd\ito-deposit
if errorlevel 1 (
    echo Build failed
    pause
    exit /b 1
)

REM 运行程序
echo Running...
bin\ito-deposit.exe
pause