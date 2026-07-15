#!/usr/bin/env bash
set -Eeuo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
cd "$repo_root"

record="docs/acceptance/PHASE-0-STEP-1-ACCEPTANCE-RECORD.md"

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

echo "== Iron Atlas Phase 0 acceptance validation =="

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
    "candidate commit is an ancestor of HEAD" \
    git merge-base --is-ancestor "$candidate_commit" HEAD

check \
    "acceptance record contains no placeholders" \
    bash -c \
        '! grep -Eiq "\b(TBD|TODO|PENDING|PLACEHOLDER)\b|<[^>]+>" "$1"' \
        _ "$record"

check \
    "acceptance decision is explicit" \
    grep -Fq \
        "Decision: Accepted as a non-production development baseline" \
        "$record"

check \
    "temporary exception does not waive two-person control" \
    grep -Fq \
        "Does not weaken the documented two-person change-management contract." \
        "$record"

check \
    "production authentication remains excluded" \
    grep -Fq \
        "Production authentication" \
        "$record"

check \
    "live SSH collection remains excluded" \
    grep -Fq \
        "Live SSH collection" \
        "$record"

check \
    "accepted tag is declared" \
    grep -Fq \
        'phase-0-repository-and-executable-baseline-complete-v1' \
        "$record"

check \
    "root README reflects accepted Phase 0 status" \
    grep -Fq \
        "Phase 0 accepted as a non-production development baseline" \
        README.md

check \
    "documentation index reflects accepted Phase 0 status" \
    grep -Fq \
        "Phase 0 accepted as a non-production development baseline" \
        docs/README.md

check \
    "roadmap reflects accepted Phase 0 status" \
    grep -Fq \
        "**Status:** Accepted as a non-production development baseline." \
        docs/roadmap/IMPLEMENTATION-ROADMAP.md

check \
    "Phase 0 implementation gate still passes" \
    ./tools/validation/phase-gates/validate_phase0_step1.sh

printf '\nPASS checks: %d\n' "$pass"
printf 'FAIL checks: %d\n' "$fail"

if (( fail != 0 )); then
    echo
    echo "Phase 0 acceptance validation FAILED."
    exit 1
fi

cat <<'MSG'

Phase 0 acceptance validation PASSED.

This accepts only the documented non-production development baseline.
It does not authorize production deployment, live collection, or an
operational infrastructure change.
MSG
