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
gomod=(root/"go.mod").read_text()
for token in ["go 1.25.0", "toolchain go1.26.5"]:
    if token not in gomod:
        errors.append(f"go.mod missing pinned toolchain token: {token}")

workflow=(root/".github/workflows/portable-validation.yml").read_text()
for token in ["GOTOOLCHAIN: go1.26.5", '"go1.26.5"']:
    if token not in workflow:
        errors.append(f"portable workflow missing toolchain token: {token}")

profile=json.loads(
    (root/"tools/environment/profiles/portable.json").read_text()
)
go_specs=[
    item for item in profile.get("required_commands", [])
    if item.get("command") == "go"
]
if len(go_specs) != 1 or "1\\.26\\.5" not in go_specs[0].get(
    "version_pattern", ""
):
    errors.append("portable profile does not require Go 1.26.5")

requirements=json.loads(
    (root/"validation/toolchain-requirements.json").read_text()
)
go_requirement=requirements.get("requirements", {}).get("go", {})
if go_requirement.get("minimum") != "1.26.5":
    errors.append("toolchain requirements do not require Go 1.26.5")
if go_requirement.get("preferred_toolchain") != "go1.26.5":
    errors.append("preferred Go toolchain is not go1.26.5")
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
