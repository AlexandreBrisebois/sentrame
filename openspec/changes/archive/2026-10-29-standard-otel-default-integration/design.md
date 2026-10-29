# Design

## Context

See [proposal](proposal.md) for motivation. This file describes proposed behavior, not current release support. Repository inspection on 2026-10-29 found:

| Area | Available in current source | Gap addressed here |
| --- | --- | --- |
| `internal/app/intake.go` | Maintained OTLP receiver; logs, spans, flattened span events; service/version, trace/span IDs; safe projection; credential-derived source | Span end, parent ID, kind and links discarded; status partly reduced to diagnostics; no retained structural trace explorer |
| `types.go`, `evidence_json.go` | One evidence projection; decimal-string time boundaries | Need exact end time and bounded structural fields, explicit projection/computation versions |
| `journey.go` | Name/source/one allowed-attribute mapping; logical `session.id` scope; browser requires `tab.id`; versioned lifecycle | No scope means unassigned; native spans do not supply timing; stage pairs require `stage.occurrence.id` |
| `query.go`, `store.go` | Owner configuration, preview, prospective definitions, immutable receipt version, saved references, durable transaction, replay conflicts | No service/signal selector in mapping; no trace-specific projection/query mode; allowlist limits selectable attributes |
| `web/src/App.svelte` | Discovered names, contextual stage/Kind/condition controls, Errors, graph and table | Names conflate carriers; no native operation tree or explicit typed links |
| Guides and canonical baseline | Ordinary OTLP receipt distinguished from journey interpretation; source-checked .NET/Next snippets | First-use workflow still escalates to custom journey/timing context; several planning-era claims are stale |

No repository AGENT/AGENTS instruction file was found in the workspace. At initial inspection canonical specs were empty and `synthetic-journey-outcomes` was active. During final checks, concurrent work archived it to `openspec/changes/archive/2026-10-29-synthetic-journey-outcomes/` and populated eight canonical capabilities. Requirement blocks were compared and are identical to the inspected baseline. This proposal now targets those canonical owners without changing them. The archived design's “documentation only/no application” statements are historical, not current implementation evidence. Read later decision-register entries and implementation contracts together. Current dirty contract/developer/agent/index edits and untracked OTel contract are inputs, not files this planning pass may replace.

## Goals / Non-Goals

**Goals:** Useful first receipt without proprietary fields; explainable native operation timing; owner-controlled business interpretation; purposeful optional enrichment; bounded private structural evidence; coherent docs and reproducible qualification.

**Non-Goals:** Automatically converting a trace into a visit; attention or page dwell inference; retaining arbitrary attributes; arbitrary queries; long-lived user/session joins; multi-app UX or the rest of roadmap B. Per owner clarification this is greenfield with no users: do not design legacy schema readers, old API aliases, historical replay-policy adapters or dual modes. Keep the useful optional application conventions as forward product choices.

## Decisions

### 1. Progressive capabilities, with independent missing-data behavior

These levels describe evidence prerequisites, not a requirement to enable all five. L3 can interpret point events from L1 without spans; L4 does not require a browser; L5 does not prove outcomes.

| Level | Minimum evidence | Resulting behavior | Missing data and limitations |
| --- | --- | --- | --- |
| L1 Receipt/discovery | Authenticated supported OTLP and configured app binding | Services, span/operation names, named logs/span events, errors and retained field availability; accepted-unassigned inspection | Missing service/name shown unknown/unnamed; omitted fields explained; metrics receipt-only; no invented journey |
| L2 Trace-scoped operation exploration | Valid trace/span IDs; native span timestamps for duration; optional parents/events/correlated logs | Partial operation tree, typed relationships, inclusive span durations and diagnostics with no `stage.*`, visit or event IDs | Missing parent is unresolved; missing/invalid end means unavailable/right-censored if a start is actually known; logs alone give no span duration; SDKs usually export only ended spans, so silence does not reveal open work |
| L3 Business interpretation | Owner-activated selector mapping over retained names/fields; trusted evidence appropriate to chosen goal | Stages, supporting activities and configured outcomes inside an explicitly labeled trace or journey scope | No mapping means technical evidence only; missing outcome evidence stays unknown; labels do not create dwell intervals or multi-request paths |
| L4 Multi-request journey | Explicit random bounded app scope; browser tab isolation; correctly propagated correlation on each participating request/event | Existing first-mapped/goal/30-minute lifecycle over supported correlated evidence | No scope means no multi-request journey; conflicting/incomplete scope remains unassigned/ambiguous; trace, user ID, IP, agent string or time proximity is insufficient |
| L5 Browser presence | L4 browser scope/tab plus actual visibility/idle observations and a supported stage interval | Separately named observed-active intervals, policy and coverage | Missing presence does not block L1–L4; hidden/idle/unobserved intervals never imply attention; no backend-derived dwell |

The default landing workflow starts in discovery, opens an operation trace where available, then offers contextual “give this business meaning.” It does not start by asking for a session ID. Existing errors remain useful without any journey.

### 2. A minimized structural projection, not raw OTLP storage

Extend the retained evidence contract with native span end time, parent span ID, span kind/status code, trace flags, ordered span-event index and bounded link target IDs. Preserve occurrence/observed/receipt distinctions and decimal nanosecond strings through persistence, protected reads and copy. Link attributes and trace state remain omitted; no baggage wholesale retention. Keep log trace/span association distinct from a parent edge. Preserve service and producer namespace on every record.

Preserve available SDK dropped-event/link/attribute counts separately from Sentrame projection omissions. A sampled flag does not establish a sampling rate or completeness. Add no automatic inverse-sampling weighting.

| Convention / native field | Standard status or boundary | Current Sentrame support | Proposed change |
| --- | --- | --- | --- |
| Trace/span IDs, log trace context, span events | Native OTel model | Retained, flattened events | Preserve typed structure and event index |
| Span end/parent/kind/links/status | Native OTel model | End/parent/kind/links omitted; status reduced to error text/severity | Minimized structural retention and qualified timing |
| `service.name`, `service.version` | Standard service metadata | Retained | Contextual discovery/selectors; no authority from values |
| HTTP method/route/status | Standard HTTP semantic attributes | Omitted by fixed projection | Explicit safe supported selectors with privacy tests |
| `session.id` | Opt-in Development client convention | Supported record-level scope; server use is Sentrame interpretation | Retain optional use with stability/scope caveats |
| `stage.name`, `stage.occurrence.id`, `tab.id`, `activity.state`, `action.id`, `event.id` | Sentrame-supported application conventions here | Projected for dependent capabilities | Retain purposefully; none required for native operation exploration |
| `session.previous_id`, arbitrary renamed correlation keys | Cross-visit linkage or unsupported app conventions | Not supported | Remain omitted; no inference or renamed mandatory schema |

Discovery exposes signal + carrier + service + name + supported fields so `POST /votes` as a server span differs from a similarly named browser log. Safe name projection remains bounded; unsafe/dynamic names require sender-side minimization and display an omission reason rather than storing raw URLs. Initial mapping predicates use exact service name, trusted source, signal/carrier, span/event name and one retained key/value. Extend the safe selectable set to `http.request.method`, route-template `http.route`, integer `http.response.status_code`, plus existing bounded fields; route templates pass explicit privacy validation, not raw URL retention. No URL query/path, request headers, database statement/parameters, IDs of people or business content. Standard attributes still require receiver implementation and privacy tests before becoming selectable.

Proposed bounds: retain at most 32 span links per span and report dropped/omitted relationship counts; reject no otherwise valid record solely because optional links were truncated. Keep current 2,000 projected records/request, 1 MiB request, 4 MiB decompression, 64-attribute boundaries and 10,000 read scan/2,000 evidence return limits unless measured work justifies a separately reviewed adjustment. Span events continue to count as records. A limit makes completeness partial and cannot imply a root, ending or direct path. Use maintained receiver decoding; this plan adds no transport parser.

Alternative rejected: retain every SDK attribute to enable arbitrary future mappings. It undermines privacy and resource bounds. Alternative rejected: synthesize start/end event rows for every span; it duplicates native evidence and changes event/path counts.

### 3. Two graph meanings, explicit relationship types

**Trace exploration:** nodes are individual span instances keyed by app + trusted producer namespace + trace/span identity. Parent edges use actual parent IDs within the trace. Span events are attached point evidence in original export order with their timestamps; correlated logs attach to exact span context, or to trace context only when span is absent. Links remain separately styled references that can cross traces; they do not merge trees, sessions or journeys. Missing targets stay unresolved. Never infer parentage from time adjacency, label equality or containment. Browser and server records can be displayed together under the same app trace context, but source authority remains separate; multiple matching parent candidates are ambiguous. Claimed parent IDs are observed structural context, not authentication or proof of a trusted business fact.

**Business path:** stages group configured semantic occurrences. Edges are chronological observed transitions within a captured journey, preserving gaps, repeats and boundary markers. Within a trace, a configured business overlay labels occurrences on the tree; it is not a path count or business visit. Parallel siblings never become A→B just because of sort order. Journey timeline order remains stable event time/source sequence with disclosed cross-clock uncertainty, not proven causality. Nested spans retain both levels of work; asynchronous linked spans stay related references, not an invented linear request. Hierarchy/force/Sankey layouts cannot change the semantic model. Trace mode has no Sankey conversion funnel; use tree/timeline and accessible table. Journey mode retains the current three layouts.

Typed read results and UI mode must distinguish `trace` from `journey`. Proposed protected `GET /api/traces` lists bounded traces and `GET /api/traces/{traceId}` returns bounded structure, evidence, limits and revision under the same app authorization. Definition/filter context accompanies the business overlay. Existing snapshot remains the journey read. New routes enter the explicit temporary-viewer read allowlist; owner-only configuration and same-origin/CSRF mutations remain enforced. SSE invalidations refresh snapshots rather than incrementing client counts.

Alternatives rejected: reuse chronological journey edges to draw traces, or turn every request into a session. Both create false order, dwell and conversion claims.

### 4. Native operation timing with explicit business semantics

Inclusive operation duration is end minus start for a qualified native span. Use the span's trace/span identity within app/producer scope as its occurrence identity. Neither `stage.occurrence.id` nor paired custom logs are needed. Do not calculate exclusive/self time or a critical path in this slice. Do not sum nested/parallel spans into unique user time.

An owner can map an existing span to a stage and select native-span timing: that stage occurrence measures the operation represented by that span only. A page GET span is request processing, not “time looking at Voting.” To measure a broader conceptual stage requires an actual encompassing span or explicit matched boundaries. Keep optional paired stage events for that case, with occurrence identity only when native identity is insufficient. One occurrence chooses one explicit timing basis; span and pair evidence are not double-pooled. Default for a mapped span is native-span; for a point event, timing is unavailable unless an explicit pair or opt-in estimate exists. Estimated and observed-active distributions remain separate. Carry basis + policy/computation + definition versions into summaries; p90 still requires 20 completed compatible samples.

Invalid/missing ends and missing clocks never become zero or negative completed durations. A present exported start without end is censored evidence, not a promise of live start export. End-before-start wall-clock evidence is timing-unavailable. No extrapolation from a backend request into browser dwell or attention. A goal on a span evaluates at the selected completion/end boundary; an event goal evaluates at that event time. A configured persistence goal must match an authoritative named fact or an operation whose documented success contract proves persistence; generic HTTP 2xx or a finished span is insufficient.

### 5. Optional stage.name simplifies configuration

`stage.name` stays a Sentrame-supported application convention, not an official OTel semantic convention. It supplies only a bounded semantic label/hint. For an unmapped selector, offer the distinct emitted label as a suggested Stage. Accepting the suggestion produces an ordinary owner mapping with a literal label and explicit conditions; it is not a dynamic rule that trusts future tag values. No tag means the owner can type a label against the span/event name. A tag never selects Kind, goal, boundary, session, authority or timing quality.

Precedence is predictable: exactly one explicit owner rule wins over any tag; two overlapping explicit rules block activation until disambiguated (no hidden specificity/order); zero rules means unassigned business meaning plus optional suggestions. A mismatching tag remains inspectable as emitted evidence without renaming the owner's stage. Multiple tag labels for one selector appear as separate conflicting suggestions; do not choose the most frequent one. Empty, invalid or absent values give no suggestion. Child spans/events/logs do not inherit stage tags merely by correlation.

Example: `vote.committed` emits `stage.name=voting`; owner maps it to `Vote saved`. The stage is **Vote saved**, and discovery can disclose the mismatch. Accepting `voting` elsewhere affects only a newly activated definition. Span/event name semantics and source rules remain explicit. Prefer short stable lowercase dot-separated values such as `voting`, `voting.submit`, `voting.persist`; avoid dynamic IDs, routes, customer text and `started`/`ended` suffixes that confuse labels with lifecycle. Preserve this guidance and concrete .NET `SetTag`/logging-scope and JS `setAttribute`/log-attribute examples in both recipes.

Alternative rejected: auto-create active stages whenever a tag appears. It lets producer-controlled content change product configuration and historical meaning.

### 6. Correlation remains optional and purposeful

Use standard W3C trace propagation for technical execution. It is not a multi-request business correlation protocol. Retain current supported record-level `session.id` for bounded logical scope and `tab.id` for browser separation. OTel's client session convention is opt-in and Development; Sentrame's additional server-operation use and stricter no-cross-visit policy are application-specific. Never imply the SDK automatically emits/propagates it or that baggage automatically becomes record attributes. Explicit approved propagation and emission are required. No `session.previous_id` linkage.

No new general-purpose correlation-key mapper is needed here. An unfamiliar attribute can remain omitted until explicit support/privacy design is separately justified. Span links preserve technical associations but do not automatically join sessions. The same valid trace context can relate logs to spans without custom action IDs; ambiguous log-to-journey association remains unassigned, especially when several scopes share context. Never promote browser context by forwarding with a server credential. Use a bounded app-origin OTLP forwarding route authenticated to Sentrame as browser, separately from trusted server emitters.

### 7. Historical behavior is a forward product contract

Record projection policy at receipt, definition at activation/capture, and computation/timing policy with derived data. All new definitions are prospective, including accepted tag suggestions. Existing captured journeys keep their definition until their known goal/inactivity closure; new mapped occurrences use the then-active definition. Trace overlays expose the selected definition, with capture-time interpretation as default and an explicit non-persisting preview for alternatives. Late evidence revises only affected retained structures under the captured definition and visible data revision; it cannot switch the goal policy silently. Saved references pin mode and relevant definition/boundary filters, never store raw detail beyond retention.

For the new model, exact span/event replay is idempotent by app/producer/signal identity; event index belongs to the exported span. Changing an accepted span under the same identity is a conflict, not an upsert that invents a missing end. Normal SDK spans export when ended; progressive span mutation support is out of scope. Anonymous logs without reliable IDs remain distinct with retry uncertainty. No cross-signal fingerprint deduplication. Sampling/drop counts, missing parents, truncation, clock skew and incomplete trace boundaries stay visible. Deletion/30-day expiry removes detail and relationship navigation; links/references cannot resurrect it. Twelve-month summaries stay text-free, identity-free and separated by policy/basis; no new long-lived trace membership index.

## Risks / Trade-offs

- [Ordinary auto-instrumentation may expose only request-level facts] → First value is operation/error exploration; require additional named outcome evidence only when the selected business question lacks it.
- [Safe projection limits selectors] → Publish retained-field availability and sender-side minimization; never silently open arbitrary retention to improve mapping convenience.
- [A partial trace looks complete] → Explicit missing-parent/link/drop/truncation indicators; no population conversions from sampled operations.
- [Inclusive duration exaggerates total work when summed] → Name the basis; keep overlapping occurrences and decline unique-time totals.
- [Conflicting tag labels or mappings] → Suggestions only, owner authority, activation conflict errors and prospective definition versions.
- [Structural indexing increases work/storage] → Bounded links/reads and repeatable synthetic resource measurements; Home Server capacity remains unqualified.
- [Greenfield schema change makes old local databases unreadable] → Deliberate test-data reset during later cutover, not automatic deletion at startup; no backward-compatibility adapter.

## Migration Plan

Use the greenfield cutover in [compatibility.md](compatibility.md). No migration executes during planning. Build and verify against a fresh disposable database. At later adoption, explicitly stop the old local test instance and select a fresh data directory or deliberately reset the chosen synthetic dataset; never reset on startup. Old binaries/data can be kept together for manual review, but no old/new database interchange is promised. Document schema mismatch as an explicit startup error. Definitions/references created in the new model retain the forward historical guarantees above.

## Review points

No unresolved product choice blocks this plan. Review the chosen defaults before apply: trace tree versus journey path separation; owner-activated label suggestions; operation-only native timing and persistence-goal evidence. Changing any of these warrants a spec revision, not an undocumented implementation choice. Exact SDK/Next runtime versions and Home Server budgets are qualification inputs to record later, not invented support today.

## Official source checks — 2026-10-29

These establish data-model/configuration facts, not Sentrame runtime support:

- [OTel tracing API](https://opentelemetry.io/docs/specs/otel/trace/api/): span timing, parent structure, events and links are separate native concepts. The projection and graph policies above are Sentrame design decisions.
- [OTel log model](https://opentelemetry.io/docs/specs/otel/logs/data-model/): optional trace/span correlation, named events, occurrence/observed timestamps and severity. These support correlation without custom IDs but do not establish a visit.
- [OTel session conventions](https://opentelemetry.io/docs/specs/semconv/general/session/): client session attributes/events are Development; `session.id` is opt-in. Retain the supported key without adopting cross-session links.
- [OTel HTTP attribute registry](https://opentelemetry.io/docs/specs/semconv/registry/attributes/http/): method, route template and response status are standard HTTP selectors; their presence in a standard does not imply current Sentrame retention or safe raw URL collection.
- [OTLP exporter configuration](https://opentelemetry.io/docs/specs/otel/protocol/exporter/): common versus signal endpoints, protocol and headers need explicit documentation and SDK checks.
- [OTel .NET exporters](https://opentelemetry.io/docs/languages/dotnet/exporters/): explicitly select HTTP Protobuf and signal paths for this receiver.
- [Next.js OTel guide](https://nextjs.org/docs/app/guides/open-telemetry): framework instrumentation guidance must be checked against the application's locked runtime; it does not qualify this repository's Node/browser recipes.
