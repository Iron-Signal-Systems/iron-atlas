# Cisco Collection Profile Catalog

## Profile Selection

The collector first establishes device identity, model, software, stack arrangement, and supported command behavior. It then selects a versioned profile. Unsupported commands are recorded; they are not silently ignored or treated as device failures.

Offline sanitized command bundles use the same profile and parser contracts that later restricted live collection uses.

## First-Value Profile Order

Development priority is:

1. Catalyst 9300L and Catalyst 9300 access switching
2. Catalyst 9500 core and distribution switching
3. Catalyst 9800 Wireless LAN Controllers
4. Catalyst 9200 compatibility
5. Catalyst 2960, 2960-S, and 2960-X compatibility

The priority reflects organizational visibility and first-product value, not a claim that older platforms are unimportant.

## Common Discovery

Candidate commands include:

```text
show version
show inventory
show switch
show module
show ip interface brief
show interfaces description
show interfaces status
show interfaces counters errors
show interfaces switchport
show interfaces trunk
show etherchannel summary
show vlan brief
show cdp neighbors detail
show lldp neighbors detail
show mac address-table dynamic
show spanning-tree
show power inline
show authentication sessions
show logging
```

Interface-specific running configuration should be collected after interface discovery, one approved interface at a time where practical.

## First-Slice Offline Bundles

A first-slice bundle should contain the minimum approved commands needed to prove:

- device and stack identity;
- interface and management identity;
- VLAN and switchport state;
- trunks, native VLANs, allowed and active VLANs, and pruning;
- CDP and LLDP topology;
- port-channel state;
- spanning-tree state;
- selected environmental and resource state;
- selected running-configuration semantics;
- Catalyst 9800 controller, HA, AP, WLAN, profile, and tag relationships; and
- parser handling for unsupported, partial, truncated, malformed, and conflicting output.

Each command record retains command identity, profile version, device identity, acquisition time, completion state, exit or transport state where applicable, output digest, truncation state, and classification.

## Monthly Comprehensive Profiles

### IOS XE 9300L/9300/9500

Base technical support plus applicable stack, StackWise Virtual, port, port-channel, PoE, platform, diagnostic, license, environment, crash, and resource packages. Feature-specific packages are conditional, not blindly executed on every device.

### Catalyst 9800

General IOS XE and wireless evidence, controller and HA state, AP inventory and join, WLAN and policy profiles, tags, RF/site/flex profiles, management interfaces, RADIUS dependencies, client aggregates, radio state, failures, resources, logs, and certificate-risk evidence where available.

### IOS XE 9200

Compatibility profile derived from the accepted IOS XE first-slice contracts, with platform-specific command and feature differences recorded explicitly.

### IOS 2960 Family

Compatibility profile for general technical support, identity, inventory, boot, interface configuration and counters, VLAN/trunk, spanning tree, EtherChannel, CDP/LLDP, MAC, PoE, stack where supported, environment, logs, reload, and failure evidence.

## Light Health Profiles

Daily or weekly profiles should be fast and low impact. Collect environment, stack/HA, interface error deltas, CPU, memory, uptime, reload reason, PoE, neighbor, AP join, and critical log summaries.

## Complementary-System Outputs

Normalized profile results may support:

- Zabbix host, interface, template, low-level discovery, map, dashboard, and report reconciliation;
- Graylog lookup and enrichment records for syslog and SNMP traps;
- Graylog query, pipeline, stream, dashboard, and report suggestions;
- Security-platform asset and topology context;
- Draw.io-compatible generated topology; and
- Atlas operational and acceptance reports.

Generated output remains distinct from applied external-system configuration.

## Safety

- Offline evidence precedes restricted live collection
- One comprehensive collection per logical device or stack
- One collection per critical redundant pair or WLC HA pair at a time
- Stagger by site and platform
- Stream output rather than write flash by default
- Bound command time, inactivity time, output size, and retries
- Never enter configuration mode or enable debugging
- Mark truncation and partial evidence explicitly
