# Operational-System Complement and Integration Model

## Status

Normative target direction. This document defines product position and future integration boundaries. It does not claim that any unimplemented external-system adapter, provisioning workflow, Cisco evidence pipeline, or production integration has been accepted.

## Decision

Iron Atlas complements established operational systems instead of replacing them merely to duplicate mature capabilities.

Atlas provides authoritative infrastructure evidence, normalized identity, topology context, documentation reconciliation, governed findings, change comparison, integration assistance, and formal acceptance history.

External systems continue to own their established operational functions:

- Zabbix owns continuous polling, availability, performance history, triggers, alerting, graphing, maintenance, escalation, and proxy behavior.
- Graylog owns centralized syslog and SNMP-trap intake, indexing, retention, search, investigation, streams, pipelines, and log-platform dashboards.
- Security Onion and similar platforms own network-security monitoring, packet analysis, detections, hunting, and security investigation.
- Cisco, Fortinet, and other infrastructure platforms own device operation, configuration enforcement, forwarding, access, and security policy execution.
- Draw.io remains a supported human-editable diagram source and publication format.

- Iron File Intelligence owns file identity, access, classification, file
  activity, audit coverage, and detailed forensic lineage.

Atlas does not become a monitoring platform, log platform, detection engine, network controller, or vendor manager merely because it integrates with one.

Atlas also does not become a duplicate file-intelligence platform or a
generic single pane of glass. Its responsibility is cross-system correlation,
reachability, dependency, change-impact, and incident-impact explanation.

## Integration Directions

### External systems into Atlas

Approved external evidence may help Atlas reconcile infrastructure state. Examples include:

- Zabbix host, interface, template, item, trigger, proxy, map, maintenance, and inventory records;
- Graylog syslog, SNMP-trap, lookup, stream, pipeline, and search metadata;
- Security Onion asset or investigation references that an authorized integration contract permits;
- Cisco and Fortinet command output, technical-support bundles, configuration exports, inventory, and health evidence; and
- curated Draw.io sources and approved documentation records.

- signed, versioned, minimized Iron File Intelligence context bundles and
  authorized IFI evidence references.

Imported records retain source identity, acquisition time, integrity information where available, classification, adapter version, and confidence. An imported external record does not silently become authoritative Atlas truth.

### Atlas into external systems

Atlas may generate or deliver reviewable artifacts that improve an external system, including:

- Zabbix host reconciliation, template suggestions, low-level discovery definitions, maps, dashboard definitions, report inputs, trapper metrics, and maintenance or inventory recommendations;
- Graylog lookup tables, enrichment data, pipeline and stream suggestions, searches, queries, dashboards, report definitions, and device or interface context for syslog and SNMP traps;
- Security Onion asset, site, role, expected-relationship, and topology context;
- Draw.io-compatible generated topology sources; and
- vendor-neutral reports and exports for human review.

The default integration mode is read, reconcile, generate, recommend, export, or deliver. It is not silent external-system modification.

## Zabbix Complement

Atlas should use Cisco and other infrastructure evidence to reconcile:

- device identity, model, serial number, software, stack, and site;
- management interfaces and addresses;
- monitored versus unmonitored devices;
- stale or duplicated hosts;
- expected templates and low-level discovery coverage;
- topology relationships suitable for Zabbix maps;
- dashboard and report context;
- collection freshness and evidence age; and
- findings that should be delivered through existing Zabbix operations.

Zabbix remains the continuous monitoring and alerting authority. Atlas remains the infrastructure-evidence and reconciliation authority.

## Graylog Complement

Atlas should use normalized device, site, interface, VLAN, neighbor, port-channel, wireless, and topology context to improve Graylog interpretation of:

- SNMP traps;
- syslog messages;
- authentication and authorization events;
- interface transitions and errors;
- spanning-tree and topology changes;
- wireless-controller and access-point events;
- configuration and administrative events; and
- infrastructure incident searches.

Atlas may generate lookup data, query examples, pipeline and stream definitions, dashboards, and reports. Graylog remains the log-retention and investigation authority.

## Security Onion Complement

Atlas may provide governed asset identity, network role, site, expected communication relationships, and topology context. Security Onion remains responsible for detection, packet analysis, alerting, and investigation.

Atlas must not represent infrastructure documentation as proof that observed network behavior is safe.

## Vendor-System Boundary

Cisco, Fortinet, and other vendor systems remain responsible for operation and enforcement.

The first Atlas product boundary is read-only assessment, documentation, reconciliation, reporting, and integration assistance. Automated configuration deployment and remediation remain deferred.

## Governed Write and Provisioning Boundary

A future adapter that writes or provisions an external system requires a separately implemented and accepted boundary that provides:

- explicit target-system identity;
- least-privileged service authentication;
- bounded object and action scope;
- human-readable preview and difference display;
- actor attribution and authorization;
- applicable independent approval;
- idempotency where practical;
- rollback or reversal where practical;
- precondition and postcondition validation;
- audit and delivery history;
- retry, backpressure, dead-letter, and duplicate controls; and
- fail-closed behavior when target state is ambiguous.

A recommendation or generated definition must not be represented as applied configuration.

## First Infrastructure-Value Slice

The first infrastructure-value slice prioritizes:

1. Catalyst 9300L/9300 access switching;
2. Catalyst 9500 core and distribution switching; and
3. Catalyst 9800 wireless controllers.

These platforms provide broad organizational visibility because users, devices, phones, access points, servers, VLANs, and network paths depend on them.

The first accepted slice should proceed through:

```text
sanitized Cisco command bundle
        ↓
evidence receipt, digest, provenance, and classification
        ↓
platform and command-profile detection
        ↓
bounded parsers with explicit unsupported and partial state
        ↓
normalized device, stack, interface, VLAN, trunk, neighbor,
port-channel, spanning-tree, wireless, and topology records
        ↓
Zabbix reconciliation and Graylog context generation
        ↓
Draw.io maps, dashboards, queries, and operational reports
        ↓
human verification, disposition, and accepted evidence
```

Offline evidence must precede restricted live collection. Live collection requires pinned SSH host keys, fixed read-only command profiles, bounded concurrency, per-command and per-device timeouts, cancellation, protected transcripts, complete provenance, and explicit operational authorization.

## Fortinet Position

Fortinet remains an important subsequent vertical slice. It adds firewall objects, services, policy order, routes, gateways, NAT, VPN, SD-WAN, zones, and traffic-boundary context.

Cisco primarily establishes what exists and how it is connected. Fortinet adds what is permitted across significant boundaries. Atlas should correlate both without requiring either vendor representation to become the canonical model.

## Non-Negotiable Boundaries

- Canonical Atlas records remain independent of an external product schema.
- External-system failure does not erase canonical Atlas evidence or block unrelated core operations.
- Imported records retain provenance and uncertainty.
- Generated definitions remain distinguishable from applied target-system state.
- Atlas never silently grants itself authority in an external system.
- UI visibility is not authorization.
- Unaccepted write integration is prohibited.
- Replacement of mature operational products is not a first-product objective.
