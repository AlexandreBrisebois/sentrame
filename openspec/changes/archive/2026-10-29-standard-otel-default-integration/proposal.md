# Proposal

## Why

Sentrame accepts standard OTLP, but currently discards span structure and duration evidence and requires logical journey correlation plus paired stage conventions for its useful timing workflow. Developers should first discover and explore ordinary OpenTelemetry, then configure business meaning in Sentrame and add instrumentation only for unanswered questions.

## What Changes

- Introduce five explicit capability levels: receipt/discovery; trace-scoped exploration/timing; configured business interpretation; explicitly correlated multi-request journeys; optional browser presence. Missing higher-level context never invalidates supported OTLP receipt.
- Retain bounded span end timestamps, parent relationships, span kind/status, span events, correlated logs and relevant links. Use native span identity and duration for operation timing; do not demand duplicate start/end logs or stage occurrence IDs for an existing span.
- Map retained span/event names and safe supported attributes to business stages, supporting activities, goals and exceptions without `stage.name`. Keep trace structure separate from chronological journey transitions and their denominators.
- Preserve optional `stage.name` as a Sentrame application convention and configuration suggestion. Explicit owner rules win; hints never activate mappings, supply authority, establish boundaries or rewrite historical interpretation.
- Keep optional paired stage timing, prospective definitions, captured boundaries, saved references, privacy projection and producer trust. Use a greenfield evidence model without legacy schema, replay or API compatibility adapters.
- Deliver coherent .NET and Next.js documentation using one progressive voting example, bounded agent instructions and separate wire/exporter/deployment qualification records.

Scope is local specification and documentation planning. No runtime edits, external-app instrumentation, telemetry export, deployment, publication, commit or archival occurs in this change's preparation. Implementation tasks require a later apply request. Roadmap B stays deferred: no multi-app workspace, richer release comparison, narrow visit grants, generic query engine, new gRPC listener, metric arithmetic, attention inference, session replay, identity tracking or AI service.

## Capabilities

### New Capabilities

- `trace-exploration`: bounded operation structure, native timing and correlated evidence without visit fields.

### Modified Capabilities

- `telemetry-intake`: incremental standard structural projection and discovery, separate from assignment.
- `journey-modeling`: configurable standard-name interpretation, optional hints, native span occurrences and safe correlation.
- `journey-exploration`: explicit trace-versus-path presentation and preserved historical references.
- `data-lifecycle`: explicit projection versions, bounded corrections and privacy-safe structural retention.
- `app-integration-guidance`: progressive .NET/Next.js recipes, document ownership and qualification discipline.
- `journey-diagnostics`: replace the obsolete legacy-projection replay scenario with same-model diagnostic replay; retain privacy, discovery and expiry behavior.

The baseline was concurrently archived and synchronized during this planning pass. The six existing paths above now target canonical `openspec/specs/` owners; unique ADDED requirements extend five of them, and one MODIFIED diagnostic requirement removes the legacy-data obligation. [Reconciliation](compatibility.md) records the observed transition. No canonical or archive files were changed by this work.

## Impact

Future implementation touches `internal/app/intake.go`, `types.go`, `evidence_json.go`, `store.go`, `journey.go`, query/reference/lifecycle/summary handlers, protected read authorization, `web/src`, fixture generation and REST requests. Reuse the maintained Collector receiver and standard exporters. No new producer SDK or runtime dependency is selected by this plan.

Required supporting deliverables: [design](design.md), [documentation change matrix and example plan](documentation-plan.md), [greenfield cutover and overlap plan](compatibility.md), [qualification matrix](qualification.md), capability deltas under `specs/`, and ordered test-first [tasks](tasks.md). Current integration-guide edits remain untouched. The owner's clarification establishes there are no users and no backward-compatibility requirement; the cutover plan replaces the original request's legacy integration migration obligations.
