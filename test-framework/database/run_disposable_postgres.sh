#!/usr/bin/env bash
set -Eeuo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

find_pg_bin() {
    local command_name="$1"
    if command -v "$command_name" >/dev/null 2>&1; then command -v "$command_name"; return; fi
    if command -v pg_config >/dev/null 2>&1; then
        local candidate="$(pg_config --bindir)/$command_name"
        [[ -x "$candidate" ]] && { printf '%s\n' "$candidate"; return; }
    fi
    return 1
}

for command_name in initdb pg_ctl psql createdb; do
    path="$(find_pg_bin "$command_name")" || {
        echo "ERROR: PostgreSQL command not found: $command_name" >&2
        echo "On Arch Linux run: sudo pacman -S --needed postgresql" >&2
        exit 1
    }
    printf -v "${command_name}_bin" '%s' "$path"
done

work="$(mktemp -d)"
data="$work/data"
socket="$work/socket"
mkdir -p "$socket"
port=$((55000 + ($$ % 900)))
started=false
cleanup() {
    if $started; then "$pg_ctl_bin" -D "$data" -m fast -w stop >/dev/null 2>&1 || true; fi
    rm -rf "$work"
}
trap cleanup EXIT

start_ns="$(date +%s%N)"
"$initdb_bin" -D "$data" --auth=trust --no-locale --encoding=UTF8 >/dev/null
"$pg_ctl_bin" -D "$data" -o "-F -h '' -k '$socket' -p $port" -w start >/dev/null
started=true

export PGHOST="$socket" PGPORT="$port" PGDATABASE=postgres PGUSER="$(id -un)"
"$psql_bin" --no-psqlrc -X -v ON_ERROR_STOP=1 -f "$repo_root/sql/bootstrap/development-roles.sql" >/dev/null
"$createdb_bin" --owner=atlas_database_owner iron_atlas_test
"$psql_bin" --no-psqlrc -X -v ON_ERROR_STOP=1 -d postgres -c \
  "GRANT CONNECT ON DATABASE iron_atlas_test TO atlas_migrator, atlas_application, atlas_test_runner; GRANT CREATE ON DATABASE iron_atlas_test TO atlas_schema_owner;" >/dev/null

export PGDATABASE=iron_atlas_test PGUSER=atlas_migrator
"$repo_root/tools/database/apply_migrations.sh" >/dev/null
"$repo_root/tools/database/apply_migrations.sh" >/dev/null

export PGUSER="$(id -un)"
"$psql_bin" --no-psqlrc -X -v ON_ERROR_STOP=1 -f "$repo_root/sql/bootstrap/runtime-grants.sql" >/dev/null

export ATLAS_PSQL="$psql_bin"
python3 "$repo_root/test-framework/database/test_database.py"

export IRON_ATLAS_TEST_DATABASE_URL="host=$socket port=$port dbname=iron_atlas_test user=atlas_application sslmode=disable"
(
  cd "$repo_root"
  go test -race -tags=integration \
    ./internal/database/postgresql \
    ./internal/change/postgresql \
    ./internal/authentication/postgresql
)

echo "PASS: Go PostgreSQL runtime integration tests"

size="$($psql_bin --no-psqlrc -X -Atqc "SELECT pg_database_size(current_database());")"
end_ns="$(date +%s%N)"
elapsed_ms=$(((end_ns-start_ns)/1000000))

echo "PostgreSQL version: $($psql_bin --version)"
echo "Disposable database bytes: $size"
echo "Database test elapsed milliseconds: $elapsed_ms"
echo "Sequential pooled identity checks: 500"
echo "Concurrent pooled identity operations: 600"
echo "Resource observation: RECORDED"
echo "Performance thresholds: NOT_EVALUATED"
