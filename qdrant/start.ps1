param()

$storagePath = Join-Path $PSScriptRoot "data"
if (-not (Test-Path $storagePath)) { New-Item -ItemType Directory -Path $storagePath -Force | Out-Null }

$env:QDRANT__STORAGE__STORAGE_PATH = $storagePath

$logFile = Join-Path $PSScriptRoot "qdrant.log"
$errFile = Join-Path $PSScriptRoot "qdrant.err"
$exe = Join-Path $PSScriptRoot "qdrant.exe"

Write-Host "Starting Qdrant (storage: $storagePath)..."

$proc = Start-Process -WindowStyle Hidden -FilePath $exe -PassThru
Write-Host "Qdrant started, PID: $($proc.Id)"
