# System Requirements

## Foundation

- `IA-FND-001`: All production capabilities shall be modular behind versioned contracts.
- `IA-FND-002`: Core canonical records shall not depend on one vendor or monitoring product.
- `IA-FND-003`: Unsupported and uncertain parser or analyzer results shall remain visible.
- `IA-FND-004`: Raw protected evidence shall remain outside Git and receive integrity protection.
- `IA-FND-005`: Collectors shall not connect directly to PostgreSQL.
- `IA-FND-006`: Atlas shall complement established operational systems and shall not duplicate mature capabilities without a separately justified and accepted decision.
- `IA-FND-007`: Imported external-system records shall retain source and confidence and shall not silently become authoritative Atlas state.
- `IA-FND-008`: Atlas shall preserve the distinction among configured, observed, calculated, inferred, unknown, and conflicting state.
- `IA-FND-009`: Vendor adapters and collectors shall be evidence sources and shall not define the complete product boundary.
- `IA-FND-010`: Atlas shall correlate supported Cisco and FortiGate evidence through a vendor-neutral network and security model.

## Query and Intelligence

- `IA-QRY-001`: Atlas shall accept direct queries for supported IP addresses, CIDRs, subnets, VLANs, devices, interfaces, routes, policies, services, ports, tunnels, findings, and changes.
- `IA-QRY-002`: Query results shall provide an answer summary, supporting relationships, evidence, confidence, assumptions, conflicts, age, unsupported areas, unknowns, and pivots.
- `IA-QRY-003`: Atlas shall support prefix containment and longest-prefix-match analysis.
- `IA-QRY-004`: Atlas shall identify supported duplicate, overlapping, conflicting, and unexpectedly broad prefixes.
- `IA-QRY-005`: Atlas shall correlate supported VLAN, interface, trunk, SVI, gateway, VDOM, VRF, routing-domain, zone, route, policy, ACL, NAT, VIP, VPN, and SD-WAN relationships.
- `IA-QRY-006`: Atlas shall provide evidence-supported reachability results that are not limited to a binary yes or no.
- `IA-QRY-007`: Reachability results shall identify each supported path decision and its evidence.
- `IA-QRY-008`: Atlas shall distinguish current, prior accepted, proposed, expected post-change, and actual post-change state.
- `IA-QRY-009`: Atlas shall provide dependency and blast-radius analysis where evidence permits.
- `IA-QRY-010`: Atlas shall provide reviewable trust-boundary and attack-path context without representing it as an unexplained detection verdict.
- `IA-QRY-011`: Missing runtime health or incomplete evidence shall result in conditional, partial, unknown, conflicting, or unsupported state rather than silent certainty.
- `IA-QRY-012`: The first complete cross-vendor slice shall accept an IP address, CIDR, or VLAN and correlate the supported Cisco and FortiGate evidence required to explain placement, routing, control, reachability, evidence, and unknowns.

## Interface and Access

- `IA-UI-001`: The system shall provide an accessible HTML5 interface.
- `IA-UI-002`: Role perspectives shall support Network Operations Teams, Security Operations Teams, operational leaders, change authorities, reviewers, auditors, and teams.
- `IA-UI-003`: The primary interface shall provide a global search or query path for common identifiers and questions.
- `IA-UI-004`: The interface shall present the most relevant answer before requiring navigation into supporting detail.
- `IA-UI-005`: Material entities and relationships shall support context-preserving pivots.
- `IA-UI-006`: Configured, observed, calculated, inferred, unknown, conflicting, stale, incomplete, unsupported, proposed, approved, denied, implemented, validated, accepted, and superseded state shall remain distinguishable.
- `IA-UI-007`: Keyboard operation, semantic HTML, visible focus, and screen-reader support shall be functional requirements.
- `IA-UI-008`: The initial accessibility target shall be WCAG 2.1 Level AA.
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
- `IA-AUTH-017`: Atlas shall normalize authentication assurance independently of provider roles and govern accepted `acr`, `amr`, `auth_time`, MFA age, role-sensitive policy, and step-up requirements without treating missing or ambiguous assurance as successful MFA.
- `IA-AUTH-018`: Production authentication shall offer a phishing-resistant MFA option for high-impact authority. WebAuthn, FIDO2, passkeys, or hardware security keys are preferred; RFC 6238 TOTP may be supported as a compatible fallback.
- `IA-AUTH-019`: Atlas-local TOTP, when implemented, shall use encrypted per-authenticator secrets, enrollment proof, replay and rate-limit controls, one-time recovery codes, governed reset, encryption-key rotation, durable audit evidence, and no silent administrator bypass.

## Change Management

- `IA-CHG-001`: A requester shall not approve the requester’s own governed infrastructure change.
- `IA-CHG-002`: Approval and denial shall retain actor, authority, stage, scope, decision, reason, and time.
- `IA-CHG-003`: High-impact changes may require multiple independent approvers.
- `IA-CHG-004`: A change shall not be accepted before post-change validation and documentation reconciliation.
- `IA-CHG-005`: Material history shall use correction and supersession rather than silent rewriting.
- `IA-CHG-006`: A proposed change shall explain the risk of approval and the risk of denial or delay.
- `IA-CHG-007`: Atlas shall provide a leadership-facing decision summary and a synchronized engineering and security view derived from the same governed evidence.
- `IA-CHG-008`: Change analysis shall identify supported reachability, attack-path, dependency, blast-radius, availability, and security effects.
- `IA-CHG-009`: Change packages shall include current state, proposed state, implementation, validation, rollback, pre-change evidence, post-change evidence, expected differences, unexpected differences, and disposition.
- `IA-CHG-010`: Atlas shall remain read-only with respect to infrastructure until a separate provisioning boundary is accepted.
- `IA-CHG-011`: Completion of implementation commands shall not by itself establish successful change acceptance.

## PostgreSQL Runtime

- `IA-DB-001`: The Go runtime shall use a least-privileged application pool and shall not own or migrate database objects.
- `IA-DB-002`: Governed mutations shall bind authenticated actor identity only within the database transaction.
- `IA-DB-003`: Actor identity shall not leak between pooled connections, committed transactions, rolled-back transactions, or failed operations.
- `IA-DB-004`: Database dependency failure shall make readiness fail closed without changing liveness behavior.
- `IA-DB-005`: Database URLs, passwords, certificates, and tokens shall remain outside Git and application logs.

## Validation Portability and Engineering Practice

- `IA-VAL-001`: Applicable validators, phase gates, helpers, and disposable test-environment scripts shall be version-controlled.
- `IA-VAL-002`: External validation toolchain requirements shall be declared in a machine-readable, verifiable repository artifact.
- `IA-VAL-003`: Pinned external dependencies shall include integrity records and shall be verified before validation.
- `IA-VAL-004`: Retained validation and acceptance evidence shall be sanitized, checksummed, validated, and committed.
- `IA-VAL-005`: No formal milestone shall be accepted until its exact pushed commit passes applicable validation from a clean clone of the canonical GitHub repository.
- `IA-VAL-006`: Repository-external secrets shall be explicit, minimal, non-retained, and prohibited from logs and Git.
- `IA-VAL-007`: Self-authored and self-tested work shall be described as self-validated and shall not be represented as independently reviewed.
- `IA-VAL-008`: Exploratory parser or model work may use proportional validation when it is isolated, non-production, truthfully labeled, and does not weaken an accepted boundary.
- `IA-VAL-009`: Security-sensitive candidates shall receive adversarial, fail-closed, secret-handling, recovery, and applicable concurrency testing.
- `IA-VAL-010`: Multiple bounded workstreams may proceed in isolated branches or worktrees with no more than one active acceptance candidate per workstream.
- `IA-VAL-011`: Adoption of external engineering standards or validators shall occur through a visible, reviewable project change and shall not silently alter Atlas.
- `IA-VAL-012`: Top-level validation shall produce a meaningful final report containing the primary failure, command, exit status, cause, retained log path, additional unique failures, cascaded failures, and skipped dependent checks.
- `IA-VAL-013`: The last line emitted by a top-level validation command shall be an actionable `FINAL RESULT` containing the result and, on failure, the primary failing check and extracted cause.

## Firewall

- `IA-FW-001`: Support FortiGate, OPNsense, and pfSense adapter boundaries.
- `IA-FW-002`: Resolve interfaces, VLANs, zones, VDOMs, routes, policies, objects, NAT, VIP, VPN, and SD-WAN relationships where supported.
- `IA-FW-003`: Preserve policy and rule evaluation order.
- `IA-FW-004`: Distinguish configured, observed, calculated, inferred, unknown, and conflicting state.
- `IA-FW-005`: Provide evidence-supported traffic-path explanation.
- `IA-FW-006`: FortiGate configuration evidence shall not be represented as proof of current interface, route, VPN, session, HA, or SD-WAN runtime state.
- `IA-FW-007`: FortiGate ingestion shall preserve unsupported sections and unresolved relationships.
- `IA-FW-008`: FortiGate operational and diagnostic profiles shall be read-only, fixed, versioned, bounded, and separately accepted before live use.
- `IA-FW-009`: Firewall analysis shall identify supported management-plane and external exposure.
- `IA-FW-010`: Firewall records shall correlate with Cisco-derived topology and routing context where evidence permits.

## Cisco

- `IA-CSC-001`: Support 2960, 2960-S, 2960-X, 9200, 9300, 9500, and Catalyst 9800 profiles.
- `IA-CSC-002`: Support a comprehensive technical-support evidence profile for periodic or targeted collection.
- `IA-CSC-003`: Provide lighter recurring health collection.
- `IA-CSC-004`: Support NPS/RADIUS and Active Directory authentication with restricted device-local command authority where that deployment model is selected.
- `IA-CSC-005`: Collect supported device, interface, description, VLAN, trunk, pruning, CDP, LLDP, MAC, ARP, spanning-tree, ACL, route, port-channel, wireless, health, and error information.
- `IA-CSC-006`: Exclude trunks from local endpoint attribution while retaining full trunk analysis.
- `IA-CSC-007`: Use counter deltas and historical baselines where comparable evidence exists.
- `IA-CSC-008`: The first Cisco infrastructure-value slice shall prioritize Catalyst 9300L/9300, Catalyst 9500, and Catalyst 9800.
- `IA-CSC-009`: Offline sanitized Cisco evidence shall be accepted before restricted live collection.
- `IA-CSC-010`: Cisco evidence shall support normalized inventory, topology, reachability context, Zabbix reconciliation, Graylog context, generated maps, operational reports, and change analysis.
- `IA-CSC-011`: Cisco operational evidence shall distinguish configured state from observed forwarding, interface, neighbor, spanning-tree, route, wireless, and health state.
- `IA-CSC-012`: Restricted live collection shall use pinned host keys, fixed command profiles, bounded authority, no configuration mode, timeouts, cancellation, protected transcripts, and complete provenance.

## Identity Attack-Graph Integration

- `IA-IDG-001`: Atlas shall support a replaceable BloodHound integration boundary without depending on BloodHound internal database schemas.
- `IA-IDG-002`: Atlas shall preserve BloodHound version, collector version, source time, query or export identity, digest, and evidence lineage.
- `IA-IDG-003`: Raw SharpHound artifacts and unredacted BloodHound exports shall remain outside Git and receive protected-evidence controls.
- `IA-IDG-004`: Atlas shall distinguish identity privilege, packet reachability, and combined operational attack-path conclusions.
- `IA-IDG-005`: Atlas shall not represent a BloodHound path as proof of network reachability or network reachability as proof of identity compromise.
- `IA-IDG-006`: Cross-system identity and asset correlation shall retain stable identifiers, matching method, confidence, time, conflicts, and disposition.
- `IA-IDG-007`: Atlas shall not silently merge records based only on short hostname, display name, or mutable IP address.
- `IA-IDG-008`: Atlas may generate a versioned BloodHound OpenGraph extension and payloads containing selected, evidence-supported network and security context.
- `IA-IDG-009`: Structural graph relationships and traversable attack relationships shall remain distinct and separately validated.
- `IA-IDG-010`: Proposed Atlas state shall not be represented as current BloodHound state.
- `IA-IDG-011`: SharpHound execution or automated BloodHound collection requires a separately accepted, authorized, bounded, version-compatible collector boundary.
- `IA-IDG-012`: BloodHound integration failure shall not erase Atlas evidence or block unrelated Atlas operations.

## External-System Integration

- `IA-EXT-001`: External adapters shall be replaceable, versioned, least privileged, and isolated from canonical Atlas authority.
- `IA-EXT-002`: Atlas shall distinguish observed external state, imported metadata, generated recommendation, exported definition, approved provisioning request, applied state, and validated result.
- `IA-EXT-003`: Atlas may generate reviewable Zabbix maps, dashboards, template and discovery recommendations, reconciliation findings, and report context without representing them as applied.
- `IA-EXT-004`: Atlas may generate Graylog lookup data, enrichment context, queries, pipelines, streams, dashboards, and report definitions without representing them as applied.
- `IA-EXT-005`: External-system writes or provisioning shall require a separately accepted boundary with preview, attribution, authorization, bounded scope, idempotency where practical, reversal where practical, and post-application validation.
- `IA-EXT-006`: External-system delivery or provisioning failure shall not erase canonical evidence or block unrelated Atlas operations.
- `IA-EXT-007`: Security Onion and other security platforms shall remain responsible for detection and investigation; Atlas context shall not be represented as a detection verdict.
- `IA-EXT-008`: BloodHound shall remain responsible for identity and privilege attack-graph semantics; Atlas shall correlate rather than replace that function.
- `IA-EXT-009`: Generated Draw.io topology shall remain separate from curated diagram sources.

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
