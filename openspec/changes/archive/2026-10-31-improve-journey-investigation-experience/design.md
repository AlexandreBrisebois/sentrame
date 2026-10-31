# Design

## Context

See `proposal.md` for motivation and the two delta specs for observable contracts. This design is grounded in the current dirty workspace, including the already archived loop-aware implementation; implementation must baseline that work again rather than assume these files are committed.

Observed source behavior:

- `Workspace.svelte` owns seven layout identifiers, filters, SSE refresh, selected stage/edge, two timeline selections, evidence dialogs and saved references. A population flag accompanies direct journey opening, but the common edge inspector prefers current evidence and does not show paired population counts.
- `JourneyViews.svelte` renders a triple-count matrix, raw repeat/return tables, repeated prefix rows, and horizontally scrolling timeline cards. It takes a population label, but no graph lens. Live and baseline are rendered as separately stacked instances with the same styles.
- `journeyViews.ts` already derives direct edges, repeat/return counts and ordered prefixes from captured steps. Source gaps prevent false adjacency. Limits are eight steps and 20 groups; the matrix takes the first 100 alphabetic stage names. These are current choices, not permanent UX requirements.
- `Graph.svelte` uses similar green/gray hues for live/baseline. Error styling paints both red. Sankey already has separate volume layers with a common scale; geometry uses combined topology without summing evidence counts.
- `query.go` provides models, revision, current time, baseline bounds, coverage and bounded detail. Baseline is app-wide for the seven days ending at `liveStart`, even when live is scoped to one visit. Its bounds are fixed, but later queries can correct evidence inside them. The existing health contract advertises a five-second freshness bound. Query limits include 10,000 model records and 2,000 returned evidence records; graph projections also have their own limits.
- SSE notifications currently trigger snapshot refreshes, not counted per-event deltas. Pause inhibits refresh, but does not itself provide a precise queued-record count or a separate expiry-safe presentation model.
- `SavedView`, reference validation and existing tests support safe layout/filter restoration. New presentation fields require a validated end-to-end contract. Existing `task check` covers Go, frontend behavior/type checks, canonical specs and whitespace; it does not prove a browser workflow or human comprehension.

## Goals / Non-Goals

**Goals:** Introduce one small shared investigation model and one visual semantics system across existing views. Preserve authoritative captured evidence while making orientation, comparison and diagnosis fast. Keep presentation changes incremental and testable without a new visualization service.

**Non-Goals:** Replace capture policy, build session replay, infer user identity/attention, add generalized cohort analysis, deploy infrastructure, or claim WCAG conformance from component tests. Implementation of this change and the four contract corrections below was approved on 2026-10-30. Sync, archive, commits and production deployment remain outside this apply; disposable local Docker qualification is authorized.

## Decisions

### 1. Task-oriented entry over a flat collection of charts

Use an always-visible task chooser with familiar visualization names retained as secondary labels. The proposed mapping is:

| Developer task | Primary presentation | Existing layout identifier |
| --- | --- | --- |
| Understand connections | Hierarchy / Force / Sankey layout choice | `hierarchy`, `force`, `sankey` |
| See where journeys go next | Transitions / heatmap / ranked destinations | `matrix` |
| Find repeated activity | Repeats and returns | `repeats` |
| Explore common routes | Expandable paths / ranked path list | `sequence` |
| Inspect one journey | Vertical journey timeline | `journey-timeline` |

Each task has a short question, current population count and one meaningful initial action. An empty timeline offers readable example selection rather than picking an arbitrary journey. Preserve the established default layout for existing entries. Show only meaningful lens controls in each task; do not silently leave a graph lens looking active over unrelated data.

Alternative: tooltips on all seven tabs. Rejected as the sole solution because they still require chart knowledge and hover discovery. Explanations remain visible and brief; detailed methodology stays expandable.

### 2. A shared investigation context with explicit population identity

Introduce a typed presentation context containing selected app/generation, request generation, revision, population (`current`, `live`, `baseline`), bounds, scope, version/definition, coverage, selected criterion and view preferences. The criterion distinguishes a direct transition from a repeated stage, qualifying return, ordered prefix or exact displayed group. Evidence examples must match that criterion, not merely include one stage.

Keep transient return frames in memory: origin task, population, criterion, ordering, expansion, scroll anchor and focus target. Open example timeline, diagnostic and retained trace through those frames. A return restores compatible state; loss of access clears it. Expiry invalidates evidence, not just its visible card. Restore by stable identity, not a mutable list index. Keep graph node positions governed by their existing per-mounted-layout rule.

Alternative: independent local state in each view. Rejected because it permits an edge selection to change populations or lose its matching rule during drilldown.

### 3. Population hue, intensity and status are independent channels

Proposed palette families are **Live: saturated teal/cyan** and **Baseline: slate/blue-gray**. These are proposed semantic roles, not yet contrast-qualified hex values. Use shared theme tokens rather than per-component literals. Historical single-population results retain their normal neutral presentation; they are not automatically called baseline.

| Meaning | Encoding | Non-color equivalent |
| --- | --- | --- |
| Live population | Teal/cyan accent | `Live` text and `L` badge |
| Baseline population | Slate accent, patterned contribution | `Baseline` text and `B` badge |
| Traffic magnitude | Sequential intensity or existing volume width | Exact number, unit, scale legend |
| Selected item | Outline and selection state | Selected label and pressed/current state |
| Recent accepted live change | Brief emphasis, at most three seconds | `Updated` marker / summarized update notice |
| Error evidence | Existing local red error core or marker | `! Error evidence`, with `L`/`B` identity retained |
| Missing or partial evidence | Hatching/marker distinct from zero | Explicit partial/unavailable text |
| Unmapped event | Existing dashed node outline | Unmapped label |

For Hierarchy/Force, preserve baseline as a dashed or otherwise patterned edge contribution and live as a solid contribution. Use badges on shared nodes; do not reuse the unmapped node outline to mean baseline. For Sankey, use patterned baseline and solid live ribbons/bars without changing their area, shared occurrence scale or captured topology. On red error edges preserve a population halo/pattern and text, not two identical red paths. Hit targets and inspector actions identify the chosen population.

Provide persistent compact legend and Live / Baseline / Both presentation controls. They affect visibility, not collection or observation start. The baseline authorization toggle still controls protected data access. A population not loaded or authorized is not shown as zero.

Alternative: only lower baseline opacity. Rejected because overlap, selection dimming and errors destroy that distinction. Also reject red/green as population colors because red already signals error. Verify tokens against actual backgrounds, hover, focus, disabled states and color-vision/monochrome views; meet relevant text and non-text contrast criteria without relying on hue differences alone.

### 4. Shared comparison context and layouts

Place a compact population header next to the active analysis:

- Live: since explicit observation start; current app or selected visit; version/definition; retained support; last displayed snapshot.
- Baseline: app-wide; explicit seven-day start/end; matching applicable filters; support/coverage; correction state.

Use aligned side-by-side heatmaps, repeat columns and path support on sufficient width; stack with persistent population headers on compact screens. Timeline has a clear population selector and independently remembered journey choice. Graphs retain an overlay and population-isolation controls. The inspector always offers labeled counts and examples for each supported side; it never chooses baseline merely because a live example is missing.

Heatmaps share a numeric domain, stage order and cell coordinates in Both mode. Each population has its own hue ramp, labeled numeric legend and exact values; equivalent values occupy equivalent positions on those ramps. A small live count against a large baseline may have a pale cell, but remains readable numerically. Live-only mode may explicitly rescale. Do not use a red/green difference heatmap in this scope.

Unequal raw windows are not directly compared as change percentages. Offer percentage-point differences only for compatible app-wide populations, definitions and complete supported denominators. A live visit versus app history, unavailable denominators, source truncation, mixed incompatible definitions and zero support disable numeric difference badges with explanations. All observations remain descriptive, with sample counts; there is no significance, anomaly or causality claim. No new matched-window/cohort query is needed for this change.

### 5. Heatmap metrics and bounded stage exploration

Default primary metric is direct transition occurrences; alternatives are affected journeys, affected visits and outgoing occurrence share. An inspector gives the sentence and all supporting units. For a row A:

`share(A → B) = observed direct occurrence count(A → B) / all observed direct outgoing occurrences(A)`.

Compute the denominator before display hiding, include self-transitions, preserve gaps and separately disclose source coverage. A zero denominator yields unavailable. Do not sum affected-journey cells into a purported distinct-row population: one journey may contribute to multiple cells.

Initially show at most 12 relevant stages, with selected stage/endpoints prioritized, remaining candidates ordered by the active population's observed traffic and stable name tie-breaks. In Both mode use one aligned stage set and disclose live-only/baseline-only stages. Search can locate every stage available within the bounded source projection, replace/include a stage, or reveal ranked destinations without expanding the whole matrix. Cap the visible matrix at 32 stages per side (1,024 cells); further exploration uses explicit stage selection and the equivalent list. These are presentation limits; source omissions remain separately disclosed.

Headers stay available while scrolling; selected row/column outlines aid orientation. Use roving focus with arrow movement, Home/End and Enter for a cell's details, with concise keyboard help. The ranked list offers an ordinary table/button alternative. Never require tabbing through every zero/unknown cell.

A linear intensity domain is initialized from the visible comparison on entry and held during observation. If exceeded, keep exact counts, mark capped intensity and offer Rescale; change metric, filter or explicit Rescale recomputes the domain with a visible indication. Graph occurrence widths retain their common linear scale; a geometry update must preserve manual positions and selected targets. Holding heatmap intensity does not freeze evidence counts.

Alternative: permanent 100-by-100 matrix or a nonlinear default scale. Rejected for scan cost and harder magnitude interpretation. Source search cannot claim access to stages omitted by the server; offer narrower authoritative queries where supported.

### 6. Repeat prevalence, path exploration and timeline

Repeat prevalence uses distinct observed repeating journeys / distinct retained journeys reaching that stage. Display numerator and denominator with the percentage. Optional intensity is extra occurrences / affected repeating journeys, explicitly labeled. Partial starts can hide prior occurrences; such metrics remain observed-only. Keep returns as the existing direct transition to a different previously observed stage. Default ranking is affected-journey count, with explicit prevalence and extra-occurrence alternatives. No automatic red severity or bug labels.

Build the sequence explorer as a prefix trie keyed by ordered stage occurrences, preserving gaps and captured boundaries. Start with three ordinal steps and five child branches per expanded parent. Each expansion may reveal more up to 32 ordinal steps and 200 visible prefix nodes total; show continuation/Other counts and use retained journey drilldown beyond the work limit. Ranked path-list rows separately show group membership and selection share. Shared-prefix count is never mislabeled as exact-group membership. Preserve the distinction among explicit closure, active/censored continuation, source gap and display cap.

Consecutive-repeat grouping is off by default for compatibility and can be enabled explicitly in sequence and timeline. Group only contiguous occurrences within the same captured journey and uninterrupted evidence segment. The compact label `Vote ×10` counts occurrences, not journeys. Preserve underlying ordinal positions and expand them on demand. A/B/A is never collapsed to A/B. More detailed paths remain available even when the visual prefix ends.

Timeline defaults to a compact vertical ordered list: stage, elapsed since first available occurrence, optional delta from prior occurrence, diagnostic markers and expandable evidence. Relative timing is observed elapsed, not measured stage duration or attention. With a partial beginning use the first available occurrence as the labeled origin; never call it the captured start. Exact timestamp, occurrence/receipt basis, full ID and policy details stay available. Diagnostics with uncertain placement remain in an explicitly unplaced section.

Alternative: only prettify current rows/cards. Rejected because it leaves duplicated prefixes, ambiguous counts and horizontal reading costs intact. A trie fits actual ordered evidence; reconstructing paths from aggregate edges remains forbidden.

### 7. Useful examples, summaries and local handoff

Use a shared example picker, defaulting to newest captured start first with journey-ID tie-break; expose its ordering. Offer time/status narrowing and path-text search over the bounded matching population. Initially show 20 rows, with explicit bounded paging/narrowing. Rows summarize start, supported outcome, partial state, short path and population. IDs are secondary. No heuristic representative sample is claimed.

A deterministic summary describes only a supported observation: `Vote repeated 9 additional times in 1 journey`. It links to the matching criterion and carries low-support/partial qualifiers. Existing rare-path policy remains authoritative; a return is not automatically friction. Goal-free navigation does not acquire a completion rate.

Add an explicit local Copy investigation summary action with safe context, bounds, criterion, counts/denominators, coverage and available evidence references. Reuse sanitized evidence handling; do not include raw attributes, invitation tokens, credentials or transient response bodies. Copy is local with visible success/failure, not a public link, issue creation or transmission. Retained trace drilldown uses existing trace IDs and structural trace mode; return restores the journey context.

### 8. Streaming updates that respect inspection

Coalesce notification-driven snapshot reads, with one in-flight refresh and one pending refresh per current context; do not queue one request per message. Target at most two visual reconciliations per second. Preserve existing generation/lifetime rejection of old responses. Derive recent-change cues by stable evidence identities/revision comparison, not total-count subtraction or the number of SSE messages. Unknown deltas produce `Updates available`, not an invented record count.

Preserve matrix ordering, prefix expansion, timeline scroll anchor, selected journey and graph positions during updates. Offer an explicit Reorder action for changed rankings. Follow latest is off by default; enabling it follows within the current authorized live selection, and user navigation suspends it. Baseline corrections use a different labeled notice and never restart the teal freshness cue.

Pause freezes presentation, not collection or authorization/lifecycle invalidation. Continue mandatory access checks and permitted invalidation handling; if retention or authorization changes, remove/invalidate affected detail even while paused. Resume fetches a consistent current snapshot and clearly distinguishes replay from reset. Follow and pause are separate controls. Do not invent a rebase-baseline control or silently restart the live window on view changes.

Track transport state, in-flight snapshot state and last successful display time separately. Treat a connected socket as transport evidence only. Use the server's five-second freshness contract for current metadata, while clearly distinguishing a successful snapshot with no new records from a stale display. Summarize announcements in a polite status region; do not announce every changing cell. Reduced motion uses static update markers.

### 9. Bounded computation and server-owned truth

Reuse captured models for supported observed projections; the client does not infer missing paths or promote a loaded page to full-population evidence. Keep numeric nanosecond identities precise using existing string/BigInt conventions. Expose independent source-record, population, graph-node/edge and evidence-list coverage. Graph caps and evidence-page limits do not make ordered journey projections incomplete; each live/baseline population identifies its own bounds, scope, definition compatibility and source coverage. Exporter loss remains unknown even when retained-source projection is complete. If a projection needs metadata not in the current snapshot, add a protected bounded field/read with current authorization, revision and query semantics rather than bypassing policy or loading the database in the browser.

Memoize projections by app/generation, normalized query, population and accepted snapshot identity (including time-dependent lifecycle state), never by revision alone; render only bounded visible stages, prefixes and example pages. Do not add a charting dependency unless implementation demonstrates a concrete need; native SVG/HTML preserves current accessibility patterns and limits migration cost.

Qualification workload: at least 300 stage names, 1,000 retained journeys, 10,000 model-record limit, long/repeated paths, and a clearly labeled client stress harness producing 20 update notifications per second for 60 seconds. Separately qualify actual OTLP intake and the existing five-second reset/snapshot SSE transport; the receiver does not emit an intake notification per record or provide replay. Routine periodic resets are reconciliation signals; reconnect or unexpected context/cursor changes indicate a continuity gap until a successful snapshot. This change does not redesign that protocol. Measure on the recorded local machine/browser; target p95 local interaction response below 200 ms for already loaded data, with no task over one second attributable to visualization work. Network completion is reported separately. If limits are hit, preserve responsive explicit omissions rather than silently reducing the source. These are acceptance targets to measure, not already established capacity claims.

### 10. Reference compatibility and accessibility as part of the flow

Add a versioned, allowlisted presentation-preference object to saved references: selected metric, bounded stage IDs, sort and repeat-grouping preference. Preserve existing layout IDs and safe filters. Reject overlong lists and invalid values server-side; old references get documented defaults. On rollback, old binaries preserve stored references and use their supported default presentation; they cannot explain unknown fields. Versions aware of preferences explain unsupported versions/values and preserve the safe query. Do not persist scroll stacks, raw examples, event payloads, paused snapshots, notification cursors or a promise to resume a previous live stream. Saving live freezes safe historical bounds and a validated `evidence-window` selection policy, so a journey begun before the bounds remains included through its clipped evidence. Reopening is historical and never resumes observation. Old references default to `captured-start`; unsupported selection policies must fail closed rather than silently broaden eligibility. Old binaries may retain the stored fields but cannot reproduce new evidence-window semantics; rollback discloses this limitation in release documentation.

Check keyboard-only full workflows, roving cell focus, dialog return focus, native controls, touch targets, relevant WCAG 2.2 AA text/non-text contrast, screen-reader names, reduced motion, monochrome, 200% zoom and 390 CSS-pixel width in both languages. The timeline and table alternatives remain readable where a spatial chart requires bounded scrolling. Manual checks and assistive-technology evidence must be recorded separately from static/component checks.

## Acceptance and evidence plan

| Journey | Required observation | Evidence |
| --- | --- | --- |
| First use | Developer names a view's question, explains its primary number and finds supporting example | Observed walkthrough; target each within 30 seconds of entering the relevant task, without coaching |
| Matrix | Nine self-transitions in one journey/visit; directed reverse cells; 6/10 and 4/10 outgoing shares; shared domain; hidden columns do not change denominator | Projection tests plus running-browser action and screenshot |
| Repeats | 20/100 versus 8/10 reaching journeys; 10 occurrences versus 9 extras; empty returns; partial beginning | Projection/component checks plus readable live UI |
| Paths | Equal-edge/different-order fixtures; shared prefix; long and uncommon branches; repeat expansion; gap interruption | Exact fixture assertions plus expand/list/example browser workflow |
| Timeline | Long repeated path, partial start, exact/relative times, placed/unplaced diagnostics, unavailable detail | Keyboard workflow and screenshots in EN/FR |
| Comparison | Shared edge, live-only, baseline-only, both-error, zero/unknown side, app baseline behind one live visit | All seven views; labels, hue/pattern, paired inspector counts and population-correct examples |
| Live stability | Burst while reading; manual reorder; follow suspension; pause/resume; old response; late baseline correction; reset | Real local OTLP intake and SSE, ordered assertions and browser recording/log |
| Lifecycles | Expiry/deletion/revocation during paused and active inspection | Protected integration tests plus visible invalidation; no stale return restoration |
| Reuse | Old/new reference round trips, invalid preferences, historical live-save, rollback fallback, viewer restrictions | Backend and browser tests; no persistence of raw examples |
| Accessibility | Full pattern-to-example-to-diagnostic-to-return by keyboard; contrast/color-independent population recognition; screen reader and reduced motion | Recorded manual/assistive checks with browser/device/software identified |
| Density | Qualification workload and response targets in decision 9 | Timings, limits, hardware/browser and passed/failed/blocked/not-run report |

Do not mark the human comprehension task complete from an agent's self-evaluation or screenshot review. Implementation can finish automated and browser work while the observed developer walkthrough remains explicitly pending. Record participant count and task outcomes; one developer review is useful acceptance feedback, not a general usability-study claim.

Market references consulted on 2026-10-30: [Amplitude Journeys](https://www.amplitude.com/docs/analytics/charts/journeys/journeys-understand-visualizations), [Mixpanel Flows](https://docs.mixpanel.com/docs/reports/flows), [GA4 Path exploration](https://support.google.com/analytics/answer/9317498?hl=en), [PostHog Paths](https://posthog.com/docs/product-analytics/paths), and [Contentsquare Journey Analysis](https://support.contentsquare.com/hc/en-us/articles/37271838986641-Introduction-to-Journey-Analysis). Borrow question-led anchors, bounded expansion and evidence drilldown; do not copy their session/count semantics. Accessibility references: [W3C Use of Color](https://www.w3.org/WAI/WCAG22/Understanding/use-of-color.html) and [Pause, Stop, Hide](https://www.w3.org/WAI/WCAG22/Understanding/pause-stop-hide.html). These inform design choices; they do not qualify this implementation.

## Risks / Trade-offs

- More requirements can create an overwhelming interface → expose one primary metric/action, retain concise visible context and progressively disclose methodology.
- Baseline dominates a common scale → exact numeric labels and explicit population isolation/rescale; never silently independently normalize Both mode.
- Error/provenance/selection encodings collide → qualify combinations, especially shared red edges and monochrome; keep population text on every drilldown.
- Optional grouping hides useful retry detail → default off, show occurrence count, preserve expansion, and never group across a gap.
- Fixed comparison bounds mistaken for immutable counts → distinguish baseline corrections, source coverage and lifecycle revisions from live arrivals.
- Stable ordering hides new high-ranked items → show an updates/reorder affordance, do not move inspection targets under the user.
- Paused snapshots outlive data/access → invalidation overrides presentation freeze and transient history cannot restore removed evidence.
- Performance targets are machine dependent → record environment, source/display limits and network versus local timings; no production capacity claim.

## Migration Plan

Implement in the numbered task groups with focused regression tests; keep existing layouts usable throughout. Introduce validated optional preference fields without changing existing values. Update frontend and protected query/reference handling together where fields are needed. Run the full local gate, integrated browser cases and human walkthrough before claiming experience acceptance.

Rollback uses the prior supported layouts and default presentation preferences without deleting references or evidence. New preferences must be safe to ignore on the previous UI; verify whether the previous server tolerates unknown optional fields and document the actual fallback. No archive, canonical-spec sync, commit, deployment or baseline redefinition is part of this planning change.

## Open Questions

- Final teal/slate token values are deferred to implementation contrast and actual-background testing; the population roles and non-color contract are decided here.
- The observed walkthrough participant and review time are scheduled during qualification. Lack of availability leaves that acceptance evidence pending; it does not justify substituting automated results.
