# Tasks

## 1. Shared evidence and references

- [x] 1.1 Implement bounded projections for direct transitions, repeats, returns and ordered prefixes over captured paths; verify counts, directed edges, gaps, partial starts, boundary clipping and equal-edge/different-order cases with focused tests.
- [x] 1.2 Carry source truncation, retention coverage and separate baseline/live counts into projections; verify partial and summary-only selections never claim complete ordered data.
- [x] 1.3 Extend saved layout validation and restoration for new views while preserving existing references, authorization and safe filters; test round trips, old layouts, viewer access and a documented rollback fallback.

## 2. Timeline and matrix

- [x] 2.1 Add the journey timeline and journey selection with repeated occurrences, diagnostic drilldown, gaps and unavailable states; verify chronological, partial and expired examples and keyboard access.
- [x] 2.2 Add the directed transition matrix with labeled counts and evidence drilldown; verify diagonal, reverse directions, empty versus unavailable cells and keyboard/table access.
- [x] 2.3 Document timeline and matrix use and English/French labels; verify both languages expose equivalent actions and units.

## 3. Repeat analysis and sequence flow

- [x] 3.1 Add ranked repeat and return analysis with explicit denominators and journey drilldown; test A-B-A-B, self-repeats, partial starts and non-judgmental labels.
- [x] 3.2 Add deterministic bounded prefix flow, labeled grouped branches and continuations with evidence drilldown; test long paths, branch limits, different orderings and journey-count widths.
- [x] 3.3 Document repeat definitions, sequence limits and coverage in English and French; verify accessible lists or tables offer equivalent counts and actions.

## 4. Integrated verification

- [x] 4.1 Verify switching among all seven views preserves applicable filters and selections, existing graphs and references continue to work, and live/baseline counts remain separate.
- [x] 4.2 Run relevant backend, frontend and browser checks, including keyboard, localization, bounded dense data and unavailable detail; record passed, failed and unrun checks.
