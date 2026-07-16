[CmdletBinding()]
param(
    [string]$VenvPath
)
$ErrorActionPreference = "Stop"
$repoRoot = (& git rev-parse --show-toplevel 2>$null)
if (-not $repoRoot) { throw "Not in a Git work tree." }
if (-not $VenvPath) { $VenvPath = Join-Path $repoRoot ".isras-tools-venv" }
$python = Get-Command python3 -ErrorAction SilentlyContinue
if (-not $python) { $python = Get-Command python -ErrorAction SilentlyContinue }
if (-not $python) { throw "Python 3 is required." }
& $python.Source -m venv $VenvPath
$venvPython = Join-Path $VenvPath "Scripts/python.exe"
& $venvPython -m pip install --upgrade pip
& $venvPython -m pip install -r (Join-Path $repoRoot "tools/requirements.txt")
& $venvPython (Join-Path $repoRoot "tools/environment/go_tools.py") bootstrap --repo-root $repoRoot
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
Write-Host "ISRAS tool environment created at $VenvPath"
Write-Host "Set ISRAS_PYTHON=$venvPython to use it."
Write-Host "Set ISRAS_GO_TOOLS_BIN=$(Join-Path $repoRoot '.isras-go-tools/bin') to use pinned Go tools."
