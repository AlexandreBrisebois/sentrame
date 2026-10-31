# Proposal

## Why

The workspace hamburger hides unrelated controls behind an unnamed visual symbol. People must inspect a long, ungrouped list to distinguish graph presentation, investigation tools and app administration; dismissal and keyboard focus also behave differently from neighboring flyouts.

## What Changes

- Replace the hamburger with a visible **Tools ▾ / Outils ▾** disclosure in the existing graph toolbar, preserving graph prominence and the separate selected-app dropdown.
- Group existing destinations under View, Investigate and Manage this app; retain every authorized action and make selected presentation state explicit.
- Use one bounded, anchored nonmodal panel with predictable keyboard order, Close, Escape, outside dismissal and focus recovery; coordinate it with other flyouts.
- Preserve the seven direct graph views, Journeys/Traces navigation, app-scoped authorization and current compact header geometry.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `journey-exploration`: discoverable grouped workspace tools and consistent disclosure interactions without increasing header rows.

## Impact

Future implementation affects `web/src/Workspace.svelte`, associated styles/localization and focused browser checks; existing Disclosure/Modal behavior can inform reuse. No API, schema, credential or telemetry change is proposed. This change contains planning and review evidence only. Runtime implementation requires Alex's subsequent approval. Canonical specifications and archived October 31 records remain untouched.
