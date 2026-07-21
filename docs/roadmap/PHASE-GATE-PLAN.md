# Phase and Gate Plan

## Purpose

This document defines the planned Iron Atlas gate sequence beneath the implementation roadmap.

It is a planning and governance artifact. A named future gate does not imply that its executable validator already exists, that its capability is implemented, or that its boundary is accepted.

Executable validators shall be created only when the corresponding contract or implementation boundary is active and sufficiently understood.

## Gate Classes

Each phase uses four gate classes.

### Phase-entry contract gate

A contract gate freezes:

- requirements;
- architecture;
- trust and privilege boundaries;
- threat model;
- data-classification rules;
- test model;
- acceptance criteria;
- explicit nonclaims;
- accepted predecessor; and
- planned evidence.

A contract gate does not accept implementation.

### Bounded implementation gate

An implementation gate proves one independently testable capability.

Each implementation gate has:

- one accepted predecessor;
- one declared candidate;
- one scope;
- one validation boundary;
- one evidence set;
- one explicit set of nonclaims; and
- one next step.

Passing an implementation gate does not establish formal phase acceptance or production readiness.

### Phase integration gate

An integration gate proves that the accepted or candidate components within a phase operate together without weakening predecessor boundaries.

It includes:

- cross-component behavior;
- predecessor revalidation;
- concurrency;
- hostile and failure-condition testing;
- resource observations;
- documentation synchronization; and
- explicit remaining limitations.

### Formal phase-acceptance gate

A phase-acceptance gate freezes:

- the exact canonical commit;
- the exact tested tree;
- retained sanitized validation evidence;
- clean canonical-clone validation;
- synchronized requirements, architecture, implementation, tests, limitations, status, and next work;
- exact accepted tag; and
- explicit nonclaims.

Self-validation shall not be represented as independent review.

## Gate Naming Convention

```text
validate_phase<N>_contract.sh
validate_phase<N>_step<M>_<capability>.sh
validate_phase<N>_integration.sh
validate_phase<N>_acceptance.sh
```

Optional modules use:

```text
validate_module_<name>_contract.sh
validate_module_<name>_implementation.sh
validate_module_<name>_integration.sh
validate_module_<name>_acceptance.sh
```

Historical gate names remain unchanged.

When successor documentation intentionally changes roadmap direction, frozen
checkpoint validators that assert the historical wording are re-run in an
isolated clone at their exact accepted or signed checkpoint. For this alignment,
the latest Step 3 implementation chain is revalidated at
`cc93fdd2311ca188ad03b0bd94293156ff243973`. Historical validators are not weakened to accept
successor claims.

---

## Phase 0 — Repository and Executable Baseline

Preserve the existing historical gates:

- `validate_phase0_step1.sh`
- `validate_phase0_acceptance.sh`

Do not rename, consolidate, or weaken the accepted Phase 0 history.

---

## Phase 1 — PostgreSQL Foundation and Governed Identity

Preserve all existing checkpoint and acceptance validators:

- `validate_phase1_step1.sh`
- `validate_phase1_step1_acceptance.sh`
- `validate_phase1_step2.sh`
- `validate_phase1_step2_acceptance.sh`
- `validate_phase1_step3_contract.sh`
- `validate_phase1_step3_authentication_foundation.sh`
- `validate_phase1_step3_governed_actor_resolution.sh`
- `validate_phase1_step3_oidc_id_token_verification.sh`
- `validate_phase1_step3_oidc_authorization_code_pkce.sh`
- `validate_phase1_step3_http_login_callback.sh`
- `validate_phase1_step3_authenticated_session.sh`
- `validate_phase1_step3_authentication_assurance.sh`

### Remaining trusted-authentication checkpoints

Planned gates:

- `validate_phase1_step3_representative_provider_compatibility.sh`
- `validate_phase1_step3_session_rotation_expiry_logout.sh`
- `validate_phase1_step3_csrf_protection.sh`
- `validate_phase1_step3_trusted_proxy_boundary.sh`
- `validate_phase1_step3_production_authenticator_wiring.sh`
- `validate_phase1_step3_emergency_and_recovery_access.sh`
- `validate_phase1_step3_integration.sh`
- `validate_phase1_step3_acceptance.sh`

The authentication-assurance checkpoint is merged and remains an implementation checkpoint rather than formal Step 3 acceptance. Atlas relies on approved external OIDC providers for primary authentication and MFA; no Atlas-local TOTP gate is planned. The listed representative-provider, session-lifecycle, CSRF, trusted-proxy, production-wiring, emergency-access, integration, and acceptance validators remain planned.

The Step 3 integration boundary proves:

```text
browser request
→ OIDC authorization request
→ authorization-code callback
→ token verification
→ governed actor resolution
→ bounded authenticated session
→ authorization context
→ logout and revocation behavior
```

### Credential delivery and PostgreSQL TLS

Planned gates:

- `validate_phase1_step4_contract.sh`
- `validate_phase1_step4_credential_delivery.sh`
- `validate_phase1_step4_credential_rotation.sh`
- `validate_phase1_step4_postgresql_tls.sh`
- `validate_phase1_step4_certificate_validation.sh`
- `validate_phase1_step4_failure_and_expiry.sh`
- `validate_phase1_step4_integration.sh`
- `validate_phase1_step4_acceptance.sh`

### Foundational recovery and resource budgets

Planned gates:

- `validate_phase1_step5_contract.sh`
- `validate_phase1_step5_backup_creation.sh`
- `validate_phase1_step5_restore_validation.sh`
- `validate_phase1_step5_connection_budgets.sh`
- `validate_phase1_step5_runtime_resource_budgets.sh`
- `validate_phase1_step5_hostile_failure_conditions.sh`
- `validate_phase1_step5_integration.sh`
- `validate_phase1_step5_acceptance.sh`

### Formal Phase 1 acceptance

Planned gate:

- `validate_phase1_acceptance.sh`

This gate revalidates the accepted Phase 1 step boundaries and freezes the complete PostgreSQL, identity, authentication, credential, TLS, recovery, and resource-governance foundation.

---

## Phase 2 — Evidence Intake, Protection, and Storage

### Contract gate

- `validate_phase2_contract.sh`

### Implementation gates

- `validate_phase2_step1_evidence_bundle_contract.sh`
- `validate_phase2_step2_intake_and_quarantine.sh`
- `validate_phase2_step3_staging_and_deduplication.sh`
- `validate_phase2_step4_protected_content_storage.sh`
- `validate_phase2_step5_classification_and_redaction.sh`
- `validate_phase2_step6_parser_isolation.sh`
- `validate_phase2_step7_recovery_and_hostile_evidence.sh`

### Integration gate

- `validate_phase2_integration.sh`

The integration path is:

```text
receipt
→ validation
→ quarantine or acceptance
→ staging
→ protected storage
→ isolated parser execution
→ normalized output
→ retained provenance
→ backup and restoration
```

### Acceptance gate

- `validate_phase2_acceptance.sh`

---

## Phase 3 — Cisco Offline Evidence and Normalization

### Contract gate

- `validate_phase3_contract.sh`

The contract defines Classic IOS, IOS XE, IOS XE wireless-controller roles, compatibility profiles, normalization, fixtures, resource limits, and unsupported behavior.

### Implementation gates

- `validate_phase3_step1_cisco_evidence_profiles.sh`
- `validate_phase3_step2_classic_ios_normalization.sh`
- `validate_phase3_step3_ios_xe_normalization.sh`
- `validate_phase3_step4_ios_xe_wireless_controller.sh`
- `validate_phase3_step5_cross_os_normalization.sh`
- `validate_phase3_step6_compatibility_profiles.sh`
- `validate_phase3_step7_adversarial_and_resource.sh`

### Integration gate

- `validate_phase3_integration.sh`

The integration gate runs common canonical queries across representative Classic IOS and IOS XE fixtures and verifies:

- deterministic normalization;
- common canonical records;
- preserved operating-system-specific meaning;
- explicit unsupported capabilities;
- compatibility-profile identity;
- fixture provenance; and
- resource governance.

### Acceptance gate

- `validate_phase3_acceptance.sh`

Phase 3 acceptance does not establish semantic topology conclusions or live collection.

---

## Phase 4 — FortiGate Offline Evidence and Normalization

### Contract gate

- `validate_phase4_contract.sh`

### Implementation gates

- `validate_phase4_step1_native_configuration.sh`
- `validate_phase4_step2_yaml_ingestion.sh`
- `validate_phase4_step3_yaml_native_equivalence.sh`
- `validate_phase4_step4_interfaces_vdoms_zones_objects.sh`
- `validate_phase4_step5_routes_policies_nat_vips.sh`
- `validate_phase4_step6_vpn_sdwan_local_in.sh`
- `validate_phase4_step7_operational_evidence.sh`
- `validate_phase4_step8_runtime_uncertainty.sh`
- `validate_phase4_step9_adversarial_and_resource.sh`

The gates preserve:

- policy order;
- configured-versus-observed state;
- unresolved object references;
- disabled-versus-absent configuration;
- unavailable runtime evidence;
- unsupported sections;
- ambiguous YAML semantics; and
- VDOM and routing-domain boundaries.

### Integration gate

- `validate_phase4_integration.sh`

The integration gate proves that native configuration, supported YAML, and operational evidence produce compatible canonical records where equivalence is claimed.

### Acceptance gate

- `validate_phase4_acceptance.sh`

OPNsense and pfSense are not required by this gate.

---

## Phase 5 — BloodHound Identity Context and Asset Correlation

### Contract gate

- `validate_phase5_contract.sh`

The contract defines approved offline evidence, source boundaries, privacy and classification rules, correlation decisions, uncertainty, and the prohibition on rebuilding BloodHound internally.

### Implementation gates

- `validate_phase5_step1_bloodhound_evidence_bundle.sh`
- `validate_phase5_step2_identity_path_normalization.sh`
- `validate_phase5_step3_asset_correlation_candidates.sh`
- `validate_phase5_step4_governed_asset_correlation.sh`
- `validate_phase5_step5_correlation_uncertainty.sh`
- `validate_phase5_step6_atlas_opengraph.sh`
- `validate_phase5_step7_privacy_and_adversarial.sh`

### Integration gate

- `validate_phase5_integration.sh`

The integration gate proves that approved BloodHound-derived identities correlate with Atlas assets while preserving:

- accepted decisions;
- rejected decisions;
- ambiguity;
- conflict;
- freshness;
- evidence;
- decision history; and
- explicit identity-versus-network separation.

### Acceptance gate

- `validate_phase5_acceptance.sh`

Direct BloodHound API access is not required.

---

## Phase 6 — Cross-Source Canonical Graph and Semantic Analysis

### Contract gate

- `validate_phase6_contract.sh`

### Implementation gates

- `validate_phase6_step1_graph_identity.sh`
- `validate_phase6_step2_layer2_graph.sh`
- `validate_phase6_step3_layer3_graph.sh`
- `validate_phase6_step4_firewall_graph.sh`
- `validate_phase6_step5_identity_asset_graph.sh`
- `validate_phase6_step6_state_and_time.sh`
- `validate_phase6_step7_uncertainty_and_provenance.sh`
- `validate_phase6_step8_semantic_analysis.sh`

### Integration gate

- `validate_phase6_cross_source_integration.sh`

Use representative Cisco, FortiGate, and BloodHound-derived fixtures together.

The graph must be:

- deterministic;
- evidence-backed;
- reproducible;
- source-traceable;
- uncertainty-preserving; and
- safe under conflicting input.

### Acceptance gate

- `validate_phase6_acceptance.sh`

---

## Phase 7 — Query, Reachability, and Attack-Path Intelligence

### Contract gate

- `validate_phase7_contract.sh`

The contract defines the answer model, reasoning model, evidence requirements, accuracy measurements, unsupported behavior, and nonclaims.

### Implementation gates

- `validate_phase7_step1_ip_cidr_vlan_intelligence.sh`
- `validate_phase7_step2_route_policy_explanation.sh`
- `validate_phase7_step3_forward_reachability.sh`
- `validate_phase7_step4_return_path.sh`
- `validate_phase7_step5_identity_path.sh`
- `validate_phase7_step6_combined_attack_path.sh`
- `validate_phase7_step7_answer_evidence.sh`
- `validate_phase7_step8_accuracy_and_adversarial.sh`

Every answer distinguishes:

- facts;
- inferred relationships;
- assumptions;
- unknowns;
- stale evidence;
- conflicts;
- unsupported behavior; and
- source evidence.

### Integration gate

- `validate_phase7_intelligence_integration.sh`

Run complete source-to-destination and identity-to-critical-asset scenarios across all accepted core evidence sources.

### Acceptance gate

- `validate_phase7_acceptance.sh`

Absence of evidence shall never be treated as proof of absence.

---

## Phase 8 — Projects, Change Impact, and Decision Support

### Contract gate

- `validate_phase8_contract.sh`

### Implementation gates

- `validate_phase8_step1_project_and_proposed_state.sh`
- `validate_phase8_step2_change_approval_governance.sh`
- `validate_phase8_step3_change_difference_analysis.sh`
- `validate_phase8_step4_dependency_and_blast_radius.sh`
- `validate_phase8_step5_approval_and_denial_risk.sh`
- `validate_phase8_step6_engineering_change_package.sh`
- `validate_phase8_step7_leadership_decision_summary.sh`
- `validate_phase8_step8_validation_and_rollback.sh`
- `validate_phase8_step9_emergency_change.sh`
- `validate_phase8_step10_post_change_disposition.sh`

### Integration gate

- `validate_phase8_change_lifecycle_integration.sh`

The integration path is:

```text
proposal
→ analysis
→ approval
→ implementation plan
→ validation
→ observed differences
→ disposition
→ rollback or acceptance
→ closure
```

The campaign includes:

- concurrent approval attempts;
- requester self-approval attempts;
- conflicting decisions;
- stale proposals;
- failed post-change validation; and
- emergency-change restrictions.

### Acceptance gate

- `validate_phase8_acceptance.sh`

---

## Phase 9 — Topology, Diagrams, and Accessible Interface

### Contract gate

- `validate_phase9_contract.sh`

The contract defines UI architecture, accessibility, diagram provenance, export behavior, evidence inspection, and performance budgets.

### Implementation gates

- `validate_phase9_step1_answer_workspace.sh`
- `validate_phase9_step2_network_views.sh`
- `validate_phase9_step3_identity_attack_path_views.sh`
- `validate_phase9_step4_change_impact_views.sh`
- `validate_phase9_step5_drawio_generation.sh`
- `validate_phase9_step6_svg_pdf_publication.sh`
- `validate_phase9_step7_diagram_provenance_and_drift.sh`
- `validate_phase9_step8_keyboard_accessibility.sh`
- `validate_phase9_step9_screen_reader_and_visual_accessibility.sh`
- `validate_phase9_step10_interface_resource_governance.sh`

### Integration gate

- `validate_phase9_user_workflow_integration.sh`

Representative workflows include:

- find an IP;
- explain a VLAN;
- investigate reachability;
- inspect an attack path;
- inspect supporting evidence;
- review a proposed change;
- export a diagram; and
- navigate without a mouse.

### Acceptance gate

- `validate_phase9_acceptance.sh`

---

## Phase 10 — Restricted Read-Only Collection and Evidence Refresh

### Contract gate

- `validate_phase10_contract.sh`

The contract defines the collector threat model, read-only restrictions, credential model, command allowlists, source protection, auditing, and explicit prohibition on modification.

### Implementation gates

- `validate_phase10_step1_collector_security_foundation.sh`
- `validate_phase10_step2_credentials_and_source_identity.sh`
- `validate_phase10_step3_command_allowlists.sh`
- `validate_phase10_step4_classic_ios_collection.sh`
- `validate_phase10_step5_ios_xe_collection.sh`
- `validate_phase10_step6_fortigate_collection.sh`
- `validate_phase10_step7_collected_bundle_replay.sh`
- `validate_phase10_step8_rate_and_resource_governance.sh`
- `validate_phase10_step9_partial_and_stale_collection.sh`
- `validate_phase10_step10_collection_audit.sh`
- `validate_phase10_step11_hostile_source_conditions.sh`

Hostile-source testing includes:

- slow responses;
- connection interruption;
- changed host keys;
- invalid certificates;
- authentication failure;
- unexpected prompts;
- paged output;
- oversized output;
- command refusal;
- source reboot;
- concurrency pressure; and
- attempted command escalation.

### Integration gate

- `validate_phase10_collection_integration.sh`

The integration path is:

```text
authorized collection
→ bounded command execution
→ evidence bundle
→ offline replay
→ normalization
→ graph update
→ complete audit
```

### Acceptance gate

- `validate_phase10_acceptance.sh`

No BloodHound API collector is required.

---

## Phase 11 — Production Security, Recovery, and Representative Deployment

### Contract gate

- `validate_phase11_contract.sh`

### Build and supply-chain gates

- `validate_phase11_step1_reproducible_build.sh`
- `validate_phase11_step2_signed_build_and_provenance.sh`
- `validate_phase11_step3_sbom_and_dependency_integrity.sh`

### Runtime-hardening gates

- `validate_phase11_step4_service_identity_and_isolation.sh`
- `validate_phase11_step5_secret_delivery_and_rotation.sh`
- `validate_phase11_step6_host_and_network_hardening.sh`

### Logging and integrity gates

- `validate_phase11_step7_off_host_logging.sh`
- `validate_phase11_step8_integrity_anchors.sh`

### Recovery gates

- `validate_phase11_step9_backup_protection.sh`
- `validate_phase11_step10_restore_validation.sh`
- `validate_phase11_step11_break_glass.sh`
- `validate_phase11_step12_trusted_rebuild.sh`
- `validate_phase11_step13_compromise_recovery.sh`

### Representative-deployment gates

- `validate_phase11_step14_representative_host.sh`
- `validate_phase11_step15_boot_and_service_recovery.sh`
- `validate_phase11_step16_storage_and_network_failure.sh`
- `validate_phase11_step17_resource_and_performance.sh`
- `validate_phase11_step18_upgrade_and_rollback.sh`

### Integration gate

- `validate_phase11_production_integration.sh`

Perform a complete representative deployment, backup, failure, restore, upgrade, rollback, break-glass, and rebuild campaign.

### Acceptance gate

- `validate_phase11_acceptance.sh`

---

## Phase 12 — Controlled Pilot and Operational Acceptance

### Pilot contract gate

- `validate_phase12_pilot_contract.sh`

The contract freezes:

- pilot authorization;
- approved environment;
- approved evidence sources;
- read-only restrictions;
- success measurements;
- prohibited actions;
- rollback and removal plan;
- privacy and handling rules;
- manual verification method; and
- known nonclaims.

### Pilot implementation and measurement gates

- `validate_phase12_step1_environment_baseline.sh`
- `validate_phase12_step2_pilot_deployment.sh`
- `validate_phase12_step3_answer_verification.sh`
- `validate_phase12_step4_accuracy_accounting.sh`
- `validate_phase12_step5_evidence_quality.sh`
- `validate_phase12_step6_operational_impact.sh`
- `validate_phase12_step7_product_value.sh`
- `validate_phase12_step8_change_package_value.sh`
- `validate_phase12_step9_recovery_exercise.sh`
- `validate_phase12_step10_residual_risk.sh`

### Pilot integration gate

- `validate_phase12_pilot_integration.sh`

The gate confirms that collected pilot evidence, accuracy measurements, resource measurements, user feedback, recovery evidence, and residual risk all describe the same pilot boundary.

### Formal operational-acceptance gate

- `validate_phase12_operational_acceptance.sh`

The final core acceptance gate verifies:

- exact canonical commit;
- exact deployed build;
- signed build provenance;
- retained pilot evidence;
- manually verified accuracy;
- false-positive and false-negative accounting;
- complete limitations;
- accepted residual risk;
- recovery proof;
- synchronized documentation; and
- no dependency on optional modules.

---

## Optional Module Gates

Optional modules are not part of the mandatory Phase 0–12 chain.

Examples include:

```text
validate_module_zabbix_contract.sh
validate_module_zabbix_implementation.sh
validate_module_zabbix_integration.sh
validate_module_zabbix_acceptance.sh

validate_module_graylog_contract.sh
validate_module_graylog_implementation.sh
validate_module_graylog_integration.sh
validate_module_graylog_acceptance.sh

validate_module_bloodhound_api_contract.sh
validate_module_bloodhound_api_implementation.sh
validate_module_bloodhound_api_integration.sh
validate_module_bloodhound_api_acceptance.sh
```

Optional module acceptance verifies that:

- the module can be disabled;
- the module can be absent;
- module failure does not stop core Atlas;
- credentials are isolated;
- data classification is enforced;
- retry and backpressure are bounded where applicable;
- failures are visible; and
- the core acceptance gate still passes without the module.

---

## Requirements for Every Implementation Gate

Where applicable, each implementation gate verifies:

### Boundary

- exact predecessor;
- exact candidate;
- declared scope;
- explicit nonclaims;
- prohibited behavior.

### Implementation

- required files;
- expected interfaces;
- dependency direction;
- configuration contract;
- migration or schema identity.

### Correctness

- positive cases;
- negative cases;
- edge cases;
- concurrency;
- cancellation;
- timeout;
- deterministic behavior.

### Security

- authorization;
- least privilege;
- secret handling;
- fail-closed behavior;
- hostile input;
- resource exhaustion;
- audit behavior.

### Evidence

- sanitized retained output;
- run identity;
- host and toolchain fingerprint;
- exact commit and tree;
- dependency versions;
- test totals;
- warnings;
- non-evaluated measurements.

### Documentation

- requirements;
- architecture;
- testing;
- traceability;
- limitations;
- roadmap status;
- changelog;
- gate index;
- manifests.

---

## Requirements for Every Formal Acceptance Gate

A formal acceptance gate verifies:

1. The exact implementation or integration gate still passes.
2. Every accepted predecessor passes at its exact historical commit.
3. The candidate exists in the canonical GitHub repository.
4. A clean canonical clone passes.
5. Retained evidence corresponds to the exact accepted commit.
6. Evidence digests are valid.
7. Requirements, architecture, implementation, tests, limitations, and status agree.
8. No undocumented capability is claimed.
9. No incomplete capability is silently accepted.
10. The accepted commit is tagged.
11. Next work and remaining nonclaims are recorded.
12. Self-validation is not represented as independent review.

## Gate-Creation Rule

Do not create all future executable gate scripts now.

For each phase:

1. create the phase contract;
2. create the first implementation gate;
3. complete and validate that capability;
4. create the next gate only when its exact boundary is understood;
5. create the integration gate after the bounded components exist; and
6. create the acceptance gate only after the full phase candidate exists.

The roadmap may list expected gates, but executable validators must represent real implemented or immediately active boundaries rather than speculative code structure.
