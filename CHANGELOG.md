# Changelog

## 0.1.2 — 2026-10-15

- Compare a selected package's direct imports, including packages with changed source files
- Filter file evidence by base only, head only or both snapshots without changing the map's conservative review scope
- Pair the same source file and import path across line moves, retain both source positions and distinguish test imports
- Keep complete source evidence accessible instead of silently showing only three locations per path link
- Fix source lines remapped by Go `//line` directives; report the physical repository file location
- Refresh the labeled synthetic report, bilingual guide and import comparison screenshots

## 0.1.1 — 2026-10-13

- Rebuilt the report around a light dependency canvas and an on-demand source inspector
- Added separate map, changed-file and verification views, with keyboard dismissal and focus restoration
- Refreshed package cards, import arrows, bilingual documentation and mobile layout

## 0.1.0

Initial working CLI, standalone offline report, bilingual interactive demo and automated validation
