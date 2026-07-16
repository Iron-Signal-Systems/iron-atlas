#!/usr/bin/env python3
from pathlib import Path
import subprocess
import sys

root = Path(__file__).resolve().parents[2]

tracked = subprocess.check_output(
    ["git", "ls-files", "--cached", "--others", "--exclude-standard"],
    cwd=root,
    text=True,
).splitlines()

files = sorted({name for name in tracked if (root / name).is_file()})

(root / "FILE-MANIFEST.txt").write_text(
    "\n".join(files) + "\n",
    encoding="utf-8",
)

subprocess.check_call(
    [
        sys.executable,
        str(root / "tools/isras/generate_source_manifest.py"),
        "--repo-root",
        str(root),
    ],
    cwd=root,
)

print(f"regenerated manifests for {len(files)} files")
