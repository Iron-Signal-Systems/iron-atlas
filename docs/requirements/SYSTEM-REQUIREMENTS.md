# System Requirements

## Foundation

- `IA-FND-001`: All production capabilities shall be modular behind versioned contracts.
- `IA-FND-002`: Core canonical records shall not depend on one vendor or monitoring product.
- `IA-FND-003`: Unsupported and uncertain parser results shall remain visible.
- `IA-FND-004`: Raw evidence shall remain outside Git and receive integrity protection.
- `IA-FND-005`: Collectors shall not connect directly to PostgreSQL.
- `IA-FND-006`: Atlas shall complement established operational systems and shall not duplicate mature capabilities without a separately justified and accepted decision.
- `IA-FND-007`: Imported external-system records shall retain source and confidence and shall not silently become authoritative Atlas state.

## Interface and Access

- `IA-UI-001`: The system shall provide an accessible HTML5 interface.
- `IA-UI-002`: Role workspaces shall support network technicians, administrators, security, reviewers, and teams.
- `IA-AUTH-001`: UI visibility shall not grant authority.
- `IA-AUTH-002`: Production authentication and authorization shall fail closed.
- `IA-AUTH-003`: Platform administration shall not automatically grant change approval.
- `IA-AUTH-004`: Production identities shall be accepted only from a configured, active, trusted authentication adapter.
- `IA-AUTH-005`: A verified provider and stable provider subject shall resolve to exactly one governed external identity and Atlas actor.
- `IA-AUTH-006`: Inactive providers and disabled or retired actors shall fail closed.
- `IA-AUTH-007`: Atlas role bindings, not provider roles or groups, shall be authoritative for request authorization.
- `IA-AUTH-008`: Request bodies, forms, query parameters, paths, and ordinary headers shall not select the production actor or roles.
- `IA-AUTH-009`: Development identity and production authentication modes shall be mutually exclusive and production mode shall not fall back.
- `IA-AUTH-010`: Browser sessions shall be server-side, opaque, bounded, rotated, revocable, and protected by secure cookie attributes.
- `IA-AUTH-011`: State-changing browser requests shall receive an accepted CSRF defense in addition to authentication.
- `IA-AUTH-012`: Trusted-proxy peers, headers, scheme, host, and bypass prevention shall be explicit and fail closed.
- `IA-AUTH-013`: Authentication, token, session, CSRF, and provider secrets shall not appear in logs, responses, Git, or retained evidence.
- `IA-AUTH-014`: Resolved actor identity shall remain immutable in server-side request context and transaction-local in PostgreSQL.
- `IA-AUTH-015`: Provider, external-identity, actor, role-binding, and session changes shall have bounded invalidation behavior.
- `IA-AUTH-016`: No adapter, provider claim, proxy, session, service, or administrator shall create an unrestricted execution context or propagate authority across boundaries.

## Change Management

- `IA-CHG-001`: A requester shall not approve the requester’s own governed change.
- `IA-CHG-002`: Approval shall retain actor, authority, stage, scope, decision, reason, and time.
- `IA-CHG-003`: High-impact changes may require multiple independent approvers.
- `IA-CHG-004`: A change shall not be accepted before post-change validation and documentation reconciliation.
- `IA-CHG-005`: Material history shall use correction and supersession rather than silent rewriting.

## PostgreSQL Runtime

- `IA-DB-001`: The Go runtime shall use a least-privileged application pool and shall not own or migrate database objects.
- `IA-DB-002`: Governed mutations shall bind authenticated actor identity only within the database transaction.
- `IA-DB-003`: Actor identity shall not leak between pooled connections, committed transactions, rolled-back transactions, or failed operations.
- `IA-DB-004`: Database dependency failure shall make readiness fail closed without changing liveness behavior.
- `IA-DB-005`: Database URLs, passwords, certificates, and tokens shall remain outside Git and application logs.

## Validation Portability

- Applicable validators, phase gates, helpers, and disposable test-environment scripts shall be version-controlled.
- External validation toolchain requirements shall be declared in a machine-readable, verifiable repository artifact.
- Pinned external dependencies shall include integrity records and shall be verified before validation.
- Retained validation and acceptance evidence shall be sanitized, checksummed, validated, and committed.
- No implementation step shall be accepted until its exact pushed commit passes applicable validation from a clean clone of the canonical GitHub repository.
- Repository-external secrets shall be explicit, minimal, non-retained, and prohibited from logs and Git.

## Firewall

- `IA-FW-001`: Support FortiGate, OPNsense, and pfSense adapter boundaries.
- `IA-FW-002`: Resolve interfaces, zones, routes, policies, objects, NAT, VPN, and SD-WAN relationships.
- `IA-FW-003`: Preserve policy and rule evaluation order.
- `IA-FW-004`: Distinguish configured, observed, calculated, inferred, unknown, and conflicting state.
- `IA-FW-005`: Provide evidence-supported traffic-path explanation.

## Cisco

- `IA-CSC-001`: Support 2960, 2960-S, 2960-X, 9200, 9300, 9500, and Catalyst 9800 profiles.
- `IA-CSC-002`: Collect a comprehensive technical-support evidence package every 30 days.
- `IA-CSC-003`: Provide lighter recurring health collection.
- `IA-CSC-004`: Support NPS/RADIUS and Active Directory authentication with restricted device-local command authority.
- `IA-CSC-005`: Collect port description, VLAN, trunk, pruning, CDP/LLDP, spanning tree, ACL, port-channel, QoS, and error information.
- `IA-CSC-006`: Exclude trunks from local endpoint attribution while retaining full trunk analysis.
- `IA-CSC-007`: Use counter deltas and historical baselines.
- `IA-CSC-008`: The first Cisco infrastructure-value slice shall prioritize Catalyst 9300L/9300, Catalyst 9500, and Catalyst 9800.
- `IA-CSC-009`: Offline sanitized Cisco evidence shall be accepted before restricted live collection.
- `IA-CSC-010`: Cisco evidence shall support normalized inventory, topology, Zabbix reconciliation, Graylog context, generated maps, and operational reports.

## External-System Integration

- `IA-INT-001`: External adapters shall be replaceable, versioned, least privileged, and isolated from canonical Atlas authority.
- `IA-INT-002`: Atlas shall distinguish observed external state, imported metadata, generated recommendation, exported definition, approved provisioning request, applied state, and validated result.
- `IA-INT-003`: Atlas may generate reviewable Zabbix maps, dashboards, template and discovery recommendations, reconciliation findings, and report context without representing them as applied.
- `IA-INT-004`: Atlas may generate Graylog lookup data, enrichment context, queries, pipelines, streams, dashboards, and report definitions without representing them as applied.
- `IA-INT-005`: External-system writes or provisioning shall require a separately accepted boundary with preview, attribution, authorization, bounded scope, idempotency where practical, reversal where practical, and post-application validation.
- `IA-INT-006`: External-system delivery or provisioning failure shall not erase canonical evidence or block unrelated Atlas operations.
- `IA-INT-007`: Security Onion and other security platforms shall remain responsible for detection and investigation; Atlas context shall not be represented as a detection verdict.
- `IA-INT-008`: Generated Draw.io topology shall remain separate from curated diagram sources.

## Telemetry

- `IA-TEL-001`: Canonical telemetry shall remain independent of Zabbix.
- `IA-TEL-002`: Zabbix delivery shall use a replaceable versioned adapter.
- `IA-TEL-003`: Delivery failure shall not block canonical recording or core operations.
- `IA-TEL-004`: Retry, backpressure, and dead-letter state shall be bounded and visible.

## Deployment

- `IA-DEP-001`: Production services shall be Go-first and deployable as signed binaries.
- `IA-DEP-002`: Rust, Node.js, and npm shall not be required without an accepted decision.
- `IA-DEP-003`: The Arch Linux host shall contain only required packages and governed administrative tools.
- `IA-DEP-004`: systemd services shall use restrictive sandboxing and dedicated identities.
