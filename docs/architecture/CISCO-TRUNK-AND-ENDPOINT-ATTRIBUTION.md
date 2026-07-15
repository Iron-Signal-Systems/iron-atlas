# Cisco Trunk and Endpoint Attribution

## Governing Rule

A trunk or infrastructure interconnect is excluded only from local endpoint attribution. It remains fully documented, analyzed, compared, validated, and diagrammed.

MAC addresses learned through trunks are not claimed as locally connected endpoints. They may be used for path tracing, MAC movement, loop investigation, and downstream topology correlation.

## Infrastructure Classification

Exclude endpoint attribution for:

- Static, negotiated, or operational trunk
- Port-channel member
- Port-channel interface
- Routed interface
- Stack interface
- StackWise Virtual link
- Fabric or infrastructure interconnect
- Firewall, router, controller, hypervisor, or access-point uplink
- Explicitly excluded port

## Required Trunk Analysis

For every trunk or infrastructure link, collect and reconcile:

- Interface description
- Administrative and operational switchport mode
- DTP and `switchport nonegotiate`
- Native VLAN
- Configured, effective, active, forwarding, blocked, and pruned VLANs
- Allowed-list additions and removals
- VTP pruning where applicable
- Port-channel membership, protocol, state, and member consistency
- Speed, duplex, MTU, UDLD, flow control, and link settings
- DHCP snooping and DAI trust
- QoS trust and service policies
- SPAN source/destination role
- CDP and LLDP neighbor identity, platform, management address, and remote interface
- Reciprocal observations when both devices are managed
- Spanning-tree mode, role, state, cost, priority, root, protections, and inconsistencies
- IPv4, IPv6, MAC, port, router, and VLAN ACL attachments and direction
- Interface error and discard deltas
- Change and accepted-baseline relationships

## Endpoint Path

Follow a learned MAC through infrastructure links until the most specific evidence-supported non-trunk edge interface is reached. Only that validated access interface may be used for physical endpoint attribution.

## Findings

- Missing or mismatched description
- Missing or unexpected neighbor
- Native VLAN mismatch
- Allowed VLAN mismatch
- Required VLAN missing or unnecessary VLAN present
- Unbounded all-VLAN trunk
- Pruning mismatch
- Unexpected spanning-tree role or block
- Root, loop, or native-VLAN inconsistency
- Port-channel member suspension or mismatch
- Speed, duplex, MTU, CRC, UDLD, or error condition
- DTP enabled contrary to standard
- ACL missing, changed, or attached in the wrong direction
- Unapproved topology movement
