#!/usr/bin/env python3
from pathlib import Path
import subprocess

ROOT = Path(__file__).resolve().parents[2]
BASE = "2347d21f779768f40496a93cb1d9140cc3b6e0ce"
passes = 0
failures: list[str] = []


def compact(value: str) -> str:
    # Static contract checks compare prose semantics, not Markdown decoration
    # or sentence-initial capitalization. Keep paths and required fragments
    # exact after normalizing those presentation-only differences.
    return " ".join(value.replace("`", "").split()).casefold()


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
    normalized = compact(text)
    for token in tokens:
        check(f"{path} contains {token}", compact(token) in normalized)
    return text


verifier = require(
    "internal/authentication/oidc/verifier.go",
    [
        "A successful provider login is not evidence of MFA",
        "!present &&",
        "claims.AuthenticationContext != \"\"",
        "len(claims.AuthenticationMethods) != 0",
        "authentication.ErrAuthenticationInvalid",
    ],
)
check(
    "verifier contains no named-provider assurance profile",
    not any(name in verifier.lower() for name in ("okta", "keycloak", "entra")),
)

assurance = require(
    "internal/authentication/assurance/assurance.go",
    [
        "func methodSetAccepted",
        "len(methods) != len(required)",
        "OutcomeAdditionalAuthenticationRequired",
        "OutcomeStepUpRequired",
    ],
)
check(
    "assurance code does not infer MFA from provider name",
    "providerid" not in assurance.lower() and "provider id" not in assurance.lower(),
)

require(
    "internal/authentication/oidc/verifier_test.go",
    [
        "TestVerifierDoesNotInferMFAFromSuccessfulLogin",
        "TestVerifierRequiresAuthenticationTimeForAssuranceEvidence",
        "urn:iron-atlas:test:mfa",
    ],
)
require(
    "internal/authentication/assurance/assurance_test.go",
    [
        "TestPolicyRejectsUngovernedAdditionalMethod",
        "provider-extra",
        "OutcomeAdditionalAuthenticationRequired",
    ],
)

architecture = require(
    "docs/architecture/PROVIDER-NEUTRAL-OIDC-ASSURANCE-EVIDENCE.md",
    [
        "Atlas-controlled evidence only",
        "acr and amr are accepted as assurance evidence only",
        "Exact governed method sets",
        "not representative-provider compatibility",
    ],
)
testing = require(
    "docs/testing/PROVIDER-NEUTRAL-OIDC-ASSURANCE-EVIDENCE-TESTING.md",
    [
        "successful login without assurance evidence never proves MFA",
        "additional ungoverned method is rejected",
        "No live or named external identity provider",
    ],
)
for name, text in (("architecture", architecture), ("testing", testing)):
    check(
        f"{name} makes no named-provider compatibility claim",
        not any(token in text.lower() for token in ("okta", "keycloak", "microsoft entra")),
    )

require(
    "README.md",
    [
        "provider-neutral assurance-evidence checkpoint is active",
        "representative-provider compatibility remains future evidence-backed work",
    ],
)
require(
    "docs/README.md",
    [
        "provider-neutral assurance evidence is the active bounded",
        "Provider-neutral OIDC assurance evidence",
    ],
)
require(
    "docs/requirements/PHASE-1-STEP-3-REQUIREMENTS-TRACEABILITY.md",
    ["Provider-neutral assurance-evidence implementation status"],
)
require(
    "docs/security/MFA-AND-AUTHENTICATION-ASSURANCE-REQUIREMENTS.md",
    [
        "Provider-neutral assurance-evidence checkpoint",
        "does not establish compatibility with any named provider",
    ],
)
require(
    "docs/roadmap/IMPLEMENTATION-ROADMAP.md",
    [
        "Provider-neutral assurance evidence is the active bounded Step 3 implementation checkpoint",
        "Representative-provider compatibility remains a future evidence-backed checkpoint",
    ],
)
require(
    "docs/roadmap/PHASE-GATE-PLAN.md",
    [
        "validate_phase1_step3_provider_neutral_assurance_evidence.sh",
        "representative-provider compatibility remains planned",
    ],
)
require(
    "docs/acceptance/PHASE-1-STEP-3-ACCEPTANCE-RECORD-TEMPLATE.md",
    ["Provider-neutral assurance-evidence checkpoint"],
)
require(
    "CHANGELOG.md",
    ["provider-neutral OIDC assurance-evidence candidate"],
)
require(
    "test-framework/authentication/test_phase1_step3_provider_neutral_assurance_evidence.sh",
    ["go test -race", "go tool govulncheck", "provider-neutral assurance evidence"],
)
require(
    "tools/validation/phase-gates/validate_phase1_step3_provider_neutral_assurance_evidence.sh",
    [BASE, "architecture-alignment evidence boundary remains valid", "bounded implementation candidate only"],
)
require(
    "test-framework/run_all.sh",
    [
        "provider-neutral assurance-evidence static validation",
        "provider-neutral assurance-evidence regression",
    ],
)
require(
    "tools/validation/validate_repository.sh",
    [
        "validate_phase1_step3_provider_neutral_assurance_evidence.py",
        "test_phase1_step3_provider_neutral_assurance_evidence.sh",
    ],
)

check(
    "signed provider-neutral implementation base is an ancestor",
    subprocess.run(
        ["git", "merge-base", "--is-ancestor", BASE, "HEAD"],
        cwd=ROOT,
        check=False,
        stdout=subprocess.DEVNULL,
        stderr=subprocess.DEVNULL,
    ).returncode == 0,
)

manifest = (ROOT / "FILE-MANIFEST.txt").read_text(encoding="utf-8").splitlines()
for path in (
    "docs/architecture/PROVIDER-NEUTRAL-OIDC-ASSURANCE-EVIDENCE.md",
    "docs/testing/PROVIDER-NEUTRAL-OIDC-ASSURANCE-EVIDENCE-TESTING.md",
    "test-framework/authentication/test_phase1_step3_provider_neutral_assurance_evidence.sh",
    "tools/validation/phase-gates/validate_phase1_step3_provider_neutral_assurance_evidence.sh",
    "tools/validation/validate_phase1_step3_provider_neutral_assurance_evidence.py",
):
    check(f"repository registration contains {path}", path in manifest)

check(
    "formal Step 3 acceptance record remains absent",
    not (ROOT / "docs/acceptance/PHASE-1-STEP-3-ACCEPTANCE-RECORD.md").exists(),
)
check(
    "abandoned representative-provider executable gate is absent",
    not (ROOT / "tools/validation/phase-gates/validate_phase1_step3_representative_provider_compatibility.sh").exists(),
)

print()
print(f"PASS checks: {passes}")
print(f"FAIL checks: {len(failures)}")
if failures:
    print("Phase 1 Step 3 provider-neutral assurance-evidence validation FAILED.")
    raise SystemExit(1)

print("Phase 1 Step 3 provider-neutral assurance-evidence validation PASSED.")
print(
    "This proves fail-closed correlation of assurance claims with explicit auth_time, "
    "exact governed method-set matching, synthetic Atlas-controlled evidence cases, "
    "and synchronized implementation, tests, documentation, and validation."
)
print(
    "It does not prove compatibility with any named provider, live hosted MFA, "
    "completed session lifecycle, formal Step 3 acceptance, or production readiness."
)
