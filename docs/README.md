# Sentrame documentation

## Use and operate the application

- [User guide](USER-GUIDE.md): graph views, selection, event mapping, time filters, highlights, references, errors, data lifecycle and temporary viewers.
- [Developer integration guide](DEVELOPER-GUIDE.md): .NET and Next.js telemetry, direct export versus Collector/proxy routing, trust boundaries, exceptions and journey context.
- [OTel event/activity contract](integration/OTEL-CONTRACT.md): standard native spans/logs first, contextual mappings, optional stage hints and capability-specific scope/presence.
- [Container deployment guide](NAS.md): Docker Hub workflow, Compose, environment configuration, filesystem bind storage, health and upgrades.
- [Local setup and demonstrations](LOCAL-REVIEW.md): source startup, synthetic scenarios and VS Code REST Client requests.
- [Docker Hub description](DOCKERHUB.md): repository description copy.

## Contracts and implementation

- [Canonical capability specifications](../openspec/specs/): current normative baseline.
- [Archived standard OTel integration change](../openspec/changes/archive/2026-10-29-standard-otel-default-integration/): completed implementation, synced capability deltas and qualification tasks.
- [Archived synthetic journey baseline](../openspec/changes/archive/2026-10-29-synthetic-journey-outcomes/): historical planning and implementation record.
- [Implementation contract](../CONTRACT.md): API/data boundaries and current integration contracts.
- [Release scope](planning/sentrame-release-scope.md) and [decision register](planning/sentrame-decisions.md): settled direction, explicit later scope additions and remaining roadmap.
- [Visual theme](design/visual-theme/README.md): preserved colors, typography and linked-thread mark.
- [Integration instructions for coding agents](integration/AGENT-INTEGRATION.md): discover an application's current instrumentation and agree an event/privacy plan before editing it. This complements the runnable-configuration discussion in the developer guide.

## Verification and historical rationale

- [Standard OTel qualification](verification/standard-otel.md): current model, progressive voting, exact wire/SDK/browser cells and remaining deployment limitations.

- [Workspace verification](verification/workspace-design.md): graph/UI, diagnostic, invitation and container checks, including failed attempts, corrections and limits.
- [Initial Release A verification](verification/release-a-local.md): earlier technical qualification; later workspace/delivery evidence extends it.
- [Technical defaults](review/release-a-technical-defaults.md): concrete receiver, timing, identity, retention and resource choices.
- [Previous MVP lessons](review/previous-mvp-lessons.md): reusable scenarios and lessons; not coverage of the current implementation.
- [Research](research/sentrame-feasibility.md), [stack/resource analysis](research/sentrame-stack-and-resources.md), [engineering harness proposal](proposals/sentrame-engineering-harness.md), and [original Home Server deployment contract](proposals/home-server-deployment-contract.md): historical rationale and proposals. Statements that no app or Compose files exist describe their original planning phase, not current repository contents.

A requirement, a code path and a passed qualification are different things. Operational guides describe current behavior; verification reports state what was actually exercised and where support remains unqualified.

[App workspaces](MULTI-APP.md) · [Multi-app qualification](verification/multi-app-workspaces.md)
