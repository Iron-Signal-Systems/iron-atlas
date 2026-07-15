# HTML5 Interface and Role Workspaces

## Purpose

Provide an accessible, responsive HTML5 interface that complements the work of network technicians, network administrators, network-security staff, reviewers, auditors, and infrastructure teams.

## Governing Rules

- The role and its work come first.
- Common work has a clear and direct path.
- Stale, incomplete, uncertain, queued, failed, conflicting, and accepted states remain distinguishable.
- Visibility of a button or record does not grant authority.
- The interface does not independently create canonical truth.
- Success is not shown before authoritative confirmation.
- Keyboard operation, semantic HTML, readable focus state, and screen-reader support are functional requirements.
- The initial accessibility target is WCAG 2.1 Level AA, with a roadmap to evaluate newer applicable requirements.

## Role Workspaces

### Network Technician

- Device and port search
- Current endpoint attachment evidence
- Access-port configuration and health
- CDP/LLDP neighbors
- Open findings and assigned remediation
- Proposed change creation
- Collection status

### Network Administrator

- Routing, trunk, VLAN, spanning-tree, port-channel, firewall, and wireless views
- Project and change planning
- Pre-change and post-change comparison
- Approval where independent and authorized
- Configuration standards and drift
- Collection scheduling and platform coverage

### Network Security

- Firewall policy and NAT analysis
- ACL attachment and semantic review
- Management-plane exposure
- AAA/NPS/RADIUS boundaries
- Security-control drift
- High-risk change approval
- Audit and evidence review

### Reviewer and Auditor

- Decision and approval history
- Evidence integrity and parser lineage
- Acceptance records
- Exceptions and supersession
- Read-only historical reconstruction

### Team View

- Shared project and change queues
- Findings by site and severity
- Upcoming collection and change windows
- Blocked work
- Documentation debt
- Accepted work and unresolved exceptions

## Implementation Direction

The initial UI uses Go `html/template`, embedded static assets, semantic HTML5, and minimal JavaScript. A large single-page application framework is not required for the first implementation and must not be introduced without a demonstrated operational benefit.
