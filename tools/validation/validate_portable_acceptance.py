#!/usr/bin/env python3
from pathlib import Path
import json
import sys

root=Path(__file__).resolve().parents[2]
errors=[]
required=[
  "validation/toolchain-requirements.json",
  "validation/evidence/README.md",
  "tools/validation/validate_toolchain.py",
  "tools/validation/redact_validation_text.py",
  "tools/validation/validate_committed_evidence.py",
  "tools/validation/record_validation_evidence.sh",
  "tools/validation/verify_canonical_clone.sh",
  "tools/isras/validate_fresh_clone.py",
  "docs/architecture/PORTABLE-VALIDATION-AND-CANONICAL-REPOSITORY-ACCEPTANCE.md",
  "docs/decisions/ADR-0005-CANONICAL-REPOSITORY-REPRODUCIBILITY.md",
  "docs/operations/CANONICAL-CLEAN-CLONE-VALIDATION.md",
]
for name in required:
    if not (root/name).is_file(): errors.append(f"missing portability artifact: {name}")
try:
    req=json.loads((root/"validation/toolchain-requirements.json").read_text())
    if req.get("canonical_repository") != "https://github.com/Iron-Signal-Systems/iron-atlas.git":
        errors.append("toolchain contract has noncanonical repository")
except Exception as exc:
    errors.append(f"invalid toolchain requirements: {exc}")
ignore=(root/".gitignore").read_text()
for token in ["!/validation/evidence/", "!/validation/evidence/**"]:
    if token not in ignore: errors.append(f".gitignore missing evidence exception: {token}")
rec=(root/"tools/validation/record_validation_evidence.sh").read_text()
for token in ["mktemp -d", "redact_validation_text.py", "validate_committed_evidence.py", "git status --porcelain"]:
    if token not in rec: errors.append(f"evidence recorder missing token: {token}")
verify=(root/"tools/validation/verify_canonical_clone.sh").read_text()
for token in ["git ls-remote", "git clone", "fetch --quiet --force --tags", "validate_toolchain.py", "go mod verify"]:
    if token not in verify: errors.append(f"canonical clone verifier missing token: {token}")
fresh=(root/"tools/isras/validate_fresh_clone.py").read_text()
for token in [
    "bootstrap_tools.sh",
    "Bootstrap-Tools.ps1",
    "ISRAS_PYTHON",
    "ISRAS_GO_TOOLS_BIN",
    "Fresh clone remains clean after ignored tool bootstrap",
]:
    if token not in fresh:
        errors.append(f"fresh-clone validator missing token: {token}")
template=(root/"docs/acceptance/PHASE-1-STEP-2-ACCEPTANCE-RECORD-TEMPLATE.md").read_text()
for token in ["Canonical clean-clone validation", "Committed evidence path", "Toolchain requirements SHA-256"]:
    if token not in template: errors.append(f"Step 2 acceptance template missing token: {token}")
if "/tmp/iron-atlas" in template: errors.append("acceptance template must not retain workstation-local transcript paths")
readme=(root/"test-framework/README.md").read_text()
if "remain excluded from Git" in readme: errors.append("test framework still says all generated results remain excluded from Git")
if errors:
    for e in errors: print(f"FAIL: {e}", file=sys.stderr)
    raise SystemExit(1)
print("validated portable validation and canonical-repository acceptance contract")
