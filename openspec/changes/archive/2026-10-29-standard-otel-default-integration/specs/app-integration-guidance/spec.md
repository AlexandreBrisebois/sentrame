# Application Integration Guidance

## ADDED Requirements

### Requirement: Progressive integration documentation
README, developer/user guides, producer contract, agent guide and deployment examples SHALL share terminology and one voting example. They SHALL distinguish current release behavior, proposed behavior, optional application conventions, source-checked examples and runtime-qualified cells.

#### Scenario: Read before implementation
- **WHEN** a developer reads planning artifacts for this change
- **THEN** proposed native timing/trace support is not presented as available in the current release, and source snippets are not described as tested exporters

### Requirement: Practical independent server recipes
The guides SHALL cover existing .NET and Next.js Node SDK/provider reuse, logs/traces, correct HTTP endpoints/protocol/auth, random token generation, secret storage, standard cross-service propagation, named outcomes, exceptions and controlled error messages. Browser instrumentation SHALL have a separate trust/lifecycle path.

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
