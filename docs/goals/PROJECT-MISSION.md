# Project Mission

Build a trustworthy infrastructure-intelligence, documentation, integration, and change-governance platform that helps network technicians, network administrators, network-security staff, reviewers, auditors, and infrastructure teams understand:

- what infrastructure exists;
- how users, devices, services, networks, and security boundaries are connected;
- what evidence supports that understanding;
- what changed, why it changed, and who approved it;
- whether monitoring, logging, security, vendor, and documentation systems agree;
- where records are missing, stale, conflicting, incomplete, or uncertain; and
- whether the resulting environment was validated and formally accepted.

The interface should support the work rather than become additional work.

## Purpose

Iron Atlas is an authoritative, version-controlled system for infrastructure documentation, diagrams, inventory, approved discovery, project tracking, change management, validation, preventive health analysis, external-system reconciliation, integration assistance, and formal acceptance.

Atlas complements systems that already serve an organization. It does not recreate mature monitoring, logging, network-security, diagramming, or vendor-management capabilities merely to own them.

Atlas shall maintain a governed relationship among:

- Current infrastructure state
- Imported and collected evidence
- Evidence source, time, digest, parser version, and collection context
- Normalized inventory
- Switch ports, trunks, neighbors, VLANs, spanning tree, port channels, and access paths
- Wireless controllers, access points, profiles, tags, and selected client context
- Firewall policy, routing, object, and traffic-boundary behavior
- Zabbix hosts, templates, discovery, maps, dashboards, and reporting context
- Graylog syslog and SNMP-trap records, lookups, pipelines, streams, searches, dashboards, and reports
- Security Onion and other security-platform asset and topology context
- Generated and curated Draw.io diagrams
- Projects and proposed target states
- Approved changes
- Pre-change and post-change validation
- Documentation reconciliation
- Operational findings and human disposition
- Formal acceptance records

## Complementary-System Principle

Atlas may consume approved evidence from an external system and may produce reviewable context, definitions, recommendations, exports, maps, dashboards, queries, lookup data, templates, and reports that improve that system.

The external system retains responsibility for its mature operational purpose. Atlas retains responsibility for its canonical evidence, normalized records, topology, governed findings, change history, and acceptance history.

No future write or provisioning adapter may silently change an external system. Such an adapter requires a separately accepted boundary with preview, attribution, authorization, bounded scope, idempotency where practical, reversibility where practical, validation, and failure isolation.

## First Operational Value

The first infrastructure-value slice prioritizes Catalyst 9300L/9300 access switching, Catalyst 9500 core and distribution switching, and Catalyst 9800 wireless controllers because organizational users, devices, access points, phones, servers, VLANs, and network paths depend on that switching and wireless fabric.

The first slice should turn sanitized offline Cisco command bundles into normalized inventory, interface, VLAN, trunk, neighbor, port-channel, spanning-tree, wireless, and topology records, then use those records to improve documentation, Zabbix reconciliation, Graylog context, maps, dashboards, queries, and operational reports.

Restricted live collection follows only after the offline evidence, provenance, parser, normalization, resource, and reporting boundaries are accepted.

## Leading Principle

> Document the environment. Strengthen the tools already in place. Govern the change. Validate the result. Preserve the record.
