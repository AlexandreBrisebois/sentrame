# Standard OpenTelemetry and optional Sentrame journey conventions

Current multi-app setup: register/select the app in the owner UI, obtain generated connection settings, and configure the existing standard exporter/provider. Credentials route the app; no Sentrame SDK, app-ID attribute or instrumentation rewrite is required. Keep server credentials out of browser recipes. Journeys is the default; explicit Traces remains available. Greenfield installation does not import legacy data/credentials. See [app workspace guidance](../MULTI-APP.md). Synthetic fixture tools register/select only their demonstration app through owner APIs and seed synthetic mappings only there; a running service and owner passphrase are required.

Use this reference to choose what to emit. The [developer guide](../DEVELOPER-GUIDE.md) explains SDK/exporter setup and authentication; the [user guide](../USER-GUIDE.md) explains mapping received records in Sentrame.

## What is standard, and what is Sentrame-specific?

**The transport and telemetry models are standard OpenTelemetry. The journey interpretation is Sentrame configuration.** You can send ordinary SDK-produced logs and spans without adding any of the application attributes below. Do not build a proprietary JSON wrapper or require a Sentrame SDK.

| Layer | Examples | Status |
| --- | --- | --- |
| OTel signal/wire model | Spans, log records, span events, timestamps, trace/span IDs, OTLP/HTTP envelopes | Standard OpenTelemetry |
| OTel semantic conventions | `service.name`, `service.version`, standard exception attributes | Defined by the applicable OTel convention/version |
| Application meaning | “Vote pressed”, “Vote saved”, “Voting page opened” | Your application's observations, emitted through standard carriers |
| Sentrame's optional interpretation | Mapping Kind **step/activity/stage start/stage end**, operational goal, `stage.name`, `stage.occurrence.id`, `activity.state` | Sentrame-supported application conventions; **not official OTel stage/activity semantics** |

In particular, OpenTelemetry does not define a general user-journey Stage or an “activity” signal meaning button clicks/presence. A .NET `Activity` is the .NET tracing API for spans. Sentrame's Kind **activity** is its own display/lifecycle choice, not another SDK carrier.

Custom application attributes can be valid OTel without being standardized semantic conventions. [OTel naming guidance](https://opentelemetry.io/docs/specs/semconv/general/naming/#recommendations-for-application-developers) recommends consulting existing conventions first and using an application/company namespace for new names. This release already recognizes the unprefixed application keys documented here; that is a supported forward application convention, not a recommendation that `stage.*` or `activity.*` are universal standards. Arbitrary renamed attributes are currently omitted by Sentrame's fixed privacy allowlist, so a new namespace needs explicit receiver support before it can replace these keys.

Start with your existing standard instrumentation and map useful names in Sentrame. Add optional context only for a question that cannot otherwise be answered. The wire examples below are **illustrative SDK output**, not a custom envelope your application must construct.

## Three meanings to keep separate

| Term | What it represents | Where it goes |
| --- | --- | --- |
| **.NET `Activity`** / OTel **span** | A technical operation with a beginning and an end, such as processing a request | SDK exports a span to `/v1/traces` |
| **Event** | Something happened at a point in time, such as `vote.committed` | Prefer a named OTel log at `/v1/logs`; a span event attached to a span also arrives through `/v1/traces` |
| Sentrame mapping Kind **activity** | Supporting activity/presence evidence, such as `activity.state=hidden` | An ordinary record, commonly a named log, mapped to Kind **activity**; it adds no graph step |

A .NET `Activity` does not automatically mean user activity. An event does not need a span, trace ID or custom Sentrame envelope. These are standard OTel carriers with optional application attributes. [.NET instrumentation](https://learn.microsoft.com/en-us/dotnet/core/diagnostics/distributed-tracing-instrumentation-walkthroughs) and the [OTel log data model](https://opentelemetry.io/docs/specs/otel/logs/data-model/) define the respective APIs/models.

## Example: navigating to Voting and pressing Vote

Start with existing request spans, child persistence operations and named post-commit logs. A request to Voting is backend request evidence; it does not establish that a page rendered or how long it was viewed.

| Observation | Existing standard evidence when available | Owner meaning / limit |
| --- | --- | --- |
| Navigate to Voting | Server `GET /vote` span | Voting request; no rendered-page or dwell claim |
| Press Vote | Browser `vote.pressed` span or named point event | Intent/supporting activity; unknown without actual browser observation |
| Process / fail / retry | Distinct `POST /votes` and `vote.persist` spans, status/exception/log context | Vote processing with native inclusive duration; repetition remains visible |
| Persist successfully | Trusted named `vote.committed` after actual persistence, or a semantically precise successful commit operation | Vote saved; browser intent or HTTP 2xx alone is insufficient |
| Publication failure / recovery | Actual post-commit error/recovery observations | Preserves friction without undoing persistence |

No stage tag, scope/tab, custom occurrence or presence stream is required for trace exploration or contextual mapping. An optional `stage.name` simplifies literal label entry for business mappings and supplies display labels for unconfigured correlated browser navigation. Add explicit random scope/tab only when the question requires a chronological path across requests; technical trace context and journey correlation remain distinct.

For page elapsed time, reuse a real encompassing browser page span. If none exists, optional `voting.page.entered` / `voting.page.exited` records can use the same fresh `stage.occurrence.id` and configured Stage with start/end Kinds. The press does not end the page interval. Require a separate action ID only for an otherwise unresolved intent/outcome relation. Actual visibility/idle observations enable a separate observed-active basis within the captured browser interval; they do not measure attention. A backend span cannot borrow browser presence.

Choose the operational goal separately. **Vote saved** closes a journey on trusted persistence; a completion-screen goal needs actual browser screen-entry evidence. A later exit cannot complete a pair in a different captured journey, and a native page interval extending beyond the captured boundary has no completed observed-active sample in that capture.

## Optional stage.name suggestions

`stage.name` is a bounded optional **Sentrame application convention**, not an official OTel semantic convention. Put it on the relevant existing span or named log/event; never set a global resource stage for an app that changes stages. Child spans, span events and logs do not inherit a label merely through correlation. Empty, whitespace, non-string or invalid values give no suggestion.

| Explicit owner matches | Emitted hint | Result |
| --- | --- | --- |
| Exactly one | Any value or absent | Literal owner Stage/Kind/goal/conditions remain authoritative |
| More than one | Any value | Activation blocks until explicit selectors are disjoint |
| None | One valid value | Optional suggestion; no mapping activates automatically |
| None | Several distinct values for the selector | Separate suggestions for review; no popularity default |
| None | Missing/invalid | Ordinary name mapping remains available |

The table describes business mapping suggestions. With an empty captured definition, browser logs with explicit session and tab IDs also use retained labels to display navigation without activating a mapping or goal.

Discovery retains at most32 distinct labels per contextual selector and discloses omissions. Accepting a suggestion copies a **literal label** into an ordinary owner draft. Preview/activation is still required and prospective. It does not create a dynamic rule reading future tag values. Example: `vote.committed` emits `stage.name=voting`; the owner's mapping says **Vote saved**. **Vote saved** wins, the emitted hint stays inspectable, and later tag changes cannot rename prior captures. A tag never supplies Kind, authority, goal, timing boundary, session or timing quality.

Prefer short stable lowercase dot-separated labels: `voting`, `voting.submit`, `voting.persist`. Keep spelling/meaning stable across releases; changed meaning needs a different label. Avoid identifiers, routes, customer content, timestamps and `started`/`ended` suffixes. A stage label is not a lifecycle event. Display labels such as **Vote saved** remain the owner's choice and may differ from the emitted value.

For actual API additions to existing .NET and Next.js emitters, see [the optional enrichment recipes](../DEVELOPER-GUIDE.md#optional-label-enrichment-in-existing-emitters). They need no new SDK/provider or proprietary wrapper. Source-checked snippets remain distinct from runtime-qualified exporters.

## Minimums depend on the question

| What you want | Evidence/context needed |
| --- | --- |
| Receive standard telemetry | A valid supported OTLP request and the configured producer bearer token. No journey attributes are required. |
| Discover a useful named event | A stable log `eventName` (or retained `event.name` fallback), span name, or span-event name. Supply source timestamps; missing time uses labeled receipt fallback. |
| Explore an operation trace | Native trace/span IDs, available start/end, optional parent/events/links and correlated logs. No stage/session/occurrence fields. |
| Configure operation business meaning | Owner mapping over service/carrier/source/name and a supported retained condition; no stage tag required. |
| Inspect an error without a journey | Standard error severity/status or exception fields. No visit ID is required. Ordinary INFO bodies are not diagnostic detail. |
| Build a **server-only journey** | Mapped named records sharing a safe random `session.id` for that operation/visit. A browser tab is not required. |
| Explore **browser navigation** without configuration | Named, non-diagnostic browser logs sharing a random `session.id` and per-tab `tab.id`, with an empty captured definition. No mapping or completion goal is required. `stage.name` supplies a label; otherwise the event name is used. Unique positive `source.sequence` values within one source/service supply ordering. Correlated browser errors attach as evidence. |
| Configure a **browser journey** | Mapped records sharing a random `session.id` and per-tab `tab.id`, received using the browser producer credential. A completion goal is optional; an explicit definition controls included steps and labels. |
| Measure a broader page stage with no encompassing native span | Mapped **stage start** and **stage end** records with the same configured Stage, `stage.occurrence.id`, journey context, trusted producer and `service.name`, and usable source timestamps. |
| Estimate observed active time | A complete native browser interval or the optional elapsed pair above, plus actual activity-state observations in the same captured journey, browser producer and service. Unobserved gaps stay unknown; this is not attention measurement. |
| Replay one log without recounting it | A unique string `event.id` per emitted record, preserved on exact retry. Optional for reception; without reliable identity retry uncertainty remains. Spans use their standard trace/span identity. |

Goal-free paths close inactive after 30 minutes without navigation activity. They never count as completed goals or enter goal-conversion denominators, and navigation logs alone do not measure stage duration or attention. Native spans, metrics and server-only logs still require their separately documented interpretation. Identified replay does not add steps.

Sentrame groups journeys using explicit permitted correlation, not a trace ID or nearby timestamps. `action.id`, `event.id`, `stage.occurrence.id` and trace/span IDs answer different questions; none is a person or a viewer credential. See the [optional-context table](../DEVELOPER-GUIDE.md#6-add-journey-context-only-when-you-need-it).

## Minimum named Event: log envelope

With the SDK configured, emit a named log. A .NET example is:

```csharp
logger.LogInformation(new EventId(0, "voting.page.entered"), "Voting opened");
```

For the pinned JS SDK example in the developer guide, the standard log API can emit:

```ts
logs.getLogger('app.journeys').emit({
  eventName: 'voting.page.entered',
  severityNumber: SeverityNumber.INFO,
});
```

The SDK supplies the transport envelope and occurrence timestamp. These snippets assume the provider/exporter setup in the developer guide. They do not create journey context themselves.

The resulting OTLP/HTTP JSON shape is below. **All JSON wire examples here are templates**, not ready-to-post payloads: replace timestamp/identity placeholders with fresh SDK-generated values, or use the repository's [fixture generator](synthetic-demo.md). Auth is the HTTP `Authorization: Bearer ...` header, never an attribute in the body.

```json
{
  "resourceLogs": [{
    "resource": {"attributes": [
      {"key": "service.name", "value": {"stringValue": "example-api"}}
    ]},
    "scopeLogs": [{
      "logRecords": [{
        "timeUnixNano": "<EVENT_TIME_NS>",
        "eventName": "voting.page.entered",
        "severityNumber": 9
      }]
    }]
  }]
}
```

`service.name` and INFO severity are useful recommended context here, not proprietary reception requirements. This record is received/discovered, but without `session.id` it remains unassigned to a journey. To make it a server journey record, add an attribute entry:

```json
{"key": "session.id", "value": {"stringValue": "operation-a"}}
```

Use a fresh random operation value in real instrumentation. For browser records add `tab.id` too and use browser authority. For configured business meaning, map the event to Kind **step**. Unconfigured correlated browser logs can already form navigation paths. A body string, `stage.name`, record ID and trace ID are not required for that graph step once mapping and safe correlation exist.

## Minimum .NET Activity: span envelope

Create a reusable `ActivitySource` and register its exact name with the application's existing tracing provider:

```csharp
// In the existing tracing-provider configuration:
traces.AddSource("Example.Voting");

// A reusable source, normally a static field/singleton:
var source = new ActivitySource("Example.Voting");

// Around the real operation:
using var activity = source.StartActivity("vote.persist");
activity?.SetTag("stage.name", "voting.persist"); // optional hint, not required
// Perform the operation here; disposing the Activity records its end.
```

This is a placement example, not a complete Program.cs. The source should be reused, not allocated per request. `StartActivity` can return null when no listener/sampling decision records it; the SDK/exporter setup and sampling matter. [Microsoft's ActivitySource guidance](https://learn.microsoft.com/en-us/dotnet/core/diagnostics/distributed-tracing-instrumentation-walkthroughs).

A span's useful minimal OTLP/HTTP JSON envelope is:

```json
{
  "resourceSpans": [{
    "resource": {"attributes": [
      {"key": "service.name", "value": {"stringValue": "example-api"}}
    ]},
    "scopeSpans": [{
      "spans": [{
        "traceId": "<TRACE_ID_32_HEX>",
        "spanId": "<SPAN_ID_16_HEX>",
        "name": "vote.persist",
        "startTimeUnixNano": "<START_TIME_NS>",
        "endTimeUnixNano": "<END_TIME_NS>",
        "attributes": [
          {"key": "stage.name", "value": {"stringValue": "voting"}}
        ]
      }]
    }]
  }]
}
```

For spans this receiver requires valid nonzero trace and span IDs; let the SDK generate them. In OTLP JSON these IDs use hexadecimal, while nanosecond integers use decimal strings. A root span does not require a parent ID. [OTLP JSON encoding](https://opentelemetry.io/docs/specs/otlp/#json-protobuf-encoding).

**Current local behavior:** native span start/end and structure support inclusive operation duration and a span-backed business occurrence in Traces. No synthetic stage-end row, paired log or custom occurrence ID is needed. A stage label alone does not measure page dwell. Explicit paired boundaries below are an alternative for a broader interval without a sufficient native span. Native occurrences also contribute to explicitly scoped journey summaries under their own named basis and computation policy; see the ledger.

## An Event attached to an Activity

A .NET `ActivityEvent` is a point inside a span:

```csharp
activity?.AddEvent(new ActivityEvent("vote.committed"));
```

On the wire it sits inside the parent span's `events` array:

```json
{
  "name": "vote.committed",
  "timeUnixNano": "<EVENT_TIME_NS>"
}
```

It is not another `resourceLogs` request. The SDK exports the parent span and its events through `/v1/traces`. Sentrame receives the named span event and retains verified parent trace/span context. Its current projection can inherit parent `session.id`, `tab.id` and compatible `action.id`; it does not generally inherit every parent attribute. Put stage-specific attributes on the event itself when needed.

Choose one intentional carrier for each business fact. Emitting both a log and a span event for the same milestone can create two observations; matching names/stack text do not establish one occurrence.

## Minimal stage-timing sequence

Suppose one visit enters Voting, performs some work and leaves. Emit these named records with usable source timestamps:

| Order | Name | Shared context | Additional attributes | Mapping |
| --- | --- | --- | --- | --- |
| 1 | `stage.entered` | `session.id=visit-a`, `tab.id=tab-a` | `stage.name=voting`, `stage.occurrence.id=occurrence-a` | Stage **Voting**, Kind **stage start** |
| 2 | `activity.changed` | Same session/tab | `activity.state=active` | Kind **activity** |
| 3 | `activity.changed` | Same session/tab | `activity.state=hidden` | Kind **activity** |
| 4 | `stage.exited` | Same session/tab | `stage.name=voting`, `stage.occurrence.id=occurrence-a` | Stage **Voting**, Kind **stage end** |

For a browser sequence all records use the browser credential and the same emitting `service.name`. For server-only stage elapsed timing use one server producer/service and omit the fictitious tab. Timestamps 1 and 4 define measured elapsed time. Activity observations 2 and 3 can additionally support observed-active time inside that interval. **Only 1 and 4 are necessary for elapsed stage timing.**

Each new entry needs a new `stage.occurrence.id`, even when revisiting Voting in the same visit. Entry and exit share that occurrence ID, but each emitted log has its own `event.id` if using record replay identity. A stage end closes that stage's timing, not necessarily the journey: the separately selected operational goal determines journey completion.

Start/end must belong to the same captured journey. Do not select the start event as the journey goal if you intend to pair an end later. Missing ends, repeated ambiguous starts, expired beginnings and receipt-time fallback are not completed measured samples.

## Minimal Sentrame activity/presence record

This is a **log record inside the `logRecords` array** from the Event envelope above:

```json
{
  "timeUnixNano": "<EVENT_TIME_NS>",
  "eventName": "activity.changed",
  "severityNumber": 9,
  "attributes": [
    {"key": "session.id", "value": {"stringValue": "visit-a"}},
    {"key": "tab.id", "value": {"stringValue": "tab-a"}},
    {"key": "activity.state", "value": {"stringValue": "active"}}
  ]
}
```

Use `active` or `visible` to begin/refresh an observed-active interval; `idle` or `hidden` stops it. This minimal record needs no `stage.name` or `stage.occurrence.id`: activity timing is scoped to the captured journey, producer, service and time interval, not selected by those activity-record attributes. If timed stages overlap in that scope, this activity timeline may contribute to more than one overlapping stage; it is not independent per-stage attention tracking.

A positive observation contributes at most 60 seconds until refreshed or stopped. Send appropriate heartbeats while the observed state remains positive if you need a longer observed interval; a single `active` record does not prove 10 minutes of activity. Browser visibility alone also does not prove engagement.

Map `activity.changed` to Kind **activity** to keep it off the graph path. Give the mapping a required Stage label such as `Presence`; that label does not create a node for this Kind. Mapped activity can start a journey and refresh the 30-minute inactivity clock, so choose this deliberately. Activity records merely received without such a mapping do not automatically get that lifecycle meaning. Unknown activity-state values do not contribute active-time calculations.

## Activity examples: when to use each meaning

### A .NET Activity/span: time a real operation

Use a span when the operation has a meaningful beginning and end. These are ordinary OTel tracing examples, independent of Sentrame's mapping Kind:

| Span example | Begin/end boundary | What it tells you |
| --- | --- | --- |
| `vote.persist` | Around the actual save operation | Backend work/latency; emit a separate committed event only after success |
| `recipes.load` | Around loading recipe data | Time spent loading, not time the user spent looking at it |
| `recommendations.calculate` | Around the computation | Technical operation duration and related errors |
| `voting.page.render` | Around server rendering | Render work, not browser dwell time or attention |

Use existing framework/HTTP/database instrumentation where it already supplies the useful span. Do not create a second span merely to duplicate it. Do not hold a backend request Activity open for the entire time a browser page is visible. Native inclusive duration is the operation metric; it is not page dwell or a sum of unique user time.

### Sentrame Kind activity: keep supporting observations off the graph path

These are ordinary named events first. **activity** is a mapping choice you make in Sentrame:

| Observation/event example | When Kind activity is useful | When to choose a visible step instead |
| --- | --- | --- |
| `vote.requested` | Preserve the press as intent while the page-oriented graph stays quiet | When examining attempted-vote paths/counts is a primary question |
| `voting.choice.changed` | Preserve a changed selection within Voting | When selection itself is a meaningful journey milestone |
| `voting.help.opened` | See that help was opened before a later result | When help usage should be an explicit branch in the map |
| `form.validation.rejected` | Retain a validation warning in the current journey without another node | When validation failures form a path you deliberately want to compare |
| `request.retry.scheduled` | Retain retry context around an existing operation | When each retry is deliberately a distinct semantic step |

These examples do **not** need `activity.state`. A press, selection, warning or retry does not establish continuous activity/attention. They need safe matching journey context to attach to that journey; otherwise reception can remain accepted-unassigned. Mapped activity refreshes inactivity, so noisy automatic events must not keep a journey alive indefinitely. Error/warning classification is independent of whether an event is a graph step or supporting activity.

For a trusted saved vote, keep a distinct `vote.committed` event from the persistence-owning server. Keeping a browser press off the graph does not turn it into proof that the vote was saved. Do not change real request/retry behavior merely to produce these example observations.

### Presence/activity-state observations: support observed active time

Only use `activity.state` when you actually observe the stated condition:

| Observed condition | Example state | Meaning/limit |
| --- | --- | --- |
| Browser document becomes visible | `visible` | Visible is not proof of engagement |
| Browser document becomes hidden | `hidden` | Stop the observed interval; it does not prove journey completion |
| App's deliberately defined activity detector enters an active state | `active` | State under that detector's policy, not a universal measure of attention |
| That detector enters idle | `idle` | Stop the observed interval; disclose the detector policy |
| Heartbeat while the observed positive state remains true | Repeat `visible` or `active` | Refresh the bounded interval; do not emit positive state while hidden/idle |

The accepted states and 60-second cap are **Sentrame's conventions**, not OTel semantic standards. Pair them with explicit stage entry/exit when asking about active stage time. A single click does not justify a later `active` heartbeat. The server cannot infer browser visibility from an HTTP request. Avoid raw mousemove/keystroke streams, entered text, coordinates and other high-volume or sensitive payloads.

### Decision rule for an instrumentation agent

1. Is this a duration-bearing operation? Reuse/create a standard span/.NET Activity where appropriate.
2. Is this a point-in-time fact? Emit one standard named log or span event.
3. Should the fact be a visible milestone in this journey? Map it to **step** (or an explicit boundary Kind where warranted).
4. Is it supporting context that should stay off the graph? Map it to **activity** only if its lifecycle effects are intended.
5. Is it an actual observed presence state? Consider the optional `activity.state` convention and disclose its limits.
6. Verify ordinary receipt first; add only the context needed for the chosen question. Do not make every emitter depend on `stage.name` or presence conventions.

## Practical rule

Emit business events for things you know happened, spans for operation timing, explicit paired events only for a broader stage without a sufficient native span, and presence signals only for observed state. Map them intentionally; keep identity, stage naming, timing and completion separate. Send no person/household IDs, credentials or arbitrary payloads. The source of truth for the current projection and derivation is [intake.go](../../internal/app/intake.go) and [journey.go](../../internal/app/journey.go); SDK/runtime/network qualification remains separate from these documented envelopes.


## Implemented structural receipt (local greenfield model)

**Available at the tested local wire seam:** native start/end decimal nanosecond strings, parent ID, kind/status/flags, ordered event indexes, correlated log trace/span IDs, SDK drop counts and at most32 link target IDs plus projection omission count. No custom stage/session/occurrence fields are required. Native span identity is immutable; changed retained content rejects the whole direct request, while exact replay remains one contribution after SQLite reopen. Link attributes, trace state and unknown attributes are omitted before storage. Fresh disposable SQLite only; unsupported old schemas fail explicitly and are never reset automatically.

HTTP method, integer response status100–599 and minimized route templates are now projected. The exact restrictive route grammar and limitations are in [the internal contract](../../CONTRACT.md). Raw URLs/queries, headers and personal content remain forbidden; a template-looking literal is still the producer's minimization responsibility.

**Available at the tested local wire/UI seam:** protected trace tree/timeline/table with native inclusive duration and typed relationships. **Available at the qualified local seam:** service/carrier mapping controls and native operation occurrences in the trace business overlay. **Available at the qualified local synthetic seam:** optional hints and span-backed scoped journeys. Actual exporter and environment cells remain separately qualified. Retention support alone does not make those interfaces available. Qualification: [standard OTel ledger](../verification/standard-otel.md). Actual .NET/Next.js exporters and every deployment route remain separately qualified; these local wire assertions do not establish them.


## Standard-name mapping and native operation occurrences

Owners can configure exact service, signal/carrier, name, trusted source and one supported retained condition without `stage.name`. Native spans use app/producer/trace/span identity and inclusive native start/end; no synthetic boundary rows or custom occurrence IDs are created. The business basis is `native-inclusive`; point-event timing is unavailable by default. Optional paired/estimated alternatives remain explicit and are not pooled into native operation duration.

An unambiguous goal mapping selects its source/service/carrier/condition. A native span goal evaluates at a valid end boundary, excluding error status or missing/incompatible clocks; an event goal evaluates at event time. A persistence goal requires a trusted named post-commit fact or an operation whose documented success proves persistence. An owner label, generic HTTP2xx or browser press supplies no such contract. `vote.committed` and `Vote saved` describe the actual commitment; failure/retry/recovery and publication errors are separately retained.

Activating a definition is prospective. Default overlays use capture-time definitions; an explicit `previewDefinition` is a non-persisting interpretation preview, never an overwrite. The trace overlay remains an operation tree, without chronological sibling transitions or journey conversion counts. Scoped journey integration is qualified at the local synthetic wire seam; exporter and deployment qualification remain separate.


## Scoped journeys and timing policy

For L4, emit explicit random `session.id`, with browser `tab.id` where applicable, across the participating requests/events. Server-only scopes need no tab. This is application propagation, not automatic baggage or trace-to-session inference. The client session convention is opt-in/Development; Sentrame's server workflow interpretation is additional application policy. No previous-session/person/IP/agent/fingerprint joins are supported.

Mapped native spans retain native identity and exact times in scoped journeys. Native goals close at qualified span end; failed/missing/invalid ends do not prove completion. Inclusive operation intervals never imply page dwell. Optional paired page boundaries and actual presence remain separate bases; conflicting bases for one occurrence are excluded until an explicit policy is chosen. Long-lived summaries identify `journey-native-v3` and separate native-inclusive, paired-wall-clock, presence-cap60 and adjacency-estimate policies; they retain no trace/scope membership or diagnostic free text.

The new schema3 model adds retained `captureDefinition` metadata for late earlier prefixes of a known captured journey, preserving `definition` as the immutable receipt version. This metadata expires with detail. Unsupported older local schema versions fail explicitly; choose a fresh directory and never reset automatically. See [qualification](../verification/standard-otel.md) for local wire and browser-UI results versus unrun exporter/route cells.
