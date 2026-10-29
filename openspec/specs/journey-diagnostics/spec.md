# journey-diagnostics Specification

## Purpose

Help a developer inspect relevant exceptions and subsequent recovery using available technical context without implying causal links unsupported by evidence. This is a proposed contract.

## Requirements

### Requirement: Evidence-based exception linkage
The system SHALL associate client/server exceptions with a journey point only through valid configured context, such as explicit app/visit/action correlation or a verified trace/span relationship. Proximity alone SHALL NOT establish linkage or cause. Unlinked exceptions SHALL remain identified as unlinked, and ambiguity SHALL remain visible.

#### Scenario: Related technical execution
- **WHEN** an exception shares a verified request trace with an action already linked to the app visit
- **THEN** it can be shown at that action with the linkage basis and available context

#### Scenario: Nearby unrelated exception
- **WHEN** another service emits an exception at a similar time without valid linkage
- **THEN** the system does not attach it as that journey's cause

### Requirement: Inspectable scrubbed diagnostic detail
The system SHALL show available exception type, scrubbed message/stack, runtime, version, service and deployment context. Missing fields SHALL remain unavailable. Detail and copy SHALL use the same permissions and versioned privacy projection, redacting known credentials and personal identifiers before storage. Raw request content SHALL remain excluded; the interface SHALL disclose limits of text scrubbing.

#### Scenario: Controlled diagnostic text
- **WHEN** a producer emits diagnostic text
- **THEN** guidance requires controlled technical messages rather than customer content and the interface discloses that pattern-based scrubbing cannot guarantee removal of every sensitive value

#### Scenario: Copy context for a developer or coding agent
- **WHEN** an authorized viewer copies an exception's diagnostic context
- **THEN** the result contains the permitted technical evidence, linkage basis and gaps without credentials, invitation tokens or personal application content

#### Scenario: Minified browser stack
- **WHEN** source maps or original frames are unavailable
- **THEN** the available stack is labeled accurately and original source locations are not fabricated

### Requirement: Recurrence and affected visits
The system SHALL distinguish received diagnostic records, qualified exception occurrences and affected anonymous visits within declared scope. Occurrence deduplication SHALL require supported identity or a qualified single carrier; a grouping fingerprint SHALL NOT prove one occurrence. Missing linkage and duplicate uncertainty SHALL remain visible without person-level counts.

#### Scenario: Same exception repeats in one visit
- **WHEN** a correlated exception occurs three times within one visit
- **THEN** recurrence can show three occurrences and one affected visit rather than three affected people

#### Scenario: Missing visit identifier
- **WHEN** an exception has no valid visit correlation
- **THEN** it contributes to unlinked diagnostic-record counts, with occurrence uncertainty, rather than an invented affected-visit count

#### Scenario: One exception in two signals
- **WHEN** a log and span event may describe one exception
- **THEN** they remain two diagnostic records unless supported occurrence identity qualifies deduplication; the interface does not claim two distinct exceptions or silently merge them by proximity

#### Scenario: Same stack occurs twice
- **WHEN** two genuine failures share an action and stack signature
- **THEN** grouping can combine their category while retaining two supported occurrences rather than deduplicating by signature alone

### Requirement: Continued activity and recovery
The system SHALL preserve activity after exceptions, show later goal evidence and distinguish recovered journeys from no further observed activity. No further activity SHALL NOT prove that an exception caused departure. A confidence badge, AI diagnosis or root-cause assertion SHALL NOT be required.

#### Scenario: Goal follows an error
- **WHEN** a linked exception is followed by another action and an evidenced goal
- **THEN** the journey shows continuation and completion with the exception retained as friction

#### Scenario: Nothing observed afterward
- **WHEN** no further events follow a linked exception within the observation window
- **THEN** the view says no further activity observed and includes collection/window limitations without claiming a causal exit

### Requirement: Bounded application-wide diagnostic discovery
The owner SHALL discover retained errors and exceptions independently of visit correlation through protected search. Standard exception attributes, failed-span status descriptions and error-log string bodies SHALL use a versioned scrubbed projection before persistence. Messages SHALL be limited to 2048 UTF-8 bytes; stacks to 8192 bytes and forty lines, processing at most 64 KiB per field. Ordinary informational bodies and unknown attributes SHALL remain omitted.

#### Scenario: Background error without a visit
- **WHEN** a standard error log arrives without a visit identifier
- **THEN** its permitted message remains searchable in application-wide errors without inventing a journey or person

#### Scenario: Exception message and stack
- **WHEN** a standard exception carries a message or stack, with or without an exception type
- **THEN** supported fields are bounded and scrubbed, unavailable fields remain explicit, source maps are not fabricated and metrics remain receipt-only

#### Scenario: Error details expire
- **WHEN** detailed evidence expires after thirty days or is deleted
- **THEN** diagnostic search excludes it immediately and twelve-month summaries retain no diagnostic free text

#### Scenario: Projection upgrade and replay
- **WHEN** identified diagnostic evidence accepted by the new model is retransmitted unchanged
- **THEN** one retained record remains without overwriting or backfilling accepted data, and conflicting retained fields reject atomically; loading or replaying pre-cutover legacy projections is not required
