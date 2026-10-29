# Sentrame decision register and review questions

Status: review draft, 2026-10-27. User intent below is trusted conversation context supplied for this task, not a transcript published as product evidence. Proposed values are not decisions. Today’s direction evolves/replaces the initial concept; the initial synthetic defaults have no continuing product authority.

## Traceability

| ID | User decision / constraint | Contract owner and acceptance evidence |
| --- | --- | --- |
| D01 | Research, challenge, find blind spots and write capability specs; no application/scaffold implementation now. Root North Star first, iteratively revisable. | README; proposal; all tasks remain unchecked. Documentation validation only. |
| D02 | WFS first; focused journey insights; low-footprint Home Server; future apps separate. | telemetry-intake application attribution; data-lifecycle resource bounds; feasibility/value and resource qualification; selected stack is recorded below. |
| D03 | Standard OTLP; existing SDK/exporter redirection; no mandatory proprietary SDK/custom transport. Ordinary optional OTel instrumentation is allowed. | telemetry-intake compatibility matrix and unmodified producer scenario; modeling missing-context behavior. |
| D04 | Anonymous random visit/session; several journeys per visit; no cross-app/visit/person tracking. | journey-modeling isolation, multiple journeys and tab ambiguity; data-lifecycle minimization. “Anonymous” product wording does not erase pseudonymous correlation risk. |
| D05 | Discover received events; user-defined steps, conditions, goals and deterministic outcome classification through simple controls; little/no AI. | journey-modeling preview/conflict handling and explicit completion/positive/failed/flaky distinction. |
| D06 | Stable step names preserve semantics across app versions; rename when semantics change. Automated technical metadata where available. | journey-modeling version compatibility; telemetry-intake unknown metadata; exploration normalized comparison. |
| D07 | Continuous collection; 30-day ordered detail with auto-delete and 12-month aggregate history, with no individual drilldown after detail expires. | data-lifecycle event-time expiry, corrections and deletion; summary grains, mergeable distributions and correction/deletion mechanics remain technical defaults. |
| D08 | Historical app/version/date map, weighted traffic; rare paths show counts/percentages; default potential dead ends show evidence. | journey-exploration counting units, low support, right-censoring, path drilldown. |
| D09 | Live SSE: selected app or isolated visit, switch freely preserving filters. The established voice choice is a toggle between a clean live window and the selected app's last-seven-day accumulated traffic with the chosen live visit over it. | live-journeys scope/window/latency/replay scenarios. Preserve compatible app/version filters and distinguish counts. This is app history, not same-visit history or seven-day detail retention. Future invitations include app history only when the owner explicitly grants it. |
| D10 | Exceptions at valid linked points, available type/message/stack/runtime/version/context, recurrence/affected sessions, scrubbed copy for developer/Codex. No confidence badge or AI root cause. | journey-diagnostics linkage, grouping/copy and missing data cases. |
| D11 | Exception does not end journey; inspect subsequent activity/recovery and eventual goal. Repeated WFS votes are healthy; measure repetitions/duration and emitted exhaustion/agreement/departure reasons. No “tired” inference from silence. | journey-modeling healthy loop/exit scenarios; diagnostics recovery; exploration default stopping evidence. |
| D12 | Passphrase login and cookie invitations for another viewer, superseding initial no-login idea. | viewer-access separate sessions, permissions, revocation and scoped live sharing. Viewer invitations are not the excluded WFS acquisition funnel. |
| D13 | Inspect WFS accessibility/auth patterns without blindly cloning; simple controls; energetic purple exploration. | exploration keyboard/table/status/EN-FR review proposal; design auth options; feasibility source evidence. Palette not finalized. |
| D14 | Today builds on/replaces initial idea; no need to preserve it as a separate design. | Intended single replacement package: broad capability draft; previous synthetic seven-day/fixed-rule semantics explicitly obsolete. The concrete planning-artifact reconciliation and documentation-only draft PR are approved. |
| D15 | The user reports that WFS moves the user to a completion screen. | Model completion-screen observation as a visit-specific outcome, separate from shared convergence, its cause and emitted end reason. This report does not establish current telemetry. Add another-participant, initial-hydration, repeated-render and departed-visit cases. |
| D16 | Use “Home Server.” Deployment is Docker Compose with the associated Compose file and `.env`, behind an existing shared Cloudflare Tunnel providing HTTPS. | Proposed deployment contract in documentation only; no actual Compose/`.env`, tunnel, credentials or network changes. Separate edge TLS, origin transport, app auth and route-specific SSE/OTLP qualification. Exact hardware and numerical budgets remain unmeasured; sufficient RAM and CPU options are reported. |
| D17 | The owner can emit events and instrument the app. Journeys map to emitted events; the owner selects a goal event from discovered events and inspects paths from journey beginnings to it. | Event discovery/goal-selection and ordered-prefix scenarios. Instrumentation ability is established; exact names, optional linkage, privacy and existing emitted coverage still need specification. No custom SDK/protocol. |
| D18 | Exceptions and warnings can be emitted as event types in app code. | Discoverable evidence/UI classifications using standard OTel carriers, names, severity and exception conventions. Severity alone does not fail/end a journey; preserve continuation, recovery and later goal success. |
| D19 | Accepted “time at stage”: node size for duration, edge thickness for traffic, median plus slower cases and safeguards for active/background/missing timing. Examples include recipe viewing and individual voting stages within healthy loops. | Stage-occurrence timing and exploration/live scenarios; sample/completeness, units, bounded legend and accessible numbers/table. Do not infer cognitive effort or classify larger as bad. Exact slower percentile/support threshold and instrumentation/activity policy remain design choices. |
| D20 | Release A is the complete useful WFS-only workflow: standard intake, discovered events/simple goal, historical/live maps, seven-day baseline, stage timing and warning/exception recovery. | Thin tracer bullets build toward A; they do not make a list-only final release. Multi-app/richer comparison/invite/config requirements stay retained roadmap B. |
| D21 | Automatic journey boundaries: first mapped event starts; selected active goal ends; next new mapped occurrence restarts; 30 minutes of inactivity closes inactive, not failed. Browser tabs stay separate. | Modeling scenarios; supported replay deduplicates before segmentation, event-time/late revisions stay honest. Server-only logical scopes do not require tabs. |
| D22 | Saved reference views support goal-backward, beginning/end segments and origin-forward exploration. | Preserve filters/definition versions, recorded boundaries and partial evidence. Save query definitions, not raw journeys beyond deletion/retention. |
| D23 | Future invitation A: one live session, read-only, owner explicitly includes/excludes seven-day app baseline and exception details; expiry/revocation, no configuration edits. Keep broader workspace B later. | Viewer-access requirements retained; invitations deferred from first release A. Exact lifetime/redemption defaults unchosen. |
| D24 | Retention B: 30 days detailed journeys, 12 months aggregated counts/conversion/timing history, no individual drilldown after 30 days. | Data lifecycle, mergeable distribution/denominator/privacy and deletion-horizon safeguards. No averaging medians or reconstructing individual paths from counters. |
| D25 | Stack A: Svelte/Vite + Go + SQLite with maintained OTLP components. RAM is reported sufficient and CPU options available. | Selected architecture, not a new ballot. Resource use remains bounded and unmeasured; no exact hardware/unlimited budget inferred. |
| D26 | External integration-agent guide must inspect the app, question developers, propose points of interest and missing warning/error/recovery instrumentation, with developer review before edits. | Separate public integration guide/capability, not internal harness or AI runtime. Explain question/value, evidence, minimal emitter boundary, standard carrier, authority, privacy/noise and synthetic validation. |
| D27 | Client-only, server-only and combined logical journeys are supported when properly emitted; combined evidence uses the same app-scoped correlation ID. | Correlation is not trace/record identity, authority, order or causality. Preserve sources, bounded lifetime and no cross-app/person/cross-visit joins. No side/tab prerequisite for basic intake. |
| D28 | On 2026-10-28, correct the implementation base to `docs/journey-capability-plan-2026-10-27` while retaining the first MVP’s visual theme and icon. | Preserve warm-paper/forest-green visual language and linked-thread mark in docs/design/visual-theme; the earlier exploratory purple direction is superseded. The prior MVP is recovery evidence, not implementation of the evolved plan. |

## Remaining technical defaults — main choices are settled

The five main decision groups are closed: release A with B retained; automatic boundaries; narrow deferred invite design; 30-day detail/12-month summaries; Svelte/Vite + Go + SQLite. Do not reopen them as unanswered choices.

| Area | Remaining focused design work |
| --- | --- |
| Event contract | Actual safe event names, selected active goal, correlation/occurrence fields, producer authority and current emitted coverage. Ability to instrument is confirmed. |
| Intake/route | Receiver embedding versus sidecar, pins, origin transport and shared-Tunnel path qualification; request/queue/body/time limits and explicit durability failure model. |
| Event-time quality | Allowed lateness/skew, missing-time policy and deterministic reconciliation of late events against accepted first-event/goal/30-minute rules. Correlation alone does not order cross-clock events. |
| Timing | Entry/exit sources, activity/idle/fallback rules, slower percentile/support threshold, visual scale and bounded refresh. Accepted median/label/node/edge safeguards are not open again. |
| Aggregates | Concrete grains/dimensions, mergeable timing distributions/precision, denominator/follow-up policy, finalization/correction horizon, small-count privacy and deletion/backup semantics after raw expiry. Twelve-month retention is accepted. |
| Future invitations | Exact expiry/use/rotation defaults and earlier-session/after-session visibility within the selected one-live-session design. Broad workspace mode remains later. |
| Resources | Actual deployment architecture/ISA/filesystem, workload and numerical qualification limits. RAM/CPU availability is reported; fit and disk growth are not measured. |
| UX/collection | Detailed accessibility/device/language checks, consent/enablement and diagnostic allowlist before real collection. Goal versus saved-query semantics remain explicit without changing captured boundaries. |

The baseline is settled as selected-app accumulated seven-day traffic with a live visit overlay and clean-live toggle. Saved views do not redefine capture boundaries. A longer-lived business correlation is not permission for person or cross-visit tracking.

## Superseded initial assumptions

| Initial draft assumption | Today's replacement |
| --- | --- |
| Synthetic-only product scope | WFS-first product direction with later separately approved synthetic qualification and real integration. No app implementation now. |
| Seven-day evidence retention | 30-day ordered detail. Seven days remains only the optional live historical baseline. |
| Discovery → first persisted vote as sole question | Configurable explicit goals and outcomes; first vote can be a milestone, healthy multi-vote loops and family agreement matter. |
| Mandatory custom event/attempt profile and rejection of unknown attributes | Standard telemetry reception without custom identifiers; optional recommended context for richer analysis; privacy filtering distinct from protocol validation. |
| One fixed outcome rule, no viewers | Versioned configured deterministic mappings; passphrase login and scoped viewer invitations. |
| Attempt-only identity | Random application visit/session can contain several journeys; no cross-visit/person joins. |

## Planning publication boundary

The original planning-file reconciliation and a documentation-only draft PR are approved. This covers the evolved proposal/design/tasks, retirement of the two obsolete deltas and the capability set including integration guidance. It does not authorize application/harness implementation, canonical spec sync/archive, merge, deployment or production tests. The historical narrow defaults are superseded, not a competing product design.

### 2026-10-28 — Graph workspace revision approved for local implementation

The owner selected a quiet full-canvas workspace with three switchable views: horizontal left-to-right hierarchy, directed force-based layout, and Sankey. Sankey bars/ribbons represent traffic volume; time-at-stage moves to the inspector and accessible tables, superseding duration-sized nodes. Keep the warm-paper/forest-green theme, linked-thread mark, small independently capped arrows and explicit healthy loops/return paths. Configure events through contextual node popups with preview and prospective activation. Explore and save reference views from stages, stage pairs or the full workspace; retain captured boundaries and filters. Show rare paths, endings, warnings/exceptions and recovery through graph lenses with support/evidence context. Clear selected journey/visit scope with a visible X without deleting data. This authorizes local synthetic implementation only; retained roadmap B and external qualification boundaries do not change.

### 2026-10-28 — Useful app-wide error diagnostics approved

The owner requested standard .NET/Next-like exceptions, stacks and custom error messages to debug the application without being the user. Local synthetic implementation now includes bounded `scrubbed-v1` diagnostic detail, protected app-wide error search, valid trace/context links and sanitized copy. Raw request collection, person joins, automatic root cause, real application integration and exporter qualification remain excluded. Text scrubbing is imperfect and controlled producer messages remain necessary; 30-day detail/12-month text-free summaries are unchanged.

### Explicit scope extensions — 2026-10-28

- The owner authorized app-wide temporary read-only viewers, including retained history/live/reference/diagnostic browsing, with expiry and revocation. This explicitly adds a broader mode; it does not widen an existing visit-only grant. Technical defaults: one-use redemption;1h/24h(default)/7d absolute expiry; separate viewer cookie; server read allowlist; active SSE and historical-session revocation checks.
- The owner requested DockerHub delivery through GitHub Actions, image `brisebois/sentrame`, multi-architecture `linux/amd64,linux/arm64`, and repository Compose/`.env.example` files for NAS deployment. Publishing/actual NAS qualification remain separately evidenced; initial no-commit/no-push instructions have not been withdrawn. WFS changes are not part of this packaging work.
