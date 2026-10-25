# Tasks

This checklist describes future work, not authorization to execute it. All items remain unchecked. Review and revise the proposed defaults before a separate implementation request; stack selection must preserve the documented behavior and architecture.

## 1. Local foundation and contract

- [ ] 1.1 Record the chosen language, transactional store and minimal interface components in design.md, including compatible OTel exporter/protocol libraries; verify the selections preserve the local single-service contract without adding a provider or broker.
- [ ] 1.2 Create the local service and synthetic producer skeleton, with loopback-only configuration and isolated synthetic storage; verify startup does not contact WFS or an external telemetry destination.
- [ ] 1.3 Document the concrete OTLP field mapping, reason-code enumeration, schema version and size/time limits; add valid and invalid fixtures and verify they exercise every required/forbidden field, including resource and scope attributes.
- [ ] 1.4 Add a minimal Taskfile exposing spec validation, focused tests and the synthetic demo using the selected stack; verify each documented command resolves and preserve the existing OpenSpec setup.

## 2. Evidence ingestion and lifecycle

- [ ] 2.1 Implement OTLP validation, per-record rejection feedback and commit-before-ack storage; verify valid/invalid mixed batches, malformed/oversized requests, unsupported versions, disallowed scopes and absence of raw payloads in errors/logs.
- [ ] 2.2 Implement event deduplication and bounded producer retries with stable IDs; verify identical replay, conflicting duplicates, lost acknowledgements, concurrent duplicate submissions and restart persistence without count inflation.
- [ ] 2.3 Implement attempt/application isolation and occurrence/receipt timestamps; verify out-of-order evidence, missing starts, stale input and future clock bounds, and document the profile's compatibility limits.
- [ ] 2.4 Implement seven-day expiry, daily physical cleanup and local attempt/application deletion with bounded deletion markers; verify expired evidence is immediately unqueryable, derived results disappear, deletion survives restart and retries cannot recreate deleted evidence within the acceptance window.

## 3. Outcome rules and insights

- [ ] 3.1 Implement rule version 1 from journey-outcomes; verify every classification scenario, including optimistic requests, terminal conflicts, persistence after an error, delayed success and identical results under delivery-order permutations.
- [ ] 3.2 Implement summary queries with started-attempt denominators and separate incomplete attempts; verify the four-start/one-orphan fixture, duplicate replay, expiry, deletion and unavailable timing for inconsistent clocks.
- [ ] 3.3 Build the single summary and expandable evidence view; verify it displays synthetic status, rule version, classification reasons, revisability, missing evidence and unknown collection completeness, without intent or abandonment claims. Document how to inspect one attempt.

## 4. End-to-end acceptance

- [ ] 4.1 Run the generated fixture suite through actual OTLP ingestion, durable storage, derivation and the interface; verify completed/interrupted/failed/unknown, delayed, duplicate and incomplete outcomes against documented expected counts and evidence.
- [ ] 4.2 Replay the suite and restart the service, then exercise deletion and expiry with a controlled clock; verify stable counts before deletion and no stale summaries afterward through the same public read path.
- [ ] 4.3 Run strict OpenSpec validation and the Taskfile checks; record passed, failed, blocked and not-run evidence, explicitly distinguishing local synthetic qualification from live WFS/browser/device qualification.

## Workflow follow-up

- Iterate on this change as product ideas are reviewed; revise proposal, specs, design and tasks together when a decision changes.
- Live integration, deployment, credentials, publication and real telemetry require separate scope and authorization; they are not completion requirements for this synthetic slice.
- Archive only after approved implementation and verification, when requested.
