# Preserved visual theme and icon

The owner explicitly requested retaining the first MVP's visual theme and icon on 2026-10-28 after correcting the implementation base to `docs/journey-capability-plan-2026-10-27`. This settles the visual direction for the replacement implementation; the earlier exploratory purple direction is superseded.

Keep the forest-green masthead/brand, warm paper background, quiet outlined cards, serif headings, readable system-sans body, restrained status colors with text/icons, generous spacing, and visible keyboard focus. The linked-thread Sentrame mark remains the brand icon. Extend this language to the selected Svelte interface and historical/live maps rather than recreating the old fixed-outcome screen as the new product scope.

| Token | Preserved value |
| --- | --- |
| Paper | `#f5f4ee` |
| Main ink | `#253b36` |
| Muted ink | `#626e66` |
| Lines | `#d9dfd5` |
| Brand green | `#153f39` |
| Card surface | `#fffefa` |
| Completed | `#1e6656` |
| Interrupted | `#8c642c` |
| Failed | `#a3493d` |
| Unknown | `#606d7c` |
| Keyboard focus | `#bd782a` |

- `theme.css`: original source stylesheet, preserved byte-for-byte as the migration reference. Extract shared tokens/components into the new interface rather than loading obsolete screen-specific classes blindly.
- `sentrame-icon.svg`: original icon geometry using `currentColor`, with an SVG namespace for standalone reuse. In the masthead it inherits the light brand foreground; elsewhere choose the appropriate contrasting color.
- `original-screen.html`: original markup reference, including the inline brand mark. Its old API/script references are historical and are not a runnable new implementation.
- `review.jpg`: captured first MVP review screen. This is appearance evidence, not proof of the evolved capability plan.

Preserve readable contrast, keyboard disclosure and focus through live updates, and equivalent list/table access. The prior browser test at 390×844 is useful baseline evidence but does not qualify the replacement Svelte maps or real devices.

The owner-approved canvas revision preserves these tokens and the mark. Hierarchy uses quiet rounded stage cards, force uses outlined circles, and Sankey uses traffic-volume bars/ribbons. Small independently sized arrows, selected-neighborhood emphasis and contextual timing replace the earlier duration-sized graph nodes. See the active journey-exploration specification for current interaction requirements.
