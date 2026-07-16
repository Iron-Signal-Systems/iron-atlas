#!/usr/bin/env bash
set -Eeuo pipefail
repo_root="$(git rev-parse --show-toplevel 2>/dev/null || true)"
[[ -n "$repo_root" ]] || { printf 'FAIL: not in a Git work tree
' >&2; exit 1; }
python_cmd="${ISRAS_PYTHON:-python3}"
exec "$python_cmd" "$repo_root/tools/environment/go_tools.py" verify --repo-root "$repo_root"
