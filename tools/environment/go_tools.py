#!/usr/bin/env python3
from __future__ import annotations

import argparse
import json
import os
import re
import shutil
import subprocess
import sys
from pathlib import Path
from typing import Any


class ToolError(RuntimeError):
    pass


def run(args: list[str], cwd: Path, env: dict[str, str] | None = None, capture: bool = False) -> subprocess.CompletedProcess[str]:
    result = subprocess.run(
        args,
        cwd=cwd,
        env=env,
        text=True,
        stdout=subprocess.PIPE if capture else None,
        stderr=subprocess.PIPE if capture else None,
        check=False,
    )
    if result.returncode:
        if capture:
            if result.stdout:
                print(result.stdout, end='')
            if result.stderr:
                print(result.stderr, end='', file=sys.stderr)
        raise ToolError(f"command failed with status {result.returncode}: {' '.join(args)}")
    return result


def load_tool(path: Path) -> dict[str, Any]:
    try:
        value = json.loads(path.read_text(encoding='utf-8'))
    except (OSError, json.JSONDecodeError) as exc:
        raise ToolError(f'cannot load Go tool lock: {path}: {exc}') from exc
    if value.get('schema_version') != 'IRON-ATLAS-GO-TOOLS-V1':
        raise ToolError('unsupported Go tool lock schema')
    tools = value.get('tools')
    if not isinstance(tools, list) or len(tools) != 1:
        raise ToolError('the current Go tool lock must contain exactly one tool')
    required = {
        'name', 'binary', 'package', 'module', 'version',
        'module_sum', 'go_mod_sum', 'minimum_go_version',
    }
    missing = required - set(tools[0])
    if missing:
        raise ToolError(f'Go tool lock is missing fields: {sorted(missing)}')
    return tools[0]


def parse_go_version(value: str) -> tuple[int, int, int]:
    match = re.search(r'go(\d+)\.(\d+)(?:\.(\d+))?', value)
    if not match:
        raise ToolError(f'cannot parse Go version: {value}')
    return tuple(int(part or '0') for part in match.groups())


def check_go(go: str, minimum: str, root: Path) -> None:
    text = run([go, 'version'], root, capture=True).stdout.strip()
    if parse_go_version(text) < tuple(int(part) for part in minimum.split('.')):
        raise ToolError(f'Go {minimum} or newer is required; actual: {text}')
    print(f'PASS: Go version satisfies minimum {minimum}: {text}')


def verify_download(go: str, tool: dict[str, Any], root: Path) -> None:
    query = f"{tool['module']}@{tool['version']}"
    record = json.loads(run([go, 'mod', 'download', '-json', query], root, capture=True).stdout)
    expected = {
        'Path': tool['module'],
        'Version': tool['version'],
        'Sum': tool['module_sum'],
        'GoModSum': tool['go_mod_sum'],
    }
    mismatch = {key: [value, record.get(key)] for key, value in expected.items() if record.get(key) != value}
    if mismatch:
        raise ToolError('downloaded module does not match lock: ' + json.dumps(mismatch, sort_keys=True))
    print(f'PASS: downloaded Go module matches lock: {query}')


def tool_binary(tool_root: Path, name: str) -> Path:
    return tool_root / 'bin' / (name + ('.exe' if os.name == 'nt' else ''))


def verify_binary(go: str, tool: dict[str, Any], root: Path, tool_root: Path) -> None:
    binary = tool_binary(tool_root, tool['binary'])
    if not binary.is_file():
        raise ToolError(f'pinned Go tool binary is missing: {binary}')
    build = run([go, 'version', '-m', str(binary)], root, capture=True).stdout
    path_match = re.search(r'(?m)^\s*path\s+(\S+)\s*$', build)
    mod_match = re.search(r'(?m)^\s*mod\s+(\S+)\s+(\S+)\s+(\S+)\s*$', build)
    if not path_match or path_match.group(1) != tool['package']:
        raise ToolError('installed package path does not match lock')
    if not mod_match:
        raise ToolError('installed module identity is unavailable')
    if mod_match.groups() != (tool['module'], tool['version'], tool['module_sum']):
        raise ToolError('installed module version or checksum does not match lock')
    version_text = run([str(binary), '-version'], root, capture=True).stdout
    expected = f"Scanner: govulncheck@{tool['version']}"
    if expected not in version_text:
        raise ToolError(f'installed scanner does not report {expected}')
    print(f"PASS: installed tool matches lock: {tool['package']}@{tool['version']}")
    print(version_text.rstrip())


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument('command', choices=('bootstrap', 'verify'))
    parser.add_argument('--repo-root', default='.')
    parser.add_argument('--lock', default='tools/go-tools.lock.json')
    parser.add_argument('--tool-root')
    args = parser.parse_args()

    root = Path(args.repo_root).resolve()
    lock = Path(args.lock)
    if not lock.is_absolute():
        lock = root / lock
    tool_root = Path(args.tool_root or os.environ.get('ISRAS_GO_TOOLS_ROOT', root / '.isras-go-tools')).resolve()
    go = shutil.which('go')
    if not go:
        raise ToolError('Go is required to bootstrap or verify pinned Go tools')
    tool = load_tool(lock)
    check_go(go, tool['minimum_go_version'], root)

    if args.command == 'bootstrap':
        verify_download(go, tool, root)
        (tool_root / 'bin').mkdir(parents=True, exist_ok=True)
        env = os.environ.copy()
        env['GOBIN'] = str(tool_root / 'bin')
        env.setdefault('GOTOOLCHAIN', 'auto')
        run([go, 'install', f"{tool['package']}@{tool['version']}"], root, env=env)

    verify_binary(go, tool, root, tool_root)
    if args.command == 'bootstrap':
        print(f'Go assurance tool environment created at {tool_root}')
    return 0


if __name__ == '__main__':
    try:
        raise SystemExit(main())
    except ToolError as exc:
        print(f'FAIL: {exc}', file=sys.stderr)
        raise SystemExit(1)
