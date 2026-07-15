#!/usr/bin/env python3
from pathlib import Path
import hashlib
import subprocess

root = Path(__file__).resolve().parents[2]
excluded = {"FILE-MANIFEST.txt", "SOURCE-SHA256SUMS.txt"}

tracked = subprocess.check_output(
    ["git", "ls-files", "--cached", "--others", "--exclude-standard"],
    cwd=root,
    text=True,
).splitlines()

files = sorted({name for name in tracked if (root / name).is_file()})
(root / "FILE-MANIFEST.txt").write_text("\n".join(files) + "\n")

records = []
for name in files:
    if name in excluded:
        continue
    digest = hashlib.sha256((root / name).read_bytes()).hexdigest()
    records.append(f"{digest}  {name}")
(root / "SOURCE-SHA256SUMS.txt").write_text("\n".join(records) + "\n")
print(f"regenerated manifests for {len(files)} files")
