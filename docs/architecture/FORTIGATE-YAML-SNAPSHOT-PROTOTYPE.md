# FortiGate YAML Snapshot Prototype

## Status

Experimental side-by-side implementation candidate. This prototype is not an accepted phase boundary, production configuration importer, or replacement for the existing FortiGate native hierarchy parser.

## Purpose

The prototype tests whether a FortiGate-generated YAML backup can be transformed into the fuller Iron Atlas firewall snapshot model while preserving policy order, object relationships, routing intent, SD-WAN configuration, original IP representations, canonical CIDR networks, and reviewable evidence.

It is intentionally runnable as a standalone Go command on Windows, WSL, or Linux so development can continue independently from the active Iron Atlas service and authentication work.

## Pipeline

```text
FortiGate-generated YAML backup
        |
        v
Bounded FortiGate YAML document parser
        |
        v
Vendor-aware FortiGate normalizer
        |
        v
Experimental vendor-independent firewall snapshot
        |
        +-- summary output
        +-- deterministic JSON output
        +-- typed reference graph
        +-- reviewable findings
```

The native FortiOS parser remains unchanged. A future accepted boundary can place both adapters behind the same normalized snapshot contract after fixture coverage and version compatibility are proven.

## Normalized domains

The snapshot separates:

- Device and source metadata
- VDOM and routing-domain scope
- Interfaces, zones, VLANs, subinterfaces, IP assignments, and subnets
- DHCP servers, ranges, reservations, relays, and options
- Address, service, schedule, user, tag, and group objects
- Static routes, policy routes, routing objects, and configured dynamic-routing intent
- SD-WAN zones, members, health checks, rules, preference, and SLA relationships
- Ordered firewall policies and every supported interface, object, NAT, DSCP, and shaper reference
- Virtual IPs, IP pools, central SNAT, and policy NAT
- IPsec VPN configuration
- Traffic shapers, per-IP shapers, shaping policies, bandwidth, VLAN CoS, and DSCP values
- Typed reference edges, unresolved references, source locations, and findings

## Address representation

Iron Atlas preserves both the source representation and a canonical network representation.

For an interface value such as:

```text
10.120.0.1 255.255.255.0
```

the normalized record contains:

```text
assigned address: 10.120.0.1
original mask:    255.255.255.0
canonical subnet: 10.120.0.0/24
```

The assigned host address is not confused with the network prefix.

## References

Policies and other records retain typed references instead of copying object contents. Examples include:

```text
policy -> source interface
policy -> destination SD-WAN zone
policy -> address group
policy -> service group
policy -> virtual IP
policy -> IP pool
policy -> traffic shaper
static route -> interface or SD-WAN zone
SD-WAN member -> physical or tunnel interface
SD-WAN rule -> member and health check
VLAN interface -> parent interface
DHCP server -> interface
```

Resolved references produce graph edges. Missing or ambiguous references produce explicit unresolved-reference records and configured-state findings; they are never silently converted into empty values.

## Parser boundary

Fortinet documents its YAML backup as a consistent machine-readable representation. The prototype nevertheless treats the input as untrusted configuration evidence.

The bounded parser supports the structures used by the FortiGate-generated block format:

- Indented mappings
- Indented sequences
- Named sequence entries
- Single-quoted, double-quoted, and unquoted scalars
- Flow sequences such as `[wan1, wan2]`
- Comments and document markers

The parser rejects anchors, aliases, custom tags, flow mappings, block scalars, tab indentation, duplicate keys, malformed indentation, and unsupported ambiguous constructs. This keeps the initial dependency boundary small and fail-closed. Broader YAML support requires a separate reviewed dependency decision or retained sanitized evidence demonstrating the need.

## Security and evidence handling

Use password-masked and sanitized backups for development whenever possible.

Do not commit:

- Raw production firewall backups
- Credentials or database URLs
- Pre-shared keys or private keys
- Certificates containing private material
- Unredacted addresses, names, or customer evidence
- Generated normalized output from production configurations

The retained repository fixture uses documentation ranges and synthetic names only.

A backup establishes `CONFIGURED` intent. It does not establish live route installation, current SD-WAN health or member selection, active VPN state, interface link state, HA state, or current sessions. Those require separately collected operational evidence.

## Validation

The prototype test fixture covers:

- Physical and VLAN interfaces
- Original netmask and canonical CIDR conversion
- DHCP range and reservation parsing
- Address and service groups
- Static default and internal routes
- SD-WAN zones, members, SLA, and ordered rule
- Firewall policies and policy order
- Virtual IP and source NAT pool
- DSCP names and numeric values
- Traffic and per-IP shapers
- IPsec interface configuration
- Complete reference resolution
- Explicit unresolved-reference findings
- Rejection of unsupported YAML features

Run:

```bash
go test ./modules/firewall/fortigate ./modules/firewall/snapshot
go run ./cmd/fortigate-inspect \
  -input modules/firewall/fortigate/testdata/fortigate-sanitized.yaml \
  -format summary
```

## Acceptance limitations

Before this can become an accepted production-facing ingestion boundary, Iron Atlas still requires:

1. Sanitized fixtures from multiple FortiOS major and minor versions.
2. Single-VDOM and multi-VDOM fixtures from actual FortiGate YAML exports.
3. Coverage for IPv6, central NAT variants, local-in policies, HA, additional VPN forms, dynamic routing details, and nested object behavior.
4. Differential comparison against the same appliance configuration exported in native FortiOS and YAML formats.
5. Hostile-input, size-limit, malformed-indentation, duplicate-key, and resource-use validation.
6. A governed decision on retaining the bounded parser or adopting a reviewed YAML dependency.
7. Integration with evidence storage, snapshot persistence, authorization, and UI boundaries.
8. Clean-clone repository validation on the exact pushed candidate commit.
