# Iron Atlas Documentation

> Owner: Iron Signal Systems
>
> Status: Active non-production development; accepted and candidate boundaries are described in the roadmap and acceptance records

## Start Here

1. [Project mission](goals/PROJECT-MISSION.md)
2. [Product vision and operating mindset](goals/PRODUCT-VISION-AND-OPERATING-MINDSET.md)
3. [Target architecture](architecture/TARGET-ARCHITECTURE.md)
4. [Query, reachability, and change-impact model](architecture/QUERY-REACHABILITY-AND-CHANGE-IMPACT-MODEL.md)
5. [BloodHound and identity attack-graph integration](architecture/BLOODHOUND-AND-IDENTITY-ATTACK-GRAPH-INTEGRATION.md)
6. [HTML5 interface and role workspaces](architecture/HTML5-INTERFACE-AND-ROLE-WORKSPACES.md)
7. [Data and record model](architecture/DATA-AND-RECORD-MODEL.md)
8. [Change management and two-person control](architecture/CHANGE-MANAGEMENT-AND-TWO-PERSON-CONTROL.md)
9. [Solo-developer operating model](engineering/SOLO-DEVELOPER-OPERATING-MODEL.md)
10. [Atlas primary-focus execution plan](roadmap/ATLAS-PRIMARY-FOCUS-EXECUTION-PLAN.md)
11. [Implementation roadmap](roadmap/IMPLEMENTATION-ROADMAP.md)

## Architecture

- [Modularity and dependency direction](architecture/MODULARITY-AND-DEPENDENCY-DIRECTION.md)
- [Operational-system complement and integration model](architecture/OPERATIONAL-SYSTEM-COMPLEMENT-AND-INTEGRATION-MODEL.md)
- [RBAC role and authority catalog](architecture/RBAC-ROLE-AND-AUTHORITY-CATALOG.md)
- [PostgreSQL migration and ownership model](architecture/POSTGRESQL-MIGRATION-AND-OWNERSHIP-MODEL.md)
- [PostgreSQL database security boundary](architecture/POSTGRESQL-DATABASE-SECURITY-BOUNDARY.md)
- [Go PostgreSQL runtime and identity context](architecture/GO-POSTGRESQL-RUNTIME-AND-IDENTITY-CONTEXT.md)
- [Trusted authentication and governed actor resolution](architecture/TRUSTED-AUTHENTICATION-AND-GOVERNED-ACTOR-RESOLUTION.md)
- OIDC authorization-code and PKCE transaction implementation *(planned implementation document; not yet present)*
- [Evidence ingestion and protection](architecture/EVIDENCE-INGESTION-AND-PROTECTION.md)
- [Firewall semantic analysis](architecture/FIREWALL-CONFIGURATION-SEMANTIC-ANALYSIS.md)
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
- [Disposable PostgreSQL testing](testing/POSTGRESQL-DISPOSABLE-DATABASE-TESTING.md)
- [Go PostgreSQL runtime integration testing](testing/GO-POSTGRESQL-RUNTIME-INTEGRATION-TESTING.md)
- [Acceptance records](acceptance/README.md)
- [Canonical clean-clone validation](operations/CANONICAL-CLEAN-CLONE-VALIDATION.md)
- [Repository creation and first push](operations/REPOSITORY-CREATION-AND-FIRST-PUSH.md)

## Documentation Synchronization Rule

A material change is not complete until the product statement, architecture, requirements, implementation, tests, validation, status, roadmap, limitations, and next-work statement describe the same repository state.

Exploratory work may remain explicitly labeled as exploratory. A candidate may be self-validated without being independently reviewed. No document may imply production readiness, complete vendor coverage, independent assurance, or operational certainty that the retained evidence does not support.
