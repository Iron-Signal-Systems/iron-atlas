# ADR-0007: Maintained YAML Decoder for FortiGate Exports

## Status

Accepted for the experimental FortiGate YAML snapshot candidate. This does not accept the adapter as a production ingestion boundary or a completed project phase.

## Context

The initial FortiGate YAML prototype implemented YAML grammar with a handwritten, physical-line parser. That boundary passed the synthetic fixture but rejected a real 13,220,044-byte FortiGate-generated export at a double-quoted scalar that continued onto the next physical line.

The exact production value remains private, and the repository does not assume that one sanitized structural sample proves every Fortinet YAML behavior. It does prove that a physical-line quote state is not a sufficient grammar boundary for the available export.

Later content-free structural diagnostics found a second boundary family: mapping values made of multiple adjacent fragments, first as double-quoted fragments and then as a double-quoted fragment followed by a bare fragment. These forms are invalid YAML. Fortinet's [FortiOS 7.4.8](https://docs.fortinet.com/document/fortigate/7.4.8/fortios-release-notes/289806/resolved-issues) and [FortiOS 7.6.3](https://docs.fortinet.com/document/fortigate/7.6.3/fortios-release-notes/289806/resolved-issues) resolved-issue notes independently describe invalid YAML generation for multi-value attributes and long strings.

Subsequent content-free composer diagnostics found unquoted literal object-name mapping keys beginning with `*`, including a name ending in `?`, while the same export quoted the first observed name in the object's `name` field. YAML treats the unquoted keys as aliases. The available export therefore requires a second narrow compatibility family for literal names beginning with YAML indicator characters without assuming an alphanumeric-only name alphabet.

After those compatibility repairs, the entire known 13,220,044-byte export passed maintained decoding and Atlas node admission. It then reached a semantic check copied from the synthetic fixture and failed because its root did not contain `global` or `vdom`. The decoder had succeeded; the remaining failure was an unsupported prototype envelope assumption. A position-based syntax diagnostic could not describe that semantic layout failure.

Maintaining YAML tokenization, quoting, escape handling, document composition, source locations, and grammar compatibility is not Iron Atlas's differentiating value. FortiGate semantics, normalized firewall modeling, evidence handling, reference resolution, findings, and operational analysis are.

## Decision

Use the YAML organization's maintained decoder:

```text
go.yaml.in/yaml/v4 v4.0.0-rc.6
```

The version and checksum are pinned in `go.mod` and `go.sum`.

Decode into the maintained library's representation-node tree, validate that tree under Atlas policy, and convert it into the existing `YAMLDocument` and `YAMLNode` types. Do not unmarshal directly into FortiGate or normalized snapshot structs.

Before maintained decoding, admit one bounded Fortinet compatibility repair for the exact invalid multi-value family: a complete mapping value containing two or more space-separated fragments is converted to a same-line flow sequence when the first fragment is double quoted and every later fragment is either double quoted or a restricted ASCII bare CLI token. Quote bare fragments during repair to preserve string semantics. Preserve physical line count; cap repaired lines and fragments; re-apply the input-byte ceiling; do not decode or log fragment contents; leave ordinary valid plain scalars unchanged; and reject all broader malformed forms.

At the same boundary, quote a complete nested-mapping key when its literal object name begins with `*`, `&`, `!`, `%`, or `@`. Require all remaining name bytes to be visible ASCII excluding whitespace, quotes, backslashes, flow delimiters, comment markers, and colons; retain other punctuation including `?`. Do not apply this rule to values, inline mappings, already quoted keys, or ordinary keys. Count each repaired key against the same rewrite and fragment ceilings. Continue to reject alias values, anchors, custom tags, and all broader excluded YAML features.

After node admission, detect either the existing canonical `global`/`vdom` fixture layout, supported FortiOS CMDB sections directly at the root, or native VDOM containers whose immediate children include supported CMDB sections. Use only the fixed public section labels already consumed by the normalizer. Do not inspect scalar content or treat arbitrary private keys as FortiGate evidence. Continue to reject unrelated YAML.

Provide an upload-safe structure format that can run even when semantic admission fails. It may report recognized public CMDB labels and aggregate node/container counts, but must omit scalar values, unknown keys, and VDOM names.

Iron Atlas retains responsibility for:

- bounded input reading;
- bounded adjacent multi-value compatibility repair beginning with a double-quoted fragment;
- bounded restricted YAML-indicator object-name mapping-key repair;
- context checks at the adapter boundary;
- exactly one document;
- unique mapping keys;
- early combined-depth and alias-expansion guards;
- node, mapping, sequence, depth, key, scalar, and aggregate limits;
- anchor, alias, custom-tag, flow-mapping, and block-scalar rejection;
- mapping and sequence order;
- source line and column;
- raw scalar text rather than unintended primitive coercion;
- normalized record, reference, and finding limits; and
- non-secret-bearing syntax, admission, and limit errors; and
- content-free semantic layout diagnostics.

The FortiGate normalization semantics and `snapshot.FirewallSnapshot` contract remain unchanged. A small layout resolver selects the global scope and VDOM nodes supplied to that normalizer.

## Rejected alternatives

### Extend the handwritten parser for multiline quotes

Rejected because it would fix one observed construct while retaining Atlas-owned responsibility for the rest of YAML grammar and its security edge cases.

### Direct typed YAML unmarshalling

Rejected because typed construction can coerce values, obscure source ordering and locations, combine syntax with vendor semantics, and make unsupported features harder to govern explicitly.

### Admit every feature understood by the dependency

Rejected because parser capability is not adapter policy. Anchors, aliases, custom tags, flow mappings, and block scalars remain outside the first admission contract unless representative sanitized evidence establishes a requirement.

### Implement a permissive Fortinet YAML parser

Rejected because the observed export defect does not justify replacing the maintained decoder or accepting arbitrary non-YAML. The compatibility pass recognizes one complete-line form and otherwise leaves syntax decisions to the maintained decoder.

### Use the frozen v3 line

Rejected for new implementation because the YAML organization directs new development and routine fixes to v4 while v3 receives security fixes only.

## Consequences

Iron Atlas gains one direct, pinned Go dependency and removes a custom grammar implementation. The dependency is a release candidate, so every update requires explicit review, complete regression, hostile-input validation, and checksum verification. Migration to a stable v4 release is required when one is available and passes the same boundary tests.

The Atlas conversion layer remains intentionally small and replaceable. A future decoder change must preserve the internal document contract and must not require a rewrite of FortiGate semantics or the vendor-independent snapshot model.

The real production-derived export remains outside Git, chat, fixtures, retained evidence, and packaged deliverables. A private local run is still required before claiming normalized real-export compatibility; maintained decoding and Atlas node admission are already proven for that file.
