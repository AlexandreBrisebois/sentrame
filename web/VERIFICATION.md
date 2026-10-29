# Initial frontend evidence — 2026-10-28

This is the initial implementation evidence log. Later canvas, timing, diagnostics, field guidance and temporary-viewer work supersede its original test counts and duration-sized-node design. The current suite has40frontend tests; current behavior and qualification are in the [user guide](../docs/USER-GUIDE.md) and [workspace verification](../docs/verification/workspace-design.md).

Scope: `web/**` only. Svelte/Vite/TypeScript, warm paper/forest green visual theme and the preserved linked-thread mark.

- `npm --prefix web test`: observed meaningful red for stale-response acceptance and unscaled node radius, then green. Later added UTC-control regression and observed red/green for histogram quantile ranks/underflow/overflow. Current result: 8 passed, 0 failed. Final summary exclusion coverage was also observed red/green, including a stage with only censored/unavailable observations.
- `npm --prefix web run build`: passed; `svelte-check` reported 0 errors, 0 warnings, followed by successful Vite production asset build. Lockfile records resolved dependencies. Approximate assets: JS 113 kB (42 kB gzip), CSS 22 kB (6 kB gzip), before any subsequent integration edits.
- `git diff --check -- web`: passed.
- No mocked API result is used as evidence of ingestion, durability or backend authorization. Browser/API end-to-end qualification is owned by the integration lead and is not established by these helper/build checks.

UI uses server-authorized snapshots and full replacement on SSE invalidation/reset; stale request generations are rejected. Display pause does not stop collection. UTC date controls preserve the stated timezone. Live baseline bounds remain fixed at observation start and counts are separate. Table/sequence alternatives preserve intermediate evidence and captured journey boundaries. Reference endpoint occurrence selectors are local display choices; saved filters/definition/endpoints persist through the API.

Median node area uses one explicit timing basis and definition across the map, with a 30-minute visual cap and numeric values. Stored summary histograms merge bucket counts by exact dimension key; bin ranges expose approximation, and p90 requires at least 20 completed samples. These UI calculations do not establish source clock/activity quality or metric arithmetic.

Main controls and standalone explanatory prose have French translations. Some small metadata phrases and backend-provided diagnostic/quality text remain English. Full English/French task parity and WCAG conformance are not claimed. Browser layout, keyboard, session-expiry and real SSE checks require the running integrated service.

Integration corrections: phone filters now use a constrained single-column grid below 480px; actual remeasurement remains with the lead. Window gaps/right censoring are explicit in journey paths. Lifecycle offers a separate typed-confirmation control for one app-wide summary UTC day using `{app:true,bucket}`; whole-app deletion remains a distinct operation. Summary timing rows merge and display the server-provided exclusions, including no-complete-sample stages and explicit unknown coverage for older summaries.

Final integration corrections: date/time rendering explicitly uses UTC (including live/baseline bounds) with a visible UTC suffix; a Toronto-timezone assertion was observed red then green. Discovery/available version/scope choices now use the app-wide bounded snapshot fields, falling back to retained evidence for an older response shape; empty-live discovery regression observed red/green.

Saved reference integration: first/last endpoint selections are persisted and restored. Selected reference ID and occurrence choices accompany snapshot queries, so the server projects map/counts/evidence consistently within captured boundaries. The frontend no longer independently clips or drops journeys in a reference result. Changing endpoint choices or clearing the reference refreshes the authoritative snapshot.

Evidence timestamp transport: time/receipt/observed accept decimal strings and remain strings in sanitized copy. Display truncation converts decimal nanos with BigInt before Date milliseconds, avoiding rounding across a millisecond boundary. The precise-value copy/display regression was observed red then green. Numeric Journey timestamps remain approximate display inputs under the shared contract.
