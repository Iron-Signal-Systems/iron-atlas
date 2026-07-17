# Cisco Network Modules

The Cisco module is an offline infrastructure-evidence candidate for Cisco IOS
and IOS XE switching. It is not a live collector, controller, configuration
engine, or production-ready compatibility implementation.

## Current Executable Boundary

`common.ParseEvidenceBundle` accepts the repository's delimited offline command
bundle format and returns an ordered, versioned `EvidenceBundle`:

```text
===== COMMAND: show version =====
sanitized command output
===== END COMMAND =====
```

The parser currently provides:

- ordered command evidence and duplicate-command preservation;
- bundle, command, and parser schema/version identities;
- exact-consumed-input, canonical-bundle, and normalized-output SHA-256 values;
- complete, incomplete, failed, unsupported, and truncation state types;
- structured diagnostics for malformed, partial, ignored, cancelled, and
  resource-limited evidence;
- explicit total-input, per-command, per-line, and command-count budgets;
- context cancellation before and between reads;
- partial evidence retention when bounded parsing can do so safely; and
- an upload-safe projection that excludes command names, raw command output,
  and diagnostic detail.

Default parser budgets are:

| Resource | Default |
| --- | ---: |
| Total input | 32 MiB |
| One command output | 8 MiB |
| One line | 1 MiB |
| Commands | 256 |

Callers with tighter trust or deployment boundaries should supply smaller
positive limits. Invalid or internally inconsistent limits are rejected.

`common.ParseCommandBundle` remains as a migration-only compatibility boundary.
It returns the historical `map[string]string`, so it cannot preserve ordering
and the last duplicate command replaces earlier duplicates. New Cisco work must
use `ParseEvidenceBundle`.

## Endpoint Attribution

`trunk.Classify` fails closed. Unknown or insufficient interface evidence is
classified as `unknown` and cannot receive local endpoint attribution.

Positive administrative or operational access-mode evidence is required before
an interface may be classified as an access endpoint. Trunks, routed
interfaces, port-channel members, stack links, fabric links, and explicitly
excluded interfaces are denied endpoint attribution. Trunk evidence takes
precedence over conflicting access evidence.

The current analyzer produces trunk findings only for interfaces positively
classified as infrastructure trunks. It does not convert unknown state into
trunk findings or endpoint attribution.

## Evidence Safety

Raw production switch output, configurations, credentials, addresses, device
names, site names, protected filenames, and employer-specific findings are
prohibited from Git and upload-safe diagnostics. Fixtures must be synthetic or
explicitly approved and sanitized.

The upload-safe projection retains structural metadata such as versions,
sequence numbers, byte and line counts, digests, completion/truncation state,
and diagnostic codes. It does not make the underlying evidence safe to commit.

## Validation

Run the current Cisco module tests with:

```bash
go test ./modules/network/cisco/... -count=1
```

See [Cisco offline evidence testing](../../../docs/testing/CISCO-OFFLINE-EVIDENCE-TESTING.md)
for the adversarial cases and current limitations.

## Explicit Nonclaims

This candidate does not yet provide:

- live SSH collection or credential and host-key management;
- command-specific IOS or IOS XE parsers;
- acquisition-time or collector identity provenance;
- Catalyst platform profiles or compatibility acceptance;
- a Cisco inspection CLI or deterministic JSON/Markdown report;
- persistence, topology generation, Zabbix, or Graylog integration;
- formal Phase 3 acceptance; or
- production readiness.
