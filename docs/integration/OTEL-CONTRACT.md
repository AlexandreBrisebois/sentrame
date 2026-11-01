# Standard OpenTelemetry and optional client journey conventions

Sentrame receives standard OTLP and interprets retained client observations. One registered app has one write-only telemetry key. No Sentrame SDK, custom wire envelope, producer-role field or server instrumentation is needed. The runtime still includes Sentrame's own Go receiver/backend and private SQLite/JSON state.

## Standard reception

Use the qualified OTLP/HTTP signal endpoints `/v1/logs`, `/v1/traces` and `/v1/metrics` with an ordinary `Authorization: Bearer <telemetry-key>` header. Base-endpoint and signal-endpoint behavior depends on the standard exporter configuration. The key selects the registered app/generation; resource attributes, names, claimed role and IDs cannot change that binding. A multi-product Collector separates requests by app key.

Supported signals/encoding/compression must be checked against the release qualification matrix. Existing local HTTP wire results do not qualify every SDK, gRPC route, direct browser exporter or public Tunnel. Unsupported cells return explicit errors. Profiles are outside the contract. Metrics remain bounded sanitized receipt evidence without unsafe value arithmetic or reconstructed individual paths.

Missing service/version/journey context is not a protocol error. Unknown or privacy-omitted values stay explicit. Native span kind, runtime name, service metadata and trace parentage do not establish physical origin, authenticated installation or persistence. Valid SERVER-kind spans can remain protocol metadata; they do not enable server-workflow journeys or dependency analysis.

## App key and assurance

The app key grants bounded ingestion permission only. It cannot read telemetry, obtain connection settings, mutate definitions or administer apps. Owner sessions and scoped temporary-viewer access remain separate. Lists mask keys; explicit reveal/copy is owner-only; keys never enter URLs, telemetry, ordinary logs or persistent browser state.

A distributed client key is extractable and permits spoofed observations. HTTPS, CORS, a label or a shared ID cannot make them truthful. Whole-app key rotation updates all senders and does not attest installations. A qualified same-origin forwarder carries the same untrusted observations without a stronger telemetry role. No server credential or trusted persistence goal remains.

## Receipt is different from a journey

| Question | Needed evidence | Honest limit |
| --- | --- | --- |
| Did compatible telemetry arrive? | Valid supported OTLP and the app key | Acceptance is not a visit, business result or complete population |
| What operation/error was observed? | Native span/log identity and supported retained structure | Inclusive duration is not page dwell; missing relationships stay unresolved |
| What client path was observed? | Explicit bounded random visit/session and tab correlation | Missing context stays received/unassigned; no person or cross-visit inference |
| Which event has journey meaning? | Owner mapping using service/carrier/name/retained fields | Configuration does not make the observation authentic or prove persistence |
| What was a goal? | Explicit selected observation and captured rule | Goals are optional; a goal-free path has no conversion success/failure |
| How long was a stage? | Supported native span or paired boundaries; optional presence | Timing bases remain distinct; none measures attention |

Use safe explicit random `session.id` and `tab.id` for journeys. Each tab stays separate. Never derive correlation from accounts, household/person IDs, IP, user-agent, fingerprints or nearby timestamps. `session.id` is an opt-in client convention, not a server-operation fallback. A trace alone is not a browsing visit. Explicitly emit approved context; standard propagation/baggage is not automatic event enrichment.

## Named observations and identity

Prefer stable semantic event names such as `screen.opened`, `vote.requested`, `vote.result.observed` and `completion.entered`. These examples are not a mandatory schema. Retain bounded service/version metadata independently of the app identity and definition version. Changed meaning needs a distinct name/step; missing versions remain unknown.

Stable record/event IDs, native trace/span IDs, actions, semantic occurrences and visit correlation are different units. Stable event IDs must be unique within their documented app/generation/signal domain. Retrying an unchanged supported identity is idempotent; conflicting retained content rejects atomically rather than overwriting accepted evidence. Anonymous identical-looking logs may be genuine repeated actions; preserve count uncertainty rather than collapsing by fingerprint.

Unique positive `source.sequence` values within the supported emitter/service scope can order observations. Otherwise occurrence time and stable receipt tie-breaks carry uncertainty. Neither correlation nor display order establishes causality.

## Automatic navigation and configured meaning

An empty captured definition can display goal-free navigation from named non-diagnostic client logs with explicit visit/tab context. A retained `stage.name` can label a step; otherwise use the event name. Diagnostic errors stay supporting evidence, not automatic terminal steps. Native spans and metrics do not become automatic navigation steps.

Owners map discovered events with service, signal/carrier, name and supported privacy-projected attribute conditions. Preview matches, missing context and conflicts before activation. There is no browser/server selector. `stage.name` is only a suggestion until an explicit mapping is accepted. Owner mappings win; competing matches block activation rather than relying on hidden specificity. Child evidence does not inherit a stage merely through parentage.

A configured journey starts on its first eligible mapped event. Its optional goal closes it; another genuine mapped occurrence can start another journey. Thirty minutes without eligible activity closes inactive with unknown outcome. A goal-free journey does not contribute completed-goal or conversion counts. Late evidence reconciles under the captured definition and visible revision. Editing a rule or accepting a hint affects future captures rather than relabeling existing journeys.

## Native structure and timing

Retain supported exact native start/end times, parent identity, kind/status, ordered events and bounded links, plus standard log context under privacy projection. Parent/child, attached events, exact log context and span links are different relationships. Missing or ambiguous targets remain unresolved; links do not join visits or create a causal story. A client span's inclusive duration is available without synthetic boundary events or stage occurrence IDs.

For point events, paired entry/exit stage timing needs the supported occurrence identity within app/visit/tab/emitter scope. Native, measured-elapsed, next-event-estimated and observed-active timings stay separate. Missing/invalid ends are censored or unavailable, not zero/negative completed samples. Do not sum overlapping inclusive durations as unique visit time. Actual page interval and scoped visibility/idle observations are required for observed-active time under the declared 60-second cap; absence means unknown, never attention.

## Diagnostics and privacy

Standard error severity, exception attributes and qualified failed-span detail remain inspectable. Link diagnostics to journeys only through supported explicit context or verified relationships; nearby time is not cause. Warning/error alone does not terminate a path or prove failed business work. Later observed recovery and goals remain visible. A shared stack groups similar failures but does not prove one occurrence or one person.

Projection precedes storage, preview, SSE and copy. Send controlled technical messages, not credentials, personal content, raw request bodies or arbitrary URL/query strings. Unknown attributes and ordinary informational bodies remain omitted. Message/stack bounds and privacy policy remain versioned. Pattern scrubbing does not guarantee perfect anonymity; missing source maps or omitted fields stay unavailable.

Direct requests validate atomically under bounded request/decompression/record/concurrency/storage budgets. A permanent invalid record rejects the whole request without partial commit; retryable storage/overload failure is not successful persistence. Collector queue acceptance and Sentrame's durable commit are distinct acknowledgement boundaries. Exporter batching, sampling, retries/loss, occurrence time, receipt time and display latency remain distinct.

## Lifecycle and compatibility

Ordered detail retains its existing 30-day lifecycle; supported text-free summaries retain their existing twelve-month lifecycle. Summary-only periods cannot reconstruct paths. Detail reads, linked navigation and copies obey app/session authorization, expiry, revocation and deletion. Links and replay never extend retention. A sample/query/display cap is not a claim of complete application traffic.

Existing two-key data is unsupported. Start with a fresh data directory and register the app for a new key; older state is not reset or migrated.

Use [the developer guide](../DEVELOPER-GUIDE.md), [agent checklist](AGENT-INTEGRATION.md) and [app workspaces](../MULTI-APP.md). The [historical qualification ledger](../verification/standard-otel.md) describes earlier tested cells and limitations; source snippets and wire fixtures are not current real-browser, deployment or server support claims.
