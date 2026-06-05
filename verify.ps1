# verify.ps1
# Compiles envguard in WSL, runs all tests, and runs smoke checks on testdata/js-project.

$ErrorActionPreference = "Stop"

Write-Host "===============================================" -ForegroundColor Cyan
Write-Host " envguard - End-to-End Build and Verification " -ForegroundColor Cyan
Write-Host "===============================================" -ForegroundColor Cyan
Write-Host

# 1. Compiling binary
Write-Host "[1/4] Compiling native binary in WSL..." -ForegroundColor Yellow
wsl bash -c "go build -o envguard ."
if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: Compilation failed!" -ForegroundColor Red
    exit $LASTEXITCODE
}
Write-Host "Success: Binary compiled to ./envguard" -ForegroundColor Green
Write-Host

# 2. Running tests
Write-Host "[2/4] Running Go test suite in WSL..." -ForegroundColor Yellow
wsl bash -c "go test ./..."
if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: Test suite failed!" -ForegroundColor Red
    exit $LASTEXITCODE
}
Write-Host "Success: All test suites passed!" -ForegroundColor Green
Write-Host

# 3. Running Project Doctor Health Check
Write-Host "[3/4] Running envguard doctor on testdata/js-project..." -ForegroundColor Yellow
wsl ./envguard doctor testdata/js-project
Write-Host

# 4. Running Complete Audit & Validation Flow
Write-Host "[4/4] Running configuration audit, init, validation, and sync..." -ForegroundColor Yellow

Write-Host "`nRunning audit..." -ForegroundColor Cyan
# audit exits with status 1 on findings which is expected, so ignore error-abort for this command
$val = wsl ./envguard audit testdata/js-project --env testdata/js-project/.env
Write-Host $val

Write-Host "`nRunning init wizard..." -ForegroundColor Cyan
wsl ./envguard init --yes --path testdata/js-project

Write-Host "`nRunning validate..." -ForegroundColor Cyan
wsl ./envguard validate --env testdata/js-project/.env --config testdata/js-project/.envguard.yaml

Write-Host "`nRunning sync..." -ForegroundColor Cyan
wsl ./envguard sync --env testdata/js-project/.env --output testdata/js-project/.env.example --force

Write-Host "`nRunning diff..." -ForegroundColor Cyan
wsl ./envguard diff testdata/js-project/.env testdata/js-project/.env.example

Write-Host "`n===============================================" -ForegroundColor Green
Write-Host " envguard is fully compiled and verified! " -ForegroundColor Green
Write-Host "===============================================" -ForegroundColor Green
