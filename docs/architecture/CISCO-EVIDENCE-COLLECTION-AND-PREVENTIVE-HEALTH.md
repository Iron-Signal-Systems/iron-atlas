# Cisco Evidence Collection and Preventive Health

## First-Value Platforms

The first infrastructure-value slice prioritizes:

- Catalyst 9300L and Catalyst 9300 access switching
- Catalyst 9500 core and distribution switching
- Catalyst 9800 Wireless LAN Controllers

These platforms provide broad organizational visibility because users, devices, phones, servers, access points, VLANs, and network paths depend on the access, core/distribution, and wireless fabric.

## Compatibility Platforms

After the first-value slice is established, profiles may extend or retain compatibility for:

- Catalyst 2960
- Catalyst 2960-S
- Catalyst 2960-X
- Catalyst 9200

Compatibility support must not dilute completion of the first 9300L/9300, 9500, and 9800 end-to-end slice.

## Evidence Sequence

Offline evidence precedes restricted live collection.

The first accepted Cisco path uses sanitized, approved command bundles to establish:

- evidence receipt, digest, provenance, classification, and parser version;
- device, model, software, stack, and command-profile detection;
- deterministic parsing and normalization;
- visible unsupported, truncated, conflicting, and partial state;
- resource, timeout, cancellation, and failure behavior;
- manually verified inventory and topology;
- Zabbix reconciliation;
- Graylog lookup, query, dashboard, and report context;
- Draw.io-compatible topology generation; and
- operational report usefulness.

Restricted live collection may follow only through a separately accepted boundary.

## Collection Modes

- Offline sanitized command-bundle import
- Scheduled 30-day comprehensive baseline
- Lighter daily or weekly health collection
- Pre-change collection
- Post-change collection
- Incident collection
- Targeted diagnostics
- Replacement acceptance
- Software-upgrade validation

A failed collection does not reset the due date.

## Thirty-Day Evidence

The platform selects a versioned platform-specific profile after detecting model, software, stack arrangement, and supported commands. The base report is the supported equivalent of `show tech-support`, augmented only by applicable stack, port, PoE, platform, diagnostic, licensing, and wireless reports.

Output streams over the authenticated SSH session. The default profile does not write temporary files to switch flash, enter configuration mode, enable debugging, reload a process, reload a device, or modify configuration.

## Service Authentication

Support Windows Server NPS/RADIUS with Active Directory for authentication and initial Cisco privilege assignment. Use a dedicated nonhuman account, dedicated AD group, dedicated NPS device-management policy, source restrictions, accounting, and a restricted local Cisco parser view or privilege level.

NPS/RADIUS must not be described as TACACS+-equivalent per-command authorization. TACACS+ remains an optional future control when centralized command authorization and command accounting are required.

## Normalized First-Slice Records

The first slice should normalize:

- device identity, model, serial number, software, boot, uptime, and license context;
- stack and StackWise Virtual membership and state;
- interfaces, descriptions, administrative state, operational state, errors, and selected counters;
- VLANs, voice VLANs, switchport mode, trunks, native VLANs, allowed and active VLANs, and pruning;
- CDP and LLDP neighbors;
- port channels and member consistency;
- spanning-tree instances, roots, roles, states, and topology-change observations;
- IP interface summaries and management identity;
- Catalyst 9800 controller and HA state;
- access-point inventory and join state;
- WLAN, policy, site, flex, and tag relationships where evidence supports them; and
- explicit unsupported, incomplete, conflicting, and uncertain output.

## Preventive Analysis

Analyze and trend:

- Running and accepted configuration differences
- AAA, SSH, SNMP, NTP, logging, and management access
- VLAN, voice VLAN, trunk, native VLAN, and port-channel state
- Spanning-tree root, roles, inconsistencies, and topology change
- Interface errors, CRC deltas, discards, resets, and link transitions
- Power supplies, fans, temperature, PoE budget, flash, and filesystem
- Stack and StackWise Virtual health
- Software, boot, license, uptime, reload, crash, CPU, and memory state
- CDP/LLDP neighbor and topology changes
- Catalyst 9800 HA, AP, WLAN, profile, tag, authentication, and selected client indicators
- Zabbix identity, monitoring, template, discovery, map, and documentation discrepancies
- Graylog context opportunities for syslog and SNMP-trap records

Cumulative counters are evaluated as deltas between comparable collections.

## Complementary Outputs

The Cisco slice may generate reviewable:

- normalized inventory and topology reports;
- Draw.io-compatible topology sources;
- Zabbix host and interface reconciliation;
- Zabbix map, dashboard, template, and discovery recommendations;
- Graylog lookup data and infrastructure-enrichment context;
- Graylog query, stream, pipeline, dashboard, and report suggestions; and
- documentation discrepancy and preventive-health findings.

These outputs do not imply that external-system configuration has been applied.

## Live-Collection Safety

Live collection requires:

- explicit authorization;
- pinned SSH host keys;
- restricted service authentication;
- fixed versioned command profiles;
- no configuration mode;
- no debugging;
- bounded concurrency;
- per-command and per-device timeouts;
- cancellation and stop controls;
- schedule jitter;
- output-size bounds;
- protected transcripts;
- complete evidence provenance; and
- one comprehensive collection against a logical device, stack, critical pair, or WLC HA pair at a time.
