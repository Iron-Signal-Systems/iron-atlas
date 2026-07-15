#!/usr/bin/env bash
set -Eeuo pipefail
repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
cd "$repo_root"

echo "== Iron Atlas Phase 0 Step 1 validation =="
./tools/validation/validate_repository.sh
./test-framework/run_all.sh

echo
cat <<'EOF'
Phase 0 Step 1 implementation candidate validation PASSED.
This proves only the repository, documentation, current Go tests, and static contracts covered by this gate.
It does not prove production authentication, PostgreSQL persistence, live collection, complete vendor analysis, or production readiness.
EOF
