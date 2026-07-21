# Validation Failure Reporting

## Status

Active project validation tooling.

## Purpose

Iron Atlas validation must help an engineer identify and correct the first
meaningful defect. A validator is incomplete when it propagates a nonzero exit
status but leaves the operator to search thousands of repeated `FAIL` lines.

## Terminal result contract

The last line emitted by every migrated top-level runner is the actionable result. No scope note, resource note, generic phase message, or separator may follow it.

A passing command ends with:

```text
FINAL RESULT: PASS — <validation title>
```

A failing command ends with:

```text
FINAL RESULT: FAIL — <primary failing check> — <extracted cause>
```

This makes terminal output, retained logs, and CI annotations identify the same result.

## Required final report

Every migrated top-level validation runner shall finish with:

- a single `FINAL RESULT`;
- pass, fail, and skipped counts;
- the **Primary failure** with its check name, command, exit status, extracted
  cause, and retained per-check log;
- **Additional unique failures** when independent causes exist;
- **Cascaded failures** grouped separately from their root cause;
- **Skipped dependent checks** and the check that blocked them;
- the final report path and per-check log directory.

Dependent checks shall not run after a required predecessor fails merely to
reproduce the same error through additional validation layers.

## Nested terminal result precedence

When a validation check executes another migrated validation runner, the caller
shall prefer the nested runner's final `FINAL RESULT: FAIL` cause over earlier
`ERROR:` or `FAIL:` lines in the nested log.

This prevents intentional failure fixtures inside successful subordinate
regression tests from being misidentified as the actual failure of the nested
runner.

## Historical checkpoints

Frozen historical validators run in an isolated clone at their exact
implementation commit. A later working tree must not rewrite an accepted
historical assertion or execute that old validator directly against successor
documentation and artifacts.

## Tool provisioning

Validation must not depend on an undeclared executable that happens to exist on
one workstation. Go-based validators are repository-managed Go tools recorded
in `go.mod`, `go.sum`, and `validation/toolchain-requirements.json`, then invoked
through `go tool`.

The authenticated-session boundary pins and invokes:

```text
go tool govulncheck
```

A new development system receives the validator through the declared Go module
graph rather than a remembered manual installation step.

## Implementation

The reusable implementation is
`tools/validation/lib/reporting.sh`. The complete repository runner and the
active authenticated-session regression and phase gate use it. Older frozen
historical gates remain unchanged and are revalidated through the existing
isolated-gate helper.
