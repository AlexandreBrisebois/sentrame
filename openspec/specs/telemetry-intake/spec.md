# telemetry-intake Specification

## Purpose

Receive ordinary OpenTelemetry without a proprietary producer dependency, while making acceptance, app attribution and interpretation limits inspectable. This is a proposed contract, not implemented behavior.

## Requirements

### Requirement: Standard compatibility contract
The system SHALL publish a versioned signal, transport, encoding and compression matrix. The target SHALL include logs, traces and metrics over OTLP/HTTP binary Protobuf, OTLP/HTTP JSON Protobuf and OTLP/gRPC Protobuf, with none and gzip compression. Unsupported combinations SHALL return explicit errors. Profiles SHALL be excluded until separately specified.

#### Scenario: Unmodified standard producer
- **WHEN** an existing SDK/exporter sends a supported combination using ordinary endpoint, protocol and authentication configuration
- **THEN** intake accepts valid telemetry without a Sentrame SDK, custom transport, schema version, event ID, session ID or journey fields

#### Scenario: Scope of compatibility claim
- **WHEN** only HTTP binary logs have been qualified
- **THEN** the published qualified matrix identifies only that combination as verified and does not claim the full target is implemented

#### Scenario: Encoding and compression
- **WHEN** conforming JSON or gzip payloads arrive on a supported endpoint
- **THEN** intake uses standard OTLP encoding semantics and bounded decompression, preserves timestamp precision and ignores unknown protocol message fields as required

### Requirement: Route-specific compatibility claims
The system SHALL distinguish receiver/origin support from support through the shared Cloudflare Tunnel route. External HTTPS SHALL NOT imply OTLP gRPC support, browser exporter compatibility or app authorization. Only qualified protocol/encoding/compression paths SHALL be advertised for that route, using ordinary standard exporter configuration.

#### Scenario: gRPC origin behind a public tunnel hostname
- **WHEN** the origin supports gRPC but the selected Tunnel public-hostname route does not
- **THEN** the route is not advertised as gRPC-compatible; an available qualified standard HTTP exporter path can be documented without replacing OTLP or exposing another network path implicitly

#### Scenario: Standard HTTP export through the tunnel
- **WHEN** an existing exporter uses the published HTTPS OTLP/HTTP route
- **THEN** qualification checks paths/content types, binary/JSON and gzip handling, body/time limits, auth, errors/retries and acknowledgements; direct-origin success alone does not establish route success

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

### Requirement: Reception differs from analysis
The system SHALL expose received signal types, sanitized event names and field availability. It SHALL distinguish accepted, privacy-filtered, rejected, unmapped and journey-analyzable records. Lack of optional journey context SHALL NOT be a protocol error or become a fabricated visit, step or outcome. Metrics SHALL NOT reconstruct individual paths.

#### Scenario: Technical telemetry only
- **WHEN** ordinary HTTP spans and metrics arrive without visit correlation or business events
- **THEN** the interface shows reception and available technical context while explaining which journey questions lack evidence

#### Scenario: Unfamiliar attributes
- **WHEN** a valid record contains an attribute outside the configured retention allowlist
- **THEN** the field is omitted under the visible privacy policy without rejecting the entire standard record solely for that attribute

### Requirement: Durable bounded acknowledgement
The system SHALL publish each acknowledgement boundary and durability failure model. Direct storage-receiver success SHALL follow a commit meeting that model. Requests, decompression, concurrency, queues and storage SHALL be bounded. Recoverable commit/overload failures SHALL use explicit retryable statuses. Queue acceptance SHALL NOT be described as final storage.

#### Scenario: Full disk or unavailable writer
- **WHEN** accepted evidence cannot be durably committed
- **THEN** the direct receiver returns a retryable failure instead of success and reports bounded collection-health information

#### Scenario: Declared crash boundary
- **WHEN** the receiver advertises process-crash durability or the stronger power-loss boundary
- **THEN** its storage configuration and qualification evidence identify that exact boundary; a process restart alone does not prove power-loss durability

#### Scenario: Intermediate collector
- **WHEN** a Collector acknowledges before Sentrame commits
- **THEN** documentation and health views identify the Collector boundary and its queue durability rather than describing that response as Sentrame persistence

### Requirement: Atomic direct-request validation
The initial direct receiver SHALL prevalidate and atomically accept or reject each storage-receiver request. Permanent record errors SHALL reject the whole request with explicit InvalidArgument/HTTP 400 and zero committed records. Missing optional business fields and privacy omissions SHALL NOT alone cause rejection. Per-item partial acceptance SHALL NOT be claimed without a qualified response integration.

#### Scenario: Mixed request
- **WHEN** a direct request contains acceptable records and a permanent record-validation error
- **THEN** zero records commit and an explicit standard non-retryable validation response is returned without unsafe payload disclosure

#### Scenario: Recoverable writer failure
- **WHEN** the writer is unavailable before commit
- **THEN** zero records commit and Unavailable/HTTP 503 requests a retry; a generic permanent error is not assumed to produce HTTP 400

#### Scenario: Collector rebatching
- **WHEN** a Collector merges or splits SDK envelopes before forwarding
- **THEN** atomicity applies only to each resulting direct storage-receiver request and the original envelopes are not described as atomic through that pipeline

### Requirement: Honest delivery identity
The system SHALL preserve source IDs and deduplicate supported identities within app, generation and signal under its published identity policy. Record identity, action correlation and semantic occurrence SHALL remain distinct. Basic reception SHALL NOT require custom IDs or claim exactly-once counts without them. Conflicts SHALL NOT overwrite accepted evidence.

#### Scenario: Retry with stable identity
- **WHEN** unchanged identified observations are retransmitted before or after restart or app-key rotation
- **THEN** one evidence item and one set of derived contributions remain under the same app identity

#### Scenario: Legitimate repeated identical logs
- **WHEN** identical-looking logs have no reliable source identity
- **THEN** intake does not silently collapse healthy repeated actions and discloses retry-related count uncertainty


### Requirement: Bounded metric receipt without unsafe arithmetic
The system SHALL receive supported metrics but initially limit their analysis to bounded sanitized receipt inspection. Metric arithmetic SHALL remain unavailable unless retained identity, type, temporality and interval semantics support it. Forbidden fields SHALL NOT be retained to rescue analysis, and privacy projection SHALL NOT silently merge distinct streams.

#### Scenario: Privacy collapses series identity
- **WHEN** omitting a forbidden attribute makes two metric streams appear identical
- **THEN** receipt remains inspectable but the affected value-derived analysis is unavailable unless a valid reaggregation policy is qualified

#### Scenario: Cumulative and delta intervals
- **WHEN** cumulative values 10 then 12, a repeated delta interval, resets or overlapping late points are received
- **THEN** receipt counts remain separate from metric values and no increment/total is claimed until the corresponding semantics and duplicate policy are qualified

#### Scenario: Metric accounting bounds
- **WHEN** a point has many histogram buckets or missing/invalid measurement time
- **THEN** documented point, bucket and byte limits govern acceptance, and missing timing does not become a fabricated interval

### Requirement: Signal-specific collection health
The system SHALL distinguish producer batching, sampling, retries and loss from receiver and live-display latency. Missing diagnostics SHALL mean unknown completeness. Configuration guidance SHALL reuse existing signal-specific SDK/exporter processors, with browser lifecycle limitations and supported exporter versions identified.

#### Scenario: Delayed exporter batch
- **WHEN** a batch appears several seconds after the observed action
- **THEN** the live view distinguishes occurrence and receipt times and does not promise instantaneous clicks

#### Scenario: Sampled spans
- **WHEN** sampled traces are available but complete business evidence is not
- **THEN** journey rates are not presented as population conversion rates derived from the sampled traces

### Requirement: Progressive evidence availability
The system SHALL distinguish receipt/discovery, trace operations, configured business meaning, explicitly correlated multi-request journeys and optional browser presence. Missing higher-level evidence SHALL disable only dependent behavior, never invalidate otherwise supported OTLP.

#### Scenario: Standard telemetry only
- **WHEN** valid supported OTLP arrives without `stage.name`, session, stage occurrence or presence fields
- **THEN** it remains accepted and discoverable with available operation/error evidence, while unavailable journey/presence capabilities explain their missing prerequisites

### Requirement: Minimized structural receipt
The system SHALL retain supported native span times, identity, parent, kind/status, ordered events and bounded link targets, plus standard log context. Projection SHALL preserve exact timestamps and compatible identity domains without retaining raw OTLP, arbitrary attributes or link payloads.

#### Scenario: Standard span structure
- **WHEN** a span carries start/end times, parent ID, events and links
- **THEN** the protected retained projection exposes supported structure and exact times, distinguishes omitted fields and does not require synthetic boundary events

#### Scenario: Structural bound exceeded
- **WHEN** more than 32 links appear on one otherwise acceptable span
- **THEN** at most 32 link targets are retained with an omission count, without rejecting the record solely for optional link truncation or implying complete structure

### Requirement: Contextual safe discovery
Discovery SHALL distinguish service, signal/carrier, span/event name and retained field availability without browser/server roles. Unknown and omitted values SHALL remain explicit. Mappings SHALL use supported privacy-projected fields; metrics SHALL remain receipt-only.

#### Scenario: Same name across carriers
- **WHEN** a log and a span use the same name
- **THEN** discovery preserves their service/carrier context for intentional owner selection without role or trusted-persistence claims

#### Scenario: Route and secret fields
- **WHEN** telemetry contains an approved route template, HTTP method/status, raw URL/query and an authorization value
- **THEN** only supported validated selectors survive projection, and secrets/raw URL content remain absent from storage, preview, errors, copies and live reads

### Requirement: Credential-derived authority
The app key SHALL authorize bounded writes to its app only. Observations SHALL remain untrusted, including through a forwarder. Payload attributes, stage labels, span kind, trace IDs, service metadata and claimed parentage SHALL NOT grant read/admin access or establish persistence, an authentic installation or a human action.

#### Scenario: Forged persistence claim
- **WHEN** accepted telemetry claims server source, a persistence label or another span's parent ID
- **THEN** it remains an untrusted app observation and cannot gain stronger assurance or bypass identity validation

### Requirement: Standard sender configuration and restart continuity
Registered apps SHALL use existing qualified OTLP endpoints, bearer headers and wire behavior with one app key. Restart and registering another app SHALL preserve settings. Older two-key installations SHALL use a fresh data directory. Subsequent rotation SHALL be an explicit owner action.

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

### Requirement: Role-free standard receipt boundary
Supported logs, traces and metrics SHALL remain receivable without a runtime-role claim. Span kind SERVER, runtime names or missing journey fields SHALL NOT alone reject otherwise valid OTLP. Reception SHALL NOT provide a server-workflow product or certify producer origin. Missing journey prerequisites SHALL remain explainable; metrics SHALL remain receipt-only.

#### Scenario: Generic exporter and no journey context
- **WHEN** valid qualified OTLP includes a SERVER-kind span or runtime metadata without visit/tab context
- **THEN** sanitized reception and available diagnostic context remain inspectable without a server journey, service dependency view or inferred user path
- **AND** documentation explains that payload fields cannot reliably enforce a client-only physical origin
