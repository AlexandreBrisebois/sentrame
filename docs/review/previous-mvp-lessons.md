# Lessons from the first MVP built against the obsolete base

On 2026-10-28 the owner corrected the implementation base from `main` to `docs/journey-capability-plan-2026-10-27` and explicitly retained the visual theme/icon. The first MVP is a recoverable engineering experiment, not an implementation of the evolved plan.

## What transfers

| Observation from actual work | Carry forward | New-plan boundary |
| --- | --- | --- |
| Expandable evidence with reason, occurrence/receipt, friction and timing made the four-count fixture inspectable. | Preserve evidence-first explanations, inert text rendering, keyboard disclosure and refresh focus. | Apply them to configurable journeys and historical/live maps; a list alone is not Release A. |
| An optimistic request and timeout remained unknown; delayed authoritative persistence revised it to completed while retaining the timeout. | Regression scenarios must distinguish requests, persistence, failure evidence and recovery. | Use versioned configured rules and trusted producer binding; a self-asserted producer attribute cannot establish server authority. |
| Identical replay did not increase counts; conflicting duplicate did not overwrite accepted evidence. | Test source-identity namespaces and semantic conflicts through real ingestion/storage/read paths. | Optional/missing source identity is accepted-unassigned where appropriate and has explicit replay limits. Custom UUID fields are not a reception prerequisite. |
| A database failure rolled back a batch and returned a retryable error. | Keep commit-before-ack tests, on-disk restart checks, concurrent writer coverage and a declared crash boundary. | The evolved contract uses whole-direct-request validation/atomic rejection. The old mixed-valid partial-success policy must not be carried forward. |
| Attempt deletion survived restart, but whole-app deletion initially erased an existing attempt tombstone. The integration test exposed it; preserving that tombstone fixed it. | Preserve deletion barriers across relevant reconfiguration and test replay after restart and scope changes. | New visit/generation identities, 30-day detail, 12-month summaries, correction/deletion horizons and app re-registration need new tests. |
| Controlled-clock queries excluded expired evidence immediately, independent of physical cleanup. | Keep injected test clocks and public-read assertions; verify lifecycle updates without new intake. | Seven-day detail is obsolete. Seven days is the live baseline; detail lasts 30 days and safe supported summaries 12 months. |
| 42 delivered REST requests were executed and checked against a real server, rather than verifying a parallel undocumented fixture path. | Keep version-controlled readable scenarios, generated official protobuf bodies, stable follow-up/replay IDs and exact-request verification. | Use standard exporters/maintained receiver qualification and the new compatibility matrix; do not reuse the restrictive old profile as the target. |
| First UI review showed baseline counts, delayed updates, isolated empty scope, keyboard expansion and no horizontal overflow at 390×844. | Preserve the visual theme/icon and repeat these focused browser checks. | New Svelte maps/SSE/access controls require their own browser and keyboard qualification. |

## Superseded implementation choices

The Python/Flask stack, mandatory `sentrame.*` event schema, fixed first-persisted-vote product question, seven-day detail retention, no viewer login, and strict rejection of every unknown attribute are not current product authority. The selected stack is Svelte/Vite + Go + SQLite; standard reception, privacy projection, configurable automatic journeys, passphrase access and lifecycle behavior follow the evolved capabilities.

The old suite passed 64 tests plus 69 subtests for its own slice. Those checks do not qualify the new plan, Go/Svelte code, historical/live maps, standard exporter cells, SSE, long-lived summaries, production routes or real WFS instrumentation. Reuse meaningful scenarios and assertions, with red → green evidence for each new vertical increment, rather than translating every old test mechanically.

## Recovery and current base

- Old uncommitted MVP: Git stash `ae83f43a86799b4aec29e971dbcf524877c8e392`, created on `main`, message `sentrame MVP from main before correcting base 2026-10-28`. Includes original code, tests, documentation and visual source.
- Old generated runtime/database files: `/private/tmp/sentrame-wrong-base-runtime-2026-10-28`. Temporary recovery location, not a durable backup policy; the Git stash is the durable source snapshot.
- New implementation branch: `codex/sentrame-mvp`, based on `origin/docs/journey-capability-plan-2026-10-27` at `4a20c953381f03b4400712891f4a7e890c340835`.
- The obsolete server was stopped. No old completed task checkboxes, application code or custom protocol were applied to the evolved branch. Its implementation roadmap remains unstarted.
- Retained visual assets: [theme reference](../design/visual-theme/README.md).

Branch correction and lesson transfer are complete. The replacement MVP still requires implementation and qualification against the evolved plan. No real WFS changes, collection, deployment or publication are authorized by this recovery work.
