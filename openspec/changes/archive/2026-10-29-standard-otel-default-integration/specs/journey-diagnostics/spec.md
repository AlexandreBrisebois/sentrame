# Journey Diagnostics

## MODIFIED Requirements

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
