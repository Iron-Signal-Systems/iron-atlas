# FortiGate YAML Snapshot Prototype

## Status

Experimental side-by-side implementation candidate. This prototype is not an accepted phase boundary, production configuration importer, or replacement for the existing FortiGate native hierarchy parser.

The maintained-decoder change replaces only the YAML grammar boundary. The FortiGate normalizer, vendor-independent firewall snapshot, evidence semantics, and native FortiOS adapter remain separate and intact.

A private, upload-safe V8 execution against a 13.2 MB FortiOS 7.2.13 export has now passed maintained decoding, Atlas admission, native direct-root layout detection, normalization, focused and full tests, race testing, vet, build, documentation validation, evidence validation, manifest validation, and checksum validation. That demonstration establishes a real working parser path, not semantic acceptance: 263 root entries remain outside recognized normalization coverage, and the first normalized run reported 1,597 unresolved references and 1,605 findings requiring aggregate classification.

## Purpose

The prototype tests whether a FortiGate-generated YAML backup can be transformed into the fuller Iron Atlas firewall snapshot model while preserving policy order, object relationships, routing intent, SD-WAN configuration, original IP representations, canonical CIDR networks, and reviewable evidence.

It is intentionally runnable as a standalone Go command on Windows, WSL, or Linux so development can continue independently from the active Iron Atlas service and authentication work.

## Pipeline

```text
Untrusted FortiGate-generated YAML bytes
        |
        v
Bounded reader and cancellation check
        |
        v
Bounded Fortinet invalid-export compatibility repair
        |
        v
Maintained go-yaml v4 node decoder
        |
        v
Iron Atlas YAML admission and resource validator
        |
        v
Stable Iron Atlas YAMLDocument/YAMLNode
        |
        v
Bounded FortiGate semantic layout detection
        |
        v
Existing vendor-aware FortiGate normalizer
        |
        v
Experimental vendor-independent FirewallSnapshot
        |
        +-- summary output
        +-- deterministic JSON output
        +-- typed resolved, built-in, unresolved, and ambiguous reference graph
        +-- upload-safe semantic-quality aggregates
        +-- reviewable findings
```

The native FortiOS parser remains unchanged. A future accepted boundary can place both adapters behind shared FortiGate semantics and the same normalized snapshot contract after fixture coverage and version compatibility are proven.

## Maintained decoder decision

The YAML syntax authority is pinned to:

```text
go.yaml.in/yaml/v4 v4.0.0-rc.6
```

Iron Atlas loads the representation-node tree rather than unmarshalling directly into normalized structs. The node boundary preserves source order, scalar text, comments, and source locations while allowing Atlas to reject features that the maintained decoder understands but the FortiGate adapter does not admit.

Private structural diagnostics established that the available FortiGate export contains both adjacent quoted fragments and a double-quoted fragment followed by a bare fragment for multi-value mapping attributes. Neither form is YAML. Fortinet's [FortiOS 7.4.8](https://docs.fortinet.com/document/fortigate/7.4.8/fortios-release-notes/289806/resolved-issues) and [FortiOS 7.6.3](https://docs.fortinet.com/document/fortigate/7.6.3/fortios-release-notes/289806/resolved-issues) release notes independently document invalid YAML generation for multi-value attributes and long strings. The adapter therefore performs one narrow repair before maintained decoding: an entire value containing two or more space-separated fragments becomes a flow sequence on the same physical line when the first fragment is double quoted and every remaining fragment is either double quoted or a restricted ASCII bare CLI token. Bare fragments are quoted during repair so their string semantics cannot be changed by YAML type resolution. Fragment contents are neither decoded nor logged by the repair. Ordinary valid plain scalars, single-quoted adjacency, unsafe bare fragments, and all other malformed forms remain unchanged for the maintained decoder to admit or reject.

A later content-free composer diagnostic established that the export also emits a literal object name beginning with `*` as an unquoted mapping key, while emitting the same name quoted in the object's `name` field. A following diagnostic established that such a name can end in `?`. YAML interprets an unquoted leading `*` as an alias rather than a literal key. The compatibility boundary therefore also quotes a complete nested-mapping key line when its name begins with `*`, `&`, `!`, `%`, or `@`. Remaining name bytes must be visible ASCII and cannot be whitespace, quotes, backslashes, flow delimiters, comment markers, or colons; punctuation such as `?` is retained. The leading characters otherwise introduce an alias, anchor, tag, directive, or reserved token. The rule does not apply to values, inline mappings, quoted keys, or ordinary keys; aliases and other excluded YAML features remain rejected by the maintained decoder and Atlas admission boundary.

The dependency is maintained by the YAML organization. Its use and the retained Atlas policy boundary are recorded in [ADR-0007](../decisions/ADR-0007-MAINTAINED-YAML-DECODER.md).

## Stable internal contract

The existing `YAMLDocument` and `YAMLNode` types remain the only YAML contract consumed by normalization:

```go
type YAMLNode struct {
    Kind   YAMLKind
    Value  string
    Map    map[string]*YAMLNode
    Order  []string
    Seq    []*YAMLNode
    Line   int
    Column int
}
```

Mapping insertion order remains in `Order`; sequence and policy order remain in `Seq`; scalar values remain strings even when they resemble booleans, numbers, octal values, IP addresses, masks, IDs, or version numbers. Line and column are copied from the maintained decoder node at the scalar or collection start.

An omitted mapping value retains the prototype's established empty-mapping representation so the normalizer contract does not change.

## FortiGate semantic layouts

The original synthetic fixture placed global configuration under `global` and VDOM configuration under `vdom`. A complete private decode proved that this envelope is not a valid admission requirement for the available native FortiGate backup. The adapter now recognizes three bounded layouts after syntax and node admission:

- the existing canonical `global`/`vdom` fixture envelope;
- supported FortiOS CMDB sections directly at the document root; and
- one or more VDOM containers whose immediate children include supported CMDB sections.

The detector uses the fixed set of public section labels already consumed by normalization, including `system_global`, `system_interface`, `firewall_address`, `firewall_policy`, `router_static`, and related supported sections. It does not infer configuration meaning from scalar values or arbitrary private keys. Unrelated YAML remains rejected. Native container names become VDOM scope only after their child structure is recognized.

## YAML admission policy

Allowed by the first maintained-decoder boundary:

- Block mappings
- Block and flow sequences
- Scalar mapping keys
- Plain, single-quoted, and double-quoted scalars
- Double-quoted scalars continued across physical lines
- Ordinary comments and a standard document marker
- Standard YAML scalar tags resolved by the maintained decoder while retaining raw scalar text
- Bounded Fortinet compatibility repairs for adjacent multi-value fragments beginning with a double-quoted value and restricted literal object-name mapping keys beginning with YAML indicator characters

Rejected before normalization:

- Duplicate mapping keys
- Anchors and aliases
- Cyclic alias graphs
- Custom tags
- Multiple documents
- Flow mappings
- Literal and folded block scalars
- Empty documents
- Non-scalar mapping keys
- Empty mapping keys
- Unsupported node kinds
- Any configured resource-limit violation

Flow mappings and block scalars remain rejected until sanitized FortiGate evidence demonstrates a requirement and defines the exact internal contract. The decoder's broader grammar support does not silently broaden the adapter's admission policy.

## Resource controls

Default limits are deliberately above the known 13,220,044-byte real export while remaining finite:

| Control | Default |
| --- | ---: |
| Input bytes | 64 MiB |
| Decoder nodes, including mapping keys | 4,000,000 |
| Mapping depth | 128 |
| Sequence depth | 128 |
| Combined collection depth | 192 |
| Entries in one mapping | 500,000 |
| Entries in one sequence | 500,000 |
| Total mapping entries | 2,000,000 |
| Total sequence entries | 2,000,000 |
| Scalar bytes | 8 MiB |
| Mapping-key bytes | 64 KiB |
| Fortinet compatibility rewrites | 100,000 |
| Fortinet compatibility fragments | 1,000,000 |
| Normalized records | 1,000,000 |
| Resolved, built-in, unresolved, and ambiguous references | 2,000,000 |
| Findings | 500,000 |

The maintained loader receives an early depth guard and rejects alias expansion. Before that load, the compatibility pass preserves line count, caps repaired lines and fragments including repaired keys, and re-applies the input-byte ceiling after inserting flow-sequence delimiters or key quotes. Iron Atlas then applies exact per-kind depth, node, collection, key, and scalar limits while converting the node tree. Normalized record, reference, and finding caps provide a final adapter boundary. Input and node caps bound the work that can reach normalization.

Context cancellation is checked during bounded input reads and throughout Atlas node admission. The maintained in-memory syntax decode itself is bounded by the byte, depth, unique-key, and alias controls.

## Error and secret handling

Syntax-decoder errors are reduced to an allowlisted decoder stage plus structured line and column when the decoder reports them. A compatibility fallback extracts positions from older error text without returning that text. The underlying decoder message is never returned because it may include source scalar content.

The inspection command's redacted summary mode replaces the input path, device identity, and FortiOS version, suppresses finding details, and retains only inventory counts plus bounded errors so the resulting log can be shared for diagnosis. It refuses JSON output because a normalized snapshot contains infrastructure names and addresses.

The separate `-format structure -redact` mode runs maintained decoding and Atlas node admission without requiring semantic layout acceptance. It reports root kind and entry count, canonical-wrapper presence, recognized public FortiOS section labels, nested-mapping counts, detected VDOM-container counts, and unknown-root-entry counts. It never reports scalar values, unrecognized keys, or VDOM names. This ensures a semantic admission failure produces an upload-safe, actionable diagnostic instead of a positionless error.

The dedicated `-format quality -redact` mode performs one bounded decode and normalization pass, then reports only aggregate Atlas-owned labels: normalized record counts by fixed kind; resolved references by fixed role and kind; built-in references by fixed kind; unresolved and ambiguous references by fixed role and kind; findings by allowlisted severity, category, and title; root-layout coverage counts; and fixed coverage-warning codes. It never emits vendor names, source-derived object identifiers, YAML paths, scalar values, finding details, VDOM names, input paths, or normalized JSON. Unknown future labels are replaced with fixed unclassified buckets rather than echoed.

Admission and resource errors report only:

- line;
- column;
- rejected construct class; and
- violated limit name.

They do not echo keys, scalar values, certificates, hashes, addresses, VPN material, comments, or other configuration evidence.

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
- Typed resolved edges, built-in references, explicit unresolved or ambiguous references, source locations, and findings

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

Resolved references produce graph edges that retain the fixed relationship role and resolved object kind. Recognized FortiOS built-ins produce explicit built-in-reference records instead of disappearing from the graph. Missing or ambiguous references produce explicit unresolved-reference records with a retained resolution classification and configured-state findings; they are never silently converted into empty values. All four classes count against the same bounded reference limit.

## Security and evidence handling

Use password-masked and sanitized backups for development whenever possible.

Do not commit:

- Raw production firewall backups
- Credentials or database URLs
- Pre-shared keys or private keys
- Certificates containing private material
- Unredacted addresses, names, or customer evidence
- Generated normalized output from production configurations

The retained repository fixtures use documentation ranges and synthetic names only. The known production-derived 13.2 MB export remains outside Git and is read in place only during user-controlled local validation.

A backup establishes `CONFIGURED` intent. It does not establish live route installation, current SD-WAN health or member selection, active VPN state, interface link state, HA state, or current sessions. Those require separately collected operational evidence.

## Validation

The prototype test fixtures and hostile-input tests cover:

- The complete existing normalized synthetic inventory and reference counts
- Multiline double-quoted scalar continuation and escaped quotes
- Mapping, sequence, and policy order
- Scalar text that resembles booleans, integers, octal values, IPs, masks, IDs, and versions
- Source line and column preservation
- Duplicate key, anchor, alias, custom tag, multi-document, flow-mapping, and block-scalar rejection
- Exact-limit acceptance and one-over-limit rejection for input, node, depth, mapping, sequence, key, and scalar controls
- Normalized record, reference, and finding limits
- Error-message scalar redaction
- Quoted-only and quoted-plus-bare adjacent multi-value repair, YAML-indicator object-name key repair including retained safe punctuation, alias-value and inline-value rejection, exact semantics, ordinary-scalar preservation, false-positive protection, and resource caps
- Canonical-envelope, native direct-section, and detected native VDOM-container layout normalization
- Content-free structure diagnostics that omit unknown keys, scalar values, and VDOM names
- Aggregate semantic-quality diagnostics with fixed-label allowlists, stable sorting, explicit built-in and ambiguous classifications, and a sentinel proving source-derived values cannot appear
- Stable normalized record-kind metrics, including zero-count kinds and dynamic-routing configuration records
- Context cancellation
- Explicit unresolved-reference findings

Run:

```bash
go test -count=1 ./modules/firewall/fortigate ./modules/firewall/snapshot
go test -count=1 -race ./modules/firewall/fortigate ./modules/firewall/snapshot
go vet ./cmd/fortigate-inspect ./modules/firewall/fortigate ./modules/firewall/snapshot
go run ./cmd/fortigate-inspect \
  -input modules/firewall/fortigate/testdata/fortigate-sanitized.yaml \
  -format summary \
  -redact
go run ./cmd/fortigate-inspect \
  -input modules/firewall/fortigate/testdata/fortigate-sanitized.yaml \
  -format structure \
  -redact
go run ./cmd/fortigate-inspect \
  -input modules/firewall/fortigate/testdata/fortigate-sanitized.yaml \
  -format quality \
  -redact
```

## Acceptance limitations

Before this can become an accepted production-facing ingestion boundary, Iron Atlas still requires:

1. Upload-safe aggregate classification and correction of the 1,597 unresolved references and 1,605 findings observed in the first private normalized run, without using private object names as fixtures or special cases.
2. Independent aggregate inventory comparison against the same appliance through FortiGate GUI totals, CLI table counts, REST count-only results, or a locally processed native FortiOS export.
3. Sanitized fixtures from multiple FortiOS major and minor versions.
4. Single-VDOM and multi-VDOM fixtures from actual FortiGate YAML exports.
5. Coverage for IPv6, central NAT variants, local-in policies, HA, additional VPN forms, dynamic routing details, and nested object behavior.
6. Differential comparison against the same appliance configuration exported in native FortiOS and YAML formats.
7. Representative resource telemetry on large sanitized inputs.
8. Review of the pre-release v4 dependency at every version update and migration to a stable v4 release when available and validated.
9. Integration with evidence storage, snapshot persistence, authorization, and UI boundaries.
10. Clean-clone repository validation on the exact future pushed candidate commit.
