@echo off
echo Starting ito-deposit service...
echo Current directory: %CD%
cd /d "%~dp0"
echo Changed to project root: %CD%

REM 检查配置文件是否存在
if exist "configs\config.yaml" (
    echo Config file found: configs\config.yaml
) else (
    echo ERROR: Config file not found: configs\config.yaml
    pause
    exit /b 1
)

REM 编译程序
echo Building application...
go build -o bin\ito-deposit.exe .\cmd\ito-deposit

REM 运行服务
echo Starting service...
bin\ito-deposit.exe
pause