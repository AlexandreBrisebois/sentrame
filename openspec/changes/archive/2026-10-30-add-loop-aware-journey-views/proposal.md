# Proposal

## Why

Journeys often revisit the same stages, making a hierarchy difficult to interpret and an aggregated network insufficient to explain sequence. Viewers need complementary ways to inspect individual paths and recurring transitions without treating repetition as failure.

## What Changes

- Add a chronological journey timeline, transition matrix, repeat-and-return analysis, and bounded sequence flow.
- Preserve the existing Hierarchy, Force and Sankey views, their behavior and saved references.
- Carry filters, selection, captured boundaries, counting units, detail availability and partial-result disclosures through the new views.
- Provide accessible, localized equivalents and drilldown to retained ordered evidence.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `journey-exploration`: Add complementary views for ordered journeys and repeated transitions.

## Impact

Journey workspace, visualization and accessible tables, snapshot-derived data, reference layout validation and persistence, localization, documentation and tests may change. Existing retention, correlation, capture and trace semantics remain authoritative. No new collection or external service is required.
