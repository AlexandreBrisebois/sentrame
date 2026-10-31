# Proposal

## Why

Persistent orientation and coverage paragraphs crowd the journey graph, while historical refresh feels jumpy and the consecutive-repeat checkbox label is compressed. Make explanation available on demand and keep investigation position stable without concealing limits or changing evidence semantics.

## What Changes

- Consolidate the graph header into two compact control rows; remove duplicate workspace headings, separate mode-navigation/support rows and empty status reservations. Use a single horizontally scrollable direct-view strip on phones.

- Remove the redundant “Choose an investigation task” dropdown; retain the seven view buttons, selected state, stable destinations and view-specific help.
- Replace repeated orientation prose with an accessible “About this view” anchored info flyout; retain task identity, units beside values, population identity and primary evidence actions.
- Replace the expanded coverage block with a compact Coverage flyout and persistent concise warnings for actual source caps, omissions, partial/censored evidence, unavailable coverage and comparison restrictions.
- Stabilize same-context historical background and explicit refresh presentation, retaining valid selection, focus, scroll, graph camera and local arrangement; preserve lifecycle invalidation and access clearing.
- Give consecutive-repeat checkboxes intrinsic width and their full labels readable wrapping in timeline and sequence controls, including French, narrow screens and zoom.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `journey-exploration`: compact accessible orientation/coverage, stable historical refresh and legible repeat controls. Existing counting, denominator, scope and retention contracts remain.

## Impact

Frontend changes are limited to Workspace and its masthead, shared flyouts, CoverageDisclosure/InvestigationHelp, InvestigationView, modal focus fallbacks, shared styles, localized labels and focused UI tests. Reuse native popover semantics and existing graph camera behavior where possible. No API, ingestion, polling interval, capture, persistence or dependency changes are proposed.

Alex approved implementation and the tighter graph-first refinement on 2026-10-31. The archived `2026-10-31-improve-journey-investigation-experience` remains closed; its three user-owned manual checks remain unpassed. Production, real telemetry, shared receiver mutations and WFS resources are out of scope.
