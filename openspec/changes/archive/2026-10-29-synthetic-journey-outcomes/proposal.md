# Proposal

## Why

WFS needs understandable evidence about meaningful journeys, healthy repeated voting, explicit exits, and recovery after errors. Existing observability and product-analytics tools cover much of the machinery, but their identity, transport, path truncation and operating assumptions must be challenged before committing to another platform.

## What Changes

- Capture a reviewable, capability-based product contract for WFS first, with separate future applications and a low-footprint Home Server goal.
- Use the owner's explicitly available standard OTel event instrumentation: discover emitted events, map them to steps, choose a goal event from the list, and inspect ordered journeys from their known beginning to it. Keep standard reception independent of optional context.
- Define configurable journey semantics, ordered paths, honest outcomes, historical/version comparisons, live views, diagnostic linkage, viewer access, privacy and bounded retention.
- Research maintained alternatives, receiver/storage/frontend tradeoffs, reasons not to build and narrow stop/go gates before stack selection.
- Propose phased end-to-end tracer bullets and a small test-first engineering harness. All tasks remain unstarted proposals.
- Record Docker Compose with its associated Compose file and `.env`, behind the existing shared Cloudflare Tunnel HTTPS route, as the deployment contract. Qualify route transports and preserve application auth; no actual setup files are created.
- Deliver the complete WFS workflow as release A, including maps, timing, saved reference views and error recovery. Keep multi-app/richer comparisons/invitations/configuration as deferred roadmap B with stable app/version identity upfront.
- Record selected automatic boundaries, narrow future invite scope, 30-day detail/12-month summaries and Svelte/Vite + Go + SQLite. Add a discovery-led integration guide for agents helping developers propose standard event/warning/error/recovery instrumentation.

This is planning only. It authorizes no application or harness implementation, deployment, merge, production tests or WFS changes. Numerical defaults marked proposed are review choices, not accepted product requirements.

## Capabilities

### New Capabilities

- `telemetry-intake`: standard signal/transport compatibility, configured app attribution, bounded durable acceptance and delivery uncertainty.
- `journey-modeling`: event discovery, optional context mapping, visit and journey boundaries, deterministic versioned steps/goals/outcomes.
- `journey-exploration`: ordered historical paths, honest counts, loops/stopping analysis and normalized version comparisons.
- `live-journeys`: selected-app or isolated-visit streaming, stable filters, clean window and separately counted fixed baseline.
- `journey-diagnostics`: linked exceptions, recurrence, continued activity/recovery and scrubbed copyable context.
- `viewer-access`: passphrase access and scoped cookie-based invitations, separated from application correlation.
- `data-lifecycle`: minimization, 30-day ordered detail, derived-data correction/deletion and bounded Home Server operation.
- `app-integration-guidance`: discovery-led developer collaboration and reviewed instrumentation plans using existing OTel, separate from the internal harness.

### Modified Capabilities

None. The durable canonical inventory is empty. This approved planning reconciliation evolves the existing `synthetic-journey-outcomes` change, retiring its obsolete journey-evidence/journey-outcomes deltas. Its seven-day detail and sole first-vote defaults are superseded. The historical change identifier remains; there is one active plan and no canonical spec sync/archive.

## Impact

Only documentation/OpenSpec planning changes. Svelte/Vite + Go + SQLite with maintained OTLP components is selected; receiver packaging, pins, numerical limits and qualification remain technical defaults. The first release is the complete WFS workflow, not a list-only scaffold. This may be published as an approved documentation-only draft PR; it implements no app/harness, actual Compose/`.env`, CI, external instrumentation, canonical sync/archive, merge or deployment.
