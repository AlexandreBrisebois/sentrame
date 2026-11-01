# Send client OpenTelemetry to Sentrame

Sentrame explains client event journeys. Register one application and use its single telemetry key with an existing standard OpenTelemetry SDK/exporter. No Sentrame SDK, custom wire envelope or server instrumentation is required. Sentrame's own Go backend, private registry, receiver and SQLite persistence still run as one service.

## Register the app and configure the exporter

1. Sign in as the owner and add/select the app.
2. Open App settings. Explicitly reveal or copy its Telemetry key and endpoint settings.
3. Configure the existing exporter with its reachable endpoint, supported OTLP/HTTP protocol and `Authorization: Bearer <telemetry-key>` header. Preserve its SDK/provider; do not initialize a duplicate provider.
4. Send a minimized synthetic observation and inspect protected received evidence before troubleshooting a missing journey.

The shared signal URLs are `/v1/logs`, `/v1/traces` and `/v1/metrics`. A base endpoint lets the exporter append its signal path; a signal-specific endpoint includes the entire path. Follow the exporter's own documented endpoint behavior. Supported encoding/compression and actual exporter versions need separate qualification. Existing local wire evidence does not establish every framework, direct browser or Tunnel path.

The key binds a request to one registered app/generation independently of `service.name`, version or claimed app IDs. A Collector exporting multiple products must separate requests by app key. New releases of the same product remain in the same app; retain bounded `service.name` and `service.version` as metadata, with unknown values left unknown. Restart, rename and app selection do not rotate the key.

## Public client key and protected access

A key distributed to a client can be extracted and used to spoof observations. It authorizes writes only; it never authorizes reads, App settings, definition changes or owner administration. Labels, IDs, span kind and HTTPS do not prove an authentic installation, a human action or a durable backend save. CORS is not authentication. Keep owner passphrases, viewer sessions and registry files private.

Use a separately qualified same-origin forwarder if the exporter requires it. Forwarding carries client observations and does not create a trusted server telemetry role. Direct cross-origin browser export, browser lifecycle/flush behavior, CSP and deployed HTTPS/Tunnel paths must be qualified explicitly; copying settings does not establish those routes.

Rotation replaces the one app key for all senders. Update each sender's exporter configuration afterward. It preserves app identity and history but cannot permanently exclude someone who can extract the new public key. Public activation, enrollment and per-install revocation are separate features and are not promised.

## Progress from receipt to a useful client journey

Ordinary valid supported logs, traces and metrics can be received without custom journey fields. Receipt does not imply that a visit or business path can be reconstructed. Metrics remain bounded receipt evidence; they do not reconstruct individual journeys or conversion populations.

For a client journey, emit safe explicit random `session.id` and `tab.id` within the registered app. Use a bounded visit and a separate tab scope; never use a person, account, household, IP, user-agent or fingerprint as correlation. Across requests, explicitly emit the approved context: trace propagation or baggage does not automatically place it on every event. Missing visit/tab context stays received/unassigned rather than becoming a server-operation journey.

With an empty captured definition, named non-diagnostic client logs can show goal-free navigation. Optional `stage.name` supplies a display suggestion; event name is the fallback. Unique positive `source.sequence` values within the supported emitter/service scope provide ordering; otherwise occurrence and receipt ordering carries uncertainty. Stable record identity supports replay; it must be unique within the documented app/signal identity domain. IDs are distinct from actions, visits and semantic occurrences.

For configured meaning, select discovered events and use service/carrier/name/retained-attribute selectors. Preview matches and conflicts before prospective activation. No browser/server selector or trusted persistence goal exists. A goal is optional; explicitly choose an emitted observation if it answers the question. `vote.result.observed` can describe the result the client saw; it cannot certify a durable save.

## A small voting example

Start with received client navigation, `vote.requested`, a request span, `vote.result.observed`, a warning/error if one occurred and `completion.entered` when that screen was actually shown. These are illustrative names, not required OTLP fields.

Keep ten real result observations and nine successive result-to-result transitions when the user legitimately repeats the action. Do not classify repetition, negative votes, silence or severity alone as failure. Keep failure → retry → recovery evidence and subsequent goals available. One intent and one result are two records, not automatically two actions or one deduplicated action.

If useful, map native client spans to operation stages. Native trace/span identity supplies duration; no paired synthetic boundary logs are needed. A point event does not gain duration merely because it has the same label. Introduce optional `stage.name` only as a label suggestion, then explicit visit/tab correlation or presence only when the desired question still lacks evidence.

## Timing, errors and privacy

Native inclusive span duration, paired elapsed stage time, next-event estimates and observed-active intervals remain separate bases. Overlapping durations cannot be added as unique visit time. Missing/invalid ends are censored or unavailable, never zero. Page dwell requires page interval evidence; observed-active time additionally requires actual scoped visibility/idle observations under the declared heartbeat policy. Neither means attention.

Send controlled technical messages and safe exception details, not request bodies, arbitrary URLs/query strings, personal data, credentials or user content. Projection and bounds apply before storage, preview, SSE and copy. Text scrubbing is useful but not perfect anonymity. Do not add unsafe attributes merely to restore an analysis.

Errors may be explored without journey correlation. Link an exception to a journey only through supported explicit context/relationships; time proximity does not establish cause. A shared stack is a grouping category, not proof that two records represent one exception. Missing source maps remain missing.

## Legacy managed state and recovery

The single-key release does not silently adopt old two-key managed state. Follow the reviewed offline migration mechanism only after separately authorizing a real cutover and taking a private coordinated backup. The recommended transition retains the old browser key as the app key and rejects the old server key. Do not import an old server key as an alias or expose it through connection settings.

Mixed historical journeys and inseparable old summaries are non-comparable to current client results. Removing B from A→B→C must not fabricate A→C. Compatible pure-client history can remain usable only with unchanged capture meaning and replay identity. Unsupported references explain the compatibility boundary; they must not broaden filters or become fake zero traffic. Existing retention/deletion rules still apply.

Back up registry, app stores, initialization/migration markers and latest revocation/deletion state together. Restoring an older image or snapshot alone can revive retired authority. Keep intake stopped when current deny state cannot be reconciled. See [app workspaces](MULTI-APP.md) and [deployment guidance](NAS.md).

## Verify the actual path

Use isolated synthetic data and the actual supported exporter → receiver → durable store → protected read seam. Cover uncorrelated receipt; visit/tab isolation; optional goals; repeated events; warning/error/recovery; replay before/after restart and rotation; missing/late timestamps; privacy; and current/legacy history boundaries. Wire fixtures, source snippets and mocked views do not qualify actual SDKs, browser lifecycle, Home Server capacity or Tunnel routes.

`task test`, `task test:delivery` and the isolated fixture commands provide repeatable local checks. Do not run reset, registration, migration or fixture commands against real retained state without separate authorization. Tests should report passed, failed, blocked and not-run cells distinctly.

The [coding-agent guide](integration/AGENT-INTEGRATION.md), [OTel contract](integration/OTEL-CONTRACT.md) and [synthetic demonstration](integration/synthetic-demo.md) use this same client-only boundary. [Historical server integration evidence](verification/historical-server-integration.md) preserves earlier .NET/Node results and limitations; those older recipes are not active server support or qualification of this release.

<a id="9-lessons-from-the-wfs-integration"></a>
<a id="qualified-wfs-net-composition--2026-10-31"></a>
<a id="presentation-and-outcome-boundaries"></a>

Historical WFS integration links resolve to the preserved [earlier integration evidence](verification/historical-server-integration.md). Those dated server recipes and qualifications describe the previous product; they are not current setup instructions.

Local qualification update (2026-11-01): [native Edge + standard OTel JS logs/traces through a disposable same-origin forwarder](verification/client-browser-otel.md) passed with the one-key receiver, replay and receiver process restart. This specific loopback HTTP Protobuf/no-compression cell does not qualify a deployed forwarder, direct cross-origin export, mobile lifecycle delivery or the shared Tunnel.
