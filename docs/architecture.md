# Architecture

Explicit repository + two refs → verified commit IDs → bounded Git blobs → Go AST imports → union of base/head package edges → multi-source reverse BFS → evidence-bearing report

Git commands use argument arrays with external diff and textconv disabled. Blob reading uses `cat-file --batch`, validates object headers and retains NUL-delimited filenames. Snapshots are read from Git objects, not the working tree. File parsing uses Go's standard parser in imports-only mode.

Changed Go sources seed the reverse traversal. Stable ordering makes shortest-chain selection deterministic when paths tie. Cycle handling uses visited nodes. Manifest changes seed every discovered package. Every edge retains evidence from base, head or both. A package can survive in the report after deletion, flagged as absent in head.

Embedded web assets render into one HTML file, with JSON script delimiters escaped. Imported JSON is shape-checked and displayed text is escaped. The browser calculates no new dependency edges. Suggested test commands are data for the user to run, never executed by the report.

The source inspector groups evidence by repository-relative file path, import path and test-file status. Each group retains sorted, distinct one-based physical lines separately for `base` and `head`; Go `//line` directives do not change repository locations. Line shifts therefore stay in the same group, while renamed files remain separate. Edge status unions all its evidence snapshots; file-source filters may show a base-only file under an edge present in both snapshots. Filters never alter edges, shortest paths, counts or test targets. Both changed and unchanged packages expose their outgoing imports, with the complete impact path available separately. All locations are accessible, without the former three-location truncation.

The JSON schema remains `patch-orbit/v1`, with no new required fields. The viewer accepts generated older reports and rejects unknown snapshot names, nonpositive source lines and invalid test flags. Presence labels describe parsed evidence only; skipped or invalid files can prevent a source from appearing in a snapshot.

Tests exercise reverse chains, deletions, tests/build variants, exclusions, cycles, manifests, import evidence, real Git snapshots with path spaces, untouched worktrees, script escaping and output protection.
