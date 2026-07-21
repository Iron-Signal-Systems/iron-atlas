#!/usr/bin/env python3
from pathlib import Path
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[2]
BASE = "a0ab1ad19cf48ba11d97b3a9e87acd7b68e1eb60"
CANDIDATE = "7bc37a94cfaff76976cb0bcc742a838df4564ca2"
MERGE = "5de9e1f5f9770f12b56a046dc735b769cc842a02"
BOUNDARY = "12569192da89a1a34f4ebfe107c4d02c60cbdb09"
WORKFLOW_RUNS = ["29792297812", "29792297864", "29792298042"]

errors = []
passes = 0

def norm(value: str) -> str:
    return " ".join(value.split()).lower()

def check(name: str, condition: bool) -> None:
    global passes
    if condition:
        print(f"PASS: {name}")
        passes += 1
    else:
        print(f"FAIL: {name}", file=sys.stderr)
        errors.append(name)

def read(path: str) -> str:
    candidate = ROOT / path
    check(f"required file {path}", candidate.is_file())
    return candidate.read_text(encoding="utf-8") if candidate.is_file() else ""

paths = {
    "module": "docs/architecture/MODULE-RUNTIME-AND-FAILURE-CONTAINMENT-MODEL.md",
    "schedule": "docs/architecture/SCHEDULED-EVIDENCE-INGESTION-MODEL.md",
    "freshness": "docs/architecture/MONITORING-ALERTING-AND-EVIDENCE-FRESHNESS-MODEL.md",
    "atomic": "docs/architecture/EVIDENCE-CANDIDATE-AND-ATOMIC-ACCEPTANCE-MODEL.md",
    "ifi": "docs/architecture/ATLAS-IFI-SNAPSHOT-INTEGRATION-CONTRACT.md",
    "failclosed": "docs/security/FAIL-CLOSED-AND-ADVERSARIAL-INVARIANT-MODEL.md",
    "mfa": "docs/security/MFA-AND-AUTHENTICATION-ASSURANCE-REQUIREMENTS.md",
    "signing": "docs/governance/SIGNED-CANDIDATE-AND-POST-MERGE-BOUNDARY-MODEL.md",
    "record": "docs/governance/ARCHITECTURE-AND-ROADMAP-ALIGNMENT-RECORD.md",
    "backlog": "docs/governance/POST-LICENSING-ALIGNMENT-BACKLOG.md",
    "readme": "README.md",
    "docs": "docs/README.md",
    "auth": "docs/architecture/TRUSTED-AUTHENTICATION-AND-GOVERNED-ACTOR-RESOLUTION.md",
    "trace": "docs/requirements/PHASE-1-STEP-3-REQUIREMENTS-TRACEABILITY.md",
    "roadmap": "docs/roadmap/IMPLEMENTATION-ROADMAP.md",
    "gates": "docs/roadmap/PHASE-GATE-PLAN.md",
    "acceptance": "docs/acceptance/PHASE-1-STEP-3-ACCEPTANCE-RECORD-TEMPLATE.md",
    "testing": "docs/testing/TESTING-AND-ADVERSARIAL-VALIDATION-MODEL.md",
    "target": "docs/architecture/TARGET-ARCHITECTURE.md",
    "changelog": "CHANGELOG.md",
}
files = {key: read(path) for key, path in paths.items()}

required = {
    "module": ["bounded queues", "last accepted state", "one-adapter failure"],
    "schedule": ["governed schedule", "concurrency key", "backoff", "jitter"],
    "freshness": ["service health", "service readiness", "evidence freshness", "unknown"],
    "atomic": ["candidate", "atomic publication", "stale candidates"],
    "ifi": ["ifi authoritative state", "query the ifi postgresql database directly", "write back to ifi"],
    "failclosed": ["fails closed", "unknown is represented as unknown", "enforcement points"],
    "mfa": [
        "approved external openid connect identity provider",
        "does not store user passwords",
        "does not store provider totp seeds",
        "representative-provider compatibility",
        "break-glass",
    ],
    "signing": [
        "ssh-signed candidate commit",
        "github merge commit",
        "ssh-signed empty post-merge boundary commit",
    ],
    "record": [
        BASE,
        CANDIDATE,
        MERGE,
        BOUNDARY,
        "accepted documentation and governance boundary",
        "representative external-provider compatibility",
        *WORKFLOW_RUNS,
    ],
}

for key, fragments in required.items():
    flat = norm(files[key])
    for fragment in fragments:
        check(f"{key} contains {fragment}", norm(fragment) in flat)

for key in ["readme", "auth", "trace", "roadmap", "gates", "acceptance"]:
    flat = norm(files[key])
    prohibited = [
        "validate_phase1_step3_totp_enrollment_verification_recovery.sh",
        "atlas-local totp is a distinct planned gate",
        "encrypted authenticator service",
        "local totp checkpoint, when included",
    ]
    check(
        f"{key} has no Atlas-local TOTP implementation commitment",
        not any(token in flat for token in prohibited),
    )

check(
    "README identifies accepted architecture alignment",
    BOUNDARY in files["readme"]
    and "architecture and roadmap alignment is accepted" in norm(files["readme"]),
)
check(
    "documentation index identifies accepted alignment",
    BOUNDARY in files["docs"]
    and "representative-provider compatibility" in norm(files["docs"]),
)
check(
    "backlog marked implemented and accepted",
    "status: implemented and accepted" in norm(files["backlog"])
    and BOUNDARY in files["backlog"],
)
check(
    "roadmap records accepted alignment and next checkpoint",
    BOUNDARY in files["roadmap"]
    and "representative-provider compatibility is the next planned bounded"
    in norm(files["roadmap"]),
)
check(
    "roadmap preserves assurance as implementation checkpoint",
    "authentication assurance remains an implementation checkpoint"
    in norm(files["roadmap"]),
)
check(
    "phase gate plan keeps assurance validator",
    "validate_phase1_step3_authentication_assurance.sh" in files["gates"],
)
check(
    "phase gate plan identifies provider compatibility next",
    "validate_phase1_step3_representative_provider_compatibility.sh"
    in files["gates"]
    and BOUNDARY in files["gates"],
)
check(
    "traceability uses representative-provider IA-AUTH-019",
    "IA-AUTH-019" in files["trace"]
    and "representative-provider compatibility" in norm(files["trace"]),
)
check(
    "acceptance template uses provider-owned MFA",
    "Provider-owned MFA assurance" in files["acceptance"],
)
check(
    "testing model adds alignment coverage",
    "Architecture-alignment adversarial coverage" in files["testing"],
)
check(
    "target architecture adds alignment contracts",
    "Alignment contracts" in files["target"],
)
check(
    "changelog records accepted alignment boundary",
    BOUNDARY in files["changelog"]
    and "accepted the architecture and roadmap alignment" in norm(files["changelog"]),
)
check(
    "alignment record preserves isolated historical revalidation",
    "cc93fdd2311ca188ad03b0bd94293156ff243973" in files["record"]
    and "isolated local clone" in norm(files["record"]),
)

def is_ancestor(commit: str) -> bool:
    return subprocess.run(
        ["git", "merge-base", "--is-ancestor", commit, "HEAD"],
        cwd=ROOT,
        check=False,
        stdout=subprocess.DEVNULL,
        stderr=subprocess.DEVNULL,
    ).returncode == 0

def parents(commit: str) -> list[str]:
    try:
        return subprocess.check_output(
            ["git", "show", "-s", "--format=%P", commit],
            cwd=ROOT,
            text=True,
        ).strip().split()
    except (OSError, subprocess.CalledProcessError):
        return []

for label, commit in [
    ("signed BUSL base", BASE),
    ("signed alignment candidate", CANDIDATE),
    ("GitHub merge commit", MERGE),
    ("signed alignment boundary", BOUNDARY),
]:
    check(f"{label} is an ancestor", is_ancestor(commit))

check(
    "merge commit has exact base and candidate parents",
    parents(MERGE) == [BASE, CANDIDATE],
)
check(
    "signed alignment boundary directly follows merge",
    parents(BOUNDARY) == [MERGE],
)

tracked = set(
    subprocess.check_output(
        ["git", "ls-files", "--cached", "--others", "--exclude-standard"],
        cwd=ROOT,
        text=True,
    ).splitlines()
)
for path in list(paths.values())[:9] + [
    "tools/validation/validate_architecture_roadmap_alignment.py",
    "test-framework/governance/test_architecture_roadmap_alignment.sh",
    "tools/validation/phase-gates/validate_architecture_roadmap_alignment.sh",
]:
    check(f"repository registration contains {path}", path in tracked)

print(f"\nPASS checks: {passes}")
print(f"FAIL checks: {len(errors)}")
if errors:
    raise SystemExit(1)

print("Architecture and roadmap alignment acceptance-evidence validation PASSED.")
print(
    "This proves synchronized normative documentation, exact candidate, PR merge "
    "lineage, signed post-merge boundary, hosted validation evidence, and "
    "repository registration."
)
print(
    "It does not implement runtime behavior, alter schema, accept Phase 1 Step 3 "
    "or Phase 2, or establish production readiness."
)
