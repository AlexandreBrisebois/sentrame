# Sentrame integration instructions for coding agents

Use this guide to help a developer instrument client observations with their existing standard OpenTelemetry SDK/exporter. Sentrame is a client journey explorer with one telemetry key per registered application. It does not provide server-only or combined telemetry journeys, trusted persistence goals, dependency topology or a custom SDK. Sentrame's own receiver/backend remains required infrastructure.

## 1. Discover the client before proposing changes

Inspect source, UI, tests, existing SDK/provider/exporter initialization, service/version metadata, privacy rules and the supported forwarding route. Trace visible navigation, actions, observed results, retry/failure/recovery and screen completion. Ask only unresolved questions about which observation answers the developer's task. Do not infer a real outcome from a button click or generic HTTP status.

Propose a small event plan before edits. For each proposed observation state source evidence, the question it answers, the client emission boundary, illustrative stable name, safe correlation/identity, severity/noise, privacy and synthetic verification. Existing telemetry should be reused; do not initialize a second provider or instrument a backend merely to satisfy Sentrame.

A plan is not permission to edit/deploy another app, configure credentials or experiment on production. Follow the user's actual authorization and preserve business behavior. Do not synthesize exceptions, swallow errors or classify expected user behavior as failure to produce telemetry.

## 2. Configure one app and one key

Register/select the app through the owner UI and obtain its one Telemetry key plus ordinary OTLP settings. Configure the external client's existing exporter with the endpoint/protocol/bearer header. App attribution comes from that key, not resource claims, app names or trace IDs. Separate products require separate app-key requests; service/version metadata does not create new apps.

The client key is extractable write permission. It cannot read history or settings, mutate definitions, authenticate an installation or prove persistence. Never expose owner passphrases or viewer tokens. A same-origin forwarder may transport client OTLP where qualified; it cannot upgrade assurance. Direct browser/CORS/CSP/lifecycle and Tunnel compatibility need their own evidence. Do not introduce a custom OTLP envelope, SDK fork, signing transport or enrollment system.

## 3. Separate receipt, diagnostics and journeys

Valid supported logs, traces and metrics are receivable without visit, tab, stage or custom event IDs. Missing optional journey context is not a protocol error. Metrics remain receipt-only. Span kind and runtime metadata cannot reliably identify physical origin and grant no trusted server authority.

Client journeys need safe explicit app-scoped random visit/session and tab context. Missing context remains received/unassigned; it cannot fall back to a server workflow. Never join people, accounts, households, devices or visits by identifiers, IP, user-agent, fingerprints or nearby clocks. `session.id` is an opt-in client convention. Explicitly emit approved context: trace propagation/baggage does not automatically create journey membership.

An empty captured definition can display named non-diagnostic client navigation logs without a goal. Optional `stage.name` is a label suggestion; explicit owner mappings win. Unique supported source sequences help order evidence. Event record identity, action identity, stage occurrence, trace/span identity and visit correlation are distinct. New stable IDs must avoid collisions within their documented app/signal domain.

Use received service/carrier/name/retained-attribute context to map events. Preview match/conflict/missing-context results and activate prospectively. A goal is optional; a configured screen/result observation is still an untrusted observation, not a durable backend save. Preserve healthy loops and error → recovery → goal evidence without inventing abandonment or causality.

## 4. Add only purposeful timing and diagnostics

Native client spans already provide operation identity and inclusive duration. Do not emit paired logs for a span solely to time it. Point-event names do not supply duration. Add actual page entry/exit or visibility/idle context only when page elapsed or observed-active time is needed. Keep native, paired elapsed, estimated and observed-active bases distinct; none measures attention. Missing ends, clock uncertainty, overlap and late arrivals remain visible.

Retain controlled exception type/message/stack only under supported privacy projection and bounds. No personal content, raw request payload, arbitrary URL/query, credentials or secrets. Scrubbing is not perfect anonymity. Link errors only through supported explicit relationships, not time proximity. Grouping fingerprints do not deduplicate exception occurrences.

## 5. Verify a minimal approved integration

Exercise real standard exporter → receiver → storage → protected read seams in disposable state. Cover ordinary uncorrelated reception; client visits/tabs; mappings with/without goals; repeated result observations; failures/retries/recovery; stable replay; out-of-order/missing ends; privacy; history/live bounds; and saved-reference return/recovery.

Report exact SDK/version, route, signal, encoding/compression, batching/flush behavior and passed/failed/blocked/not-run checks. A wire test or source snippet does not qualify a framework, actual browser shutdown delivery, Home Server or deployed Tunnel. Use bounded supported processor/exporter options; no global ambient configuration changes, forks or custom transport to hide unsupported limits.

Existing two-key data needs the explicit reviewed offline transition. Do not operate on real retained state, read/reveal credentials, retire keys, delete data or restore snapshots merely because this guide describes the mechanism. Preserve mixed history without inventing client transitions, and never broaden old role-dependent mappings/references silently.

See [the developer guide](../DEVELOPER-GUIDE.md), [OTel contract](OTEL-CONTRACT.md) and [app workspaces](../MULTI-APP.md). Historical server integration reports are evidence of an earlier product boundary, not current integration instructions.

<a id="apply-the-wfs-lessons-to-the-next-integration"></a>

Earlier WFS references are preserved in [historical integration evidence](../verification/historical-server-integration.md); use the client discovery/verification checklist above for current work.
