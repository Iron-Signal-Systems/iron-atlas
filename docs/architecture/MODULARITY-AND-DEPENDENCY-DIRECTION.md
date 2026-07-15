# Modularity and Dependency Direction

## Principle

Every vendor, collector, parser, analyzer, diagram renderer, identity integration, and monitoring destination must be replaceable behind a versioned contract.

## Module Contract

Each module declares:

- Stable module identifier
- Supported vendor, platform, versions, and evidence formats
- Input contract
- Normalized output contract
- Parser/analyzer version
- Unsupported and partially supported sections
- Security classification
- Resource limits
- Tests and fixtures
- Compatibility and deprecation policy

## Dependency Direction

```text
Core records and governance
          ↓
Capability contracts
          ↓
Vendor modules
          ↓
External devices and products
```

Core packages must not import Fortinet-, Cisco-, Zabbix-, or other vendor-specific implementation packages. Vendor modules may depend on common contracts.

## Go-First Decision

Production services and modules should be implemented in Go where practical. Go supplies the HTTP server, HTML templating, XML parsing, cryptography, concurrency, networking, and binary deployment requirements without a large runtime dependency chain.

Rust is not required by the initial architecture. It may be introduced only when a measured requirement cannot be met safely and maintainably in Go.

Python is permitted for development validation, fixture generation, migration checks, and carefully isolated tooling. Python is not required for the core production request path.
