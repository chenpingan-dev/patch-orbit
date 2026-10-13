# Architecture

Explicit repository + two refs → verified commit IDs → bounded Git blobs → Go AST imports → union of base/head package edges → multi-source reverse BFS → evidence-bearing report

Git commands use argument arrays with external diff and textconv disabled. Blob reading uses `cat-file --batch`, validates object headers and retains NUL-delimited filenames. Snapshots are read from Git objects, not the working tree. File parsing uses Go's standard parser in imports-only mode.

Changed Go sources seed the reverse traversal. Stable ordering makes shortest-chain selection deterministic when paths tie. Cycle handling uses visited nodes. Manifest changes seed every discovered package. Every edge retains evidence from base, head or both. A package can survive in the report after deletion, flagged as absent in head.

Embedded web assets render into one HTML file, with JSON script delimiters escaped. Imported JSON is shape-checked and displayed text is escaped. The browser calculates no new dependency edges. Suggested test commands are data for the user to run, never executed by the report.

Tests exercise reverse chains, deletions, tests/build variants, exclusions, cycles, manifests, import evidence, real Git snapshots with path spaces, untouched worktrees, script escaping and output protection.
