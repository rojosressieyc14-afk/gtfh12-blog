@echo off
cd /d "%~dp0"
powershell.exe -ExecutionPolicy Bypass -File "start.ps1"
pause
