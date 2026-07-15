#!/usr/bin/env bash
set -Eeuo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
cd "$repo_root"

record="docs/acceptance/PHASE-1-STEP-1-ACCEPTANCE-RECORD.md"
predecessor_tag="phase-0-repository-and-executable-baseline-complete-v1"
predecessor_commit="6b24494e1e443eb8175af204e5a2e8ff66b2a2c6"
accepted_tag="phase-1-step-1-postgresql-governance-foundation-complete-v1"

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

candidate_commit="$(
    sed -n \
        's/^- Candidate implementation commit: `\([0-9a-f]\{40\}\)`$/\1/p' \
        "$record" 2>/dev/null |
    head -n 1
)"

recorded_artifact_sha256="$(
    sed -n \
        's/^- Candidate Git archive SHA-256: `\([0-9a-f]\{64\}\)`$/\1/p' \
        "$record" 2>/dev/null |
    head -n 1
)"

echo "== Iron Atlas Phase 1 Step 1 acceptance validation =="

check \
    "current branch is dev" \
    bash -c '[[ "$(git branch --show-current)" == "dev" ]]'

check \
    "working tree is clean" \
    bash -c '[[ -z "$(git status --porcelain)" ]]'

check \
    "acceptance record exists" \
    test -f "$record"

check \
    "acceptance index exists" \
    test -f docs/acceptance/README.md

check \
    "candidate commit is a full SHA" \
    bash -c '[[ "$1" =~ ^[0-9a-f]{40}$ ]]' _ "$candidate_commit"

check \
    "candidate commit exists" \
    git cat-file -e "${candidate_commit}^{commit}"

check \
    "acceptance commit directly follows candidate" \
    bash -c '[[ "$(git rev-parse HEAD^)" == "$1" ]]' _ "$candidate_commit"

check \
    "accepted predecessor tag exists" \
    git rev-parse --verify -q "${predecessor_tag}^{commit}"

check \
    "accepted predecessor tag is unchanged" \
    bash -c '[[ "$(git rev-parse "$1^{commit}")" == "$2" ]]' \
        _ "$predecessor_tag" "$predecessor_commit"

check \
    "candidate descends from accepted predecessor" \
    git merge-base --is-ancestor "$predecessor_commit" "$candidate_commit"

check \
    "acceptance record contains no placeholders" \
    bash -c \
        '! grep -Eiq "\b(TBD|TODO|PENDING|PLACEHOLDER)\b|<[^>]+>" "$1"' \
        _ "$record"

check \
    "acceptance decision is explicit" \
    grep -Fq \
        "Decision: Accepted as a non-production PostgreSQL governance foundation" \
        "$record"

check \
    "accepted tag is declared" \
    grep -Fq "$accepted_tag" "$record"

recompute_artifact_sha256() {
    local actual
    actual="$(
        git archive \
            --format=tar \
            --prefix=iron-atlas-phase1-step1/ \
            "$candidate_commit" |
        sha256sum |
        awk '{print $1}'
    )"
    [[ "$actual" == "$recorded_artifact_sha256" ]]
}
check \
    "candidate Git archive hash matches acceptance record" \
    recompute_artifact_sha256

check \
    "temporary development exception preserves operational two-person control" \
    grep -Fq \
        "Does not weaken the documented two-person operational change-control contract." \
        "$record"

check \
    "root README reflects accepted Phase 1 Step 1 status" \
    grep -Fq \
        "Phase 1 Step 1 PostgreSQL governance foundation accepted" \
        README.md

check \
    "documentation index reflects accepted Phase 1 Step 1 status" \
    grep -Fq \
        "Phase 1 Step 1 PostgreSQL governance foundation accepted" \
        docs/README.md

check \
    "roadmap reflects accepted Phase 1 Step 1 status" \
    grep -Fq \
        "**Status:** Step 1 accepted under tag" \
        docs/roadmap/IMPLEMENTATION-ROADMAP.md

check \
    "acceptance index links the accepted record" \
    grep -Fq \
        "PHASE-1-STEP-1-ACCEPTANCE-RECORD.md" \
        docs/acceptance/README.md

check \
    "Phase 1 Step 1 implementation gate still passes" \
    ./tools/validation/phase-gates/validate_phase1_step1.sh

printf '\nPASS checks: %d\n' "$pass"
printf 'FAIL checks: %d\n' "$fail"

if (( fail != 0 )); then
    echo
    echo "Phase 1 Step 1 acceptance validation FAILED."
    exit 1
fi

cat <<'MSG'

Phase 1 Step 1 acceptance validation PASSED.

This accepts only the documented non-production PostgreSQL migration,
ownership, governed-identity, approval, append-only history, and disposable
validation boundary. It does not authorize production deployment, production
credentials, live collection, or operational infrastructure changes.
MSG
