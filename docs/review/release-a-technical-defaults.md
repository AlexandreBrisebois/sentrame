# Release A technical defaults

Reviewed 2026-10-28 for the explicitly authorized local synthetic implementation.
These choices resolve implementation defaults; they are not runtime qualification. The initial prerequisite blocker recorded below was resolved later in this session; see the final verification and current implementation notes at the end.
The implementation request supersedes historical planning-only authorization text.
Release A remains the complete workflow in the release scope; roadmap B is deferred.

## Base and preservation

- Branch: `codex/sentrame-mvp`; HEAD and `origin/docs/journey-capability-plan-2026-10-27` both resolve to `4a20c953381f03b4400712891f4a7e890c340835`. Ancestry passed against the remote-tracking ref; there is no local planning branch. No fetch or remote freshness claim.
- One worktree: `/Users/alex/Code/sentrame`. No applicable ancestor or repository AGENTS/AGENT file was found. Use the repository's OpenSpec apply skill and engineering-harness discipline.
- Existing modified README, decision register and design, plus untracked visual-theme assets and previous-MVP lessons, are user work. Preserve them. Do not restore the old stash.
- Read the README, decision register, release scope, proposal/design/tasks, all eight capability specs, engineering-harness proposal, previous-MVP lessons, visual-theme reference, deployment contract and external integration guide.
- Preserve warm paper `#f5f4ee`, forest green `#153f39`, serif headings, system-sans body, spacing, outlined cards, visible focus and the linked-thread SVG. Maps require equivalent tables and ordered sequences.

## Synthetic semantic contract

Use a synthetic app named `Sentrame demonstration`, with a random immutable registration ID and intake generation. Illustrative event names are local fixture choices, not claims about WFS emissions:

| Event | Meaning / owning producer |
| --- | --- |
| `discovery.entered` | Actual screen entry, browser |
| `vote.requested` | Intent only, browser |
| `vote.committed` | Persistence milestone, trusted synthetic server |
| `publication.failed` | Post-commit friction, trusted synthetic server |
| `publication.recovered` | Explicit recovery, trusted synthetic server |
| `completion.entered` | Actual completion-screen entry, browser; initial convergence permitted |
| `journey.exited` | Explicit `exhausted`, `agreed` or `departed` reason from its configured owner |
| `stage.entered` / `stage.exited` | Semantic stage occurrence boundaries from the same producer/clock |
| `activity.changed` | Explicit visible/hidden/idle/active observation, browser |

Initial active goal: browser `completion.entered`. Commitment, positive shared agreement, explicit exit and screen observation stay distinct. Ten distinct commitments are healthy; nine successive commitment transitions are possible. Severity never closes a journey. Shared convergence never identifies another participant or proves that the observed visit caused it.

Mapping configuration selects existing safe attribute keys for random logical scope, per-tab scope, record identity, action identity, stage occurrence and optional source sequence. Those fields remain optional for standard reception. Browser evidence lacking safe tab scope is accepted-unassigned; server-only logical operations require no tab. Combined evidence uses an explicitly propagated app-scoped correlation plus tab scope where applicable, never proximity, household or person joins.

Producer authority comes from configured credential-to-app/generation/namespace bindings. Separate browser and server credentials and namespaces; browser credentials cannot support server persistence even when telemetry says `producer=server`. Browser keys are untrusted public claims. The local synthetic server credential is never included in browser code. No forwarding or real-app instrumentation is authorized.

## Receiver, limits and durability

Choose the embedded maintained Collector OTLP receiver with a synchronous custom storage consumer, without an asynchronous batch processor ahead of storage. This gives one lifecycle and one commit boundary. Candidate pinned receiver release: `v0.164.0` (stable Collector modules `v1.70.0` where required by that release), checked against the [official release](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.164.0). Candidate Go toolchain: `1.27.3`, listed on the [official downloads page](https://go.dev/dl/). Source inspection, dependency resolution, exact transitive checksums and standard-exporter pins remain blocked until Go/network prerequisites are available; no module compatibility claim yet. Lock the resolved graph in `go.mod`/`go.sum` and frontend dependencies in `package-lock.json` before qualification.

Use Svelte/TypeScript + Vite static assets served by Go and SQLite on a local filesystem. Prefer a maintained SQLite Go driver; resolve/pin it with the first tracer rather than inventing an unverified module version. A sidecar is a measured fallback only if embedding is materially costly; an ad hoc protocol handler is not the default.

| Boundary | Local default |
| --- | --- |
| Binding | Loopback only; intended viewer `127.0.0.1:8096`, OTLP HTTP `127.0.0.1:4318`; check ownership before binding |
| Direct request | 1 MiB compressed, 4 MiB decompressed, 2,000 total projected records/points; reject entire request on overflow |
| Structured cost | 64 attributes per record/resource/scope, 64 KiB projected record, 2,048 histogram buckets per point; bound aggregate bytes before allocation growth |
| Writer | One synchronous writer, maximum 8 concurrent intake requests, no unbounded waiting queue; overload returns retryable status |
| Deadlines | Headers 5 seconds, intake 10 seconds, write transaction 5 seconds; no stream-long DB transaction |
| Disk | Local synthetic DB ceiling 1 GiB including accounted WAL; stop intake before exceeding the bound; expose health, never silently discard acknowledged data |
| Query/render | 30-day detail window, 100 graph nodes, 300 edges, 200 sequence rows/page, explicit limits/continuation; no silent complete-count claim |
| Live | 8 connections, 16 invalidations/connection, 15-second heartbeat, 5-second lifecycle freshness bound; overflow/restart emits reset |

Prevalidate the whole direct request, separating protocol validation from privacy projection. Unknown attributes are omitted, not rejected merely for being unknown. Missing optional metadata stays unknown. Use explicit InvalidArgument/HTTP 400 for permanent validation and Unavailable/HTTP 503 for commit/overload failures; qualify receiver error adaptation before advertising it. No payloads in errors or access logs.

SQLite uses WAL, foreign keys, `synchronous=FULL` and bounded busy timeout. A direct success follows an atomic evidence/revision commit. Initial qualification target is process restart/crash recovery; OS crash, power loss and storage-controller flush behavior remain unqualified even with FULL. Upstream exporter/Collector enqueue or rebatching acknowledgement is a separate boundary and cannot guarantee final storage or original SDK-envelope atomicity.

Target logs/traces/metrics, HTTP binary/JSON and none/gzip, then origin gRPC only with separate tests. Advertise zero compatibility cells until a pinned standard producer passes the real receiver/store/read seam. Metric receipt is not arithmetic; profiles are excluded. Tunnel and browser exporter/CSP/CORS/lifecycle support require their own qualification and are unclaimed.

## Privacy and identity

Default deny raw bodies, request/response content, URLs, headers, arbitrary messages/stacks, personal/household/recipe/device identifiers and unknown attributes. Permit bounded safe event/type/reason codes, random ephemeral context and allowlisted technical service/version/runtime metadata. For synthetic diagnostic examples use controlled exception type and frame symbols; arbitrary diagnostic free text remains omitted until a separately reviewed scrub policy. UI/copy/SSE use the exact same sanitized projection and inert text rendering. Report omissions and unavailable detail.

Deduplicate only supported identity under `(app, generation, trusted producer namespace, signal, source record identity)`. Preserve span identity under qualified trace/span rules. Conflict rejects the direct request atomically and never overwrites evidence. Compare sanitized semantic content, not unsafe raw payload hashes. Without a stable identity, retain legitimate repeats and disclose retry uncertainty. A shared action is neither a duplicate nor authority. Cross-carrier exception occurrence counting requires explicit supported occurrence identity; a stack signature is grouping only.

## Time, boundaries and statistics

- UTC nanosecond occurrence/observed/receipt timestamps; display timezone explicitly. Accept valid occurrence times up to 30 days late, reject more than 5 minutes in the future. Missing/zero time uses first receipt for a known identity with a fallback flag; unidentified retries cannot promise original-age preservation.
- Definition revisions activate prospectively at receipt. Capture the active definition on first accepted evidence and keep each captured journey pinned; late events reconcile retained evidence within that definition. No silent retrospective reinterpretation. Saved views never activate definitions or resegment.
- First mapped occurrence starts; goal or explicit configured exit closes; next distinct mapped occurrence starts again. Thirty-minute gaps close inactive, not failed. Late data can revise provisional inactivity within retained detail. Replay deduplication precedes segmentation.
- Prefer reliable sequence within a single source/clock; otherwise event-time order plus stable receipt-row tie-break, visibly uncertain. Cross-source timestamps cannot prove causality. Unavailable or negative clock measurements never become zero durations.
- Pair explicit stage entry/exit by stage occurrence within app/logical scope/tab/producer. Elapsed and observed-active are separate. Active intervals require explicit visible/active evidence, close on hidden/idle and cap missing-heartbeat activity at 60 seconds. No activity evidence means active duration unavailable. Overlaps are separate occurrences, never summed into a unique attention total.
- Next-event fallback is opt-in, estimated, only for genuine adjacent evidence within the same captured boundary/source clock and at most 30 minutes. Filtering cannot manufacture adjacency. Ongoing/missing/ambiguous pairs remain censored/unavailable and stay out of completed percentiles.
- Exact median for retained eligible samples; p90 nearest rank only at >=20 completed samples of the same timing kind/policy. Show numerator, exclusions and completeness. Historical summaries use mergeable distributions, never averaged medians.
- Graph node **area** encodes median duration, capped at 30 minutes with explicit clipping and numeric values; missing is unavailable. Edge width encodes transition occurrences with a bounded legend. Counts remain distinct from affected journeys/visits. No color-only status or cognitive/causal inference.

Historical cohorts use recorded journey starts in `[from,to)` with follow-up through snapshot time, and separate partial starts/right censoring. Version is initial mapped app version, with mixed/unknown flagged and per-event versions preserved. Default low support is fewer than 20 eligible journeys; show rare path counts/denominators rather than bug labels. Hidden intermediate steps preserve gaps. Repeated endpoints in saved views use explicit occurrence selection within recorded segments, with boundary markers and unavailable-start explanations.

Live baseline freezes `[T0-7 days,T0)` at observation start; live uses `[T0,now]`. Late pre-T0 data corrects baseline, not live counts. Baseline is app-wide under compatible version/definition filters, while live can isolate one visit. Snapshot/cursor share a committed revision. SSE sends bounded invalidations; clients refetch full authorized snapshots and reject superseded subscription generations, preventing additive double counts. Missed notifications/restart require explicit reset.

## Access, retention and deletion

Owner passphrase via environment/config secret input with no default, a maintained bcrypt verifier at cost 12, rate limit 5 failed attempts/minute per local address plus 30/minute globally. Independent random 256-bit cookie tokens are hashed server-side; HttpOnly, host-only, SameSite=Strict, Secure outside explicit loopback HTTP mode. Idle timeout 30 minutes, absolute lifetime 8 hours. Logout/passphrase rotation revokes sessions and active streams. Cookie mutations require exact allowed Origin and CSRF token; reject missing proof. No credentials in URLs, telemetry or logs. Invitations remain B.

Detail eligibility ends at valid occurrence +30 days, or supported first-receipt fallback. Reads exclude expired rows immediately; lifecycle invalidates within 5 seconds and bounded cleanup removes physical rows within 60 seconds under healthy operation. Declare health failures instead of promising cleanup deadlines under unavailable storage. SQLite pages/WAL/filesystem snapshots are not forensic erasure.

Aggregate at UTC day × app generation × definition × initial app version × supported stage/timing kind or outcome. Store counts, conversion numerators/denominators, censoring/coverage and mergeable logarithmic duration histograms (ratio 1.1, 1 ms–24 h plus under/overflow; approximate bucket quantiles explicitly labeled). No visit IDs, ordered paths, exception payload or individual membership after detail expiry. Rebuild retained contributions transactionally for late/replay/deletion corrections; finalize expired-day summaries before raw purge. Expire summaries at UTC calendar date minus 12 months, not an implied fixed 365-day interval. Unsupported filters/statistics are unavailable. No path reconstruction from edge counts.

Support explicit visit or app deletion. Journey-only deletion is unavailable without stable replay membership. Visit tombstones retain only generation-scoped keyed identity and expire after 31 days; timestamped evidence outside the 30-day horizon is rejected. Beyond that horizon, unidentified/missing-time replay has no original-age guarantee. App deletion disables credentials permanently for that generation; explicit re-registration uses a fresh generation/credentials and preserves earlier tombstones. Within retained detail, subtract/rebuild affected contributions; after individual membership expires, selective summary subtraction is unavailable and requires explicit app/time-bucket deletion. No automatic backups in this local synthetic release. A future restore must apply current deletion ledger/expiry before enabling queries.

## Local evidence and gates

Provisional synthetic qualification envelope: 10 records/second sustained, 100/second for 10 seconds, 100,000 retained records, 2 active viewers, <512 MiB process RSS, <2 seconds p95 snapshot and <1 second p95 direct commit, <120 seconds warm build/test and <10 minutes cold dependency/build. These are local test targets, not measured capacities or owner-approved Home Server budgets. Record actual measurements and adapt only with evidence.

Observed on this Mac: `uname -m` returned `arm64`; `df -h .` reported ~514 GiB available. Node `v26.12.0`, npm, Task and OpenSpec `1.14.3` are available. `sysctl` resource inspection was denied; RAM/CPU/model remain unmeasured. No Go executable in PATH or checked conventional install paths. Direct `curl` to `https://go.dev/dl/?mode=json` failed with exit 6 (DNS resolution). This is a prerequisite blocker, not a red behavior assertion. No runtime, build, exporter, database or browser verification has run.

Home Server ISA/filesystem/shared workload/numerical budget remain unknown and must not block authorized local synthetic semantics. Task 1.3 and Home Server G3 qualification stay unchecked. Local G3b/G3c are blocked before the tracer, not passed. Tunnel, Compose, real WFS integration, real telemetry, deployment/publication and power-loss qualification are unclaimed and outside this request.

## Implementation and review sequence after prerequisites

1. Shared contract and one observed red/green vertical tracer: pinned standard log exporter → maintained receiver → sanitized SQLite commit → protected read after restart → accessible Svelte sequence row.
2. Missing optional context/privacy, atomic failures/replay and local measurement before expansion.
3. Lifecycle/healthy loops, diagnostics/recovery, timing/maps, saved views, SSE/baseline and aggregate/deletion increments, each with focused red/green evidence through the same seams.
4. Delegate only after shared API/storage/fixture contracts are written, with exclusive source ownership and one integration lead. No all-backend/all-frontend split before the tracer exists.
5. Generate official protobuf bodies from readable fixtures with fresh IDs/timestamps; retain a scenario manifest for replay/follow-ups. Version-controlled REST Client requests and the verification runner must execute the same literal requests against the running service. No direct DB demonstration preload.
6. Taskfile: setup, start, test, fixtures, demo, scoped synthetic reset and strict spec validation. Record actual resource/build/restart/concurrency/browser/keyboard/language evidence. No complete-release claim until the full A workflow passes.

Highest-value product review choices once running: whether completion-screen observation is the preferred initial goal, whether conservative omitted free-text diagnostics are sufficient, whether estimated timing is useful beside explicit stage timing, and whether cohort/low-support explanations make potential dead ends understandable. None blocks local synthetic implementation.

## Current implementation reconciliation

The Mac Go/network prerequisite is resolved: Go1.27.3 and approved dependency downloads were used. The reviewed choices above preceded implementation; candidate pins now resolve in `go.mod`/`go.sum` and `web/package-lock.json`. Receiver0.164.0 / stable Collector1.70.0 and SQLite driver modernc1.61.0 are in use. Actual direct-origin HTTP signal/encoding/compression qualification is recorded in the integration guide; the initial zero-cell statement is historical.

Runtime details refine the review defaults: public loopback viewer+OTLP port8096 proxies to an ephemeral private loopback maintained receiver; gRPC disabled. Eight intake and eight stream slots are bounded; SSE has no accumulating per-viewer queue and emits a full-snapshot reset every5seconds, with a5-second write deadline and revocation check. SQLite writes have5-second transaction contexts. One-second lifecycle work advances revision on status/expiry changes. Reads release store locks before network serialization. Detail response limit2,000; indexed recent derivation/discovery limit10,000; graph100nodes/300edges; UI sequences expand in50-row increments. Oversized/truncated views remain labeled.

Byte accounting refuses new intake above896MiB across DB/WAL/SHM, reserving space below the1GiB local envelope; SQLite max_page_count200,000 and journal_size_limit16MiB complement the100,000-record cap. Full-capacity disk-pressure/retention soak is not qualified. Local storage test evidence does not prove OS/power-loss durability or forensic erasure.

Correction/finalization ends at captured-start+30days. Expiry preserves only bounded partial-boundary markers while related detail remains eligible; late retained evidence inside an already-finalized boundary cannot create eligible completions or a second summary. Summaries store no individual membership after detail expiry and support app-wide complete UTC day/version/definition counts and duration histograms, including exclusions. Visit-scoped/partial-day/new-dimension summary queries are unavailable. Saved references persist first/last endpoint choice and project graphs inside captured boundaries. Evidence occurrence/observed/receipt nanoseconds use exact decimal JSON strings; graph timestamp displays are at millisecond precision.

Authentication is owner-only, with a bcrypt12 verifier, independent hashed cookie sessions, exact Origin+CSRF mutation proof,5failed-logins/minute and two concurrent password verifications. Restarts/rotation revoke all in-memory sessions. No global30-attempt policy or proxy/client-IP trust qualification is claimed for this loopback build; public auth deployment remains separate.

The measured Mac is arm64,24GiB RAM,10logical CPUs, modelMac17,2, with about514GiB initially free. Isolated synthetic qualification:1,100JSON log requests (10/s10s and100/s10s), one snapshot viewer, one SSE viewer and active cleanup; all1,100acknowledged records were publicly readable, p95ingest~3.3ms, p95snapshot~7.6ms, peak observedRSS~37MiB and DB+WAL~4.4MiB. Warm build~0.75s; empty compiler cache~8.15s with modules already downloaded. Home Server/Tunnel budgets and qualification remain unknown/unclaimed; see the final verification for limitations and task disposition.

## Owner-approved diagnostic-detail revision — 2026-10-28

The owner explicitly requested useful .NET/Next-like exceptions and custom error messages for app-wide debugging without being the tracked user. This supersedes the initial blanket omission of diagnostic messages/stacks for newly accepted local synthetic evidence. No real app is connected by this decision.

- `scrubbed-v1` retains standard `exception.message`, `exception.stacktrace`, failed-span status descriptions and string bodies on error/exception logs (`log.message` is an internal projection field). Standard OTLP remains sufficient; no producer SDK/schema change is required merely for reception. INFO bodies, arbitrary request content and unknown attributes remain omitted.
- Each field processes at most64 KiB. Retained message is2048 UTF-8 bytes; stack8192 bytes/40lines. Controls/ANSI are normalized; known credentials, cookies/auth, API/token patterns, emails, URL queries/paths and homeusername paths are redacted. Useful function symbols, source basenames and reported line numbers remain where possible. Truncation and receipt-policy version are visible. Regex scrubbing is imperfect: producers must emit controlled diagnostic text, not customer content. No blanket anonymity or perfect PII-removal claim applies to arbitrary text.
- Diagnostic search is owner-protected, app-wide by default, independent of journey correlation, with optional occurrence-time/version/record-definition/scope filters and160-character search. It scans newest10,000 retained records, returns at most100, and marks partial results. Counts are records, not unique incidents or people. Trace/span IDs support technical context; time proximity does not establish cause.
- Details still expire after30 days and are excluded from12-month summaries. App/visit deletion applies through the existing store. Worst-case detail size increases; the existing100,000-record quota is not a new disk/resource qualification. Home Server and full-capacity retention remain unqualified.
- Existing rows retain their original projection. Known legacy replays compare only fields retained by their receipt policy; no text is backfilled or overwritten. New-policy message changes still cause identity conflicts.
