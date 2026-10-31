# Apply verification — 2026-10-30

Status: **11 of 11 tasks complete**. Tasks 4.1 and 4.2 were completed in the final continuation below. Earlier incomplete-status entries are historical. No archive, commit, or deployment was performed. Existing uncommitted implementation and artifacts were preserved.

## Fixes during this apply

- Timeline diagnostics now expand to retained event details and sanitized attributes, with an explicit unavailable state.
- Timeline journey selection survives switches through graph views, independently for live and baseline.
- Timeline labels include the journey ID and localized outcome. Missing nanosecond bounds represented by `"0"` no longer appear as January 1970; absent ends use the last mapped occurrence.
- Matrix cells disclose partial coverage for gaps, partial starts, and clipped continuations. A legend distinguishes observed zeroes from unavailable coverage.
- Repeat analysis names the retained-journey population, without inventing a rate or classifying repetition as a defect.
- Added French ordering-uncertainty copy and updated English/French usage guidance.
- Separated inspector selection counts from the individual journey opened in the evidence dialog. Automated checks passed after this final fix; the follow-up browser regression retest below passed.

## Passed

### Automated

- `npm --prefix web run build`: Svelte type/accessibility analysis found 0 errors and 0 warnings; production assets built.
- `npm --prefix web test`: 71 tests passed, including five new server-rendered component tests. These verify rendered text/structure, not browser input or screen-reader behavior.
- New component cases: repeated timeline occurrences; captured boundaries and missing-zero bounds; gaps/partial beginnings; diagnostic detail and missing evidence; expired selection; directed and diagonal matrix counts; unavailable versus zero cells; repeat rankings and journey actions; English/French labels and units.
- Dense component fixture: 120 journeys with 13 steps and 1,441 distinct stages. Matrix displays 100 stage columns; sequence displays 20 groups, names 100 omitted groups/journeys, caps paths at eight steps, labels continuation, and keeps prefix widths based on journey counts.
- Existing projection tests cover A-B-A-B and self-repeats, partial starts, gaps, clipping, equal-edge/different-order paths, branch limits, prefix widths, retention coverage and truncation.
- `openspec validate add-loop-aware-journey-views --strict`: passed.

### Browser / running service

Used the previous isolated synthetic database under `/tmp/sentrame-loop-NPCzPc` with the current backend and rebuilt frontend, on loopback port 18976. Added synthetic A-B-A-B journeys through actual OTLP intake; no normal workspace database was used.

- Edge keyboard traversal reached the journey selector. Arrow keys and Enter selected a journey; Tab and Enter expanded occurrence details, linked diagnostic details, and opened ordered evidence.
- Diagnostic drilldown exposed the retained synthetic exception and sanitized attributes.
- Timeline showed ten separate Vote occurrences in chronological order. After rebuilding the timestamp fix, the diagnostic journey showed its correct captured bounds without a spurious January 1970 evidence window.
- French timeline and matrix exposed matching actions and count units. Matrix accessibility exposed origin/destination headers, labeled cells and zeroes. Tab and Enter selected a matrix cell and revealed supporting journeys.
- A selected Discovery → Completion transition survived Force, Sankey and Hierarchy switches. The selected timeline journey also survived graph switches. Repeats and sequence were opened with the shared selection still present.
- Synthetic Discovery → Recovered → Discovery → Recovered produced one additional occurrence for each repeated stage and one return in each direction. Repeat evidence was opened by keyboard.
- Sequence accessible list exposed ordinal positions, prefix journey counts, eight-step continuation labels, and keyboard journey drilldown.
- Existing saved `Expired detail test` reference restored matrix layout and its expired historical bounds. Matrix disclosed partial coverage and no retained detail; switching to timeline explained that expired detail cannot be reconstructed.
- Live and baseline were separately labeled with **1 live journey / 9 baseline journeys**. Matrix and sequence showed distinct populations. Stage → Stage had 19 occurrences in the baseline only; its journey action opened the retained baseline journey instead of the live journey.

## Blocked / unrun

- Chrome initially worked, then screen capture failed again with `SCStreamErrorDomain -3812`. Native Chrome and later Edge input also failed with `noWindowsAvailable` or produced accessibility state without the corresponding visible action. The in-app browser fallback reported `The browser access host is unavailable`.
- The final inspector-count fix passed its running-browser regression retest during the follow-up below.
- Complete the remaining 4.1 browser matrix after that fix: non-default version/definition/scope/time filters through all seven views; save and reopen a fresh new-view reference with those filters; recheck independent live/baseline selections and counts, including repeats.
- Partial/gapped and dense high-cardinality cases passed component/projection checks but were not fully exercised visually in the browser. Complete those browser checks before checking off 4.2.
- Screen-reader software and other device/browser combinations were not run. No production, NAS, or external telemetry qualification was attempted.

## Final repository gate

`task check` passed after the final implementation changes: all Go packages passed, all 71 frontend tests passed, Svelte reported 0 errors and 0 warnings, all 10 canonical OpenSpec capabilities validated, and `git diff --check` passed. The change-specific strict validation also passed. No automated failures remain. Browser checks listed above remain blocked or unrun and are not covered by this gate.

## Follow-up browser verification — 2026-10-30

Read repository guidance, all change artifacts and this ledger before resuming. `openspec instructions apply --change add-loop-aware-journey-views --json` confirmed 9/11 tasks complete. Preserved every pre-existing tracked and untracked change; no implementation edits were needed for the checks completed here.

Rebuilt the backend into `/tmp/sentrame-loop-NPCzPc/sentrame` and production frontend, then restarted only that isolated synthetic registry with `SENTRAME_DB=/tmp/sentrame-loop-NPCzPc/telemetry.db`, loopback address `127.0.0.1:18976` and matching origin. No normal workspace database was opened. Sent another synthetic Discovery → Recovered → Discovery → Recovered → Completion journey through actual OTLP intake. Browser observations below used native Microsoft Edge against this service.

### Passed — running browser

- **Inspector-count regression:** opened live sequence journey `f66a4a74`, closed its individual evidence dialog, switched to matrix, selected baseline-only Stage → Stage, and confirmed **19 occurrences / 1 journey / 1 visit**, **Baseline only**, and **Show journeys (1)**. That action opened **1 of 9** baseline journeys, showing the expected Discovery → Stage → Stage… path, definition v2 and version `synthetic-a1`, instead of the previously opened live journey.
- **Separate populations:** live matrix showed Discovery → Recovered **4 / 2 / 2** and Recovered → Discovery **2 / 2 / 2**; baseline showed the same pairs as **2 / 1 / 1** and **1 / 1 / 1**, respectively. Footer separately reported **2 live journeys / 9 baseline journeys**. No combined population was substituted.
- **Repeats and returns:** live repeat rows showed Discovery and Recovered with **2 additional occurrences / 2 affected journeys** each; baseline showed Stage **19 / 1**, Vote **9 / 1**, and Discovery and Recovered **1 / 1** each. Live return rows each showed **2 / 2**, versus **1 / 1** in baseline. The baseline Stage repeat journey action opened the expected `synthetic-a1` evidence, **1 of 1**.
- **Independent timeline selections:** selected live `f66a4a74` and baseline `dfcf49f7` using native dropdowns, arrow keys and Enter. The former showed the five-step repeated path and live evidence bounds; the latter showed Discovery → Completion and separate seven-day baseline bounds. Both selections survived switching to Force and back to timeline. Force also preserved the selected Stage → Stage transition with **0 live / 19 baseline** transitions.
- Native Edge screenshots worked during this interval and visibly showed the live timeline, ordering-uncertainty disclosure, separate population footer and baseline-only inspector with **Show journeys (1)**. English workspace actions rendered after switching from French.

### Failed / blocked — browser tooling

- Browser inventory reported `The browser access host is unavailable`; explicitly retrying an in-app browser tab returned the same error.
- Chrome native accessibility actions did not update the visible view; foreground/coordinate input failed with `noWindowsAvailable`. This is not evidence of an application regression.
- Edge worked for the checks above, then became intermittent during filter setup. Native coordinate input failed with `noWindowsAvailable`, one accessibility click returned `elementHasNoFrame`, and subsequent clicks/keyboard actions left the dialog unchanged. Selecting version `synthetic-a1` succeeded, but definition remained **All definitions** and advanced time/scope controls did not open. No completed all-filter preservation claim is made from this partial setup.
- No application failure was established by the completed checks. Tool failures and unchanged UI state prevented further trustworthy browser verification.

### Still unrun / incomplete

- Non-default **version + definition + scope + exact time** together through all seven views, and save/reopen of a fresh new-view reference with those filters.
- Full browser partial/gapped, expired-detail and bounded dense-data cases in this follow-up; prior component/projection and earlier expired-reference browser evidence remain separate above.
- Full post-fix keyboard and English/French parity sweep. Native arrow/Enter journey selection passed, but that does not establish the complete keyboard/localization matrix.
- Screen-reader software, additional devices, production/external telemetry and deployment qualification.

### Passed — follow-up automated gate and cleanup

- `go build -o /tmp/sentrame-loop-NPCzPc/sentrame ./cmd/sentrame` and `npm --prefix web run build` passed; Svelte reported 0 errors and 0 warnings and production assets built.
- `task check` passed: Go packages passed, all 71 frontend tests passed with 0 failures, Svelte reported 0 errors and 0 warnings, all 10 canonical capabilities validated, and `git diff --check` passed. Full command output is in `/tmp/sentrame-loop-followup-check.log`.
- `openspec validate add-loop-aware-journey-views --strict` passed. These automated checks do not substitute for the unfinished browser matrix.
- Stopped only this run's isolated service (PID 8824); its process exited successfully and no listener remained on port 18976. Preserved `/tmp/sentrame-loop-NPCzPc` for a subsequent browser continuation.

Tasks **4.1 and 4.2 remain unchecked** because their specified integrated browser verification is incomplete. No commit, deployment or archive was performed.

## Second browser continuation — 2026-10-30

Read `CONTRACT.md`, `README.md`, `Taskfile.yml`, OpenSpec configuration, the apply skill, all change artifacts and this ledger; inspected the dirty worktree before proceeding. There is no repository `AGENTS.md` or `AGENT.md`. Apply instructions again confirmed **9/11**, with only 4.1 and 4.2 unchecked. No implementation files or existing task checkboxes were changed.

Restarted the existing binary and production assets from the preceding verified build, using only `SENTRAME_DB=/tmp/sentrame-loop-NPCzPc/telemetry.db`, address `127.0.0.1:18976`, matching origin and a synthetic owner passphrase. No normal workspace database was used. No new telemetry or reference was created in this continuation.

### Passed — partial running-browser observations

- Native Edge reloaded the service, displayed owner login, accepted the synthetic login and opened the workspace. Native arrow keys and Enter selected `synthetic-a1` and `v2 · completion.entered`; the advanced disclosure opened and accepted synthetic scope `da131391bee0e7106b5d2aaee458152a`.
- After applying and closing filters, the workspace showed that scope, version and definition together, with **1 journey / 1 completed goal**. Hierarchy showed Discovery → Stage → Vote (ten occurrences) → Completion and Vote → Vote **9 transitions**. Selecting that diagonal transition and switching to Force and Sankey retained the inspector with **9 occurrences / 1 journey / 1 visit**.
- Timeline retained the same transition inspector and filter chips. Keyboard selection of journey `da131391` showed all 13 occurrences, captured bounds **2026-10-30 14:13:41–14:14:11 UTC**, completion and ordering-uncertainty disclosure.
- Matrix retained the selected Vote → Vote cell with **9 / 1 / 1**, distinct origin/destination headers, observed zeroes, and a `Journey da131391` evidence action. Repeats subsequently showed Vote **9 additional occurrences / 1 journey**, a retained-journey population of **1**, and the same transition inspector. An Edge screenshot visibly confirmed the repeat table, scope/version/definition chips and inspector counts.
- The service remained available during the tooling failures: `/ready` and `/` both returned HTTP 200. These HTTP probes establish service availability only, not browser rendering correctness.

### Failed / blocked — browser tooling and incomplete observations

- Browser inventory and an explicit in-app browser creation retry both returned **The browser access host is unavailable**.
- Exact-time editing displayed **2026-10-29 00:00:00.000 → 2026-10-31 00:00:00.000 UTC** in the native fields and screenshot. However, the date picker's parent accessibility value still reported the earlier start and no end value. No independent restoration or evidence-bound check established that both entered exact bounds reached application state. Therefore this is **not a pass for the all-four-filter case**.
- An Edge action returned `SCStreamErrorDomain -3812` (invalid screen-capture parameter). Later sequence-switch results were inconsistent: accessibility text exposed an eight-step prefix, **8 / 20** limits and **1 continuation**, while the screenshot still showed Repeats selected and the repeat table. This stale/inconsistent result is **not a sequence-view browser pass**. Coordinate retry returned **noWindowsAvailable**.
- Native Chrome was retried. Reload removed the page content tree; capture returned a blank white image. Keyboard navigation to `http://127.0.0.1:18976/` changed the address but still exposed no page content. Chrome provided no usable application verification. No application regression was established from these tool failures.

### Still unrun / incomplete

- Confirm all four non-default filters, including exact application time bounds, through all seven views; finish compatible timeline/transition selection round trips.
- Save and reopen a fresh reference using a new view and those filters; confirm restored layout, filters and evidence.
- Finish remaining live/baseline selection/count checks across views. Earlier verified comparisons above remain valid historical evidence; none were repeated here.
- Run partial/gapped, expired-detail and bounded dense-data browser cases, including no invented adjacency across gaps, matrix stage omissions and sequence branch omissions. Earlier component/projection and expired-reference checks remain separate.
- Complete keyboard and English/French parity sweeps including evidence drilldowns. This continuation verified keyboard dropdown selection only.
- Automated suites and builds were **not rerun**: no implementation changed, and the preceding passing build, `task check`, 71 frontend tests and Go tests remain recorded above. They do not cover these browser blockers.

Tasks **4.1 and 4.2 remain unchecked**. Cleanup and final validation results follow below.

### Passed — final validation and cleanup

- `openspec validate add-loop-aware-journey-views --strict` and `git diff --check` passed.
- Reran `openspec instructions apply --change add-loop-aware-journey-views --json`: **9 complete / 11 total / 2 remaining**, state `ready`.
- Sent SIGTERM only to this continuation's service, PID **16760**. The process exited with code 0; `lsof` reported no listener on port 18976. Preserved `/tmp/sentrame-loop-NPCzPc` and all existing uncommitted work.
- Only this verification ledger was edited. No commit, deployment or archive was performed.


## Completion continuation — 2026-10-30

Resumed the same change and apply workflow, preserving all pre-existing tracked and untracked changes. Used only the isolated registry at `/tmp/sentrame-loop-NPCzPc`, with matching loopback address/origin `127.0.0.1:18976`. The original synthetic demonstration app was preserved. Registered a second synthetic app, **Loop browser bounds**, in that same isolated registry through the owner API; generated cases entered through actual OTLP/HTTP intake. No normal workspace database or direct SQL fixture writes were used.

### Passed — integrated running-browser checks (4.1)

- **All four non-default filters:** selected `synthetic-a1`, definition `v2 · completion.entered`, scope `da131391bee0e7106b5d2aaee458152a`, and exact UTC bounds **2026-10-29 00:00:00.000 → 2026-10-31 00:00:00.000**. Closed and reopened the filter dialog to remount its native date controls; both parent values now confirmed those application bounds. This resolves the previous stale-parent-value uncertainty.
- **Seven-view sweep:** Hierarchy → Force → Sankey → Journey timeline → Transition matrix → Repeats and returns → Sequence flow all retained the scoped population (**1 journey / 1 completed goal**), version/definition chips and compatible Vote → Vote selection. The inspector remained **9 occurrences / 1 journey / 1 visit**; matrix kept its selected diagonal, timeline exposed the expected 13 occurrences, repeats showed **9 additional Vote occurrences / 1 journey**, and sequence showed the same journey with **8 / 20** limits and **1 continuation**. The save dialog at the end of the sweep still displayed the exact selected bounds.
- **Fresh reference:** saved **Loop verification sequence 2026-10-30** from Sequence flow. Cleared the reference, scope, version and definition and switched to Force, then reopened it using References. Sequence layout, scope, version, definition and the expected one-journey evidence were restored. Repeated opening after a full page reload/rebuilt frontend; the reference again restored the correct layout and expected evidence. Keyboard Tab/Enter expanded the ordered journey.
- **Independent comparison selections:** in the second app, selected live `live-gap-check` and baseline `dense-00` independently. Switched through Force, Sankey and Hierarchy, then returned to timeline; both selectors still displayed their respective IDs, paths and captured bounds. Footer kept **1 live journey / 25 baseline journeys**. Sequence exposed separately labeled populations, live gap continuation and baseline omitted groups. The earlier demonstration-app matrix/repeat counts and baseline-only inspector regression above remain the supporting checks for nonzero repeat/return comparisons; the final fixes change disclosure text and absent timestamp handling, not projections or population selection.
- Existing **Expired detail test** was already verified above; newly saved sequence and additional expired/partial references also opened through the journey References UI. The `reference` URL parameter belongs to the trace entry flow; journey references in this sweep were opened using the journey UI.

### Failed, fixed, and retested — application issues found here

- **Absent zero bounds in evidence:** expanding the restored reference exposed `"0"` window bounds formatted as January 1970. The ordered-evidence dialog now suppresses an absent window and treats either zero endpoint as unavailable, matching timeline handling. Reloaded the built frontend, reopened the saved reference and expanded evidence using Tab/Enter: the spurious window and January 1970 text were gone.
- **Absent end in the bottom boundary:** the active journey card's bottom marker also formatted string `"0"` as January 1970. It now uses the last mapped occurrence for an absent end. Final rebuilt-browser retest on `repeat-keyboard` showed **Limites enregistrées · Actifs · Oct 30, 2026, 4:38:18 p.m. UTC**, with no January 1970 timestamp.
- **French coverage/drilldown text:** navigation, pending/right-censored, omitted beginning/intermediate steps, reference clipping and endpoint warnings lacked translations. Added them; localized the `[window gap]` marker and captured-status/bottom-boundary labels. In the rebuilt browser, French partial and gapped timelines and expanded ordered evidence exposed the translated warnings, **[lacune de la fenêtre]**, and **État enregistré: Actifs**. The English equivalents remained available.
- **Explicit matrix omissions:** extended the existing bounded-stage disclosure to name the omitted count. Running browser showed **100/297 · 197 omitted stages** and the equivalent **100/297 · 197 étapes omises**. A screenshot confirmed the final English disclosure and matrix cells on screen.

### Passed — bounded, partial, unavailable and keyboard/localization browser checks (4.2)

- **Dense data:** sent 24 distinct 13-step navigation journeys sharing stage A, with 289 distinct stages, plus a separate three-step ordering case. Browser initially reported **25 journeys**, matrix **100/292 stages**, and sequence exactly **20 displayed groups**, **8-step** prefixes, **5 omitted groups / 5 journeys**, continuation labels and supporting-journey actions. A screenshot at the end of the sequence visibly showed ordinal steps 1–8, broad shared-A bars versus narrow single-journey branches, continuation and omitted counts. Source work-limit/partial-coverage warnings remained visible. After adding small loop/live cases, final matrix disclosed **100/297 / 197 omitted**; a later sequence observation disclosed **7 omitted groups / 7 journeys**. These are bounded synthetic examples, not a production capacity qualification.
- **Actual gap:** began a live window, then ingested source-sequenced A, B, C where B's occurrence time lay before that window. Protected snapshot confirmed retained A/C with `gapBefore=true`, partial coverage and right censoring. Running timeline showed **Live Gap A**, a missing-intermediate-step marker, then **Live Gap C**; matrix showed **—** in every A/C cell and no invented A → C edge or evidence action. Sequence stopped its contiguous prefix at A and labeled **1 continuation / 1 partial**. Repeated with fresh synthetic scopes during localization retests; final French timeline and ordered evidence preserved and translated the gap. Historical query setup alone returns whole captured journeys, so the initially prepared historical “Gapped …” references were not used as proof of clipping; the live-window cases above provide that evidence.
- **Partial beginning and continuation:** opened **Partial retained segment**, a scoped reference from B00-03 through B00-08 within `dense-00`. Browser showed six occurrences, partial/earlier-visit and continuation disclosures, preserved original captured bounds and a separate evidence window. English and French exposed equivalent warnings and evidence actions. Keyboard Tab/Enter expanded the first occurrence and showed its retained `stage.name` attribute; screenshot confirmed French warnings, focus and details.
- **Expired range:** opened **Expired browser**, whose fixed bounds precede the 30-day detail cutoff. French matrix disclosed partial coverage and no retained detail; timeline explained that expired detail cannot be reconstructed. English exposed the same state. Repeats named a retained population of **0**, and sequence stayed empty with its display limits. No summaries were turned into paths. This verifies an out-of-retention selection, not a newly observed timed expiry/deletion event.
- **Repeat/return fixture:** A → B → A → B → B showed B **2 additional occurrences / 1 journey**, A **1 / 1**, and each directed non-self return **1 / 1**. English and French exposed matching rows, units, definitions and supporting evidence. From References focus, six Tabs reached the first English repeat action; Enter opened the expected five-step journey. In French, ten Tabs from the selected repeat-view control reached the equivalent action; Enter opened **1 sur 1** of the same evidence. Escape closed the evidence dialog and restored focus.
- **Matrix keyboard access:** French matrix labeled A → B **2 / 1 / 1**, B → A **1 / 1 / 1**, B → B **1 / 1 / 1**, and A → A as observed zero. Eleven Tabs from the matrix-view control reached A → B; Enter selected it and exposed the **2 / 1 / 1** inspector. Three further Tabs reached its supporting journey; Enter opened the expected evidence. Earlier English matrix keyboard checks remain recorded above.
- **Timeline and sequence keyboard/drilldowns:** arrow/Enter journey selection worked for filtered, partial, live and baseline timelines. Tab/Enter expanded occurrence details and reached ordered evidence; final French gap evidence opened and expanded entirely with keyboard actions. English restored-reference and French dense-sequence evidence both exposed the complete ordered supporting path beyond the eight-step visual prefix; Tab/Enter expanded the evidence card. Earlier diagnostic keyboard drilldown passed above; no diagnostic behavior changed in this continuation.

### Passed — final automated gate

- Final `npm --prefix web run build`: production assets built, Svelte **0 errors / 0 warnings**.
- Final `task check`: all Go packages passed; **72 frontend tests passed / 0 failed / 0 skipped**; Svelte **0 errors / 0 warnings**; all **10 canonical OpenSpec capabilities passed**; `git diff --check` passed. Full output: `/tmp/sentrame-loop-completion-check.log`.
- Added a rendered French timeline regression covering eight navigation/partial/window/reference warnings, missing-step text and continuation; strengthened dense component coverage for the explicit omitted-stage count. These are component checks; the browser checks above independently exercised actual input, state and drilldowns.
- Task 4.1 was checked after its sweep/reference/comparison checks; rerunning apply instructions confirmed **10/11 complete**. Task 4.2 is checked after the final fixes, browser retests and passing gate.

### Tool failures, remaining scope and cleanup

- Earlier in-app browser/Chrome failures remain historical tool failures, not application passes. Native Edge recovered after a supported **Raise** accessibility action and provided current screenshots and changing UI state throughout this continuation. No required integrated browser case remains blocked.
- **Not run / outside this change's local completion gate:** dedicated screen-reader software, additional devices/browser combinations, production/NAS capacity, external telemetry/exporters, deployment and timed retention-expiry qualification. Component/static checks are not substitutes for these environments.
- Preserved the original synthetic app, the new synthetic app, generated fixtures and saved references under `/tmp/sentrame-loop-NPCzPc`. No commit, deployment or archive was performed. Final service-stop, strict validation and apply results are recorded below.

### Final completion and cleanup

- `openspec validate add-loop-aware-journey-views --strict` passed.
- Apply instructions returned **11 complete / 11 total / 0 remaining**, state **all_done**. Output: `/tmp/sentrame-loop-completion-apply.json`.
- Stopped only this continuation’s isolated service, PID **18517**, using SIGTERM. Its process exited with code 0; no listener remained on port 18976. The isolated registry and both synthetic apps remain preserved.
- Final `git diff --check` passed. All original dirty files remain present; no commit, deployment or archive was performed.
