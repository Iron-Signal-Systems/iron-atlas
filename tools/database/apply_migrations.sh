#!/usr/bin/env bash
set -Eeuo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
manifest="$repo_root/sql/schema/manifests/atlas.manifest"

find_pg_bin() {
    local command_name="$1"
    if command -v "$command_name" >/dev/null 2>&1; then
        command -v "$command_name"
        return
    fi
    if command -v pg_config >/dev/null 2>&1; then
        local candidate
        candidate="$(pg_config --bindir)/$command_name"
        [[ -x "$candidate" ]] && { printf '%s\n' "$candidate"; return; }
    fi
    return 1
}

psql_bin="$(find_pg_bin psql)" || {
    echo "ERROR: psql was not found. Install PostgreSQL client tools." >&2
    exit 1
}

[[ -f "$manifest" ]] || { echo "ERROR: missing migration manifest" >&2; exit 1; }

work="$(mktemp -d)"
script="$work/run.sql"
trap 'rm -rf "$work"' EXIT

cat >"$script" <<'SQL'
\set ON_ERROR_STOP on
SET lock_timeout = '5s';
SET statement_timeout = '5min';
SET idle_in_transaction_session_timeout = '5min';
SELECT pg_advisory_lock(732041501);
SET ROLE atlas_schema_owner;
SQL

while IFS= read -r entry; do
    entry="${entry%%#*}"
    entry="$(printf '%s' "$entry" | xargs)"
    [[ -n "$entry" ]] || continue
    migration="$repo_root/sql/schema/$entry"
    [[ -f "$migration" ]] || { echo "ERROR: missing migration $entry" >&2; exit 1; }
    filename="$(basename "$migration")"
    version="${filename%%_*}"
    version=$((10#$version))
    digest="$(sha256sum "$migration" | awk '{print $1}')"
    prepared="$work/$filename"
    python3 - "$migration" "$prepared" "$version" "$filename" "$digest" <<'PY_PREPARE'
from pathlib import Path
import sys
source, destination, version, filename, digest = sys.argv[1:]
text = Path(source).read_text().rstrip()
if not text.endswith("COMMIT;"):
    raise SystemExit(f"migration does not end with COMMIT: {source}")
text = text[:-len("COMMIT;")].rstrip() + f"""

INSERT INTO atlas.schema_migration(
    migration_version, migration_filename, content_sha256, applied_by
) VALUES (
    {int(version)}, '{filename}', '{digest}', session_user
);

COMMIT;
"""
Path(destination).write_text(text)
PY_PREPARE

    cat >>"$script" <<SQL
\echo checking migration $filename
SELECT to_regclass('atlas.schema_migration') IS NOT NULL AS history_exists \gset m${version}_
\if :m${version}_history_exists
DO \$verify\$
DECLARE recorded_hash text;
BEGIN
    SELECT content_sha256 INTO recorded_hash
    FROM atlas.schema_migration
    WHERE migration_version = ${version};
    IF recorded_hash IS NOT NULL AND recorded_hash <> '${digest}' THEN
        RAISE EXCEPTION 'migration ${version} checksum mismatch';
    END IF;
END
\$verify\$;
SELECT EXISTS (
    SELECT 1 FROM atlas.schema_migration WHERE migration_version = ${version}
) AS applied \gset m${version}_
\else
\set m${version}_applied false
\endif
\if :m${version}_applied
\echo migration $filename already applied
\else
\ir $prepared
\echo applied migration $filename
\endif
SQL
done < "$manifest"

cat >>"$script" <<'SQL'
RESET ROLE;
SELECT pg_advisory_unlock(732041501);
SQL

"$psql_bin" --no-psqlrc --set ON_ERROR_STOP=1 --file "$script"
