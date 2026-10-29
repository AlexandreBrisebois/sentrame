# Data Lifecycle

## ADDED Requirements

### Requirement: Explicit greenfield cutover
The new model SHALL require no legacy schema or API compatibility. Unsupported old local databases SHALL produce an explicit startup error and documented fresh-directory/reset guidance. Startup SHALL NOT silently delete or reinterpret existing data.

#### Scenario: Old test database selected
- **WHEN** the new binary opens an unsupported pre-change test database
- **THEN** it reports the schema mismatch without serving misinterpreted results or performing an automatic reset

### Requirement: Versioned retained structural policy
New evidence SHALL identify its projection policy; derived results SHALL identify definition and computation/timing policy. Privacy filtering SHALL precede storage, preview, live reads and copy. Projection changes SHALL NOT silently backfill or reinterpret retained history.

#### Scenario: Omitted field requested by a mapping
- **WHEN** an owner previews a condition over a field not retained by the applicable projection
- **THEN** the field is unavailable with a reason rather than read from hidden raw telemetry or invented by correlation

### Requirement: Stable native replay identity
Exact replay under supported app/producer/signal/span identity SHALL remain idempotent. Conflicting identified evidence SHALL reject atomically without overwriting accepted structure. Anonymous logs SHALL retain retry uncertainty rather than being collapsed by fingerprints.

#### Scenario: Replay complete span and its events
- **WHEN** the same span and ordered events are retransmitted unchanged after restart
- **THEN** retained evidence, graph relationships and timing samples do not multiply

#### Scenario: Changed end under same identity
- **WHEN** a span is retransmitted with a different end timestamp or conflicting retained parent/event content
- **THEN** the request rejects as an identity conflict with zero new committed records, without treating progressive span updates as supported

#### Scenario: Identical anonymous logs
- **WHEN** two legitimate point logs have identical names and timestamps but no supported stable record ID
- **THEN** both remain possible distinct records with count uncertainty rather than being silently deduplicated

### Requirement: Bounded structural lifecycle
Trace structure and relationship targets SHALL follow 30-day detail expiry and supported deletion/correction horizons. Twelve-month summaries SHALL retain only supported identity-free, text-free aggregates separated by timing/definition policy. Structural links SHALL NOT extend detail retention.

#### Scenario: Late correction and expiry
- **WHEN** a missing parent arrives within the retained horizon and the linked detail later expires
- **THEN** bounded reads first reconcile at a new revision and later remove expired navigation/detail, without extending the lifetime through replay or summaries

#### Scenario: Deletion while following a trace
- **WHEN** an owner deletes the supported correlated scope or application containing a trace
- **THEN** subsequent reads and live refreshes exclude its detail and dependent relationships under the existing replay-protection horizon, with no promise of unsupported individual deletion from finalized summaries
