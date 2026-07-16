# External-System-Independent Telemetry

## Decision

Iron Atlas owns canonical health, metric, finding, evidence, delivery-intent, reconciliation, and operational records independently of Zabbix, Graylog, Security Onion, or any other external product.

An external product may remain the primary operational interface for the capability it performs well. It is not Atlas's authorization source, canonical schema, infrastructure-evidence authority, governed-change authority, or formal-acceptance history.

## Complementary Direction

Atlas integrates to strengthen systems already serving an organization.

- Zabbix remains responsible for continuous monitoring, alerting, graphing, escalation, maintenance, availability history, and proxy behavior.
- Graylog remains responsible for log and SNMP-trap collection, indexing, retention, search, investigation, streams, pipelines, and log dashboards.
- Security Onion and similar platforms remain responsible for network-security monitoring, packet analysis, detection, and investigation.
- Atlas provides normalized infrastructure identity, site, interface, VLAN, topology, evidence age, findings, and documentation context that those systems may consume.

Atlas may also ingest approved metadata from those systems for reconciliation. Imported metadata retains its source and does not silently become canonical truth.

## Initial Adapters

- Native Go Zabbix sender-protocol adapter
- Future Zabbix API adapter for governed reconciliation and provisioning where justified
- Graylog/syslog and SNMP-trap context exports
- Future Graylog lookup, pipeline, stream, query, dashboard, and report-definition adapter
- OpenMetrics-compatible endpoint
- Syslog
- Webhook
- SIEM and security-platform context delivery

## Delivery Contract

Each delivery contains:

- Integration contract version
- Destination ID
- Canonical metric, event, record, or artifact ID
- Device, site, module, interface, VLAN, topology, and finding context as applicable
- Value or generated artifact and timestamp
- Evidence and confidence references where applicable
- Classification
- Attempt count and next attempt
- Delivery result
- Applied-state status when the destination supports governed provisioning

Retries are bounded with backpressure and dead-letter handling. Destination failure does not block collection, parsing, change management, canonical recording, or unrelated integrations.

## Generated and Applied State

Atlas distinguishes:

- observed external state;
- imported external metadata;
- generated recommendation;
- exported definition;
- approved provisioning request;
- applied target-system state; and
- post-application validation result.

A generated map, dashboard, query, pipeline, template, lookup table, or report definition must not be represented as applied until the target system confirms it and Atlas records the result through an accepted integration boundary.

## Zabbix Direction

The initial implementation sends values to existing Zabbix trapper items using the sender protocol directly from Go, avoiding a required `zabbix_sender` package.

Later versioned work may add:

- host and interface reconciliation;
- template and low-level discovery recommendations;
- topology-derived maps;
- dashboard and report definitions;
- inventory and maintenance recommendations;
- TLS profiles;
- governed API provisioning; and
- applied-state validation.

Atlas should not recreate mature Zabbix polling, alerting, graphing, escalation, maintenance, proxy, and availability capabilities without a separately justified and accepted decision.

## Graylog Direction

Atlas may provide normalized infrastructure lookup and enrichment data for Graylog syslog and SNMP-trap records. It may generate reviewable searches, queries, pipeline and stream definitions, dashboards, and reports that use device, site, interface, VLAN, wireless, and topology context.

Graylog remains the log-retention and investigation authority.

## Governed Write Boundary

Any adapter that changes an external system requires separate acceptance, least privilege, explicit target identity, previewable differences, actor attribution, applicable independent approval, bounded scope, idempotency where practical, reversal where practical, and post-application validation.
