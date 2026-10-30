# app-integration-guidance Specification

## Purpose

Provide reusable discovery-led instructions for coding agents helping developers instrument other applications with existing standard OpenTelemetry, separate from Sentrame's own engineering harness.

## Requirements

### Requirement: Progressive browser navigation without a goal
Integration guidance SHALL explain automatic browser navigation for an empty captured definition with safe explicit session and tab IDs. It SHALL distinguish emitted labels and ordering from configured business meaning, persistence, duration and goal conversion. A completion goal SHALL be optional for both navigation and mapping activation.

#### Scenario: Existing browser exporter supplies navigation context
- **WHEN** an application already emits correlated opening, screen and error logs
- **THEN** guidance requests no invented completion event to connect the navigation and explains that explicit mappings remain authoritative once configured
- **AND** it documents optional `stage.name` labels, unique increasing `source.sequence` ordering, event identity for replay, 30-minute inactivity and separate app/session/tab boundaries

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

### Requirement: Progressive integration documentation
README, developer/user guides, producer contract, agent guide and deployment examples SHALL share terminology and one voting example. They SHALL distinguish current release behavior, proposed behavior, optional application conventions, source-checked examples and runtime-qualified cells.

#### Scenario: Read before implementation
- **WHEN** a developer reads planning artifacts for this change
- **THEN** proposed native timing/trace support is not presented as available in the current release, and source snippets are not described as tested exporters

### Requirement: Practical independent server recipes
The guides SHALL cover existing .NET and Next.js Node SDK/provider reuse, logs/traces, correct HTTP endpoints/protocol/auth, UI app registration, generated credential retrieval and secret storage, standard cross-service propagation, named outcomes, exceptions and controlled error messages. Browser instrumentation SHALL have a separate trust/lifecycle path.

#### Scenario: Existing provider and browser bundle
- **WHEN** an agent finds existing OTel initialization and a Next.js Client Component
- **THEN** it reuses the existing provider, keeps server credentials out of browser code, validates Node versus Edge constraints and introduces no mandatory proprietary wrapper

### Requirement: Three progressive voting increments
Examples SHALL first use ordinary OTel plus owner mapping, then optional `stage.name`, then purposeful journey correlation/presence. Both language paths SHALL show stable stage naming and distinguish navigation, button intent, trusted persistence, failure/retry/recovery and configured goal evidence.

#### Scenario: Ordinary OTel answers the question
- **WHEN** existing request/persistence spans and correlated named outcomes already answer an operation question
- **THEN** the recipe requests no stage tags, stage occurrence IDs or paired boundary logs, and states that page dwell and multi-request paths remain unknown

#### Scenario: Optional enrichment
- **WHEN** a developer adds `stage.name=voting.persist` to an existing span or named commitment event
- **THEN** both .NET and Next.js examples show a label suggestion with owner precedence, no automatic goal/boundary/session semantics and no required exporter replacement

### Requirement: Troubleshooting by evidence level
Guidance SHALL distinguish transport rejection, privacy omission, sampling/batching loss, unmapped/unassigned evidence, missing span timing and broken correlation. It SHALL recommend inspecting the earliest failed seam before requesting new instrumentation.

#### Scenario: Record accepted but no journey
- **WHEN** a standard record is received without supported scope
- **THEN** troubleshooting verifies discovery/trace behavior and explains the missing journey prerequisite instead of treating valid ingestion as a protocol failure

### Requirement: Reproducible qualification evidence
Implementation qualification SHALL use repeatable synthetic fixtures and version-controlled REST Client requests through actual ingestion, durable storage, derivation and protected reads. Reports SHALL separate wire cells, actual .NET/Next.js exporters, browser behavior, Mac, container, Home Server and Tunnel, with passed/failed/blocked/not-run evidence.

#### Scenario: Wire succeeds but exporters unrun
- **WHEN** a wire fixture passes while no actual .NET or Next.js exporter has run
- **THEN** only the wire cell is qualified, actual exporters remain not-run and documentation/mocks do not establish framework or deployed-route support

### Requirement: Bounded instrumentation agent checklist
The agent guide SHALL require source/SDK discovery, explicit question and evidence gaps, approved bounded emitter plan, minimal code changes and real verification before support claims. It SHALL prohibit duplicate initialization, credential exposure, fabricated outcomes and unnecessary custom context.

#### Scenario: Unfamiliar voting application
- **WHEN** an agent prepares integration in an unfamiliar application
- **THEN** it locates persistence/retry/error boundaries, inventories current exporters, proposes only missing observations, checks privacy/authority and reports actual versus unrun evidence without changing business behavior

### Requirement: Configuration-only app onboarding
User/developer/deployment guides, producer contract and agent instructions SHALL explain UI registration, JSON credential persistence and standard exporter configuration. Registered senders SHALL need no changes during restart. Legacy adoption is outside this greenfield change. New-app setup SHALL require no Sentrame SDK, app-ID field or instrumentation rewrite. Secret handling SHALL distinguish Sentrame's registry from each sender's own configuration.

#### Scenario: Owner adds a product
- **WHEN** the owner follows onboarding for a second app
- **THEN** registration generates credentials and the UI supplies copyable OTLP endpoint/header settings without editing Sentrame's environment or running a secret-generation command
- **AND** guidance explains configuring the external sender's existing exporter or secret store with those values

#### Scenario: Greenfield deployment
- **WHEN** the owner follows the deployment guide
- **THEN** guidance uses a fresh persistent directory, UI registration and generated credentials without producer environment values or legacy import

#### Scenario: Browser and server settings
- **WHEN** the owner connects a browser frontend and server to the same app
- **THEN** guidance assigns their distinct generated credentials to their existing supported trust paths and never places server credentials in client code
- **AND** multi-app support is not presented as qualification of a new browser exporter, CORS route, transport or deployed Tunnel

### Requirement: Sender upgrade qualification
Qualification SHALL demonstrate actual unchanged qualified OTLP requests before and after restart without producer environment tokens, and while a second app is active. Evidence SHALL separate wire compatibility from real exporter and deployed-route qualification. Static checks or mocked UI tests SHALL NOT establish sender compatibility.

#### Scenario: Exercise unchanged and new senders
- **WHEN** the multi-app implementation is qualified locally
- **THEN** the report records unchanged endpoint/header/payload reception and stable-ID replay for the first registered app, plus isolated reception using a second generated app credential
- **AND** it records JSON persistence, absence of database credential copies and explicit passed/failed/blocked/not-run results for exporter, browser, container and Tunnel checks
