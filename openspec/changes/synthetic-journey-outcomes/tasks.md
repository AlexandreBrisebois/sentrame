# Tasks

Implementation roadmap only. Every task is unstarted and requires a subsequent implementation request. Release A is the complete useful WFS workflow in [release scope](../../../docs/planning/sentrame-release-scope.md); B remains retained scope. A thin tracer proves a seam and is not the final release. The five main product decisions are settled; remaining technical defaults receive one focused design review.

## 1. Review technical defaults before implementation

- [ ] 1.1 Review concrete event names/active goal, safe logical and occurrence context, trusted producer namespaces, consent/allowlist and WFS emitter boundaries against the accepted completion-screen/healthy-loop cases. Do not change external apps under this planning task.
- [ ] 1.2 Propose and review receiver packaging/pins, request/queue/decompression limits, lateness/skew, failure/durability model, timing activity/fallback/slow-percentile rules and aggregate grains/distributions/correction/deletion horizons. Preserve the selected stack and product policies.
- [ ] 1.3 Record Home Server architecture/filesystem, workload, shared capacity and numerical resource/build limits. RAM/CPU availability is reported; qualify actual fit under later authorized G3b/G3c work rather than inventing limits or reopening the stack ballot.

## 2. First thin end-to-end tracer

- [ ] 2.1 Observe a meaningful failing behavior assertion for one pinned standard log exporter → app binding → sanitized durable store → protected owner read → accessible sequence row. Implement minimally, observe green through the real store and after restart, then refactor. Qualify only that compatibility cell.
- [ ] 2.2 Extend the same path for standard records missing optional journey/version context, unknown metadata and unsafe extra attributes. Verify accepted-unassigned receipt without custom identifiers and no raw payload leakage to storage/logs/SSE/copy.
- [ ] 2.3 Verify whole-direct-request prevalidation/atomic rejection, explicit permanent versus transient status, failed commit, bounded overload and namespaced replay. Distinguish direct-storage durability from upstream Collector enqueue/rebatch acknowledgement and declare unknown-ID limitations.
- [ ] 2.4 Measure this authorized isolated tracer and cold/warm build/test costs against the agreed G3a envelope. Record G3b passed/failed/blocked before expansion; this does not establish integrated deployment fit.

## 3. Release A journey workflow

- [ ] 3.1 Add discovered events, simple mappings/preview and active goal configuration. Verify client-only, server-only and combined evidence, first mapped event, selected goal, next-new-occurrence restart, per-tab browser separation and 30-minute inactive-not-failed closure including no-new-input transitions.
- [ ] 3.2 Verify healthy repeated votes, distinct commitment/intent/completion-screen/reason semantics, other-participant convergence, initial convergence, supported goal replay, exhaustion/departure, silence and post-commit publication errors. No person/household joins or invented cause.
- [ ] 3.3 Add warning/exception context and recovery through eventual completion. Verify linked versus nearby-unlinked evidence, one occurrence/two carriers, two real same-stack failures, affected logical scopes, inert scrubbed detail/copy and no severity-based termination.
- [ ] 3.4 Implement stage-occurrence timing and historical/live graph encodings: node size for time at stage, edge thickness for traffic, median plus supported slower cases. Verify active/background/idle/missing/clock/overlap cases, repeated occurrences, estimated versus measured basis, sample/completeness, bounded legend and numeric/table alternatives.
- [ ] 3.5 Add historical maps, default potential dead ends, path/loop/outlier counts and simple filters. Verify cohorts/follow-up/units, hidden X in A-X-B, missing/expired beginnings and equal-edge/different-sequence paths without invented adjacency or bug labels.
- [ ] 3.6 Add saved goal-backward, bounded beginning/end and origin-forward reference views over recorded boundaries. Verify repeated endpoints, unavailable starts, cross-boundary markers, definition/filter persistence and no expired-detail resurrection or query-driven resegmentation.
- [ ] 3.7 Add consistent-revision snapshot/cursor and bounded SSE with filter-preserving scope switches. Verify reconnect/restart/missed notification, stale filter updates, slow-viewer reset, timeout/expiry without new input and no stream-long DB transaction.
- [ ] 3.8 Add the clean-live versus fixed selected-app seven-day baseline toggle, highlighting the chosen live visit. Verify compatible app/version/definition filters, fixed baseline bounds, separate live/baseline counts and provisional ongoing timing. Qualify routed SSE/HTTP behavior separately from origin OTLP; public Tunnel gRPC is not assumed.
- [ ] 3.9 Implement 30-day detail and 12-month summary lifecycle using supported mergeable counts/denominators/timing distributions, never averaged medians. Verify late/replay/rule corrections, bounded physical cleanup, stable deletion scopes/tombstones/re-registration, summary expiry/privacy and disabled individual drilldown after detail expiry.
- [ ] 3.10 Qualify release A intake cells using pinned standard producers, including HTTP content/encoding/gzip/error behavior and supported traces/logs/metrics receipt. Advertise tested cells only; metrics arithmetic and profiles are outside the initial analysis. Preserve app/version identity without building the multi-app experience upfront.
- [ ] 3.11 Review the external integration-agent guide with an approved synthetic app case: discovery-led questions, prioritized event plus new warning/error/recovery plan, developer approval before edits, standard SDK context, privacy/noise and honest support validation. Do not claim the illustrative WFS plan is deployed integration.
- [ ] 3.12 Under later authorized isolated integrated scope, qualify concurrent ingest/query/cleanup/live/browser operation, declared crash boundary, retention/storage growth, resource/build/test costs and keyboard/focus/language recovery. Record G3c evidence and revisit cheaper adaptation/value before any deployment.

## 4. Retained roadmap B

- [ ] 4.1 Extend stable registration to the multi-app experience and richer normalized version comparisons, preserving app isolation, unknown/mixed versions, definition identity and denominators.
- [ ] 4.2 Implement the selected narrow read-only one-live-session invitations with explicit baseline/diagnostic grants, expiry/revocation and no configuration rights. Verify atomic redemption, CSRF/origin checks, credential-free URLs, inert text and active SSE revocation; choose bounded lifetime/use defaults first.
- [ ] 4.3 Review broader read-only app workspace sharing and richer configuration as separate later extensions. Do not silently widen an existing narrow grant or delete these retained requirements.

## Engineering discipline

Use one red → green → refactor vertical increment with proportionate checks; tool/prerequisite failures are blocked, not red. Reuse maintained components and existing test filtering; introduce no harness framework before needed. Proposed harness/Compose/agent instructions are documentation, not implemented files. Real WFS edits, collection, deployment, merge and canonical spec sync/archive require their own authorized work. This draft PR completes planning only.
