# Zabbix Integration Contract

## Purpose

Complement the existing Zabbix environment with Atlas infrastructure evidence, normalized identity, topology, reconciliation, findings, and generated operational artifacts without making Zabbix the canonical Atlas database or making Atlas a replacement monitoring platform.

## Ownership

Atlas owns:

- source evidence and provenance;
- normalized infrastructure identity and relationships;
- topology and documentation context;
- finding state and confidence;
- change and acceptance history;
- reconciliation results;
- generated recommendations and artifact history; and
- delivery and provisioning-attempt history.

Zabbix owns:

- item and trend history;
- polling and proxy behavior;
- availability and performance monitoring;
- triggers and event processing;
- alerting and escalation;
- maintenance;
- graphs;
- operational dashboards and maps after they are accepted and applied in Zabbix; and
- Zabbix-side permissions and configuration state.

## Initial Sender Metrics

Examples:

```text
atlas.service.health
atlas.collection.last_success_unixtime
atlas.collection.overdue
atlas.collection.partial
atlas.device.findings.critical
atlas.device.findings.high
atlas.device.findings.moderate
atlas.device.config_drift
atlas.device.evidence_age_seconds
atlas.integration.delivery_failures
```

## Reconciliation

Atlas should compare approved Cisco and other infrastructure evidence with selected Zabbix records, including:

- host identity and duplicate state;
- management address and interface identity;
- device model, serial number, software, stack, and site;
- expected versus present monitoring;
- stale or missing hosts;
- expected template coverage;
- low-level discovery coverage;
- proxy and site alignment where authorized; and
- topology relationships suitable for maps and reports.

A mismatch is a governed reconciliation finding, not automatic proof that either system is wrong.

## Generated Assistance

Atlas may generate reviewable:

- host and interface proposals;
- template and low-level discovery recommendations;
- topology-derived maps;
- dashboard definitions;
- report inputs and report definitions;
- inventory updates;
- maintenance recommendations; and
- trapper-item and finding-delivery definitions.

Generated artifacts remain distinct from applied Zabbix configuration.

## Delivery

The initial Go adapter implements the Zabbix sender protocol for pre-created trapper items.

Later work may add TLS profiles, host and item lifecycle management, template provisioning, low-level discovery, map and dashboard provisioning, and API-based reconciliation through separate versioned adapters.

## Governed Provisioning

Any Zabbix API write requires a separately accepted boundary with:

- least-privileged API identity;
- explicit target Zabbix instance;
- previewable proposed differences;
- actor attribution and authorization;
- applicable independent approval;
- bounded object types and actions;
- idempotency where practical;
- rollback or reversal where practical;
- post-write readback and validation;
- retry, duplicate, and dead-letter controls; and
- audit history.

## Failure

Failed delivery or provisioning enters bounded retry and eventually dead-letter state. It does not discard canonical Atlas evidence, findings, or reconciliation records and does not block collection or analysis.

## Replacement

The same canonical event, record, or context can be delivered to OpenMetrics, Graylog/syslog, webhooks, SIEM destinations, or future adapters without changing the source schema.
