# Executive design review and handoff — 2026-10-31

## Recommendation

Restore the known-ID workflow with compact typed Find ID access and Open/Filter actions beside IDs. The session/operation field survives under Time & filters → Advanced → Correlation ID in current source and inspected history. A recent removal is not established; Matching journeys' path search does not search these IDs. Distinguish session/operation scope, journey and trace IDs rather than silently treating every code as a session.

Make Matching journeys select supporting evidence for a criterion, and Timeline explain the selected occurrence with next/previous match, Follow scope, diagnostic drilldown and reversible return. A useful task outcome is finding evidence and resuming the question; no invented conversion goal is needed. Preserve configured goal meaning and goal-free observations.

Use a bounded directed prefix hierarchy. The current graph already has arrows and SCC-based layering. Prefix identity preserves A-B-A and prevents aggregate edges from inventing A-B-C. Shared later stages appear separately in different prefixes, trading compact rejoins for honest path provenance. Keep stage-aggregate Force/Sankey views and mark changed units. Source/history evidence, prioritized remediation matrix and primary product/algorithm references are in [research.md](research.md).

## Trade-offs and remaining review inputs

Typed lookup adds an explicit choice, reduced by typed actions beside supplied IDs. Same-tab local state supports reload/backtracking without adding raw IDs to URLs, while saved references retain explicit owner-controlled reusable context. Prefix nodes duplicate suffixes and need bounds/expansion. The first release's exact wording/branch count and usability participants remain deferrable; architectural recommendations are resolved in design.md for review.

Retain Sentrame's forest/paper visual language, Tools/Outils, seven direct views and approximately165px desktop/201px phone graph start. Context remains in accessible disclosures; meaningful warnings remain concise and visible. No permanent explanatory wall or extra header row is proposed.

## Verification and exact blocker

- [Proposal](proposal.md), [design](design.md), [capability delta](specs/journey-exploration/spec.md), [tasks](tasks.md), [research/matrix](research.md).
- `openspec validate actionable-investigation-and-directed-hierarchy --strict`: passed. Status: 4/4 planning artifacts complete; future implementation boxes unchecked.
- Fresh Vite build at `c7288d5` succeeded; source/asset hashes and command are recorded in research.md. This was not a full typecheck, production test or runtime interaction verification.
- **Environment-denied action:** launching `/tmp/sentrame-investigation-research/observe.cjs` to bind a new ephemeral `127.0.0.1` synthetic preview failed with `listen EPERM` before browser launch. No listener, browser screenshot or API call completed. No escalation/automatic-approval review was requested; no pending or uncertain mutation remains. Existing19010 was not used.
- Findings remain source/history-backed. Native200% zoom, EN/FR screen reader, physical touch and an uncoached walkthrough are unobserved. Earlier pending native checks remain pending; build/spec validation does not establish accessibility conformance.

## Repository and implementation boundary

Entry/exit main at `c7288d572811979fecc2f3f856a6a072d85395f1`, initially clean. Only two new change directories are untracked; tracked/index diffs empty. No runtime code edits, staging, commits, pushes, branch changes, production tests, receiver mutations, WFS writes, credentials, sync or archive. Tools was already archived here and was preserved. No applicable repo/ancestor AGENTS.md or local designer skill was found; the requested WFS designer references were read as task/evidence guidance, with Sentrame's own visual language controlling presentation.

Review the completed artifacts; implementation and deployment remain separate future instructions. The companion [server/API review](../server-api-journeys-and-dependency-timing/REVIEW-HANDOFF.md) adds request/dependency semantics without authorizing WFS changes. A third public/keyless-ingestion exploration has a separate fresh-session brief; it is not part of this change.
