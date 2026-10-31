# Design

## Context

The snapshot already supplies captured journeys, ordered paths and weighted nodes/edges; the workspace shares stage and edge selection among three graph layouts. Ordered evidence expires before supported summaries. Saved references validate layout values and restore safe filters. See the proposal and delta for the desired behavior.

## Goals / Non-Goals

**Goals:** Build complementary analyses over the same bounded, filtered captured evidence, keeping existing views and references usable. Make partial coverage and count units explicit.

**Non-Goals:** Change journey segmentation, infer user identity or attention, reconstruct expired paths, change trace layouts, or classify repetition as a problem.

## Decisions

1. **Four views, shared context.** Add timeline, matrix, repeats and sequence flow to the journey workspace while retaining the existing graph layouts and current default. Preserve compatible selections and filters; an individual journey selection is required for timeline and can be chosen from retained journeys. This avoids forcing an arbitrary journey choice. Alternatives that replace hierarchy or force a new default would change established behavior unnecessarily.
2. **Use captured paths as authority.** Matrix cells use observed direct edges. Repeat and sequence analyses derive from ordered paths within each captured journey, never from aggregate edge counts. Evaluate paths after the existing reference/window projection, preserving gaps and boundary markers; do not join across gaps or boundaries. The return definition in the spec is deliberately observable and excludes self-transitions, which are reported as repeats. An aggregate-only reconstruction would invent order.
3. **Bound and disclose.** Define deterministic, documented step and branch limits for sequence flow, group omitted branches with counts, and label continuation beyond the step limit. Use stable ordering and make drilldown available through the existing bounded ordered-evidence mechanism. Disclose source query truncation, clipping, partial starts and missing retained detail; never substitute daily summaries for paths. Keep baseline and live datasets distinct. This trades an exhaustive visual for an honest, usable display.
4. **Share contract and access.** Extend supported saved journey layout values end-to-end, maintaining old values and references; preserve viewer read-only authorization and existing protected query boundaries. Keep per-journey selection ephemeral unless an existing safe scope or reference represents it. Reuse snapshot data where sufficient; if bounded data prevents accurate aggregates, add a protected bounded projection with explicit coverage rather than silently computing full-population claims from a paged subset. Do not persist raw evidence in references.
5. **Accessible equivalents.** Use labeled lists/tables for matrix, ranking and sequence and an ordered list for timeline, with matching actions, keyboard focus, non-color indicators and English/French copy. Existing evidence and inspectors remain the drilldown destinations.

## Risks / Trade-offs

- Large path sets and high-cardinality matrices → bound work and display, report omitted counts and source truncation, test representative dense and long data.
- Partial starts can undercount returns → label observed-only classification and coverage, avoid claims about earlier visits.
- Per-position flow can be mistaken for unique complete paths → label prefix units and grouped/continued branches; drill down to actual journeys.
- Layout additions can break old references → test existing values, new round trips, validation and rollback compatibility.

## Migration Plan

Extend validation and UI together without rewriting stored references. Existing layouts retain their identifiers and default behavior. On rollback, references saved with new layout identifiers may need an explicit supported-layout fallback; verify existing server behavior and implement a safe documented fallback before release, without deleting references or evidence.
