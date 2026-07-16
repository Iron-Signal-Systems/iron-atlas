# Firewall Snapshot Model

This package contains the experimental vendor-independent normalized firewall configuration model used by the FortiGate YAML side-by-side prototype.

The model keeps major domains separate while preserving typed references among policies, interfaces, zones, VLANs, subnets, routes, SD-WAN, objects, NAT, VPN, DHCP, DSCP, and QoS records. Source locations and evidence state remain attached to normalized records and findings.

The package is additive and does not replace `modules/firewall/common` during the current accepted project phase.
