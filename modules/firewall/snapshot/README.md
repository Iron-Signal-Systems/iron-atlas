# Firewall Snapshot Model

This package contains the experimental vendor-independent normalized firewall configuration model used by the FortiGate YAML side-by-side prototype.

The model keeps major domains separate while preserving typed references among policies, interfaces, zones, VLANs, subnets, routes, SD-WAN, objects, NAT, VPN, DHCP, DSCP, and QoS records. The reference graph distinguishes resolved edges, recognized built-ins, unresolved references, and ambiguous references; resolved edges retain their fixed object kind. Source locations and evidence state remain attached to normalized records and findings.

The package also exposes stable aggregate normalized-record counts for every fixed record kind, including zero-count kinds. These metrics support upload-safe semantic-quality reporting without inspecting or emitting source-derived names or values.

The package is additive and does not replace `modules/firewall/common` during the current accepted project phase.
