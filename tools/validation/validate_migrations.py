#!/usr/bin/env python3
from pathlib import Path
import sys

root = Path(__file__).resolve().parents[2]
manifest = root / 'sql/schema/manifests/atlas.manifest'
required = [
    "BEGIN;",
    "SET LOCAL lock_timeout = '5s';",
    "SET LOCAL statement_timeout = '1min';",
    "SET LOCAL idle_in_transaction_session_timeout = '1min';",
    "COMMIT;",
]
errors=[]
for line in manifest.read_text().splitlines():
    line=line.strip()
    if not line or line.startswith('#'): continue
    path=manifest.parent.parent / line
    if not path.exists():
        errors.append(f'missing migration: {line}')
        continue
    text=path.read_text()
    positions=[]
    for token in required:
        if token not in text: errors.append(f'{line}: missing {token}')
        else: positions.append(text.index(token))
    if len(positions)==len(required) and positions != sorted(positions):
        errors.append(f'{line}: execution-safety tokens are out of order')
if errors:
    print('\n'.join(errors), file=sys.stderr); raise SystemExit(1)
print('migration manifest and execution-safety headers validated')
