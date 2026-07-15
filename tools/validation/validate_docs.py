#!/usr/bin/env python3
from pathlib import Path
import re, sys

root = Path(__file__).resolve().parents[2]
errors = []
link_re = re.compile(r'\[[^\]]+\]\(([^)]+)\)')
for doc in list(root.rglob('*.md')):
    if '.git' in doc.parts:
        continue
    text = doc.read_text(encoding='utf-8')
    for target in link_re.findall(text):
        if '://' in target or target.startswith('#') or target.startswith('mailto:'):
            continue
        clean = target.split('#', 1)[0]
        if not clean:
            continue
        resolved = (doc.parent / clean).resolve()
        try:
            resolved.relative_to(root.resolve())
        except ValueError:
            errors.append(f'{doc.relative_to(root)}: link escapes repository: {target}')
            continue
        if not resolved.exists():
            errors.append(f'{doc.relative_to(root)}: missing link target: {target}')

if errors:
    print('\n'.join(errors), file=sys.stderr)
    raise SystemExit(1)
print(f'validated Markdown links across {len(list(root.rglob("*.md")))} files')
