$ErrorActionPreference = "Stop"

$root = Split-Path -Parent $MyInvocation.MyCommand.Path
$go = "C:\Program Files\Go\bin\go.exe"
$tmp = Join-Path $root "tmp"

function Write-Step($message) {
  Write-Host ""
  Write-Host "==> $message" -ForegroundColor Green
}

function Stop-DemoProcesses() {
  Get-Process user-service, gamification, bff, node, python, py -ErrorAction SilentlyContinue |
  Stop-Process -Force
}

function Wait-ForHttp($url, $label, $timeoutSeconds = 20) {
  for ($attempt = 0; $attempt -lt $timeoutSeconds; $attempt++) {
    try {
      $response = Invoke-WebRequest $url -UseBasicParsing -TimeoutSec 3
      if ($response.StatusCode -ge 200 -and $response.StatusCode -lt 500) {
        Write-Host "$label ready: $url" -ForegroundColor DarkGreen
        return $true
      }
    }
    catch {
      Start-Sleep -Milliseconds 900
    }
  }

  return $false
}

function Show-LogTail($path, $label) {
  if (Test-Path $path) {
    Write-Host ""
    Write-Host "Last lines from $label ($path)" -ForegroundColor Yellow
    Get-Content $path -Tail 25
  }
}

function Build-Binaries() {
  Write-Step "Building user-service"
  Push-Location (Join-Path $root "services\user-service")
  & $go build -o (Join-Path $tmp "user-service.exe") ./cmd/user-service
  Pop-Location

  Write-Step "Building gamification-service"
  Push-Location (Join-Path $root "services\gamification")
  & $go build -o (Join-Path $tmp "gamification.exe") ./cmd/gamification
  Pop-Location

  Write-Step "Building BFF"
  Push-Location (Join-Path $root "bff")
  & $go build -o (Join-Path $tmp "bff.exe") ./cmd/bff
  Pop-Location
}

function Start-Backend($postgresHost, $postgresPort, $postgresDb, $postgresUser, $postgresPassword) {
  Stop-DemoProcesses

  Write-Host "Trying PostgreSQL DSN: $postgresUser@${postgresHost}:$postgresPort/$postgresDb" -ForegroundColor Cyan

  $userProcess = Start-Process -PassThru -FilePath "cmd.exe" -ArgumentList "/c", "set USER_SERVICE_ADDRESS=:15051&&set POSTGRES_HOST=$postgresHost&&set POSTGRES_PORT=$postgresPort&&set POSTGRES_DB=$postgresDb&&set POSTGRES_USER=$postgresUser&&set POSTGRES_PASSWORD=$postgresPassword&&`"$tmp\user-service.exe`"" -RedirectStandardOutput (Join-Path $tmp "user-service.log") -RedirectStandardError (Join-Path $tmp "user-service.err.log") -WindowStyle Hidden
  Start-Sleep -Milliseconds 700

  $gamificationProcess = Start-Process -PassThru -FilePath "cmd.exe" -ArgumentList "/c", "set GAMIFICATION_SERVICE_ADDRESS=:15052&&set POSTGRES_HOST=$postgresHost&&set POSTGRES_PORT=$postgresPort&&set POSTGRES_DB=$postgresDb&&set POSTGRES_USER=$postgresUser&&set POSTGRES_PASSWORD=$postgresPassword&&`"$tmp\gamification.exe`"" -RedirectStandardOutput (Join-Path $tmp "gamification.log") -RedirectStandardError (Join-Path $tmp "gamification.err.log") -WindowStyle Hidden
  Start-Sleep -Milliseconds 700

  $bffProcess = Start-Process -PassThru -FilePath "cmd.exe" -ArgumentList "/c", "set USER_SERVICE_GRPC=127.0.0.1:15051&&set GAMIFICATION_SERVICE_GRPC=127.0.0.1:15052&&set BFF_ADDRESS=:18080&&set BFF_METRICS_ADDRESS=:19090&&set POSTGRES_HOST=$postgresHost&&set POSTGRES_PORT=$postgresPort&&set POSTGRES_DB=$postgresDb&&set POSTGRES_USER=$postgresUser&&set POSTGRES_PASSWORD=$postgresPassword&&`"$tmp\bff.exe`"" -RedirectStandardOutput (Join-Path $tmp "bff.log") -RedirectStandardError (Join-Path $tmp "bff.err.log") -WindowStyle Hidden

  if (Wait-ForHttp "http://127.0.0.1:18080/healthz" "BFF" 15) {
    return @{
      User = $userProcess
      Game = $gamificationProcess
      Bff  = $bffProcess
      Db   = $postgresDb
      PgUser = $postgresUser
    }
  }

  Show-LogTail (Join-Path $tmp "user-service.log") "user-service.log"
  Show-LogTail (Join-Path $tmp "gamification.log") "gamification.log"
  Show-LogTail (Join-Path $tmp "bff.log") "bff.log"
  Stop-DemoProcesses
  return $null
}

function Build-Frontend() {
  Write-Step "Building frontend"
  Push-Location (Join-Path $root "frontend")
  cmd /c "set VITE_API_BASE_URL=http://127.0.0.1:18080&&npm run build"
  Pop-Location
}

function Start-Frontend() {
  Write-Step "Starting frontend static server"
  $frontendProcess = Start-Process -PassThru -FilePath "cmd.exe" -ArgumentList "/c", "node `"$root\frontend\server.cjs`"" -RedirectStandardOutput (Join-Path $tmp "frontend.log") -RedirectStandardError (Join-Path $tmp "frontend.err.log") -WindowStyle Hidden
  if (-not (Wait-ForHttp "http://127.0.0.1:4173" "Frontend" 15)) {
    Show-LogTail (Join-Path $tmp "frontend.err.log") "frontend.err.log"
    throw "Frontend did not start on http://127.0.0.1:4173"
  }
  return $frontendProcess
}

New-Item -ItemType Directory -Force -Path $tmp | Out-Null

Write-Step "Stopping old processes"
Stop-DemoProcesses

Build-Binaries

$candidates = @()

if ($env:POSTGRES_HOST -and $env:POSTGRES_PORT -and $env:POSTGRES_DB -and $env:POSTGRES_USER) {
  $candidates += @{
    Host = $env:POSTGRES_HOST
    Port = $env:POSTGRES_PORT
    Db   = $env:POSTGRES_DB
    User = $env:POSTGRES_USER
    Pass = $env:POSTGRES_PASSWORD
  }
}

$candidates += @(
  @{ Host = "127.0.0.1"; Port = "5432"; Db = "postgres";       User = "postgres"; Pass = "1234" },
  @{ Host = "127.0.0.1"; Port = "5432"; Db = "cdek_platform";  User = "cdek";     Pass = "cdek" },
  @{ Host = "127.0.0.1"; Port = "5432"; Db = "postgres";       User = "postgres"; Pass = "" }
)

$backend = $null
foreach ($candidate in $candidates) {
  $backend = Start-Backend $candidate.Host $candidate.Port $candidate.Db $candidate.User $candidate.Pass
  if ($backend -ne $null) {
    break
  }
}

if ($backend -eq $null) {
  throw "Could not start backend with available PostgreSQL credentials. Set POSTGRES_* env vars and run script again."
}

$frontendProcess = Start-Frontend

Write-Host ""
Write-Host "Project started successfully" -ForegroundColor Cyan
Write-Host "Frontend: http://127.0.0.1:4173"
Write-Host "BFF: http://127.0.0.1:18080"
Write-Host "Metrics: http://127.0.0.1:19090"
Write-Host "Backend DB user: $($backend.PgUser), database: $($backend.Db)"
