# Journey Outcomes

## Purpose

Answer a declared product question from retained evidence, preserving uncertainty and distinguishing successful contribution from time spent in an application.

## ADDED Requirements

### Requirement: Explicit milestone definition
The initial journey SHALL ask whether a Discovery attempt reaches its first persisted vote. A start SHALL mean an observed entry into Discovery, not inferred voting intent. Completion SHALL require an authoritative persistence event for that attempt. “First” SHALL refer only to the attempt, not lifetime usage. The rule SHALL have a visible version.

#### Scenario: Optimistic action
- **WHEN** an attempt contains a start and a vote-requested event but no persistence evidence
- **THEN** the milestone is not classified as completed

#### Scenario: Persisted negative preference
- **WHEN** an authoritative event confirms that a vote was persisted
- **THEN** the milestone is completed regardless of whether the vote expressed a positive or negative preference

### Requirement: Evidence-based outcome classification
The system SHALL classify each attempt as completed when persistence is evidenced, otherwise unknown when terminal evidence conflicts, otherwise failed when terminal failure is evidenced, otherwise interrupted when explicit interruption is evidenced, and otherwise unknown. A request error, timeout, hidden tab or silence SHALL NOT alone establish terminal failure or interruption.

#### Scenario: Quick successful departure
- **WHEN** a persistence event is followed by no further observed activity
- **THEN** the milestone remains completed and subsequent activity is described as unobserved, not abandonment

#### Scenario: Error after persistence
- **WHEN** an attempt contains both persistence evidence and a request error
- **THEN** it is completed with the request error retained as separate friction evidence

#### Scenario: Terminal failure
- **WHEN** an authoritative terminal-failure event exists without persistence or conflicting interruption evidence
- **THEN** the attempt is classified as failed with its bounded reason code

#### Scenario: Explicit interruption
- **WHEN** an explicit attempt-interrupted event exists without persistence or terminal failure
- **THEN** the attempt is classified as interrupted without inferring why the person stopped

#### Scenario: Silence or transport error
- **WHEN** the only evidence after start is a timeout, request error, or no activity
- **THEN** the outcome is unknown and no claim of abandonment or persistence failure is made

#### Scenario: Conflicting terminal evidence
- **WHEN** interruption and terminal-failure events coexist without persistence
- **THEN** the outcome is unknown and both observations are shown as a conflict

### Requirement: Revisable outcomes
The system SHALL recompute an attempt's outcome when accepted evidence changes and SHALL produce the same result for the same evidence set regardless of delivery order or duplicates. Results SHALL remain labeled revisable within retention. Earlier incomplete classifications SHALL not be treated as immutable facts.

#### Scenario: Delayed success
- **WHEN** persistence evidence arrives for an attempt previously classified as unknown, failed or interrupted
- **THEN** the result becomes completed and retains the earlier friction evidence

### Requirement: Explainable summary and evidence
The interface SHALL show the product question, rule version, retained time range, outcome counts, denominator and incomplete attempts. Each attempt SHALL expose the allowed evidence and classification reason. It SHALL distinguish missing observations from reported collection errors and SHALL NOT claim that telemetry establishes intent, satisfaction or meaningful value.

#### Scenario: Summary denominator
- **WHEN** a fixture contains four attempts with start evidence and one persistence-only attempt
- **THEN** the started-attempt denominator is four and the persistence-only attempt is shown separately as incomplete, not silently included in that rate

#### Scenario: Optional timing
- **WHEN** start and persistence have consistent occurrence timestamps
- **THEN** elapsed time to the first observed persistence can be shown as a secondary measure without ranking longer visits as more successful

#### Scenario: Missing collection evidence
- **WHEN** no producer-side delivery diagnostics are available
- **THEN** the interface states that collection completeness is unknown rather than reporting zero lost events

### Requirement: Repeatable synthetic demonstration
The demonstration SHALL exercise real local ingestion, storage, derivation and presentation with generated data. It SHALL include completed, interrupted, failed, unknown, delayed, duplicate and incomplete attempts with documented expected results. Fixture data SHALL be clearly identified as synthetic.

#### Scenario: Repeatable run
- **WHEN** the same fixture event identifiers are replayed
- **THEN** outcome counts and denominators remain unchanged and the evidence view explains the expected classifications
