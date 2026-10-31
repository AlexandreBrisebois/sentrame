# Review handoff — public telemetry activation and identity

Date: 2026-10-31. Research/spec only; no architecture selected or implementation authorized.

## Result and reading order

1. `proposal.md`: scope and candidate capability.
2. `research.md`: executive findings, current-code references and dated primary links (Sentry SDK/client-key API, Application Insights, Azure Monitor OTLP, OTel, SPIRE, GitHub OIDC, RFC 9449).
3. `design.md`: option comparison, conditional recommendation, trust/abuse analysis, caching/enrollment/recovery cases and exact product questions.
4. `specs/telemetry-activation/spec.md`: nine candidate requirements with testable scenarios; scope/assurance boundaries, not a selected wire protocol.
5. `tasks.md`: twelve unchecked decision/refinement/qualification-planning tasks. None authorizes implementation. OpenSpec's complete/4-of-4 means artifacts exist, not decisions approved or features implemented.

## Main finding

A stable public app handle plus explicit opt-in HTTPS discovery can eliminate manual installer key acquisition and return current ordinary OTLP settings. It cannot authenticate an authentic installation. A public distributed signing key is available to spoofers. A self-generated installation key proves possession/continuity; trusted-server-signed enrollment proves only the issuer's admission claims under its protected policy. None proves truthful events or a real user action. Public intake must never grant reads, owner operations or protected server authority.

Prefer public discovery conditionally for low-risk central community diagnostics, with isolated public assurance and strict abuse budgets. Consider per-install enrollment for individual revocation and protected issuer admission for stronger operator/workload confidence. Installer-owned receiving remains a separate ownership option. No custom SDK fork is needed; bootstrap/refresh integration still needs language-specific qualification.

Version aggregation belongs to stable registry app identity plus bounded service/version metadata, independent of credential sharing. Current registry recovery is a private coordinated backup, not arbitrary key import/rebind. A durable alias must not revive deleted generations or revoked permission. A stale backup cannot know future revocations; missing latest deny state must fail closed.

## Ask Alex before selecting architecture

- Whose receiver should collect: the developer's central receiver, each installer's receiver, or an explicit choice? Who pays/operates it?
- What scale/budget and trust are expected: low-risk diagnostics or evidence supporting business/persistence decisions?
- What explicit opt-in actor/flow, minimized fields, purpose, disablement and destination-change consent are acceptable?
- Is public app/version attribution enough, or is pseudonymous installation continuity, per-install revocation or verified operator/workload identity required?
- Does forever restore/rebind mean stable identity/configuration continuity or reactivating revoked permission? Which latest revocation/deletion state survives disaster?
- What cache/offline/loss bounds are acceptable, and how do product releases map to service versions?

Do not interpret silence as a choice. Exact URLs, field schemas, token type, gateway location, enrollment issuer and numeric quotas remain unresolved.

## Verification and preservation

- Branch `main`; HEAD `c7288d572811979fecc2f3f856a6a072d85395f1` verified at start and finish.
- `openspec status --change public-telemetry-activation-and-identity --json`: planning complete, proposal/specs/design/tasks all done.
- `openspec validate public-telemetry-activation-and-identity --strict`: passed.
- Both preceding changes also pass strict validation. SHA-256 comparison confirms all 18 existing files in their directories unchanged. Their original handoffs were preserved.
- Tracked worktree and index remain clean; only the two preceding and this new change directory are untracked. No staging, commit, push, branch switch or archive.
- No implementation, credentials/signatures/key generation, credential configuration, live registry/data access, service startup, production tests or shared-resource mutation. Current Tools implementation unchanged.

## Limitations / blockers

No operational blocker or denied action occurred in this session. Material product decisions above remain open by design and block an implementation choice. No runtime, abuse/load, browser, route, recovery or SDK-refresh qualification was performed; tests were read only. Prior ephemeral loopback `listen EPERM` was not retried. One Sentry overview web fetch rejected Markdown content; maintained source/API evidence was used instead. Moving upstream source links must be rechecked/pinned during any later implementation proposal.

Next authorized step is review/discussion of these artifacts, not apply or deployment.
