# Docker Hub description

Short description:

> Self-hosted OpenTelemetry journey maps for understanding app flows, timing, errors, and recovery.

## Sentrame

Sentrame turns OpenTelemetry events into interactive journey maps, helping you understand how your application is used and what happens when things go wrong.

Explore historical and live activity through hierarchy, force-directed, and Sankey views. Inspect incoming and outgoing paths, repeated steps, stage timing, exceptions, and recovery—with ordered evidence behind the graph.

Map discovered events to meaningful journey stages, save reference views, and share read-only access through temporary invitations with expiry and revocation.

Sentrame is self-hosted, built with Go, Svelte, and SQLite, and designed for Docker Compose. It retains sanitized event detail for 30 days and supported summaries for 12 months. Standard OTLP/HTTP ingestion works without a proprietary Sentrame SDK.

Image target: Linux AMD64 and ARM64 under one Docker Hub tag. Platform runtime verification is enforced by the publishing workflow; publish this claim with an image only after both platform checks pass.

For setup, environment variables, persistent storage and supported transport boundaries, see [NAS/container delivery](NAS.md). Registry publication, a working NAS route and an application integration are separate verification steps.
