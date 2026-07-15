# Cisco Collection Profile Catalog

## Profile Selection

The collector first establishes device identity, model, software, stack arrangement, and supported command behavior. It then selects a versioned profile. Unsupported commands are recorded; they are not silently ignored or treated as device failures.

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

## Monthly Comprehensive Profiles

### IOS 2960 Family

General technical support, identity, inventory, boot, interface configuration and counters, VLAN/trunk, spanning tree, EtherChannel, CDP/LLDP, MAC, PoE, stack where supported, environment, logs, reload, and failure evidence.

### IOS XE 9200/9300/9500

Base technical support plus applicable stack, StackWise Virtual, port, port-channel, PoE, platform, diagnostic, license, environment, crash, and resource packages. Feature-specific packages are conditional, not blindly executed on every device.

### Catalyst 9800

General IOS XE and wireless evidence, controller/HA state, AP inventory and join, WLAN and policy profiles, tags, RF/site/flex profiles, management interfaces, RADIUS dependencies, client aggregates, radio state, failures, resources, logs, and certificate-risk evidence where available.

## Light Health Profiles

Daily or weekly profiles should be fast and low impact. Collect environment, stack/HA, interface error deltas, CPU, memory, uptime, reload reason, PoE, neighbor, AP join, and critical log summaries.

## Safety

- One comprehensive collection per logical device or stack
- One collection per critical redundant pair at a time
- Stagger by site and platform
- Stream output rather than write flash by default
- Bound command time, inactivity time, output size, and retries
- Never enter configuration mode or enable debugging
- Mark truncation and partial evidence explicitly
