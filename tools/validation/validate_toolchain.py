#!/usr/bin/env python3
from __future__ import annotations

import json
from pathlib import Path
import re
import shutil
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[2]
REQ = json.loads((ROOT / "validation/toolchain-requirements.json").read_text())
errors: list[str] = []


def version_tuple(text: str) -> tuple[int, ...]:
    m = re.search(r"(\d+)(?:\.(\d+))?(?:\.(\d+))?", text)
    if not m:
        return ()
    return tuple(int(x) for x in m.groups(default="0"))


def output(*cmd: str) -> str:
    try:
        return subprocess.check_output(cmd, text=True, stderr=subprocess.STDOUT).splitlines()[0]
    except (OSError, subprocess.CalledProcessError) as exc:
        errors.append(f"cannot execute {' '.join(cmd)}: {exc}")
        return ""


def minimum(name: str, actual_text: str, required: str) -> None:
    actual = version_tuple(actual_text)
    needed = version_tuple(required)
    if not actual or actual < needed:
        errors.append(f"{name} {actual_text!r} is below required {required}")

req = REQ["requirements"]
minimum("bash", output("bash", "--version"), req["bash"]["minimum"])
minimum("git", output("git", "--version"), req["git"]["minimum"])
go_text = output("go", "version")
minimum("go", go_text, req["go"]["minimum"])
if version_tuple(go_text) >= version_tuple(req["go"]["maximum_exclusive"]):
    errors.append(f"Go {go_text!r} is outside tested range below {req['go']['maximum_exclusive']}")
minimum("python", sys.version.split()[0], req["python"]["minimum"])

for command in req["postgresql"]["commands"]:
    if shutil.which(command) is None:
        errors.append(f"missing PostgreSQL command: {command}")
minimum("PostgreSQL", output("psql", "--version"), req["postgresql"]["minimum"])

for command in req["utilities"]:
    if shutil.which(command) is None:
        errors.append(f"missing required utility: {command}")

if errors:
    for error in errors:
        print(f"FAIL: {error}", file=sys.stderr)
    raise SystemExit(1)

print("validated declared external toolchain requirements")
