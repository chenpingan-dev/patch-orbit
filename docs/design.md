# Report design

Patch Orbit is a spatial review tool. Its primary interaction is following imports from a changed package to its dependents.

- Use a white canvas, soft silver surroundings, blue selection and a restrained amber marker for changed packages
- Open directly on the report and map; avoid a large marketing hero inside the workspace
- Keep package navigation on the left and the map dominant; open source evidence only after selecting a package
- Separate the map, changed files and test checklist into focused views
- Use fine import arrows and readable package cards; retain visible status labels instead of relying only on color
- Microsoft YaHei comes first in the font stack; use short Chinese and English action labels
- At narrow widths, keep scrolling within the graph and place the inspector and package list below it
- Preserve keyboard node selection, Escape to dismiss evidence and focus restoration
- A selected package shows its direct imports, including changed packages; keep the full dependent-to-change path in an expandable section
- Compare file sources with explicit base-only, head-only and both-snapshot labels, using 14px source text and visible before/after line positions
- Filter file evidence without changing the map or checklist; distinguish an edge's combined status from each file's status and explain scan limitations

The interface must remain different from Agent Rewind's recording player. Do not copy its filmstrip or return to a dark green dashboard with multiple boxed metric panels.

Before publishing UI changes, inspect both languages at desktop and phone sizes. Check selection, search, scope, source evidence, view switching, import/export, empty reports and the 80-node display limit. Keep analysis limitations accessible and synthetic fixtures labeled. The graph represents conservative package import scope, not a prediction of runtime failures.
