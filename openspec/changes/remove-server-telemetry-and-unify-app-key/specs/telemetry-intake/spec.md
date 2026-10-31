# Telemetry intake delta

## MODIFIED Requirements

### Requirement: Configured application attribution
The system SHALL bind each active telemetry key to exactly one stable app and generation. Standard intake SHALL resolve that binding independently of service/resource claims. Optional metadata SHALL NOT be required for attribution. Service selectors SHALL NOT move data into another app, and unauthorized keys SHALL reject intake. Intake SHALL assign no browser/server role.

#### Scenario: Several services for one app
- **WHEN** components use the same app key with different service names
- **THEN** accepted observations belong to that registered app and retain bounded service metadata without separate role credentials

#### Scenario: Same service name in different apps
- **WHEN** two authorized keys receive the same service name for separate applications
- **THEN** their events, queries, diagnostics and counts remain separate

#### Scenario: Missing service metadata
- **WHEN** a valid standard record has no service or version metadata
- **THEN** it is accepted under its key-bound app with unknown metadata rather than invented values

#### Scenario: Same endpoint and different app credentials
- **WHEN** A and B use the same OTLP signal URL with overlapping service, trace, span, session and event IDs
- **THEN** each request reaches only its key-bound app and payload claims cannot override the app or generation

#### Scenario: Multiple products in a collector
- **WHEN** a Collector exports multiple products
- **THEN** requests are separated by app key before export; resource groups cannot distribute one request across apps

### Requirement: Honest delivery identity
The system SHALL preserve source IDs and deduplicate supported identities within app, generation and signal under its published identity policy. Record identity, action correlation and semantic occurrence SHALL remain distinct. Basic reception SHALL NOT require custom IDs or claim exactly-once counts without them. Conflicts SHALL NOT overwrite accepted evidence. Removing roles SHALL NOT collapse legacy identity domains.

#### Scenario: Retry with stable identity
- **WHEN** unchanged identified observations are retransmitted before or after restart or app-key rotation
- **THEN** one evidence item and one set of derived contributions remain under the same app identity

#### Scenario: Legitimate repeated identical logs
- **WHEN** identical-looking logs have no reliable source identity
- **THEN** intake does not silently collapse healthy repeated actions and discloses retry-related count uncertainty

#### Scenario: Equal IDs from different producers
- **WHEN** old client and server records shared the same local ID
- **THEN** compatibility preserves their distinct identities and only eligible client evidence can contribute to current analysis
- **AND** replay under the retained app key cannot resurrect excluded server evidence or duplicate already accepted client contributions

### Requirement: Contextual safe discovery
Discovery SHALL distinguish service, signal/carrier, span/event name and retained field availability without browser/server roles. Unknown and omitted values SHALL remain explicit. Mappings SHALL use supported privacy-projected fields; metrics SHALL remain receipt-only.

#### Scenario: Same name across carriers
- **WHEN** a log and a span use the same name
- **THEN** discovery preserves their service/carrier context for intentional owner selection without role or trusted-persistence claims

#### Scenario: Route and secret fields
- **WHEN** telemetry contains an approved route template, HTTP method/status, raw URL/query and an authorization value
- **THEN** only supported validated selectors survive projection, and secrets/raw URL content remain absent from storage, preview, errors, copies and live reads

### Requirement: Minimized structural receipt
The system SHALL retain supported native span times, identity, parent, kind/status, ordered events and bounded link targets, plus standard log context. Projection SHALL preserve exact timestamps and compatible identity domains without retaining raw OTLP, arbitrary attributes or link payloads.

#### Scenario: Standard span structure
- **WHEN** a span carries start/end times, parent ID, events and links
- **THEN** the protected retained projection exposes supported structure and exact times, distinguishes omitted fields and does not require synthetic boundary events

#### Scenario: Structural bound exceeded
- **WHEN** more than 32 links appear on one otherwise acceptable span
- **THEN** at most 32 link targets are retained with an omission count, without rejecting the record solely for optional link truncation or implying complete structure

### Requirement: Credential-derived authority
The app key SHALL authorize bounded writes to its app only. Observations SHALL remain untrusted, including through a forwarder. Payload attributes, stage labels, span kind, trace IDs, service metadata and claimed parentage SHALL NOT grant read/admin access or establish persistence, an authentic installation or a human action.

#### Scenario: Forged persistence claim
- **WHEN** accepted telemetry claims server source, a persistence label or another span's parent ID
- **THEN** it remains an untrusted app observation and cannot gain stronger assurance or bypass identity validation

### Requirement: Standard sender configuration and restart continuity
Registered apps SHALL use existing qualified OTLP endpoints, bearer headers and wire behavior with one app key. Restart and registering another app SHALL preserve settings. Existing two-key installations SHALL require the explicit lifecycle cutover before role removal. Subsequent rotation SHALL be an explicit owner action.

#### Scenario: Unchanged sender across restart
- **WHEN** an eligible client sender uses a qualified request before and after restart
- **THEN** its endpoint, header and payload continue to work for the same app and identified replay stays idempotent

#### Scenario: Configure a newly registered app
- **WHEN** a second app uses an existing standard exporter
- **THEN** ordinary endpoint/protocol/header configuration is sufficient for supported reception without a Sentrame SDK, app-ID field, custom envelope or journey fields
- **AND** journey analysis separately requires explicit visit/tab correlation and the applicable mappings

#### Scenario: Rotation and intentional reassignment
- **WHEN** an owner explicitly rotates a key or configures a sender for another app
- **THEN** guidance identifies the exporter configuration change without implying that stored data moves or instrumentation must change

## ADDED Requirements

### Requirement: Role-free standard receipt boundary
Supported logs, traces and metrics SHALL remain receivable without a runtime-role claim. Span kind SERVER, runtime names or missing journey fields SHALL NOT alone reject otherwise valid OTLP. Reception SHALL NOT provide a server-workflow product or certify producer origin. Missing journey prerequisites SHALL remain explainable; metrics SHALL remain receipt-only.

#### Scenario: Generic exporter and no journey context
- **WHEN** valid qualified OTLP includes a SERVER-kind span or runtime metadata without visit/tab context
- **THEN** sanitized reception and available diagnostic context remain inspectable without a server journey, service dependency view or inferred user path
- **AND** documentation explains that payload fields cannot reliably enforce a client-only physical origin
