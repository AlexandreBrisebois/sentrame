# Minimal engineering harness — proposal only

This document proposes future engineering rules/entrypoints. It installs no AGENTS file, Taskfile, scripts, CI, dependencies or application scaffold. Creating a planning artifact does not authorize its implementation.

## What actual WFS patterns teach

WFS's [contract/testing owner](https://github.com/AlexandreBrisebois/whats-for-supper/blob/0dc61dfce9a7bf7d55eafff08f57c10e3f237f7c/.agents/core/contract-testing.md) explicitly requires contract → tests → implementation, observed red → green → refactor, a first tracer bullet through affected real seams, scoped follow-on slices and honest check evidence. Its [mission](https://github.com/AlexandreBrisebois/whats-for-supper/blob/0dc61dfce9a7bf7d55eafff08f57c10e3f237f7c/.agents/core/mission.md) favors existing owners/reuse and rejects speculative abstractions. Its [context-loading guidance](https://github.com/AlexandreBrisebois/whats-for-supper/blob/0dc61dfce9a7bf7d55eafff08f57c10e3f237f7c/.agents/core/context-loading.md) loads focused context and preserves unrelated work.

The inspected [execution harness](https://github.com/AlexandreBrisebois/whats-for-supper/blob/0dc61dfce9a7bf7d55eafff08f57c10e3f237f7c/.agents/core/execution-harness.md) classifies planning documents as documentation and separates static checks from real API/database evidence. [Taskfile](https://github.com/AlexandreBrisebois/whats-for-supper/blob/0dc61dfce9a7bf7d55eafff08f57c10e3f237f7c/Taskfile.yml) offers focused tests and explicit broader verification; some development gates also stop processes. These are source observations, not commands to run in Sentrame. Sentrame has none of that harness today.

## Proposed compact agent rules

1. Read the selected capability, task and immediate source/config/test owners; inspect worktree state. Load more only to resolve a named uncertainty. User scope and the selected OpenSpec workflow control writes; artifacts are not implementation permission.
2. State one observable scenario, its evidence and stopping point. Reuse an existing rule/module/library before creating another owner. Keep changes minimal; no generic framework merely for anticipated future apps.
3. Define/reconcile public request/response/error/SSE behavior before logic. Use official OTLP schemas and maintained receivers; a proprietary JSON approximation cannot qualify OTLP. Introduce a small explicit viewer API contract only when there is an API to build.
4. **Red:** write the smallest meaningful behavior test, run it before production logic changes, and observe the intended missing/incorrect behavior. A compiler/tool/environment failure is not red. If the test already passes, establish the missing behavior first.
5. **Green:** implement only that scenario, run its focused test and affected seams against a real temporary store. **Refactor:** improve within the selected scope while preserving green; rerun checks affected by the refactor. Record actual red/green commands/results.
6. Start with a thin tracer bullet across only necessary layers: standard exporter → receiver → durable sanitized evidence → protected read path → one explainable list/sequence view. Do not build all backend, then all UI, then integrate. External/provider doubles stay at the boundary; mocked UI API results do not prove receiver persistence.
7. Add one behavior at a time: healthy loop, explicit exit, recovery after error, late correction, replay, expiry, live scope. Each slice includes its own tests and documentation. Do not defer all tests or privacy checks to a final phase.
8. Bound test fixtures, request sizes, clocks, temp directories, ports and owned processes. Never reset shared data or stop unrelated services to obtain green. A worktree does not isolate ports/DBs.
9. Select proportionate validation from actual impact. Prose changes get link/format/spec checks; deterministic logic gets focused cases; intake/storage/access/SSE changes get real affected seam tests. Broaden for shared/unknown impact, but do not rerun the entire browser suite after every tiny edit without a reason.
10. Before completion, review actual diff/ownership and name checks as passed, failed, blocked, not-run or not-applicable. Include tested content/config identity when useful. Do not weaken a contract/test or call a tooling failure passed. No automatic apply, sync/archive, merge or deployment follows successful tests.

These rules borrow engineering posture, not WFS's generated client tooling, framework, mock-file layout, deployment scripts, large fixture inventory or access identity. A small repository should not need a harness whose maintenance exceeds its behavior code.

## Proposed entrypoints after implementation approval

| Logical command | Scope and evidence | Introduce when |
| --- | --- | --- |
| `spec:check` | Pinned OpenSpec strict validation for the selected change, plus full-repository reconciliation once concurrent drafting ends. | Reuse installed CLI now; only add a task wrapper when harness implementation is approved. |
| `test:unit <filter>` | Deterministic mapping/lifecycle/counting/time cases; no services. Observed red/green. | First selected behavior. |
| `test:seam <scenario>` | Official exporter into real receiver, temporary DB and same public read/SSE path; owned loopback endpoints and controlled clock. | First tracer bullet. |
| `check` | Affected format/types/static checks and focused tests; no implicit installs, service resets or process kills. | Toolchain chosen and app code exists. |
| `verify <scope>` | Explicit union of required checks for reviewed changed paths and shared contract/DB/auth impact. Report exclusions and prerequisites. | More than one slice makes check selection valuable. |
| `demo:synthetic <case>` | One bounded generated case with expected outcomes through the real local path; never production credentials/WFS endpoints. | Thin tracer exists and is approved. |

Taskfile, Make or plain documented scripts are alternatives. No need to introduce a dependency merely to implement these names. Do not build an adaptive test selector/cache before the checks are costly enough to justify it. Use existing test-runner filtering first. Avoid storing verbose logs, entire source dumps or repeated narrative in task evidence; a command/result, relevant identity and exception is enough.

## Acceptance examples for the harness itself

- A deliberately broken goal rule causes its targeted behavior assertion to fail before the fix; a missing database process is recorded blocked rather than red.
- The first tracer test reads the committed evidence through the protected query API after restart; a fixture preloaded directly into a mocked frontend is insufficient.
- A docs-only capability edit runs strict OpenSpec/link checks without launching browsers or databases.
- A viewer scope change adds server-side query and SSE-replay authorization checks; hiding a button alone does not satisfy it.
- A failed or unavailable required check leaves the selected task unverified, with the blocker named. Passing unrelated checks cannot compensate.

All examples are proposed future checks. This research pass verifies document structure and consistency only.

## Release and integration boundaries

A row/sequence tracer is the first engineering increment toward the complete WFS historical/live graph release, including timing, recovery and 12-month summaries. It does not replace accepted release A or delete retained roadmap B. The [external integration-agent guide](../integration/AGENT-INTEGRATION.md) serves developers of other apps through discovery, questions and reviewed event/warning/error plans; it is separate from this internal harness proposal and adds no AI runtime.
