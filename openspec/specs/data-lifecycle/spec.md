# data-lifecycle Specification

## Purpose

Keep detailed journey evidence and its derivatives private, bounded and consistently correct through late arrivals, deletion, expiry and Home Server resource pressure. This is a proposed contract.

## Requirements

### Requirement: Minimized retained telemetry
The system SHALL apply a visible configurable retention allowlist and scrubbing policy before persistence, diagnostics and live publication. It SHALL exclude personal application content, household/person identifiers, credentials, sensitive URL components and fingerprinting. Random visit IDs SHALL be treated as correlatable pseudonymous data, without cross-visit identity joins.

#### Scenario: Sensitive exception message
- **WHEN** an accepted record includes a message, stack or URL carrying a credential or personal value
- **THEN** only the permitted scrubbed projection is retained and neither rejection logs nor live/copy paths expose the original

#### Scenario: Retention policy changed
- **WHEN** an operator tightens the allowed fields
- **THEN** the interface states whether existing retained detail has been scrubbed or deleted; it does not imply a prospective filter erased history

### Requirement: Thirty-day ordered detail
The system SHALL expire detailed journey evidence after 30 days under a documented event-time policy, with bounded fallback for missing/invalid timestamps. Expired detail SHALL be excluded immediately from reads and purged within a documented cleanup bound. Known retries SHALL NOT extend lifetime. Unidentifiable retries with missing source time SHALL expose age/deduplication uncertainty.

#### Scenario: Delayed record
- **WHEN** an event arrives 29 days after its valid occurrence time
- **THEN** it has at most one remaining day of detail eligibility rather than another 30 days from arrival

#### Scenario: Expiry while viewing
- **WHEN** detail expires while a view or replay cursor is active
- **THEN** subsequent reads omit it and affected cached detail is invalidated; previously disclosed client copies are outside server erasure guarantees

### Requirement: Consistent derived corrections
The system SHALL correct affected paths, outcomes, counters and caches within their supported evidence horizon after late data, replay, deletion or expiry, identifying definition/computation versions. It SHALL retain aggregate history for 12 months while raw detail expires at 30 days. Longer summary retention SHALL NOT permit individual journey drilldown or unsupported reprocessing after raw evidence is gone.

#### Scenario: Late success changes a stopping rate
- **WHEN** valid goal evidence revises a provisionally stopped journey
- **THEN** the goal denominator, stopping count, path summary and relevant live snapshot reconcile without duplicate contributions

#### Scenario: Detail expired but aggregate retained
- **WHEN** a query addresses history older than 30 days but within 12 months
- **THEN** supported counts, conversion and timing summaries remain available with their coverage/version basis, while individual ordered journey/exception drilldown is disabled

### Requirement: Mergeable private summary statistics
Twelve-month summaries SHALL retain sufficient numerator/denominator, coverage and duration-distribution structures for advertised queries at declared grains/dimensions. They SHALL exclude individual journey/correlation payloads and preserve semantic/timing versions. Combined medians SHALL NOT be computed by averaging bucket medians; approximations and unsupported queries SHALL be explicit.

#### Scenario: Merge duration buckets
- **WHEN** a requested period combines several stored timing buckets
- **THEN** median/slower statistics are derived from a qualified mergeable distribution or marked unavailable/approximate, never an average of stored medians; censoring and timing kinds remain separate

#### Scenario: Summary deletion after detail expiry
- **WHEN** an individual contribution can no longer be identified after the detail/correction horizon
- **THEN** the interface does not promise selective subtraction it cannot perform; it states the limitation and offers an explicitly selected app/time-bucket deletion policy without retaining hidden long-lived individual membership

#### Scenario: Summary retention boundary
- **WHEN** aggregate history reaches 12 months under the declared calendar/time policy
- **THEN** it expires from query/cache results and is purged within the documented cleanup bound; saved views cannot extend that retention

### Requirement: Deletion and replay protection
The system SHALL support deletion by app or identifiable visit, removing dependent detail and contributions. Journey-only deletion SHALL require stable replay membership; inferred boundaries alone SHALL NOT qualify. Replay protection SHALL state its identity/time horizon. App deletion SHALL disable intake, and re-registration SHALL prevent stale queued evidence from silently returning.

#### Scenario: Delete then retry
- **WHEN** a deleted visit's events are retried within the accepted lateness window
- **THEN** its evidence, exception counts and path contributions do not reappear

#### Scenario: Restore an older backup
- **WHEN** an approved backup is restored
- **THEN** the restore process applies expiry and deletion state before making data queryable and cannot silently reintroduce erased scope

#### Scenario: Inferred journey boundary changes after deletion
- **WHEN** a requested journey-only deletion lacks stable replay membership and late evidence could change its boundary
- **THEN** the UI explains that limitation and offers explicit broader visit/app scope without silently deleting the broader scope or promising journey-only replay protection

#### Scenario: Replay at the protection horizon
- **WHEN** identifiable or unidentifiable timestamp-free records arrive near or after a tombstone horizon
- **THEN** the declared policy states which retries can be recognized and no original-age/no-resurrection guarantee is claimed beyond supported identity and time bounds

#### Scenario: Application re-registration
- **WHEN** an application is explicitly re-registered after deletion
- **THEN** an intake generation or equivalent replay barrier prevents old queued data/credentials from becoming new evidence without an explicit reviewed policy

### Requirement: Honest ordering and timing
The system SHALL preserve occurrence and receipt timestamps, reliable source sequence/tab context where present, and a deterministic display tie-break. It SHALL expose clock-skew, missing-context and causal-order uncertainty. Invalid clocks SHALL NOT produce fabricated negative durations or a claim that arrival order is user-action order.

#### Scenario: Equal or skewed timestamps
- **WHEN** events tie or clocks disagree without a reliable causal sequence
- **THEN** display remains stable but order and duration uncertainty are visible

### Requirement: Bounded Home Server work
The system SHALL document limits for intake, storage, query ranges, graph response size, streams and background cleanup. It SHALL expose resource and collection health without retaining unsafe raw payloads. Before deployment qualification, measured idle/ingestion/query/build/test costs SHALL be compared with an owner-supplied Home Server budget and workload.

#### Scenario: Query exceeds bounds
- **WHEN** a query or graph would exceed its configured work/size bound
- **THEN** the system returns a bounded result with an explicit limit, paging or narrowing path rather than silently reporting complete counts

#### Scenario: Resource pressure
- **WHEN** queue, memory or disk limits are reached
- **THEN** intake applies documented backpressure/rejection, slow viewers cannot block collection and the UI exposes loss or unknown completeness

### Requirement: Compose configuration and transport boundaries
The deployment contract SHALL use Docker Compose with a matching Compose file and `.env`, behind the existing shared Cloudflare Tunnel HTTPS route. It SHALL document persistent storage, required settings, origin transport and trusted proxy boundaries without credential defaults. Tunnel TLS SHALL NOT replace application auth or establish routed SSE/OTLP compatibility.

#### Scenario: Configuration or origin transport mismatch
- **WHEN** a future setup lacks required values or the configured origin transport does not match the connector path
- **THEN** validation/readiness reports the mismatch without exposing secrets, and public HTTPS is not described as proof of authenticated access or origin TLS

#### Scenario: Shared connector and persistent state
- **WHEN** a future Compose app restarts or updates behind the existing connector
- **THEN** the deployment contract preserves documented state/lifecycle behavior and does not take ownership of sibling tunnel routes or silently provision/change the shared connector

### Requirement: Explicit greenfield cutover
Fresh installations SHALL initialize without importing unsupported legacy single-app state. Existing managed two-key installations SHALL use an explicit reviewed transition to one app key; automatic reuse without that transition SHALL fail closed. Unsupported state SHALL report a recovery error without resetting data or regenerating credentials.

#### Scenario: Old test database selected
- **WHEN** startup finds an unsupported single-app database
- **THEN** it requests a fresh directory or supported recovery path without modifying the old database

#### Scenario: Managed installation awaiting cutover
- **WHEN** the new release finds an existing two-key registry without a completed migration marker
- **THEN** it requests the reviewed offline transition and leaves state unchanged rather than starting intake with guessed credentials

### Requirement: Versioned retained structural policy
New evidence SHALL identify its projection policy; derived results SHALL identify definition and computation/timing policy. Privacy filtering SHALL precede storage, preview, live reads and copy. Projection changes SHALL NOT silently backfill or reinterpret retained history.

#### Scenario: Omitted field requested by a mapping
- **WHEN** an owner previews a condition over a field not retained by the applicable projection
- **THEN** the field is unavailable with a reason rather than read from hidden raw telemetry or invented by correlation

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

### Requirement: Bounded structural lifecycle
Trace structure and relationship targets SHALL follow 30-day detail expiry and supported deletion/correction horizons. Twelve-month summaries SHALL retain only supported identity-free, text-free aggregates separated by timing/definition policy. Structural links SHALL NOT extend detail retention.

#### Scenario: Late correction and expiry
- **WHEN** a missing parent arrives within the retained horizon and the linked detail later expires
- **THEN** bounded reads first reconcile at a new revision and later remove expired navigation/detail, without extending the lifetime through replay or summaries

#### Scenario: Deletion while following a trace
- **WHEN** an owner deletes the supported correlated scope or application containing a trace
- **THEN** subsequent reads and live refreshes exclude its detail and dependent relationships under the existing replay-protection horizon, with no promise of unsupported individual deletion from finalized summaries

### Requirement: Greenfield persistent initialization
A new installation SHALL start with an empty private registry and no producer environment tokens. Owner registration SHALL generate all app credentials. Subsequent starts SHALL use JSON and SHALL NOT import environment values. Missing or inconsistent registry state SHALL fail closed.

#### Scenario: Fresh installation
- **WHEN** no registry, initialization marker, prior database or managed stores exist
- **THEN** owner login opens app registration with no synthetic app or producer tokens required

#### Scenario: Obsolete environment values
- **WHEN** a registered installation restarts with obsolete producer token environment values
- **THEN** only the JSON credentials authorize intake and no credentials are overwritten

### Requirement: Complete app state isolation
Every app SHALL own its evidence, discovery, mappings, definition versions, traces, journeys, diagnostics, references, health counters, revisions, summaries, deletion barriers and invitations. Equal local identifiers SHALL NOT join, deduplicate, resolve links or affect another app. Cleanup, replay and deletion SHALL remain app-scoped with existing retention and privacy semantics.

#### Scenario: Deliberately colliding identities
- **WHEN** A and B contain the same service, trace, span, event, session, view, definition, invitation and revision identifiers
- **THEN** reads, derivation, replay and lifecycle operations on A cannot disclose or change B's state
- **AND** a reference or trace link cannot resolve across the app boundary

#### Scenario: Delete or expire one app's evidence
- **WHEN** cleanup expires detail or the owner deletes a supported visit in A
- **THEN** A's own detail, aggregates and replay barriers reconcile under existing rules while B remains unchanged

### Requirement: Coordinated persistent state recovery
Backup/restore guidance SHALL treat the JSON registry and its referenced app stores as one protected state set. Restores SHALL preserve app-to-store identity, credentials, generations and deletion state, apply expiry before reads, and reject mismatched state. Registry failure SHALL NOT trigger reconstruction or regeneration from database contents.

#### Scenario: Restore mismatched files
- **WHEN** a restored registry identifies a generation or app that does not match its database
- **THEN** the affected app is unavailable with a sanitized recovery error and is not served under another registration

#### Scenario: Restore predating deletion or rotation
- **WHEN** an older backup predates app deletion, credential rotation or invitation revocation
- **THEN** restore guidance requires reconciliation against the latest deletion and credential state before enabling access
- **AND** an old backup alone is not claimed to preserve later revocations

### Requirement: Bounded multi-app operation
The installation SHALL enforce shared intake, query and streaming concurrency limits alongside per-app storage/query bounds. App registration SHALL have a documented finite limit, and cleanup SHALL schedule work fairly across apps. Adding apps SHALL NOT multiply receiver listeners, unbounded store handles or existing global resource budgets.

#### Scenario: One noisy app
- **WHEN** one app exceeds an intake or query bound while another is active
- **THEN** backpressure remains explicit and work stays within installation limits, with app-attributed health feedback and no cross-app data fallback

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
