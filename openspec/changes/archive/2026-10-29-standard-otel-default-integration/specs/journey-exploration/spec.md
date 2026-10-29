# Journey Exploration

## ADDED Requirements

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
