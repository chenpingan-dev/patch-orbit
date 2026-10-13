# Patch Orbit

**Turn a Go change into a review map — with the import chain that explains each package**

[Live map](https://chenpingan-dev.github.io/patch-orbit/) · [中文说明](README.zh-CN.md) · [MIT](LICENSE)

![CI](https://github.com/chenpingan-dev/patch-orbit/actions/workflows/ci.yml/badge.svg)
![Patch Orbit impact map](assets/screenshot.png)

Change a shared model and the API may need review three imports away. Patch Orbit reads two committed snapshots, walks reverse package imports and shows the shortest chain to a changed package. Each edge links back to a source file, import line and base/head snapshot.

## Try it

Go 1.23+ and Git. No model API, database, indexing service or external Go module dependencies.

```sh
go install github.com/chenpingan-dev/patch-orbit/cmd/orbit@v0.1.0
orbit -repo /path/to/your/go-repo -base HEAD~1 -head HEAD -html impact.html -json impact.json
```

Or download a binary from [Releases](https://github.com/chenpingan-dev/patch-orbit/releases). Open `impact.html` in your browser. Specify a root module repository; both commits need a root `go.mod`. Existing outputs need `-force`.

## What you get

- A changed / dependent / unmatched package map, with search and clickable nodes
- Shortest reverse-import chains and exact import evidence from both commits
- Deleted packages and removed imports retained from the base snapshot
- A full test command plus a focused package checklist; generating reports runs **no tests**
- CN/EN switch, report JSON import/export, Markdown checklist and portable offline HTML
- A synthetic demo built by analyzing two actual local fixture commits, not hand-drawn data

## What the map means

This is a **conservative review scope**, not a runtime bug prediction or proof that unmatched packages are safe. Base and head imports are unioned, including `_test.go` and all build-tag variants. Go module manifest changes put all discovered packages in scope.

v0.1 operates on package imports within one root module. It excludes vendor, hidden/underscore directories, testdata and nested modules. Runtime calls, reflection, embedded resources, external modules and uncommitted files are not modeled. Renames appear as delete/add. Symlinks are not dereferenced. Invalid import syntax and unresolved local imports become scan notes. No worktree files are executed, no fetch or Go builds are run by the analyzer.

Limits per snapshot: 3,000 selected files and 32 MiB total; blobs above 1 MiB are skipped with notes. The diagram shows up to 80 packages; the searchable list and JSON retain the full result. The file panel shows the first 30 changes with its displayed/total count. Reports contain module and file paths: inspect before sharing.

## Develop

```sh
git clone https://github.com/chenpingan-dev/patch-orbit.git
cd patch-orbit
go test ./...
go vet ./...
go run ./cmd/demo
```

The builder creates and removes a fresh synthetic repository under ignored `build/`, then emits `build/index.html` and `build/report.json`. Fixture history stays local. GitHub Pages serves the generated offline report.

[Architecture](docs/architecture.md) · [Roadmap](docs/roadmap.md) · [Contributing](CONTRIBUTING.md)

Related work: [GitNexus](https://github.com/abhigyanpatwari/GitNexus) provides broader code intelligence and impact tooling. Patch Orbit focuses on committed Go package imports, a single executable and a report with explicit import evidence. Independent, AI-assisted implementation; no upstream code is included.
