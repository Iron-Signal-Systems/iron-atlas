# Validation Environment Model

Repositories define portable, canonical, specialized, release, and recovery
environment profiles as applicable.

Profiles declare:

- permitted operating systems and architectures;
- required commands and command-version patterns;
- pinned Python validation modules;
- required environment-variable names without secret values;
- capabilities and explicit non-capabilities;
- whether containers are required by an accepted product boundary.

The environment doctor emits a human-readable result and can write a
machine-readable fingerprint. Command presence without version verification is
not sufficient for canonical or formal acceptance.

Containers remain optional unless the accepted product deployment is itself
container-native. Native hosts and disposable virtual machines are first-class
validation environments.

## Pinned Go toolchain

Iron Atlas retains Go language compatibility at `1.25.0` while
declaring `go1.26.5` as the preferred build and validation toolchain in
`go.mod`.

The exact toolchain is security-relevant because `govulncheck` evaluates
reachable vulnerabilities in the standard library supplied by the selected Go
toolchain. Hosted validation sets `GOTOOLCHAIN=go1.26.5` and the Go
command downloads and checksum-verifies that official toolchain when it is not
already available on the host.

The portable environment profile and project toolchain requirements reject Go
toolchains older than `1.26.5`. A host-specific patched build may
include a suffix after `go1.26.5` but must report the same base release.
