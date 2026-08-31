@echo off
set QDRANT__STORAGE__STORAGE_PATH=%~dp0data
start /B "" "%~dp0qdrant.exe"
echo Qdrant starting on port 6333...
echo Storage: %~dp0data
