# Data and Record Model

## Canonical Record Families

- Organization, site, building, room, rack, and location
- Device, component, software, license, and lifecycle
- Interface, VLAN, subnet, zone, routing domain, circuit, and tunnel
- Neighbor, link, port-channel, spanning-tree, MAC observation, and endpoint attribution
- Firewall object, policy, route, NAT, SD-WAN, VPN, and traffic path
- WLC, AP, WLAN, profile, tag, radio, and client observation
- Evidence bundle, artifact hash, parser run, warning, and lineage
- Finding, observation, trend, exception, remediation, and risk
- Project, change, approval, validation, acceptance, and supersession
- Canonical metric, health event, delivery destination, outbox, and delivery attempt

## Identity

Every durable record receives a stable identifier. Vendor IDs and names are attributes, not the sole canonical identifier.

## Time

Distinguish:

- Source-device time
- Collector-received time
- Ingestion time
- Parser time
- Observation-valid time
- Acceptance time
- Supersession time

## History

Material changes retain lineage. Accepted state is derived from governed records; it is not silently overwritten by the newest collection.

## Storage

PostgreSQL is the planned authoritative store for normalized and governed records. Raw evidence remains in protected content-addressed storage and is referenced by immutable hashes and storage references.
