# Spec Delta

## MODIFIED Requirements

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

## ADDED Requirements

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
