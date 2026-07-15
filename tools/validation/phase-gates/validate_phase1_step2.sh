#!/usr/bin/env bash
set -Eeuo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
cd "$repo_root"
predecessor_tag="phase-1-step-1-postgresql-governance-foundation-complete-v1"
predecessor_commit="f41d932beff01e3faf1aeb73d6386a22c95cdda8"

pass=0
fail=0
check() {
  local name="$1"; shift
  if "$@"; then printf 'PASS: %s\n' "$name"; pass=$((pass+1));
  else printf 'FAIL: %s\n' "$name"; fail=$((fail+1)); fi
}

echo "== Iron Atlas Phase 1 Step 2 validation =="
check "current branch is dev" bash -c '[[ "$(git branch --show-current)" == "dev" ]]'
check "working tree is clean" bash -c '[[ -z "$(git status --porcelain)" ]]'
check "accepted predecessor tag exists" git rev-parse --verify -q "${predecessor_tag}^{commit}"
check "accepted predecessor tag is unchanged" bash -c '[[ "$(git rev-parse "$1^{commit}")" == "$2" ]]' _ "$predecessor_tag" "$predecessor_commit"
check "accepted predecessor is ancestor" git merge-base --is-ancestor "$predecessor_commit" HEAD

source "$repo_root/tools/validation/lib/isolated_gate_revalidation.sh"

revalidate_predecessor() {
  isolated_gate_revalidate \
    "$repo_root" \
    "$predecessor_commit" \
    tools/validation/phase-gates/validate_phase1_step1_acceptance.sh
}
check "accepted Phase 1 Step 1 predecessor revalidates in isolation" revalidate_predecessor
check "portable canonical-repository validation contract" python3 tools/validation/validate_portable_acceptance.py
check "repository validation" ./tools/validation/validate_repository.sh --skip-database
check "complete test framework" ./test-framework/run_all.sh

printf '\nPASS checks: %d\nFAIL checks: %d\n' "$pass" "$fail"
if (( fail != 0 )); then
  echo "Phase 1 Step 2 validation FAILED."
  exit 1
fi
cat <<'MSG'

Phase 1 Step 2 validation PASSED.

This proves the documented Go PostgreSQL adapter, bounded pool,
transaction-local identity context, persistence, rollback, readiness,
disposable integration-test boundary, and repository-portable validation contract. It does not prove production
authentication, credential delivery, TLS provisioning, backup recovery,
high availability, live collection, or production readiness. Canonical clean-clone acceptance is a separate post-push requirement.
MSG
