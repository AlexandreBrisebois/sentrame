# Sentrame local Release A review build

This is the authorized synthetic implementation on `codex/sentrame-mvp`, based on planning commit `4a20c95`. The historical planning-only passages describe the previous phase. No WFS app was contacted or changed; no real telemetry, deployment, publication, merge or spec archive is included. Roadmap B remains deferred. Local qualification and the remaining OpenSpec finish gates are recorded separately in [verification](verification/release-a-local.md).

## Canvas workspace update

The owner-approved workspace is implemented with hierarchy, force and Sankey views. The original `.sentrame/local.db` was already empty and intake-disabled when this work began; it is preserved. The current review uses `.sentrame/workspace-review.db`. To resume this exact workspace after stopping the service:

```sh
task start:workspace
```

The URL remains [http://127.0.0.1:8096](http://127.0.0.1:8096). Refresh an already-open tab to load new built assets. A restart requires owner login again with the existing local passphrase.

- Switch **Hierarchy / Force / Sankey** in the top bar. Select a node for incoming/outgoing traffic, timing and evidence; select a connection for observed transitions. Small arrowheads retain screen size independently of traffic. Use trackpad **pinch** to zoom at the pointer and two-finger scrolling to pan, or **Fit**, +/- and keyboard controls; live snapshots do not continually rearrange the camera.
- **Errors** opens application-wide diagnostics without requiring a visit. Search by message, exception, service or trace ID; expand a record for a scrubbed message/stack, technical IDs and sanitized copy. Details are bounded, expire with30-day evidence, and scrubbing is imperfect. Actual .NET/Next.js exporter qualification remains separate.
- **Configure event…** opens a node popup. Preview actual matching records before activating a copied, prospective definition. Existing captured journeys retain their version. **Mappings & goal** is in the workspace menu; goal changes remain explicit.
- **Paths into here / Paths from here** preview recorded projections without saving. **Use as beginning…**, then select an end stage for a segment. **References** saves the current investigation or overview with exact time bounds, visit/version/definition filters, layout and highlight. It preserves captured boundaries, not expired raw evidence. A live reference freezes its live observation window; the baseline remains a separate comparison layer.
- Error-related observed connections are red by default, with evidence labels; healthy continuation is separate. Choose **Rare paths**, **Endings**, **Warnings & exceptions** or **Recovery** in Highlights. Rarity uses full closed, non-partial/uncensored sequences at <=10% share with >=20 comparable journeys per definition/version. Shared highlighted edges can also occur in common paths. The base demo has low support; `task demo:workspace` adds a qualified rarity example.
- **Historical / Live** opens the event-volume timeline. Select a context (1h/24h/7d/30d), drag a span or move either handle; the shared green track shows only the included interval. Quick ranges and advanced exact UTC bounds remain available. Choose a visit from a node’s journey details using **Follow this visit / operation**; optional exact correlation-ID search is under advanced filters, with no opaque-ID dropdown. a selected visit has an **X** to return to all traffic while preserving other filters. **Table** and **Journeys & evidence** are always available. Discovery and lifecycle tools live in the workspace menu.

```sh
task demo:diagnostics         # standard trace exceptions, JS-like stacks and custom error bodies
task demo:workspace           # generate, ingest and verify 20 common + 1 rare journey
task fixtures:workspace       # prepare fresh protobuf/JSON files only
task demo:workspace:verify    # execute/replay the same literal requests
```

For VS Code REST Client, generate fixtures (`task fixtures:workspace` or `task fixtures:diagnostics`), open `rest/workspace.http` or `rest/diagnostics.http` and send its requests in order using the adjacent ignored `rest/.env`. The script executes those same requests and verifies the real protected snapshot. It does not preload the database. Detailed verification and remaining limits: [workspace verification](verification/workspace-design.md).

## Start and access

From `/Users/alex/Code/sentrame`:

```sh
task setup
task start
```

Open [http://127.0.0.1:8096](http://127.0.0.1:8096). Use that exact origin, not `localhost`, because mutations enforce the configured Origin. `task setup` creates ignored `.sentrame/local.env` with a random owner passphrase and separate browser/server producer credentials. Open that local file to obtain `SENTRAME_PASSPHRASE`; keep it private. No default credential is embedded in source or browser assets. The current review tab may already be signed in. Restart/logout invalidates its cookie and requires login again.

Prerequisites used here: Go 1.27.3, Node 26.12.0/npm, Task and OpenSpec 1.14.3. Exact Go modules/checksums are in `go.mod`/`go.sum`; frontend dependencies resolve from `web/package-lock.json` through `npm ci`. Go serves built Svelte/Vite assets and viewer/API/OTLP HTTP on loopback port 8096. A private ephemeral loopback listener hosts the maintained Collector receiver; its consumers independently validate producer credentials. gRPC is not exposed.

Configuration inputs: `SENTRAME_PASSPHRASE` (at least 12 characters), distinct nonempty `SENTRAME_BROWSER_TOKEN` and `SENTRAME_SERVER_TOKEN`, optional `SENTRAME_ADDR`, `SENTRAME_ORIGIN`, `SENTRAME_DB`. Defaults are loopback `127.0.0.1:8096`, origin `http://127.0.0.1:8096` and `.sentrame/local.db`. Do not repurpose this as a public deployment; Home Server/HTTPS/proxy configuration is not qualified.

## Repeatable demonstration

In another terminal, while the service runs:

```sh
task demo
task demo:replay
```

The tool generates real official protobuf bodies from `rest/fixtures/scenarios.json`, then executes the version-controlled [REST Client requests](../rest/demo.http) literally. It verifies protected snapshots, not a parallel database preload. Fresh generation creates new synthetic scopes; replay preserves original IDs/bytes. Existing demonstration generations remain visible until retention or explicit reset.

1. **Explore:** inspect the map, map/timing table and ordered sequences. The healthy case has ten trusted Vote commitments and nine Vote→Vote transitions. Intent and completion-screen observation are separate facts.
2. **Recovery:** expand Discovery→Intent→Vote→Publication error→Recovered→Completion. The earlier error remains; it never undoes the persisted vote or proves an exit cause. Copy exposes only the retained sanitized projection.
3. **Timing:** compare measured elapsed and observed-active values under the same definition. Unknown time stays unavailable, not zero. p90 appears only with at least 20 completed compatible samples; exact sample/exclusion counts accompany it. Timing is in the inspector/table; edge width is transition traffic and Sankey bars use volume. The owner-approved canvas update supersedes duration-sized nodes.
4. **Configure:** use the workspace menu for discovered events, then select a node and its configuration popup to inspect producer roles and missing context. Preview mappings and select an active mapped goal. Activation creates a new definition; ongoing captured journeys stay pinned to their original definition. Browser claims cannot establish server commitments.
5. **Reference views:** save goal-backward, origin-forward or bounded start/end views and choose first/last repeated endpoints. Reopening applies captured definition/date/version filters and projects the map/sequences inside each recorded journey. Missing/cross-boundary endpoints remain marked; queries never join captured journeys or change their operational goal.
6. **Live:** switch to Live, use the clean window or toggle the fixed seven-day app baseline, select one anonymous visit and pause/resume rendering. Counts remain separate; baseline remains app-wide. SSE reconnect/reset replaces snapshots instead of adding counts.
7. **Data lifecycle:** inspect supported daily summaries, exclusions and explicit deletion limits. Whole-app deletion disables this intake generation; one-day summary deletion is a distinct operation. Try destructive controls only on disposable synthetic data.

For a fresh visible live continuation:

```sh
task demo:live
```

Select the newly printed/generated live scope in the browser, then start Live observation before sending:

```sh
task demo:followup
```

First manual follow-up gets a fresh occurrence time while retaining correlation/identity. Ambiguous retries reuse the prepared bytes. An identified replay must not create another goal. See [the full REST/fixture guide](integration/synthetic-demo.md).

## REST Client

Run `task fixtures`, open `rest/demo.http` in VS Code with REST Client, and use the adjacent generated ignored `rest/.env`. Send owner login to obtain the cookie, then the ingestion and protected read requests in order. Binary bodies are loaded from `rest/generated/*.bin`; paired `.json` files are readable official OTLP JSON. Re-send the same POST for replay; regenerate for fresh IDs/timestamps. Close the final SSE request when finished. `task demo:verify` executes these same requests automatically; direct VS Code extension operation is not claimed as tested.

## Commands and reset

| Command | Scope |
| --- | --- |
| `task setup` | Resolve pinned dependencies, create local synthetic credentials and build |
| `task start` | Build and run the local service; Ctrl+C stops the process |
| `task test` | Go behavior/seam tests, UI helper tests and Svelte/TypeScript checks |
| `go test -race ./... -count=1` | Explicit broader concurrent access check |
| `task fixtures` | Generate fresh binary/JSON bodies, identities and ignored REST credentials |
| `task demo` / `task demo:verify` | Generate+execute, or execute existing literal REST requests |
| `task demo:live` / `task demo:followup` | Real intake with stable follow-up identity |
| `task demo:replay` | Replay exact identified ten-commitment bytes |
| `task qualify:local` | Isolated temporary loopback DB; 1,100 synthetic requests with concurrent query/SSE/cleanup |
| `task spec:validate` | Strict selected-change OpenSpec validation |
| `task check` | Tests, type checks, selected spec validation and whitespace |
| `task reset` | Refuses while port 8096 is serving; deletes only named local DB/WAL/config files and rotates credentials |

`task reset` is intentionally destructive to disposable local synthetic state. Stop Sentrame first. It does not restore the old Python stash, touch other applications, stop unrelated processes, or clean arbitrary directories. Reset/re-registration rotates credentials and creates a fresh persisted app/generation identity, so old queued requests cannot silently enter the replacement dataset.

## Limits for review

- Qualified intake is direct-origin HTTP only. Logs/traces/metrics binary none/gzip passed actual standard Go SDK exporters; JSON none/gzip passed official wire fixtures. No gRPC, Tunnel, browser/.NET exporter or upstream Collector-ack guarantee. Metrics are receipt-only; profiles are excluded.
- Detailed evidence expires at occurrence+30 days (supported missing-time fallback disclosed). Reads exclude expired data immediately; healthy cleanup runs each second. UTC day/version/definition summaries retain supported counts/conversion and mergeable timing histograms for 12 calendar months. No individual paths, visits or arbitrary new filters after detail expiry.
- Correction ends at captured-start+30 days. After finalization, individual subtraction/reprocessing is unsupported; use explicit app/UTC-day deletion. Retained partial tails and late evidence cannot restore started-journey eligibility or multiply finalized contributions. Saved queries cannot resurrect detail. No automatic backups/restore workflow or forensic erasure claim.
- Counts without reliable record IDs have retry uncertainty; clocks/sampling/collection completeness remain explicit. Controlled diagnostic messages/stacks and error-log string bodies now use bounded scrubbed-v1 projection. Ordinary bodies, raw request content and unknown attributes remain omitted; scrubbed text can still contain application data. Full source-mapped stacks and perfect text scrubbing are not implemented.
- At most 100,000 records in this local dataset; bounded recent 10,000-record derivation, 2,000 detail response rows, 100 graph nodes and 300 edges. Limits are visible; a truncated result is not a complete count. This short workload does not qualify full-capacity retention growth or sustained 30-day operation.
- Owner-only passphrase access is implemented. Invitations, multi-app experience and richer parallel comparisons remain B. Main English/French controls and explanatory text are available, but some metadata/diagnostic text remains English; comprehensive language/assistive-technology conformance is not claimed.

The most useful product review choices are the initial completion-screen goal, the value of conservative diagnostic omission, estimated versus explicitly measured timing, and how well the cohort/partial/censoring explanations help distinguish potential dead ends from unknown evidence.
