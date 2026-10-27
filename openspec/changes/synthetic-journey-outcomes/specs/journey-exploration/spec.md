# Journey Exploration

## Purpose

Let viewers explain historical observed paths and release differences using ordered evidence, honest denominators and accessible alternatives to the graph. This is a proposed contract.

## ADDED Requirements

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

### Requirement: Node size represents stage time
The graph SHALL encode the selected time-at-stage metric in relative node size while edge thickness represents traffic. It SHALL show a bounded scale legend, statistic and units, with numeric/table equivalents. Missing duration SHALL be unavailable rather than fast. Larger nodes SHALL NOT automatically be marked bad, and comparison scales SHALL remain explicit and compatible.

#### Scenario: Duration without cognitive inference
- **WHEN** a recipe-viewing node is larger than a voting node
- **THEN** the legend and numbers explain the measured/estimated time metric, without claiming more thinking, hesitation, failure or dissatisfaction

#### Scenario: Bounds and accessibility
- **WHEN** a duration reaches the visual cap or the graph cannot be used
- **THEN** exact values and sample/completeness remain accessible numerically in the table, clipping is disclosed and size/color is not the only way to compare

### Requirement: Saved reference journey views
The owner SHALL be able to save reusable views for goal-backward exploration, beginning/end-event segments and origin-forward exploration. Views SHALL preserve safe query filters and definition/version context without retaining raw journeys beyond lifecycle limits. They SHALL query recorded boundaries rather than change them, and SHALL expose partial evidence without fabricating direct edges.

#### Scenario: Three reusable exploration modes
- **WHEN** the owner saves a goal event, a beginning/end pair or an origin event
- **THEN** reopening the view applies its declared filters/definition context and respectively explores paths to the goal, the bounded segment or subsequent paths from the origin

#### Scenario: Saved endpoints span captured boundaries
- **WHEN** a saved segment crosses an operational goal/timeout boundary or selects another retrospective goal
- **THEN** the result preserves boundary/gap markers and does not merge captured journeys or silently change their active segmentation definition

#### Scenario: Detail expired or deleted
- **WHEN** a saved view is reopened after matching raw detail expires or is deleted
- **THEN** the definition remains reusable but cannot resurrect that detail; only supported retained aggregate queries are available, with individual/ordered drilldown disabled where unavailable
