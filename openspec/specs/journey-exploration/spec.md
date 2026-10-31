# journey-exploration Specification

## Purpose

Let viewers explain historical observed paths and release differences using ordered evidence, honest denominators and accessible alternatives to the graph. This is a proposed contract.

## Requirements

### Requirement: Historical application map
The system SHALL filter a weighted journey map by registered app, application version and explicit date range. Edges SHALL show transition counts with a named counting unit. Journey, visit and occurrence counts SHALL remain distinguishable. Filters and definition versions SHALL accompany results. Empty and partial coverage SHALL differ from zero traffic.

#### Scenario: Repeated transition
- **WHEN** one journey traverses Vote to Vote nine times
- **THEN** transition occurrences can be nine while affected-journey and affected-visit counts remain one

#### Scenario: Correlated navigation with no completion goal
- **WHEN** correlated browser logs produce an unconfigured navigation path
- **THEN** the map shows its stage labels and observed connections, identifies the path as goal-free navigation and excludes it from goal conversion
- **AND** event types already represented by path steps are not duplicated as disconnected discovery nodes in the same selection
- **AND** references into, from or between navigation stages distinguish repeated event names by the selected stage labels and preserve recorded boundaries

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

### Requirement: Discoverable completion goal configuration
The configuration interface SHALL provide a dedicated completion-goal section showing the current goal or an explicit unset state. It SHALL distinguish successful completion from an exit or dead end, and explain that a goal is optional and separate from event Kind. Owners SHALL be able to select received unmapped events as well as existing mappings, review live matches and explicitly activate the goal. Existing goal state SHALL reflect its exact selector.

#### Scenario: First goal before any mappings
- **WHEN** an owner opens Mappings & goal with received events and no mappings
- **THEN** the completion-goal section offers those events, explains selection, live checking and activation, and does not require a prior mapping

#### Scenario: Review a completion goal
- **WHEN** an owner selects an event as the completion goal
- **THEN** its editor clearly identifies successful completion, explains replacement of any current goal for new journeys, and labels the activation action Activate completion goal
- **AND** selecting exit as its Kind blocks activation with guidance to use Step for a successful milestone

#### Scenario: Inspect or replace the current goal
- **WHEN** an owner edits the active goal or another mapping sharing its event name
- **THEN** only the exact goal selector is marked as the current goal, and the interface explains how to replace it by selecting another event

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

### Requirement: Automatic mapping validation beside received evidence
The event editor SHALL automatically validate the complete draft definition when a mapping, condition, kind or goal selection changes. Received-record results SHALL serve as the preview without a separate Preview matches action or duplicate bottom summary. Activation SHALL require successful validation of the exact current draft. Pending, failed and stale checks SHALL NOT authorize activation; retry SHALL be available after a check fails. Activation SHALL remain explicit.

#### Scenario: Change a received-data selection
- **WHEN** an owner selects a received field value or changes a draft field
- **THEN** matching examples and complete-definition validation update automatically, and activation remains disabled until the latest draft passes
- **AND** an older response cannot overwrite the current selection or enable its activation

#### Scenario: Journey context in matching records
- **WHEN** matching records are displayed
- **THEN** the same section explains whether session IDs and applicable browser tab IDs are available, using plain wording without a zero-missing-context warning
- **AND** missing context does not block saving a valid mapping

#### Scenario: Failed validation request
- **WHEN** the automatic check fails
- **THEN** the editor explains the failure, offers retry and keeps activation disabled without requiring the owner to find a separate preview section

### Requirement: Chronological journey timeline
The workspace SHALL offer a timeline for a selected retained journey showing ordered stage occurrences and linked diagnostic evidence. Repeated occurrences SHALL remain individually inspectable, including when consecutive repeats are reversibly grouped. It SHALL identify the selected journey, captured boundaries, partial portions and ordering uncertainty, and provide access to ordered evidence.

#### Scenario: Returning to a page
- **WHEN** a retained journey follows Home, Recipe, Planner, Recipe, Planner
- **THEN** the timeline shows all five occurrences in recorded order and permits inspection of each occurrence

#### Scenario: No selected or retained journey
- **WHEN** no journey is selected or its detail has expired
- **THEN** the view offers a way to select an available journey or explains that detail is unavailable, without reconstructing it from summaries

#### Scenario: Compact consecutive repeats
- **WHEN** a viewer groups consecutive Vote occurrences
- **THEN** the group names the stage and occurrence count and expands to every retained occurrence in order
- **AND** grouping does not combine nonconsecutive returns, cross gaps or captured boundaries, or change counts and evidence

### Requirement: Direct transition matrix
The workspace SHALL offer a matrix of observed direct transitions with origins as rows and destinations as columns. Cells SHALL identify transition occurrences, affected journeys and visits, and provide access to available supporting evidence. Empty cells SHALL be distinguished from unavailable or partial coverage.

#### Scenario: Self transition and return
- **WHEN** journeys contain Vote to Vote and Recipe to Planner to Recipe
- **THEN** the diagonal and directional cells show their respective observed transition counts without combining directions

#### Scenario: Missing intermediate evidence
- **WHEN** an intermediate step is hidden or unavailable
- **THEN** the matrix does not count an invented direct transition across it

### Requirement: Repeat and return analysis
The workspace SHALL offer ranked stages and directed stage pairs by repeated activity, with explicit occurrence and affected-journey counts and a named denominator for any rate. A repeat SHALL mean an additional occurrence of a stage within a captured journey; a return SHALL mean a direct transition to a different stage previously observed in that same captured journey. Supporting retained journeys SHALL be inspectable.

#### Scenario: Repeated and returned stages
- **WHEN** one journey follows A, B, A, B
- **THEN** A and B each have one additional occurrence, B to A and A to B each have one return transition, and the journey is counted once per applicable row
- **AND** no repeat or return is labeled a defect solely because it occurred

#### Scenario: Partial beginning
- **WHEN** a journey beginning is unavailable
- **THEN** return classification is limited to the observed sequence and the view discloses that earlier visits may be missing

### Requirement: Ordered sequence flow
The workspace SHALL offer a bounded, expandable ordered-prefix flow and an equivalent ranked path list from retained journeys. Ordinal occurrences SHALL remain distinct beneath reversible consecutive-repeat grouping. Flow widths SHALL count journeys sharing each prefix, not transition occurrences. Limits, omitted paths, partial boundaries and supporting retained journeys SHALL remain accessible.

#### Scenario: Different ordering with identical edges
- **WHEN** journeys share aggregate edges but follow different ordered sequences
- **THEN** the flow represents their observed prefixes separately and does not infer paths from edge counts

#### Scenario: Long or uncommon paths
- **WHEN** paths exceed the displayed step or branch limit
- **THEN** omitted continuation or grouped paths are labeled with journey counts and do not appear to terminate at that display limit

#### Scenario: Explore shared beginnings
- **WHEN** several journeys share Home then Recipe before branching to Planner or Search
- **THEN** the flow shares that prefix and exposes each supported continuation on expansion
- **AND** the list separately identifies group membership counts and shared-prefix counts without inviting their addition as disjoint populations

#### Scenario: Bounded expansion
- **WHEN** the viewer requests more steps or less common branches
- **THEN** the view reveals retained evidence within declared work limits, keeps omitted counts and further evidence actions available, and preserves the selected prefix
- **AND** a gap or exhausted detail never becomes an invented continuation or a proven journey ending

### Requirement: Shared context for loop-aware views
The new views SHALL respect the existing app, version, definition, time, scope and reference projection filters and captured journey boundaries. Switching views SHALL preserve compatible stage or transition selection and filters. Baseline and live observations SHALL remain separately counted, and limited or expired detail SHALL be disclosed.

#### Scenario: Switch from an existing graph
- **WHEN** a viewer selects a transition in Force and opens the matrix
- **THEN** the same directed transition and filters remain selected and its counts are available

#### Scenario: Existing reference and layouts
- **WHEN** an existing saved reference is opened or an existing graph is selected
- **THEN** Hierarchy, Force and Sankey remain available with their established behavior and saved layout meaning
- **AND** newly saved view choices reopen with their corresponding view and safe filters

#### Scenario: Truncated or summary-only selection
- **WHEN** the source query is truncated or ordered detail is unavailable for part of the selected range
- **THEN** each affected view discloses partial coverage and does not present unsupported sequence or return counts as complete

### Requirement: Accessible loop-aware exploration
The new views SHALL provide keyboard-operable controls and structured text or tables with equivalent counts, sequence evidence and drilldowns. English and French labels SHALL have task parity, and status SHALL NOT rely on color alone.

#### Scenario: Matrix without color or pointer
- **WHEN** a viewer navigates the matrix using a keyboard or assistive technology
- **THEN** origin, destination, counts and evidence actions are available as labeled text

### Requirement: Question-led journey orientation
Each journey view SHALL state the question it answers, its primary counting unit and a visible next action. A task-oriented chooser SHALL distinguish aggregate exploration from individual-journey investigation while retaining all seven existing view identifiers and a direct way to switch views.

#### Scenario: First entry to transitions
- **WHEN** a developer opens the transition view without previous chart knowledge
- **THEN** the interface explains that it answers where journeys go directly after a stage and offers stage selection and supporting examples

#### Scenario: View names and existing references
- **WHEN** a saved matrix, sequence or graph reference is opened
- **THEN** its established view identifier resolves to the corresponding task and the developer can switch to every other supported view
- **AND** explanatory names do not invalidate existing references

### Requirement: Visible analysis population
Each analysis SHALL prominently identify app, current time bounds, scope, version, definition, retained-journey count and relevant coverage. A one-journey selection SHALL be recognizable beside the analysis and offer a direct timeline action. Clearing an intentional narrowing SHALL preserve other filters and access limits.

#### Scenario: One-journey aggregate
- **WHEN** a selected visit leaves one retained journey in a heatmap or path analysis
- **THEN** the view states that population, offers its timeline, and offers clearing the visit filter when authorized
- **AND** it does not imply that this example describes all app traffic

#### Scenario: No retained detail
- **WHEN** records exist but no retained ordered journeys support the selection
- **THEN** the view distinguishes unassigned records, expired detail, restrictive filters and a failed request where known, with a relevant recovery action
- **AND** it does not turn retained summaries into ordered examples

### Requirement: One primary transition metric
The matrix SHALL display one selected primary metric, defaulting to transition occurrences, while keeping affected-journey and visit counts available in details. Available metrics SHALL include those three counts and share of observed outgoing transition occurrences. Units and denominators SHALL accompany values; journeys and visits SHALL NOT be described as people.

#### Scenario: Repeated transition explanation
- **WHEN** Vote to Vote occurs nine times in one retained journey and one visit
- **THEN** the default cell displays nine with an occurrence unit and its details explain nine occurrences within one journey and one visit
- **AND** switching the primary metric changes its label and scale without changing the selected transition

#### Scenario: Outgoing occurrence share
- **WHEN** A has six observed direct transitions to B and four to C
- **THEN** outgoing occurrence share shows 60% for B and 40% for C, naming ten observed outgoing transition occurrences as the denominator
- **AND** omitted display columns remain in that denominator, partial source coverage is disclosed, and an empty denominator is unavailable rather than zero percent
- **AND** the percentage is not labeled conversion, people, or the probability a journey reaches B

### Requirement: Legible transition heatmap
Matrix intensity SHALL encode only the selected metric through a labeled sequential scale with numeric equivalents. Zero, unavailable, source-partial and display-omitted states SHALL be distinguishable. Selected rows, columns and cells SHALL remain locatable without relying on hue, and selection SHALL NOT change the metric's visual meaning.

#### Scenario: Find a dominant destination
- **WHEN** a developer selects an origin with several observed destinations
- **THEN** the heatmap and equivalent ranked destination list expose the same metric and supporting transitions, with pinned or otherwise persistently available origin and destination labels

#### Scenario: Known count under partial coverage
- **WHEN** a transition has observed occurrences but the source query is partial
- **THEN** its count remains readable as observed retained evidence with partial coverage, while unsupported absent cells remain distinct from known zero
- **AND** the legend identifies the metric, endpoints and any capped intensity

### Requirement: Searchable bounded stage discovery
The matrix SHALL offer stage search, explicit stage inclusion and traffic-based or alphabetical ordering. Display limits SHALL disclose visible and omitted stage counts and preserve a way to find omitted stages. Hiding a stage SHALL remain a display operation without manufacturing adjacency or silently changing metric denominators.

#### Scenario: Important stage outside initial display
- **WHEN** a frequently used or searched stage falls outside the initial displayed set
- **THEN** the developer can bring it into view without changing captured journeys and can identify what was omitted

#### Scenario: Stable matrix during inspection
- **WHEN** traffic changes while a cell or header is focused
- **THEN** row and column positions remain stable until an explicit reorder or filter action, and any updated rank is available without moving the focused target

### Requirement: Repeat prevalence and intensity
Repeat analysis SHALL distinguish extra occurrences from affected journeys. Repeat prevalence SHALL use journeys with at least two observed occurrences of a stage divided by retained journeys reaching that stage. Intensity SHALL name its selected population. Return counts SHALL retain their existing observed-only definition and SHALL NOT imply defect, failure or user intent.

#### Scenario: Compare a busy and a less common stage
- **WHEN** A repeats in 20 of 100 journeys reaching A and B repeats in eight of ten journeys reaching B
- **THEN** the view shows 20/100 and 8/10, offers prevalence and extra-occurrence sorting, and does not rank A as more problematic solely by raw count
- **AND** partial populations are described as observed-only and low denominators remain visible

#### Scenario: Repetition intensity
- **WHEN** one journey contains ten Vote occurrences
- **THEN** its row distinguishes ten total occurrences, nine extra occurrences and one affected journey, with any average naming whether it uses all reaching or only repeating journeys

#### Scenario: No returns observed
- **WHEN** available evidence contains repeats but no qualifying non-self return
- **THEN** the return section explains that no qualifying returns were observed and describes the definition instead of showing unexplained empty headers
- **AND** partial coverage prevents a claim that returns never occurred

### Requirement: Prefix and path counting clarity
Path exploration SHALL distinguish the journeys sharing a prefix from membership of a displayed path group and from the full retained selection. Counts, shares, grouping and terminal labels SHALL expose their denominators and coverage. Active, censored, display-limited, gap-stopped and explicitly ended paths SHALL remain distinguishable.

#### Scenario: Shared prefix in multiple groups
- **WHEN** two displayed groups share their first three stages
- **THEN** the shared prefix's journey count is labeled as shared support and each group's distinct membership count remains available
- **AND** the interface does not sum repeated shared support into a larger journey population

#### Scenario: Consecutive grouping in a flow
- **WHEN** Vote times ten is displayed as one expandable group
- **THEN** its label distinguishes ten stage occurrences from the number of journeys supporting that sequence
- **AND** expanding restores original ordinal positions and does not join across missing evidence

### Requirement: Readable individual journey story
The timeline SHALL provide a compact vertical reading mode with outcome or navigation-only context, observed span, occurrence count and provenance. Relative elapsed time SHALL aid scanning while exact timestamps and clock uncertainty remain available. Linked diagnostics SHALL appear at supported positions without inventing order or causality.

#### Scenario: Scan a long journey
- **WHEN** a developer opens a journey with many repeated stages
- **THEN** the stage story and supported diagnostic markers can be read without horizontal scrolling, and exact evidence remains expandable
- **AND** any elapsed value names its basis rather than claiming attention or time spent on a page

#### Scenario: Diagnostic cannot be placed precisely
- **WHEN** a linked diagnostic has uncertain ordering or no exact stage attachment
- **THEN** it remains inspectable with that uncertainty instead of being inserted at a fabricated precise position
- **AND** a later successful outcome does not erase the earlier error

### Requirement: Readable supporting journey selection
Aggregate drilldowns SHALL offer a bounded, searchable list of matching retained journeys with readable time, path summary, status, partial markers and population identity. IDs SHALL remain secondary and copyable. The matching rule, deterministic ordering and any omitted count SHALL be visible; opening an example SHALL NOT silently change the selected population.

#### Scenario: Choose an example of a repeat
- **WHEN** a developer activates inspect examples for a repeated stage
- **THEN** the list contains journeys satisfying that repeat criterion, offers an explicit default ordering and available time/status filters, and opens the chosen timeline
- **AND** no arbitrary example is presented as representative of every affected journey

#### Scenario: Many supporting journeys
- **WHEN** an aggregate has more matching journeys than the visible list limit
- **THEN** the interface provides bounded paging or further narrowing, discloses omitted matches, and does not render an unbounded row of opaque ID links

### Requirement: Reversible investigation navigation
Moving from a graph, heatmap, repeat or path group into examples, a timeline or linked diagnostics SHALL preserve the originating filters, population, selected criterion and return position. Back navigation SHALL restore compatible sort, expansion and scroll state without resurrecting expired evidence or moving keyboard focus unpredictably.

#### Scenario: Pattern to diagnosis and back
- **WHEN** a developer selects a baseline transition, opens a matching journey, inspects a diagnostic and returns
- **THEN** every step identifies baseline provenance and Back restores that transition and its prior investigation context

#### Scenario: Target expires while open
- **WHEN** selected detail expires or is deleted during an investigation
- **THEN** the interface explains the unavailable target and permits returning to authorized available context without stale raw-detail restoration

### Requirement: Relevant controls and local recovery
Controls SHALL communicate which active view and population they affect. Unsupported highlight or comparison controls SHALL be hidden or explain their unavailability. Loading, failed, empty and stale results SHALL differ; retry SHALL preserve valid context and unrelated successful population data.

#### Scenario: Unsupported graph highlight
- **WHEN** a developer switches to a view without a meaningful implementation of the selected graph lens
- **THEN** the control is unavailable with an explanation or absent, and no active-looking control falsely promises that highlighting is applied

#### Scenario: One comparison population fails
- **WHEN** refreshing baseline fails while live evidence remains available
- **THEN** baseline displays a labeled failure or stale state with retry, live remains separately identified, and failed baseline reads do not become zero traffic

### Requirement: Developer investigation handoff
The workspace SHALL offer explicit local copying of a concise investigation summary containing safe filters, population bounds, selected observation, units, denominators, coverage and available evidence references. Retained linked traces SHALL be reachable where supported. Copying SHALL NOT transmit data, include credentials or silently embed raw payloads, and unavailable references SHALL remain explicit.

#### Scenario: Prepare a reproducible bug investigation
- **WHEN** a developer copies the summary for a selected transition and opens its retained diagnostic trace
- **THEN** the summary names the exact analysis context and evidence limits, the trace keeps its structural interpretation, and returning restores the journey investigation
- **AND** clipboard success or failure is visible without creating a public link or sending a message

#### Scenario: Copy after detail expiration
- **WHEN** an earlier evidence target is no longer retained
- **THEN** the summary describes its unavailability and current safe context rather than restoring or copying expired detail

### Requirement: Evidence-backed observation summaries
Views SHALL offer concise deterministic summaries of supported observations with units, population and example actions. Summaries SHALL disclose partial or low-support evidence and SHALL NOT infer causality, abnormality, abandonment, defects or attention from traffic patterns alone.

#### Scenario: Describe a repeated stage
- **WHEN** retained evidence supports nine extra Vote occurrences in one journey
- **THEN** the summary states that observation and offers its supporting example without labeling the repeat a bug

#### Scenario: Navigation without a goal
- **WHEN** captured journeys are goal-free navigation
- **THEN** summaries describe observed paths and explicit outcomes without manufacturing conversion or failed-goal rates

### Requirement: Safe reusable investigation presentation
Saved references SHALL preserve supported metric, stage-display, sort and repeat-grouping preferences alongside existing safe context. Older references SHALL use documented defaults. Preference validation and rollback fallback SHALL preserve references and evidence; raw examples and transient live-update state SHALL NOT be persisted as reference content.

#### Scenario: Reopen a heatmap investigation
- **WHEN** an owner saves a matrix with affected journeys as its primary metric and a selected stage set
- **THEN** reopening restores those safe preferences under current authorization, retention and query limits
- **AND** an older reference with no preferences remains usable and an invalid preference is rejected or safely defaulted with a visible explanation

#### Scenario: Read-only investigation
- **WHEN** a temporary viewer changes metrics, stage display, grouping or investigation focus
- **THEN** read-only exploration is available within granted scope but saving or changing definitions still requires owner authority

### Requirement: Equivalent accessible exploration
Heatmaps, trees, timelines and comparison controls SHALL provide keyboard-operable and structured-text equivalents with English/French task parity. Metric, population, selection, error and unavailable meanings SHALL survive monochrome and reduced motion. Dense views SHALL support compact screens and zoom without losing labels or primary actions.

#### Scenario: Heatmap without a pointer or color
- **WHEN** a developer uses keyboard navigation or the ranked transition list
- **THEN** origin, destination, selected metric, denominator, population and example actions are available without reading heatmap intensity
- **AND** focus remains visible, keyboard movement has a documented bounded model, and unavailable cells do not create thousands of dead-end tab stops

#### Scenario: Zoom and localization
- **WHEN** the interface is viewed at 200% zoom or a 390 CSS-pixel viewport in French
- **THEN** primary controls, context and evidence actions remain reachable, spatial charts may scroll within labeled regions, and the vertical timeline and list alternatives remain readable

#### Scenario: Return from a dialog
- **WHEN** an example or diagnostic dialog closes
- **THEN** focus returns to its invoking control or a clearly identified successor if that evidence is no longer available

### Requirement: Bounded responsive investigation
Visualizations and evidence lists SHALL bound rendering and query work while exposing limits, omissions and recovery. Expansion, search and live refresh SHALL preserve responsive navigation and discard superseded results. Display limits SHALL NOT be represented as complete source coverage or journey endings.

#### Scenario: Dense stage population
- **WHEN** an authorized selection exceeds the declared stage, branch or evidence-list limit
- **THEN** the interface remains navigable, reveals each applied limit and provides bounded search, expansion or narrowing without creating an unbounded table

#### Scenario: Rapid investigation changes
- **WHEN** a slow result arrives after a newer app, filter or population selection
- **THEN** it cannot overwrite that selection or open a journey from the previous context

### Requirement: Fixed historical evidence-window references
Saved live investigations SHALL retain validated fixed historical evidence-window eligibility, including continuations begun before the window. Older references SHALL default to captured-start cohorts. Reopening SHALL not resume live collection or silently substitute a different eligibility policy. Unsupported policies SHALL fail closed.

#### Scenario: Save a live continuation
- **WHEN** an observed journey began before observation and contributes retained steps inside the saved bounds
- **THEN** reopening its historical reference includes those clipped steps with original boundaries and partial markers, without including out-of-window evidence

### Requirement: Independent projection coverage
Source-record, population, graph and evidence-list limits SHALL be independently disclosed. Presentation caps SHALL not change authoritative transition denominators or imply source loss.

#### Scenario: Evidence page is capped
- **WHEN** the evidence list is capped but the selected ordered journeys are fully projected
- **THEN** ordered metrics retain their supported denominators while the evidence list alone discloses its omitted records
