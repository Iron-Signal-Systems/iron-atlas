# FortiGate Module

This package contains two separate experimental syntax paths:

- the initial native FortiOS hierarchy parser; and
- the FortiGate-generated YAML adapter that uses a pinned maintained YAML node decoder behind Iron Atlas admission limits.

The YAML adapter converts admitted decoder nodes into the stable internal `YAMLDocument` and `YAMLNode` contract, then uses the existing FortiGate normalizer and vendor-independent firewall snapshot model. It preserves mapping and sequence order, scalar text, and source locations while rejecting aliases, anchors, custom tags, multiple documents, flow mappings, block scalars, duplicates, and configured resource-limit violations.

After node admission, the adapter detects the semantic layout instead of requiring one synthetic root envelope. It accepts the existing `global`/`vdom` fixture layout, native direct FortiOS CMDB sections such as `system_global` and `firewall_policy`, and native VDOM containers whose children contain supported CMDB sections. Detection uses only the public section labels already consumed by the normalizer. Unrelated YAML remains rejected.

Before maintained decoding, a bounded Fortinet-specific compatibility pass repairs two vendor export-defect families. First, two or more adjacent fragments used for one multi-value mapping attribute are converted to a YAML flow sequence when the first fragment is double quoted and each remaining fragment is either double quoted or a restricted bare CLI token. Second, a literal object-name key beginning with `*`, `&`, `!`, `%`, or `@` is quoted when it occupies a complete nested-mapping key line. The remaining name may use visible ASCII other than whitespace, quotes, backslashes, flow delimiters, comment markers, or colons; punctuation such as a trailing `?` is retained. This prevents YAML from interpreting the name as an alias, anchor, tag, directive, or reserved token. The repair does not change physical line count or reinterpret ordinary valid plain scalars. Alias values and all broader invalid YAML still fail closed.

The adapters remain experimental. Native and YAML syntax should converge on shared FortiGate semantics rather than duplicate interface, route, SD-WAN, object, policy, NAT, VPN, QoS, and reference normalization.

`fortigate-inspect -format summary -redact` produces an upload-safe diagnostic summary: it omits the input path, device identity, FortiOS version, and finding details while retaining inventory counts and bounded decoder stage/position errors. `fortigate-inspect -format structure -redact` decodes without semantic admission and reports only root counts, canonical-wrapper presence, recognized public CMDB section labels, and detected-container counts. `fortigate-inspect -format quality -redact` performs one bounded decode and normalization pass and reports only aggregate fixed-label record, reference, finding, and coverage counts. The quality report requires `-redact`; unknown roles, kinds, severities, categories, or titles collapse into fixed unclassified buckets instead of being echoed. None of the redacted modes emit private scalar values, source-derived object identifiers, YAML paths, VDOM names, input paths, or normalized JSON.

The normalized reference graph now retains resolved object kinds, explicit built-in references, and explicit unresolved-versus-ambiguous classification. Resolved, built-in, unresolved, and ambiguous references all count against the same resource ceiling.

See [the FortiGate YAML architecture](../../../docs/architecture/FORTIGATE-YAML-SNAPSHOT-PROTOTYPE.md) and [ADR-0007](../../../docs/decisions/ADR-0007-MAINTAINED-YAML-DECODER.md).
