# sentrame
Find the meaning between events : “sen-tram” Inspired by sens and trame: meaning, direction and the thread connecting events.

## Local synthetic review build

The first Svelte/TypeScript + Go + SQLite review build runs locally at [http://127.0.0.1:8096](http://127.0.0.1:8096). For the current canvas workspace, run `task setup`, `task start:workspace`, then `task demo` and `task demo:workspace` in another terminal. The workspace review uses its own synthetic database and preserves the original database. Owner access uses the generated private `.sentrame/local.env` passphrase. See the [launch/demo/REST guide](docs/LOCAL-REVIEW.md), [technical defaults](docs/review/release-a-technical-defaults.md) and [verification with remaining limits](docs/verification/release-a-local.md).

The current request authorizes local synthetic implementation. Historical planning-only authorization/status passages below describe the preceding phase. They do not restrict this local build. WFS changes, real collection, deployment, publication, merge and specification archival remain excluded; roadmap B remains deferred except for the explicitly authorized app-wide temporary viewers and NAS/container delivery preparation described below.

## North Star

Sentrame helps a small application team understand observed journeys: what happened, where people continued or stopped, which goals were evidenced, and how errors affected the path. What's for Supper (WFS) is the first integration to learn from; future applications remain separate. It is intended for a low-footprint Home Server deployment using Docker Compose, configured by an associated Compose file and `.env`, behind the existing shared Cloudflare Tunnel with HTTPS. A local synthetic application is implemented; Home Server deployment and the shared Tunnel remain unqualified.

The useful questions are concrete: Which paths reach a configured goal? Where do visits stop without a known outcome? Are repeated votes healthy exploration, recovery after an error, or an explicitly reported exit? Did a release change those patterns? Can a developer inspect the relevant exception and scrubbed context without guessing a cause?

## Confirmed direction

- Receive standard OpenTelemetry through OTLP using existing SDKs/exporters and ordinary configuration. No mandatory Sentrame SDK, custom transport, or custom fields merely to ingest telemetry. Product interpretation needs sufficient optional context.
- Keep stable app/version identities from day one. Support client-only, server-only and correctly correlated combined journeys. Browser visits stay separate per tab; server operations need no fictitious browser/tab. Shared app-scoped correlation links observations, not authority, total order, duplicates or person identity. No cross-app/cross-visit/person joins.
- The owner can emit application events and add ordinary standard OTel instrumentation. Journey steps map to those events. The UI lists discovered events, lets the owner choose a goal event, and shows how journeys reached that event from their beginning, marking missing/truncated prefixes honestly. Stable semantics support version comparisons; HTTP success and silence do not establish business intent.
- Warnings and exceptions are discoverable event classifications too, carried through standard OTel logs/events or qualified span events. Severity alone does not fail or end a journey: preserve continuation, recovery and later goal success. No custom SDK or protocol is required.
- Show **time at stage**, such as recipe viewing or deciding on a vote: timing is inspected contextually and edge thickness represents traffic; the switchable Sankey view uses volume-sized bars. Show median time plus supported slower cases, sample/completeness and numeric/table equivalents. Distinguish elapsed, observed-active and estimated time, including background/idle/tab changes, repeated occurrences and missing ends. Larger does not mean worse or more cognitive effort; exact timing/percentile policies remain open.
- Start at the first mapped event, end at the active selected goal, and start another journey at the next new mapped occurrence. After 30 minutes without mapped activity, close inactive—not failed. Preserve event-time/late corrections and supported duplicate-goal handling.
- Save reusable reference views: explore backward to a goal, between beginning/end events, or forward from an origin event. Views preserve recorded boundaries and do not keep deleted/expired raw journeys.
- Keep ordered detail for 30 days and count/conversion/timing summaries for 12 months. After 30 days, compare supported trends without individual journey drilldown. Preserve version/denominator/privacy and mergeable duration statistics; never average medians.
- Provide historical traffic maps and live views for a selected application or anonymous visit. Preserve filters when switching. A toggle offers a clean live window or the selected app's accumulated traffic over the last seven days, with the chosen live visit highlighted over that baseline. Keep compatible app/version filters and baseline/live counts distinct. Show potential dead ends by default with evidence and uncertainty.
- Keep exceptions attached only where linkage is valid. An exception does not end a journey. Repeated voting can be healthy. The user reports WFS moves the user to a completion screen: observing that screen is a visit outcome, separate from emitted end reason and who caused shared agreement. Suitable telemetry is not yet established. Silence cannot establish tiredness.
- Owner access uses a passphrase. The owner-approved extension adds temporary whole-app read-only invitations with expiry/revocation, including history, live maps, saved references and sanitized diagnostics. Narrower one-visit invitations with optional baseline/diagnostic grants remain deferred; viewer identity stays separate from tracked correlation.

## Principles and boundaries

Evidence before inference; plain-language explanations; privacy by minimization; bounded storage and work; reuse maintained standards and libraries; keyboard and accessible alternatives to graph-only interaction; minimal or no AI. The owner selected preservation of the first MVP’s warm-paper/forest-green visual theme and linked-thread icon; see the [visual reference](docs/design/visual-theme/README.md).

This is a focused journey-insight tool, not a full observability platform, person-tracking system, session replay service, automatic root-cause engine, or proof of satisfaction. Research and specifications do not authorize application/scaffold implementation, deployment, merge, or production testing.

## Proposed success measures and stop/go gates

These are review targets, not measured results or approved hardware budgets:

1. A developer can redirect an unmodified standard exporter and inspect accepted signals; unsupported transport/signal combinations are reported honestly. Richer journey interpretation has explicit missing-context states.
2. A WFS example distinguishes healthy voting, completion-screen observation, an explicit end reason, post-commit errors/recovery and unknown outcome from silence. Counts are replay-stable where qualified source identity permits; shared agreement is not attributed to every visit.
3. One reviewer can explain a path, denominator, version comparison and associated exception from the interface, using keyboard controls and an equivalent list/table.
4. Detail expires at 30 days and summaries at 12 months, with correction/deletion horizons made explicit. Unknown-ID/missing-time retries do not support unlimited original-age guarantees. Summary precision/grains and backup policy remain technical design defaults.
5. Measure the selected Svelte/Vite + Go + SQLite stack in Compose on the Home Server. The owner reports sufficient RAM and CPU options; hardware specifics, workload, disk growth and numerical bounds remain unmeasured. Tunnel HTTPS does not establish origin TLS, app authorization or routed SSE/OTLP compatibility.
6. Stop or adapt an existing project if a cheaper maintained combination answers these questions, ordinary WFS telemetry cannot support them, or the operational cost exceeds the eventual budget.

## Decisions, proposals and open questions

The five main choices are settled: a complete WFS release A with B preserved; automatic journey boundaries; narrow deferred invitations; 30-day detail/12-month summaries; Svelte/Vite + Go + SQLite with maintained OTLP components. [Release scope](docs/planning/sentrame-release-scope.md) keeps multi-app experience, richer version comparisons and broader sharing/configuration as retained roadmap B. A list-only tracer is an engineering step, not the final graph release.

Today's direction replaces the initial narrow concept. The [active OpenSpec change](openspec/changes/synthetic-journey-outcomes/proposal.md) contains the evolved capabilities, design and roadmap under its historical identifier. [Research](docs/research/sentrame-feasibility.md), [resource tradeoffs](docs/research/sentrame-stack-and-resources.md), [decisions](docs/planning/sentrame-decisions.md), [review disposition](docs/review/post-review-disposition.md), [internal harness proposal](docs/proposals/sentrame-engineering-harness.md) and [external integration-agent instructions](docs/integration/AGENT-INTEGRATION.md) support it.

The [deployment contract](docs/proposals/home-server-deployment-contract.md) is documentation only: no actual Compose/`.env` or tunnel setup is created. Integration agents inspect the target app, question the developer and propose useful events plus new warnings/errors/recovery before approved edits; this is separate from Sentrame's internal harness and adds no AI runtime. Remaining technical defaults will receive one focused design review before implementation. Planning-file reconciliation and a documentation-only draft PR are approved; no app/harness implementation, canonical sync/archive, merge or deployment is authorized.

## Authorized delivery and sharing extensions

The owner authorized temporary read-only whole-app invitations with expiry/revocation, and Docker Compose delivery with a linux/amd64+linux/arm64 DockerHub publishing workflow targeting `brisebois/sentrame`. These extend the local synthetic review; narrower per-visit invitations, multi-app experience and richer comparison/configuration remain deferred. No local implementation or CI definition alone proves image publication, NAS operation, shared Tunnel transport or real WFS export. Follow [NAS delivery instructions](docs/NAS.md) for the Compose/environment contract and qualification status.

Publishing status: the owner authorized committing/pushing the implementation branch after both local container smoke checks passed. DockerHub publication/release tagging follows registry-secret setup; no image has been published by this work yet.
