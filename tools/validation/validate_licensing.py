#!/usr/bin/env python3
from __future__ import annotations

import hashlib
import json
from pathlib import Path
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[2]
errors: list[str] = []
passes = 0

EXPECTED_BASE = "cc93fdd2311ca188ad03b0bd94293156ff243973"
EXPECTED_BSD = """BSD 3-Clause License

Copyright (c) 2026, Iron Signal Systems
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its
   contributors may be used to endorse or promote products derived from
   this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
"""

def check(name: str, condition: bool, detail: str = "") -> None:
    global passes
    if condition:
        print(f"PASS: {name}")
        passes += 1
    else:
        message = f"FAIL: {name}"
        if detail:
            message += f": {detail}"
        print(message, file=sys.stderr)
        errors.append(message)

def text(path: str) -> str:
    p = ROOT / path
    check(f"required file {path}", p.is_file())
    if not p.is_file():
        return ""
    return p.read_text(encoding="utf-8")

def normalized(value: str) -> str:
    return " ".join(value.split())

license_text = text("LICENSE")
licensing = text("LICENSING.md")
status = text("docs/governance/LICENSING-STATUS.md")
transition = text("docs/governance/LICENSING-TRANSITION-RECORD.md")
trademark = text("docs/governance/TRADEMARK-AND-BRANDING-POLICY.md")
backlog = text("docs/governance/POST-LICENSING-ALIGNMENT-BACKLOG.md")
history = text("docs/records/licensing/IRON-ATLAS-BSD-3-CLAUSE-BEFORE-BSL.txt")
readme = text("README.md")
docs_index = text("docs/README.md")
changelog = text("CHANGELOG.md")
requirements_text = text("validation/toolchain-requirements.json")

required_license_fragments = [
    "Business Source License 1.1",
    "Licensor: John Wood",
    "Licensed Work: Iron Atlas",
    "Additional Use Grant: None",
    "Change Date: 2030-07-18",
    "Change License: GNU Affero General Public License, Version 3 only",
    'The Business Source License (this document, or the "License") is not an Open',
    "make non-production use of the Licensed Work",
    "or the fourth anniversary of the first publicly",
    "must purchase a commercial license",
    "does not grant you any right in any trademark or logo",
    "Covenants of Licensor",
    "Not to modify this License in any other way.",
]
license_flat = normalized(license_text)
for fragment in required_license_fragments:
    check(f"LICENSE contains {fragment}", normalized(fragment) in license_flat)

check(
    "LICENSE does not retain BSD grant text",
    "Redistribution and use in source and binary forms" not in license_text,
)

check(
    "historical BSD license preserved exactly",
    history == EXPECTED_BSD,
    f"sha256={hashlib.sha256(history.encode()).hexdigest()}",
)

for document_name, document in [
    ("LICENSING.md", licensing),
    ("licensing status", status),
    ("transition record", transition),
]:
    check(f"{document_name} identifies BUSL-1.1", "BUSL-1.1" in document)
    check(
        f"{document_name} identifies no Additional Use Grant",
        "Additional Use Grant" in document and "None" in document,
    )
    check(
        f"{document_name} identifies AGPLv3-only change license",
        "AGPL-3.0-only" in document
        or "GNU Affero General Public License, Version 3 only" in document,
    )
    check(
        f"{document_name} preserves prospective BSD boundary",
        EXPECTED_BASE in document,
    )

check(
    "licensing status says BSL is not Open Source before change",
    "not an Open Source license" in normalized(status),
)
check(
    "transition record prohibits retroactive BSD revocation",
    "does not revoke" in status or "does not revoke" in transition,
)
check(
    "trademark policy reserves Iron Atlas branding",
    "Iron Atlas" in trademark and "does not grant rights" in trademark,
)
check(
    "trademark policy allows accurate factual reference",
    "nominative" in trademark and "non-misleading" in trademark,
)

required_backlog_fragments = [
    "not intended to become a local password or TOTP identity provider",
    "representative-provider compatibility testing",
    "module-runtime-and-failure-containment-model.md",
    "scheduled-evidence-ingestion-model.md",
    "monitoring-alerting-and-evidence-freshness-model.md",
    "evidence-candidate-and-atomic-acceptance-model.md",
    "atlas-ifi-snapshot-integration-contract.md",
    "fail-closed-and-adversarial-invariant-model.md",
    "mfa-and-authentication-assurance-requirements.md",
    "query the IFI PostgreSQL database directly",
    "preserve all historical accepted tags",
]
backlog_flat = normalized(backlog)
for fragment in required_backlog_fragments:
    check(
        f"alignment backlog contains {fragment}",
        normalized(fragment) in backlog_flat,
    )

check(
    "README identifies Business Source License 1.1",
    "Business Source License 1.1" in readme and "BUSL-1.1" in readme,
)
check(
    "README no longer declares current BSD license",
    "BSD 3-Clause. See [LICENSE](LICENSE)." not in readme,
)
check(
    "README rejects local password and TOTP credential ownership",
    "does not store local user passwords or TOTP seeds" in readme,
)
check(
    "documentation index links licensing governance",
    "LICENSING-STATUS.md" in docs_index
    and "POST-LICENSING-ALIGNMENT-BACKLOG.md" in docs_index,
)
check(
    "changelog records prospective BSL transition",
    "prospectively transitioned" in normalized(changelog).lower()
    and "BUSL-1.1" in changelog,
)

try:
    requirements = json.loads(requirements_text)
except json.JSONDecodeError as exc:
    check("toolchain requirements JSON parses", False, str(exc))
    requirements = {}
else:
    check("toolchain requirements JSON parses", True)

metadata = requirements.get("license", {})
expected_metadata = {
    "name": "Business Source License 1.1",
    "spdx_identifier": "BUSL-1.1",
    "licensor": "John Wood",
    "additional_use_grant": "None",
    "change_date": "2030-07-18",
    "change_license_name": "GNU Affero General Public License, Version 3 only",
    "change_license_spdx_identifier": "AGPL-3.0-only",
    "historical_license_spdx_identifier": "BSD-3-Clause",
    "signed_transition_base": EXPECTED_BASE,
}
for key, expected in expected_metadata.items():
    check(
        f"machine-readable license metadata {key}",
        metadata.get(key) == expected,
        f"expected {expected!r}, got {metadata.get(key)!r}",
    )

try:
    ancestor = subprocess.run(
        ["git", "merge-base", "--is-ancestor", EXPECTED_BASE, "HEAD"],
        cwd=ROOT,
        check=False,
    ).returncode == 0
except OSError:
    ancestor = False
check("signed transition base is an ancestor", ancestor)

tracked = set(
    subprocess.check_output(
        ["git", "ls-files", "--cached", "--others", "--exclude-standard"],
        cwd=ROOT,
        text=True,
    ).splitlines()
)
for path in [
    "LICENSE",
    "LICENSING.md",
    "docs/governance/LICENSING-STATUS.md",
    "docs/governance/LICENSING-TRANSITION-RECORD.md",
    "docs/governance/TRADEMARK-AND-BRANDING-POLICY.md",
    "docs/governance/POST-LICENSING-ALIGNMENT-BACKLOG.md",
    "docs/records/licensing/README.md",
    "docs/records/licensing/IRON-ATLAS-BSD-3-CLAUSE-BEFORE-BSL.txt",
    "tools/validation/validate_licensing.py",
    "test-framework/governance/test_business_source_license_transition.sh",
    "tools/validation/phase-gates/validate_business_source_license_transition.sh",
]:
    check(f"repository registration contains {path}", path in tracked)

print(f"\nPASS checks: {passes}")
print(f"FAIL checks: {len(errors)}")
if errors:
    raise SystemExit(1)
print("Business Source License 1.1 transition validation PASSED.")
print(
    "This proves the declared prospective BUSL-1.1 parameters, preserved "
    "historical BSD text, machine-readable metadata, trademark separation, "
    "alignment backlog, repository registration, and signed-base ancestry."
)
print(
    "It does not provide legal advice, grant production-use rights, complete "
    "architecture alignment, alter runtime behavior, or establish production readiness."
)
