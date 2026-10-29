# Documentation change plan

This is a required implementation deliverable plan. Operational guides are not changed during this planning pass; existing edits are preserved. New behavior must stay labeled **Proposed by this change** until implemented. Each recipe must separately state **Source-checked example** or **Runtime-qualified integration**, including exact versions and route.

## Ownership and terminology

The producer-facing source of truth is `docs/integration/OTEL-CONTRACT.md`: standard carrier, optional convention, retained-field support, capability prerequisites and limitations. Capability specs govern future behavior; `CONTRACT.md` governs implemented internal/API boundaries; the developer guide owns runnable setup; the user guide owns contextual mapping; the qualification report owns support claims. Link these sources rather than duplicating slightly different attribute tables everywhere.

Use: **span/operation** for .NET Activity and JS span; **point event** for a named log or attached span event; **supporting activity** for Sentrame Kind activity; **stage** for configured business meaning; **trace** for execution context; **journey** for configured lifecycle within supported correlation; **observed-active** for explicit presence measurement. Never use “session” for a trace or “active” for generic server duration.

Status labels in every affected guide:

- **Available in current release:** supported standard OTLP/HTTP receipt, discovered names/errors, current configured correlated journeys and paired timing; cite the actual tested cell.
- **Proposed by this change:** structural trace exploration, native timing, richer selectors and tag suggestions until verified implementation replaces the label.
- **Optional application convention:** `stage.name`, `tab.id`, `stage.occurrence.id`, `activity.state`, action/record IDs where used. `session.id` is an opt-in Development OTel client convention with a separate Sentrame server-workflow interpretation.
- **Source-checked example:** inspected docs/code and locked API family only. **Runtime-qualified integration:** actual named exporter/runtime/path exercised through durable protected reads. Wire-shaped .NET-like data is not .NET SDK qualification.

## Concrete file matrix

| File | Required update | Verification / owner task |
| --- | --- | --- |
| `README.md` | Default flow: export → discover → explore trace → configure meaning → optional enrichment. Link five-level table and progressive voting example; avoid unconditional support claims. | Follow first-use links; task 2.3 and 6.4 |
| `docs/DEVELOPER-GUIDE.md` | Retain practical independent .NET/Next setup, protocol/full-path/header/token guidance; restructure sections 6–7 around levels; add both tag examples and native timing limits; one voting story replaces unrelated checkout examples. | Compile/run separately qualified recipes; tasks 3.3, 4.4, 6.1–6.4 |
| `docs/integration/OTEL-CONTRACT.md` | Primary producer contract; replace blanket paired-event timing advice with native span-first choice; five levels; supported fields versus proposed fields; stage hint precedence; stable naming; scoped correlation and presence. | Compare every field/carrier against projection tests and language examples; tasks 1.3, 3.3, 4.4, 5.3 |
| `docs/integration/AGENT-INTEGRATION.md` | Lead with reuse/discovery, then smallest evidence gap; operation versus business stage versus presence; approved emitter checklist; remove universal pair/occurrence prescriptions and stale duration-sized-node text. | Walk checklist on isolated example without wrapper or duplicate provider; tasks 4.4, 6.3 |
| `docs/USER-GUIDE.md` | Distinguish trace tree from chronological graph; choose service/carrier/name/condition; preview hint conflicts; native operation time versus dwell; capture versions and reference modes. Retain accessible tables and read-only restrictions. | Real browser/API acceptance for all five levels; tasks 2.3, 3.3, 4.4, 5.3 |
| `CONTRACT.md` | New evidence/timestamp fields and schema versions, selectors, goal boundary and timing policies; trace read routes and viewer allowlist; new-model replay/conflict rules. Remove old local-format migration promises at greenfield cutover. | Contract/API/serialization tests; tasks 1.3, 2.3, 3.3 |
| `docs/README.md` | Add proposed change beside canonical specs and archived baseline, describe precedence/current-vs-proposed status, link producer contract and qualification source. | Local-link checker and cross-document review; task 6.4 |
| `docs/planning/sentrame-decisions.md` | Append the owner’s greenfield clarification and selected trace/hint/timing defaults when applying; distinguish old planning snapshots and later approvals. | Trace decisions to this package; task 1.3 |
| `docs/integration/synthetic-demo.md`, `synthetic-event-plan.md` | Add progressive standard voting fixture, trace relationships, exact replay, failure/retry/commit/recovery goals; keep wire versus exporter distinction. | Literal REST verifier assertions; tasks 2.3, 4.4, 5.3 |
| `rest/standard-otel.http`, `rest/fixtures/standard-otel.json` (new) | Version-controlled requests/readable fixture sources with standard signals and safe placeholders; generated exact protobuf/JSON bodies and secret manifest remain ignored/private. | Same requests executed by verifier and manually usable in REST Client; tasks 1.1–5.3 |
| `docs/LOCAL-REVIEW.md`, `Taskfile.yml` | New explicit fixture/verify tasks, isolated fresh-data instructions; add validation of this change because current `task spec:validate` still names the now-archived baseline. | Run literal commands in isolated qualification; tasks 2.3, 6.4 |
| `docs/NAS.md`, `docs/DOCKERHUB.md` | Direct/container/Tunnel endpoint distinction, new capability version minimum, safe test-data cutover and exporter qualification links. Keep infrastructure results scoped to actual observations. | Documentation/config checks now; runtime each route separately later; task 6.4 |
| `.env.example`, `compose.yaml`, `compose.cloudflare.yaml` | Audit comments/examples for correct origin and sender endpoint. Add only useful non-secret guidance if absent; no new token format or connector. | Config validation without printing env values; task 6.4 |
| `docs/verification/standard-otel.md` (new at apply) | Durable evidence ledger for exact versions, fixture hashes, request statuses, persisted reads and all qualification cells. | Complete task 7.1; do not overwrite old reports |

## Shared practical setup path

Reuse the existing developer-guide package/source research. Revalidate against each application's lockfiles before editing snippets; do not blindly upgrade .NET/JS families. A .NET app reuses its hosted tracing/logging provider with ASP.NET Core and outbound HTTP instrumentation; a Next Node app reuses its current framework/Node provider initialized once from `instrumentation.ts`. Choose one initialization path; do not combine `@vercel/otel` and a second NodeSDK. Edge and browser need independent supported setups. Add no proprietary wrapper or generic telemetry framework; a small app logging adapter is optional and cannot hide the standard carrier.

Use ordinary OTLP/HTTP Protobuf. Common `OTEL_EXPORTER_OTLP_ENDPOINT` is the origin/base; explicit or signal-specific endpoints contain `/v1/traces` or `/v1/logs`. Verify the chosen SDK's header parser, using `Authorization: Bearer <producer-token>` and its documented env encoding. Metrics path is `/v1/metrics`, receipt-only. Distinguish same-Mac `127.0.0.1:8096`, another container's reachable service address, Home Server and HTTPS Tunnel routes. Do not default to gRPC or treat a proxy HTML login page as OTLP success. An optional Collector is a separate queue/acknowledgement and trust boundary.

Document two independently generated `openssl rand -hex 32` values for server/browser tokens, private secret environment storage and coordinated rotation. Never execute generation just to populate docs. Server tokens/viewer passphrases/invites never enter browser bundles, `NEXT_PUBLIC_*`, URLs, fixtures or Git. Browser SDKs use a reviewed bounded app-origin forwarder that forwards standard OTLP under browser authority; the server producer remains separate. Check consent/collection decisions, sender minimization, CSP/CORS and lifecycle flushing before real browser collection.

Cross-service tracing uses normal W3C propagation on approved destinations. Do not add the Sentrame ingestion token to business API requests. Existing HTTP instrumentation may already inject context; verify before adding manual injection. Standard logs should carry active trace/span context where the SDK supports it; plain console output is not automatically OTLP. Use actual exception/logging APIs, controlled messages such as `Vote persistence unavailable; retry permitted`, and error status where the operation actually failed. Do not change retry/error behavior or throw solely to manufacture evidence.

## One voting story in three increments

### Increment 1: ordinary OTel, configured in Sentrame

Story: navigate to Voting; press Vote; server processes/persists; first attempt fails; a retry later persists successfully; the configured goal is evidenced.

| Observation | Evidence already usable if emitted | Owner mapping / what remains unknown |
| --- | --- | --- |
| Navigate | Existing document/navigation or server `GET /vote` span | Label Voting request; backend span is not rendered-page/dwell evidence |
| Press | Existing browser interaction span or named `vote.pressed` point event | Intent/supporting activity; absent browser instrumentation means press unknown |
| Process | Existing server `POST /votes` span with child persistence operation | Label Vote processing; native elapsed duration, no stage pair |
| Failure | Failed operation status, standard exception event or correlated error log | Error evidence remains attached; generic request failure does not prove no commit |
| Retry | Distinct request/processing spans, optional existing retry event | Repetition remains visible; without supported shared scope, attempts are separate traces |
| Persist/recover | Trusted `vote.committed` named log/span event after actual commit; optional actual recovery event | Goal `Vote saved` selects the authoritative fact, not browser press or HTTP 2xx |

If existing telemetry contains a semantically precise successful commit span, its completion can serve that mapped operation outcome; if there is no such fact, request exactly one named post-commit observation through the existing SDK. That is ordinary application instrumentation, not a mandatory Sentrame schema. A publication error after commit does not undo the saved vote. Completion-screen navigation is a separate browser goal only if actually observed and configured. Do not fabricate it from persisted data.

Both recipes must show selecting service + carrier + operation/event name, optionally an approved retained HTTP field, entering Stage/Kind, previewing matches/missing context and activating a definition. Useful operation exploration requires no new attributes. No `stage.name`, occurrence ID, session ID or presence stream is demanded at this increment. Navigation, press and retry may be disconnected traces; this increment does not promise a full multi-request story.

### Increment 2: optional stage.name reduces labeling effort

Illustrative additions to existing emitters, **not runtime-qualified examples**:

```csharp
// Existing operation span; no new ActivitySource/provider required for this tag.
Activity.Current?.SetTag("stage.name", "voting.persist");
// For an existing post-commit ILogger event, use the configured scope exporter:
using (logger.BeginScope(new Dictionary<string, object?> {
    ["stage.name"] = "voting.persist"
})) {
    logger.LogInformation(new EventId(0, "vote.committed"), "Vote persisted");
}
```

```ts
// Existing active span under the app's existing provider:
trace.getActiveSpan()?.setAttribute('stage.name', 'voting.persist');
// Alternative: add to the existing named log's attributes, preserving context:
logger.emit({ eventName: 'vote.committed', context: context.active(),
  attributes: { 'stage.name': 'voting.persist' } });
```

The final guides must include appropriate imports and pinned logger/SDK APIs; these fragments show the intended bounded edit only. Tag only the real semantic boundary; never tag an entire request as persistence if that is not its meaning. Stable lowercase dot-separated labels (`voting`, `voting.submit`, `voting.persist`) simplify suggestions. No IDs, recipe names or free text. Changed meaning needs a different value; a label is not a start/end event.

Preview offers `voting.persist`, but an existing owner label **Vote saved** wins. A missing tag keeps ordinary mapping working. Several emitted labels prompt review; no active mapping is automatically selected. This increment reduces configuration entry, not the evidence needed for goal authority, correlation or timing.

### Increment 3: richer context only for the unanswered question

For one path spanning navigation, button, first request, retry and commitment, explicitly propagate a fresh bounded random `session.id` and browser `tab.id` on participating records using reviewed application boundaries. Verify that server events retain server authority and the browser cannot forge persistence. OTel context propagation relates technical calls; session copying is an explicit app convention, not automatic baggage behavior. Do not emit `session.previous_id` or join later visits. Use action IDs only if a specific intent/outcome relationship cannot be resolved by existing trace/span context.

For page elapsed time, use a real encompassing page interval or optional paired page-stage events with distinct occurrence identity if no native span represents it. For observed-active time, add actual visibility/idle observations under the documented 60-second heartbeat cap and stop positive intervals while hidden/idle. Native backend durations remain useful without this. Neither visibility nor a click proves attention; browser intent still does not prove persistence.

## Agent implementation checklist

1. Confirm authorized app/scope; read its instructions and lockfiles. Inventory existing providers, instrumentation, exporters, processors, logging bridge, propagation, sampling, shutdown and secret management. Preserve unrelated work.
2. Trace navigation→intent→server transaction→failure/retry→commit→recovery in actual source/tests. Identify which question ordinary telemetry can already answer and list missing facts separately from missing Sentrame retention/mapping support.
3. Prepare a concrete prioritized plan: question, existing carrier, minimum new emitter location if any, authority, retained selectors, volume/privacy, expected evidence and exact synthetic acceptance. Obtain plan approval when not already covered by authorization.
4. Reuse initialization. Make bounded changes only to necessary exporter setup or missing facts. No wrapper framework, duplicate SDK, raw request capture, altered transaction/retry/error flow or browser server-token exposure.
5. First verify receipt and durable protected discovery; then trace timing/relations; then owner mappings/hints; then only required journey/presence context. Synthetic unit tests first for actual changes, then real exporter seam qualification under authorization.
6. Exercise missing data, failure after commit, failed attempt/retry/success, replay and privacy/producer forgery. Verify actual SDK batching/flush and cross-service parentage, not simply shared trace IDs.
7. Record versions, process/runtime, destination/route, hashes/IDs without secrets, request statuses, persisted read evidence and passed/failed/blocked/not-run cells. Do not claim deployment, framework or browser support from source review.

## Troubleshooting order

| Symptom | Inspect first | Additional context only if justified |
| --- | --- | --- |
| No receipt | Endpoint/protocol/full path, auth/proxy status, SDK processor, sampling and flush | Do not add `stage.name` to fix transport |
| Receipt but no named event | Supported event-name API/log bridge, carrier, batching and sanitized discovery | Controlled named log/span event if no actual semantic event exists |
| Accepted-unassigned | Mapping preview, retained fields, selected mode and correlation prerequisites | Scope/tab only for the requested journey; trace exploration still works |
| Timing unavailable | Actual retained start/end, invalid clocks, selected basis and partial/exported span state | Pair only for a broader stage without a sufficient span |
| Broken correlation | Standard parent propagation, exact log context, scope/tab copying, sampling and clock gaps | Explicit approved journey context; no IP/time/fingerprint inference |
| Stage label differs | Explicit owner mapping and recorded definition, then emitted hint | Resolve suggestions prospectively; do not rename historical captures |
