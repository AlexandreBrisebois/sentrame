# Design

## Context

See `proposal.md` for the product scope and `inventory.md` for source evidence and previous-change reconciliation. The baseline is `main` at `c7288d5`. The source currently binds each app's two JSON credentials to browser/server evidence; that source participates in identity, mapping and journey behavior. Unsupported nonempty SQLite schemas are rejected, so existing-state migration requires deliberate work rather than relying on startup to convert it.

**Review status:** the deltas express a recommended design for review. Alex approved specification work only. The three product decisions below must be accepted or revised before an apply instruction is actionable; their presence is not permission to revoke keys or change data. There is no implementation in this change.

## Goals / Non-Goals

**Goals:** one active write key and one client-observation model per registered app; standard existing SDK/exporter integration; understandable journeys and diagnostics; explicit history/replay compatibility; lean deployment and accessible investigation.

**Non-goals:** remove Sentrame's backend, receiver, storage or authentication; add a custom SDK; implement public activation/enrollment/installation identity; instrument another app; add APM, dependency topology or request waterfalls; delete historical data; alter viewer grants or publish a new public route. Existing proxy/forwarder infrastructure can remain where it carries client telemetry without instrumenting the host application's server.

## Decisions

Existing canonical scenario names are retained in deltas where OpenSpec requires traceability. Their updated assertions define the replacement behavior; a legacy server-related scenario title does not retain that feature.

### 1. Remove the telemetry role, retain evidence distinctions

Use app identity/generation for authorization and service, version, signal/carrier and retained evidence for interpretation. New registry/API/types/UI expose one `telemetryKey` concept and no producer role. The specific JSON field spelling is an implementation contract to finalize before coding; tests must lock it consistently. Connection remains an owner-only read; rotation remains an owner-only origin/CSRF-protected mutation. Lists stay masked and keys stay out of URLs, logs, telemetry and persistent browser state.

Remove active browser/server mapping fields, trusted-source badges/help, role-specific reveal/copy/rotate cards, Follow operation labels and server/combined journey derivation. Retain the word browser where it describes real browser behavior, exporter lifecycle, tabs, CORS or zoom. Do not mechanically delete references to the server that enforces Sentrame's authorization or runs Compose.

Alternative: keep two backend credentials hidden by one UI card. Rejected because role authority and revocation remain confusing and violate one-key behavior. Alternative: collapse every `source` field blindly. Rejected because it merges identity domains and silently broadens old rules.

### 2. Keep OTel reception independent of physical origin — product decision A

**Recommend:** retain currently qualified OTLP logs/traces/metrics and ordinary bearer configuration. Do not filter records based on `SpanKind.SERVER`, SDK runtime name or absence of browser-specific fields: these are payload claims, not reliable origin evidence, and valid client instrumentation/forwarding can include them. This is receipt compatibility, not advertised server telemetry support. Uncorrelated records remain sanitized received evidence and optional diagnostic context; no server operation journey is constructed. Metrics remain receipt-only.

Client journeys require explicit app-scoped visit and tab correlation. Automatic navigation keeps named non-diagnostic client logs under an empty definition; configured spans/events require applicable mapping and the same client scope. A trace alone never becomes a visit. Preserve native span timing, events/log relationships and errors as diagnostic evidence; keep Traces as a deliberate supporting destination rather than expand dependencies. All observations have one untrusted assurance. A user can name an observed result as a goal, but that proves only the retained observation under the rule, never durable persistence or truthful human action.

Alternatives: logs-only reception is smaller but breaks qualified exporters and discards useful client trace/error evidence without a stated need. Rejecting all SERVER-kind spans creates false origin guarantees and unnecessary OTLP incompatibility. Removing Traces altogether simplifies navigation further but loses independent client diagnostics; Alex must choose if that additional reduction is intended. No new protocol/encoding or browser/Tunnel support is qualified by this spec.

### 3. Retain the existing browser key at a reviewed cutover — product decision B

**Recommend:** for existing managed apps, preserve the existing browser key bytes as the single app key. Retire the server key only at the explicitly approved coordinated migration; do not publish it, use it as an alias or give it client meaning. Fresh registrations generate one random key. Rotation after cutover replaces that key for all senders, preserves app identity/generation and history, and maintains authority locking through commit.

This avoids asking working client exporters to change configuration. It also deliberately stops any old server sender. Inventory affected sender configurations during the later authorized migration without printing values. There is no dual-key grace interval: that would perpetuate removed authority. Return ordinary invalid-credential responses without revealing which retired role a caller used. Connection/rotation APIs return the new shape; old `role` mutation inputs fail explicitly, never silently rotate the entire app.

Alternatives: generate a fresh app key at cutover for maximum separation, requiring coordinated client updates; temporarily accept both old keys, minimizing short-term disruption but violating the selected one-key boundary and permitting backend observations to reenter as client history. Fresh-key cutover is reasonable if Alex prefers a deliberate clean break, but changes migration acceptance scenarios and must be selected before apply. Neither choice makes public clients authentic.

### 4. Preserve state, isolate incompatible history — product decision C

**Recommend:** keep retained legacy records and definition/reference metadata under their existing retention and deletion policies. Do not create an in-app legacy server explorer. Version the new projection/derivation and report incompatible pre-cutover evidence as unavailable/non-comparable with a useful recovery action.

- Old server records never become client observations or contribute to new graphs, timeline bins, diagnostics, traces, discovery, error counts or live/baseline metrics. Apply eligibility before source scan/display caps so excluded records cannot starve current reads or distort coverage denominators.
- A mixed capture remains incompatible as a whole. Removing B from A→B→C must not invent A→C. Do not drop its original boundary and then regroup by timestamps, session or tab. Unknown provenance fails closed for current analysis.
- Pure client history is usable only when capture definitions, original boundaries and supported replay identities remain equivalent. Otherwise explain the compatibility limit. Historic presentation must not expose a live role selector merely to preserve an old wire field.
- Preserve old definitions. Preview any source-selector removal; server-dependent goals and newly overlapping rules block automatic activation. Owners create a prospective role-free definition after review. Do not convert old persistence goals to screen goals silently.
- Preserve saved references. An unsupported reference gets a clear message, return action and explicit selection of a current client period. Never rewrite it, ignore its role filter or return a fake zero. Invalid active role query fields receive an explicit unsupported error.
- Old summaries lacking a separable client dimension cannot be accurately subtracted. Preserve them under retention and label the interval non-comparable; current comparisons must not pool them. Do not force 12-month history into the new computation version.

Alternative: continue full mixed-history rendering in a hidden legacy engine. It preserves more in-app access but maintains a server product indefinitely. Alternative: delete all old state/start fresh. Simpler technically but loses identity/history and is not authorized. Alternative: rederive filtered history. Rejected because it invents paths, goals and denominators. A minimal archive-only transition that makes all pre-cutover analysis unavailable is safer if equivalence for pure client history proves too costly; selecting it changes the deltas and must return for review.

### 5. Keep replay and lifecycle authority stable

The new namespace is app/generation/signal with signal-specific record/span identity, not the rotating key. Existing client identities need a versioned deterministic alias to their old browser namespace or equivalent preserved keys; server namespace aliases are never added to the new client domain. A role label retained only in an old schema decoder is compatibility bookkeeping, not active authority or UI. Do not accept caller-chosen namespaces as authenticated identity. Stable event IDs must be unique in their documented app/signal domain; conflicting content fails atomically. Anonymous duplicates remain uncertain. Cross-app collisions remain isolated.

Preserve original IDs/times, tombstones, supported replay suppression horizons and deleted generations. Key rotation cannot make old evidence new. Native relationships resolve only within the appropriate app/generation and compatible identity domain; equal legacy IDs cannot produce a client-to-server link or ambiguous attachment resolved by proximity. Preserve privacy projection and protected read checks across traces, copies, SSE and cached selections.

### 6. Keep deployment and investigation lean

Reuse the current Go receiver/SQLite store and one Sentrame Compose service; no new bootstrap service, enrollment issuer, mandatory Collector or separate trust dashboard. Resource limits, acknowledgment boundary, health/readiness, private registry persistence, owner sessions and temporary viewers remain required. Any public-ingestion hardening needed for the retained route must use explicit bounded request/decompression/concurrency/rate/cardinality/storage policies; current source limits are not proof of Internet-scale capacity.

Preserve Tools/Outils and the graph-first compact header, seven direct journey views, stable refresh, keyboard/focus return, 44px targets, English/French parity and existing ~165px desktop/~201px phone fixture canvas budgets. Reuse existing on-demand help and evidence actions. Old reference, missing context, empty app and expired evidence states must offer a concrete return, retry, filter or mapping action.

## Risks / Trade-offs

- Extractable client key permits spoofed observations and quota abuse → write-only permission, protected reads/admin, explicit untrusted wording and finite receiver budgets; CORS/TLS/key rotation do not attest a real installation.
- One-key compromise affects all client senders → owner-visible whole-key rotation and clear configuration update instructions. No per-install revocation is promised.
- Legacy analysis becomes unavailable even though files remain → preview counts/categories without secret values, explain intervals/references and preserve coordinated backup. Do not report archival files as a working in-app history view.
- Deleting source fields changes identity/mapping semantics → signal-specific aliases, mixed-history fixtures, prospective definitions and fail-closed conflicts.
- Optional SDK metadata cannot prove runtime origin → document product scope honestly rather than inventing origin validation.
- Hidden role dependencies survive in tests/docs/API consumers → use the inventory and an explicit residual audit; exempt only historical evidence, Sentrame infrastructure language and scoped legacy readers.

## Migration Plan

Planning only. No commands here are an instruction to operate existing state now.

1. Review decisions A–C and inspect the missing server proposal if it becomes available. Lock public API field names, version markers and compatibility errors before implementation.
2. Implement and qualify a bounded offline preflight/staging converter using disposable coordinated synthetic state. Check exact supported registry/store versions, free space and history categories. Emit counts/reasons, never credentials. Unsupported single-app state remains rejected.
3. For a later owner-authorized cutover, stop intake/writers and take a private coordinated backup of registry, stores, markers and latest revocation/deletion state. Validate app/generation consistency; preserve originals until staging is verified. Do not copy secrets into SQLite or ordinary reports.
4. Stage new registry/store metadata and identity aliases with an interrupted-migration marker; validate all apps before publishing the completed authority version. Preserve IDs and privacy data. An interruption keeps intake stopped and requires deterministic resume or verified rollback; filesystem-wide atomicity across JSON and several databases must not be claimed.
5. Verify new code with synthetic before/after/replay, retired-key rejection, mixed captures, definitions/references, summaries, deletion and rollback cases. A successful publication exposes one active key. Reconcile latest denies before reopening intake.
6. Follow client history/goals/diagnostics through actual retained standard exporter and browser paths in isolated state. Record Mac, container/Home Server and Tunnel results separately; do not claim unrun cells.
7. Rollback stops intake first. Older code must reject newer authority/state versions. Restore only a verified coordinated compatible set while preserving latest key retirement, rotation and deletion denies; if those cannot be represented safely, keep intake disabled and recover forward. A blind old image or snapshot restore is unsafe.

## Required review decisions

These are acceptance gates, not deferrable implementation questions:

1. Keep standard signals and supporting client Traces as recommended, or request a further reduction? Generic valid OTLP remains receivable, with no backend journey/dependency interpretation.
2. Reuse each existing browser key as the single app key and retire its server key at the authorized cutover, or choose fresh-key replacement with client reconfiguration?
3. Preserve stored history with incompatible mixed/summary intervals unavailable in current analysis, and keep pure client history only where semantics are equivalent? A full in-app mixed-history viewer is deliberately excluded.

Exact visible labels, review participants and migration preflight layout are deferrable within these boundaries. Public activation receiver ownership, opt-in, scale/budget, stable alias and installation enrollment decisions remain separate research questions and do not block this spec review. The unavailable prior server proposal limits exact reconciliation, not the clarity of Alex's superseding product direction.
