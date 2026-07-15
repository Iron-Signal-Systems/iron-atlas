# Firewall Configuration Semantic Analysis

## Initial Platforms

- FortiGate native FortiOS configuration
- FortiGate YAML when the source version supports it
- OPNsense XML configuration backup
- pfSense XML configuration backup

## Required Understanding

The firewall module must reconstruct relationships among:

- Physical and logical interfaces
- VLANs and zones
- Virtual routing domains and FortiGate VDOMs
- Connected, static, default, policy, and dynamic-routing configuration
- Gateways, gateway groups, and tunnel interfaces
- SD-WAN members, zones, rules, health checks, and performance objectives
- Firewall policies and evaluation order
- Address, service, schedule, user, and security-profile objects
- Source NAT, destination NAT, port translation, and address pools
- VPN configuration and route relationships
- Management and local-in exposure
- High-availability configuration

## Evidence States

Every conclusion is marked:

- `CONFIGURED`
- `OBSERVED`
- `CALCULATED`
- `INFERRED`
- `UNKNOWN`
- `CONFLICTING`

A backup proves configured intent. It does not prove current gateway health, installed dynamic routes, active SD-WAN member selection, live VPN state, interface link state, or current sessions.

## Traffic-Path Query

Given source interface, source address, destination, protocol, port, identity, and time context, the analyzer should explain:

1. Routing domain or VDOM
2. Ingress interface and zone
3. Policy order and matching objects
4. Policy route or gateway selection
5. SD-WAN rule and member preference
6. Longest-prefix route
7. Egress zone or interface
8. NAT
9. Security inspection
10. Logging
11. Return-path uncertainty
12. Evidence supporting each conclusion

## Required Findings

- Missing object references
- Empty groups
- Route to disabled or missing interface
- Unreachable next hop
- Policy without route
- Route without usable policy
- Shadowed or overly broad rules
- Unused or orphaned objects
- Management exposure on external interfaces
- SD-WAN member or SLA inconsistencies
- Default route bypassing intended SD-WAN
- NAT or return-path risk
- Unexplained drift from last accepted state

Findings begin as reviewable evidence, not automatic declarations that the configuration is wrong.
