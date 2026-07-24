#!/usr/bin/env bash
set -Eeuo pipefail

canonical_repository="https://github.com/Iron-Signal-Systems/atlas.git"
canonical_branch="dev"

usage() {
  echo "usage: $0 EXPECTED_COMMIT VALIDATOR [ARG ...]" >&2
  exit 2
}
(( $# >= 2 )) || usage
expected_commit="$1"
shift
validator="$1"
shift
[[ "$expected_commit" =~ ^[0-9a-f]{40}$ ]] || { echo "expected commit must be a full SHA-1" >&2; exit 2; }
[[ "$validator" == ./* && "$validator" != *".."* ]] || { echo "validator must be a repository-relative ./ path" >&2; exit 2; }

remote_commit="$(git ls-remote "$canonical_repository" "refs/heads/$canonical_branch" | awk '{print $1}')"
[[ "$remote_commit" == "$expected_commit" ]] || {
  echo "canonical $canonical_branch is $remote_commit, expected $expected_commit" >&2
  exit 1
}

tmp="$(mktemp -d)"
cleanup() { local rc=$?; trap - EXIT; rm -rf "$tmp"; exit "$rc"; }
trap cleanup EXIT

git clone --quiet --branch "$canonical_branch" "$canonical_repository" "$tmp/repository"
git -C "$tmp/repository" fetch --quiet --force --tags origin
[[ "$(git -C "$tmp/repository" rev-parse HEAD)" == "$expected_commit" ]] || {
  echo "clean clone HEAD does not match expected commit" >&2
  exit 1
}
[[ -z "$(git -C "$tmp/repository" status --porcelain)" ]] || {
  echo "canonical clone is not clean" >&2
  exit 1
}
cd "$tmp/repository"
[[ -x "$validator" ]] || { echo "validator is missing or not executable: $validator" >&2; exit 1; }
python3 tools/validation/validate_toolchain.py
go mod download
go mod verify
"$validator" "$@"
printf '\nCanonical clean-clone validation PASSED.\nRepository: %s\nBranch: %s\nCommit: %s\nValidator: %s\n' \
  "$canonical_repository" "$canonical_branch" "$expected_commit" "$validator"
