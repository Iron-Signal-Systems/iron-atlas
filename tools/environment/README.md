# Environment Profiles and Pinned Tools

Environment profiles declare capabilities rather than relying on machine memory.

## Bootstrap

Linux and macOS:

```bash
chmod +x tools/environment/bootstrap_tools.sh
./tools/environment/bootstrap_tools.sh

export ISRAS_PYTHON="$PWD/.isras-tools-venv/bin/python"
export ISRAS_GO_TOOLS_BIN="$PWD/.isras-go-tools/bin"
export PATH="$ISRAS_GO_TOOLS_BIN:$PATH"
```

Windows PowerShell:

```powershell
.\tools\environment\Bootstrap-Tools.ps1
$env:ISRAS_PYTHON = Join-Path $PWD ".isras-tools-venv/Scripts/python.exe"
$env:ISRAS_GO_TOOLS_BIN = Join-Path $PWD ".isras-go-tools/bin"
$env:PATH = "$env:ISRAS_GO_TOOLS_BIN$([IO.Path]::PathSeparator)$env:PATH"
```

Python packages are pinned in `tools/requirements.txt`. Go assurance tools are pinned in `tools/go-tools.lock.json`.

## Go-tool integrity

`tools/environment/go_tools.py` verifies the module version and both Go module checksums before installation. It then verifies the installed binary's embedded package path, module version, and checksum.

Wrappers:

- `bootstrap_go_tools.sh` / `Bootstrap-GoTools.ps1`
- `verify_go_tools.sh` / `Verify-GoTools.ps1`

The repository-local `.isras-go-tools/` directory is ignored and must not be committed.

## Environment doctor

```bash
./tools/environment/doctor.sh portable
```

or:

```powershell
.\tools\environment\Doctor.ps1 portable
```

Canonical and specialized profiles must be customized before their results are used for acceptance claims.

## Pinned Go toolchain

Iron Atlas retains Go language compatibility at `1.25.0` while
declaring `go1.26.5` as the preferred build and validation toolchain in
`go.mod`.

The exact toolchain is security-relevant because `govulncheck` evaluates
reachable vulnerabilities in the standard library supplied by the selected Go
toolchain. Hosted validation sets `GOTOOLCHAIN=go1.26.5` and the Go
command downloads and checksum-verifies that official toolchain when it is not
already available on the host.

The portable environment profile and project toolchain requirements reject Go
toolchains older than `1.26.5`. A host-specific patched build may
include a suffix after `go1.26.5` but must report the same base release.
