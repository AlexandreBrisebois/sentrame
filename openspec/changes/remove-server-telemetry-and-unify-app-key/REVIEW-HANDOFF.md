# Review handoff — remove server telemetry and unify app key

Prepared 2026-10-31. Specification/design review only. No implementation or operational cutover is authorized by these artifacts.

## Result

Change: `remove-server-telemetry-and-unify-app-key`.

Read [proposal.md](proposal.md) for scope, [design.md](design.md) for decisions/trade-offs and migration, [inventory.md](inventory.md) for observed source and earlier-proposal reconciliation, and [tasks.md](tasks.md) for 25 unchecked future tasks. Eight capability deltas are under `specs/`: app-management, telemetry-intake, journey-modeling, journey-exploration, trace-exploration, journey-diagnostics, data-lifecycle and app-integration-guidance.

## Recommended review choices

1. **Client-only journey product, standard receipt:** remove server-only/combined journeys, persistence assurance, producer-role selectors and server recipes. Keep ordinary qualified OTLP logs/traces/metrics, useful client trace/errors and existing SDKs/exporters. Generic payloads cannot reliably prove physical client origin; valid SERVER-kind spans remain metadata, not a server product. Decide whether Alex also wants to remove the supporting Traces destination; the candidate retains it.
2. **One key:** retain each app's current browser key as its app telemetry key at a later explicitly authorized migration, and reject its old server key then. No dual-key grace period. Alternatively choose a fresh app key and coordinated client reconfiguration before implementation. App identity/generation/history do not change with key rotation.
3. **Preserve history safely:** keep legacy records, definitions/references and aggregates under their existing lifecycle, but exclude mixed journeys and inseparable summary intervals from current comparisons with explicit compatibility states. Keep pure client history only where original boundaries/meaning/identity are equivalent. Do not fabricate A→C by deleting server B from A→B→C. No legacy server explorer is proposed.

These choices are stated explicitly in the design and gate future apply. Silence is not approval to retire credentials, migrate state or make incompatible history unavailable. The one key is write-only app attribution for spoofable observations; it cannot prove a genuine installation, truthful human action or durable backend commitment. Reads, administration and temporary-viewer grants remain independently protected.

## What stays useful

Sentrame's own Go backend, OTLP receiver, SQLite state, owner authentication, private JSON registry and lean Home Server Compose deployment remain. Preserve app isolation, optional goals, goal-free navigation, healthy loops, timing limits, client diagnostics, historical/live bounds and baselines. Preserve shipped Tools/Outils, all seven journey views, compact graph space, stable refresh, EN/FR and accessible evidence navigation.

The investigation/hierarchy proposal remains independently useful after removing server-operation/dependency examples. Public activation research still matters for opt-in, receiver ownership, scale/budget, public spoofability, stable identity versus versions, revocation and coordinated recovery; public discovery/enrollment/per-install keys are not included in this change. Existing proposals are left untouched.

## Exact limitation

Action: inspect and reconcile `openspec/changes/server-api-journeys-and-dependency-timing/` in full. Target is absent from this checkout, archives and reachable Git path history, and was not found in `/tmp` or local Codex worktrees. The actionable-investigation handoff contains a dangling link to it. Outcome: **blocked by unavailable source**, not denied, pending approval or uncertain execution. Its requirements cannot be claimed reviewed. The new explicit direction supersedes its server/API/dependency expansion scope; exact requirement disposition awaits the actual files. Other specification work is complete.

## Validation and preservation evidence

- `openspec validate remove-server-telemetry-and-unify-app-key --strict`: passed after preserving canonical scenario identifiers and updating their assertions for the replacement behavior.
- `openspec status --change remove-server-telemetry-and-unify-app-key`: 4/4 planning artifacts complete. This establishes artifact presence, not approved decisions or implemented features.
- Eight delta files: 40 requirement entries (including the removed server-recipe requirement), 103 scenario entries. All 25 future tasks remain unchecked.
- SHA-256 comparison of all 25 baseline files across canonical `openspec/specs` and the two pre-existing untracked change directories: unchanged.
- Branch/HEAD remain `main` / `c7288d572811979fecc2f3f856a6a072d85395f1`; tracked diff and staged diff are empty. Only the same two prior proposal directories and this new proposal directory are untracked.
- Whitespace checked on the new files; `git diff --check` is clean for tracked state. No application test run is claimed: this is documentation-only planning. Runtime, migration, standard-exporter, actual browser, native zoom, container, Home Server, Tunnel, abuse/load and recovery qualification remain future checks.
- No application edits, staging, commit, push, branch switch, deployment, credential reveal/generation/rotation/revocation, database read/migration/deletion, canonical sync or archive. No service was launched and no approval denial occurred.

Next step is Alex's design review and decisions A–C. A separate later instruction is required to implement the reviewed scope.
