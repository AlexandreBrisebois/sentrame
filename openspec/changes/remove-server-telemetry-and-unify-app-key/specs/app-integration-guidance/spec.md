# App integration guidance delta

## MODIFIED Requirements

### Requirement: Standard carriers and scoped correlation
The guidance SHALL separate supported standard OTLP intake from sufficient optional journey evidence. It SHALL cover client journeys, stable event/step meaning, bounded app-scoped correlation, timing and goals without a mandatory Sentrame SDK/protocol. Correlation SHALL NOT imply authentic origin, ordering, deduplication or person/cross-visit identity.

#### Scenario: Server-only or combined integration
- **WHEN** a developer correlates client observations across requests
- **THEN** the plan uses explicit bounded visit/tab context and preserves clocks and distinct identities; missing context does not invalidate supported OTLP reception

### Requirement: Three progressive voting increments
Examples SHALL first use ordinary OTel plus owner mapping, then optional `stage.name`, then purposeful journey correlation/presence. Client examples SHALL show stable stage naming and distinguish navigation, button intent, observed result, failure/retry/recovery and configured goal evidence without certifying persistence.

#### Scenario: Ordinary OTel answers the question
- **WHEN** existing client request spans and correlated result observations already answer an operation question
- **THEN** the recipe requests no stage tags, stage occurrence IDs or paired boundary logs, and states that page dwell and multi-request paths remain unknown

#### Scenario: Optional enrichment
- **WHEN** a developer adds `stage.name=voting.result` to an existing span or named result-observation event
- **THEN** the client examples show a label suggestion with owner precedence, no automatic goal/boundary/session semantics and no required exporter replacement

### Requirement: Reproducible qualification evidence
Implementation qualification SHALL use repeatable synthetic fixtures and version-controlled REST Client requests through actual ingestion, durable storage, derivation and protected reads. Reports SHALL separate wire cells, actual retained client exporters, browser behavior, Mac, container, Home Server and Tunnel, with passed/failed/blocked/not-run evidence.

#### Scenario: Wire succeeds but exporters unrun
- **WHEN** a wire fixture passes while no actual client exporter has run
- **THEN** only the wire cell is qualified, actual exporters remain not-run and documentation/mocks do not establish framework or deployed-route support

### Requirement: Bounded instrumentation agent checklist
The agent guide SHALL require source/SDK discovery, explicit question and evidence gaps, approved bounded emitter plan, minimal code changes and real verification before support claims. It SHALL prohibit duplicate initialization, credential exposure, fabricated outcomes and unnecessary custom context.

#### Scenario: Unfamiliar voting application
- **WHEN** an agent prepares integration in an unfamiliar application
- **THEN** it locates client navigation/result/retry/error boundaries, inventories current exporters, proposes only missing observations, checks privacy/authority and reports actual versus unrun evidence without changing business behavior

### Requirement: Configuration-only app onboarding
User/developer/deployment guides, producer contract and agent instructions SHALL explain UI registration, JSON credential persistence and standard exporter configuration. Registered client senders SHALL need no changes during restart. Existing two-key deployments SHALL follow the explicit reviewed lifecycle cutover. New-app setup SHALL require no Sentrame SDK, app-ID field or instrumentation rewrite. Key handling SHALL distinguish the private registry/owner interface from a potentially extractable client write key.

#### Scenario: Owner adds a product
- **WHEN** the owner follows onboarding for a second app
- **THEN** registration generates one app key and the UI supplies copyable OTLP endpoint/header settings without editing Sentrame's environment or running a secret-generation command
- **AND** guidance explains configuring the external sender's existing exporter or secret store with those values

#### Scenario: Greenfield deployment
- **WHEN** the owner follows the deployment guide
- **THEN** guidance uses a fresh persistent directory, UI registration and the generated app key without producer environment values or legacy import

#### Scenario: Browser and server settings
- **WHEN** the owner connects client components to one registered app
- **THEN** guidance uses the same app key with existing supported exporter paths, explains public-key spoofability and grants no stronger producer assurance
- **AND** multi-app support is not presented as qualification of a new browser exporter, CORS route, transport or deployed Tunnel

## REMOVED Requirements

### Requirement: Practical independent server recipes
**Reason**: Server instrumentation, trusted persistence and cross-service backend recipes are outside the client journey product.
**Migration**: Replace active onboarding with the retained client recipe below. Keep historical verification reports labeled as historical; forwarding infrastructure may carry client OTLP without becoming a supported server-telemetry journey.

## ADDED Requirements

### Requirement: Single-key client integration recipe
Active documentation and agent guidance SHALL use one app registration and telemetry key with ordinary qualified OTel exporter configuration. They SHALL cover base versus signal endpoints, optional journey context, goals, privacy, lifecycle and troubleshooting without a custom SDK, producer-role setup or server instrumentation. Direct browser/CORS/Tunnel support SHALL require its own qualification.

#### Scenario: Existing instrumented client
- **WHEN** an agent configures a client that already uses OTel
- **THEN** it reuses existing initialization and supported exporter settings, requests only missing approved observations and validates reception before diagnosing absent journey context
- **AND** a same-origin forwarder, if needed, forwards the client signal without granting trusted persistence or requiring backend instrumentation

#### Scenario: Extracted app key
- **WHEN** a caller obtains the key from distributed client code
- **THEN** guidance explains that it can submit spoofed app observations but cannot read history, obtain owner settings or administer the app
- **AND** rotation is described as replacing all senders' write permission, not authenticating installations or permanently excluding an attacker who can obtain the replacement

#### Scenario: Documentation and agent consistency
- **WHEN** the user follows README, user/developer/Home Server guides, OTEL-CONTRACT, AGENT-INTEGRATION or synthetic examples
- **THEN** the active path uses the same one-key client vocabulary and preserves the distinction between source inspection, local wire tests, actual browser/exporter tests and deployed-route qualification

