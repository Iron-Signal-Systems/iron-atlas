#!/usr/bin/env bash
set -Eeuo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
source "$repo_root/tools/validation/lib/isolated_gate_revalidation.sh"

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

work="$(mktemp -d)"
trap 'rm -rf "$work"' EXIT
source_repo="$work/source"
mkdir -p "$source_repo/tools/validation/phase-gates" "$source_repo/tools/validation/static"
git -C "$source_repo" init --quiet
git -C "$source_repo" config user.name "Iron Atlas Gate Test"
git -C "$source_repo" config user.email "gate-test@example.invalid"

cat > "$source_repo/tools/validation/phase-gates/pass.sh" <<'PASS'
#!/usr/bin/env bash
exit 0
PASS
cat > "$source_repo/tools/validation/phase-gates/fail.sh" <<'FAIL'
#!/usr/bin/env bash
exit 23
FAIL
cat > "$source_repo/tools/validation/static/pass.py" <<'PY_PASS'
raise SystemExit(0)
PY_PASS
cat > "$source_repo/tools/validation/static/fail.py" <<'PY_FAIL'
raise SystemExit(29)
PY_FAIL
chmod +x \
    "$source_repo/tools/validation/phase-gates/pass.sh" \
    "$source_repo/tools/validation/phase-gates/fail.sh"
git -C "$source_repo" add -A
git -C "$source_repo" commit --quiet -m "test isolated gate"
commit="$(git -C "$source_repo" rev-parse HEAD)"

check \
    "successful isolated predecessor validator returns success" \
    isolated_gate_revalidate \
        "$source_repo" \
        "$commit" \
        tools/validation/phase-gates/pass.sh

failing_validator_is_rejected() {
    if isolated_gate_revalidate \
        "$source_repo" \
        "$commit" \
        tools/validation/phase-gates/fail.sh; then
        return 1
    fi
}
check \
    "failing isolated predecessor validator returns failure" \
    failing_validator_is_rejected

missing_validator_is_rejected() {
    if isolated_gate_revalidate \
        "$source_repo" \
        "$commit" \
        tools/validation/phase-gates/missing.sh; then
        return 1
    fi
}
check \
    "missing isolated predecessor validator returns failure" \
    missing_validator_is_rejected


check \
    "successful isolated Python validator returns success" \
    isolated_python_validator_revalidate \
        "$source_repo" \
        "$commit" \
        tools/validation/static/pass.py

failing_python_validator_is_rejected() {
    if isolated_python_validator_revalidate \
        "$source_repo" \
        "$commit" \
        tools/validation/static/fail.py; then
        return 1
    fi
}
check \
    "failing isolated Python validator returns failure" \
    failing_python_validator_is_rejected

python_parent_traversal_is_rejected() {
    if isolated_python_validator_revalidate \
        "$source_repo" \
        "$commit" \
        ../outside.py; then
        return 1
    fi
}
check \
    "isolated Python validator rejects parent traversal" \
    python_parent_traversal_is_rejected

printf '\nPASS checks: %d\n' "$pass"
printf 'FAIL checks: %d\n' "$fail"
(( fail == 0 ))
