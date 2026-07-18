# Iron Atlas Documentation

> Owner: Iron Signal Systems
>
> Status: Active non-production development; accepted and candidate boundaries are described in the roadmap and acceptance records

## Start Here

## Brand Assets

- [Iron Atlas brand assets](branding/IRON-ATLAS-BRAND-ASSETS.md)

1. [Project mission](goals/PROJECT-MISSION.md)
2. [Product vision and operating mindset](goals/PRODUCT-VISION-AND-OPERATING-MINDSET.md)
3. [Target architecture](architecture/TARGET-ARCHITECTURE.md)
4. [Query, reachability, and change-impact model](architecture/QUERY-REACHABILITY-AND-CHANGE-IMPACT-MODEL.md)

5. [Compromise blast-radius and incident-impact intelligence](architecture/COMPROMISE-BLAST-RADIUS-AND-INCIDENT-IMPACT-INTELLIGENCE.md)
6. [Iron File Intelligence integration](architecture/IRON-FILE-INTELLIGENCE-INTEGRATION.md)
7. [BloodHound and identity attack-graph integration](architecture/BLOODHOUND-AND-IDENTITY-ATTACK-GRAPH-INTEGRATION.md)
8. [HTML5 interface and role workspaces](architecture/HTML5-INTERFACE-AND-ROLE-WORKSPACES.md)
9. [Data and record model](architecture/DATA-AND-RECORD-MODEL.md)
10. [Change management and two-person control](architecture/CHANGE-MANAGEMENT-AND-TWO-PERSON-CONTROL.md)
11. [Solo-developer operating model](engineering/SOLO-DEVELOPER-OPERATING-MODEL.md)
12. [Atlas primary-focus execution plan](roadmap/ATLAS-PRIMARY-FOCUS-EXECUTION-PLAN.md)
13. [Implementation roadmap](roadmap/IMPLEMENTATION-ROADMAP.md)

## Architecture

- [Modularity and dependency direction](architecture/MODULARITY-AND-DEPENDENCY-DIRECTION.md)
- [Operational-system complement and integration model](architecture/OPERATIONAL-SYSTEM-COMPLEMENT-AND-INTEGRATION-MODEL.md)
- [RBAC role and authority catalog](architecture/RBAC-ROLE-AND-AUTHORITY-CATALOG.md)
- [PostgreSQL migration and ownership model](architecture/POSTGRESQL-MIGRATION-AND-OWNERSHIP-MODEL.md)
- [PostgreSQL database security boundary](architecture/POSTGRESQL-DATABASE-SECURITY-BOUNDARY.md)
- [Go PostgreSQL runtime and identity context](architecture/GO-POSTGRESQL-RUNTIME-AND-IDENTITY-CONTEXT.md)
- [Trusted authentication and governed actor resolution](architecture/TRUSTED-AUTHENTICATION-AND-GOVERNED-ACTOR-RESOLUTION.md)
- [OIDC authorization-code and PKCE transaction implementation](architecture/OIDC-AUTHORIZATION-CODE-AND-PKCE-TRANSACTION-IMPLEMENTATION.md)
- [Evidence ingestion and protection](architecture/EVIDENCE-INGESTION-AND-PROTECTION.md)

- [External evidence context-bundle contract](architecture/EXTERNAL-EVIDENCE-CONTEXT-BUNDLE-CONTRACT.md)
- [Blast-radius result contract](architecture/BLAST-RADIUS-RESULT-CONTRACT.md)
- [ADR-0008 — IFI remains an external authoritative system](decisions/ADR-0008-IFI-EXTERNAL-AUTHORITATIVE-SYSTEM.md)
- [Firewall semantic analysis](architecture/FIREWALL-CONFIGURATION-SEMANTIC-ANALYSIS.md)
- [FortiGate YAML snapshot prototype](architecture/FORTIGATE-YAML-SNAPSHOT-PROTOTYPE.md)
- [ADR-0007 — maintained YAML decoder](decisions/ADR-0007-MAINTAINED-YAML-DECODER.md)
- [BloodHound and identity attack-graph integration](architecture/BLOODHOUND-AND-IDENTITY-ATTACK-GRAPH-INTEGRATION.md)
- [Firewall traffic path and SD-WAN model](architecture/FIREWALL-TRAFFIC-PATH-AND-SDWAN-MODEL.md)
- [Cisco evidence collection and preventive health](architecture/CISCO-EVIDENCE-COLLECTION-AND-PREVENTIVE-HEALTH.md)
- [Cisco trunk and endpoint attribution](architecture/CISCO-TRUNK-AND-ENDPOINT-ATTRIBUTION.md)
- [Cisco collection profile catalog](architecture/CISCO-COLLECTION-PROFILE-CATALOG.md)
- [Cisco NPS/RADIUS service authentication](architecture/CISCO-NPS-RADIUS-SERVICE-AUTHENTICATION.md)
- [Topology and Draw.io governance](architecture/TOPOLOGY-AND-DRAWIO-GOVERNANCE.md)
- [External-system-independent telemetry](architecture/EXTERNAL-SYSTEM-INDEPENDENT-TELEMETRY.md)
- [Zabbix integration contract](architecture/ZABBIX-INTEGRATION-CONTRACT.md)
- [Project and portfolio tracking](architecture/PROJECT-AND-PORTFOLIO-TRACKING.md)
- [Minimal Arch Linux deployment](architecture/MINIMAL-ARCH-LINUX-DEPLOYMENT.md)
- [Verification, validation, and acceptance](architecture/VERIFICATION-VALIDATION-AND-ACCEPTANCE.md)
- [Portable validation and canonical repository acceptance](architecture/PORTABLE-VALIDATION-AND-CANONICAL-REPOSITORY-ACCEPTANCE.md)

## Requirements, Testing, and Operations

- [System requirements](requirements/SYSTEM-REQUIREMENTS.md)
- [Testing and adversarial validation model](testing/TESTING-AND-ADVERSARIAL-VALIDATION-MODEL.md)

- [Compromise blast-radius correlation testing](testing/COMPROMISE-BLAST-RADIUS-CORRELATION-TESTING.md)
- [Compromise blast-radius program](roadmap/COMPROMISE-BLAST-RADIUS-PROGRAM.md)
- [ISRAS 0.1.1 adoption readiness](engineering/ISRAS-0.1.1-ADOPTION-READINESS.md)
- [Disposable PostgreSQL testing](testing/POSTGRESQL-DISPOSABLE-DATABASE-TESTING.md)
- [Go PostgreSQL runtime integration testing](testing/GO-POSTGRESQL-RUNTIME-INTEGRATION-TESTING.md)
- [Acceptance records](acceptance/README.md)
- [Canonical clean-clone validation](operations/CANONICAL-CLEAN-CLONE-VALIDATION.md)
- [Repository creation and first push](operations/REPOSITORY-CREATION-AND-FIRST-PUSH.md)

## Documentation Synchronization Rule

A material change is not complete until the product statement, architecture, requirements, implementation, tests, validation, status, roadmap, limitations, and next-work statement describe the same repository state.

Exploratory work may remain explicitly labeled as exploratory. A candidate may be self-validated without being independently reviewed. No document may imply production readiness, complete vendor coverage, independent assurance, or operational certainty that the retained evidence does not support.
