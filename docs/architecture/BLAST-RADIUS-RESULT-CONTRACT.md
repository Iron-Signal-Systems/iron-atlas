# Blast-Radius Result Contract

## Status

Planning contract.

## 1. Purpose

Define the deterministic, evidence-backed result produced for a compromise,
failure, or governed change-impact calculation.

## 2. Input identity

Every result binds:

- result ID;
- incident, change, or analysis ID;
- subject identities;
- requested time window;
- effective evidence window;
- accepted-state cutoff;
- organization and scope;
- analyzer release;
- policy release;
- canonical model release;
- evidence-set digest; and
- request and authorization context where retained.

## 3. Result categories

The result contains separate sections for:

- observed impact;
- potential impact;
- prevented impact;
- unknown impact;
- network radius;
- identity and privilege radius;
- data radius;
- infrastructure-control radius;
- operational radius;
- governance and notification radius;
- containment effects;
- evidence conflicts; and
- additional evidence required.

No section may be inferred from another merely to complete the report.

## 4. Relationship result

Each relationship result includes:

- source entity;
- destination entity;
- capability or dependency;
- direction;
- valid time;
- impact category;
- evidence state;
- protocol, port, right, privilege, permission, or trust requirement;
- supporting evidence IDs;
- confidence and correlation quality;
- assumptions;
- conflicting evidence;
- coverage limitations;
- calculated path or rule identity; and
- explanation.

## 5. Coverage summary

The coverage summary reports:

- expected sources;
- available sources;
- last accepted source time;
- stale sources;
- missing sources;
- collection gaps;
- unsupported record families;
- unresolved identities;
- conflicting identities;
- clock uncertainty;
- retained evidence limits; and
- effect on the answer.

## 6. Confidence

Confidence is explanatory, not a substitute for evidence state.

A confidence value records:

- method;
- inputs;
- thresholds;
- release;
- reasons increasing confidence;
- reasons decreasing confidence; and
- human disposition.

## 7. Risk presentation

Atlas may provide ordered severity or urgency when governed criteria exist.

Atlas shall not provide one unexplained risk score as the complete answer.

## 8. Determinism and history

The same:

- accepted evidence set;
- canonical model;
- analyzer release;
- policy release;
- subject;
- scope; and
- time window

must produce the same semantic result digest.

A later evidence set or analyzer creates a new result linked to the previous
result. It never silently rewrites the earlier result.

## 9. Canonical output

The canonical output is strict versioned JSON.

Human-readable HTML, Markdown, PDF, diagram, or API views are derived
representations and identify the canonical result digest.

## 10. Required human-readable answer

The first screen or report page presents:

- what happened;
- what could have happened;
- what was prevented;
- what remains unknown;
- affected systems and data;
- operational consequences;
- containment status;
- required owners and reviewers;
- evidence quality; and
- immediate unanswered questions.

## 11. Nonclaims

A valid result does not prove:

- complete source coverage;
- malicious intent;
- successful data exfiltration;
- legal notification applicability;
- absence of persistence;
- absence of credential theft;
- absence of an unknown path; or
- production readiness of an unaccepted adapter.
