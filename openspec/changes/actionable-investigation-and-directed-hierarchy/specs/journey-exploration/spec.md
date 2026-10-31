# Spec Delta

## MODIFIED Requirements

### Requirement: Visual occurrence-time selection
The owner SHALL select historical time ranges using retained-record histogram bars and a shared two-handle selection track. Green SHALL denote only the included interval. Keyboard controls, quick ranges and optional exact UTC inputs SHALL remain available. Histogram context SHALL be independent of graph from/to bounds, use protected bounded reads and disclose units, expiry and partial results. Known identifiers SHALL be available through a compact Find ID action alongside optional advanced exact-time controls; discovering traffic SHALL not require knowing an identifier.

#### Scenario: Slice observed traffic
- **WHEN** the owner drags a range or adjusts a keyboard handle and releases it
- **THEN** the selected UTC interval updates the same graph query while preserving other filters and captured journey boundaries

#### Scenario: Discover a visit without knowing its ID
- **WHEN** the owner selects a stage and inspects journeys by path, time and outcome
- **THEN** Follow this visit or operation scopes the graph without typing an identifier, and the X returns to all traffic


## ADDED Requirements

### Requirement: Typed exact identifier investigation
The workspace SHALL offer discoverable exact-ID lookup for anonymous session/operation scope, journey and trace IDs. Input SHALL be bounded, type-labeled and scoped to the authorized app/generation and visible query constraints. Displayed usable IDs SHALL offer Copy and the applicable Open or Filter action. Lookup SHALL never use fuzzy matching or broaden authorization.

#### Scenario: Session ID supplied by the app
- **WHEN** a viewer copies a supported anonymous session ID and chooses its ID type in Find ID
- **THEN** the graph filters by exact scope equality, displays the accepted ID type, app and remaining constraints, and offers a reversible clear action
- **AND** distinct tab/journey boundaries remain visible rather than merging everything into one person or path

#### Scenario: Journey and trace IDs
- **WHEN** the viewer opens a displayed journey ID or pastes a known retained trace ID
- **THEN** the exact authorized journey or trace opens with its evidence/interpretation context and a route back, without treating that ID as a session scope

#### Scenario: Invalid or absent match
- **WHEN** input is invalid or no retained evidence matches current constraints
- **THEN** the interface distinguishes invalid input from no authorized retained match, names active constraints and offers explicit safe recovery
- **AND** it neither enumerates other apps nor claims the session never existed, and known expiry, missing modeling and limited coverage are explained only when evidence supports them

### Requirement: Persistent scoped investigation state
Accepted ID type/value, query, criterion, selected view and evidence focus SHALL survive supported view switches, backtracking and reload within the current authorized app context. Persistence SHALL not place raw IDs in ordinary URL/query/referrer paths or retain evidence payloads. App/access changes and expiry SHALL invalidate obsolete state; saved references SHALL revalidate safe query context.

#### Scenario: Follow evidence and return
- **WHEN** a viewer narrows by scope, opens a matching occurrence, inspects its trace and returns
- **THEN** the graph/timeline restores the accepted filters, matching criterion, selection and valid focus/scroll without converting trace exploration into a new session query

#### Scenario: Reload or change app
- **WHEN** a viewer reloads a supported investigation or switches app/generation
- **THEN** reload restores only valid same-context state and app changes cannot apply or disclose a previous app's ID or evidence
- **AND** stale asynchronous results cannot overwrite the new context

### Requirement: Task-directed matching and timeline
Matching journeys and Journey timeline SHALL carry the originating question, criterion and matching occurrences into selectable ordered evidence. They SHALL offer follow-scope, occurrence/diagnostic drilldown and reversible return. Goal-free observations SHALL remain valid; configured completion SHALL identify its rule and evidence rather than imply success from the last event.

#### Scenario: Investigate a repeat
- **WHEN** a viewer selects a repeated stage or path and opens a matching example
- **THEN** the timeline opens that journey at its matched occurrence, preserves all repeat identities and offers its relevant trace/diagnostic and return to the same criterion
- **AND** newest-first examples remain labeled as examples rather than representative samples

#### Scenario: No preselected journey or configured goal
- **WHEN** timeline has no exact journey selection or matching journeys have no configured goal
- **THEN** timeline offers bounded criterion-aware selection with a clear task, and matching evidence states goal-free observation without inventing conversion or forcing goal setup

#### Scenario: Recovery and disappearing evidence
- **WHEN** an error is followed by recovery or selected detail expires
- **THEN** supported continuation remains inspectable, unsupported diagnostic causality stays unresolved, and expiry explains the unavailable evidence with a safe backtrack or retry action

### Requirement: Directed occurrence hierarchy
Hierarchy SHALL represent uninterrupted observed ordered prefixes with directed edges, separate semantic-stage and path-node identity, and retained provenance. Repeated stages SHALL remain distinct occurrences. Diverged prefixes SHALL not merge into a fabricated rejoin or imply an unobserved full path. Layers SHALL mean progression, not elapsed time or causality.

#### Scenario: Return loop and shared suffix
- **WHEN** retained paths include A to B to A to C and A to D to C
- **THEN** the returned A is a distinct path occurrence, later C nodes remain distinct by prefix and stage-wide cross-reference identifies their common stage without joining the paths

#### Scenario: Aggregate path ambiguity
- **WHEN** separate journeys supply A to B and X to B to C
- **THEN** hierarchy cannot imply that a journey followed A to B to C, and selecting an edge opens only its supporting prefix members

#### Scenario: Parallel request projection
- **WHEN** an API journey supplies uncertain or concurrent request order
- **THEN** hierarchy identifies display-order uncertainty or concurrency and does not turn an ordering choice into a causal parent/child claim

### Requirement: Bounded hierarchy evidence and endings
Hierarchy SHALL disclose branch/depth/source/display limits, populations, units and versions. Start, configured goal, explicit exit, inactivity, missing continuation and display cutoff SHALL remain distinct. Omitted branches and rare paths SHALL remain reachable through bounded expansion or exact lookup. Grouping SHALL preserve original occurrences and gaps.

#### Scenario: Deep rare branch
- **WHEN** the default branch or depth limit hides a rare matching path
- **THEN** its omitted support and cutoff are explicit, expansion or exact-ID lookup reaches retained evidence, and no cutoff becomes a journey ending or zero traffic

#### Scenario: Switch layout and access alternative
- **WHEN** a viewer switches from occurrence hierarchy to stage-aggregate Force/Sankey or uses the keyboard/list alternative
- **THEN** semantic stage, filters and criterion remain identifiable, the changed node/relationship units are explicit and exact occurrence evidence remains accessible

### Requirement: Compact actionable investigation presentation
ID access, investigation purpose, state and recovery SHALL fit the existing compact workspace/disclosure patterns with EN/FR task parity and accessible controls. Material warnings SHALL remain concise and visible; details SHALL be available on demand. Tools destinations, direct views and graph-first space SHALL be preserved.

#### Scenario: Compact and accessible investigation
- **WHEN** Find ID, matching evidence and hierarchy are used at desktop, 390 CSS pixels or 200% native zoom in either language
- **THEN** labels, 44px targets, keyboard focus, dismissal/return and structured evidence remain usable without adding permanent header rows or explanatory walls
- **AND** the normal historical fixture retains the existing 165px desktop and approximately 201px phone canvas-start budget, with intentional narrowing/material warnings accounted for separately
