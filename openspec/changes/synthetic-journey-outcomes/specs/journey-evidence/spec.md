# Journey Evidence

## Purpose

Provide bounded, inspectable event evidence for product journey analysis without collecting personal content or assuming that correlation identifies a person.

## ADDED Requirements

### Requirement: Bounded synthetic ingestion
The initial system SHALL accept only explicitly configured synthetic application scopes through a documented OTLP logs profile. It SHALL validate record types, required fields, supported schema versions, allowed attributes and request size before retaining evidence. Invalid records SHALL receive protocol-compatible rejection feedback without raw payload logging.

#### Scenario: Valid synthetic evidence
- **WHEN** a configured synthetic producer submits a valid attempt-start record
- **THEN** the system stores its allowed fields with a server receipt timestamp and acknowledges acceptance only after durable storage

#### Scenario: Unsupported or unsafe input
- **WHEN** a record has an unsupported version, unknown attribute, personal-content field or unconfigured application scope
- **THEN** the record is rejected without storing its payload or echoing it in diagnostics

#### Scenario: Oversized request
- **WHEN** a request exceeds the documented size limit
- **THEN** the request is rejected before event persistence and the response does not claim acceptance

### Requirement: Attempt-scoped correlation
The system SHALL correlate events only within an application scope and random attempt identifier. It SHALL NOT require persistent person, household, recipe or device identifiers, or link attempts across visits. Correlation identifiers SHALL be described as pseudonymous, not anonymous.

#### Scenario: Application boundaries
- **WHEN** two configured applications emit the same attempt identifier
- **THEN** their evidence and outcomes remain separate

#### Scenario: Missing start
- **WHEN** an otherwise valid outcome event arrives without a start event
- **THEN** its evidence is retained as an incomplete attempt without inventing a start timestamp

### Requirement: Idempotent event delivery
The system SHALL deduplicate records by application scope and event identifier. Identical retransmission SHALL not add evidence or increase outcome counts. Conflicting content under the same key SHALL be rejected without overwriting the accepted record. Delivery retries SHALL preserve the original event identifier.

#### Scenario: Lost acknowledgement
- **WHEN** an accepted event is retransmitted after its acknowledgement is lost
- **THEN** it is acknowledged without creating another event or outcome

#### Scenario: Conflicting duplicate
- **WHEN** a previously accepted event identifier is reused with different content
- **THEN** the original evidence remains unchanged and the conflict is reported without the raw payload

### Requirement: Delayed evidence and schema compatibility
The system SHALL preserve occurrence time and receipt time separately, accept supported delayed events within the retention window, and derive results independently of delivery order. It SHALL reject unsupported schema versions rather than interpreting them under another version. Clock inconsistency SHALL not produce fabricated timing measurements.

#### Scenario: Outcome arrives before start
- **WHEN** a persistence event arrives before its corresponding start event
- **THEN** later arrival of the start reconciles the same attempt without duplicating it

#### Scenario: Clock inconsistency
- **WHEN** a persisted event has an occurrence time earlier than its start
- **THEN** completion evidence is retained but elapsed duration is unavailable and the inconsistency is visible

### Requirement: Retention and deletion
The initial system SHALL expire evidence seven days after occurrence time and remove derived data that depends on expired evidence. It SHALL support deletion by application or attempt scope, remove that scope from summaries, and prevent late delivery from recreating a deleted attempt during the remaining acceptance window. No longer-lived aggregates SHALL be retained.

#### Scenario: Expiration
- **WHEN** evidence reaches seven days of age
- **THEN** it is excluded from queries and purged with dependent derived data by the next daily cleanup

#### Scenario: Deletion followed by retry
- **WHEN** an operator deletes an attempt and a producer retries its events within the acceptance window
- **THEN** the deleted evidence and outcome do not reappear

#### Scenario: Stale replay
- **WHEN** a record older than seven days is submitted
- **THEN** it is rejected and cannot extend its retention through retransmission

### Requirement: Default-off collection boundary
The synthetic slice SHALL NOT instrument WFS or collect real user telemetry. Any future real integration SHALL require separate approval, explicit opt-in and default-off collection, with participant consent scope reviewed before instrumentation. Disabling collection SHALL stop emission and clear pending telemetry buffers without implying historical deletion.

#### Scenario: Initial demonstration
- **WHEN** the synthetic demonstration runs
- **THEN** only generated fixture data is sent and no WFS endpoint, account or user device is contacted

#### Scenario: Future integration remains disabled
- **WHEN** no explicit opt-in exists for a future approved integration
- **THEN** no product journey telemetry is emitted or buffered
