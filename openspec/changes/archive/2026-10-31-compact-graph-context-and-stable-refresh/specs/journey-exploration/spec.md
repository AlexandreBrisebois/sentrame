# Journey exploration delta

## MODIFIED Requirements

### Requirement: Question-led journey orientation
Each journey view SHALL expose its question and explanatory context through an accessible on-demand flyout, with counting units beside values and a visible next evidence action. The seven direct view controls SHALL replace the redundant task dropdown, distinguish aggregate exploration from individual journeys and preserve existing view identifiers.

#### Scenario: First entry to transitions
- **WHEN** a developer opens the transition view without previous chart knowledge
- **THEN** the view offers stage selection and supporting examples, and its named help flyout explains the direct-transition question and counting units

#### Scenario: View names and existing references
- **WHEN** a saved matrix, sequence or graph reference is opened
- **THEN** its established view identifier resolves to the corresponding task and the developer can switch to every other supported view
- **AND** explanatory names do not invalidate existing references

#### Scenario: Remove duplicate navigation safely
- **WHEN** the task dropdown is absent
- **THEN** Hierarchy, Force, Sankey, Journey timeline, Transition matrix, Repeats and returns and Sequence flow remain directly keyboard/touch reachable with selected state and localized help
- **AND** expired-target and dialog/trace return fallbacks focus an available active-view control without relying on the removed dropdown

#### Scenario: Open explanatory context
- **WHEN** the viewer activates Help by keyboard, touch or pointer
- **THEN** the question, units, selection instructions and applicable population/error legend explanation become available in an anchored labeled flyout with an announced expanded state
- **AND** closing returns to compact presentation without losing evidence selection or obscuring units and population marks


### Requirement: Independent projection coverage
Source-record, population, graph and evidence-list limits SHALL be independently disclosed. Presentation caps SHALL not change authoritative transition denominators or imply source loss. A compact Coverage flyout SHALL retain all details while visible concise warnings identify material limits, omissions, partial/censored evidence and unknown coverage for each affected scope.

#### Scenario: Evidence page is capped
- **WHEN** the evidence list is capped but the selected ordered journeys are fully projected
- **THEN** ordered metrics retain their supported denominators while the evidence list alone discloses its omitted records

#### Scenario: Closed coverage with material limits
- **WHEN** Coverage is closed and a source scan, graph projection or record list is limited, or evidence is materially partial/censored or coverage is unavailable
- **THEN** concise visible text identifies the affected source, population or display and offers full details
- **AND** source loss, graph omission and record-list omission remain distinct, including local graph display caps

#### Scenario: Inspect all coverage details
- **WHEN** the viewer opens Coverage
- **THEN** source-record scan, current and applicable baseline population, graph nodes/connections, received-record list, unknown exporter completeness and denominator explanations remain available with their values and independent meanings
- **AND** opening or closing changes no counts, supported denominators, filters or evidence

#### Scenario: No actual display omissions
- **WHEN** no cap or material omission applies
- **THEN** zero-omission detail can remain collapsed, collection completeness remains explicitly unknown and no summary claims complete application traffic


### Requirement: Visible analysis population
Each analysis SHALL keep app identity, short current bounds, scope, retained-journey counts and active narrowing visible in compact context. Exact timestamps and full version/definition context SHALL remain available in a named Context flyout. Live and baseline populations SHALL remain separately labeled and scoped. A one-journey selection SHALL be recognizable beside the analysis and offer a direct timeline action. Clearing intentional narrowing SHALL preserve other filters and access limits.

#### Scenario: One-journey aggregate
- **WHEN** a selected visit leaves one retained journey in a heatmap or path analysis
- **THEN** the view states that population, offers its timeline, and offers clearing the visit filter when authorized
- **AND** it does not imply that this example describes all app traffic

#### Scenario: No retained detail
- **WHEN** records exist but no retained ordered journeys support the selection
- **THEN** the view distinguishes unassigned records, expired detail, restrictive filters and a failed request where known, with a relevant recovery action
- **AND** it does not turn retained summaries into ordered examples


## ADDED Requirements

### Requirement: Stable historical refresh presentation
Same-context historical refresh SHALL preserve valid evidence selection, focus, scroll, graph viewport and local arrangement. Routine loading status SHALL NOT move or resize the canvas or remount its graph. Refreshing, failed, stale and first-load states SHALL remain distinguishable; access loss and lifecycle invalidation SHALL override preserved state.

#### Scenario: Unchanged background validation
- **WHEN** a periodic historical validation reads the same retained evidence
- **THEN** canvas bounds, mounted graph identity, pan/zoom, arrangement, selected target, focus and scroll remain stable through loading and completion
- **AND** repetitive loading announcements do not interrupt reading

#### Scenario: Explicit refresh and recoverable failure
- **WHEN** a same-context explicit Refresh starts or fails
- **THEN** authorized last-successful evidence remains identified as refreshing or stale with local retry, without blanking the graph or clearing valid selection
- **AND** an older response cannot overwrite newer filters or access state

#### Scenario: Evidence or authorization no longer valid
- **WHEN** selected evidence expires, is deleted or access ends during refresh
- **THEN** unavailable detail is removed under existing lifecycle/access rules, remaining valid context is preserved where allowed and any lost focus target receives a clearly identified successor
- **AND** prior-query evidence is never presented as the current selection after an app/filter change

### Requirement: Readable repeat-grouping controls
Consecutive-repeat controls SHALL show a complete readable label and an intrinsic-size checkbox in timeline and sequence views. Labels SHALL wrap without compression or overlap at compact widths and zoom, retain English/French meaning and support keyboard operation and a full-label touch target. Grouping SHALL preserve original occurrence counts and evidence semantics.

#### Scenario: Compact and zoomed controls
- **WHEN** either repeat-grouping control appears at 390 CSS pixels or 200% browser zoom in English or French
- **THEN** its full label remains readable beside or wrapped around a normally sized checkbox and adjacent controls remain reachable
- **AND** activating its label toggles grouping without altering underlying evidence

### Requirement: Graph-first compact workspace header
The workspace SHALL use compact navigation/control rows and on-demand anchored Help, Coverage and Context. Duplicate headings and empty standalone status rows SHALL NOT consume graph space. All seven views SHALL remain directly reachable; phone view controls SHALL scroll horizontally with keyboard focus revealing each target. Controls SHALL retain readable labels and44px touch targets.

#### Scenario: Open and dismiss context without moving the graph
- **WHEN** a viewer opens Help, Coverage or Context by keyboard, pointer or touch
- **THEN** a labeled flyout remains within viewport bounds without resizing or moving the graph
- **AND** Close, Escape and outside dismissal are available, focus returns appropriately, and opening another flyout replaces the first

#### Scenario: Phone view navigation
- **WHEN** a viewer navigates the seven-view strip at390 CSS pixels in English or French
- **THEN** every view remains reachable by horizontal touch scrolling or keyboard focus without page horizontal overflow or a tall button grid
- **AND** the selected view and active narrowing remain visible or are revealed on focus

#### Scenario: Routine status within compact rows
- **WHEN** a same-context read starts, completes or reports available updates
- **THEN** status changes within existing row space without empty dedicated rows or canvas movement, while meaningful state changes remain announced
- **AND** material evidence/read warnings remain visible even when flyouts are closed
