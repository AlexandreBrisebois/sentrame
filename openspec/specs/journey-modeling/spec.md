# journey-modeling Specification

## Purpose

Turn available telemetry into explicitly configured, explainable journey evidence without inventing business meaning or tracking people across visits. All requirements are proposed behavior.

## Requirements

### Requirement: Separate application visits and journeys
The system SHALL support client-only, server-only and combined journeys within one registered app and bounded logical correlation scope. Browser visits SHALL use separate random tab scopes; server workflows SHALL NOT require a browser/tab. Combined observations SHALL use the same correctly scoped correlation while preserving sources. Person/cross-visit/app joins SHALL NOT be inferred, and a request trace SHALL NOT automatically become a business journey.

#### Scenario: Session identifier collision across apps
- **WHEN** two applications emit the same session value
- **THEN** their journeys remain separate and no cross-application path is constructed

#### Scenario: Several goals in one visit
- **WHEN** a visit completes one configured journey and later starts meaningful activity
- **THEN** the next new mapped occurrence starts another journey without fabricating a new person or merging the outcomes

#### Scenario: Concurrent tabs
- **WHEN** concurrent tab activity shares a visit identifier without reliable ordering context
- **THEN** the browser evidence is marked ambiguous/unassigned until separate tab scopes are available; the system does not merge concurrent tabs into a certain user path

#### Scenario: Server-only journey
- **WHEN** mapped events identify a bounded server operation with no browser
- **THEN** automatic first-event/goal/inactivity rules apply within that operation scope without a fabricated tab, person or browser visit

#### Scenario: Combined client and server evidence
- **WHEN** properly emitted client/server observations share app-scoped logical correlation
- **THEN** they can form one journey while retaining producer authority, clock/order uncertainty and distinct record/action identities; the shared ID alone establishes neither duplicates nor causality

### Requirement: Self-explanatory event classification
The system SHALL list sanitized discovered events and fields so an authorized owner can map emitted events to stable journey steps, choose a goal event from that list and define explicit start/exit conditions. It SHALL preview matches, missing context and conflicts before definition activation. Ordinary standard OTel instrumentation SHALL suffice; a custom SDK, executable script or AI SHALL NOT be required.

#### Scenario: Mapping a received event
- **WHEN** an editor selects a received event, names its step and chooses an attribute condition
- **THEN** the editor sees representative sanitized matches, the plain-language rule and missing required context before saving a new definition version

#### Scenario: Conflicting mappings
- **WHEN** two mappings would assign incompatible semantics to the same evidence
- **THEN** activation requires explicit deterministic precedence or conflict resolution rather than hidden evaluation order

#### Scenario: Select an emitted goal event
- **WHEN** the owner selects a discovered event as the goal to inspect
- **THEN** the view identifies matching journeys and shows their ordered evidence from each known beginning to the selected occurrence, preserving filters and marking unavailable starts or omitted steps rather than inventing a prefix

### Requirement: Standard event and severity classification
The system SHALL distinguish OTel signal/carrier, event name, severity and owner-assigned journey meaning. Named log events, warnings and exceptions SHALL be discoverable without new transport types. Qualified span events SHALL retain source context. Warning/error severity or an exception type SHALL NOT alone establish journey failure, termination or business success.

#### Scenario: Warning followed by goal success
- **WHEN** a visit emits a warning event, continues and reaches its selected goal event
- **THEN** the warning remains in sequence and the evidenced goal can succeed without a severity-based failed-journey label

#### Scenario: Exception is an investigation goal
- **WHEN** the owner selects a discovered exception event to inspect how journeys reached it
- **THEN** the prefix is inspectable and later recovery remains available; event selection does not silently activate a terminal/failure rule

### Requirement: Explicit journey lifecycle
The system SHALL start at the first mapped event in a logical scope, end at the active goal and restart at the next new mapped occurrence. After 30 minutes without mapped activity it SHALL close inactive, not failed. Browser tabs SHALL remain separate; server-only scopes need no tab. Event-time/late revisions SHALL be explicit. Severity and saved-view queries SHALL NOT implicitly change boundaries.

#### Scenario: Timeout without terminal evidence
- **WHEN** a journey reaches 30 minutes without mapped activity or a goal
- **THEN** it becomes provisionally inactive with unknown outcome and does not acquire an invented abandonment reason

#### Scenario: Goal followed by new activity
- **WHEN** the active goal closes a journey and a new mapped occurrence follows in that tab
- **THEN** it starts another journey; supported replay of the prior goal is deduplicated first and does not create another journey

#### Scenario: Late goal after timeout
- **WHEN** a valid goal occurred in the prior active interval but arrives after provisional inactivity closure
- **THEN** the outcome and affected counts are revised consistently, with revision visible

#### Scenario: Explicit goal and continuing activity
- **WHEN** a goal is evidenced and activity continues
- **THEN** the next new mapped occurrence starts a new journey and the completed goal remains evidenced

### Requirement: Deterministic evidence-based outcomes
The system SHALL separate completion, positive business outcome, observed failure and configured flaky behavior. Each classification SHALL expose its rule version and evidence. A negative vote SHALL NOT imply a negative journey outcome. HTTP success, optimistic UI, repetition, silence or an exception alone SHALL NOT establish goal completion or failure.

#### Scenario: Recovered operation
- **WHEN** a request error is followed by valid persistence and later goal evidence
- **THEN** the journey shows the completed goal plus the earlier friction rather than ending at the error

#### Scenario: Publication fails after commit
- **WHEN** a vote commits before a follow-up publication fails and the client receives an error
- **THEN** valid commitment evidence still establishes the saved contribution without requiring a second vote, while the error and any later completion-screen observation remain separate

#### Scenario: Flaky classification
- **WHEN** a configured retry/error condition matches within its declared denominator and time window
- **THEN** the interface shows the matched condition and counts; without such a rule it does not label the journey flaky automatically

#### Scenario: Contradictory evidence
- **WHEN** terminal evidence conflicts under the active rule
- **THEN** the system shows the conflict and deterministic unresolved result rather than inventing a preferred business story

### Requirement: Healthy voting loops and explicit exit reasons
The system SHALL preserve vote-loop repetitions and valid duration. WFS completion-screen observation SHALL be a visit outcome distinct from shared agreement, its cause and emitted end reason. Mappings SHALL distinguish recipes exhausted, family agreement and departure before completion where emitted. Missing evidence SHALL NOT imply tiredness, frustration or a defect.

#### Scenario: Repeated healthy voting
- **WHEN** a journey records ten distinct valid vote commitments and then observes the completion screen
- **THEN** it shows ten votes, nine successive vote-to-vote transitions and one evidenced visit completion without a repetition-based failure label or a claim that this visit caused shared agreement

#### Scenario: Exhaustion differs from agreement
- **WHEN** emitted evidence says recipes exhausted without agreement
- **THEN** that exit reason is shown separately from family agreement and from departure before completion

#### Scenario: Another participant triggers shared convergence
- **WHEN** another participant changes shared state and this visit observes the completion screen
- **THEN** the current visit can satisfy its screen-observation goal without identifying or joining the other participant, attributing causality to this visit or counting another household agreement

#### Scenario: Initial convergence and repeated hydration
- **WHEN** a zero-vote visit emits a mapped completion-screen entry or a supported goal occurrence is replayed during hydration
- **THEN** the entry follows first-event/goal rules and identified replay cannot multiply completions; without reliable occurrence identity, uncertainty is shown rather than inventing distinct entries

#### Scenario: Visit left before convergence
- **WHEN** a visit has no completion-screen observation before leaving and shared convergence happens later
- **THEN** the later shared state does not fabricate completion-screen evidence for the departed visit

### Requirement: Evidence authority and semantic counting
The system SHALL evaluate each claim under its configured owning boundary: client evidence can establish a rendered screen, while committed operations require trusted persistence evidence. Transport records, action intents, committed actions and completed journeys SHALL have distinct counting units. Correlation SHALL NOT automatically deduplicate records or establish authority.

#### Scenario: Intent and commitment describe one vote
- **WHEN** one browser intent and one trusted server commitment describe the same supported action
- **THEN** a vote-commitment count does not count both records as two votes; unsupported identity remains visible as uncertainty

### Requirement: Stable semantics across versions
The system SHALL retain stable step identity for unchanged semantics and preserve application version separately from journey-definition version. Changed semantics SHALL require a distinct step identity/name rather than silently joining unlike steps. Missing release metadata SHALL remain unknown. Definition edits SHALL preserve the interpretation version used by comparisons.

#### Scenario: Same behavior across app releases
- **WHEN** two app versions emit events mapped to the same step semantics
- **THEN** their normalized journey comparison can use the same stable step

#### Scenario: First vote versus family agreement
- **WHEN** a milestone changes from first persisted vote to family agreement
- **THEN** the definition uses distinct semantics and does not compare the milestones as identical completion

### Requirement: Stage occurrence timing
The system SHALL support stage durations from configured entry/exit or start/end evidence emitted through standard OTel instrumentation. It SHALL pair occurrences within app/logical-correlation/source scope, with tab context where applicable, without double counting. Measured elapsed, observed-active and next-event-estimated time SHALL remain distinct; server elapsed time SHALL NOT imply human attention or cognitive effort.

#### Scenario: Recipe viewing and repeated votes
- **WHEN** a visit enters/exits recipe viewing and then starts/ends several occurrences of the same voting stage
- **THEN** each qualified occurrence has its own duration and the healthy loop retains its repetitions; duplicate delivery does not add a second duration where identity permits deduplication

#### Scenario: Background, idle and overlapping stages
- **WHEN** visibility, idle or tab-switch evidence exists, or stages overlap
- **THEN** the declared policy distinguishes elapsed from observed-active time, preserves tab/occurrence scope and does not add overlapping durations as a unique visit total; without activity evidence it does not label elapsed time active

#### Scenario: Missing end or delayed evidence
- **WHEN** an end is missing, arrives before its start, or has inconsistent clock/overlap context
- **THEN** the occurrence remains ongoing/right-censored or timing-unavailable until reconciled; it is not zero duration, a negative measurement or a completed observation in a duration percentile

#### Scenario: Next-event fallback
- **WHEN** the definition permits the next observed event as a stage-end fallback
- **THEN** that duration is explicitly estimated under a named adjacency/window policy and is not silently pooled with measured entry/exit durations

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
