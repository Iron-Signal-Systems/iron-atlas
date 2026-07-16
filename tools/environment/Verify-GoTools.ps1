[CmdletBinding()]
param()
$ErrorActionPreference = "Stop"
$repoRoot = (& git rev-parse --show-toplevel 2>$null)
if (-not $repoRoot) { throw "Not in a Git work tree." }
$pythonPath = $env:ISRAS_PYTHON
if (-not $pythonPath) {
    $python = Get-Command python3 -ErrorAction SilentlyContinue
    if (-not $python) { $python = Get-Command python -ErrorAction SilentlyContinue }
    if (-not $python) { throw "Python 3 is required." }
    $pythonPath = $python.Source
}
& $pythonPath "$repoRoot/tools/environment/go_tools.py" verify --repo-root "$repoRoot"
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
