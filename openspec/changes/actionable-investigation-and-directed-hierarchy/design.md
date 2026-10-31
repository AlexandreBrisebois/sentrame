# Design

## Context

See proposal.md for why and research.md for the source/history audit, remediation matrix, dated primary sources and exact fresh-build observation blocker. Baseline: `main` at `c7288d572811979fecc2f3f856a6a072d85395f1`. The earlier investigation/compact-header work and Tools are archived in this checkout; pending native acceptance remains pending. This change neither reopens nor rewrites those records.

The exact session filter survives under **Time & filters → Advanced → Correlation ID**. Matching journeys searches path labels only. Journey timeline/story display identifiers without a consistent follow-scope action. Current hierarchy already has arrows and SCC-based layered stage layout. The design addresses hidden capability and evidence semantics, not an established deletion or absence of arrowheads.

The WFS designer skill and judgment reference inform task → state → action → recovery review. Its visual reference is WFS-specific; Sentrame retains its own warm-paper/forest-green theme, serif headings, outlined surfaces and linked-thread icon from docs/design/visual-theme/README.md. No style migration is proposed.

## Goals / Non-Goals

**Goals:** complete “paste this ID and investigate,” “show the evidence for this pattern,” and “explain which directed path occurred.” Every view must reach useful evidence, a safe next action or an explained stopping point.

**Non-goals:** invent conversion goals, a new identity system, replay collection, full APM, force a linear chronology on concurrent APIs, redesign all seven views, or implement anything now. The source-backed matrix is a scoped audit, not proof that all app interactions were runtime-tested.

## Decisions

### 1. Find ID is visible; discovery remains ID-free

Add a compact **Find ID / Rechercher un identifiant** action in the existing toolbar/disclosure pattern, with an explicit **Session or operation / Journey / Trace** type selector. Reuse the exact `scope` query for anonymous session/operation; add bounded exact journey lookup only where protected retained-query context supports it; use protected exact trace detail instead of depending on a capped dropdown. Keep the old advanced scope field as a synchronized alias during migration, not a second independent query state. Remove it only if equivalent access and keyboard tasks are verified.

Use exact case-sensitive opaque identity. Strip only accidental leading/trailing whitespace with visible accepted value; reject control characters and over-limit input (proposed maximum 256 UTF-8 bytes, never truncating to another ID). Report existing narrower server constraints honestly. No regex, fuzzy, cross-app or automatic type guessing. Type-specific actions beside IDs prefill the correct type, so common tasks avoid choosing it manually. A journey ID opens a derived journey under its capture/version context; it is not assumed globally timeless. Trace IDs open technical structure, not a visit.

Submitting shows a compact accepted narrowing chip with type, safe value/copy, clear action and visible app/date/version constraints. Choosing a session scope must preserve browser tabs and multiple captured journeys rather than merge them. Clear removes only that narrowing. Large values wrap in detail and use a bounded accessible representation in the toolbar.

| State | What the user can know | Next action |
| --- | --- | --- |
| Invalid type/format/length | Input cannot be submitted under this type | Correct it; retain other filters |
| No authorized retained match under query | Current app/time/version constraints found no evidence | Inspect constraints; explicitly widen retained dates or clear one filter |
| Retained records without journey projection | Ingestion succeeded, modeling prerequisites absent | Open received evidence/discovery; owner may inspect mappings |
| Known expired/deleted target | Only if current authorized lifecycle evidence establishes it | Return to query or supported aggregates; never resurrect raw evidence |
| Unknown completeness/capped source | The current query cannot establish full coverage | Narrow/expand supported bounds; show which cap applies |
| Read/auth failure | Search result is unavailable, not zero | Retry or reauthenticate; remove inaccessible cached detail |

Do not issue “try every app/all time” queries to diagnose an ID. No-match messages cannot promise a cause (expiry, sender absence) without evidence. Do not reveal whether an ID exists outside access. Session IDs are pseudonymous correlators; search history must not become an identity directory.

### 2. Persist the question, not raw evidence

Single accepted-query state owns app/generation, ID type/value, time/version/definition, population/window policy, view, criterion and selected occurrence. Draft edits remain distinct from accepted filters. In-memory navigation and `history.state` preserve back/forward. A minimal versioned `sessionStorage` record may restore reload in the same browser tab, keyed by authorized app/generation and cleared on sign-out, revoke, app/generation change or invalid schema. No raw payloads, journey snapshots or credentials enter this state. Revalidate at every read; session storage is not authorization.

Do not put raw scope/journey IDs in new query strings, fragments or referrer-bearing links. Existing app/trace/reference link behavior is not expanded silently; exact-ID reload uses the above local mechanism. Explicit saved references can persist supported safe ID type/filter and criterion under existing owner authorization, lifecycle and version rules. Reference identity is a reusable query, not an immutable evidence archive or a public share. A temporary viewer can explore/copy locally but cannot save a definition or broaden scope.

Trace drilldown carries a return token/context; app-wide trace browsing clearly states its broader **authorized** context. Returning restores valid selection, graph viewport, scroll and focus. On refresh, late data cannot overwrite newer query/access state; expiry may remove selected detail and returns focus to a labeled surviving control. Baseline/live state remains separate, and saving a live selection retains its existing fixed historical eligibility policy.

### 3. Give each investigation surface a completion task

| Surface | Question | Evidence/action that completes it |
| --- | --- | --- |
| Graph/hierarchy | Which observed path/pattern warrants inspection? | Select node/edge/prefix → supporting journeys with exact criterion |
| Matching journeys | Which retained example supports/refutes this observation? | Show matched occurrence(s), choose an example → focused story, follow its scope, return |
| Journey timeline | What happened around this occurrence, including retry/recovery? | Open the chosen journey at the matching occurrence → inspect diagnostic/request → return |
| Trace/dependencies | What technically executed, failed or overlapped? | Native structure/timing → safe evidence copy and return to journey |
| Received/unmapped evidence | Was data received but excluded from meaning? | Explain missing context/policy/mapping; owner-only configuration stays explicit |
| Saved/reference view | Can I resume the same question? | Restore query/criterion/version, explain changed coverage or expired detail |

Timeline with no selected journey shows a bounded criterion-aware chooser, not an empty unexplained panel. Graph stage/transition selection carries into it. Do not auto-select a “representative” journey; existing newest-first examples remain explicit and deterministic. Matching rows show enough path/time/outcome plus match occurrence count to choose; IDs/technical clocks stay in accessible details with meaningful Open/Copy/Follow actions.

Timeline focuses matched occurrence(s), retains chronological context and offers previous/next match. Grouping repeats remains expandable with original count/evidence; it cannot hide the matched event. Error → recovery remains traversable; unplaced diagnostics do not acquire causal placement. If no linked trace exists, show that fact and next useful received-evidence action, not a disabled unexplained button.

Success means Alex can explain the selected observation, open its supporting evidence and return or copy a safe reproducible summary. A configured goal appears by name/rule/evidence. Goal-free observations are valid; last retained event, inactivity, exit and cutoff are not conversion success/failure. “Matching journeys” is evidence selection, not a new goal configuration popup. Configuration stays optional and owner-scoped.

### 4. Recommend a directed ordered-prefix hierarchy

Use the existing bounded prefix-tree evidence model as the basis, reconciled with graph selection and stage-wide identity. Nodes are **full uninterrupted prefix + cohort/definition/projection context**, not merely stage ID. Edges advance one retained ordered occurrence and carry supporting membership. Display stage identity separately so selection can distinguish “this prefix occurrence” from “all occurrences of this stage.” Canvas x-rank is ordinal progression, not elapsed time or causal depth. Node position may be adjusted vertically under existing hierarchy arrangement rules without changing evidence.

Example:

```text
Start → A → B → A → C → [captured boundary]
          ↘ D → C → [captured boundary]
```

The two C nodes are the same semantic stage in different prefixes. The second A is an occurrence of A, not a cycle folded into one SCC column. No rejoin connector merges the Cs. A stage-wide cross-reference can list both occurrences and their paths without creating a traversable link. For A→B and X→B→C, the B under A cannot gain C simply because aggregate stage edges include B→C.

| Alternative | Useful aspect | Why not primary hierarchy |
| --- | --- | --- |
| Current directed layered aggregate graph | Compact topology, branch/rejoin overview, existing renderer | Collapses repeats and can suggest unsupported cross-prefix paths; retain as explicitly stage-level transition overview where useful |
| Stage × ordinal DAG | Displays repeats in forward layers with shared stages | Ordinal merging still risks joining disjoint support; lineage would need extra provenance |
| Sankey | Traffic volume under a clear linear counting scale | Cycles require unfolding/custom return handling; volume is not session chronology or goal probability |
| Prefix trie (recommended) | Exact prefixes, repeats, counts and evidence aligned with ID lookup | Duplicates suffixes and grows; accept this cost with bounds and drilldown |

A general layered graph can contain cycles; the layout's cycle-breaking choice cannot reverse evidence arrows. Ordinary D3 Sankey expects a DAG, whereas Sentrame already has custom loop handling. No new layout library is selected or required. A prefix trie is a deliberate projection, not a claim the application process itself is a hierarchy.

Start markers distinguish captured start from left-truncated retained start. End markers distinguish configured goal, explicit exit, inactivity, active/right-censored and source/display cutoff. Goal marker only follows the configured meaning. Gaps terminate an uninterrupted prefix; hiding a stage cannot create direct observed adjacency. Do not put all dropped branches into a semantically traversable Other node.

Start with existing limits (32 displayed steps, 200 visible prefix nodes, five initial branches), expose each limit separately from source/population limits and provide bounded expansion/paging. Preserve counts/denominators using the authoritative selection, not only displayed nodes. Exact-ID investigation must reach retained rare/deep paths without first finding a top branch. Version/time aggregation groups compatible semantics and shows window/population; low-support/outlier claims remain observations, not defects.

An exact-session filter can yield several captured journeys/tabs; show them as separately labeled roots rather than stitch across boundaries. New API concurrency groups from the sibling change preserve uncertain display order; they cannot become a causal tree. If unsupported by this hierarchy projection, offer the ordered interval/list detail with an explicit limit rather than flatten concurrent spans into successive user steps.

Switching to Force/Sankey preserves semantic stage/criterion/filter, while marking those as aggregate transition views. Prefix occurrence focus is carried in the inspector or explained as not a single aggregate node. Keyboard/list alternatives expose full prefix, units, limits, matching evidence and actions. This preserves equivalent investigation without requiring graph gestures.

### 5. Preserve the graph-first interaction shell

Keep seven direct view controls, Tools/Outils grouping, one-open disclosure behavior and current graph bounds. No permanent purpose paragraph or extra control row above the canvas. Put the question in a compact local view/inspector caption and explain details in Help/Context/Coverage. Material sampling/partial/cap warnings remain concise and visible. Find ID opens a bounded labeled disclosure with Close/Escape/outside dismissal, type/input labels, visible focus and valid return. Targets at least 44px; EN/FR labels can wrap. Do not copy WFS visual tokens.

## Risks / Trade-offs

- [Duplicated suffix nodes use more space] → initial bounds, clear same-stage labels, exact-ID bypass and stage-wide evidence cross-reference; no false rejoin.
- [ID types increase choices] → typed actions beside app-provided IDs and a persistent explicit selector; no ambiguous silent guessing.
- [Local persistence retains a correlator] → minimal tab-scoped state with explicit clear/access invalidation; no new raw-ID URL paths.
- [App already meets some contracts in source] → retain working path filters/paging/context, implement only evidenced gaps and verify rendered behavior later.
- [Unavailable runtime evidence] → do not call source findings reproduced UI defects; acceptance requires a verified isolated build and actual interactions.

## Migration Plan

Following a separate implementation instruction, first unify typed query state and ID actions without changing mappings/boundaries. Then wire matching/timeline/diagnostic return context, then prefix hierarchy with structured alternative. Migrate saved preference schema additively; older references default to supported original behavior and invalid preferences explain fallback. Keep aggregate layout/reference readers usable. Rollback new presentation/state fields without deleting references or evidence; ignore unknown optional fields safely. Native acceptance from earlier archived changes stays independently pending and is not checked off by this work.

## Open Questions

Deferrable wording/presentation choices: final EN/FR labels; initial branch-expansion count within the declared cap; and the first representative usability participants/tasks. The model, exact typed lookup and privacy-preserving persistence are recommendations fixed for this proposal, not unresolved architecture. Alex may revise them during spec review before implementation; doing so requires reconciling deltas/tasks. Physical-device/native accessibility checks and the denied local observation are verification gaps, not product approval gates to writing these specs.

## Acceptance tasks and evidence quality

Priorities follow research.md: first exact-ID → scoped graph and recovery; then relevant matching example → focused timeline → trace → back; then prefix evidence without invented paths; then reusable references and less frequent path-table routes. Verify with synthetic IDs shared across two apps, multiple tabs/journeys in a scope, expired/unmodeled/missing data, references/reload, retries/loops, separate branches sharing suffixes, low support, caps and revoked access. Record task outcome, steps/time and wrong turns in an uncoached walkthrough; do not invent a numerical usability result now. Native EN/FR screen-reader, 200% native zoom and physical touch remain explicit future verification cells. A successful build or mocked DOM check does not establish those outcomes or blanket WCAG conformance.
