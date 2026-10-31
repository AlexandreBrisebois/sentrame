# Tasks

Alex approved this reviewed implementation on 2026-10-31. Alex also approved the graph-first flyout refinement after measured header review. Prior software evidence is retained in IMPLEMENTATION-HANDOFF.md; the approved refinement has final-code evidence in GRAPH-FIRST-HANDOFF.md. Task 5.2 remains pending actual native zoom and physical screen-reader acceptance. Prior archived manual checks remain separate and user-owned.

## 1. Compact orientation and direct navigation

- [x] 1.1 Remove the task dropdown, retain all seven view buttons with pressed state and a single horizontally scrollable phone strip, and replace every investigation-task focus fallback with a stable active-view target; verify saved view destinations, trace returns, example-dialog close and unavailable-target recovery in focused navigation tests.
- [x] 1.2 Move repeated question/unit/selection and applicable L/B/error explanatory prose into a named info flyout using existing visual tokens; verify Enter/Space, expanded/controls state, hidden-content tab order, visible focus and unchanged metric/population/evidence actions in EN/FR render/browser checks. Update user-guide navigation/help wording in the same patch and check links.

## 2. Compact coverage with visible material caveats

- [x] 2.1 Replace permanent coverage paragraphs with a compact flyout preserving all values and denominator/exporter-completeness explanations; verify closed/open reading order and overlay geometry, keyboard/touch activation and retained open state across refresh, with no count/filter changes.
- [x] 2.2 Derive concise visible scope-specific warnings for source caps, unknown metadata, current/baseline partial or censored evidence, graph-local/server omissions and record-list omissions; verify each independently, combined and zero-omission cases, unchanged ordered denominators, local error/population identity and existing stale/comparison restrictions in focused tests. Document the compact Coverage behavior alongside the patch.

## 3. Historical refresh stability

- [x] 3.1 Put routine refresh status into allocated compact space and retain successful evidence during explicit same-context Refresh; use delayed synthetic reads to verify stable canvas bounds, SVG identity and camera before/during/after unchanged background and explicit reads, plus first load and failed-read/retry distinctions.
- [x] 3.2 Preserve valid selected evidence, open detail, focus, page/container scroll and local layout positions across same-context refresh; verify selected node/edge/example cases, all three graph layouts, changed revision, superseded filters, deletion/expiry and access loss. Confirm invalid targets clear with a valid focus successor, with no old-query evidence or stale raw-detail resurrection. Record whether any changed-data re-layout issue requires a separate proposal rather than an unsupported graph rewrite.

## 4. Repeat-control sizing

- [x] 4.1 Scope full-width input rules to appropriate field types and give repeat checkboxes intrinsic size/nonshrinking space with wrapping label text in timeline and sequence; verify computed checkbox/text geometry, full-label activation, keyboard toggling and unchanged original ordinals/counts in EN/FR at desktop and 390px. Keep changes scoped to investigation controls.

## 5. Integrated review

- [x] 5.1 Run frontend tests, Svelte checks/build, strict change validation and whitespace checks; record results and browser asset/source identity. Exercise all seven direct view buttons and help/Coverage with relevant user-guide wording in a synthetic isolated browser; verify no unavailable focus IDs or missing actions.
- [ ] 5.2 Coordinate a native browser slot before interactive QA; observe actual 200% browser zoom for new wrapping/disclosures and physical EN/FR screen-reader operation of the changed controls and refresh status. Record performed checks separately from unperformed checks and from the three user-owned archived acceptance tasks; do not mark proxies as native or screen-reader passes.

## 6. Approved graph-first refinement

- [x] 6.1 Consolidate mode/account/app navigation and remove repeated workspace headings; keep all seven view destinations in a single responsive strip and References/highlights available in the existing menu. Verify keyboard-revealed phone views, account actions and all destination/return paths in EN/FR.
- [x] 6.2 Convert Help/Coverage and detailed Context to viewport-clamped anchored flyouts with44px controls. Verify Enter/Space, Tab, Close, Escape, outside/touch dismissal, opening another flyout, focus return and no canvas movement, including bounds/scroll on390px.
- [x] 6.3 Consolidate visible population/bounds/active narrowing and refresh/pause status into compact rows, with no dedicated empty reservations; verify unchanged/changed reads, two populations, actual partial/limit/error warnings and retained denominators with flyouts closed.
- [x] 6.4 Measure and visually inspect final1440×900 and390×844 EN/FR layouts against468/842/900px graph-start baselines. Record actual graph room, no page overflow, all minimum target sizes and refresh/flyout stability. Update user-guide wording and durable final handoff with exact checks and remaining native limits.

## Archive disposition — 2026-10-31

Alex explicitly instructed “sync and archive the spec” after the implementation report disclosed12/13 tasks and pending native zoom/screen-reader acceptance. This authorizes archival with that disclosed task outstanding; it does not verify task5.2. Its checkbox remains unchecked. Actual native200% browser zoom and physical EN/FR screen-reader observations remain pending. Software/headless results are separately scoped in GRAPH-FIRST-HANDOFF.md. See ARCHIVE-NOTE.md for the supported closeout and exception; older archives and their manual checks are unchanged.
