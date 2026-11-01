# Data lifecycle delta

## MODIFIED Requirements

### Requirement: Stable native replay identity
Exact replay under supported app/generation/signal/span identity and explicit legacy identity aliases SHALL remain idempotent. Conflicting identified evidence SHALL reject atomically without overwriting accepted structure. Anonymous logs SHALL retain retry uncertainty rather than being collapsed by fingerprints.

#### Scenario: Replay complete span and its events
- **WHEN** the same span and ordered events are retransmitted unchanged after restart
- **THEN** retained evidence, graph relationships and timing samples do not multiply

#### Scenario: Changed end under same identity
- **WHEN** a span is retransmitted with a different end timestamp or conflicting retained parent/event content
- **THEN** the request rejects as an identity conflict with zero new committed records, without treating progressive span updates as supported

#### Scenario: Identical anonymous logs
- **WHEN** two legitimate point logs have identical names and timestamps but no supported stable record ID
- **THEN** both remain possible distinct records with count uncertainty rather than being silently deduplicated

### Requirement: Explicit greenfield cutover
Fresh installations SHALL initialize without importing unsupported legacy single-app state. Existing managed two-key installations SHALL use an explicit reviewed transition to one app key; automatic reuse without that transition SHALL fail closed. Unsupported state SHALL report a recovery error without resetting data or regenerating credentials.

#### Scenario: Old test database selected
- **WHEN** startup finds an unsupported single-app database
- **THEN** it requests a fresh directory or supported recovery path without modifying the old database

#### Scenario: Managed installation awaiting cutover
- **WHEN** the new release finds an existing two-key registry without a completed migration marker
- **THEN** it requests the reviewed offline transition and leaves state unchanged rather than starting intake with guessed credentials

## ADDED Requirements

### Requirement: Explicit one-key authority transition
A two-key transition SHALL require owner-approved policy before mutation. It SHALL retain app identity/generation, stage and validate coordinated state, publish one authority atomically and reject retired credentials after success. The proposed policy SHALL retain the existing browser key as the app key without minting or exposing a replacement. Failure or uncertain publication SHALL fail closed.

#### Scenario: Approved retained-key cutover
- **WHEN** the owner explicitly approves and completes the reviewed transition
- **THEN** existing client exporter settings keep working, the old server key rejects, one key is presented and app identity, deletion barriers and viewer grants retain their valid bindings

#### Scenario: Interrupted transition
- **WHEN** validation, storage capacity or durable publication fails or is uncertain
- **THEN** intake remains stopped until a verified old or new coordinated state is selected, with no partial conversion or automatic key regeneration

### Requirement: Non-destructive legacy evidence boundary
Cutover SHALL NOT delete or relabel legacy server records as client observations. Existing expiry/deletion rules SHALL continue. Mixed journeys and inseparable summaries SHALL be non-comparable to current client results. Client-only history SHALL remain usable only where its original interpretation and identity can be preserved; missing or unsupported history SHALL NOT become zero traffic.

#### Scenario: Mixed path A server B client C
- **WHEN** removing a legacy record would connect client A directly to client C
- **THEN** current graphs, counts, timeline, traces and diagnostics do not reconstruct that mixed journey as an A-to-C client transition
- **AND** an explicit compatibility explanation replaces unsupported old results without modifying retained evidence

#### Scenario: Summary cannot be separated
- **WHEN** an old aggregate lacks sufficient dimensions to isolate client contributions
- **THEN** it remains preserved under existing retention and is excluded as non-comparable rather than subtracted approximately, merged with current counts or displayed as zero

#### Scenario: Current live baseline crosses cutover
- **WHEN** the selected historical or baseline interval includes incompatible legacy evidence
- **THEN** affected populations identify the excluded interval and compatibility limitation independently from query truncation and exporter completeness

### Requirement: Prospective role-free definitions and references
Cutover SHALL preserve old definitions and references. It SHALL preview removal of obsolete selectors and block ambiguous or persistence-dependent rules until owner review. New role-free definitions SHALL activate prospectively; old mixed captures SHALL NOT be resegmented. Unsupported references SHALL explain recovery without silently broadening filters or overwriting the original.

#### Scenario: Conflicting rules after source removal
- **WHEN** two previously role-disjoint mappings would overlap or a goal required trusted server evidence
- **THEN** the owner sees the conflict and affected goal before activation; no general client rule or replacement goal activates automatically

#### Scenario: Pure client history and replay
- **WHEN** a historical client capture can be read with its unchanged meaning and an identified record is replayed using the retained key
- **THEN** the capture and supported contributions remain single, while aliases never join an equal-ID excluded legacy record

#### Scenario: Unsupported reference reopened
- **WHEN** a saved view depends on obsolete mixed semantics
- **THEN** the original remains stored and a localized explanation offers return or explicit current-client-period selection without creating an empty-success result

### Requirement: Revocation-safe migration recovery
Recovery SHALL verify registry, stores, generations, migration version and latest revocation/deletion barriers as coordinated state. Old code or backups SHALL NOT revive retired keys or deleted apps. Missing current authority SHALL keep intake disabled. New state SHALL be protected from writers that cannot understand its version.

#### Scenario: Rollback after retiring server key
- **WHEN** an operator selects an older image or backup after cutover
- **THEN** intake stays stopped unless latest deny state and compatible coordinated state are verified, rather than enabling the retired key from the old snapshot
