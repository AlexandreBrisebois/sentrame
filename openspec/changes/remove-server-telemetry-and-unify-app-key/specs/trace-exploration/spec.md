# Trace exploration delta

## MODIFIED Requirements

### Requirement: Trace-scoped first value
The system SHALL explore app-scoped traces without session, stage, action or custom event identifiers. It SHALL distinguish trace exploration from business journeys and identify every record as an untrusted app observation.

#### Scenario: Ordinary request without enrichment
- **WHEN** standard client spans and correlated error logs arrive without custom attributes
- **THEN** protected discovery and trace reads expose services, operations, available diagnostics and native timing while no visit or journey is fabricated

### Requirement: Typed structural relationships
The system SHALL distinguish parent-child edges, attached span events, correlated logs and span links. Unknown targets SHALL remain unresolved; timing proximity and matching labels SHALL NOT manufacture edges. Cross-trace links SHALL NOT merge journeys.

#### Scenario: Parallel children and asynchronous link
- **WHEN** a parent has two overlapping children and one child links to another trace
- **THEN** the children remain siblings, the link is a separate reference and no sibling-to-sibling chronological transition or shared session is asserted

#### Scenario: Missing and ambiguous parent
- **WHEN** a declared parent is missing or several compatible app-scoped identity candidates match its ID
- **THEN** the child remains visible with unresolved or ambiguous parentage rather than acquiring the nearest timestamped parent

#### Scenario: Events and logs
- **WHEN** exported span events are out of timestamp order and logs contain exact span or trace-only context
- **THEN** event export order and timestamps remain inspectable, exact logs attach to their span and trace-only logs remain trace-level evidence without invented parentage

## ADDED Requirements

### Requirement: Client diagnostic trace boundary
Trace detail SHALL retain supported native span timing, events, logs and links as diagnostic context for app observations. It SHALL NOT present a server dependency map, backend persistence assurance or automatic server-workflow journey. Physical producer origin SHALL NOT be inferred from span kind. Legacy excluded records SHALL follow lifecycle compatibility rules.

#### Scenario: Client request trace
- **WHEN** client spans expose overlapping operations and an error
- **THEN** native inclusive duration, supported relationships and uncertainty remain inspectable with protected evidence navigation and no invented causal sequence or page-dwell measurement

#### Scenario: Generic SERVER-kind span
- **WHEN** a valid received span declares SERVER kind
- **THEN** kind remains protocol metadata and grants no extra assurance, server-workflow support or role-based UI
