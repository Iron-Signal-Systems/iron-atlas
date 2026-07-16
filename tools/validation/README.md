# Validation Entrypoints

- `validate_current.sh` — current checkout policy and portable validation
- `validate_portable.sh` / `Validate-Portable.ps1` — portable project checks
- `validate_fresh_clone.sh` — canonical remote completeness
- `validate_checkpoint.sh` — isolated historical checkpoint
- `validate_canonical.sh` — project-specific canonical environment

The bootstrap portable validator detects common project types. Replace or extend
it with explicit project checks before formal repository-assurance acceptance.

## Required bootstrap

Before portable validation:

```bash
chmod +x tools/environment/bootstrap_tools.sh
./tools/environment/bootstrap_tools.sh

export ISRAS_PYTHON="$PWD/.isras-tools-venv/bin/python"
export ISRAS_GO_TOOLS_BIN="$PWD/.isras-go-tools/bin"
export PATH="$ISRAS_GO_TOOLS_BIN:$PATH"
```

Portable validation verifies the exact pinned `govulncheck` binary and does not silently accept an unrelated ambient installation.

## Manifest regeneration

After changing repository files, regenerate both repository manifests:

```bash
python3 tools/validation/regenerate_manifests.py
git add FILE-MANIFEST.txt SOURCE-SHA256SUMS.txt
```

`FILE-MANIFEST.txt` records the complete tracked and non-ignored working file
set. `SOURCE-SHA256SUMS.txt` uses the adopted ISRAS generator and hashes every
tracked source file except the checksum manifest itself. This includes and
cryptographically binds `FILE-MANIFEST.txt`.

Manifest regeneration must be the final content-changing operation before
validation and candidate review.

## Fresh-clone tool bootstrap

Fresh-clone validation checks out the exact candidate commit from the canonical
remote and then bootstraps the pinned Python and Go assurance tools inside that
temporary clone. It does not depend on an ambient workstation installation.

The temporary clone sets clone-local `ISRAS_PYTHON`,
`ISRAS_GO_TOOLS_BIN`, and `PATH` values before portable validation. The
repository must remain clean because `.isras-tools-venv/` and
`.isras-go-tools/` are ignored repository-local tool directories.
