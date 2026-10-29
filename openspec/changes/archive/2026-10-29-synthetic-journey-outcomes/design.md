# Design

## Context

See [proposal.md](proposal.md) for purpose and [README](../../../README.md) for the North Star. This is the evolved planning proposal, not implementation authorization. All architecture choices below are proposals. The repository currently contains documentation/OpenSpec only; canonical specs are empty, and there is no application, harness or measured deployment.

Research: [feasibility and comparables](../../../docs/research/sentrame-feasibility.md), [stack/resource analysis](../../../docs/research/sentrame-stack-and-resources.md), [decision traceability](../../../docs/planning/sentrame-decisions.md). These distinguish source observations, user decisions and proposed defaults. The original narrow concept is superseded by this approved planning reconciliation. The historical change identifier is retained; it is one evolved plan, not a competing implementation.

## Goals / Non-Goals

**Goals:** Keep standard reception independent of optional product semantics; make a small sequence store sufficient for real path evidence; make every derived count explainable/correctable; bound memory/disk/work; separate ingestion trust, anonymous application context and viewer authorization.

**Non-Goals:** A telemetry vendor platform, identity resolution, arbitrary graph queries, session replay, automated intent/root-cause inference, a general SQL/dashboard builder, custom OTel transport or mandatory SDK. No scaffold, app, agent harness, benchmark service, live WFS change, deployment or production test is created in this pass.

## Decisions

### Status of choices and activation gates

“Confirmed” means D01–D27 in the decision register. The complete WFS release A, preserved roadmap B, automatic boundaries, narrow deferred invitations, 30-day detail/12-month summaries, Svelte/Vite + Go + SQLite, Home Server Compose/shared Tunnel, event-driven goal discovery, timing and client/server correlation are selected. The owner reports sufficient RAM and CPU options. Technical limits, pins, summary precision, timing thresholds and invitation expiry remain a focused design review before implementation; no hardware or measured fit is invented.

The requested scope permits research and specification now. Material choices gate their dependent future implementation tasks; they do not gate finishing this package. Independent Astra review is a separate review step. Planning reconciliation and a documentation-only draft PR are approved; application implementation, canonical sync/archive, merge and deployment are not.

### A. Selected stack and bounded deployment

Selected stack: Svelte/TypeScript + Vite static frontend, Go serving assets/APIs/SSE, SQLite with indexed sequence rows and bounded JSON, and maintained OTLP components. Comparative alternatives remain research rationale, not an open stack ballot. Keep stable app/version/definition identity from day one without implementing roadmap B upfront or creating speculative extension frameworks. The first release includes working historical/live maps, timing and recovery; a list/row tracer proves seams only.

The deployment direction is Docker Compose with a corresponding Compose file and `.env`, behind the existing shared Cloudflare Tunnel providing HTTPS. The [Home Server deployment contract](../../../docs/proposals/home-server-deployment-contract.md) is documentation only: this pass creates no deployment/configuration files, connector, credentials or network changes. The shared connector's reachability and lifecycle are outside the app's ownership. Browser-to-edge TLS, connector-to-origin HTTP/HTTPS and Sentrame's passphrase/producer authorization are distinct boundaries.

Reuse the official Collector OTLP receiver unchanged, with a synchronous storage consumer and atomic whole-direct-request validation/commit if embedding is proportionate. Its ordinary consumer path does not provide per-item partial-success counts. Alternative: an official Collector sidecar, forwarding to a standards-compliant app receiver that still must exist. Choose from lifecycle complexity and measured resource/build cost. A direct handler implementation using official Protobuf definitions is a third choice only if maintained receiver overhead proves material or a separately accepted response contract requires it.

```text
Existing SDK/exporter                         Viewer browser
  | standard OTLP                               | cookie + filters
  v                                             v
OTLP receiver --> app attribution --> scrub --> transactional evidence
                       |                              |
                       | config + policy              v
                       +-----------------------> bounded derivation
                                                      |
                                          query snapshot + SSE updates
```

An optional Collector sits before this path, with its own documented acknowledgement/queue boundary. A batch processor can enqueue and acknowledge before export, then merge/split SDK envelopes. No asynchronous batch/queue before storage may retain a final-commit acknowledgement claim. Atomicity covers one direct storage-receiver request, not the original SDK envelope after rebatching. Existing SDK processors still own producer batching; do not introduce a proprietary producer layer.

### B. Transport, signal and compatibility plan

The proposed target includes logs, traces and metrics for OTLP/HTTP binary Protobuf, HTTP JSON Protobuf and gRPC Protobuf, including none/gzip compression. Profiles are excluded. See [telemetry intake](specs/telemetry-intake/spec.md) and the research matrix. Each cell must be qualified before being advertised; an early logs-only tracer bullet cannot claim full OTLP support.

Publish a second qualification view for the actual tunnel route. Receiver/origin gRPC support does not make it work through a public Tunnel hostname; current Cloudflare documentation excludes that path. Use qualified standard OTLP/HTTP where applicable and retain origin/private gRPC as a distinct target without opening an alternate route implicitly. Browser exporters still need their own supported HTTP/CSP/CORS/lifecycle checks. SSE must stream `text/event-stream` and qualify flushing, heartbeat/idle behavior, reconnect, auth and revocation through the shared tunnel; HTTPS reachability alone proves none of these. Shared edge/body/time/cache policies can change path behavior. No actual tunnel configuration has been inspected or changed.

Recommended qualification sources: pinned official JS browser HTTP exporters, the actual WFS .NET/Serilog packages and a pinned official Collector forwarding path. Test endpoints, mixed-resource requests, zero-commit whole-request rejection, unsupported signals, gzip bounds, empty requests and retry/error behavior. Explicitly map permanent validation to InvalidArgument/HTTP 400 and recoverable writer/overload errors to Unavailable/HTTP 503: the inspected receiver maps a generic permanent consumer error to Internal/HTTP 500, so do not rely on its stale comments. Per-item acceptance/counts require a separately qualified response integration. Browser JS uses HTTP; server gRPC support does not imply browser support. Preserve nanosecond precision.

Initial metrics support means bounded sanitized receipt inspection, with no value-derived summaries. Privacy omission can collapse stream identity; temporality, resets, overlaps and histogram semantics must be qualified before arithmetic, without retaining forbidden attributes. Receiver-health counts are distinct from summing application metrics. Account separately for log records, spans, span events, metric points and bucket/attribute bytes. Traces may supply mapped evidence/exception context; sampled traces are not an exact business denominator. Logs/events are the recommended product carrier without custom fields required merely for reception. Unmapped technical telemetry stays received/unassigned.

### C. Application, producer and viewer trust boundaries

Application registration has immutable app identity, display name, allowed producer routes/credentials, configured service/resource selectors and metadata roles. Multiple services can map to WFS. Generic `service.name` is not globally unique and self-asserted resource attributes are not authorization. Intake route/credential binding establishes permitted app scope; selector ambiguity cannot merge apps. Unknown service/version values remain unknown.

Separate server-authenticated observations from browser-reported observations using trusted configured producer namespaces that a resource attribute cannot self-assert. Browser evidence can establish its own rendered-screen claim, while a committed vote requires trusted persistence evidence. Browser-embedded keys cannot establish server authority; CORS is not producer authentication. Proposed web path: app-origin forwarding with ordinary OTLP configuration and restricted routes/quotas, or deliberately public intake scoped to untrusted browser claims. Forwarding must preserve the authority distinction. Review reachability/CSP/CORS and never embed a viewer passphrase or server-authority credential in a browser SDK.

An emitted “vote committed” from an approved server origin can support a persistence milestone. A similarly named browser event cannot acquire that authority. Server persistence precedes follow-up publications in inspected WFS source, so an HTTP error can coexist with a saved vote. Commit-to-emission failures still leave an observation gap. Add an outbox only if a later requirement needs that stronger guarantee; there is no exactly-once product promise now.

Viewer access is an independent session/grant domain. Tracking context never contains viewer cookie/token/person or WFS household identity. Collector/receiver ports, protected viewer APIs and administration should have distinct access policies even if hosted by one process.

### D. Optional recommended evidence mapping

The controls map existing safe attributes/event names, rather than requiring a new Sentrame wire schema. A reference contract for a later WFS integration can recommend:

The owner explicitly can emit events and instrument the app accordingly. This is settled capability, not an open feasibility question. Journey steps map to emitted events; the UI lists discovered events, the owner chooses a goal event, and the view traces matching journeys from their known beginning to that occurrence. Missing/expired starts or hidden intermediates remain marked, and later continuation/recovery remains inspectable. Selecting a diagnostic event for investigation does not silently activate a terminal/failure rule; changes to lifecycle definitions remain explicit and versioned.

In standard OTel terms, named events are represented in the logs data model; warnings use severity and may have an event name, and exception conventions provide fields such as type/message/stack. Span events remain part of traces with trace context. SDK APIs and event-name support vary, so use existing standard logging/event/exception APIs and configured safe-field mappings where needed. Event, Warning and Exception can be discoverable UI classifications; they are not invented OTLP signals or transports. A WARN/ERROR record or exception is evidence at a point in a journey, not automatic failure/termination. Metrics remain aggregate signal data, not individual journey events.

| Meaning | Optional context and constraint |
| --- | --- |
| Visit | Random app-scoped session/visit value, fresh per visit, no persistent person/device key; no `session.previous_id` linkage. |
| Journey | Optional random journey value or configured start/goal/restart derivation within one visit. Several journeys may occur in that visit. |
| Step | Stable semantic event/step name with documented meaning; unchanged across app versions only when behavior is unchanged. |
| Action linkage | Optional random operation/action context linking browser intent, server commit and exception; standard trace/span relationship when verified. Never recipe/household/member ID. |
| Replay identity | Optional stable event identity preserved across retry; receiver does not make it mandatory for ordinary logs. |
| Exit/goal | Bounded explicit reason/goal code from the owning operation; candidate WFS reasons: recipes exhausted, agreement reached, explicit departure before completion. |
| Timing/order | Occurrence timestamp, optional source sequence and tab context; preserve standard source-observed time and add receiver receipt time. |
| Technical metadata | Service/resource, SDK/language/runtime, service version, release/build/deployment where actual SDK/detector/build metadata supplies them. |

Auto-instrumented request telemetry is useful for reception/technical context but does not manufacture a goal or visit. Existing WFS frontend/server identity must not be repurposed as the random visit identifier. Missing recommended fields disable only dependent interpretation and remain explained in event discovery.

### E. Journey lifecycle and deterministic classification

Each definition version owns event mappings and the active goal. Automatic boundaries are settled: first mapped event starts, selected active goal ends, next genuine mapped occurrence starts another, and 30 minutes without mapped activity closes inactive, not failed. Browser scopes stay per-tab. Server-only operations use bounded logical correlation without a tab; combined client/server events use the same correctly scoped correlation while preserving sources. Severity and saved reference queries do not implicitly change boundaries. Lateness/skew limits and pairing precision remain technical defaults.

Proposed lifecycle:

```text
Unassigned evidence --valid context/start--> Active journey
Active --goal evidence--> Goal completed
Active --explicit exit--> Explicitly ended (reason retained)
Active --inactivity fallback--> Provisionally inactive / outcome unknown
Any retained state --valid late evidence--> Revised result under its policy
Exception --observation only--> Current journey continues
```

This is a model for review, not a chosen timeout. Goal completion and positive business outcome are separate: first committed vote may complete a contribution milestone while agreement remains unobserved. A negative recipe vote is still a successful recorded action. Failed requires the configured authoritative evidence; flaky requires a visible deterministic retry/error condition and denominator. Unknown/conflicting evidence stays unknown/conflicting.

WFS completion-screen observation is a visit outcome, separate from first persisted contribution, shared agreement/cause and emitted exhaustion/agreement/departure reason. Another participant can trigger it; no household/person/round join is needed. Automatic first-mapped-event/active-goal rules apply to an actual mapped screen-entry event, including initial entry into already-converged state. Supported replay of one goal occurrence is deduplicated before segmentation; instrumentation should distinguish actual entry from render/hydration noise. Without identity, uncertainty remains visible. Ten distinct committed votes are healthy, with nine successive vote transitions; post-commit publication failure does not undo commitment.

Browser journeys use separate random per-tab correlation; server-only workflows use their own bounded operation correlation. Combined evidence can share a logical journey ID across client/server through existing context propagation, but trace/span/request ID, event identity and business correlation remain distinct. Correlation alone does not authorize a source, order cross-clock records or make them duplicates. Never reuse it as household/person identity or join unrelated visits/apps. Missing either source side does not block supported base OTLP intake.

### F. Ordered evidence and derived storage

Recommended logical records, not a SQL schema implementation:

- App registrations and permitted producer bindings.
- Immutable journey definitions/mapping revisions and activation metadata.
- Sanitized evidence rows: app, source/signal, occurrence/observed/receipt times, effective expiry, optional source ID, visit/journey/action/trace/span context, step candidates, app/service versions and bounded allowed attributes.
- Sequence memberships: ordered evidence references per app/visit/journey/definition, source ordering quality and provisional/complete status.
- Rebuildable path/edge/outcome/exception read models tagged by definition and computation revision.
- Separate viewer sessions, invitation grants and minimal bounded deletion markers.

Index known query dimensions; do not scan arbitrary JSON for every filter. Store sequence rows, not only pairwise edge counts or one growing session JSON document. Display ordering uses reliable source sequence where available, otherwise event time with a stable receiver tie-break and visible uncertainty. Event times from unrelated clocks cannot establish causal order. Unknown or skewed timing yields unavailable duration rather than a negative value.

Deduplication keys include app, trusted configured producer namespace, signal and supported source identity. Equal browser/server IDs coexist; cross-source action correlation is not deduplication or authority. Span identity follows a qualified trace/span policy within that namespace. Logs without reliable IDs cannot distinguish retries from legitimate identical repeats. Do not fingerprint away healthy votes. Browser intent and server commitment are different records for one supported action, not two votes. A log and span event can describe one exception, but a grouping signature cannot prove one occurrence; show diagnostic-record counts and qualify occurrence counts separately. Metric arithmetic remains unavailable until identity/temporality survives projection. Expose uncertainty rather than claiming exactly-once semantics.

Direct intake sanitizes and prevalidates the entire permitted projection, then commits all accepted evidence/contribution invalidations atomically before success. A permanent record error commits zero records; privacy omission and missing optional context do not alone make valid telemetry erroneous. A bounded writer queue can serialize work only while ack remains synchronous to commit. Specify process/OS/power-loss durability separately: a restart test proves only its tested boundary. If power-loss durability is selected for SQLite, use an appropriate synchronous setting such as FULL and documented storage/fsync assumptions; do not weaken it to improve a benchmark. SSE must not hold a DB read transaction for a connection lifetime.

Snapshot data and cursor represent one consistent committed revision. After-commit updates include lifecycle invalidations: a bounded in-process timer or query-driven mechanism must advance timeout/expiry/deletion state even without new telemetry. Define a freshness bound, replay/reset after restart or missed notification, and subscription identity so in-flight updates from old filters cannot contaminate a new view. A separate scheduler service is not required. Mark stale/recomputing state rather than mixing revisions invisibly.

Start with bounded indexed queries and small caches. Add transactional contribution tables/counters only when actual query measurements demand them. A late event, duplicate resolution, expiry or deletion invalidates/rebuilds the affected journey and read-model contribution; blind additive counters cannot do this correctly. No generic arbitrary-path graph query engine is warranted.

### G. Windows, comparisons, rare paths and stopping evidence

Historical filters identify app, app version, explicit timezone/date range and definition version. Recommend UTC storage with displayed timezone; month/day buckets must name their timezone and calendar boundaries. A session spanning app deployments keeps per-event/service version evidence; select a declared version attribution rule and show mixed/unknown versions. Do not assign a frontend app release from an unrelated API version.

Each query names its unit and selection policy: start-cohort eligibility versus event-window clipping, display-only step hiding, follow-up/censoring and version attribution. Hiding X from A-X-B may produce a labeled projection, never an observed direct A-B transition. Expired/missing/out-of-window starts remain partial even if a recent goal survives. Do not splice a mixed-version visit into fictitious paths. Repeated-step stopping rates declare occurrence versus journey units; ten votes are not ten journeys. Comparisons require compatible eligibility and follow-up and show numerator/denominator/support/coverage; a rate change is not release causality. Select the initial policy profile before the historical slice.

Rare-path thresholds are visible settings and low-support results remain counts/examples, not automatic bugs. Graph display limits use an explicit “other” summary plus path-list drilldown; rare evidence is not silently excluded. Default stopping view distinguishes recent pending/censored journeys, provisional inactivity, explicit exits, observed failure and later recovery. “No further activity observed” is acceptable; “user tired/abandoned because of this error” is unsupported.

### G1. Time at stage, without cognitive inference

The user accepted **time at stage**, node size for duration, edge thickness for traffic, median duration plus slower cases, and safeguards around active/background/missing timing. The statistic describes observed timing, not measured thought or cognitive effort. Viewing a recipe and deciding/voting on one item are separate stage semantics; repeated occurrences of a vote step remain a healthy loop. A large node is not automatically friction or failure.

Recommend explicit stage-entry/exit or start/end events through the existing standard SDK, optionally carrying safe random occurrence/tab linkage. A qualified semantic stage span can be another carrier, but technical request latency is not user-stage dwell, and sampled traces limit coverage. No mandatory Sentrame SDK, custom transport, recipe/person identity or new universal ingestion fields are introduced. Pair each occurrence within app/logical-correlation/source scope, with tab where applicable; stable supported record identity makes replay idempotent. Do not collapse consecutive occurrences of the same semantic step into one long occurrence or double-count intent and commitment.

Retain timing basis and quality alongside the evidence: measured elapsed start-to-end, observed-active intervals only where visibility/idle/tab evidence exists, or explicitly estimated next-event fallback under a named adjacency/window policy. Plain elapsed can include hidden/background/idle time. Exact activity/idle and fallback rules are unchosen. Never infer attention from a missing visibility event. Overlapping/nested stages retain distinct occurrence scope and cannot be summed into a non-overlapping visit total by default.

An end received before its start can reconcile by valid source evidence; incompatible clocks, jumps or ambiguous overlapping pairs yield unavailable/qualified timing, not fabricated negative duration. An absent end is ongoing/right-censored, not zero. Query filtering must not erase a true middle/exit event and then manufacture a fallback duration. Durations across window/version boundaries require the same declared cohort/follow-up policy used by comparisons.

Show median plus slower cases with eligible sample/completeness and separate timing kinds; do not silently pool estimates with measurements or insert ongoing occurrences into completed percentiles. A slower percentile such as p90 and a minimum-support rule are recommendations still to choose. Low-support cases can remain individual observations. Historical/version comparisons show censoring, source/sampling/activity policy and shared visual scale; differences do not establish release causality or more thinking.

Use a bounded, labeled duration scale, recommending node area rather than an unexplained radius transform. Exact scale/bounds are not chosen. Show units/statistic, clipped extremes, missing-data state and numeric/table equivalents; missing duration is not encoded as fast. In live mode show provisional elapsed/observed-active time so far with latency/clock uncertainty, under bounded refresh. Closing or late-correcting an occurrence adds at most one qualified completed contribution. Exporter/tunnel delay is not stage dwell. No timing calculation or rendering is implemented here.

### H. Live views and consistency

SSE is same-origin and cookie-authorized. Each connection snapshots app/visit/filter scope and grant. Snapshot + monotonic committed update cursor prevents gap/duplicate ambiguity; reconnect replays authorized updates within a bounded window or emits a reset requiring a snapshot. Cursor reuse cannot widen authorization, and revoked sessions close streams. Slow consumers have bounded buffers and disconnect/reset instead of blocking ingestion.

The established voice decision is a toggle between a clean live window and the selected app's accumulated traffic over the last seven days, with the chosen live anonymous visit highlighted over it. At observation start T0, app-baseline time is `[T0-7 days,T0)` and live occurrence time is `[T0,now]`. The baseline does not inherit the isolated visit filter. Compatible app/version/definition filters remain common and visible; baseline and live counts stay separate. App-history permission is independent: a visit-only invite gains no app baseline merely by toggling it. Invitation permission, earlier visit-detail visibility and access after the visit ends still require explicit grants. Bounds freeze while late evidence may correct values. Propose event-time counts with receipt-time streaming; pre-T0 late data updates history rather than live action totals. Mark timestamp fallback and avoid double counting. Same-visit history is not an unresolved alternative to this confirmed baseline.

Changing app/visit preserves applicable controls and shows incompatible filters as inactive or explicitly reset, never a silently broadened grant. Pausing rendering does not pause collection. Display freshness distinguishes source batching, receipt delay, connection state and paused motion. No promise of instant clicks.

### I. Exception context and accessible UI

Attach only through verified app/visit/action or trace/span relationships. A trace relationship without visit mapping can support technical context but not an invented visit. Group recurrence using a versioned sanitized fingerprint such as type + normalized frame signature; dynamic messages and minified stacks need explicit grouping limitations. Count occurrences separately from affected visits; unlinked records remain separate. Original source mapping is unavailable unless actual source maps/release artifacts are supplied under a later design.

Detail/copy exposes sanitized allowed type/message/stack/runtime/service/release/context, linkage basis, later activity/goals and known gaps. No confidence badge/AI diagnosis is necessary. A copy action is a human-mediated export, so permission and privacy apply to that path too. Do not send it to another person or coding service automatically.

Use self-explanatory controls for app/version/date, historical/live, all-traffic/visit, baseline on/off, step mapping and goal conditions. Offer table/sequence equivalents for canvas selection/counts/diagnostics. Preserve focus on updates; do not announce every streamed event to assistive technology. Review pause/reduced motion, keyboard, focus, contrast, error recovery and English/French task parity against the cited proposed WCAG 2.2 AA target. Preserve the first MVP’s warm-paper/forest-green theme and linked-thread Sentrame icon, per the owner’s 2026-10-28 instruction and the [visual reference](../../../docs/design/visual-theme/README.md). State must not be encoded by color alone. No visual mockup/application is produced here.

### J. Viewer access options

Release A uses owner passphrase access. Invitations are deferred, with their narrow design selected: read-only one live anonymous session, owner explicitly includes/excludes its seven-day app baseline and exception details, no configuration edits, expiry and revocation. Broad read-only app-workspace/history/diagnostic sharing stays retained roadmap B. Exact expiry/use defaults remain to review. Permission checks govern queries, saved views, replay and detail disclosure.

Recommend a high-entropy one-use, short-lived bearer redeemed into a revocable server-side session with HttpOnly/Secure/SameSite cookie, explicit allowed local redirect and bounded login attempts. Hash stored invite/session secrets as appropriate; do not persist the passphrase in telemetry. Use a maintained password verifier/KDF appropriate to the selected auth design rather than WFS's direct reusable secret comparison by analogy. Exact expiry/idle/lifetime and invite use count remain technical defaults; illustrative “15-minute redemption / one-hour watch” values are not requirements.

A reusable invite is an alternative for group viewing, but scope/use count/revocation must be explicit. Atomically consume one-use bearers; simultaneous redemption cannot create two grants. Define an explicit CSRF/origin policy for cookie-authenticated mutations, prohibit destructive GETs and account for untrusted sibling Home Server origins and host-scoped cookies. Remove redeemed bearers from URLs/logs/referrers and copied view state. Check grants on queries, replay and active streams, including rotation/revocation. Render even scrubbed telemetry as inert untrusted text. Diagnostic detail permission is the disclosure boundary: disabling the copy button cannot prevent manual copying of content already shown. Shared passphrase access cannot establish per-person accountability.

### K. Privacy, 30-day detail and cleanup

Apply an explicit allowed projection before persistence/SSE/diagnostics. Prefer bounded safe event names/codes and controlled technical attributes; unknown raw bodies, URL queries, headers, request content, household/person/recipe identifiers and secrets are not retained for “future analysis.” Available stacks/messages must be scrubbed; scrubbers are imperfect, so choose what diagnostic content is permitted before real collection. Raw intake buffers are transient/bounded and never payload-logged. Disabling collection and deleting history are distinct actions.

Recommended lifetime: valid occurrence time + 30 days; absent/unusable time uses first receipt for a known source identity and marks fallback. Reject stale data, bound future skew and quarantine/reject impossible times under an explicit policy rather than giving arbitrary future timestamps unlimited retention. Exact lateness/skew/cleanup limits are activation decisions. Unknown-ID missing-time retransmissions cannot be recognized perfectly; document this limit rather than claiming exact duplicate or original-age detection.

Queries enforce expiry immediately; bounded lifecycle invalidation refreshes open views without new intake. Physical cleanup removes sequence memberships/contributions and invalidates replay/cache paths. SQLite free pages/WAL/backups need separate handling before physical-erasure claims. Initial targeted deletion is app or identifiable visit. Journey-only deletion requires stable replay membership, not a mutable inferred journey ID; otherwise explain and offer an explicitly selected broader scope. Tombstones retain minimal identity/scope/expiry through the accepted replay horizon and survive rule/boundary revision. Unknown-ID/missing-time replay cannot support an unlimited original-age guarantee. Whole-app deletion disables intake; re-registration needs a new generation or equivalent stale-queue/credential barrier. Missing visit context can only be deleted by an actually known scope.

The selected retention is 30 days of detail and 12 months of aggregate counts/conversion/timing history. After detail expires, no individual journey/exception/ordered-path drilldown is available. Store sufficient mergeable numerators/denominators, coverage and duration distributions at declared safe grains/dimensions; never average medians. A bounded histogram or mergeable quantile structure is a technical option, with precision/approximation disclosed. Summaries retain semantic/timing versions and cannot satisfy arbitrary new filters or replay new definitions after source expiry. Keep no hidden long-lived individual membership merely to promise selective subtraction: correct identifiable contributions within the documented horizon; afterward offer explicit app/time-bucket deletion where necessary. Twelve-month expiry applies to summary caches too. Backup/restore policy remains a technical choice, not permission to retain detail longer.

### L. Minimal engineering harness proposal

The [harness proposal](../../../docs/proposals/sentrame-engineering-harness.md) adapts actual WFS patterns: contract before code, explicit observed red/green/refactor, one thin tracer bullet through affected real seams, focused discovery and proportionate verification. It does not copy WFS's large Taskfile, Kiota/client/mock stack, process-killing gates or multi-language fixtures into this empty repository.

Future task entrypoints should be a handful of explicit commands for spec validation, focused behavior tests, isolated seam tests, static checks and a controlled synthetic scenario. Define public app APIs before implementation; use official protocol definitions for OTLP. Keep test doubles at external boundaries; a mocked frontend API does not prove durable receiver/storage behavior. Preserve unrelated work and report passed/failed/blocked/not-run separately.

## Risks / Trade-offs

- [Semantics cannot be obtained from standard HTTP spans] → Keep reception useful but unassigned; require optional sufficient context for business claims. Stop at G1 if this is not feasible.
- [Low volume yields a compelling but misleading map] → Show support/counts and bounded examples; require a real decision at G0.
- [All accepted technical signals exhaust disk] → Use owner-approved retention projection/quotas and volume model; do not silently discard accepted data or describe loss as zero activity.
- [Receiver reuse brings large build/runtime dependencies] → Compare embedded versus official sidecar using actual qualification; retain protocol compatibility whichever wins.
- [Counter drift under retries/late data/deletion] → Immutable evidence, known identity, revisioned derivation and rebuildable contributions; admit unknown-ID limitations.
- [SQLite read scans starve checkpoint/ingest] → Bound query work, local-volume validation and measured G3; PostgreSQL is the first operational fallback.
- [Viewer share leaks diagnostics/history] → Server-side scoped grants and explicit copy/baseline permissions; revoke active SSE as well as cookies.
- [Full stacks compromise minimization] → Allowlisted technical projection, visible gaps and narrowly chosen diagnostics; do not claim universal scrubbing.
- [Changing definition retroactively changes comparisons] → Pin definition/recomputation revision and request reinterpretation explicitly; do not mutate historical meanings invisibly.

## Premortem

| Failure six months later | Earliest falsifying observation | Response / stop gate |
| --- | --- | --- |
| Beautiful graph, no decision value | Reviewers cannot say what to change from the four WFS examples. | Stop at G0; prefer periodic counts/evidence. |
| Healthy vote loops labeled broken | First fixture with ten votes and agreement becomes failed/flaky. | Block the rule; preserve business goal and loop semantics before expansion. |
| Ingest works, no real journeys | All received WFS signals lack visit/goal linkage. | Stop at G1; optional standard instrumentation design before backend scope growth. |
| Home Server runs hot/fills disk | 30-day projected detail exceeds owner budget or ingest degrades queries/cleanup. | Narrow allowed retained telemetry, use existing capacity or stop at G3. |
| Retry/late corrections destroy trust | Counts change on replay without explanation, or deleted visits reappear. | Block expansion; qualify identities, corrections and deletion through the same read API. |
| Invite behaves like admin access | Restricted viewer reaches another app/history/diagnostic copy or stream survives revocation. | Block viewer release; explicit scope tests before collaboration. |
| Error context cannot explain recovery | A post-commit error ends the path, hides later agreement or is joined by timestamp proximity. | Block classification/linkage; use verified action/trace context and continued-activity cases. |
| Maintenance becomes an observability project | New requests focus on generic metric explorers/trace query engines rather than WFS decisions. | Revisit G0/G2 and integrate an existing observability tool. |

## Migration Plan

There is no deployed service/data to migrate. The approved documentation reconciliation replaces the initial proposal/design/tasks and obsolete deltas in the existing change. Publication is a documentation-only draft PR; it performs no canonical sync/archive, app/harness setup, implementation, merge or deployment. Every implementation task remains unchecked.

After a future approved implementation decision, start with a single synthetic standard-exporter tracer bullet and temporary local storage; qualify protocol, ordering/replay, unknown outcomes and expiry before live WFS integration. Real telemetry, deployment, production testing and GitHub publication remain separate explicit scope. Rollback plans must be concrete at that later stage: stop service, revoke intake/viewer grants, preserve or delete data under the chosen lifecycle, and reverse a backward-compatible schema/config change only with verified data handling.

## Open Questions

The main product choices are settled. Remaining implementation-level defaults are receiver packaging/pins, ingress and lateness/skew bounds, durability failure model, correlation/clock validation, summary grains/quantile precision/correction horizon, stage activity/fallback/slow-percentile policy, invitation expiry/use and deployment origin settings. Present these together for a focused design review before code; do not reopen release scope, automatic boundaries, selected stack or retention.


### Saved reference views and operational goals

Save safe query definitions for goal-backward paths, beginning/end-event segments and origin-forward exploration, including app/filter/definition version and declared time semantics. Saved views neither preserve raw data beyond retention nor alter automatic capture boundaries. A retrospective goal differing from the active operational goal queries existing segments; if endpoints cross goal/timeout boundaries, render those boundaries rather than inventing a continuous captured journey. Summary-only dates support only their retained dimensions/statistics, not full path reconstruction.

### External integration-agent instructions

The public integration guide is separate from the internal engineering harness. It discovers the target app and existing OTel, asks developers focused success/friction/failure questions, proposes minimal prioritized event and new warning/error/recovery instrumentation with evidence/value/authority/privacy/noise/tests, and obtains event-plan authorization before edits. It provides no AI runtime, proprietary client abstraction or untested runnable SDK recipe. Client-only, server-only and combined correlation, stage timing, healthy loops, safe diagnostics and synthetic support-matrix qualification are covered.

### Owner-approved canvas workspace revision — 2026-10-28

The owner explicitly approved implementation of the graph-centered workspace, superseding duration-sized nodes. Three shared-context views are horizontal hierarchy, directed settled-force layout and cycle-capable Sankey. Warm paper, forest branding, linked-thread icon, serif headings and visible focus remain. Traffic controls line/ribbon weight and Sankey bars; timing, measurement quality and supported median/p90 stay in the inspector and tables. Arrowheads stay approximately six screen pixels independently of weight/zoom. Skip-stage edges route around intervening nodes; healthy loops and recovery remain visible.

Selection exposes incoming/outgoing neighbors and ordered evidence. Event mapping opens a native modal from the node, previews real records and prospectively activates a full copied definition. Operational-goal change is explicit. Unmapped discovery nodes remain disconnected absent supported relationships. Rare/ending/exception/recovery lenses use actual captured journeys, preserve gaps and explain sample support. The rare lens is <=10% per full sequence with >=20 closed, non-partial/uncensored comparable journeys per definition/version; individual highlighted links can also occur in common paths.

Reference exploration is contextual and can precede saving. Overview, goal-backward, origin-forward and segment references preserve bounds and safe filters plus layout/lens. Live references freeze the observation window; the seven-day baseline remains a comparison layer with separate counts. Saved reference IDs default their stored filters; explicit empty query filters clear those defaults. No view changes capture boundaries or outlives retention. A removable visit chip clears view scope without deleting evidence. Tables/ordered sequences remain first-class alternatives. Lifecycle/deletion/collection health reside in workspace menus, outside contextual event controls.

### Owner-approved diagnostic detail and discovery — 2026-10-28

The owner approved bounded scrubbed messages/stacks and custom standard error log bodies for application-wide debugging. `scrubbed-v1` replaces blanket diagnostic-text omission for new receipts; the concrete byte/line/redaction/query limits and residual risk are in the technical-default revision. Errors are searchable without a visit; valid correlation exposes related context rather than a person or cause. Failed spans, standard exception attributes and error log bodies use the same sanitized storage/read/copy path. Older records are never backfilled; replay comparison preserves the original receipt projection. Source-map reconstruction and real .NET/Next exporter qualification remain unclaimed. Diagnostic text expires with 30-day detail and never enters 12-month summaries.

### Owner-approved node arrangement — 2026-10-28
Drag placement is local to the mounted graph, keyed by layout and semantic node ID. Force allows both axes; hierarchy/Sankey preserve column x coordinates and allow vertical adjustment within reserved routing bounds. The graph reroutes observed edges and independent comparison layers from moved node positions without rerunning the force simulation per pointer event or altering volume/evidence. A five-screen-pixel threshold separates clicks from drags; pointer capture completes drags outside the node, cancellation restores the prior placement. Shift-arrow movement and a quiet Reset layout control provide keyboard access and recovery. Reloading clears manual positions; saved reference APIs and stored evidence remain unchanged.

### Authorized app-wide viewers and container delivery — 2026-10-28
The owner explicitly brought app-wide temporary read-only invitations forward from roadmap B. These are separate from the retained narrower visit-only grants. Invitations last1h,24h by default, or7d from creation and are single-redemption. Store only a bearer digest and immutable non-reused grant ID bound to app registration; use a distinct viewer cookie, strict server read allowlist and persistent expiry/revocation checks per request and SSE tick. Owner capabilities stay separate. Sessions are in-memory and end on restart; consumed invitations remain consumed. Fragments are stripped before API activity, and redemption is explicit rather than automatic. A viewer can see all retained app evidence and sanitized diagnostics, but cannot mutate shared data/configuration or administer invitations. UI session checks clear mounted evidence on access loss even in historical/paused modes; previously copied data cannot be erased remotely.

The owner also authorized preparing a Docker Compose deployment and DockerHub publishing workflow for brisebois/sentrame, with linux/amd64 and linux/arm64 under one image tag. Compose and an environment example are required repository artifacts. Existing no-commit/no-push instructions remain until explicitly changed. Container build, registry publication, Home Server fit and shared Tunnel route qualification are independent checks; do not label unrun cells as passed. No WFS emitter edits or actual export are implied merely by packaging Sentrame.
