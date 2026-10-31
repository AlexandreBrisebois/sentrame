# Spec Delta

## ADDED Requirements

### Requirement: Chronological journey timeline
The workspace SHALL offer a timeline for a selected retained journey showing its ordered stage occurrences and linked diagnostic evidence. Repeated stages SHALL appear separately. It SHALL identify the selected journey, captured boundaries, partial portions and ordering uncertainty, and provide access to ordered evidence.

#### Scenario: Returning to a page
- **WHEN** a retained journey follows Home, Recipe, Planner, Recipe, Planner
- **THEN** the timeline shows all five occurrences in recorded order and permits inspection of each occurrence

#### Scenario: No selected or retained journey
- **WHEN** no journey is selected or its detail has expired
- **THEN** the view offers a way to select an available journey or explains that detail is unavailable, without reconstructing it from summaries

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
The workspace SHALL offer a bounded sequence flow built from retained ordered journey paths, with a stage represented separately at each ordinal position. Widths SHALL count journeys following each displayed prefix, not pairwise transition occurrences. Limits, omitted paths and partial boundaries SHALL be visible, and displayed paths SHALL link to supporting retained journeys.

#### Scenario: Different ordering with identical edges
- **WHEN** journeys share aggregate edges but follow different ordered sequences
- **THEN** the flow represents their observed prefixes separately and does not infer paths from edge counts

#### Scenario: Long or uncommon paths
- **WHEN** paths exceed the displayed step or branch limit
- **THEN** omitted continuation or grouped paths are labeled with journey counts and do not appear to terminate at that display limit

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
