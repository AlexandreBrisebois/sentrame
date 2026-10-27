# Journey Diagnostics

## Purpose

Help a developer inspect relevant exceptions and subsequent recovery using available technical context without implying causal links unsupported by evidence. This is a proposed contract.

## ADDED Requirements

### Requirement: Evidence-based exception linkage
The system SHALL associate client/server exceptions with a journey point only through valid configured context, such as explicit app/visit/action correlation or a verified trace/span relationship. Proximity alone SHALL NOT establish linkage or cause. Unlinked exceptions SHALL remain identified as unlinked, and ambiguity SHALL remain visible.

#### Scenario: Related technical execution
- **WHEN** an exception shares a verified request trace with an action already linked to the app visit
- **THEN** it can be shown at that action with the linkage basis and available context

#### Scenario: Nearby unrelated exception
- **WHEN** another service emits an exception at a similar time without valid linkage
- **THEN** the system does not attach it as that journey's cause

### Requirement: Inspectable scrubbed diagnostic detail
The system SHALL show available exception type, sanitized message, stack, runtime, app version, service and deployment context. Missing fields SHALL remain unavailable. Detail and copy actions SHALL use the same privacy and permission policy and SHALL NOT expose secrets, household/person identifiers or unsanitized request content.

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
