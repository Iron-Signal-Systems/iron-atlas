#!/usr/bin/env bash
set -Eeuo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
source "$repo_root/tools/validation/lib/reporting.sh"

work="$(mktemp -d)"
trap 'rm -rf "$work"' EXIT

validation_report_init "intentional failure reporting test" "$work/fail-checks"
validation_note "test note retained before terminal result"

validation_run "successful prerequisite" bash -c 'printf "ready\n"'
if validation_run \
    "root failing check" \
    bash -c 'printf "ERROR: exact root cause\n" >&2; exit 7'
then
    printf 'FAIL: intentional root failure unexpectedly passed\n' >&2
    exit 1
fi

validation_skip \
    "dependent integration check" \
    "blocked by root failing check"

if validation_report_finish "$work/fail-report.txt"; then
    printf 'FAIL: failure report unexpectedly returned success\n' >&2
    exit 1
fi

grep -q '^PRIMARY FAILURE$' "$work/fail-report.txt"
grep -q '^Check: root failing check$' "$work/fail-report.txt"
grep -q '^Exit status: 7$' "$work/fail-report.txt"
grep -q '^Cause: ERROR: exact root cause$' "$work/fail-report.txt"
grep -q '^SKIPPED DEPENDENT CHECKS$' "$work/fail-report.txt"
grep -q '^REPORT NOTES$' "$work/fail-report.txt"

test "$(tail -n 1 "$work/fail-report.txt")" = \
    "FINAL RESULT: FAIL — root failing check — ERROR: exact root cause"

cat > "$work/nested-failure.sh" <<'EOF'
#!/usr/bin/env bash
printf 'ERROR: unrelated earlier diagnostic\n' >&2
printf 'FINAL RESULT: FAIL — nested static validation — FAIL: exact nested cause\n'
exit 9
EOF
chmod +x "$work/nested-failure.sh"

validation_report_init     "nested failure reporting test"     "$work/nested-checks"

if validation_run "nested runner" "$work/nested-failure.sh"; then
    printf 'FAIL: nested failure unexpectedly passed\n' >&2
    exit 1
fi

if validation_report_finish "$work/nested-report.txt"; then
    printf 'FAIL: nested failure report unexpectedly returned success\n' >&2
    exit 1
fi

grep -q '^Cause: FAIL: exact nested cause$'     "$work/nested-report.txt"

test "$(tail -n 1 "$work/nested-report.txt")" = \
    "FINAL RESULT: FAIL — nested runner — FAIL: exact nested cause"

validation_report_init "passing reporting test" "$work/pass-checks"
validation_run "passing check" bash -c 'printf "PASS payload\n"'
validation_report_finish "$work/pass-report.txt"

test "$(tail -n 1 "$work/pass-report.txt")" = \
    "FINAL RESULT: PASS — passing reporting test"

printf 'PASS: terminal actionable validation result reporting\n'
