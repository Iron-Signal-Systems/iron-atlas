#!/usr/bin/env bash
set -Eeuo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
cd "$repo_root"
predecessor_tag="phase-0-repository-and-executable-baseline-complete-v1"
predecessor_commit="6b24494e1e443eb8175af204e5a2e8ff66b2a2c6"

pass=0
fail=0
check() {
  local name="$1"; shift
  if "$@"; then printf 'PASS: %s\n' "$name"; pass=$((pass+1));
  else printf 'FAIL: %s\n' "$name"; fail=$((fail+1)); fi
}

echo "== Iron Atlas Phase 1 Step 1 validation =="
check "current branch is dev" bash -c '[[ "$(git branch --show-current)" == "dev" ]]'
check "working tree is clean" bash -c '[[ -z "$(git status --porcelain)" ]]'
check "accepted predecessor tag exists" git rev-parse --verify -q "${predecessor_tag}^{commit}"
check "accepted predecessor tag is unchanged" bash -c '[[ "$(git rev-parse "$1^{commit}")" == "$2" ]]' _ "$predecessor_tag" "$predecessor_commit"
check "accepted predecessor is ancestor" git merge-base --is-ancestor "$predecessor_commit" HEAD

revalidate_predecessor() {
  local tmp
  tmp="$(mktemp -d)"
  git clone --quiet --shared --no-checkout "$repo_root" "$tmp/repo"
  git -C "$tmp/repo" switch --quiet -C dev "$predecessor_commit"
  "$tmp/repo/tools/validation/phase-gates/validate_phase0_acceptance.sh"
  rm -rf "$tmp"
}
check "accepted Phase 0 predecessor revalidates in isolation" revalidate_predecessor
check "repository validation" ./tools/validation/validate_repository.sh --skip-database
check "complete test framework" ./test-framework/run_all.sh

printf '\nPASS checks: %d\nFAIL checks: %d\n' "$pass" "$fail"
if (( fail != 0 )); then
  echo "Phase 1 Step 1 validation FAILED."
  exit 1
fi
cat <<'MSG'

Phase 1 Step 1 validation PASSED.

This proves the documented migration, ownership, governed-identity,
approval, append-only history, and disposable PostgreSQL test boundary.
It does not prove production authentication, application persistence,
live collection, backup recovery, high availability, or production readiness.
MSG
