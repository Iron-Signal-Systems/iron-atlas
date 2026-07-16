[CmdletBinding()]
param()
$ErrorActionPreference = "Stop"
$repoRoot = (& git rev-parse --show-toplevel 2>$null)
if (-not $repoRoot) { throw "Not in a Git work tree." }
$pythonPath = $env:ISRAS_PYTHON
if (-not $pythonPath) {
    $python = Get-Command python3 -ErrorAction SilentlyContinue
    if (-not $python) { $python = Get-Command python -ErrorAction SilentlyContinue }
    if (-not $python) { throw "Python 3 is required by the baseline portable validator." }
    $pythonPath = $python.Source
}
$goToolsBin = $env:ISRAS_GO_TOOLS_BIN
if (-not $goToolsBin) { $goToolsBin = Join-Path $repoRoot ".isras-go-tools/bin" }
$env:PATH = "$goToolsBin$([System.IO.Path]::PathSeparator)$env:PATH"
& $pythonPath "$repoRoot/tools/environment/go_tools.py" verify --repo-root "$repoRoot"
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
& $pythonPath "$repoRoot/tools/isras/doctor.py" --repo-root "$repoRoot" --profile portable
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
& $pythonPath "$repoRoot/tools/isras/validate_policy.py" --repo-root "$repoRoot"
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
& $pythonPath "$repoRoot/tools/isras/portable_project_checks.py" --repo-root "$repoRoot"
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
Write-Host ""
Write-Host "Portable validation PASSED."
