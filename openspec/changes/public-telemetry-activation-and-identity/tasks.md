# Tasks — decision and research follow-up only

This is an exploration, not an implementation backlog or authorization to apply. All checkboxes are future decision/qualification planning work. The research and four planning artifacts are complete; unchecked tasks record remaining product decisions, not missing code. Do not generate credentials, start services, run production tests, change shared resources, stage, commit, push or archive under this request.

## 1. Product ownership and trust decision

- [ ] 1.1 Ask Alex whose receiver receives data and who operates/pays for it; verify completion by an explicit recorded answer distinguishing central and installer-owned modes.
- [ ] 1.2 Agree expected installations, event/byte rates, burst/cardinality and cost ceilings plus diagnostic versus trusted-decision use; verify a concrete workload/budget and acceptable loss statement, not an inferred Internet-scale capacity claim.
- [ ] 1.3 Agree acceptable attribution confidence and whether per-install revocation is necessary; verify the selected assurance names what is authenticated and what remains forgeable, including admission evidence if D/E is selected.
- [ ] 1.4 Agree opt-in actor/flow, minimized fields, purpose, disablement and destination-change consent; verify a reviewable consent/data example without instrumenting or sending telemetry.
- [ ] 1.5 Clarify stable alias/identity recovery versus revoked-token reactivation and receiver ownership changes; verify an explicit lifecycle decision covering stale backups, deleted apps and source of latest deny state.

## 2. Refine only the selected candidate after answers

- [ ] 2.1 Record Alex's chosen option(s) and rejected alternatives in the design, keeping current implementation authorization unchanged; verify every material question above has an answer or is explicitly excluded from scope.
- [ ] 2.2 Specify finite activation/cache/refresh/offline/queue policies and any enrollment admission/replay/retry rules for that option; verify a paper state table for first activation, expiry, 401/403, overload, lost response and withdrawal.
- [ ] 2.3 Define the selected option's public-assurance propagation, fairness/rate/cardinality/storage budgets and emergency disable behavior; verify a paper abuse matrix against spoofed names/versions, colliding IDs, signup Sybils and resource floods without contacting a receiver.
- [ ] 2.4 Define authorized same-identity recovery, rebind, stale backup and rollback behavior with current generations/tombstones; verify a paper failure matrix including missing deny state, revoked token, deleted app, mismatched stores and changed receiver.

## 3. Prepare a separately authorized implementation proposal

- [ ] 3.1 After choices are explicit, draft targeted amendments to existing app-management/intake/lifecycle and assurance-aware analysis contracts where necessary; verify no conflict with mandatory current browser/server bearer binding or secret handling, and strictly validate the revised artifacts.
- [ ] 3.2 Name target SDK/language versions and standard exporter/Collector refresh approach; verify a documented compatibility plan for endpoint/header encoding, startup/reconfiguration, binary/JSON/gzip and route/browser limitations without claiming unrun tests.
- [ ] 3.3 Prepare synthetic acceptance cases for opt-in, protected API denial, assurance separation, scope/replay, rotation/cache, abuse budgets and coordinated recovery; verify cases map to delta scenarios and require isolated disposable resources only under later authorization.

## Workflow follow-up

- Present the selected/refined planning artifacts to Alex. Implementation requires a new explicit request; OpenSpec completeness is not approval.
- Do not apply, merge into main specs, archive, stage, commit, push, configure credentials or mutate receivers as part of this exploration.
