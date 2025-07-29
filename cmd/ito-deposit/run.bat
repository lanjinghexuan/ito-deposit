@echo off
echo Starting ito-deposit from cmd directory...
echo Current directory: %CD%

REM 检查配置文件是否存在
if exist "..\..\configs\config.yaml" (
    echo Config file found: ..\..\configs\config.yaml
    ito-deposit.exe -conf ../../configs
) else (
    echo ERROR: Config file not found: ..\..\configs\config.yaml
    echo Please run from project root directory instead
)
pause