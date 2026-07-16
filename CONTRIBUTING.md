# Contributing

## Development workflow

1. Start from the current canonical `dev` branch.
2. Create a purpose-named work branch.
3. Submit a pull request back to `dev`.
4. Do not push material development directly to `dev` or `main`.
5. Keep modules replaceable and dependency direction explicit.

## Required change set

A material change includes all applicable requirements, architecture,
implementation, tests, hostile cases, fixtures, expected outcomes, validation,
environment changes, synchronized documentation, retained evidence, warnings,
limitations, non-claims, and approval records.

Documentation is part of the same change set, not follow-up cleanup.

## Validation

Before review, bootstrap the repository-owned tools and run:

```bash
chmod +x tools/environment/bootstrap_tools.sh
./tools/environment/bootstrap_tools.sh

export ISRAS_PYTHON="$PWD/.isras-tools-venv/bin/python"
export ISRAS_GO_TOOLS_BIN="$PWD/.isras-go-tools/bin"
export PATH="$ISRAS_GO_TOOLS_BIN:$PATH"

python3 tools/validation/validate_toolchain.py
chmod +x test-framework/run_all.sh
./test-framework/run_all.sh
./tools/validation/validate_portable.sh
```

Run the applicable Iron Atlas phase gate for bounded implementation or
acceptance work.

On Windows, the ISRAS portable entrypoint is:

```powershell
.\tools\validation\Validate-Portable.ps1
```

Before formal acceptance, perform applicable fresh-clone, canonical,
specialized, and historical predecessor validation. Acceptance requires the
exact pushed commit to pass from a clean canonical GitHub clone.

## Security and evidence

- Do not commit credentials, raw infrastructure evidence, production data,
  unrestricted logs, private keys, or credential-bearing connection strings.
- Commit only sanitized, checksummed evidence in approved repository paths.
- Do not rely on workstation-only scripts, services, logs, or undeclared
  packages.
- Governed production-impacting changes require the applicable independent
  approval.
