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
The system SHALL register stable application identities and map authenticated intake routes and configured resource/service selectors to them. Missing optional metadata SHALL NOT prevent intake when configuration establishes the application. Ambiguous or unauthorized mappings SHALL NOT mix application data. Resource claims alone SHALL NOT grant access.

#### Scenario: Several services for one app
- **WHEN** a configured browser service and API service map to WFS
- **THEN** both belong to the same registered app while retaining service identity

#### Scenario: Same service name in different apps
- **WHEN** two authorized routes receive the same service name for separate applications
- **THEN** their events, queries, diagnostics and counts remain separate

#### Scenario: Missing service metadata
- **WHEN** an authorized route binds a valid standard record to an app but its service or version is absent
- **THEN** the record is accepted with unknown metadata rather than invented values

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
The system SHALL preserve source IDs and deduplicate supported identities within app, trusted configured producer namespace and signal. Source-record identity, action correlation and semantic occurrence SHALL remain distinct. Basic ingestion SHALL NOT require custom IDs or claim exactly-once counts without them. Conflicts SHALL NOT overwrite accepted evidence or retain unsafe payloads.

#### Scenario: Retry with stable identity
- **WHEN** the same app, trusted producer namespace, signal and stable source identity are retransmitted unchanged
- **THEN** one evidence item and one set of derived contributions remain after replay and restart

#### Scenario: Legitimate repeated identical logs
- **WHEN** two identical-looking log records have no reliable source identity
- **THEN** intake does not silently collapse a potentially healthy repeated action and marks retry-related count uncertainty

#### Scenario: Equal IDs from different producers
- **WHEN** browser and server records use the same local ID in different configured producer namespaces
- **THEN** they coexist and the browser record cannot suppress or acquire the authority of the server record; forwarding preserves the declared trust distinction

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
The system SHALL retain supported native span times, identity, parent, kind/status, ordered events and bounded link targets, plus standard log context. Projection SHALL preserve exact timestamps and producer namespaces without retaining raw OTLP, arbitrary attributes or link payloads.

#### Scenario: Standard span structure
- **WHEN** a span carries start/end times, parent ID, events and links
- **THEN** the protected retained projection exposes supported structure and exact times, distinguishes omitted fields and does not require synthetic boundary events

#### Scenario: Structural bound exceeded
- **WHEN** more than 32 links appear on one otherwise acceptable span
- **THEN** at most 32 link targets are retained with an omission count, without rejecting the record solely for optional link truncation or implying complete structure

### Requirement: Contextual safe discovery
Discovery SHALL distinguish service, signal/carrier, span/event name, trusted producer and retained field availability. It SHALL expose unknown/omitted values honestly. Mappings SHALL use only supported privacy-projected fields; metrics SHALL remain receipt-only.

#### Scenario: Same name across carriers
- **WHEN** a browser log and server span have the same name
- **THEN** discovery separates their carrier/source context and permits intentional owner selectors without presenting the browser record as server evidence

#### Scenario: Route and secret fields
- **WHEN** telemetry contains an approved route template, HTTP method/status, raw URL/query and an authorization value
- **THEN** only supported validated selectors survive projection, and secrets/raw URL content remain absent from storage, preview, errors, copies and live reads

### Requirement: Credential-derived authority
Producer authority SHALL derive from authenticated intake configuration and SHALL NOT come from telemetry attributes, stage labels, trace IDs or claimed parentage. Browser forwarding SHALL preserve browser authority.

#### Scenario: Forged persistence claim
- **WHEN** browser-authenticated telemetry claims server source, `stage.name=voting.persist` and a trusted span's parent ID
- **THEN** it remains browser evidence and cannot satisfy a server-only persistence goal or suppress an equal-ID server record
