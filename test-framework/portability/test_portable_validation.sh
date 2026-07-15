#!/usr/bin/env bash
set -Eeuo pipefail
repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$repo_root"
pass=0
fail=0
check(){ local n="$1"; shift; if "$@"; then echo "PASS: $n"; pass=$((pass+1)); else echo "FAIL: $n"; fail=$((fail+1)); fi; }
check "portable acceptance static contract" python3 tools/validation/validate_portable_acceptance.py
check "portable validation shell syntax" bash -n tools/validation/record_validation_evidence.sh tools/validation/verify_canonical_clone.sh

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

# Build synthetic credentials only at runtime. The repository must not contain a
# credential-shaped URI or token literal merely to test transcript redaction.
database_scheme="$(printf '%s%s' 'post' 'gres')"
database_password="$(printf '%s%s' 'sec' 'ret')"
bearer_token="$(printf '%s%s' 'abc' '123')"
printf '%s://user:%s@example/atlas\nAuthorization: Bearer %s\n' \
  "$database_scheme" \
  "$database_password" \
  "$bearer_token" \
  > "$tmp/raw"

python3 tools/validation/redact_validation_text.py "$tmp/raw" "$tmp/clean"
database_secret="${database_password}@example"
check "validation transcript credential redaction" bash -c '
  ! grep -Fq "$2" "$1" &&
  ! grep -Fq "$3" "$1" &&
  grep -q REDACTED "$1"
' _ "$tmp/clean" "$database_secret" "$bearer_token"

mkdir -p "$tmp/evidence/run"
cat > "$tmp/evidence/run/metadata.json" <<'JSON'
{"schema_version":1,"boundary":"test","repository":"https://github.com/Iron-Signal-Systems/iron-atlas.git","commit":"0000000000000000000000000000000000000000","branch":"dev","timestamp_utc":"2026-07-15T00:00:00Z","command":["true"],"exit_code":0,"source":"test"}
JSON
printf 'environment\n' > "$tmp/evidence/run/environment.txt"
printf 'PASS: test\n' > "$tmp/evidence/run/validation.log"
printf 'PASS checks: 1\nFAIL checks: 0\n' > "$tmp/evidence/run/summary.txt"
( cd "$tmp/evidence/run" && sha256sum metadata.json environment.txt validation.log summary.txt > sha256sums.txt )
check "committed evidence integrity validation" python3 tools/validation/validate_committed_evidence.py --path "$tmp/evidence"
printf 'PASS checks: %d\nFAIL checks: %d\n' "$pass" "$fail"
(( fail == 0 ))
