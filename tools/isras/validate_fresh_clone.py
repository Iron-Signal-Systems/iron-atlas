#!/usr/bin/env python3
from __future__ import annotations

import argparse
import os
import sys
import tempfile
from pathlib import Path

from common import ISRASError, git, load_json, print_result, repository_root, run


def clone_environment(clone: Path) -> dict[str, str]:
    environment = os.environ.copy()
    if os.name == "nt":
        python_path = clone / ".isras-tools-venv" / "Scripts" / "python.exe"
    else:
        python_path = clone / ".isras-tools-venv" / "bin" / "python"

    go_tools_bin = clone / ".isras-go-tools" / "bin"
    environment["ISRAS_PYTHON"] = str(python_path)
    environment["ISRAS_GO_TOOLS_BIN"] = str(go_tools_bin)
    environment["PATH"] = (
        str(go_tools_bin)
        + os.pathsep
        + environment.get("PATH", "")
    )
    return environment


def bootstrap_clone(clone: Path) -> dict[str, str]:
    if os.name == "nt":
        bootstrap = clone / "tools" / "environment" / "Bootstrap-Tools.ps1"
        if not bootstrap.is_file():
            raise ISRASError(
                f"fresh-clone bootstrap is missing: {bootstrap.relative_to(clone)}"
            )
        run(
            ["pwsh", "-NoProfile", "-File", str(bootstrap)],
            cwd=clone,
        )
    else:
        bootstrap = clone / "tools" / "environment" / "bootstrap_tools.sh"
        if not bootstrap.is_file():
            raise ISRASError(
                f"fresh-clone bootstrap is missing: {bootstrap.relative_to(clone)}"
            )
        run(["bash", str(bootstrap)], cwd=clone)

    environment = clone_environment(clone)
    python_path = Path(environment["ISRAS_PYTHON"])
    go_tools_bin = Path(environment["ISRAS_GO_TOOLS_BIN"])
    binary_name = "govulncheck.exe" if os.name == "nt" else "govulncheck"

    if not python_path.is_file():
        raise ISRASError(
            f"fresh-clone Python environment was not created: {python_path}"
        )
    if not (go_tools_bin / binary_name).is_file():
        raise ISRASError(
            "fresh-clone pinned govulncheck installation was not created"
        )

    print_result("Fresh-clone pinned tool bootstrap", True)
    return environment


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--repo-root", default=".")
    parser.add_argument("--allow-dirty", action="store_true")
    args = parser.parse_args()

    repo_root = repository_root(args.repo_root)
    manifest = load_json(repo_root / "REPOSITORY-ASSURANCE.json")

    status = git(repo_root, "status", "--porcelain")
    if status and not args.allow_dirty:
        raise ISRASError("fresh-clone validation requires a clean working tree")

    commit = git(repo_root, "rev-parse", "HEAD")
    origin = manifest["canonical_origin"]
    portable = (
        manifest["validation"].get("portable_powershell")
        if os.name == "nt"
        else manifest["validation"].get("portable_shell")
    )
    if not portable:
        raise ISRASError(
            "portable validation entrypoint is not configured for this host"
        )

    with tempfile.TemporaryDirectory(
        prefix="isras-fresh-clone-"
    ) as temporary:
        clone = Path(temporary) / "repository"
        run(["git", "clone", "--no-local", origin, str(clone)])

        exists = (
            run(
                ["git", "cat-file", "-e", f"{commit}^{{commit}}"],
                cwd=clone,
                check=False,
            ).returncode
            == 0
        )
        print_result(
            "Exact local commit exists in canonical remote clone",
            exists,
            commit,
        )
        if not exists:
            raise ISRASError(
                "exact commit is not present in the canonical remote; "
                "push it before fresh-clone acceptance"
            )

        run(["git", "checkout", "--detach", commit], cwd=clone)
        actual = git(clone, "rev-parse", "HEAD")
        if actual != commit:
            raise ISRASError(
                "fresh clone did not check out the requested commit"
            )

        command_path = clone / portable.removeprefix("./")
        if not command_path.is_file():
            raise ISRASError(
                f"portable validator is missing from fresh clone: {portable}"
            )

        environment = bootstrap_clone(clone)
        if os.name == "nt":
            run(
                ["pwsh", "-NoProfile", "-File", str(command_path)],
                cwd=clone,
                env=environment,
            )
        else:
            run(
                ["bash", str(command_path)],
                cwd=clone,
                env=environment,
            )

        clean = git(clone, "status", "--porcelain") == ""
        print_result(
            "Fresh clone remains clean after ignored tool bootstrap",
            clean,
        )
        if not clean:
            raise ISRASError(
                "fresh clone became dirty after tool bootstrap and validation"
            )

        print_result("Fresh-clone portable validation", True)

    print("\nFresh-clone and remote-completeness validation PASSED.")
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except ISRASError as exc:
        print(f"FAIL: {exc}", file=sys.stderr)
        raise SystemExit(1)
