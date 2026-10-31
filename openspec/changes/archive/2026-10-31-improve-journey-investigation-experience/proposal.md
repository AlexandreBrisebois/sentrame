# Proposal

## Why

The new journey views expose evidence but leave developers to discover their purpose, decode counts and choose examples without enough guidance. Similar baseline and live styling makes it difficult to tell whether an action just appeared, belongs to historical context, or changed because the view refreshed.

## What Changes

- Organize all seven journey views around developer questions, with concise explanations, population context, relevant controls and recoverable empty states.
- Turn the transition matrix into a metric-selectable heatmap with a ranked destination alternative, explicit scales and bounded stage discovery.
- Make repeats useful through prevalence, intensity, readable examples and neutral observation summaries.
- Replace repeated sequence rows with expandable ordered-prefix exploration and an equivalent ranked path list; provide optional reversible consecutive-repeat grouping.
- Make an individual journey readable through a compact vertical timeline, relative time, inline diagnostics and direct entry from aggregate evidence.
- Apply consistent baseline/live identity across graphs, heatmaps, lists, inspectors and evidence: distinct hues, explicit labels and non-color encodings. Keep error, selection, recency and traffic-volume meanings independent.
- Make live observation calm and controllable: stable positions and ordering, bounded update cues, follow controls, pause/resume, stale/reset disclosure and preserved investigation context.
- Expose both populations' windows, scopes, filters, units and coverage. Permit only supported comparisons, with no combined counts or misleading raw percentage changes between unequal populations.
- Preserve safe reference compatibility, read-only access, English/French task parity, accessible alternatives and bounded rendering/query work.
- Qualify the experience with deterministic fixtures, integrated browser checks and an observed developer task walkthrough.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `journey-exploration`: Question-led orientation, heatmap and metrics, repeat interpretation, ordered-prefix exploration, timeline, cross-view evidence navigation, references and accessible developer workflows.
- `live-journeys`: Distinct population identity, truthful comparison context, stable live updates, observation controls, freshness and provenance-preserving drilldown.

## Impact

Primarily `web/src/Workspace.svelte`, `JourneyViews.svelte`, `Graph.svelte`, `journeyViews.ts`, visualization helpers, styles, localization, reference types and relevant tests. Protected snapshot/reference projections in `internal/app/query.go` and `reference.go` may need bounded metadata or validated presentation preferences. Documentation and browser verification will change. No new telemetry, third-party service or charting dependency is required by this proposal.

This change plans an investigation experience over existing captured evidence. It does not change journey capture, goals, identity, retention, authorization, trace relationship semantics, or provide statistical causality claims. It does not add session recording, alerting, arbitrary cohort comparison, public sharing or external-route qualification.
