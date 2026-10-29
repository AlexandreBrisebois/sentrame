# Qualification and acceptance matrix

Planning observation date: 2026-10-29. No telemetry, runtime test, package installation, deployment or external application change is authorized or performed in this planning pass. A listed test is a future acceptance obligation unless explicitly recorded as a planning check below.

## Evidence status rules

**Passed** means the exact named check ran and produced supporting output. **Failed** means an executed behavior/assertion failed; prerequisites/tool failures are **Blocked**. **Not run** means no execution, including checks intentionally excluded by this planning scope. Do not mark an unimplemented feature blocked merely because apply has not started. A previously recorded check is historical evidence, never a new pass for this model.

## Current planning evidence

| Check | Status | Evidence / practical limit |
| --- | --- | --- |
| Repository root and workflow | Passed | `openspec list --json`, `openspec context --json`: repo-local root and spec-driven schema |
| Canonical capability inventory | Passed | Initially empty; final `openspec list --specs` exposes eight canonical owners after concurrent archival; requirement blocks compared against inspected baseline |
| Implementation/source comparison | Passed | Intake, evidence types/serialization, journey derivation, store, query validation, UI and fixtures inspected; native ends/parents/links are missing today |
| Official OTel source checks | Passed | Linked sources in design checked on planning date; `session.id` is opt-in/Development, not stable mandatory session support |
| New change strict validation | Passed | `openspec validate standard-otel-default-integration --strict` after reconciliation; all 4 planning artifact types complete, 7 capability delta files and 23 unchecked implementation tasks |
| Baseline validation attempt | Failed | `openspec validate synthetic-journey-outcomes --strict` returned Unknown item after concurrent archive; this is a stale-target failure, not a product regression |
| Canonical strict validation | Passed | `openspec validate --specs --strict`: 8 passed, 0 failed; no canonical files edited here |
| Local links, fences, scenario shape and file preservation | Passed | Package check: 13 Markdown files, all local links resolve, fences balanced, scenarios contain WHEN/THEN, files resolve inside change root; 17 input files hash-identical including archived counterparts |
| Whitespace | Passed | `git diff --check`; new planning files separately checked for trailing whitespace |
| Runtime, SDK and route behavior for this change | Not run | No runtime code changed or telemetry sent |

Validation correction: the first diagnostic replacement renamed the legacy scenario header, which strict validation rejected as an omitted canonical scenario. The replacement now preserves the original `Projection upgrade and replay` header while explicitly changing its behavior to replay within the new model, with no pre-cutover compatibility obligation. The complete requirement and other three scenarios remain present. Strict validation passed after this correction. No blocked checks remain in planning; runtime cells below remain not-run.

## Prior evidence, not requalified here

`docs/integration/synthetic-demo.md` records direct Mac origin tests for logs/traces/metrics across HTTP Protobuf/JSON and none/gzip (12 wire cells) and Go SDK binary none/gzip (6 exporter cells). `internal/app/protocol_test.go` is the implementation reference. These establish the earlier projection/read path only. Existing container reports record separate AMD64/ARM64 smoke results, not the new trace model. .NET-like and Next-like diagnostic fixtures are standard wire shapes, not those frameworks' actual exporters. No previous result is copied into a new-model pass cell.

## Repeatable acceptance fixtures through real seams

Proposed version-controlled sources: `rest/fixtures/standard-otel.json`, `rest/standard-otel.http`, `cmd/fixtures` generation/verifier extension. Generate official pdata/proto requests with fresh synthetic timestamps and stable scoped identities; preserve exact bytes and hashes for replay. Ignore generated bodies/credential manifests and keep them private. Execute literal REST Client requests through the actual receiver→SQLite commit→derivation→protected read chain. No direct DB insertion or precomputed graph fixture qualifies the seam.

| ID | Synthetic case | Required protected read assertions | New-model status |
| --- | --- | --- | --- |
| Q01 | Standard spans/logs, no stage/session/occurrence/presence attributes | Accepted; service/name/carrier discovery; one trace structure; timing; no invented journey | Not run |
| Q02 | Nested, parallel and repeated spans; async link; correlated logs/events | Exact parent edges; no sibling adjacency; separately typed links; exported event order; trace-only logs not invented as children | Not run |
| Q03 | Span start/end with nanosecond precision; mapped native stage | Exact string timestamps, expected inclusive duration, no custom pair/ID; separate sample basis | Not run |
| Q04 | Missing parent/end; bad clock; sampling/drop counts; limit overflow | Partial/unresolved/unavailable/censored; no zero/negative sample; bounded disclosure; no live span invented if no export | Not run |
| Q05 | Standard-name/service/carrier/allowed-field mappings, no hint | Owner preview and activation; stage/goal matches and missing journey prerequisite; unsupported field unavailable | Not run |
| Q06 | Optional hint, owner override, missing/conflicting labels, overlapping rules | Literal accepted suggestion; explicit rule wins; two explicit matches block; no automatic label/Kind/goal/authority | Not run |
| Q07 | Voting failure→retry→trusted commit/recovery; publication error after commit | Preserve friction and success; browser intent and HTTP status alone not persistence; separate diagnostic records/occurrences | Not run |
| Q08 | Explicit combined scope/tab over several requests; server-only; missing scope | Correct lifecycle where sufficient; no fabricated tab server-side; missing evidence accepted-unassigned; no IP/person/fingerprint joins | Not run |
| Q09 | Exact native span/event replay; restart; conflicting span replay; anonymous repeats | One identified contribution, conflict atomicity/zero new rows, anonymous uncertainty; no overwrite/progressive span upsert | Not run |
| Q10 | Parent/goal late and out of order; new definition during active journey | Visible revision corrections under captured policy; no resegmentation from new label; saved reference mode/bounds preserved | Not run |
| Q11 | Optional explicit page pair and visible/hidden/idle intervals | No mandatory presence; capped observed-active separate from elapsed/native duration; no attention inference; no duplicate basis contribution | Not run |
| Q12 | Browser forges source/tag/trace; secrets in attrs/links/errors/routes | Namespace prevents authority/dedup collision; forbidden sentinels absent in stored projection, preview/read/SSE/copy; scrubbed diagnostic bounds retained | Not run |
| Q13 | Unauthorized/revoked viewer, owner mutation CSRF, linked target navigation | Protected trace/list/copy; viewer read-only; revoked stream/read denied; no leakage through errors/links | Not run |
| Q14 | Expiry/deletion/replay and summary boundaries; old schema startup | No expired detail resurrection or ID/text summaries; documented replay horizon; unsupported old schema error and no auto-reset | Not run |

Each verifier records declared response statuses plus read assertions; successful ingestion alone is insufficient. Demonstrate persistence by stopping/reopening the isolated service and rereading. A process restart proves only that tested durability boundary, not power-loss behavior. The literal REST verifier executes the `.http` file; report VS Code extension interaction separately if not manually exercised.

## Qualification dimensions — never collapse these rows

| Environment / producer | Required future check | Status / limit |
| --- | --- | --- |
| Mac direct origin, official wire fixtures | Q01–Q14 through HTTP Protobuf and JSON, none/gzip per claimed signal | Not run for new model; metrics receipt-only |
| Mac .NET actual exporter | Locked runtime/packages, `ILogger` events/error messages/stacks, Activity spans, shutdown, auth and exact parent propagation into/out of Next | Not run; source recipes alone insufficient |
| Mac Next.js Node actual exporter | Locked framework/Node/SDK, one provider, real Route Handler/server event/error path, log context, flush and request propagation | Not run; no Edge/browser inference |
| Actual browser exporter/forwarder | Controlled navigation/press/error, CSP/CORS/collection decision, browser namespace, absence of server secrets, lifecycle loss and optional presence | Not run |
| Cross-service Next→.NET and .NET→Next where used | Real standard propagation, parent IDs and log context; no duplicate initialization or ingestion token on business request | Not run; same trace string alone insufficient |
| Optional Collector path | Routing, namespace preservation, batching/retry and enqueue versus commit acknowledgement | Not run; direct receiver evidence does not qualify forwarding |
| Local container Linux ARM64 | New model, durable bind directory/restart, protected reads, Q01–Q14, readiness | Not run; previous image smoke is historical |
| Local container Linux AMD64 | Same checks; record native versus emulated execution | Not run; ARM64 pass cannot fill this cell |
| Native Home Server | Selected hardware/filesystem/ISA/workload, memory/CPU/disk growth/concurrent queries/cleanup | Not run; budget unavailable; no capacity claim |
| Existing shared Tunnel HTTPS route | Auth, paths/content types/gzip/bounds/errors/retries and real exporter receipt; SSE flush/heartbeat/reconnect/revocation | Not run; no connector/DNS changes in scope |
| Origin/public-route gRPC | Separate future scope; not implemented by this change | Not run; do not advertise from HTTP results |
| Production/WFS integration | Outside this change's implementation/qualification plan | Not run; no external app modified or contacted |

For every executed cell record commit/build identity, runtime and package lock, platform/ISA and emulation, route, transport/content encoding, fixture hashes, request/read outcomes, storage/restart evidence, privacy checks, failure reason and command. Keep secret values and personal data out of the report. After a failed behavior, retain the failure and link the corrected rerun; after a blocked prerequisite, record the exact blocker without calling it red coverage.
