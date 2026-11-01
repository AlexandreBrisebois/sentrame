# Local synthetic client demonstration

Use only an explicitly selected disposable loopback installation. These commands create/select the demonstration app and submit synthetic observations through actual APIs. They must not target real retained data without separate authorization. They never preload a separate graph database.

## Generate and run

Start the built service with a private synthetic owner passphrase, matching loopback origin/address and a fresh data folder. From the same configured shell, use the Task fixture commands:

```sh
task demo
task demo:workspace
task demo:diagnostics
task demo:standard
```

The tool registers/selects **Sentrame demonstration**, reads its one `telemetryKey` through owner connection settings and configures its own sender in memory. The generated `rest/.env`, login bodies, protobuf files and manifests are private ignored local output. Never paste their values into reports or commit them. No per-role environment credentials are needed. Existing SDK/exporter/route qualification is separate from official wire fixture success.

## Product walkthrough

The current client scenario uses explicit random visit/tab context, navigation, intents, `vote.result.observed`, optional goal `completion.entered`, healthy repeated results, client errors/recovery and distinct timing bases. Missing-context observations remain received/unassigned; a claimed runtime origin grants no stronger assurance. Stable event IDs are unique within the app/signal identity domain; exact replay reuses identical bytes and does not invent additional actions.

Use **Tools**, the seven direct journey views, retained evidence and client diagnostic Traces to inspect the question and return to its context. Goals are optional. Native spans do not create visits automatically; HTTP success/result labels do not prove durable persistence. Live current/baseline populations and timing units remain separate.

For the live demo, `task demo:live` starts a new synthetic capture and `task demo:followup` sends its later screen observation. Replay preserves prepared identity/bytes. Do not reuse a follow-up manifest whose acknowledgement state was changed outside the CLI; start a new disposable case instead.

## VS Code REST Client

The `.http` requests use one `{{telemetryKey}}` plus owner-session CSRF for protected mutations. Use the generated private environment and exact files. The CLI executes the same literal requests for reproducibility; this does not claim that the VS Code extension itself, real browsers or deployed routes were tested.

## Standard OTel voting operations

`task demo:standard` exercises ordinary spans/logs, native structure/timing and standard one-key receipt, then service/carrier/name mappings, optional label hints and an explicitly scoped client retry/result/screen story. Navigation/result stream filenames identify fixture contents, not producer roles. The qualified client journey computation is `journey-client-v5`.

The retained wire matrix and actual exporter/browser/container/Home Server/Tunnel cells must be reported separately after running this source version. Earlier verification below is preserved historical evidence of the previous browser/server model; it is not active server integration support or qualification of the single-key release.

## Recorded local verification

On 2026-10-28 at 20:00 UTC, `go run ./cmd/fixtures demo` executed all **20 requests in the supplied `rest/demo.http`** against `http://127.0.0.1:8096`. All declared statuses passed: unauthenticated read 401; owner login/read/SSE and valid binary-log ingestion 200; mixed valid/future-invalid direct request 400. Per-scope protected reads confirmed healthy=16 records including ten trusted commitments, recovery=6, initial convergence=1, server-only=3, forged authority=1 unassigned record, inactivity=1, historical baseline=2, live with follow-up=2 and atomic-invalid=0. Healthy and recovery scopes each derived one completed journey; silence derived inactive; browser-forged commitment derived no journey. The repeated identified server request did not multiply retained votes.

An earlier run returned every expected request status but failed the additional live count assertion because its follow-up timestamp was two seconds in the future. The fixture now uses an already-occurred start and current follow-up, with live T0 two seconds before generation; the full request sequence and assertions then passed. This was a fixture timing correction, not a weakened read assertion.

Focused generator/parser tests passed. The shared-record identity regression was observed red (browser/server source IDs differed), then green after preserving a readable source record label across producer namespaces. The first generator test initially hit a missing dependency checksum, which was a tooling prerequisite failure and is not counted as behavior-red evidence.

These are direct local service observations. Tunnel, Home Server, VS Code extension interaction and actual browser exporter operation were not exercised by this verifier. The CLI executed the same version-controlled HTTP requests; it does not claim the VS Code extension itself was tested. Earlier failed/successful demonstration generations may remain as separate synthetic scopes unless the owner invokes the explicit local reset command.

## Direct-origin protocol qualification

`internal/app/protocol_test.go` qualifies the actual maintained receiver, sanitizer, SQLite writer and protected snapshot API using isolated local services:

| Signal | HTTP binary, none/gzip | HTTP JSON, none/gzip |
| --- | --- | --- |
| Logs | Passed with official pdata wire fixtures and actual Go `otlploghttp v0.23.0` / SDK log `v1.47.0` exporters | Passed with official pdata JSON fixtures; no JSON SDK/exporter claim |
| Traces, including span-event receipt | Passed with official pdata wire fixtures and actual Go `otlptracehttp v1.47.0` / SDK `v1.47.0` exporters | Passed with official pdata JSON fixtures; no JSON SDK/exporter claim |
| Metrics, receipt inspection only | Passed with official pdata wire fixtures and actual Go `otlpmetrichttp v1.47.0` / SDK metric `v1.47.0` exporters | Passed with official pdata JSON fixtures; no metric arithmetic or JSON SDK/exporter claim |

The wire test covers all twelve signal/encoding/compression cells, unknown protocol fields, missing optional business context, forbidden raw body/resource/point/span-event attributes and protected sanitized reads. Six binary SDK cells use ordinary endpoint, authorization header and compression options on unmodified maintained exporters. No proprietary producer library is involved. SDK tests assert retained public-read evidence rather than accepting exporter construction as success.

Negative tests cover absent/invalid producer credentials, malformed protobuf/JSON, unsupported content type/encoding, corrupt gzip, the 1 MiB request and 4 MiB decompression bounds, zero records from those rejected requests and bounded overload with retryable HTTP 503. Rejections carry safe standard `google.rpc.Status` encoded as protobuf or JSON to match the request. A regression test first observed HTTP 400 with missing RPC code (`0`); after the lead fixed sanitized response encoding, InvalidArgument code `3` and the other standard error-body checks passed.

These checks establish the listed direct local cases only. gRPC is not exposed/qualified. Public Tunnel, browser JavaScript/.NET exporters, forwarding Collector batching/acknowledgement, browser lifecycle/CSP/CORS and power-loss durability remain not run. The direct-origin HTTP support above does not establish those paths. The documented attribute-bound regression checks are now green: >64 resource attributes across logs/traces/metrics; >64 instrumentation-scope attributes across all three signals; >64 span or span-event attributes; and >64 metric-point attributes for gauge, sum, histogram, exponential histogram and summary. Each invalid direct request returns 400 and commits zero records, including when a valid metric point precedes the invalid point. The initial resource-bound test exposed traces/metrics incorrectly returning 200 with new records; after the lead added prevalidation, that regression and the expanded tests passed. This is evidence for these concrete limits, not exhaustive protocol conformance.

Latest focused verification after the integrated fixes: `go test ./internal/app -run TestOTLP -count=1` passed (1.773 seconds package time). All direct-origin protocol, SDK, standard error-body and resource/scope/span-event/metric-point bound tests in this file are green.

Manual live timing regression verification: the new behavior test first failed because the follow-up retained its generation-time timestamp, then passed after timestamp-only preparation and stable retry support. A fresh `live` case was sent to the running service, followed later by `followup`; protected reads confirmed two actual records and one completed journey, with completion time after generation. Repeating `followup` preserved exact protobuf bytes and manifest hash. The final full `demo` at 2026-10-28T20:23:26Z passed all 20 literal requests and all scoped assertions. Its verifier set the sent marker; subsequent `replay` and `followup` returned success with unchanged bytes, and protected reads still showed healthy=16 and live=2 records, each one completed journey. The latest generated manifest is this complete verified demonstration. No reset was performed.

## Twenty-sample timing demonstration

The current readable fixture set additionally includes `timing-samples`: discovery, twenty uniquely identified stage entry/exit pairs lasting 1–20 seconds, and completion, totaling 42 actual log records. Every timestamp is already in the past at generation, in a single interval of less than six minutes; no inactivity boundary is crossed. The source file lists every pair explicitly. The generator test observes that the scenario exists, all 20 pairing identities are distinct, each exit preserves its entry identity, all records are past and all durations match 1–20 seconds.

The added literal ingestion and protected scoped snapshot bring the **current `rest/demo.http` to 22 requests**. After generation, REST Client's `timingScope` variable selects this scope. The verifier asserts 42 retained records, one completed captured journey and the Stage node's `measured-elapsed` summary: sample count 20, median 10,500 ms, p90 18,000 ms and completeness 1. This demonstrates supported slower cases through the same ingestion/storage/derivation/read path; it never inserts precomputed timing summaries. The historical 20-request runs recorded above describe their earlier file revision and are preserved as historical evidence.

The new fixture test first failed because the readable scenario was absent, then passed after the explicit source scenario was added. `go test ./cmd/fixtures -count=1` passed (0.234 seconds package time) with the new sample checks. The expanded 22-request sequence awaits the lead's final demo run against the restarted service; no new generated data was written during this addition.


## Final expanded request verification

On 2026-10-28 at20:36UTC, `task demo` passed all22 current literal requests, including the42-record/20-occurrence timing case. Protected API reads verified median10.5s,p90=18s and100% completed timing coverage. After the final source-aligned restart at20:51UTC, `task demo:verify` passed all22 again with the same identified bodies; healthy remained16 and live2 records, timing42. No DB preload or reset. Final backend checks also preserve canonical trace/span hex IDs, reject invalid standard span IDs atomically, and inherit only verified parent logical context for span events.


## Standard OTel voting operations

The [readable source](../../rest/fixtures/standard-otel.json) is ordinary OTLP with relative timestamp offsets: `POST /votes`, parallel repeated `vote.persist` operations, a trusted `vote.committed` span event and correlated named log, trace-only log, missing parent/end, invalid clock, dropped counts and34 links. It contains no stage/session/occurrence fields. These are wire observations, not actual .NET/Next.js exports or proof of a real saved vote. The first increment explores operations; it does not join attempts into a multi-request journey.

With a separately configured disposable loopback service and its environment exported, run `task demo:standard`. This generates fresh native trace IDs and timestamps, preserves exact bytes/hashes under ignored `rest/generated`, and executes [rest/standard-otel.http](../../rest/standard-otel.http) through actual ingestion/storage/protected reads. `task fixtures:standard` prepares the files without sending; `task demo:standard:verify` reuses them. Private `rest/.env` contains the REST Client variables. Never paste producer tokens into tracked requests. Do not rerun the anonymous log fixture against an existing dataset when checking log counts: logs without stable IDs intentionally remain distinct. Native span replay is idempotent.

The verifier checks native timing, link bounds, contextual evidence and replay through HTTP binary, JSON span replay and gzip binary span replay. Full framework/exporter/protocol matrices remain separate in [the ledger](../verification/standard-otel.md). Owner mappings and optional stage enrichment are later increments; browser intent and generic HTTP2xx do not prove persistence.


## Scoped voting story

The third increment uses [standard-journey.json](../../rest/fixtures/standard-journey.json) and explicit owner [mappings](../../rest/fixtures/standard-journey-mapping.json). It adds random scope/tab only for the multi-request journey. It sends browser navigation/intent, two separate native server request traces (failure then successful retry), trusted `vote.committed`, publication failure/recovery, explicit page pair and optional synthetic presence observations. The separately configured browser completion-screen goal arrives after another definition activates. Protected reads prove the existing capture keeps its definition, earlier friction and one trusted commitment; HTTP2xx does not invent the screen goal. Exact goal replay adds no journey.

The local verifier now executes37 literal steps, including trace-mode references and separate native/pair/observed-active timing. Native durations are200ms and300ms, the page pair is1s and the **synthetic** visibility/idle intervals total800ms. These are distinct bases, never attention or a unique-time sum. Actual browser visibility/idle collection and lifecycle loss are not qualified by these wire observations. Use a new disposable database/generation for a full repeat, preserving prior local data.
