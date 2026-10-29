# Tasks

Implementation roadmap for the evolved journey plan. The 2026-10-28 request authorizes local synthetic implementation; checked tasks have recorded evidence. Historical planning-only authorization text describes the prior phase. Release A is the complete useful WFS workflow in [release scope](../../../docs/planning/sentrame-release-scope.md); B remains retained scope. A thin tracer proves a seam and is not the final release. The five main product decisions are settled; remaining technical defaults receive one focused design review.

## 1. Review technical defaults before implementation

- [x] 1.1 Review concrete event names/active goal, safe logical and occurrence context, trusted producer namespaces, consent/allowlist and WFS emitter boundaries against the accepted completion-screen/healthy-loop cases. Do not change external apps under this planning task.
- [x] 1.2 Propose and review receiver packaging/pins, request/queue/decompression limits, lateness/skew, failure/durability model, timing activity/fallback/slow-percentile rules and aggregate grains/distributions/correction/deletion horizons. Preserve the selected stack and product policies.
- [ ] 1.3 Record Home Server architecture/filesystem, workload, shared capacity and numerical resource/build limits. RAM/CPU availability is reported; qualify actual fit under later authorized G3b/G3c work rather than inventing limits or reopening the stack ballot.

## 2. First thin end-to-end tracer

- [x] 2.1 Observe a meaningful failing behavior assertion for one pinned standard log exporter → app binding → sanitized durable store → protected owner read → accessible sequence row. Implement minimally, observe green through the real store and after restart, then refactor. Qualify only that compatibility cell.
- [x] 2.2 Extend the same path for standard records missing optional journey/version context, unknown metadata and unsafe extra attributes. Verify accepted-unassigned receipt without custom identifiers and no raw payload leakage to storage/logs/SSE/copy.
- [x] 2.3 Verify whole-direct-request prevalidation/atomic rejection, explicit permanent versus transient status, failed commit, bounded overload and namespaced replay. Distinguish direct-storage durability from upstream Collector enqueue/rebatch acknowledgement and declare unknown-ID limitations.
- [ ] 2.4 Measure this authorized isolated tracer and cold/warm build/test costs against the agreed G3a envelope. Record G3b passed/failed/blocked before expansion; this does not establish integrated deployment fit.

## 3. Release A journey workflow

- [x] 3.1 Add discovered events, simple mappings/preview and active goal configuration. Verify client-only, server-only and combined evidence, first mapped event, selected goal, next-new-occurrence restart, per-tab browser separation and 30-minute inactive-not-failed closure including no-new-input transitions.
- [x] 3.2 Verify healthy repeated votes, distinct commitment/intent/completion-screen/reason semantics, other-participant convergence, initial convergence, supported goal replay, exhaustion/departure, silence and post-commit publication errors. No person/household joins or invented cause.
- [x] 3.3 Add warning/exception context and recovery through eventual completion. Verify linked versus nearby-unlinked evidence, one occurrence/two carriers, two real same-stack failures, affected logical scopes, inert scrubbed detail/copy and no severity-based termination.
- [x] 3.4 Implement stage-occurrence timing and historical/live graph encodings: time at stage in contextual inspector/table (duration-sized nodes superseded by the owner-approved workspace revision), edge thickness for traffic, median plus supported slower cases. Verify active/background/idle/missing/clock/overlap cases, repeated occurrences, estimated versus measured basis, sample/completeness, bounded legend and numeric/table alternatives.
- [x] 3.5 Add historical maps, default potential dead ends, path/loop/outlier counts and simple filters. Verify cohorts/follow-up/units, hidden X in A-X-B, missing/expired beginnings and equal-edge/different-sequence paths without invented adjacency or bug labels.
- [x] 3.6 Add saved goal-backward, bounded beginning/end and origin-forward reference views over recorded boundaries. Verify repeated endpoints, unavailable starts, cross-boundary markers, definition/filter persistence and no expired-detail resurrection or query-driven resegmentation.
- [x] 3.7 Add consistent-revision snapshot/cursor and bounded SSE with filter-preserving scope switches. Verify reconnect/restart/missed notification, stale filter updates, slow-viewer reset, timeout/expiry without new input and no stream-long DB transaction.
- [x] 3.8 Add the clean-live versus fixed selected-app seven-day baseline toggle, highlighting the chosen live visit. Verify compatible app/version/definition filters, fixed baseline bounds, separate live/baseline counts and provisional ongoing timing. Qualify routed SSE/HTTP behavior separately from origin OTLP; public Tunnel gRPC is not assumed.
- [x] 3.9 Implement 30-day detail and 12-month summary lifecycle using supported mergeable counts/denominators/timing distributions, never averaged medians. Verify late/replay/rule corrections, bounded physical cleanup, stable deletion scopes/tombstones/re-registration, summary expiry/privacy and disabled individual drilldown after detail expiry.
- [x] 3.10 Qualify release A intake cells using pinned standard producers, including HTTP content/encoding/gzip/error behavior and supported traces/logs/metrics receipt. Advertise tested cells only; metrics arithmetic and profiles are outside the initial analysis. Preserve app/version identity without building the multi-app experience upfront.
- [x] 3.11 Review the external integration-agent guide with an approved synthetic app case: discovery-led questions, prioritized event plus new warning/error/recovery plan, developer approval before edits, standard SDK context, privacy/noise and honest support validation. Do not claim the illustrative WFS plan is deployed integration.
- [ ] 3.12 Under later authorized isolated integrated scope, qualify concurrent ingest/query/cleanup/live/browser operation, declared crash boundary, retention/storage growth, resource/build/test costs and keyboard/focus/language recovery. Record G3c evidence and revisit cheaper adaptation/value before any deployment.

## 4. Retained roadmap B

- [ ] 4.1 Extend stable registration to the multi-app experience and richer normalized version comparisons, preserving app isolation, unknown/mixed versions, definition identity and denominators.
- [ ] 4.2 Implement the selected narrow read-only one-live-session invitations with explicit baseline/diagnostic grants, expiry/revocation and no configuration rights. Verify atomic redemption, CSRF/origin checks, credential-free URLs, inert text and active SSE revocation; choose bounded lifetime/use defaults first.
- [ ] 4.3 Review broader read-only app workspace sharing and richer configuration as separate later extensions. Do not silently widen an existing narrow grant or delete these retained requirements.


## 5. Owner-approved workspace revision (2026-10-28)

- [x] 5.1 Implement and verify the shared hierarchy/force/Sankey canvas, bounded arrowheads, routed loops/skip-stage edges, truthful separate baseline/live volumes and keyboard/table alternatives.
- [x] 5.2 Implement and verify contextual mapping preview/activation, unsaved reference exploration and persisted overview/endpoint references with captured bounds and safe filters, graph evidence highlights and removable visit selection.

- [x] 5.3 Implement and verify the protected event-volume timeline with shared two-handle slicer, contextual visit following with advanced ID search, local red error connections, and the pinch/focus/ribbon refinements.

- [x] 5.4 Implement and verify protected application-wide errors with bounded scrubbed standard exception messages/stacks and custom error bodies, privacy/expiry/deletion, original-policy replay, inert inspection/copy and repeatable actual OTLP diagnostic demonstrations.

- [x] 5.5 Implement and verify per-layout node dragging with rerouted connections, preserved evidence/volume, click-versus-drag behavior, keyboard movement and reset.

- [x] 5.6 Implement and verify explicit whole-app temporary read-only invitations, bounded lifetime/one-use redemption, server permissions, revocation/expiry/SSE, clear role-aware UI and safe fragment links. Keep narrower visit grants and richer configuration deferred.

## 6. Owner-authorized NAS delivery preparation (2026-10-28)

- [x] 6.1 Add a multi-architecture container build, Compose file, environment example and DockerHub GitHub Actions workflow for brisebois/sentrame. Verify source/configuration checks and accurately record container-runtime, registry and NAS qualification gaps; do not commit or push.

- [x] 6.2 Replace the planning-era README with the implemented product overview and add user/developer documentation, including source-checked .NET/Next.js OTLP recipes, routing/trust distinctions and explicit unrun SDK qualification. Validate local links/code fences and record published-container evidence.

## Engineering discipline

Implementation verification: [local review candidate, passed checks and retained finish gates](../../../docs/verification/release-a-local.md). The initial Go/network blocker is resolved. Task 2.4 remains unchecked because measurement was recorded after expansion; task 3.12 remains partial beyond the short local qualification. Roadmap B remains deferred.

Technical review evidence: [local release-A defaults and prerequisite status](../../../docs/review/release-a-technical-defaults.md), 2026-10-28. Tasks 1.1/1.2 record reviewed choices, not implementation or dependency/runtime qualification. Task 1.3 remains unverified because Home Server measurements are unavailable; it does not prevent the separately authorized local synthetic work. The first tracer is now qualified; the initial Go/network blocker was resolved. The current implementation request authorizes local synthetic application/harness work despite historical planning-only language below.

Use one red → green → refactor vertical increment with proportionate checks; tool/prerequisite failures are blocked, not red. Reuse maintained components and existing test filtering; introduce no harness framework before needed. Proposed harness/Compose/agent instructions are documentation, not implemented files. Real WFS edits, collection, deployment, merge and canonical spec sync/archive require their own authorized work. No publication, commit, merge, deployment or canonical sync/archive follows local checks.

Workspace tasks 5.1–5.4: [observed automated/browser/actual REST evidence and limits](../../../docs/verification/workspace-design.md). Physical trackpad/assistive-technology and actual .NET/Next exporters are not qualified. The final message-only exception predicate is source-tested and loads at the next normal restart; the demonstrated typed diagnostic cases are verified live. No original Home Server/resource or roadmap-B gate is marked complete by these checks.

App-wide invitation extension5.6 is verified by handler/SQLite and helper tests plus owner-interface browser review, as recorded in [workspace evidence](../../../docs/verification/workspace-design.md). This does not complete deferred narrow visit grants4.2 or richer configuration4.3. Container source preparation6.1 is distinct from any unrun runtime/publish/NAS gate.
