# Sentrame

**Find the meaning between events.** Pronounced “sen-tram,” inspired by *sens* and *trame*: meaning, direction, and the thread connecting events.

Sentrame is a self-hosted OpenTelemetry journey explorer. It helps you see how application activity connects: which paths reach a goal, where visits stop, how often steps repeat, and what happens before and after an error.

The application is implemented with **Svelte/TypeScript, Go and SQLite**. It runs as one service with a built-in OTLP/HTTP receiver and a graph-first workspace.

## Start here

| I want to… | Guide |
| --- | --- |
| Understand the graph, configure events, investigate errors or share a view | [User guide](docs/USER-GUIDE.md) |
| Send telemetry from a .NET or Next.js application | [Developer integration guide](docs/DEVELOPER-GUIDE.md) |
| Run the container with Compose and a chosen data folder | [Container deployment guide](docs/NAS.md) |
| Run from source, load synthetic scenarios or use VS Code REST Client | [Local setup and demonstrations](docs/LOCAL-REVIEW.md) |
| Review requirements, decisions and verification | [Project documentation](docs/README.md) |

## What the workspace does

- Switch between left-to-right **Hierarchy**, directed **Force** and **Sankey** views. Select nodes and connections for incoming/outgoing traffic, counts, timing and ordered evidence. Drag nodes to improve the arrangement.
- Discover incoming event names and configure their journey meaning in context. Preview a mapping before activating a new definition; captured journeys keep their original boundaries.
- Explore historical traffic or watch live updates. Select time ranges with an event-volume histogram, follow a visit from its evidence, and compare live activity with a separately counted seven-day baseline.
- Highlight rare paths, endings, warnings/exceptions and recovery. Repetition can be healthy; an error does not automatically terminate a journey.
- Inspect application-wide errors, scrubbed messages and stacks without needing a session ID. Save and reopen reference views, or delete a saved view without deleting telemetry.
- Invite a temporary **read-only viewer** to browse the whole application. Links are redeemable once, expire after one hour, 24 hours or seven days, and can be revoked.

The workspace preserves the warm-paper background, forest-green branding and linked-thread icon. Tables and ordered sequences complement the graph; keyboard controls remain available.

## Run with Docker Compose

Image repository: [`brisebois/sentrame`](https://hub.docker.com/r/brisebois/sentrame). The first release candidate is `brisebois/sentrame:v0.1.0-rc.1`. The publishing workflow builds and smoke-tests **Linux AMD64 and ARM64** before publishing one multi-architecture tag. Check [available tags](https://hub.docker.com/r/brisebois/sentrame/tags) or the [Container workflow](https://github.com/AlexandreBrisebois/sentrame/actions/workflows/container.yml) for publication status.

Choose [compose.yaml](compose.yaml) for direct local access, or [compose.cloudflare.yaml](compose.cloudflare.yaml) for the existing `cloudflare-proxy` network. Both are standalone examples and use [.env.example](.env.example). Copy the selected Compose file and environment example into your deployment folder. Copy `.env.example` to `.env`, then set:

- `SENTRAME_IMAGE`: a published tag or immutable image digest.
- `SENTRAME_DATA_DIR`: the existing local filesystem folder to bind into `/data`, writable by UID/GID `10001`.
- `SENTRAME_ORIGIN`: the exact browser origin, without a trailing path, such as `https://sentrame.example.com`.
- A fresh owner passphrase and distinct browser/server producer tokens. The example contains no working credentials; keep `.env` private.

After preparing the folder and configuration:

```sh
docker compose config --quiet
docker compose pull
docker compose up -d --wait
```

The standard example publishes host loopback `127.0.0.1:8096`. For Cloudflare, use `docker compose -f compose.cloudflare.yaml up -d --wait`; that standalone file joins the existing external `cloudflare-proxy` network and publishes no host port. Configure the proxy target as `http://sentrame:8096` and set `SENTRAME_ORIGIN` to the public HTTPS origin. See the [deployment guide](docs/NAS.md) for folder ownership, HTTPS/Tunnel integration, health checks, upgrades and rollback.

## Try it locally from source

Prerequisites: Go **1.27.3** (see [go.mod](go.mod)), Node.js **26** with npm, and [Task](https://taskfile.dev/). OpenSpec CLI is needed for specification validation, not to run the service.

```sh
task setup
task start:workspace
```

Open [http://127.0.0.1:8096](http://127.0.0.1:8096). Find the generated owner passphrase in the private, ignored `.sentrame/local.env` file. Leave the server running, then load synthetic examples from another terminal:

```sh
task demo
task demo:workspace
task demo:diagnostics
```

These commands send actual standard OTLP through ingestion, storage and protected reads. They do not preload a separate graph database. `start:workspace` uses `.sentrame/workspace-review.db`; demonstrations expect the seeded mappings. See the [local guide](docs/LOCAL-REVIEW.md) before experimenting with definitions or resetting data.

## Connect your application

.NET and Next.js **servers can export independently** to Sentrame using standard OpenTelemetry exporters. A shared OpenTelemetry Collector is optional; .NET does not need to relay every other application's telemetry. Browser telemetry needs a separate untrusted producer path and must never contain the server ingestion token or viewer passphrase.

Use **OTLP/HTTP**, with ordinary authentication headers, at `/v1/logs`, `/v1/traces` or `/v1/metrics`. No Sentrame SDK or proprietary fields are required for valid reception. Meaningful journey maps depend on event mappings and sufficient safe correlation; technical records without that context remain accepted and unassigned. Metrics are receipt-only.

The [developer guide](docs/DEVELOPER-GUIDE.md) explains transport, credentials, .NET/Next.js setup, named events, exceptions and the path from first receipt to a useful journey. SDK-specific recipes and a deployed network route need their own verification; standard wire compatibility alone does not qualify every framework/exporter.

## Data and interpretation

Retain ordered detail for **30 days** and supported summaries for **12 months**. Summary-only periods cannot reconstruct individual journeys. Scrubbed diagnostic text is bounded, but minimization is not a promise of perfect anonymity: send controlled technical data, not credentials, personal information or request bodies.

Counts describe observed evidence. Sampling, missing identifiers, clock differences and query limits affect interpretation. Timing distinguishes measured elapsed, observed-active and estimated durations; it does not measure attention. See [data lifecycle and limits](docs/USER-GUIDE.md#retention-collection-health-and-limits).

## Development and current scope

```sh
task test                 # Go behavior/seam tests, frontend tests and type checks
task test:viewer-access   # Temporary viewer access boundaries
task test:delivery        # Probe/configuration checks
task container:smoke      # Local ARM64 container build and smoke test
task container:smoke PLATFORM=linux/amd64
task spec:validate        # Strict validation of the active OpenSpec change
```

Core journey exploration, diagnostics, temporary app-wide viewers and container delivery are implemented. Multi-app workspace management, richer version comparisons/configuration and narrower visit-only invitation grants remain deferred. Local/container checks do not establish native Home Server capacity, shared Tunnel behavior, actual .NET/Next.js exporter qualification or a deployed WFS integration.

Specifications live in the [active OpenSpec change](openspec/changes/synthetic-journey-outcomes/); its historical name does not limit the app to the obsolete first-vote design. Planning/research documents record how decisions evolved. Use the user/developer/deployment guides for current operation and the [documentation index](docs/README.md) to find the underlying contracts and evidence.
