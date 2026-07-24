#!/usr/bin/env bash
set -Eeuo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
canonical_repository="https://github.com/Iron-Signal-Systems/atlas.git"

usage() {
  echo "usage: $0 BOUNDARY COMMAND [ARG ...]" >&2
  exit 2
}

(( $# >= 2 )) || usage
boundary="$1"
shift
[[ "$boundary" =~ ^[a-z0-9][a-z0-9._/-]*$ && "$boundary" != *".."* ]] || {
  echo "invalid evidence boundary: $boundary" >&2
  exit 2
}
cd "$repo_root"
[[ -z "$(git status --porcelain)" ]] || {
  git status --short
  echo "working tree must be clean before evidence recording" >&2
  exit 1
}

commit="$(git rev-parse HEAD)"
branch="$(git branch --show-current)"
short="$(git rev-parse --short=12 HEAD)"
timestamp="$(date --utc '+%Y-%m-%dT%H:%M:%SZ')"
run_id="${IRON_ATLAS_VALIDATION_RUN_ID:-$(date --utc '+%Y%m%dT%H%M%SZ')-${short}}"
target="$repo_root/validation/evidence/$boundary/$run_id"
[[ ! -e "$target" ]] || { echo "evidence target already exists: $target" >&2; exit 1; }

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT
raw="$tmp/validation.raw.log"

set +e
"$@" 2>&1 | tee "$raw"
rc=${PIPESTATUS[0]}
set -e

python3 tools/validation/redact_validation_text.py "$raw" "$tmp/validation.log"
rm -f "$raw"

{
  grep -nE '^== |^PASS:|^FAIL:|^PASS checks:|^FAIL checks:|^Correctness result:|^Resource observation:|^Performance thresholds:|validation PASSED|validation FAILED' "$tmp/validation.log" || true
} > "$tmp/summary.txt"

{
  echo "repository=$canonical_repository"
  echo "commit=$commit"
  echo "branch=$branch"
  echo "timestamp_utc=$timestamp"
  . /etc/os-release
  echo "operating_system=$PRETTY_NAME"
  echo "kernel=$(uname -srmo)"
  echo "cpu_count=$(nproc)"
  echo "memory_kib=$(awk '/^MemTotal:/ {print $2}' /proc/meminfo)"
  git --version
  go version
  python3 --version
  postgres --version 2>&1 || true
  psql --version 2>&1 || true
} > "$tmp/environment.txt"

python3 - "$tmp/metadata.json" "$boundary" "$canonical_repository" "$commit" "$branch" "$timestamp" "$rc" "${IRON_ATLAS_VALIDATION_SOURCE:-local}" "$@" <<'PY_METADATA'
import json
from pathlib import Path
import sys
path=Path(sys.argv[1])
path.write_text(json.dumps({
  "schema_version": 1,
  "boundary": sys.argv[2],
  "repository": sys.argv[3],
  "commit": sys.argv[4],
  "branch": sys.argv[5],
  "timestamp_utc": sys.argv[6],
  "exit_code": int(sys.argv[7]),
  "source": sys.argv[8],
  "command": sys.argv[9:],
}, indent=2)+"\n")
PY_METADATA

(
  cd "$tmp"
  sha256sum metadata.json environment.txt validation.log summary.txt > sha256sums.txt
)
python3 tools/validation/validate_committed_evidence.py --path "$tmp"
mkdir -p "$(dirname "$target")"
mv "$tmp" "$target"
trap - EXIT
printf 'Recorded validation evidence: %s\n' "${target#$repo_root/}"
exit "$rc"
