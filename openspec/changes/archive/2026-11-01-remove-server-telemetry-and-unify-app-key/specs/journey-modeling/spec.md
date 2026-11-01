# Journey modeling delta

## MODIFIED Requirements

### Requirement: Separate application visits and journeys
The system SHALL derive client journeys within one registered app and bounded explicit visit/tab scope. Tabs SHALL remain separate. Server-only and combined client/server journeys SHALL be unsupported. Person, cross-visit and cross-app joins SHALL NOT be inferred; request traces SHALL NOT automatically become journeys.

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
- **WHEN** a received record lacks the required visit or tab context
- **THEN** it remains received evidence without a fabricated journey or a server-operation fallback

#### Scenario: Combined client and server evidence
- **WHEN** retained evidence previously formed a mixed journey
- **THEN** the lifecycle cutover marks it non-comparable rather than flattening it into a new client path

### Requirement: Explicit journey lifecycle
The system SHALL start at the first mapped event or eligible unconfigured client log within explicit visit/tab scope. Configured goals SHALL close journeys; otherwise 30 minutes without eligible activity SHALL close them inactive. Goal-free journeys SHALL NOT enter completed-goal counts or goal conversion. Browser tabs SHALL remain separate. Severity and saved-view queries SHALL NOT change boundaries.

#### Scenario: Mappings without a completion goal
- **WHEN** an owner previews and activates a valid mapping definition without a goal
- **THEN** activation succeeds, mapped events form paths, and no event implicitly becomes completion
- **AND** journeys remain active until inactivity or an explicitly mapped exit; missing client context remains unassigned and late revisions remain explicit

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

### Requirement: Correlated browser navigation without configuration
With an empty captured definition, the system SHALL derive navigation from named non-diagnostic client logs sharing safe explicit session and tab IDs within one app. It SHALL label steps with retained `stage.name`, otherwise the event name. Explicit definitions SHALL retain precedence. Navigation SHALL NOT imply business meaning, goals or duration, or turn uncorrelated logs, metrics or native spans into automatic steps.

#### Scenario: Next.js navigation retains document correlation
- **WHEN** Home, Search and Planner browser logs share session and tab IDs, increasing source sequences, and an empty captured definition
- **THEN** they form one ordered path with Home-to-Search and Search-to-Planner transitions without requiring a mapping or completion goal
- **AND** correlated browser errors remain evidence attached to the path rather than automatic steps or terminal outcomes
- **AND** unique positive `source.sequence` values within one source/service determine order; otherwise occurrence time and stable receipt tie-breaks are used with order uncertainty disclosed
- **AND** exact identified replay does not multiply steps

#### Scenario: Missing correlation and concurrent visits
- **WHEN** browser logs lack session or tab context, or belong to different apps, sessions or tabs
- **THEN** missing-context records remain unassigned and different scopes remain separate without proximity joins

#### Scenario: Configured meaning overrides emitted labels
- **WHEN** a captured definition explicitly maps an event to a label or goal
- **THEN** its literal mappings and lifecycle apply, and automatic navigation does not add excluded events or replace owner labels

### Requirement: Deterministic evidence-based outcomes
The system SHALL separate completion, positive business outcome, observed failure and configured flaky behavior. Each classification SHALL expose its rule version and evidence. A negative vote SHALL NOT imply a negative journey outcome. HTTP success, optimistic UI, repetition, silence or an exception alone SHALL NOT establish goal completion or failure.

#### Scenario: Recovered operation
- **WHEN** a client request error is followed by an observed recovery and later configured goal evidence
- **THEN** the journey shows the completed goal plus the earlier friction rather than ending at the error

#### Scenario: Publication fails after commit
- **WHEN** a client observes a vote result and then a follow-up error
- **THEN** the result observation, error and later completion-screen evidence remain separate; none independently proves durable backend persistence

#### Scenario: Flaky classification
- **WHEN** a configured retry/error condition matches within its declared denominator and time window
- **THEN** the interface shows the matched condition and counts; without such a rule it does not label the journey flaky automatically

#### Scenario: Contradictory evidence
- **WHEN** terminal evidence conflicts under the active rule
- **THEN** the system shows the conflict and deterministic unresolved result rather than inventing a preferred business story

### Requirement: Healthy voting loops and explicit exit reasons
The system SHALL preserve vote-loop repetitions and valid duration. WFS completion-screen observation SHALL be a visit outcome distinct from shared agreement, its cause and emitted end reason. Mappings SHALL distinguish recipes exhausted, family agreement and departure before completion where emitted. Missing evidence SHALL NOT imply tiredness, frustration or a defect.

#### Scenario: Repeated healthy voting
- **WHEN** a journey records ten distinct qualified client vote-result observations and then observes the completion screen
- **THEN** it shows ten observed vote results, nine successive result-to-result transitions and one evidenced visit completion without a repetition-based failure label or a claim that this visit caused shared agreement

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
The system SHALL describe untrusted client observations, including rendered screens, intents and reported results. It SHALL NOT certify backend persistence or truthful user actions. Records, action observations and completed configured journeys SHALL retain distinct units; correlation alone SHALL NOT deduplicate them.

#### Scenario: Intent and commitment describe one vote
- **WHEN** one client intent and one result observation describe the same supported action
- **THEN** the interface distinguishes two records from one qualified interaction, labels unsupported identity as uncertain and makes no trusted-commitment claim

### Requirement: Stage occurrence timing
The system SHALL support stage durations from configured entry/exit or start/end evidence emitted through standard OTel instrumentation. It SHALL pair occurrences within app/visit/tab and supported emitter context, without double counting. Measured elapsed, observed-active and next-event-estimated time SHALL remain distinct; operation elapsed time SHALL NOT imply human attention or cognitive effort.

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
Owners SHALL map standard span/event names and supported attributes using explicit service, carrier and retained attribute selectors without `stage.name`. They SHALL preview stages, supporting activity, goals, exceptions and missing prerequisites before prospective activation.

#### Scenario: Ordinary voting operation
- **WHEN** a client exports a request span and named `vote.result.shown` event without Sentrame attributes
- **THEN** the owner can label the operation and map the result observation to a milestone without requiring a producer role or creating a journey from a trace alone

#### Scenario: Unproven persistence
- **WHEN** only a browser button event and generic HTTP success exist
- **THEN** configuration explains that received observations do not certify persistence; role selection cannot manufacture stronger evidence

### Requirement: Optional stage hint suggestions
The system SHALL retain optional `stage.name` as a Sentrame application convention and suggest labels for unmapped evidence. Business mapping suggestions SHALL require owner acceptance into literal versioned rules; the attribute alone SHALL NOT assign Kind, goal, session, boundary, authority or timing quality. Unconfigured correlated browser navigation may use the retained label for display without activating a business rule.

#### Scenario: Enriched voting event
- **WHEN** an unmapped client `vote.result.shown` event carries `stage.name=voting.result`
- **THEN** configuration suggests that label and shows its source, and activation requires explicit owner preview/acceptance rather than silently creating a graph stage

#### Scenario: Missing and conflicting labels
- **WHEN** the same selector emits absent, empty, invalid or several distinct stage labels
- **THEN** missing/invalid values yield no suggestion, distinct valid labels remain separate suggestions and the owner can choose a literal label without a popularity-based default activation

### Requirement: Purposeful multi-request correlation
Client journeys spanning requests SHALL require explicit bounded app/visit/tab correlation. `session.id` SHALL remain an opt-in client convention, never server-operation scope or person identity. No IP, user-agent, fingerprint, timestamp or automatic cross-session inference SHALL join evidence.

#### Scenario: Multiple request traces
- **WHEN** navigation, button press, failed request, retry and observed result span different traces without a shared supported scope
- **THEN** each trace remains explorable but no single browsing session or complete multi-request path is invented

#### Scenario: Explicit combined correlation
- **WHEN** those records carry correctly propagated random `session.id` and matching browser `tab.id`, as client observations
- **THEN** the configured journey preserves observed intent, failure, retry, recovery and configured goal under first-mapped/goal/30-minute rules without joining unrelated tabs, apps or visits
