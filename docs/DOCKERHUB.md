# Docker Hub description

Short description:

> Self-hosted OpenTelemetry journey maps for understanding app flows, timing, errors, and recovery.

## Sentrame

Sentrame turns OpenTelemetry events into interactive journey maps, helping you understand how your application is used and what happens when things go wrong.

Explore historical and live activity through hierarchy, force-directed, and Sankey views. Inspect incoming and outgoing paths, repeated steps, stage timing, exceptions, and recovery—with ordered evidence behind the graph.

Map discovered events to meaningful journey stages, save reference views, and share read-only access through temporary invitations with expiry and revocation.

Sentrame is self-hosted, built with Go, Svelte, and SQLite, and designed for Docker Compose. It retains sanitized event detail for 30 days and supported summaries for 12 months. Standard OTLP/HTTP ingestion works without a proprietary Sentrame SDK.

Images are available for Linux AMD64 and ARM64 under one Docker Hub tag. Release candidate `v0.1.0-rc.1` passed both architecture smoke checks and publication.

For setup, environment variables, persistent storage and supported transport boundaries, see [NAS/container delivery](NAS.md). Registry publication, a working NAS route and an application integration are separate verification steps.

User documentation: [user guide](USER-GUIDE.md) · [developer integration guide](DEVELOPER-GUIDE.md).

The published `v0.1.0-rc.1` is the historical baseline. Standard trace structure/native timing, contextual business mapping and optional hints are implemented in local source (schema3 / structural-v2 / journey-native-v3), with [separate qualification evidence](verification/standard-otel.md). No new image is published by that work. Preserve existing data and select a fresh dedicated directory for the greenfield cutover; unsupported older schemas fail without automatic migration or reset. Native trace trees and chronological business journeys remain distinct.
