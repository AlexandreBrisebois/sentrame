# Local synthetic demonstration

These demonstrations use the same authenticated standard OTLP ingestion, SQLite commit, derivation and protected read paths as other local producers. They do not instrument or contact WFS. Readable source cases are in `rest/fixtures/scenarios.json`; they are semantic fixture descriptions, not a proprietary telemetry protocol. The generator translates them to official OTLP Protobuf messages using the pinned `go.opentelemetry.io/proto/otlp` schema and `google.golang.org/protobuf` runtime.

Run from the repository root after starting the local service with the same environment:

```sh
go run ./cmd/fixtures generate
go run ./cmd/fixtures verify
```

`SENTRAME_PASSPHRASE`, `SENTRAME_BROWSER_TOKEN` and `SENTRAME_SERVER_TOKEN` must be available to the generator. `SENTRAME_URL` defaults to `http://127.0.0.1:8096`. No default credentials ship with the fixtures. Keep the generated `rest/.env`, `rest/generated/login.json`, binary files and manifest private; their directory is ignored by Git and generated files use mode 0600. The fixture tool prints statuses and bounded scope counts, never credentials or telemetry bodies.

`generate` creates fresh occurrence timestamps and random logical, tab, record, action and stage identities. Cases with the same source scope label share a generated logical/tab scope; action and stage identities remain distinct. `rest/generated/manifest.json` preserves the generated identities and SHA-256 hashes. Re-running `generate` creates a fresh demonstration; it does not reset existing data. Re-running an identified POST or `replay` uses the exact existing protobuf bytes and tests idempotent delivery. Records deliberately lacking identity cannot promise deduplication.

The paired generated `.json` files are readable official OTLP JSON encodings of the binary requests. The HTTP requests send `.bin` bodies with `application/x-protobuf`; they do not send JSON renamed as Protobuf. Generation never writes to the database.

## VS Code REST Client

1. Start the local service and generate the fixtures with its local synthetic environment.
2. Open `rest/demo.http` with the REST Client extension. The adjacent ignored `.env` supplies the endpoint, producer tokens and generated view variables; `generated/login.json` supplies the owner login body.
3. Send the first unauthenticated snapshot request and confirm 401. Send owner login, allowing REST Client to retain its session cookie. The later protected requests use that independent owner session.
4. Send the ingestion requests in order. The identical server replay should return success without adding commitments. The healthy browser entry and first server commitment deliberately share one local record ID; trusted producer namespaces keep them distinct. The mixed valid/future-invalid request should return 400 and commit zero records.
5. Open the healthy scoped snapshot, then the fixed seven-day app baseline behind the live scope. The final SSE request streams authenticated invalidations; stop the request to close it.

The verifier parses this exact `.http` file, substitutes the same generated variables, reads the same binary/login files and executes its literal methods, paths, headers and bodies. It checks the declared status of every request, privacy sentinel absence in public responses, initial SSE delivery, then checks scoped evidence/model results through protected snapshots. Those additional assertions are API reads, not alternate fixture insertion.

## Product walkthrough

- **Healthy loop:** browser entry and recipe-stage timing surround ten trusted server commitments. A completion-screen observation ends one journey. Inspect ten Vote occurrences and nine Vote→Vote transitions; replay must leave them unchanged.
- **Supported slower cases:** the `timing-samples` scope contains one completed journey with twenty explicit Stage entry/exit pairs. Its measured durations are 1 through 20 seconds; the scoped map/table should show 20 samples, median 10.5 seconds, p90 18 seconds and 100% completed timing coverage. These are stage occurrences within one journey, not twenty independent visits.
- **Post-commit recovery:** distinguish browser intent, persisted server vote, publication exception, recovery and later completion. The error remains friction and never undoes the commitment. Unsafe message/stack/household/URL sentinel attributes are omitted.
- **Initial convergence:** a completion screen without votes still evidences that visit's screen goal. It does not prove who caused agreement.
- **Server-only:** an operation has no invented browser tab and ends with an explicit exhausted reason.
- **Missing context and authority:** ordinary telemetry lacking session/tab/event IDs remains accepted-unassigned. A browser `vote.committed` carrying self-asserted server attributes cannot establish trusted persistence.
- **Silence:** a 31-minute-old entry closes inactive, not failed or tired. No fabricated exit reason is emitted.
- **Historical versus live:** a prior-day case belongs to the app baseline. A new live scope is separate and receives a completion follow-up with the same logical identity. The fixed baseline remains app-wide while live is scoped.

For a fresh manual live observation, run:

```sh
go run ./cmd/fixtures live
```

This regenerates fresh synthetic identities and sends only `live-start.bin` through actual ingestion. Select the new active live scope from `rest/generated/manifest.json` in the UI and start live observation. Then run:

```sh
go run ./cmd/fixtures followup
```

On its first manual send, `followup` refreshes only occurrence/observed timestamps to the current time, so the completion occurs after the UI's live T0. It preserves the scope, tab, action and event ID, rewrites the paired protobuf/JSON body and updates its manifest hash. The manifest is marked **prepared before the HTTP attempt**; retries then reuse exact bytes, including after an ambiguous network failure. A successful response marks it sent. The literal REST verifier also marks its successful follow-up sent; running the CLI after verification therefore replays unchanged bytes instead of conflicting with accepted evidence.

Start a fresh `live` case for the manual walkthrough rather than reusing a completion already sent independently through VS Code: externally clicking a REST request cannot update the CLI's local acknowledgement marker. The follow-up uses the actual receiver/storage/derivation path and does not rebuild a separate in-memory journey. To replay the ten commitments:

```sh
go run ./cmd/fixtures replay
```

## Qualification boundaries

The fixture requests can verify direct-origin OTLP/HTTP binary logs and the actual owner/SSE/read seams. They do not by themselves qualify a standard SDK/exporter, HTTP JSON, gzip, traces, metrics, gRPC, browser SDK lifecycle/CORS/CSP, a Collector forwarding acknowledgement, power-loss durability, Home Server capacity or Tunnel routing. Record those separately in release verification. Repeated use retains multiple demonstration scopes until explicit synthetic reset or retention; the verifier checks only its newly generated scopes.

The `.http` file is version-controlled; generated credentials and protobuf bodies are intentionally reproducible local outputs. Never substitute real app credentials or personal content to run this demonstration.

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
