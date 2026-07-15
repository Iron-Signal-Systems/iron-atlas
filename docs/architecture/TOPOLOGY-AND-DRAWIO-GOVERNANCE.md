# Topology and Draw.io Governance

## Record Types

- Machine-readable inventory contains objective facts.
- Markdown explains purpose, ownership, dependencies, risks, and decisions.
- `.drawio` is the editable source for curated visual composition.
- SVG and PDF are published representations.
- Generated topology is rebuilt from normalized records and never overwrites curated diagrams.

## Diagram IDs

Examples:

- `ARCH-001` — Atlas system context
- `NET-LOG-001` — Logical topology
- `NET-PHY-001` — Physical topology
- `NET-VLAN-001` — VLAN and subnet relationships
- `NET-ROUTE-001` — Routing
- `NET-SEC-001` — Security zones
- `NET-WAN-001` — WAN and SD-WAN
- `NET-WLAN-001` — Wireless
- `RACK-001` — Rack elevation
- `PRJ-2026-001` — Project target state

## Diagram Contract

Every accepted curated diagram has:

- Identifier
- Editable `.drawio` source
- Companion Markdown record
- Published SVG
- PDF when printing or approval requires it
- Status, classification, owner, review date, and evidence references
- Explicit current-state or target-state designation

## Style

- Consistent flow direction
- Visible legend
- Explicit trust boundaries
- Distinguishable physical and logical connections
- No unexplained colors
- No secrets
- Inventory IDs on governed objects
- Dashed planned or conditional paths
- Readable grayscale and print behavior
