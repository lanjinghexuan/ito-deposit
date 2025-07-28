@echo off
echo Starting ito-deposit service...
echo Current directory: %CD%
cd /d "%~dp0"
echo Changed to project root: %CD%
kratos run
pause