# Business Source License 1.1 Transition Record

## Decision

Iron Atlas transitions prospectively from BSD 3-Clause to Business Source
License 1.1 (`BUSL-1.1`) to support a sustainable commercial product while
retaining source availability, non-production use, modification, and a defined
future Open Source conversion.

## Exact predecessor

The transition candidate must descend directly from the SSH-signed
authentication-assurance boundary:

```text
cc93fdd2311ca188ad03b0bd94293156ff243973
establish SSH-signed post-merge validation boundary for authentication assurance
```

## Historical license preservation

The exact BSD 3-Clause text from the predecessor tree is retained without
modification at:

```text
docs/records/licensing/IRON-ATLAS-BSD-3-CLAUSE-BEFORE-BSL.txt
```

Versions distributed under that BSD license remain available under it.

## BSL parameters

```text
Licensor: John Wood
Licensed Work: Iron Atlas
Additional Use Grant: None
Change Date: 2030-07-18
Change License: GNU Affero General Public License, Version 3 only
SPDX identifier: BUSL-1.1
```

## Commercial boundary

Because there is no Additional Use Grant, production use before the applicable
change event requires a separate commercial license. Non-production use remains
available under BSL 1.1.

## Change-License boundary

Each BSL-licensed version changes to AGPLv3-only on the declared Change Date or
the fourth anniversary of that version's first public BSL distribution,
whichever occurs first.

## Acceptance sequence

1. Apply this transition from the exact signed predecessor.
2. Run the dedicated licensing static validator and regression.
3. Run the complete Iron Atlas test framework and repository validation.
4. Create one SSH-signed candidate commit.
5. Push the candidate branch and open a dedicated PR.
6. Require all GitHub validation workflows to pass.
7. Merge with a merge commit.
8. Revalidate the merged `dev` tree.
9. Create and push an SSH-signed post-merge licensing boundary.
10. Record the concrete candidate, PR, merge, and post-merge hashes in retained
    evidence or a later acceptance record.

## Explicit exclusions

This transition does not:

- claim that BSL 1.1 is an Open Source license;
- retroactively relicense historical BSD distributions;
- grant trademark or logo rights;
- add a free production-use exception;
- modify source-code behavior;
- implement local credentials or TOTP;
- complete authentication, architecture alignment, or Phase 1 Step 3;
- establish production readiness; or
- replace review by qualified legal counsel.
