#!/usr/bin/env bash
set -Eeuo pipefail
repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$repo_root"
pass=0
fail=0
check() { local name="$1"; shift; if "$@"; then printf 'PASS: %s\n' "$name"; pass=$((pass+1)); else printf 'FAIL: %s\n' "$name" >&2; fail=$((fail+1)); fi; }
check "alignment static validator" python3 tools/validation/validate_architecture_roadmap_alignment.py
check "alignment phase-gate syntax" bash -n tools/validation/phase-gates/validate_architecture_roadmap_alignment.sh
check "alignment regression syntax" bash -n test-framework/governance/test_architecture_roadmap_alignment.sh
check "no Atlas-local TOTP gate" bash -c '! grep -R -nF "validate_phase1_step3_totp_enrollment_verification_recovery.sh" docs README.md'
check "external-provider MFA" grep -Fq "approved external OpenID Connect identity provider" docs/security/MFA-AND-AUTHENTICATION-ASSURANCE-REQUIREMENTS.md
check "IFI direct database coupling prohibited" grep -Fq "query the IFI PostgreSQL database directly" docs/architecture/ATLAS-IFI-SNAPSHOT-INTEGRATION-CONTRACT.md
check "signed post-merge boundary governed" grep -Fq "SSH-signed empty post-merge boundary commit" docs/governance/SIGNED-CANDIDATE-AND-POST-MERGE-BOUNDARY-MODEL.md
printf '\nPASS checks: %d\nFAIL checks: %d\n' "$pass" "$fail"
(( fail == 0 ))
