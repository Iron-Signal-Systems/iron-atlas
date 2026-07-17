# Cisco Offline Evidence Testing

> Status: bounded offline evidence-foundation candidate; not formally accepted;
> no live Cisco collection or production compatibility claim

## Purpose

This document defines the tests for the first shared Cisco IOS and IOS XE
evidence boundary. The boundary must safely retain ordered, sanitized offline
command evidence before Atlas adds command-specific parsers, platform profiles,
reporting, persistence, or restricted live collection.

## Test Command

```bash
go test ./modules/network/cisco/... -count=1
```

This targeted command is necessary for candidate development but is not a
formal phase-acceptance gate. The complete repository test framework remains
required before the bounded change is committed and proposed for review.

## Ordered Evidence Contract

Tests prove that the parser:

- retains command execution order;
- retains repeated executions of the same command as separate records;
- assigns stable one-based command sequence numbers;
- identifies bundle, command, and parser versions;
- normalizes command output consistently;
- produces deterministic SHA-256 values for identical input and parser
  versions; and
- retains partial command evidence with an explicit non-complete state when
  structurally possible.

The compatibility `ParseCommandBundle` test protects existing callers during
migration. It does not grant the legacy map representation any ordering or
duplicate-preservation claim.

## Structural and Resource Adversarial Cases

The candidate test matrix includes:

| Case | Required result |
| --- | --- |
| Ordered multi-command bundle | Preserve sequence and normalized output |
| Duplicate command | Preserve both executions |
| Nested command section | Return malformed status, diagnostic, and safe partial evidence |
| End marker without command | Return malformed status and diagnostic |
| Unclosed command | Retain partial command as incomplete and diagnose it |
| Empty command name | Retain failed command metadata and diagnose it |
| Nonempty text outside sections | Ignore raw text and emit a non-leaking diagnostic |
| Total-input limit | Stop at the hard budget and mark truncation |
| Per-command limit | Retain bounded output and mark the command truncated/incomplete |
| Per-line limit | Stop scanning and report bounded truncation |
| Command-count limit | Retain only the allowed ordered commands and report truncation |
| Cancelled context | Stop with the context error and cancellation diagnostic |

Default limits are 32 MiB total input, 8 MiB per command, 1 MiB per line,
and 256 commands. Tests use deliberately smaller limits to exercise each
boundary without large fixtures.

Cancellation is observed before parsing and before or after reader calls. A
caller that supplies potentially blocking I/O must still own transport-level
deadlines or another mechanism capable of interrupting that I/O.

## Upload-Safe Diagnostics

The upload-safe projection is serialized during testing with protected marker
values. The serialized result must not contain raw text from outside command
sections, command names, or command output.

Permitted structural fields include schema and parser versions, command
sequence, byte and line counts, digests, completion and truncation state, and
diagnostic code, severity, stage, line, and command sequence.

An upload-safe diagnostic summary is not permission to upload or commit the raw
evidence bundle.

## Fail-Closed Endpoint Attribution

Tests prove:

- unknown or insufficient interface state denies endpoint attribution;
- positive configured or operational access-mode evidence permits attribution;
- trunk evidence overrides conflicting access evidence;
- trunks, routed interfaces, port-channel members, stack links, fabric links,
  and explicitly excluded interfaces deny attribution; and
- an unknown interface does not produce trunk-specific findings.

MAC-address presence alone is never positive local-endpoint evidence.

## Current Limitations

These tests do not claim coverage for real Catalyst output variants, IOS or IOS
XE command parsing, platform detection, acquisition provenance, live collection,
topology analysis, persistence, reporting, external-system integrations, or
production operation. Those boundaries require later separately reviewed and
accepted changes.
