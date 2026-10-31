# Proposal

## Why

Open-source installers should be able to opt into telemetry without obtaining Sentrame credentials manually. A public repository or container cannot keep embedded signing material secret, so convenient activation and stable release aggregation need an explicit distinction between public intake and trusted producer authority.

## What Changes

**Exploration only, 2026-10-31.** These are candidate behavioral boundaries for review, not approval to implement or a selected enrollment protocol.

- Compare keyless routing by public app identifier, baked public write-only identifiers, HTTPS activation that returns current public ingestion configuration, per-installation enrollment, and enrollment authorized by a trusted server.
- Separate stable registered app identity, current intake generation, optional installation identity, service/version metadata and revocable write permission. Sharing a credential is unnecessary for cross-version aggregation.
- Define opt-in, standard OTLP compatibility, public-versus-trusted authority, bounded abuse handling and recovery boundaries that any selected option must meet.
- Record current-code evidence, primary-source research, a conditional recommendation and unresolved product choices in research/design/handoff artifacts.
- Preserve existing owner/read APIs, private server credentials, generation/deletion barriers and current Tools implementation. No implementation, credential creation/configuration, live receiver traffic, shared-resource mutation, staging, commit or push.

## Capabilities

### New Capabilities

- `telemetry-activation`: Proposed consented public configuration discovery and optional enrollment, with explicit trust, lifecycle and ordinary-exporter boundaries. Public app aliases are not owner credentials.

### Modified Capabilities

None in this exploration. Existing `app-management`, `telemetry-intake` and `data-lifecycle` contracts remain authoritative. The proposed capability is additive and disabled unless later approved; choosing an implementation will require targeted amendments where current intake only accepts browser/server credentials. See design.md for this compatibility gap rather than treating the new delta as an implemented route.

## Impact

Potential later integration points are `internal/app/registry.go`, `apps_http.go`, `http.go`, `intake.go`, authority-aware analysis and integration guidance. Public exposure may need an operator-owned admission gateway and separate budgets; current Home Server limits do not qualify Internet-scale intake. No dependency, API, deployment or data schema changes are made here. Product decisions about receiver ownership, scale/trust, consent, acceptable attribution and restore/revocation semantics remain open.
