# journey-exploration Specification

## Purpose

Let viewers explain historical observed paths and release differences using ordered evidence, honest denominators and accessible alternatives to the graph. This is a proposed contract.

## Requirements

### Requirement: Historical application map
The system SHALL filter a weighted journey map by registered app, application version and explicit date range. Edges SHALL show transition counts with a named counting unit. Journey, visit and occurrence counts SHALL remain distinguishable. Filters and definition versions SHALL accompany results. Empty and partial coverage SHALL differ from zero traffic.

#### Scenario: Repeated transition
- **WHEN** one journey traverses Vote to Vote nine times
- **THEN** transition occurrences can be nine while affected-journey and affected-visit counts remain one

#### Scenario: Retention boundary
- **WHEN** a selected range extends outside retained detail
- **THEN** the interface identifies unavailable dates and does not present incomplete counts as full-range evidence

### Requirement: Selection preserves path evidence
The system SHALL distinguish journey/cohort eligibility, evidence-window clipping and display-only step hiding. It SHALL declare start, follow-up and version-attribution policies. Hidden or missing intermediate steps SHALL NOT create observed direct transitions. Startless/left-truncated and right-censored journeys SHALL remain separate from complete-journey denominators.

#### Scenario: Hidden middle step
- **WHEN** a path A to X to B is displayed with X hidden
- **THEN** any A-to-B projection is labeled as spanning omitted steps and is not counted as an observed direct A-to-B transition

#### Scenario: Missing or expired start
- **WHEN** a goal remains but its start is outside the cohort/window or has expired
- **THEN** the goal evidence remains inspectable as partial and does not silently enter a started-journey completion rate

#### Scenario: Mixed version and unequal follow-up
- **WHEN** a journey crosses a deployment or comparison cohorts have unequal follow-up
- **THEN** the declared attribution/window policy and censoring are shown without splicing per-event version fragments into fictitious paths

### Requirement: Ordered path and loop inspection
The system SHALL retain sufficient ordered detail to inspect complete retained paths, repeated steps, loops and exceptions. A graph summary or pairwise counts SHALL NOT substitute for path evidence. Truncation, aggregation and display limits SHALL be explicit, with a drilldown to retained sequence detail.

#### Scenario: Pairwise ambiguity
- **WHEN** two evidence sets have equal edge counts but different path order or goal outcomes
- **THEN** path drilldown distinguishes the sequences and does not reconstruct an invented unique path from those edges

#### Scenario: Long journey
- **WHEN** a journey exceeds the graph display limit
- **THEN** the interface marks the omitted display portion and offers paged ordered evidence without implying that the journey ended there

#### Scenario: Discover paths to a selected goal event
- **WHEN** the owner chooses a goal event from the discovered-event list
- **THEN** matching journeys expose their known ordered prefixes from journey start through that occurrence, with stable step mappings, cohort/window eligibility and any truncated/missing beginning explicit; later events remain available for recovery inspection

### Requirement: Rare paths with traffic context
The system SHALL show rare-path counts, percentages, selected denominator and coverage. It SHALL distinguish low-support observations from sufficient traffic under a visible configured threshold. Rarity SHALL NOT automatically imply failure or a bug. Hidden paths SHALL remain discoverable through a list or drilldown.

#### Scenario: One out of three
- **WHEN** a path appears once among three comparable journeys
- **THEN** the interface shows one of three and low support instead of a confident outlier defect claim

### Requirement: Default stopping and friction evidence
The system SHALL show potential dead ends by default with stopping rates, observation window, pending/censored cases, retries, loops and known outcomes. It SHALL distinguish explicit exits from inactivity and unavailable telemetry. It SHALL NOT describe silence or missing events as proven failure or abandonment.

#### Scenario: Recently active visit
- **WHEN** the selected window ends before a journey has enough follow-up time
- **THEN** its stopping evidence is pending/censored rather than an established dead end

#### Scenario: Repeated-step stopping denominator
- **WHEN** one journey visits Vote ten times and stops there
- **THEN** a stopping rate explicitly uses eligible step occurrences or eligible journeys, and does not treat ten occurrences as ten independent journeys

#### Scenario: Error and later goal
- **WHEN** a journey resumes after an exception and reaches its goal
- **THEN** its path remains visible as recovery and completion, including the earlier retry/loop evidence

### Requirement: Normalized application version comparison
The system SHALL compare compatible step semantics across application versions using explicit comparable windows, denominators and definition versions. It SHALL show sample sizes, missing metadata and sampling/coverage differences. Absolute traffic and normalized rates SHALL be separate. A rate difference SHALL NOT be claimed to establish release causality.

#### Scenario: Unequal release traffic
- **WHEN** one version has 100 eligible journeys and another has 10
- **THEN** the comparison shows both counts and rates using compatible definitions rather than ranking raw counts as better performance

#### Scenario: Different eligibility rules
- **WHEN** versions use incompatible goal definitions or collection coverage
- **THEN** the interface prevents a silently equivalent comparison and explains the mismatch

### Requirement: Accessible self-explanatory exploration
The interface SHALL support keyboard operation, visible focus, plain-language labels, recoverable errors and an equivalent path/count table. Status SHALL use more than color. Live motion SHALL be pausable. English/French meaning and task parity SHALL be reviewed. WCAG 2.2 AA SHALL be the proposed verification target, with checks distinguished from conformance claims.

#### Scenario: Graph unavailable
- **WHEN** a viewer cannot use the canvas, pointer or animation
- **THEN** the same filters, path evidence, counts and diagnostic actions remain usable through labeled keyboard controls and a structured list/table

#### Scenario: Language and filter changes
- **WHEN** a viewer changes language or corrects an invalid date range
- **THEN** labels and recovery instructions remain understandable, focus is managed and valid filter selections are preserved

### Requirement: Transparent time-at-stage summaries
The system SHALL show median time at stage, slower cases when adequately supported, units, eligible occurrence count and completeness. The slower statistic/support policy SHALL be explicit. Measured elapsed, observed-active and estimated time SHALL remain distinct; missing/right-censored data SHALL remain visible. Comparisons SHALL use compatible stage/clock/window/activity policies; slower time SHALL NOT alone imply a defect.

#### Scenario: Completed and incomplete samples
- **WHEN** a stage has completed measured occurrences plus ongoing, missing-end or estimated occurrences
- **THEN** the summary names its selected sample and exclusions, exposes completeness and does not silently treat missing ends as zero or mix timing kinds

#### Scenario: Comparing release durations
- **WHEN** recipe-viewing or per-vote durations are compared across releases
- **THEN** the view shows sample/support, follow-up and policy differences; repeated votes are occurrences within journeys rather than independent people, and unequal censoring is not hidden

### Requirement: Quiet switchable graph workspace
The interface SHALL provide hierarchy, directed force and Sankey views preserving selection, filters and captured boundaries. Timing SHALL remain in the inspector and numeric tables. Sankey bars and ribbons SHALL use a common linear occurrence scale with baseline/live counts separate. Arrows SHALL remain small independently of weight/zoom. Loops and skipped-stage links SHALL remain explicit without implying intermediate adjacency.

#### Scenario: Switch views while inspecting a stage
- **WHEN** the owner selects a stage and switches between hierarchy, force and Sankey
- **THEN** the same selection, filters and incoming/outgoing evidence remain available, with timing basis, samples and exclusions in the inspector

#### Scenario: Heavy traffic and healthy loops
- **WHEN** a heavily traversed edge or repeated stage appears
- **THEN** its traffic remains visible without oversized arrowheads, hidden loops or invented adjacency

#### Scenario: Trackpad and focus
- **WHEN** the owner pinches or scrolls with two fingers over the canvas
- **THEN** zoom preserves its gesture focal point and scrolling pans; pointer selection does not add a whole-canvas focus box, while keyboard focus remains visible

### Requirement: Contextual event meaning and graph highlights
The owner SHALL configure discovered events through a node popup with real-record preview and prospective activation preserving other mappings. Unmapped nodes SHALL remain unconnected absent evidence. The graph SHALL offer rare-path, endings, warning/exception and recovery lenses with evidence drilldown. Rarity SHALL use comparable full closed, non-partial, uncensored sequences; warnings SHALL NOT terminate journeys.

#### Scenario: Configure one discovered event
- **WHEN** the owner opens its popup, previews and activates a mapping
- **THEN** all other mappings remain, captured journeys keep their definition and producer authority remains server-validated

#### Scenario: Rare sequence support
- **WHEN** a full sequence has at most ten percent share among at least twenty closed, non-partial, uncensored journeys with the same definition and initial version
- **THEN** the rare lens may highlight it with counts/support, without implying its links occur exclusively in rare paths or that rarity proves a defect

#### Scenario: Named ending states
- **WHEN** endings are highlighted
- **THEN** pending, inactive, explicit exit and completed outcomes remain distinguishable

#### Scenario: Clear a selected visit
- **WHEN** the owner uses the X on the selected journey/visit filter
- **THEN** the workspace returns to all traffic under its other filters and layout without deleting evidence

### Requirement: Saved reference journey views
The owner SHALL be able to save reusable views for the current workspace, goal-backward exploration, beginning/end-event segments and origin-forward exploration. Views SHALL preserve safe query filters and definition/version context without retaining raw journeys beyond lifecycle limits. They SHALL query recorded boundaries rather than change them, and SHALL expose partial evidence without fabricating direct edges.

#### Scenario: Three reusable exploration modes
- **WHEN** the owner saves a goal event, a beginning/end pair or an origin event
- **THEN** reopening the view applies its declared filters/definition context and respectively explores paths to the goal, the bounded segment or subsequent paths from the origin

#### Scenario: Saved endpoints span captured boundaries
- **WHEN** a saved segment crosses an operational goal/timeout boundary or selects another retrospective goal
- **THEN** the result preserves boundary/gap markers and does not merge captured journeys or silently change their active segmentation definition

#### Scenario: Detail expired or deleted
- **WHEN** a saved view is reopened after matching raw detail expires or is deleted
- **THEN** the definition remains reusable but cannot resurrect that detail; only supported retained aggregate queries are available, with individual/ordered drilldown disabled where unavailable

#### Scenario: Explore first, save later
- **WHEN** the owner explores into or from a selected stage, or chooses beginning/end stages
- **THEN** the same read path previews a projection over captured journeys without saving or resegmenting; saving later retains exact time bounds, visit/version/definition filters, layout and highlight

#### Scenario: Delete a saved reference
- **WHEN** the owner confirms deletion of a saved reference
- **THEN** only that saved query is removed, while collected evidence, captured journeys and other references remain
- **AND** deletion requires owner authorization and origin/CSRF proof, survives restart, and offers a cancel action before removal

### Requirement: Visual occurrence-time selection
The owner SHALL select historical time ranges using retained-record histogram bars and a shared two-handle selection track. Green SHALL denote only the included interval. Keyboard controls, quick ranges and optional exact UTC inputs SHALL remain available. Histogram context SHALL be independent of graph from/to bounds, use protected bounded reads and disclose units, expiry and partial results. Opaque correlation IDs SHALL be optional advanced search, not a primary dropdown.

#### Scenario: Slice observed traffic
- **WHEN** the owner drags a range or adjusts a keyboard handle and releases it
- **THEN** the selected UTC interval updates the same graph query while preserving other filters and captured journey boundaries

#### Scenario: Discover a visit without knowing its ID
- **WHEN** the owner selects a stage and inspects journeys by path, time and outcome
- **THEN** Follow this visit or operation scopes the graph without typing an identifier, and the X returns to all traffic

### Requirement: Local error evidence marking
The graph SHALL mark error-related observed connections red with a text or symbol equivalent and ordered-evidence access. Error classification SHALL use emitted error severity or exception evidence, preserve gaps and keep baseline/live evidence separate. It SHALL NOT mark every connection of an affected journey as failed or change capture outcomes.

#### Scenario: Error followed by recovery
- **WHEN** an observed error stage has incoming and outgoing transitions followed by healthy recovery and goal completion
- **THEN** error-related transitions are red, subsequent healthy connections remain distinct and the journey can remain completed

### Requirement: Manual graph arrangement
The owner SHALL be able to reposition nodes without changing evidence, counts, boundaries or mappings. Force SHALL allow both axes; hierarchy and Sankey SHALL preserve columns and volume encodings. Connections SHALL follow moved nodes. Clicks SHALL remain distinct from drags. Keyboard movement and resetting SHALL be available. Positions SHALL remain local presentation state per mounted layout, not persisted reference or capture configuration.

#### Scenario: Move and inspect a node
- **WHEN** the owner drags a node beyond the click threshold and releases it
- **THEN** its observed connections follow its position, the evidence remains unchanged and no inspector opens from that drag
- **AND** a subsequent ordinary click still opens the node inspector

#### Scenario: Restore automatic arrangement
- **WHEN** the owner moves nodes, switches layouts and returns, then chooses Reset layout
- **THEN** each layout initially retains its separate manual positions, including through live updates, and resetting restores only the current layout's automatic positions

#### Scenario: Keyboard arrangement
- **WHEN** the owner focuses a node and presses Shift plus an arrow key
- **THEN** it moves within the same bounds as dragging, with horizontal movement limited to Force view

### Requirement: Distinct trace and journey views
The interface SHALL name its current scope and relationship units. Trace views SHALL show operation instances and typed structural relations; journey graphs SHALL show configured chronological stage transitions within captured boundaries. Switching layouts SHALL NOT change relationship meaning or denominators.

#### Scenario: Parallel trace versus journey
- **WHEN** a trace has nested and parallel operations but no journey scope
- **THEN** tree/timeline/table exploration remains useful without Sankey funnel counts, visit counts or fabricated chronological sibling edges

#### Scenario: Journey with a technical trace
- **WHEN** a supported correlated journey contains several technical operations
- **THEN** opening one trace preserves its operation structure and returning restores the journey's filters, path boundaries and separate counting units

### Requirement: Contextual meaning without mandatory enrichment
The interface SHALL offer configuration from discovered services, spans, events and errors using retained examples. It SHALL explain missing context for each capability and distinguish optional stage hints from active owner mappings. Authorized viewers SHALL inspect but SHALL NOT activate rules.

#### Scenario: First standard export
- **WHEN** an owner discovers a standard request span without stage or session attributes
- **THEN** the owner can explore timing, preview a business label and see why multi-request paths remain unavailable without being told to add stage attributes first

### Requirement: Trace and journey reference context
Saved references SHALL retain explicit exploration mode, safe filters, definition context and supported captured endpoints. Queries and mode changes SHALL NOT rewrite boundaries or manufacture adjacency across hidden, missing, expired or separate evidence.

#### Scenario: Reopen after a definition change
- **WHEN** a reference created under definition A is reopened after B activates
- **THEN** its recorded interpretation context remains explicit and A's captured paths are not silently recalculated using B

#### Scenario: Expired linked target
- **WHEN** retained evidence points to a trace or journey endpoint that has expired
- **THEN** the reference identifies unavailable detail without resurrecting it or inventing a direct edge through the gap

### Requirement: Explicit timing basis in exploration
Inspectors and accessible tables SHALL distinguish native inclusive span duration, configured stage elapsed, estimated and observed-active time, including policy, sample count, exclusions and censoring. Incompatible bases SHALL NOT be pooled.

#### Scenario: Mixed operation and page observations
- **WHEN** the Voting stage has native request timing and separately observed browser intervals
- **THEN** each basis appears separately with its prerequisites and unknowns; no single duration claims page attention
