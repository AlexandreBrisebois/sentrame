# Sentrame integration instructions for coding agents

Current multi-app setup: register/select the app in the owner UI, obtain generated connection settings, and configure the existing standard exporter/provider. Credentials route the app; no Sentrame SDK, app-ID attribute or instrumentation rewrite is required. Keep server credentials out of browser recipes. Journeys is the default; explicit Traces remains available. Greenfield installation does not import legacy data/credentials. See [app workspace guidance](../MULTI-APP.md). Synthetic fixture tools register/select only their demonstration app through owner APIs and seed synthetic mappings only there; a running service and owner passphrase are required.

Use these instructions when a developer asks you to prepare or add Sentrame journey evidence to **their application**. They are separate from Sentrame's own [engineering harness proposal](../proposals/sentrame-engineering-harness.md). Sentrame stays deterministic; discovery assistance runs in the developer's coding agent, not a new AI runtime service.

This guide is a planning deliverable. It has not instrumented any application or installed a package. Examples are semantic plans, not tested runnable SDK recipes. Read the application's own instructions and current SDK documentation before choosing real APIs.

## Standard telemetry and application interpretation

Read the [standard OTel versus Sentrame convention boundary](OTEL-CONTRACT.md#what-is-standard-and-what-is-sentrame-specific) before proposing fields. OTel spans/logs/events and applicable semantic conventions remain the producer foundation. Sentrame Kind **activity**, `stage.name`, `stage.occurrence.id` and `activity.state` are optional application conventions, not official OTel user-journey semantics or mandatory reception fields.

Use the contract's [activity decision examples](OTEL-CONTRACT.md#activity-examples-when-to-use-each-meaning) and [stage-value naming convention](OTEL-CONTRACT.md#naming-convention-for-stagename-values). Choose a graph step versus supporting activity according to the developer's question; do not call every click a .NET Activity, manufacture attention from a click, or require stage tags on ordinary logs/spans. Existing valid standard telemetry must remain accepted even when journey context is missing.

## 1. Discover the app before proposing emitters

Establish the developer's authorized scope first. Read relevant source, navigation/UI states, tests and current instrumentation; preserve unrelated work. Identify existing OpenTelemetry SDK versions, logging/event APIs, resource detectors, exporter/processor configuration, emitted signals, sampling/batching, standard propagation and deployment/consent controls. Reuse them. Do not replace the SDK or introduce a mandatory Sentrame client, custom telemetry protocol or generic wrapper framework.

Trace a few meaningful paths from entry through decisions, repetition, completion and recovery. Inspect state/persistence boundaries rather than assuming a click or HTTP status proves success. Inspect catch/error handling, retry exhaustion, fallback, validation, cancellation and recovery boundaries, including places where a useful warning/error does not yet exist. Check where the app knows an outcome versus where the UI merely anticipates it.

Ask small, adaptive questions that the source cannot answer. For example, ask what success means for one path and which observation would change a product decision; then ask how expected loops or degraded-but-successful recovery should appear. Do not send a fixed questionnaire or ask the developer to repeat facts already verified. Do not instrument everything.

## 2. Present a prioritized event plan before edits

For each proposed point of interest, show:

| Field | Required explanation |
| --- | --- |
| Question/value | What developer/product question does this event answer, and why is it worth the volume? |
| Evidence/location | Observed source/test/UI boundary and the smallest emitter location; distinguish verified source behavior, developer report and a hypothesis needing confirmation. |
| Meaning/authority | Intent, committed outcome, rendered screen, explicit exit, warning, failure, fallback or recovery; who actually owns that fact? |
| Standard carrier/name | Existing standard log/event/exception or qualified span-event API; stable semantic event name, SDK/version constraints and standard exporter path. |
| Context | Minimum safe app-scoped correlation/occurrence context and optional technical enrichment; never person, household, content or long-lived device identity. |
| Severity/noise | Why warning versus error versus ordinary information; expected frequency, retry reporting owner, duplicate/sampling policy and cardinality bound. |
| Privacy/test | Allowed fields, redaction/content risks and synthetic assertions through the real exporter/receiver/read path. |

Prioritize a small useful plan. Explain uncertainties in text; no confidence badge or automated root-cause claim is needed. The developer reviews/approves the event plan before code edits. Existing authorization counts when it explicitly covers that concrete plan; otherwise obtain it. A plan approval is not deployment or production-experiment approval.

## 3. Separate standard intake from sufficient journey evidence

**Base reception:** valid supported OTLP can ingest with standard SDK/exporter configuration and configured app attribution, without custom visit/event fields, a Sentrame SDK or both a browser and server. Missing enrichment is not a protocol error.

**Journey interpretation:** enough context is needed to group and order the intended logical journey. Propose that context through the developer's existing SDK and agreed safe field mappings; missing context stays unassigned/partial rather than being invented.

| Semantic context | Recommended meaning, not a mandatory wire schema |
| --- | --- |
| App/source | Registered app and trusted producer namespace come from permitted configuration; resource claims and correlation IDs do not grant authority. |
| Logical correlation | Random app-scoped ID for a bounded visit or server workflow/operation. Client-only, server-only and combined journeys are supported. Never reuse it as a person or cross-visit key. |
| Browser/tab | Use a fresh per-tab random visit scope for browser journeys. A server-only operation needs no browser or tab. Combined events carry the same correctly scoped logical correlation while preserving producer/tab context where applicable. |
| Event and occurrence | Stable semantic event name; optional supported stable record/occurrence identity for replay and repeated-stage pairing. Record ID, action ID, journey correlation and trace/span IDs are different concepts. |
| Goal | Owner chooses an emitted event from the discovered list. A screen observation is not proof of server persistence or of who caused shared agreement. |
| Timing | Source occurrence/observed time and safe start/end or entry/exit occurrence context. Monotonic measured duration can help where the existing API supports it; document clock basis. Receiver arrival time is separate. |
| Technical enrichment | Existing service/runtime/language/release/deployment fields when available; missing values do not block intake. No visit IDs as metric labels. |

Use existing [OTel context propagation](https://opentelemetry.io/docs/concepts/context-propagation/) and supported context APIs to carry approved correlation between the app's components. Do not assume one HTTP trace covers the full business journey, that baggage is automatically copied into every event, or that a shared ID establishes order, duplicates or causality. [Baggage guidance](https://opentelemetry.io/docs/concepts/signals/baggage/) matters: restrict propagation to permitted app boundaries and avoid forwarding identity/secrets to unrelated services. Received correlation remains untrusted input, never producer authentication.

## 4. Map emitted events to journeys and reusable views

Use simple, stable semantic names. Keep the same step name across releases only while meaning remains the same; changed meaning gets a new name/definition version. Do not embed recipe IDs, URLs, free-text content or ever-growing dynamic values into event names.

Automatic boundaries start at the first mapped event, end at the active goal and start another journey at the next genuine mapped occurrence. Thirty minutes without mapped activity closes inactive, not failed. Browser tabs stay separate; server operations use their bounded logical scope. Supported duplicate goal records are deduplicated before segmentation. Late/out-of-order evidence uses the documented event-time/revision policy, not arrival order.

The owner can save goal-backward, beginning/end-segment and origin-forward views. These are queries over recorded boundaries, not new instrumentation or permission to resegment history. If selected endpoints cross boundaries, show those boundaries. Missing/expired beginnings remain partial; never invent a direct edge across hidden/missing steps. A saved view does not retain raw detail beyond its 30-day lifecycle; 12-month summaries cannot reconstruct individual paths.

## 5. Instrument time at stage honestly

First reuse a native span for a duration-bearing operation. No paired logs, stage attributes or custom occurrence ID are needed for that native operation. Propose semantic entry/exit pairs only for a broader page/conceptual interval without a sufficient encompassing span. Keep each repeated occurrence distinct using safe occurrence context, not recipe/person identity. Duplicate delivery must not create another measurement where supported identity permits deduplication.

Distinguish measured elapsed, observed-active and estimated timing. Active viewing needs actual visibility/idle/tab evidence; elapsed alone can include background time. Server-only elapsed timing does not imply human attention. Next-event fallback is explicitly estimated under a defined policy. Missing ends are ongoing/right-censored, not zero. Incompatible clocks, clock jumps, ambiguous overlaps and ends received before starts must reconcile or be unavailable, never produce invented negative durations. Do not sum overlapping stages as a unique visit total.

The accepted display keeps timing in inspectors and accessible numeric tables, separated by native-inclusive, paired elapsed, estimated and observed-active basis. Journey edge thickness represents traffic; trace parent/event/log/link relationships retain their own meaning. The slower percentile, support threshold, idle/fallback policy and scale are technical defaults to review, not silently selected by this guide. Live ongoing time is provisional and separate from completed distributions. Neither native duration nor longer observed dwell proves cognitive effort, frustration or a defect.

## 6. Propose missing warnings/errors and recovery evidence

Ask which degraded states matter and inspect the paths that know them: failed operations, exhausted retries, degraded fallbacks, unexpected validation/service failures and successful recovery. Propose new warning/error events where useful, not only hooks around existing thrown exceptions. Explain the question, evidence, appropriate severity, reporting owner, safe context and synthetic validation for each proposal.

Use standard OTel/logging severity and exception conventions. Warning/Exception can be discoverable UI classifications; they are not new OTLP signals. A warning or exception is evidence at a journey point, not an automatic failed or ended journey. Preserve subsequent activity, recovery and eventual goal success. Expected cancellation, a rejected user choice, ordinary invalid input and healthy repetition are not automatically errors; classify them according to the actual operation contract. [OTel error-recording guidance](https://opentelemetry.io/docs/specs/semconv/general/recording-errors/) is a reference, not permission to change app behavior.

Do not throw an exception merely to emit a warning, swallow or rethrow differently to manufacture telemetry, change retry/fallback behavior, or log unbounded raw messages/stacks. Prefer safe type/reason codes and allowed diagnostic fields. Decide one reporting owner or supported occurrence identity so logging plus span events do not double-count one exception. A stack fingerprint groups similar failures; it does not prove one occurrence. If sampling/identity is insufficient, disclose count uncertainty rather than claiming exact totals. A new recovery event should state the recovered operation/goal under safe correlation without asserting root cause.

## 7. Illustrative WFS discovery plan

Static source research found optimistic Discovery UI changes before background vote requests, server persistence before later publication, and a shared convergence/completion display. These are source observations, not proof of emitted integration telemetry. The developer reports that the app moves to a completion screen. Ask which observation supports the selected product question, then propose the smallest plan:

| Proposed observation | Question and minimal boundary | Important distinction |
| --- | --- | --- |
| Discovery entered | Where do mapped journeys begin? Emit at an actual entry transition. | A render loop is not repeated entry or intent to vote. |
| Vote requested / vote committed | Did an interaction become a committed contribution? Intent at UI action; commitment at the owning server commit boundary. | One action can have several technical records; do not count intent and commit as two votes. |
| Recipe-stage entry/exit and per-vote occurrence timing | How long was each observed stage? Emit actual boundaries with safe random occurrence scope. | No recipe ID/content is needed; repeated voting is healthy. |
| Completion screen reached | Did this visit observe completion? Emit the real screen/state transition once per supported occurrence. | Another participant may cause shared convergence; observing it does not identify the cause or create another household agreement. |
| Explicit exhausted/agreed/departed reason | Why did the app explicitly report voting ended? Emit only where known. | Silence is not tiredness or abandonment. |
| Post-commit publication failure / recovered delivery | Did a committed vote encounter downstream friction, and did it recover? Propose warning/error at the owning publication boundary and explicit recovery when known. | The vote remains committed even if the client sees an error; do not require a second vote to prove it. |

Possible event names such as `wfs.discovery.vote_committed` and `wfs.discovery.completion_screen_reached` are illustrative proposals, not required standard names or existing emissions. Select/approve actual names and locations with the developer. Pinned [Discovery UI](https://github.com/AlexandreBrisebois/whats-for-supper/blob/0dc61dfce9a7bf7d55eafff08f57c10e3f237f7c/pwa/src/app/(app)/discovery/page.tsx), [vote service](https://github.com/AlexandreBrisebois/whats-for-supper/blob/0dc61dfce9a7bf7d55eafff08f57c10e3f237f7c/api/src/RecipeApi/Services/DiscoveryService.cs) and [Discovery specification](https://github.com/AlexandreBrisebois/whats-for-supper/blob/0dc61dfce9a7bf7d55eafff08f57c10e3f237f7c/openspec/specs/recipe-discovery/spec.md) ground this example; inspect the app's current revision before any future edit.

## 8. Validate synthetically and report limits

After developer approval, use that application's test-first workflow and minimal changes. First inspect the receiver's actual qualified support matrix, SDK version/configuration, app mapping and route. Standard OTLP/HTTP and origin gRPC support are separate from shared Cloudflare Tunnel compatibility; do not assume browser gRPC or public-hostname tunnel gRPC. No secrets or production experiments are needed to prepare the plan.

Propose synthetic cases for base intake with missing journey fields; one client-only, server-only and correctly correlated combined path; cross-app/tab isolation; healthy repeats; goal/next-event restart and 30-minute inactivity; duplicate/late/out-of-order events; warning/error then recovery/goal; post-commit error; stage visibility/missing ends; saved-view partial evidence; 30-day detail/12-month summary boundaries; and privacy/producer-namespace checks. Use real standard exporter/receiver/storage/read seams when separately authorized; a mocked UI or direct database fixture is not enough to qualify integration.

Report tested SDK/config/content identity and passed, failed, blocked, not-run or not-applicable checks. Do not claim compatibility, privacy certification or production readiness from a plan. Stop at the approved scope; request separate developer authorization for deployment or real collection where it is not already granted. This guide performs none of those actions itself.


## Progressive voting checklist

1. Inventory existing providers, logging bridges, auto-instrumentation, exporters, sampling, propagation and shutdown. Reuse one initialization path; do not add a proprietary wrapper or duplicate SDK/provider.
2. Trace actual navigation → intent → processing/persistence → failure/retry → commit → recovery boundaries in source/tests. Ordinary operations/errors come first. An HTTP request is not browser page dwell, intent is not persistence and a publication error after commit does not undo the saved vote.
3. Verify actual receipt/storage/protected discovery and native trace structure/timing. Then configure service/carrier/source/name/retained-field mappings. Require a trusted named commitment or documented successful operation for a persistence goal, not generic HTTP2xx.
4. Optional `stage.name=voting.persist` on an existing span/named log only offers a label suggestion. Explicit **Vote saved** wins over `voting`. Accept into a literal prospective definition; no child inheritance, automatic Kind/goal/authority/session or historical remapping. Stable lowercase dotted labels contain no content, routes or IDs. Both [language additions](../DEVELOPER-GUIDE.md#optional-label-enrichment-in-existing-emitters) reuse the existing provider.
5. Add explicit random scope/tab only for the multi-request question; server-only workflows need no fabricated browser tab. `session.id` is an opt-in Development client convention with a separate Sentrame server-workflow interpretation. Approved propagation must explicitly emit it; baggage and shared trace IDs do not automatically create journeys.
6. Add real page boundaries/presence only if that question remains unanswered. Visibility/idle observations under the60-second cap do not establish attention. Keep browser forwarding under browser authority, review collection/CSP/CORS/lifecycle loss and keep server tokens out of bundles.
7. Exercise missing evidence, replay, late arrivals, privacy, browser forgery, failed attempt/retry/commit and recovery through the actual exporter → receiver → SQLite → protected read path. Record exact versions/route/protocol/compression and batching/flush behavior. Wire fixtures, source review and mocks cannot qualify .NET, Next.js, browser or deployment routes. No external application instrumentation, collection or deployment follows this checklist without its own authorization.
