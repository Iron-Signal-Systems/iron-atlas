# Iron Atlas Documentation

> Owner: Iron Signal Systems
>
> Status: Phase 1 Step 2 is accepted and the Phase 1 Step 3 trusted-authentication contract is integrated; the typed authentication-mode and immutable request-identity foundation is the active implementation candidate; no external provider, session, CSRF, or trusted-proxy implementation is accepted; not ready for production use

## Start Here

- [Project mission](goals/PROJECT-MISSION.md)
- [Target architecture](architecture/TARGET-ARCHITECTURE.md)
- [Modularity and dependency direction](architecture/MODULARITY-AND-DEPENDENCY-DIRECTION.md)
- [HTML5 interface and role workspaces](architecture/HTML5-INTERFACE-AND-ROLE-WORKSPACES.md)
- [Change management and two-person control](architecture/CHANGE-MANAGEMENT-AND-TWO-PERSON-CONTROL.md)
- [RBAC role and authority catalog](architecture/RBAC-ROLE-AND-AUTHORITY-CATALOG.md)
- [PostgreSQL migration and ownership model](architecture/POSTGRESQL-MIGRATION-AND-OWNERSHIP-MODEL.md)
- [PostgreSQL database security boundary](architecture/POSTGRESQL-DATABASE-SECURITY-BOUNDARY.md)
- [Go PostgreSQL runtime and identity context](architecture/GO-POSTGRESQL-RUNTIME-AND-IDENTITY-CONTEXT.md)
- [Trusted authentication and governed actor resolution](architecture/TRUSTED-AUTHENTICATION-AND-GOVERNED-ACTOR-RESOLUTION.md)
- [ADR-0004 — pgx PostgreSQL runtime driver](decisions/ADR-0004-PGX-POSTGRESQL-RUNTIME-DRIVER.md)
- [Cisco NPS/RADIUS service authentication](architecture/CISCO-NPS-RADIUS-SERVICE-AUTHENTICATION.md)
- [Cisco collection profile catalog](architecture/CISCO-COLLECTION-PROFILE-CATALOG.md)
- [Firewall traffic path and SD-WAN model](architecture/FIREWALL-TRAFFIC-PATH-AND-SDWAN-MODEL.md)
- [Data and record model](architecture/DATA-AND-RECORD-MODEL.md)
- [Zabbix integration contract](architecture/ZABBIX-INTEGRATION-CONTRACT.md)
- [Project and portfolio tracking](architecture/PROJECT-AND-PORTFOLIO-TRACKING.md)
- [Evidence ingestion and protection](architecture/EVIDENCE-INGESTION-AND-PROTECTION.md)
- [Firewall semantic analysis](architecture/FIREWALL-CONFIGURATION-SEMANTIC-ANALYSIS.md)
- [Cisco evidence collection and preventive health](architecture/CISCO-EVIDENCE-COLLECTION-AND-PREVENTIVE-HEALTH.md)
- [Cisco trunk and endpoint attribution](architecture/CISCO-TRUNK-AND-ENDPOINT-ATTRIBUTION.md)
- [Topology and Draw.io governance](architecture/TOPOLOGY-AND-DRAWIO-GOVERNANCE.md)
- [External-system-independent telemetry](architecture/EXTERNAL-SYSTEM-INDEPENDENT-TELEMETRY.md)
- [Minimal Arch Linux deployment](architecture/MINIMAL-ARCH-LINUX-DEPLOYMENT.md)
- [Verification, validation, and acceptance](architecture/VERIFICATION-VALIDATION-AND-ACCEPTANCE.md)
- [Portable validation and canonical repository acceptance](architecture/PORTABLE-VALIDATION-AND-CANONICAL-REPOSITORY-ACCEPTANCE.md)
- [ADR-0005 — canonical repository reproducibility](decisions/ADR-0005-CANONICAL-REPOSITORY-REPRODUCIBILITY.md)
- [Requirements catalog](requirements/SYSTEM-REQUIREMENTS.md)
- [Phase 1 Step 3 requirements traceability](requirements/PHASE-1-STEP-3-REQUIREMENTS-TRACEABILITY.md)
- [Atlas primary-focus execution plan](roadmap/ATLAS-PRIMARY-FOCUS-EXECUTION-PLAN.md)
- [Implementation roadmap](roadmap/IMPLEMENTATION-ROADMAP.md)
- [Testing model](testing/TESTING-AND-ADVERSARIAL-VALIDATION-MODEL.md)
- [Disposable PostgreSQL testing](testing/POSTGRESQL-DISPOSABLE-DATABASE-TESTING.md)
- [Go PostgreSQL runtime integration testing](testing/GO-POSTGRESQL-RUNTIME-INTEGRATION-TESTING.md)
- [Trusted authentication and governed actor resolution testing](testing/TRUSTED-AUTHENTICATION-AND-GOVERNED-ACTOR-RESOLUTION-TESTING.md)
- [Acceptance records](acceptance/README.md)
- [Repository creation and first push](operations/REPOSITORY-CREATION-AND-FIRST-PUSH.md)
- [Canonical clean-clone validation](operations/CANONICAL-CLEAN-CLONE-VALIDATION.md)

## Documentation Synchronization Rule

A phase is not complete until the root README, documentation indexes, architecture status, requirements, test documentation, validation tooling, committed retained evidence, acceptance record, terminology, version statements, and next-work statement describe the same repository state. No step may be accepted until the exact pushed commit passes applicable validation from a clean canonical GitHub clone.
