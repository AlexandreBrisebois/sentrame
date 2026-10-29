# Tasks

All tasks below are future implementation, not work authorized by this planning pass. Read [design](design.md), [cutover/reconciliation](compatibility.md), [documentation matrix](documentation-plan.md) and [qualification](qualification.md). Use test-first vertical slices: record a meaningful behavior failure before implementation, then focused green evidence; tooling failure is blocked, not red. No legacy compatibility layer, production collection, external-app edit, deployment, publishing, commit or archival is implied. Roadmap B remains deferred.

## 1. Standard evidence projection and durable receipt

- [x] 1.1 Add failing isolated receiver→SQLite→protected-read tests and readable standard OTLP fixture sources for native end/parent/kind/status/flags, span events/log correlation and bounded links, without custom journey fields (Q01–Q04). Verify failures are missing behavior rather than setup errors; extend the existing fixture generator, not a parallel harness.
- [x] 1.2 Implement bounded structural projection, exact nanosecond serialization, safe HTTP selectors, explicit schema/projection identity and native span/event replay/conflict behavior (Q09/Q12/Q14). Verify red tests pass, forbidden sentinels never persist, >32 links disclose omissions, direct-request conflicts commit zero and unsupported old schemas fail without reset. Reuse current receiver/auth/request limits; do not build legacy readers.
- [x] 1.3 Update `CONTRACT.md`, producer contract and decision register for this slice, preserving concurrent edits and recording greenfield clarification. Verify retained fields match passing tests and HTTP route-template privacy rules, all proposed-but-unimplemented features remain labeled, and original receipt/auth/metrics behavior tests pass.

## 2. Protected trace exploration and first value

- [x] 2.1 Add failing structural/API/UI tests for parallel/nested/repeated operations, missing/ambiguous parents, out-of-order span events, exact versus trace-only log context, links to absent/other traces and sampling/truncation (Q01–Q04/Q13). Verify no chronological sibling edge, person count or fictitious journey is expected.
- [x] 2.2 Implement bounded protected trace list/detail reads, typed relationship model, native inclusive timing and tree/timeline/table UI with revision invalidation. Verify focused tests and actual isolated protected reads, viewer allowlist/revocation, missing/invalid end exclusions, exact time strings and no cross-source authority promotion; avoid critical-path/exclusive-time scope.
- [x] 2.3 Add literal `rest/standard-otel.http` login/ingest/trace-read/unauthorized/replay requests and verifier support through actual receiver/storage/read paths; verify Q01–Q04/Q13 and restart persistence using the same bytes. Update README, user guide, local demo guide and explicit fixture Taskfile tasks with available trace behavior and limits. Record wire results separately from exporters.

## 3. Configured business interpretation using standard evidence

- [x] 3.1 Add failing mapping/derivation/read/UI tests for service+carrier+name+source+supported-condition selectors, overlapping rules, native-span stages, semantic goal authority and end-boundary evaluation (Q03/Q05/Q07). Include ordinary telemetry without `stage.name`, browser intent and HTTP 2xx that cannot prove persistence.
- [x] 3.2 Implement contextual discovery/preview/activation and trace business overlay, native occurrence identity/timing basis, explicit paired/estimated alternatives and unambiguous goal selectors. Verify tests pass, no synthetic boundary rows/custom occurrence IDs for spans, span duration does not imply dwell, and unsupported fields remain unavailable rather than rehydrated from raw payloads.
- [x] 3.3 Update developer/producer/user/internal contracts for ordinary OTel mapping and timing. Verify UI labels/table bases, named-event/operation semantics and goal boundary examples match Q03/Q05/Q07; extend literal REST assertions to prove configured outcome evidence rather than merely accepted export.

## 4. Optional stage.name enrichment

- [x] 4.1 Add failing tests for stage hint suggestions, absent/empty/invalid values, multiple labels, explicit owner override, overlapping-rule rejection, no child inheritance and no tag-implied Kind/goal/authority/session (Q06/Q12). Verify the tests distinguish a suggestion from an activated literal mapping.
- [x] 4.2 Implement bounded suggestions and explicit acceptance into prospective definitions. Verify owner **Vote saved** wins over emitted `voting`, missing hints preserve mapping, two explicit matches block activation and emitted label changes do not rename prior captures.
- [x] 4.3 Extend voting REST fixtures/read assertions with standard and enriched variants, rule conflict, accepted suggestion and historical-definition/reference checks (Q06/Q10). Verify exact replay yields unchanged graph/timing counts and rules affect only intended future captures.
- [x] 4.4 Update both .NET and Next.js recipes with stable `stage.name` naming, existing-span and named-log additions, producer-contract precedence table, user preview behavior and agent guidance. Verify examples request no new SDK/provider/wrapper and clearly label optional application convention versus standard OTel; retain source-only status until actual exporter checks pass.

## 5. Richer journeys, presence and forward history

- [x] 5.1 Add failing tests for explicit multi-request scope/tab, server-only scope, absent/ambiguous context, goal/inactivity/late evidence, captured definition/reference modes, optional page pair/presence intervals, expiry/deletion and identity-free summaries (Q07–Q14). Verify no IP/agent/time/fingerprint inference or historical remapping expectations.
- [x] 5.2 Integrate span-backed occurrences with existing scoped journey lifecycle, explicit timing basis, trace-versus-journey references and summary policy separation. Verify late/replay revisions, no duplicate pair/span sample, no hidden-gap adjacency, no negative/missing completed samples, no expired-link resurrection and no new long-lived identity index. Preserve owner/read-only/CSRF boundaries.
- [x] 5.3 Expand literal REST voting story through navigation/intent/failure/retry/trusted commit/recovery/goal with explicit scope only in this increment; add optional actual visibility/idle observations. Verify Q07–Q14 through durable reads and update user/producer/agent/demo docs with what each richer signal enables and cannot establish.

## 6. Real recipe qualification and deployment documentation

- [x] 6.1 In a later authorized disposable .NET sample, pin compatible runtime/packages and reuse one provider; exercise standard logs/spans, EventId names, exceptions/custom safe messages, optional stage enrichment, correlated context, batching and shutdown. Verify actual exporter→receiver→durable protected read evidence and record exact protocol/compression cells; update the .NET recipe to runtime-qualified only for passed cells.
- [x] 6.2 In a later authorized disposable Next.js Node sample, pin framework/Node/SDK and reuse one initialization path; exercise real route/server operation/errors, named logs, optional tag and flush/shutdown. Verify actual exporter reads and cross-service W3C parent/log relationships where used; keep Edge/browser separate and update the Next recipe only for passed cells.
- [x] 6.3 Under separately authorized isolated browser qualification, verify reviewed app-origin OTLP forwarding, browser namespace, no server credentials in bundles, CSP/CORS, collection choice, lifecycle loss and optional presence. Walk the agent checklist on the sample, verifying minimum changes and no duplicate SDK/proprietary wrapper; report browser and agent evidence separately from server checks.
- [x] 6.4 Reconcile the complete documentation matrix: README, developer/producer/agent/user/internal guides, index, local review, NAS/DockerHub and deployment-example comments as needed. Add this change to explicit spec validation tasks because the existing Taskfile still targets the now-archived baseline. Verify local links/fences, exact endpoint/protocol/auth/token guidance, consistent voting names and current/proposed/source/runtime labels; never print secret env values or claim unrun routes.

## 7. Integrated acceptance and evidence delivery

- [x] 7.1 Run proportionate focused suites followed by `task test`, strict validation of this change and canonical specs, whitespace/link checks and the literal REST verifier on a fresh isolated service. Populate `docs/verification/standard-otel.md` with red/green history, Q01–Q14 statuses, fixture hashes, exact commands, durable restart reads and failed/blocked/not-run limits. Static checks alone cannot close runtime cells.
- [x] 7.2 When local container execution is authorized and available, qualify Linux ARM64 and AMD64 separately with the new model, persisted bind storage/read/restart and protected trace behavior. Verify results per platform and record native versus emulated execution; leave unexecuted cells not-run, without publishing an image.
- [x] 7.3 Review the new-model qualification matrix and developer workflow end to end: standard export first, optional hints second, richer context only where required. Verify that Home Server/Tunnel/Collector/gRPC/external-app cells remain explicitly unqualified unless separately authorized and exercised; record owner-supplied budgets as prerequisites rather than claiming infrastructure support.

## Workflow follow-up

- Review the trace-tree/path separation, suggestion activation and operation-only timing defaults before apply. The owner has already resolved backward compatibility: none is required.
- A later separately authorized canonical sync/archive must retain unrelated canonical requirements and roadmap B, add these unique deltas and apply the diagnostic replacement that removes only the obsolete legacy-data gate. No archive, publish or deployment task is part of this checklist.
