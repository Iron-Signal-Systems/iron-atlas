#!/usr/bin/env python3
from __future__ import annotations

import argparse
import hashlib
import json
from pathlib import Path
import re
import sys

ROOT = Path(__file__).resolve().parents[2]
parser = argparse.ArgumentParser()
parser.add_argument("--path", type=Path, default=ROOT / "validation/evidence")
args = parser.parse_args()
base = args.path.resolve()
errors: list[str] = []
runs = 0
secret_patterns = [
    re.compile(r"(?i)postgres(?:ql)?://[^\s:@/]+:[^@\s]+@"),
    re.compile(r"-----BEGIN [^-]*PRIVATE KEY-----"),
    re.compile(r"(?im)^authorization:\s*(?:bearer|basic)\s+\S+"),
]

for metadata in base.rglob("metadata.json") if base.exists() else []:
    runs += 1
    run = metadata.parent
    required = ["metadata.json", "environment.txt", "validation.log", "summary.txt", "sha256sums.txt"]
    for name in required:
        if not (run / name).is_file():
            errors.append(f"{run}: missing {name}")
    try:
        data = json.loads(metadata.read_text())
    except Exception as exc:
        errors.append(f"{metadata}: invalid JSON: {exc}")
        continue
    for key in ["schema_version", "boundary", "repository", "commit", "branch", "timestamp_utc", "command", "exit_code", "source"]:
        if key not in data:
            errors.append(f"{metadata}: missing field {key}")
    if not re.fullmatch(r"[0-9a-f]{40}", str(data.get("commit", ""))):
        errors.append(f"{metadata}: commit must be a full SHA-1")
    if data.get("repository") != "https://github.com/Iron-Signal-Systems/iron-atlas.git":
        errors.append(f"{metadata}: noncanonical repository")
    sums = run / "sha256sums.txt"
    if sums.is_file():
        for line in sums.read_text().splitlines():
            if not line.strip():
                continue
            try:
                expected, name = line.split("  ", 1)
            except ValueError:
                errors.append(f"{sums}: malformed line")
                continue
            target = run / name
            if not target.is_file():
                errors.append(f"{sums}: missing hashed file {name}")
                continue
            actual = hashlib.sha256(target.read_bytes()).hexdigest()
            if actual != expected:
                errors.append(f"{sums}: checksum mismatch for {name}")
    for path in run.iterdir():
        if not path.is_file() or path.name == "sha256sums.txt":
            continue
        text = path.read_text(errors="replace")
        for pattern in secret_patterns:
            if pattern.search(text):
                errors.append(f"{path}: possible secret material")

if errors:
    for error in errors:
        print(f"FAIL: {error}", file=sys.stderr)
    raise SystemExit(1)
print(f"validated {runs} committed validation evidence run(s)")
