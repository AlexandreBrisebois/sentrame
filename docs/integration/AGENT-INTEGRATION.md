# Sentrame integration instructions for coding agents

Use these instructions when a developer asks you to prepare or add Sentrame journey evidence to **their application**. They are separate from Sentrame's own [engineering harness proposal](../proposals/sentrame-engineering-harness.md). Sentrame stays deterministic; discovery assistance runs in the developer's coding agent, not a new AI runtime service.

This guide is a planning deliverable. It has not instrumented any application or installed a package. Examples are semantic plans, not tested runnable SDK recipes. Read the application's own instructions and current SDK documentation before choosing real APIs.

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

Propose semantic entry/exit or start/end pairs for recipe viewing, deciding on one vote, or a server operation stage. Keep each repeated occurrence distinct using safe occurrence context, not recipe/person identity. Duplicate delivery must not create another measurement where supported identity permits deduplication.

Distinguish measured elapsed, observed-active and estimated timing. Active viewing needs actual visibility/idle/tab evidence; elapsed alone can include background time. Server-only elapsed timing does not imply human attention. Next-event fallback is explicitly estimated under a defined policy. Missing ends are ongoing/right-censored, not zero. Incompatible clocks, clock jumps, ambiguous overlaps and ends received before starts must reconcile or be unavailable, never produce invented negative durations. Do not sum overlapping stages as a unique visit total.

The accepted display is **time at stage**: median plus supported slower cases, node size for duration, edge thickness for traffic, a bounded legend and numeric/table access. The slower percentile, support threshold, idle/fallback policy and scale are technical defaults to review, not silently selected by this guide. Live ongoing time is provisional and separate from completed distributions. Neither larger nodes nor longer dwell proves cognitive effort, frustration or a defect.

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
