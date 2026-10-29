# Trace Exploration

## Purpose

Expose bounded technical execution structure and timing from standard spans and correlated evidence without requiring business journey instrumentation.

## ADDED Requirements

### Requirement: Trace-scoped first value
The system SHALL explore app-scoped traces without session, stage, action or custom event identifiers. It SHALL distinguish trace exploration from business journeys and preserve producer authority on every record.

#### Scenario: Ordinary request without enrichment
- **WHEN** standard server spans and correlated error logs arrive without custom attributes
- **THEN** protected discovery and trace reads expose services, operations, available diagnostics and native timing while no visit or journey is fabricated

### Requirement: Typed structural relationships
The system SHALL distinguish parent-child edges, attached span events, correlated logs and span links. Unknown targets SHALL remain unresolved; timing proximity and matching labels SHALL NOT manufacture edges. Cross-trace links SHALL NOT merge journeys.

#### Scenario: Parallel children and asynchronous link
- **WHEN** a parent has two overlapping children and one child links to another trace
- **THEN** the children remain siblings, the link is a separate reference and no sibling-to-sibling chronological transition or shared session is asserted

#### Scenario: Missing and ambiguous parent
- **WHEN** a declared parent is missing or several producer-scoped candidates match its ID
- **THEN** the child remains visible with unresolved or ambiguous parentage rather than acquiring the nearest timestamped parent

#### Scenario: Events and logs
- **WHEN** exported span events are out of timestamp order and logs contain exact span or trace-only context
- **THEN** event export order and timestamps remain inspectable, exact logs attach to their span and trace-only logs remain trace-level evidence without invented parentage

### Requirement: Native span timing
The system SHALL measure inclusive operation duration from valid native start/end timestamps using span identity. It SHALL NOT require paired custom events or stage occurrence IDs, compute exclusive time in this scope, or present operation duration as browser dwell, attention or unique total time.

#### Scenario: Nested repeated operations
- **WHEN** two distinct spans named `vote.persist` repeat under one request and overlap another operation
- **THEN** each span retains its own duration and occurrence identity without summing inclusive durations as nonoverlapping journey time

#### Scenario: Missing or invalid end
- **WHEN** a retained span lacks an end or has an end preceding its start
- **THEN** it is censored or timing-unavailable, contributes no completed duration sample and is neither a zero nor negative measurement

#### Scenario: Unexported unfinished work
- **WHEN** an SDK has exported no record for an unfinished operation
- **THEN** Sentrame does not infer a live start or running span from the absence of an exported end

### Requirement: Partial trace disclosure
Trace reads SHALL expose sampling/completeness uncertainty, dropped relationships, missing targets and bounded-result truncation. Trace counts SHALL describe observed evidence, never population conversion or people. Late evidence SHALL reconcile at a visible revision without changing configured historical meaning.

#### Scenario: Child arrives before parent
- **WHEN** a child and error arrive before the parent, which arrives later within the retention horizon
- **THEN** the initial partial structure gains its supported parent relationship at a new revision without duplicating the child or inventing a business boundary

### Requirement: Protected bounded trace access
Trace lists, detail, relationship navigation and copy SHALL enforce the same app/session permissions, expiry, privacy and revocation as other protected reads. Read-only viewers SHALL gain no configuration permission from trace exploration.

#### Scenario: Read without authorization
- **WHEN** an unauthenticated or revoked viewer requests a trace, linked target or copied diagnostic context
- **THEN** access is denied even when the requester knows a valid trace identifier
