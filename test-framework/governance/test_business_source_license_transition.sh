#!/usr/bin/env bash
set -Eeuo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$repo_root"

pass=0
fail=0

check() {
    local name="$1"
    shift
    if "$@"; then
        printf 'PASS: %s\n' "$name"
        pass=$((pass + 1))
    else
        printf 'FAIL: %s\n' "$name" >&2
        fail=$((fail + 1))
    fi
}

check \
    "licensing static validator" \
    python3 tools/validation/validate_licensing.py

check \
    "licensing phase-gate shell syntax" \
    bash -n tools/validation/phase-gates/validate_business_source_license_transition.sh

check \
    "licensing regression shell syntax" \
    bash -n test-framework/governance/test_business_source_license_transition.sh

check \
    "current license is BUSL-1.1" \
    grep -Fqx "Additional Use Grant: None" LICENSE

check \
    "historical BSD record remains exact" \
    python3 - <<'PY'
from pathlib import Path
historical = Path(
    "docs/records/licensing/IRON-ATLAS-BSD-3-CLAUSE-BEFORE-BSL.txt"
).read_text(encoding="utf-8")
current = Path("LICENSE").read_text(encoding="utf-8")
raise SystemExit(
    0
    if historical.startswith("BSD 3-Clause License\n")
    and "Redistribution and use in source and binary forms" in historical
    and current.startswith("Business Source License 1.1\n")
    and "Redistribution and use in source and binary forms" not in current
    else 1
)
PY

check \
    "README does not advertise BSD as current" \
    bash -c '! grep -Fq "BSD 3-Clause. See [LICENSE](LICENSE)." README.md'

check \
    "local TOTP implementation removed from required direction" \
    grep -Fq \
        "does not store local user passwords or TOTP seeds" \
        README.md

printf '\nPASS checks: %d\nFAIL checks: %d\n' "$pass" "$fail"
(( fail == 0 ))
