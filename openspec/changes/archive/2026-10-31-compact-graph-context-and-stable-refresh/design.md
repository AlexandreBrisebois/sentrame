# Design

## Context

See proposal.md for scope. Review baseline: clean `main`, HEAD `44e4f1430ec3f923a5e5934a1dda0b624cf80bad`, 2026-10-31. No applicable AGENTS.md or designer skill found in the project/ancestor chain; `.agents/skills/openspec-propose/SKILL.md` governs planning. Preserve the forest-green/paper visual language and visible focus in `docs/design/visual-theme/README.md`. Canonical journey-exploration and live-journeys contracts retain separate populations, visible material caveats and unchanged denominators.

Source evidence:

- `web/src/Workspace.svelte:334`: the task dropdown maps connections to lastGraph and four other options to matrix/repeats/sequence/journey-timeline. The adjacent seven buttons already reach every destination; the dropdown adds question wording and a shortcut to the last graph, not unique evidence actions.
- `Workspace.svelte:365`: three orientation paragraphs plus comparison text; `CoverageDisclosure.svelte:8` renders all coverage details continuously. Population identity/bounds are separately displayed at `Workspace.svelte:355` and must remain visible.
- `Workspace.svelte:212`, `:359`, `:302`: every five-second poll sets currentReadLoading and inserts a loading paragraph above the canvas. `:235` retains the previous historical snapshot for validation-only reads without lost evidence. `:354` explicitly clears snapshot before the updates-available Refresh; the graph remains mounted but receives an empty model (`:168`), so that path can empty/re-layout the graph. It was not exercised in this review.
- `Graph.svelte:32`, `:45`, `:111`: geometry derives from input; camera/placements are local state, initial fit happens once, and userCamera preserves subsequent pan/zoom. No periodic camera reset is apparent.
- `InvestigationView.svelte:58`, `:79` and `style.css:106-108`: checkbox labels share flex sizing with selects, while all inputs receive width:100%. The checkbox consumes the available label width and compresses the remaining text.

Isolated reproduction: Chrome 156.0.8078.74, current built assets `index-C0hJ1JIW.js` / `index-BjJcuUb_.css`, synthetic in-memory API on an ephemeral loopback port, separate temporary profile; no receiver writes. Three snapshot reads (initial plus two polls), constant revision, 700ms artificial read delay and a deliberately panned camera. Across 210 samples at 50ms, canvas top alternated **475.578125 → 513.828125px** during loading (**38.25px shift**, 28 loading samples), height stayed 260px, SVG identity and camera transform remained unchanged. This supports loading/layout shift for this case, not remount or camera reset. Delay makes the existing insertion observable; it does not establish real receiver latency. No changed-revision, deletion, selected-evidence/focus, real trackpad or native-zoom qualification was performed.

Sequence controls measured the checkbox at **605px** within a 605px label at 1280px, and **362px** within a 362px label at 390px; label height was 409.0625px in both cases. This independently reproduces the reported squashing in current assets. Timeline shares the rule but was not separately measured.

Existing preview19010 was read-only verified against the relinquished stage5 handoff: container `d29a9c9f2a7deae02921c88a94cb59d76c697047053f879643e101cad9aa6ee9`, image `sha256:5ff1b05e898a66b5845757d6d092f2e673ff28959ad72db795d5f8ad9dd139f8`, loopback mapping19010→8096 and matching asset names. The reproduction used neither its registry nor credentials. WFS resources/native browser session were not used.

## Goals / Non-Goals

**Goals:** graph-first default space, two compact control rows, stable refresh and accessible on-demand context. Alex approved this refinement after the measured current header remained468px desktop and842/900px phone EN/FR. `.task/compact-graph-header-review/REVIEW.md` retains measurements and trade-offs.

**Non-Goals:** change graph physics, query semantics, denominators, capture, authorization or retention; deploy to shared preview19010; waive unperformed native acceptance.

## Decisions

### Compact controls and visible context

Remove repeated app/workspace headings and the separate Journeys/Traces navigation strip; retain mode switching in the compact masthead. Fold account/settings/language actions into a named account flyout while retaining app identity/selection. The first graph control row keeps all seven view buttons, current state, Time & filters and the existing workspace menu. Move References and highlight settings into that menu. On phones the seven buttons form a single44px horizontal strip with a visible overflow affordance; keyboard focus reveals each button. Do not reintroduce the rejected task dropdown.

The second row shows short current bounds, scope, separate retained population counts, active filter/clear chips and Context, Help, Coverage and refresh/pause controls. Compact Context contains exact timestamps, version/definition details, saved-view identity and Live controls. Active narrowing remains visible; default All versions and long definition explanations do not require permanent paragraphs. Live/Baseline identity and scope differences remain visible. Units remain beside values and direct evidence actions remain available.

A fixed-size Refresh/pause control shares existing row space. Loading changes its label/state without inserting a new row; meaningful updates/stale states are announced politely, routine polling is not. Remove standalone empty44px update and18.6px population-loading reservations. Failed/partial/unauthorized states remain distinct, and material failures may intentionally take warning space.

### Anchored flyouts

Help, Coverage and detailed Context use labeled native nonmodal auto popovers, activated by real buttons with expanded/controls state. The flyout is positioned within viewport margins near its trigger, with bounded scrolling, readable type and44px Close. Enter/Space, Tab, Escape, outside dismissal, touch and return focus are verified. Popovers overlay the graph without changing canvas bounds. Opening one closes another; opening existing modal tools first dismisses flyouts. Preserve open flyouts during same-context refresh where valid.

Coverage retains every source/population/server-graph/record-list detail, exporter completeness unknown and authoritative denominator explanation. Render material warning summaries separately from its trigger so warnings remain visible while the flyout is closed. Keep local graph-cap warnings in their existing footer. Generic collection uncertainty is in the Coverage trigger description/details, never a complete-traffic assurance. Navigation-without-goal and generic comparison prose move into Context; unsupported comparison, partial/censored evidence, unknown metadata, omissions and stale/error warnings remain visibly scope-specific and may wrap.

### Stable evidence and repeat controls

Retain the existing reviewed same-context snapshot/selection/focus/scroll/camera preservation and generation/lifecycle precedence. No graph remount key or automatic Fit is introduced. Keep intrinsic18px repeat checkboxes, full wrapping labels and44px label targets in timeline/sequence; do not shrink controls or typography to meet a height target.

### Alternatives and trade-offs

Inline disclosures are simpler but consume graph height when open; Alex approved flyouts after reviewing this trade-off. Popovers cover part of the graph intentionally and require dismissal/focus testing. A phone view grid makes all destinations immediately visible but costs178px; a scroll strip retains direct actions with a swipe/focus-scroll trade-off. No destination becomes a hidden dropdown-only action.

## Risks / Trade-offs

- Context takes an extra activation → preserve concise population/bounds, active narrowing, material caveats and identifiable help buttons.
- Flyouts may run offscreen or obscure dismissal → viewport-clamped placement, bounded content and persistent Close; resize and phone checks.
- Header consolidation may break focus/return → stable active-view IDs and trace/example fallback regression tests.
- Compact live context may mix populations → explicit L/B counts, labels and scope, with detailed windows one activation away.

## Verification and migration

Default target budgets are approximately160–200px desktop and180–230px phone above the graph, excluding necessary material warnings/active narrowing. These are design targets, not evidence; measure actual1440×900 and390×844 EN/FR against468/842/900px baselines. Verify no page overflow, current assets, all destinations, account/context actions, flyout keyboard/touch/outside dismissal, stable refresh, errors and changed/removed evidence. No server/data migration. Roll back presentation without changing references or evidence. No sync/archive while native acceptance remains pending.

## Open Questions

Actual native200% zoom and physical EN/FR screen-reader acceptance remain to be observed in a coordinated native session. Prior archived tasks remain separate and user-owned. Alex’s originally viewed URL/version remains unknown; current-code measurements independently justified this approved refinement.
