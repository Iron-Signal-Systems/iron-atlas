# System Requirements

## Foundation

- `IA-FND-001`: All production capabilities shall be modular behind versioned contracts.
- `IA-FND-002`: Core canonical records shall not depend on one vendor or monitoring product.
- `IA-FND-003`: Unsupported and uncertain parser results shall remain visible.
- `IA-FND-004`: Raw evidence shall remain outside Git and receive integrity protection.
- `IA-FND-005`: Collectors shall not connect directly to PostgreSQL.

## Interface and Access

- `IA-UI-001`: The system shall provide an accessible HTML5 interface.
- `IA-UI-002`: Role workspaces shall support network technicians, administrators, security, reviewers, and teams.
- `IA-AUTH-001`: UI visibility shall not grant authority.
- `IA-AUTH-002`: Production authentication and authorization shall fail closed.
- `IA-AUTH-003`: Platform administration shall not automatically grant change approval.

## Change Management

- `IA-CHG-001`: A requester shall not approve the requester’s own governed change.
- `IA-CHG-002`: Approval shall retain actor, authority, stage, scope, decision, reason, and time.
- `IA-CHG-003`: High-impact changes may require multiple independent approvers.
- `IA-CHG-004`: A change shall not be accepted before post-change validation and documentation reconciliation.
- `IA-CHG-005`: Material history shall use correction and supersession rather than silent rewriting.

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
