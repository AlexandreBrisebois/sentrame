# Tasks

Implementation and the four reviewed contract corrections were approved on 2026-10-30. Read `design.md` for decisions, fixture cases and evidence boundaries. Preserve the existing dirty implementation and archived change; do not treat this plan as permission to commit, deploy, sync or archive.

## Archive disposition — 2026-10-31

Alex explicitly requested: “I'll handle the checks. sync and archive the spec for sentrame. then give me a git commit comment”. This later instruction supersedes the earlier requirement to wait for every manual acceptance check before syncing and archiving this change. It does not declare those checks passed or authorize a commit, push or deployment.

Implementation/software acceptance remains **30/33**, with these tasks deliberately unchecked and delegated to Alex:

| Task | Remaining user-owned observation | Status |
| --- | --- | --- |
| 5.5 | Actual native browser 200% zoom; the documentation, compact-width and software workflow portions already have separate evidence | Delegated to Alex; native zoom not verified passed |
| 8.4 | Physical screen-reader session in EN/FR across all seven views, plus outstanding native zoom observations | Delegated to Alex; dedicated assistive checks not verified passed |
| 8.5 | Observed uncoached developer walkthrough; record participant count, per-task time against the 30-second target, mistakes and follow-up | Delegated to Alex; no observed participant result |

The supported archive-with-incomplete-tasks path is the repository archive skill's confirmation flow and the CLI's `--yes` confirmation flag. Alex's explicit request provides that confirmation. Canonical deltas are synced and verified before `openspec archive improve-journey-investigation-experience --yes --skip-specs`; `--skip-specs` prevents applying those same deltas a second time, while normal change validation remains enabled. No checkbox is converted to complete, no requirement is removed, and no validation bypass or unsupported metadata is used. See `ARCHIVE-NOTE.md` for the closeout record. Archive-completion lint intentionally still reports these three unchecked tasks; that is an honest outstanding-check report, not a claim of 33/33 acceptance.

## 1. Shared investigation context and orientation

- [x] 1.1 Add typed population/criterion/return-context helpers and stale-context rejection around the existing workspace state; verify focused tests distinguish direct edges, repeat membership, returns, prefixes and exact groups and never cross app/generation/population boundaries.
- [x] 1.2 Introduce the task-oriented chooser, visible question/unit/next action and population context while retaining all seven layout identifiers; verify existing layout entry and references plus first-use, one-journey and no-detail cases in component/browser checks.
- [x] 1.3 Make lens availability, loading/error/empty/stale states and local retry meaningful for the active task; verify a nonfunctional highlight cannot appear applied and a failed population does not become zero or clear valid unrelated context.
- [x] 1.4 Document the task mapping, context meanings and recovery actions in the user guide and add English/French labels in this group; verify label parity and keyboard entry to each task.

## 2. Population identity and comparison

- [x] 2.1 Introduce shared Live/Baseline tokens, labels, badges and legends across all seven views; verify actual-background contrast, monochrome and color-vision checks, preserving unmapped, selection and local-error meanings including an edge with errors in both populations.
- [x] 2.2 Add Live/Baseline/Both presentation controls, aligned comparison context and independently labeled inspector/example actions; verify shared, live-only, baseline-only and unauthorized/absent baseline fixtures without changing live start or authorization.
- [x] 2.3 Implement common heatmap/volume scale metadata and comparison eligibility; verify equal values use equivalent common-domain positions, scoped or incomplete populations disable difference badges, and compatible outgoing shares use percentage points with visible denominators rather than raw-window change percentages.
- [x] 2.4 Document the population palette, patterns, fixed baseline bounds versus corrections, comparison restrictions and responsive paired/stacked layout; verify EN/FR help and a running-browser example with one live visit against app-wide history.

## 3. Transition heatmap and destination exploration

- [x] 3.1 Add selected transition metric and precise outgoing-occurrence denominators using captured direct evidence; extend `journeyViews.test.ts` for self/reverse edges, gaps, nine/one/one, six/ten and four/ten shares, display-hidden stages, zero denominators and partial source coverage.
- [x] 3.2 Build the labeled heatmap and equivalent ranked destination list with stable headers, selected row/column, numeric cells, intensity legend and explicit rescale; verify rendering and browser behavior for zero, unavailable, capped, partial and independently colored population states.
- [x] 3.3 Add bounded stage search/inclusion, traffic/alphabetic ordering and explicit reorder without moving focused targets on updates; verify initial 12-stage and maximum 32-stage display limits against at least 300 stage names, selected endpoints and omitted source/display distinctions.
- [x] 3.4 Add roving keyboard movement and ordinary table/list access, plus EN/FR metric/scale/denominator help and user-guide examples; verify arrow/Home/End/Enter behavior, screen-reader names, visible focus and matching example actions without thousands of tab stops.

## 4. Repeat interpretation and useful examples

- [x] 4.1 Implement prevalence and explicitly named intensity, with affected-journey, prevalence and extra-occurrence sorts; verify 20/100 versus 8/10, ten total/nine extra occurrences, self-repeats, non-self returns, partial starts and zero-denominator cases without defect classifications.
- [x] 4.2 Replace unbounded ID-link lists with the shared matching-example picker, readable rows, deterministic newest-start ordering, time/status/path narrowing and bounded 20-row pages; verify exact criterion membership, population provenance, omitted counts and honest empty returns.
- [x] 4.3 Add deterministic observation summaries and EN/FR definitions/examples in the UI and guide; verify summaries preserve units, small support, partial coverage and goal-free navigation without inventing conversion or causality.

## 5. Expandable paths, timeline and reversible diagnosis

- [x] 5.1 Replace repeated sequence rows with the captured-prefix trie and equivalent ranked list; verify equal-edge/different-order paths, shared prefix versus exact-group counts, three-step/five-child defaults, 32-step/200-node caps, Other groups and evidence drilldown beyond visual limits.
- [x] 5.2 Add optional consecutive-repeat grouping, default off, across timeline and sequence; verify Vote times ten preserves original ordinal expansion and support counts while A/B/A, source gaps and separate captured journeys never collapse together.
- [x] 5.3 Build the compact vertical timeline with readable summary, relative elapsed basis, exact timestamps and placed/unplaced diagnostic markers; verify long paths, partial starts, absent/zero bounds, clock uncertainty, recovery and unavailable detail without attention claims.
- [x] 5.4 Connect examples to timelines, retained diagnostics and supported trace mode with return-context restoration; verify keyboard focus, sort/expansion/scroll restoration, shared-edge population choice, expired targets and no silent fallback to another population.
- [ ] 5.5 Document path/group/prefix counting, continuation states, grouping and timeline timing in EN/FR and the guide; verify 200% zoom, 390 CSS-pixel width, keyboard evidence expansion and a full pattern-to-diagnosis-to-return browser workflow.

## 6. Stable live observation

- [x] 6.1 Coalesce notification refreshes and reconcile only current-context snapshots while preserving stage order, prefix expansion, journey selection, scroll anchors and manual graph positions; verify one in-flight/one pending refresh, superseded-response rejection and no per-notification count inflation under bursts.
- [x] 6.2 Add bounded recent-live-change markers, explicit reorder and follow-latest controls with manual-navigation suspension; verify three-second cue expiry, reduced-motion equivalents, unchanged live identity after freshness clears and baseline corrections never appearing as live arrivals.
- [x] 6.3 Separate transport, snapshot freshness, pause, waiting/reconnect/reset and failed-read states; verify pause/resume with unknown queued counts, last successful display time, late event-time classification and consistent snapshot reconciliation without duplicate counts or shifted baseline bounds.
- [x] 6.4 Ensure expiry, deletion, revocation and old-app teardown override paused/returned evidence; extend protected lifecycle/stream integration tests and browser checks to prove removed evidence cannot be restored from transient navigation state.
- [x] 6.5 Document live status, pause versus collection, follow, resets, baseline corrections and not-observed-in-baseline wording in EN/FR and the guide; verify actual local OTLP intake/SSE across all seven views rather than simulating notifications alone.

## 7. Reference reuse and developer handoff

- [x] 7.1 Add bounded allowlisted presentation preferences to reference types, server validation and save/restore behavior; verify old defaults, metric/stage/sort/grouping round trips, invalid/overlong values, owner-only writes, historical bounds when saving live and safe documented rollback fallback.
- [x] 7.2 Add explicit local Copy investigation summary with exact safe context, selected observation and coverage; verify clipboard success/failure, sanitized content, no credentials/raw payloads/public transmission, and unavailable references after expiry.
- [x] 7.3 Document preferences, safe copying, trace return navigation and rollback defaults in EN/FR and the guide; verify an existing reference and a fresh heatmap investigation after reload with owner and temporary-viewer access.

## 8. Integrated qualification and experience acceptance

- [x] 8.1 Run the full local gate (`task check`, production frontend build and strict validation of this change) on a stable content identity; record passed, failed, blocked and not-run evidence separately from integrated browser and human results.
- [x] 8.2 Execute the design acceptance matrix in an isolated synthetic registry using real local intake/SSE, covering population/error combinations, late baseline corrections, gaps, expiry/unavailability, existing references and all seven view switches; record exact bounds, scope, fixtures, screenshots and regression results while preserving normal databases.
- [x] 8.3 Measure the documented density/burst workload and p95 local interaction target on the identified machine/browser; record source/display caps, a separately labeled client harness at 20 notifications/second for 60 seconds, actual five-second reset/SSE intake, UI reconciliation cadence, long tasks and network timing separately, with responsive omissions rather than hidden failures.
- [ ] 8.4 Run keyboard, screen-reader, contrast, monochrome/color-vision, reduced-motion, zoom and compact-width checks across both languages; record software/device and actual observations, keeping dedicated assistive checks pending when not exercised rather than claiming conformance.
- [ ] 8.5 Conduct an observed developer walkthrough using the question/number/example tasks and a live-versus-baseline identification task, targeting 30 seconds per task without coaching; record participant count, time, mistakes and follow-up fixes, leaving human acceptance pending until actually observed and retested.

## Workflow follow-up

- Review the completed artifacts before an explicit apply request; this planning turn does not start implementation.
- The earlier all-checks-before-archive condition is superseded for this change by Alex’s explicit 2026-10-31 archive instruction above. Preserve tasks 5.5, 8.4 and 8.5 as user-owned, unchecked follow-up; archive status does not imply their acceptance.
- Commit, deployment, public sharing, new telemetry and external-route qualification remain separately scoped work.
