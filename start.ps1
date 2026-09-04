param([switch]$NoAdmin)

$Root = Split-Path -Parent $MyInvocation.MyCommand.Path

# Read DB password from environment variable; prompt if not set
$mysqlPass = $env:PULSEBLOG_DB_PASSWORD
if (-not $mysqlPass) {
    $mysqlPass = Read-Host 'Enter MySQL root password'
    if (-not $mysqlPass) {
        Write-Host '[FAIL] MySQL password is required' -ForegroundColor Red
        exit 1
    }
}

Write-Host '=== PulseBlog Startup ===' -ForegroundColor Cyan

# 1. kill old processes on target ports
Write-Host '[..] releasing ports...' -ForegroundColor Yellow
$ports = @(8080, 5173, 5174)
$used = netstat -ano | Select-String ':8080|:5173|:5174' | Select-String 'LISTENING'
foreach ($line in $used) {
    if ($line -match 'LISTENING\s+(\d+)$') {
        Stop-Process -Id $Matches[1] -Force -ErrorAction SilentlyContinue
    }
}
Start-Sleep -Seconds 1
Write-Host '[OK] ports released' -ForegroundColor Green

# 2. check MySQL
$mysql = Get-Service MySQL80 -ErrorAction SilentlyContinue
if (-not $mysql -or $mysql.Status -ne 'Running') {
    Write-Host '[..] starting MySQL80...' -ForegroundColor Yellow
    Start-Service MySQL80 -ErrorAction Stop
    Start-Sleep -Seconds 3
}
Write-Host '[OK] MySQL80 running' -ForegroundColor Green

# 2b. ensure Redis
$redisProc = Get-Process redis-server -ErrorAction SilentlyContinue
if (-not $redisProc) {
    Write-Host '[..] starting Redis...' -ForegroundColor Yellow
    $redisExe = Get-Command redis-server -ErrorAction SilentlyContinue
    if ($redisExe) {
        $null = Start-Process -WindowStyle Hidden -FilePath 'redis-server' -PassThru
        Start-Sleep -Seconds 2
        Write-Host '[OK] Redis started' -ForegroundColor Green
    } else {
        Write-Host '[SKIP] redis-server not found in PATH (install or use docker-compose up -d)' -ForegroundColor DarkYellow
    }
} else {
    Write-Host '[OK] Redis already running' -ForegroundColor Green
}

# 3. ensure Qdrant
$qdrantExe = Join-Path $Root 'qdrant\qdrant.exe'
$qdrantData = Join-Path $Root 'qdrant\data'
if (Test-Path $qdrantExe) {
    $qdrantProc = Get-Process qdrant -ErrorAction SilentlyContinue
    if (-not $qdrantProc) {
        Write-Host '[..] starting Qdrant...' -ForegroundColor Yellow
        if (-not (Test-Path $qdrantData)) { New-Item -ItemType Directory -Path $qdrantData -Force | Out-Null }
        $env:QDRANT__STORAGE__STORAGE_PATH = $qdrantData
        $null = Start-Process -WindowStyle Hidden -FilePath $qdrantExe -PassThru
        Start-Sleep -Seconds 3
        Write-Host '[OK] Qdrant started' -ForegroundColor Green
    } else {
        Write-Host '[OK] Qdrant already running' -ForegroundColor Green
    }
} else {
    Write-Host '[SKIP] Qdrant binary not found at qdrant\qdrant.exe' -ForegroundColor DarkYellow
}

# 4. ensure .env
$envFile = Join-Path $Root '.env'
if (-not (Test-Path $envFile)) {
    Write-Host '[..] creating .env from environment or prompts...' -ForegroundColor Yellow
    $jwtSecret = $env:PULSEBLOG_JWT_SECRET
    if (-not $jwtSecret) { $jwtSecret = [guid]::NewGuid().ToString('N') + [guid]::NewGuid().ToString('N') }
    $adminUser = $env:PULSEBLOG_ADMIN_USERNAME
    if (-not $adminUser) { $adminUser = 'admin' }
    $adminPass = $env:PULSEBLOG_ADMIN_PASSWORD
    if (-not $adminPass) { $adminPass = Read-Host 'Enter default admin password' }
    @"
DB_PASSWORD=$mysqlPass
DB_NAME=blog_system
JWT_SECRET=$jwtSecret
DEFAULT_ADMIN_USERNAME=$adminUser
DEFAULT_ADMIN_PASSWORD=$adminPass
WEB_ORIGIN=http://localhost:5173
ADMIN_ORIGIN=http://localhost:5174
GIN_MODE=debug
UPLOAD_DIR=./uploads
DEEPSEEK_API_KEY=
DEEPSEEK_BASE_URL=https://api.deepseek.com/v1
QDRANT_ADDR=http://localhost:6333
QDRANT_API_KEY=
VITE_ADMIN_URL=http://localhost:5174
API_ENCRYPTION_KEY=
SMTP_HOST=
SMTP_PORT=587
SMTP_USER=
SMTP_PASSWORD=
SMTP_FROM=
"@ | Set-Content -Path $envFile
}
Write-Host '[OK] .env ready' -ForegroundColor Green

# 5. start backend (new window, port 8080)
Write-Host '[..] starting backend...' -ForegroundColor Yellow
# Read JWT secret from .env if not in environment
if (-not $env:JWT_SECRET) {
    $envMatch = Select-String -Path $envFile -Pattern '^JWT_SECRET=(.+)$'
    if ($envMatch) { $env:JWT_SECRET = $envMatch.Matches[0].Groups[1].Value }
}
$env:DB_PASSWORD = $mysqlPass
$env:DB_NAME = 'blog_system'
$env:GIN_MODE = 'debug'
$env:WEB_ORIGIN = 'http://localhost:5173'
$env:ADMIN_ORIGIN = 'http://localhost:5174'
$serverDir = Join-Path $Root 'server'
$backend = Start-Process -WindowStyle Normal -FilePath 'go' -ArgumentList 'run','./cmd/api' `
    -WorkingDirectory $serverDir -PassThru
Write-Host ('[OK] backend PID={0}' -f $backend.Id) -ForegroundColor Green

# 6. start frontend (new window, port 5173)
Write-Host '[..] starting frontend...' -ForegroundColor Yellow
$webDir = Join-Path $Root 'web'
$frontend = Start-Process -WindowStyle Normal -FilePath 'npx.cmd' -ArgumentList 'vite' `
    -WorkingDirectory $webDir -PassThru
Write-Host ('[OK] frontend PID={0}' -f $frontend.Id) -ForegroundColor Green

# 7. optional: admin (new window, port 5174)
if (-not $NoAdmin) {
    Write-Host '[..] starting admin...' -ForegroundColor Yellow
    $adminDir = Join-Path $Root 'admin'
    $adminProcess = Start-Process -WindowStyle Normal -FilePath 'npx.cmd' -ArgumentList 'vite' `
        -WorkingDirectory $adminDir -PassThru
    Write-Host ('[OK] admin PID={0}' -f $adminProcess.Id) -ForegroundColor Green
}

Start-Sleep -Seconds 3

# 8. verify with netstat
Write-Host "`n=== Status ===" -ForegroundColor Cyan
$listeners = netstat -ano | Select-String ':8080 |:5173 |:5174 ' | Select-String 'LISTENING'
if (Get-Process qdrant -ErrorAction SilentlyContinue) { Write-Host '[OK] Qdrant   -> http://localhost:6333' -ForegroundColor Green }
else { Write-Host '[SKIP] Qdrant   not running' -ForegroundColor DarkYellow }
if ($listeners -match ':8080 ') { Write-Host '[OK] backend  -> http://localhost:8080' -ForegroundColor Green }
else { Write-Host '[FAIL] backend  not listening' -ForegroundColor Red }

if ($listeners -match ':5173 ') { Write-Host '[OK] frontend -> http://localhost:5173/PulseBlog/' -ForegroundColor Green }
else { Write-Host '[FAIL] frontend not listening' -ForegroundColor Red }

if (-not $NoAdmin) {
    if ($listeners -match ':5174 ') { Write-Host '[OK] admin    -> http://localhost:5174' -ForegroundColor Green }
    else { Write-Host '[FAIL] admin    not listening' -ForegroundColor Red }
}
