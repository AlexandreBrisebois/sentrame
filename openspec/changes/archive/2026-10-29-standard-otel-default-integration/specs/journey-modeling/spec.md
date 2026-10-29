# Journey Modeling

## ADDED Requirements

### Requirement: Standard-name business mappings
Owners SHALL map standard span/event names and supported attributes using explicit service, carrier and trusted-source selectors without `stage.name`. They SHALL preview stages, supporting activity, goals, exceptions and missing prerequisites before prospective activation.

#### Scenario: Ordinary voting operation
- **WHEN** a server exports `POST /votes`, a persistence span and a named `vote.committed` event without Sentrame attributes
- **THEN** the owner can label the operations and map the trusted commitment to a business milestone within the supported trace context, without creating a multi-request journey from that trace alone

#### Scenario: Unproven persistence
- **WHEN** only a browser button event and generic HTTP success exist
- **THEN** configuration explains the absent trusted business-outcome evidence rather than asserting that a vote persisted

### Requirement: Optional stage hint suggestions
The system SHALL retain optional `stage.name` as a Sentrame application convention and suggest labels for unmapped evidence. Suggestions SHALL require owner acceptance into literal versioned rules; the attribute alone SHALL NOT assign Kind, goal, session, boundary, authority or timing quality.

#### Scenario: Enriched voting event
- **WHEN** an unmapped `vote.committed` event carries `stage.name=voting.persist`
- **THEN** configuration suggests that label and shows its source, and activation requires explicit owner preview/acceptance rather than silently creating a graph stage

#### Scenario: Missing and conflicting labels
- **WHEN** the same selector emits absent, empty, invalid or several distinct stage labels
- **THEN** missing/invalid values yield no suggestion, distinct valid labels remain separate suggestions and the owner can choose a literal label without a popularity-based default activation

### Requirement: Authoritative explicit mappings
Exactly one explicit matching owner mapping SHALL override any emitted stage hint. Overlapping explicit rules SHALL block activation until disambiguated. Missing hints SHALL NOT disable explicit rules; child evidence SHALL NOT inherit labels merely through trace correlation.

#### Scenario: Owner label differs
- **WHEN** a rule maps `vote.committed` to `Vote saved` while the event emits `stage.name=voting`
- **THEN** the stage remains `Vote saved`, the emitted hint remains inspectable and no historical label changes

#### Scenario: Two explicit matches
- **WHEN** a generic rule and service-specific rule overlap
- **THEN** preview reports the conflict and activation requires disjoint selectors rather than hidden rule order or implicit specificity

### Requirement: Span-backed stage occurrences
A configured native-span stage SHALL use native span identity and its inclusive duration without paired events or custom occurrence IDs. Point-event stages SHALL NOT acquire duration from their label. Optional paired boundaries and estimates SHALL remain explicit separate timing policies.

#### Scenario: Repeated vote processing
- **WHEN** several distinct spans map to Voting without `stage.occurrence.id`
- **THEN** each span is a distinct operation-stage occurrence with qualified native timing, not a whole page-dwell interval

#### Scenario: Span and pair describe one operation
- **WHEN** native span and explicit paired events both describe a configured occurrence
- **THEN** its selected timing policy contributes at most one completed observation for that occurrence and basis, with evidence links and no implicit pooling

### Requirement: Purposeful multi-request correlation
Multi-request journeys SHALL require explicit bounded app-scoped correlation and browser tab separation where applicable. `session.id` support SHALL disclose its opt-in Development client convention and Sentrame's server-scope use. No person, IP, user-agent, fingerprint, timestamp or automatic cross-session inference SHALL join evidence.

#### Scenario: Multiple request traces
- **WHEN** navigation, button press, failed request, retry and commitment span different traces without a shared supported scope
- **THEN** each trace remains explorable but no single browsing session or complete multi-request path is invented

#### Scenario: Explicit combined correlation
- **WHEN** those records carry correctly propagated random `session.id` and matching browser `tab.id`, with trusted source retained
- **THEN** the configured journey preserves intent, failure, retry, recovery and valid goal under first-mapped/goal/30-minute rules without joining unrelated tabs, apps or visits

### Requirement: Optional observed browser activity
Observed-active stage time SHALL require actual scoped browser presence evidence and a supported stage interval. It SHALL remain separate from span duration, stage elapsed time and attention. Missing observations SHALL disable only observed-active measurement.

#### Scenario: Request ends while page stays open
- **WHEN** a voting-page request ends and no visibility/idle observations or page interval exist
- **THEN** operation timing remains available but page dwell and observed-active time remain unknown

#### Scenario: Visible hidden and idle transitions
- **WHEN** a properly scoped browser interval has actual visible/active observations followed by hidden/idle states
- **THEN** only supported positive intervals under the disclosed heartbeat policy contribute observed-active time, with no backend or click-based attention inference

### Requirement: Prospective interpretation versions
Definitions and timing policies SHALL preserve capture-time meaning. Hint acceptance and rule changes SHALL affect new captures prospectively; late evidence SHALL reconcile under recorded definitions with visible revisions, and saved references SHALL preserve captured boundaries.

#### Scenario: Label changes during an active journey
- **WHEN** the owner activates a new suggested label while a captured journey remains active
- **THEN** that journey retains its original definition until closure, later journeys use the new definition and a saved historical view does not silently rename or resegment the earlier journey
