# Application Integration Guidance

## Purpose

Provide reusable discovery-led instructions for coding agents helping developers instrument other applications with existing standard OpenTelemetry, separate from Sentrame's own engineering harness.

## ADDED Requirements

### Requirement: Discovery-led developer event plan
The guidance SHALL inspect relevant app source/UI/tests and current SDK/resource/signal/exporter setup before recommending changes. It SHALL ask focused developer questions about success, friction and useful observations, then present prioritized event proposals with evidence, question/value, emitter boundary, authority, privacy/noise and synthetic acceptance expectations for approval before edits.

#### Scenario: Developer asks for Sentrame instrumentation
- **WHEN** the guide is used on an app with existing OTel and an uncertain business journey
- **THEN** the agent first discovers current behavior, asks only unresolved useful questions and proposes a minimal reviewed plan rather than replacing instrumentation or emitting everything

### Requirement: Standard carriers and scoped correlation
The guidance SHALL separate supported standard OTLP intake from sufficient optional journey evidence. It SHALL cover client-only, server-only and combined journeys, stable event/step meaning, bounded app-scoped correlation, timing and goals without a mandatory Sentrame SDK/protocol. Correlation SHALL NOT imply source authority, ordering, deduplication or person/cross-visit identity.

#### Scenario: Server-only or combined integration
- **WHEN** a developer instruments a server workflow without a browser, or correlates browser/server observations
- **THEN** the plan needs no fictitious tab for the server-only path and uses the same correctly scoped logical correlation for the combined path while preserving sources, clocks and distinct identities

### Requirement: Proposed warning error and recovery instrumentation
The guidance SHALL inspect existing and missing failure/retry/fallback/validation boundaries and propose useful warning/error and recovery events with severity rationale, safe context, reporting ownership and noise/deduplication controls. It SHALL NOT recommend synthetic throws, swallowed errors, changed app behavior or automatic error severity for expected user behavior merely to collect telemetry.

#### Scenario: Missing degraded-fallback warning
- **WHEN** source inspection finds a successful but degraded fallback without useful telemetry
- **THEN** the agent proposes an appropriately classified event and recovery observation, explains its value/privacy/test plan and waits for the developer's event-plan approval before edits

### Requirement: Honest examples and authorized validation
The guide SHALL label illustrative names/pseudocode and actual verified SDK examples distinctly. It SHALL require developer-authorized minimal changes and synthetic validation before compatibility claims, preserve privacy and avoid production experiments/secrets. It SHALL distinguish its external-app integration purpose from Sentrame's internal harness and require no AI runtime service.

#### Scenario: Planning-only use
- **WHEN** only discovery/event planning is authorized
- **THEN** the agent returns evidence, questions and the reviewed plan without installing packages, editing the app, deploying or claiming a tested integration
