# ADR-0001 — Go-First Standard-Library Baseline

## Status

Accepted initial direction.

## Decision

Use Go for production services, collection orchestration, parsers, analysis, HTML rendering, API delivery, and initial external adapters. Prefer the standard library until a dependency provides a measured and governed benefit.

Rust is not required initially. Python remains development and validation tooling.

## Consequences

- Small deployable binaries
- No Node.js/npm runtime for the initial UI
- Fewer package and supply-chain dependencies
- Some formats, such as general YAML, require a reviewed dependency or isolated adapter in a later phase
