#!/usr/bin/env bash
set -Eeuo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
cd "$repo_root"

record="docs/acceptance/PHASE-1-STEP-2-ACCEPTANCE-RECORD.md"
predecessor_tag="phase-1-step-1-postgresql-governance-foundation-complete-v1"
predecessor_commit="f41d932beff01e3faf1aeb73d6386a22c95cdda8"
expected_implementation="a56de3f9a9859f8daffea29b772091359dd3c3c9"
expected_local_evidence_commit="777019f9a55b5cb2988a13496005a40c9789a47a"
expected_evidence_boundary="8ea3e715aed989c742e6c9231614417122ead24d"
accepted_tag="phase-1-step-2-go-postgresql-runtime-and-identity-context-complete-v1"
local_run="validation/evidence/phase-1-step-2/local-phase-gate/20260715T154526Z-a56de3f9a985"
canonical_run="validation/evidence/phase-1-step-2/canonical-clean-clone/20260715T154723Z-777019f9a55b"

pass=0
fail=0

check() {
    local name="$1"
    shift
    if "$@"; then
        printf 'PASS: %s\n' "$name"
        pass=$((pass + 1))
    else
        printf 'FAIL: %s\n' "$name"
        fail=$((fail + 1))
    fi
}

record_value() {
    local label="$1"
    awk -v prefix="- ${label}: \`" 'index($0, prefix) == 1 { value = substr($0, length(prefix) + 1); sub(/`$/, "", value); print value; exit }' "$record"
}

implementation_commit="$(record_value 'Candidate implementation commit')"
evidence_boundary_commit="$(record_value 'Repository-complete evidence boundary commit')"
recorded_archive_sha256="$(record_value 'Candidate Git archive SHA-256')"
recorded_toolchain_sha256="$(record_value 'Toolchain requirements SHA-256')"

if git rev-parse --verify -q "${accepted_tag}^{commit}" >/dev/null; then
    acceptance_commit="$(git rev-parse "${accepted_tag}^{commit}")"
else
    acceptance_commit="$(git rev-parse HEAD)"
fi

metadata_matches() {
    local path="$1"
    local expected_commit="$2"
    local expected_source="$3"
    python3 - "$path" "$expected_commit" "$expected_source" <<'PY'
import json
from pathlib import Path
import sys
p = Path(sys.argv[1])
data = json.loads(p.read_text())
assert data["schema_version"] == 1
assert data["repository"] == "https://github.com/Iron-Signal-Systems/iron-atlas.git"
assert data["branch"] == "dev"
assert data["commit"] == sys.argv[2]
assert data["source"] == sys.argv[3]
assert data["exit_code"] == 0
PY
}

summary_has_contract() {
    local path="$1"
    grep -Fq 'PASS checks: 46' "$path" &&
    grep -Fq 'PASS checks: 14' "$path" &&
    grep -Fq 'PASS checks: 9' "$path" &&
    grep -Fq 'FAIL checks: 0' "$path" &&
    grep -Fq 'Correctness result: PASS' "$path" &&
    grep -Fq 'Resource observation: RECORDED_BY_DATABASE_TEST' "$path" &&
    grep -Fq 'Performance thresholds: NOT_EVALUATED' "$path" &&
    grep -Fq 'Phase 1 Step 2 validation PASSED.' "$path"
}

recompute_archive_sha256() {
    local actual
    actual="$(git archive --format=tar --prefix=iron-atlas-phase1-step2/ "$implementation_commit" | sha256sum | awk '{print $1}')"
    [[ "$actual" == "$recorded_archive_sha256" ]]
}

recompute_toolchain_sha256() {
    local actual
    actual="$(sha256sum validation/toolchain-requirements.json | awk '{print $1}')"
    [[ "$actual" == "$recorded_toolchain_sha256" ]]
}

echo "== Iron Atlas Phase 1 Step 2 acceptance validation =="

check "current branch is dev" bash -c '[[ "$(git branch --show-current)" == "dev" ]]'
check "working tree is clean" bash -c '[[ -z "$(git status --porcelain)" ]]'
check "acceptance record exists" test -f "$record"
check "acceptance index exists" test -f docs/acceptance/README.md
check "implementation commit is a full SHA" bash -c '[[ "$1" =~ ^[0-9a-f]{40}$ ]]' _ "$implementation_commit"
check "implementation commit is exact" bash -c '[[ "$1" == "$2" ]]' _ "$implementation_commit" "$expected_implementation"
check "evidence boundary is a full SHA" bash -c '[[ "$1" =~ ^[0-9a-f]{40}$ ]]' _ "$evidence_boundary_commit"
check "evidence boundary is exact" bash -c '[[ "$1" == "$2" ]]' _ "$evidence_boundary_commit" "$expected_evidence_boundary"
check "acceptance commit directly follows evidence boundary" bash -c '[[ "$(git rev-parse "$1^")" == "$2" ]]' _ "$acceptance_commit" "$expected_evidence_boundary"
check "acceptance commit is contained in current history" git merge-base --is-ancestor "$acceptance_commit" HEAD
check "implementation directly follows accepted Step 1" bash -c '[[ "$(git rev-parse "$1^")" == "$2" ]]' _ "$expected_implementation" "$predecessor_commit"
check "local evidence commit directly follows implementation" bash -c '[[ "$(git rev-parse "$1^")" == "$2" ]]' _ "$expected_local_evidence_commit" "$expected_implementation"
check "evidence boundary directly follows local evidence" bash -c '[[ "$(git rev-parse "$1^")" == "$2" ]]' _ "$expected_evidence_boundary" "$expected_local_evidence_commit"
check "accepted predecessor tag exists" git rev-parse --verify -q "${predecessor_tag}^{commit}"
check "accepted predecessor tag is unchanged" bash -c '[[ "$(git rev-parse "$1^{commit}")" == "$2" ]]' _ "$predecessor_tag" "$predecessor_commit"
check "acceptance record contains no placeholders" bash -c '! grep -Eiq "\b(TBD|TODO|PENDING|PLACEHOLDER)\b|<[^>]+>" "$1"' _ "$record"
check "acceptance decision is explicit" grep -Fq 'Decision: Accepted as a non-production Go PostgreSQL runtime, transaction-local identity-context, and portable-validation boundary' "$record"
check "accepted tag is declared" grep -Fq "$accepted_tag" "$record"
check "candidate Git archive hash matches acceptance record" recompute_archive_sha256
check "toolchain requirements hash matches acceptance record" recompute_toolchain_sha256
check "local evidence metadata matches" metadata_matches "$local_run/metadata.json" "$expected_implementation" "local-candidate"
check "canonical evidence metadata matches" metadata_matches "$canonical_run/metadata.json" "$expected_local_evidence_commit" "canonical-clean-clone"
check "local evidence checksums pass" bash -c 'cd "$1" && sha256sum -c sha256sums.txt' _ "$local_run"
check "canonical evidence checksums pass" bash -c 'cd "$1" && sha256sum -c sha256sums.txt' _ "$canonical_run"
check "local evidence proves the Step 2 contract" summary_has_contract "$local_run/summary.txt"
check "canonical evidence proves the Step 2 contract" summary_has_contract "$canonical_run/summary.txt"
check "canonical evidence declares clean-clone success" grep -Fq 'Canonical clean-clone validation PASSED.' "$canonical_run/summary.txt"
check "committed evidence integrity validator passes" python3 tools/validation/validate_committed_evidence.py
check "portable acceptance invariant is recorded" grep -Fq 'No implementation step may be accepted unless a clean clone from the canonical GitHub repository can execute its applicable validation using only version-controlled project artifacts, declared and verifiable external toolchain requirements, disposable test environments, and explicitly supplied non-repository secrets.' "$record"
check "temporary exception preserves operational two-person control" grep -Fq 'Does not weaken the documented two-person operational change-control contract.' "$record"
check "root README reflects accepted Step 2 status" grep -Fq 'Phase 1 Step 2 Go PostgreSQL runtime, identity-context, and portable-validation boundary accepted' README.md
check "documentation index reflects accepted Step 2 status" grep -Fq 'Phase 1 Step 2 Go PostgreSQL runtime, identity-context, and portable-validation boundary accepted' docs/README.md
check "roadmap reflects accepted Step 2 status" grep -Fq 'Steps 1 and 2 accepted; Step 2 is frozen under tag' docs/roadmap/IMPLEMENTATION-ROADMAP.md
check "acceptance index links Step 2 record" grep -Fq 'PHASE-1-STEP-2-ACCEPTANCE-RECORD.md' docs/acceptance/README.md
check "runtime architecture reflects acceptance" grep -Fq 'Accepted as the non-production Phase 1 Step 2 runtime and identity-context boundary' docs/architecture/GO-POSTGRESQL-RUNTIME-AND-IDENTITY-CONTEXT.md
check "pgx ADR reflects acceptance" grep -Fq 'Accepted for the non-production Phase 1 Step 2 boundary' docs/decisions/ADR-0004-PGX-POSTGRESQL-RUNTIME-DRIVER.md
check "runtime testing reflects acceptance" grep -Fq 'Accepted as the Phase 1 Step 2 disposable runtime-integration test boundary' docs/testing/GO-POSTGRESQL-RUNTIME-INTEGRATION-TESTING.md
check "portable validation document records accepted application" grep -Fq 'Phase 1 Step 2 is the first accepted application of this invariant.' docs/architecture/PORTABLE-VALIDATION-AND-CANONICAL-REPOSITORY-ACCEPTANCE.md
check "phase-gate index links acceptance validator" grep -Fq 'validate_phase1_step2_acceptance.sh' tools/validation/phase-gates/README.md
check "Phase 1 Step 2 implementation gate still passes" ./tools/validation/phase-gates/validate_phase1_step2.sh

printf '\nPASS checks: %d\n' "$pass"
printf 'FAIL checks: %d\n' "$fail"

if (( fail != 0 )); then
    echo
    echo "Phase 1 Step 2 acceptance validation FAILED."
    exit 1
fi

cat <<'MSG'

Phase 1 Step 2 acceptance validation PASSED.

This accepts only the documented non-production Go PostgreSQL adapter,
bounded pool, transaction-local identity context, governed persistence,
rollback, readiness, disposable integration-test, and repository-portable
validation boundary. It does not authorize production authentication,
credentials, TLS deployment, live collection, backup recovery, high
availability, or production use.
MSG
