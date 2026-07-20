#!/usr/bin/env python3
from pathlib import Path
import json
import re

ROOT = Path(__file__).resolve().parents[2]
failures: list[str] = []
passes = 0


def check(name: str, condition: bool) -> None:
    global passes
    if condition:
        passes += 1
        print(f"PASS: {name}")
    else:
        failures.append(name)
        print(f"FAIL: {name}")


def require(path: str, tokens: list[str]) -> str:
    candidate = ROOT / path
    check(f"required file {path}", candidate.is_file())
    if not candidate.is_file():
        return ""
    text = candidate.read_text(encoding="utf-8")
    compact = " ".join(text.split())
    for token in tokens:
        check(f"{path} contains {token}", " ".join(token.split()) in compact)
    return text


require(
    "tools/validation/lib/reporting.sh",
    [
        "validation_report_init",
        "validation_run",
        "validation_skip",
        "validation_report_finish",
        "validation_note",
        "FINAL RESULT: PASS",
        "PRIMARY FAILURE",
        "ADDITIONAL UNIQUE FAILURES",
        "CASCADED FAILURES",
        "SKIPPED DEPENDENT CHECKS",
        "Per-check logs",
        "FINAL RESULT: FAIL",
    ],
)

run_all = require(
    "test-framework/run_all.sh",
    [
        'source "$repo_root/tools/validation/lib/reporting.sh"',
        "isolated_gate_revalidate",
        "IRON_ATLAS_HTTP_PREDECESSOR_ALREADY_VALIDATED",
        "validation_report_finish",
        "validation reporting static validation",
        "validation reporting regression",
    ],
)
check(
    "later complete runner does not execute frozen HTTP static validator directly",
    "validate_phase1_step3_http_login_callback.py" not in run_all,
)
check(
    "later complete runner does not execute frozen HTTP regression directly",
    "test_phase1_step3_http_login_callback.sh" not in run_all,
)

require(
    "tools/validation/phase-gates/validate_phase1_step3_authenticated_session.sh",
    [
        'source "$repo_root/tools/validation/lib/reporting.sh"',
        'source "$repo_root/tools/validation/lib/isolated_gate_revalidation.sh"',
        "isolated_gate_revalidate",
        "IRON_ATLAS_HTTP_PREDECESSOR_ALREADY_VALIDATED=1",
        "validation_skip",
        "validation_report_finish",
        "HTTP login and callback checkpoint remains valid",
    ],
)

regression = require(
    "test-framework/authentication/test_phase1_step3_authenticated_session.sh",
    [
        'source "$repo_root/tools/validation/lib/reporting.sh"',
        "go tool govulncheck",
        "repository-managed govulncheck preflight",
        "validation_report_finish",
    ],
)

require(
    "test-framework/validation/test_validation_reporting.sh",
    [
        "ERROR: exact root cause",
        "Exit status: 7",
        "SKIPPED DEPENDENT CHECKS",
        "FINAL RESULT: PASS",
    ],
)

go_mod = require("go.mod", ["tool", "golang.org/x/vuln/cmd/govulncheck"])
toolchain = json.loads(
    (ROOT / "validation/toolchain-requirements.json").read_text(encoding="utf-8")
)
tools = toolchain.get("go_tools", [])
check(
    "toolchain declares pinned govulncheck",
    any(
        tool.get("command") == "govulncheck"
        and tool.get("version") == "v1.6.0"
        and tool.get("invocation") == "go tool govulncheck"
        for tool in tools
    ),
)

bare_invocations: list[str] = []
bare_govulncheck_command = re.compile(
    r"(?:^\\s*|[;&|]\\s*|\\b(?:if|then|while|until|do)\\s+)"
    r"(?:(?:command|env)\\s+)?"
    r"(?:[A-Za-z_][A-Za-z0-9_]*=[^\\s;|&]+\\s+)*"
    r"govulncheck(?=\\s|$)"
)

for path in ROOT.rglob("*.sh"):
    if ".git" in path.parts:
        continue
    for number, line in enumerate(
        path.read_text(encoding="utf-8").splitlines(), start=1
    ):
        executable_text = line.split("#", 1)[0]
        if bare_govulncheck_command.search(executable_text):
            bare_invocations.append(f"{path.relative_to(ROOT)}:{number}")
check(
    "shell validators contain no PATH-dependent bare govulncheck invocation",
    not bare_invocations,
)
if bare_invocations:
    print("Bare govulncheck invocations:")
    for item in bare_invocations:
        print(f"  {item}")

require(
    "docs/testing/VALIDATION-FAILURE-REPORTING.md",
    [
        "Primary failure",
        "Additional unique failures",
        "Cascaded failures",
        "Skipped dependent checks",
        "repository-managed Go tools",
        "isolated clone",
    ],
)

require(
    "docs/requirements/SYSTEM-REQUIREMENTS.md",
    ["IA-VAL-012", "primary failure", "retained log path"],
)

manifest = (ROOT / "FILE-MANIFEST.txt").read_text(encoding="utf-8").splitlines()
for path in (
    "tools/validation/lib/reporting.sh",
    "tools/validation/validate_validation_reporting.py",
    "test-framework/validation/test_validation_reporting.sh",
    "docs/testing/VALIDATION-FAILURE-REPORTING.md",
):
    check(f"file manifest contains {path}", path in manifest)

print()
print(f"PASS checks: {passes}")
print(f"FAIL checks: {len(failures)}")
if failures:
    print("Validation failure-reporting static validation FAILED.")
    raise SystemExit(1)

print("Validation failure-reporting static validation PASSED.")
