# Release scope and preserved roadmap

The main choices are settled. These are planning requirements, not implemented capabilities or authorization to start application work.

## Release A: complete useful WFS journey workflow

- Standard OTLP intake with an honest origin/routed support matrix, discovered events and simple event-to-step/goal mapping.
- Historical and live **maps**, ordered evidence, rare/stopping context and warning/exception continuation/recovery through completion.
- Clean-live toggle and the selected app's accumulated seven-day traffic behind the selected live anonymous visit, with distinct counts.
- Accepted time-at-stage analysis: node size for duration, edge thickness traffic, median plus slower cases, timing quality and accessible numeric/table alternatives.
- Automatic first-mapped-event start, active-goal end, next-event restart and 30-minute inactivity closure as inactive, not failed. Browser tabs stay separate; server-only logical scopes do not need tabs; correctly correlated client/server evidence is supported.
- Saved goal-backward, bounded beginning/end-segment and origin-forward reference views over recorded boundaries, with honest partial/expired evidence.
- Owner passphrase access, privacy minimization, 30-day detailed journeys and 12-month count/conversion/timing summaries without individual drilldown after detail expires.
- A discovery-led public integration guide for agents helping developers instrument other apps, separate from Sentrame's own engineering harness.

An early one-row/list tracer bullet proves engineering seams only. It does **not** replace the agreed first-release graph workflow or declare the release finished.

## Roadmap B: retained, deferred requirements

Multi-app experience, richer parallel version comparisons, viewer invitations and broader configuration remain specified for later. Keep stable app/version/definition identities and compatibility from day one; do not build an extension/plugin framework merely to anticipate them.

The selected future invitation design is read-only access to **one live session**, with the owner explicitly choosing whether its app baseline and exception details are included, no configuration edits, expiry and revocation. Broader read-only app-workspace/history/diagnostic access remains a later mode. Invitation lifetimes and redemption defaults remain technical proposals. Invitations are deferred from A even though their permission design is already selected.

## Selected architecture and deployment

Svelte/Vite static frontend + Go + SQLite, reusing maintained OTLP components, is selected. Docker Compose with an associated Compose file/`.env`, behind the existing shared Cloudflare Tunnel HTTPS route, is selected. RAM is reported sufficient and CPU options available; no specific hardware, unlimited budget or measured fit is inferred. Bounded queues, storage, query/stream/render work and later qualification remain required.

Receiver embedding versus a maintained Collector sidecar, image pins, origin transport, exact intake/late-data limits, summary grains/distribution precision, invitation defaults and timing thresholds remain implementation-level design review items. No Compose/`.env`, credentials, network, app, harness or CI setup is implemented by this package.

## Retention and query limits

Raw ordered evidence expires at 30 days. Aggregates retain 12 months of supported counts, conversion rates and timing distributions. Do not average medians to combine periods or imply that pairwise counters reconstruct paths. Saved views retain query definitions, not raw journeys. Beyond detail retention, show supported summaries and disable individual drilldown. Declare correction/deletion horizons and summary dimensions; do not retain hidden person/visit membership merely to promise unlimited selective deletion.

Selecting an operational goal changes a versioned segmentation definition. Choosing a different goal, segment or origin in a saved retrospective view queries existing boundaries; it does not silently resegment or join captured journeys. If endpoints cross a goal/timeout boundary, preserve that boundary and explain the limited result. This narrow interpretation must remain visible in the interface and integration guidance.

## Owner-approved extensions on2026-10-28

The owner explicitly authorized whole-app temporary read-only invitations with expiry and revocation. This separate mode includes retained historical/live data, references and sanitized diagnostics. It does not silently widen the deferred visit-only grant, and does not include shared configuration/saving/deletion rights. Multi-app and richer comparisons/configuration remain deferred.

The owner also requested multi-architecture DockerHub delivery (`brisebois/sentrame`), a GitHub Actions publishing workflow, simple Compose and an environment example using a chosen local filesystem bind folder. These source/configuration artifacts do not by themselves establish publication, container runtime support, NAS fit, real application export or the shared Tunnel route. The initial no-commit/no-push constraint remains until separately changed.
