#!/usr/bin/env python3
from pathlib import Path
import hashlib
import sys

root = Path(__file__).resolve().parents[2]
record = root / "SOURCE-SHA256SUMS.txt"
errors=[]
for line_number, raw in enumerate(record.read_text().splitlines(), 1):
    if not raw.strip():
        continue
    try:
        digest, filename = raw.split("  ", 1)
    except ValueError:
        errors.append(f"line {line_number}: malformed checksum record")
        continue
    path=root/filename
    if not path.is_file():
        errors.append(f"missing checksummed file: {filename}")
        continue
    actual=hashlib.sha256(path.read_bytes()).hexdigest()
    if actual != digest:
        errors.append(f"checksum mismatch: {filename}")
if errors:
    print("\n".join(f"FAIL: {e}" for e in errors), file=sys.stderr)
    raise SystemExit(1)
print("source SHA-256 records validated")
