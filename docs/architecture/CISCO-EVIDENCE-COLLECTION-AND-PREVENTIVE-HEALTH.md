# Cisco Evidence Collection and Preventive Health

## Initial Platforms

- Catalyst 2960
- Catalyst 2960-S
- Catalyst 2960-X
- Catalyst 9200
- Catalyst 9300
- Catalyst 9500
- Catalyst 9800 Wireless LAN Controller

## Collection Modes

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
- Catalyst 9800 HA, AP, WLAN, profile, tag, authentication, and client indicators

Cumulative counters are evaluated as deltas between comparable collections.

## Safety

Collections are staggered and bounded by device, stack, site, critical pair, and WLC HA pair. Only one comprehensive collection runs against a logical device at a time.
