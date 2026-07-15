# Phase 0 Step 1 Acceptance Record

## Record Status

- Decision: Accepted as a non-production development baseline
- Product: Iron Atlas
- Repository: `Iron-Signal-Systems/iron-atlas`
- Branch: `dev`
- Acceptance date: `2026-07-15T00:55:18Z`
- Accepted tag: `phase-0-repository-and-executable-baseline-complete-v1`

## Candidate

- Candidate implementation commit: `e74e4e0dcab239483327b5b329ce1f396c9837fa`
- Candidate short commit: ``
- Artifact method: deterministic Git archive of the candidate commit
- Artifact SHA-256: `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`
- Go version: `go version go1.26.5-X:nodwarf5 linux/amd64`
- Python version: `Python 3.14.6`
- C compiler: `gcc (GCC) 16.1.1 20260625`
- Host fingerprint: `Arch Linux; Linux 7.1.3-arch1-3 x86_64 GNU/Linux; CPUs=2; MemTotal=3866668 KiB`

## Accepted Scope

Phase 0 accepts the following non-production development boundary:

- Iron Atlas mission, purpose, terminology, and repository identity.
- Modular Go-first architecture and dependency direction.
- Embedded HTML5 interface candidate.
- Role-aware workspace and RBAC contract foundation.
- Memory-backed requester and approver independence candidate.
- Duplicate concurrent approval suppression.
- Module registry and vendor-adapter boundaries.
- Initial FortiGate hierarchical configuration parser.
- Initial OPNsense and pfSense XML import probes.
- Initial Cisco command-bundle and trunk-attribution analysis.
- Native Go Zabbix sender-protocol adapter.
- Minimal Arch Linux deployment documentation and service candidates.
- Repository validation, Go tests, race-enabled tests, and Phase 0 gate.
- Documentation, testing, validation, and acceptance structure.
- Prohibition against committing raw infrastructure evidence or secrets.

## Validation

- Repository validation: PASS
- Go formatting validation: PASS
- Go vet: PASS
- Race-enabled Go tests: PASS
- Documentation-link validation: PASS
- Migration-contract validation: PASS
- Draw.io XML validation: PASS
- Phase 0 Step 1 gate: PASS
- Correctness result: PASS
- Resource observation: NOT_RECORDED
- Performance thresholds: NOT_EVALUATED
- Validation transcript: `/tmp/iron-atlas-phase0-step1-validation.log` on the acceptance host

## Review and Approval

- Requester and implementer: John Wood
- Approval authority: John Wood, repository owner
- Independent human reviewer: Not assigned for this non-production single-maintainer development baseline
- Conflicts checked: A temporary single-maintainer exception is recorded below
- Operational two-person approval: Not exercised by this development acceptance

### Temporary Single-Maintainer Exception

Phase 0 contains no production authentication, live infrastructure collection,
production database, operational evidence, or authority to execute infrastructure
changes.

For this reason, the repository owner accepts this development baseline under a
temporary single-maintainer exception.

This exception:

- Does not authorize a production deployment.
- Does not authorize collection from live infrastructure.
- Does not authorize a governed infrastructure change.
- Does not permit a requester to approve their own operational change.
- Does not weaken the documented two-person change-management contract.
- Expires before production-boundary acceptance or operational use.
- Must not be used as precedent for bypassing independent operational approval.

## Known Limitations

The following capabilities are explicitly outside the accepted Phase 0 boundary:

- Production authentication
- PostgreSQL runtime persistence
- Production database ownership and privilege boundaries
- Live SSH collection
- NPS/RADIUS integration
- SSH host-key pinning
- Durable evidence ingestion and storage
- Complete FortiGate semantic analysis
- Complete OPNsense or pfSense semantic analysis
- Complete Cisco IOS or IOS XE collection
- Complete Catalyst 9800 collection
- Production Zabbix delivery
- Production deployment security
- Backup and restoration acceptance
- Compromise recovery
- Production performance budgets
- Operational infrastructure-change approval

## Security Assumptions

- Raw firewall configurations are prohibited from Git.
- Cisco technical-support reports are prohibited from Git.
- Credentials, private keys, shared secrets, and unredacted evidence are prohibited from Git.
- Development identity headers are not a production authentication boundary.
- Current parsers operate only within the limited Phase 0 candidate contract.
- A passing parser test does not prove a live network is correct or secure.

## Decision

The Phase 0 repository and executable baseline represented by candidate commit
`` is accepted for continued non-production development.

The annotated acceptance tag shall point to the acceptance commit containing
this record.

## Next Work

Begin Phase 1 Step 1:

- Establish the manifest-driven PostgreSQL migration framework.
- Define database ownership and production role topology.
- Establish the governed identity and authority persistence boundary.
- Add disposable PostgreSQL validation.
- Preserve requester and approver independence across independent database connections.
