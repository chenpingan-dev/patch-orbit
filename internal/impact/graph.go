package impact

import (
	"bufio"
	"fmt"
	"go/parser"
	"go/token"
	"path"
	"sort"
	"strconv"
	"strings"
)

type Change struct {
	Path     string `json:"path"`
	Status   string `json:"status"`
	Package  string `json:"package,omitempty"`
	GoSource bool   `json:"go_source"`
}

// Evidence locates one parsed import in a committed repository snapshot.
type Evidence struct {
	File     string `json:"file"`
	Line     int    `json:"line"` // One-based physical line; //line directives do not remap it.
	Import   string `json:"import"`
	Snapshot string `json:"snapshot"` // "base" or "head", independent of whether the file changed.
	Test     bool   `json:"test"`
}

type Edge struct {
	From     string     `json:"from"`
	To       string     `json:"to"`
	Evidence []Evidence `json:"evidence"`
}

type Node struct {
	ID       string   `json:"id"`
	Files    []string `json:"files"`
	Present  bool     `json:"present"`
	Changed  bool     `json:"changed"`
	Distance *int     `json:"distance"`
	Chain    []string `json:"chain"`
}

type Summary struct {
	Packages         int `json:"packages"`
	ChangedFiles     int `json:"changed_files"`
	ChangedPackages  int `json:"changed_packages"`
	AffectedPackages int `json:"affected_packages"`
	Direct           int `json:"direct"`
	Transitive       int `json:"transitive"`
	Unaffected       int `json:"unaffected"`
}

type Report struct {
	Schema        string            `json:"schema"`
	Title         string            `json:"title"`
	TitleI18n     map[string]string `json:"title_i18n,omitempty"`
	Module        string            `json:"module"`
	Base          string            `json:"base"`
	Head          string            `json:"head"`
	Synthetic     bool              `json:"synthetic"`
	ModuleChanged bool              `json:"module_changed"`
	Nodes         []Node            `json:"nodes"`
	Edges         []Edge            `json:"edges"`
	Changes       []Change          `json:"changes"`
	Summary       Summary           `json:"summary"`
	Notes         []string          `json:"notes"`
	TestCommand   string            `json:"test_command"`
}

type Snapshot struct {
	Files map[string][]byte
	Notes []string
}

type snapshotGraph struct {
	module string
	nodes  map[string][]string
	edges  map[string]*Edge
	notes  []string
}

func modulePath(data []byte) (string, error) {
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(strings.SplitN(scanner.Text(), "//", 2)[0])
		parts := strings.Fields(line)
		if len(parts) == 2 && parts[0] == "module" {
			value := parts[1]
			if strings.HasPrefix(value, "\"") {
				var err error
				value, err = strconv.Unquote(value)
				if err != nil {
					return "", fmt.Errorf("invalid module path: %w", err)
				}
			}
			if value != "" {
				return value, nil
			}
		}
	}
	return "", fmt.Errorf("snapshot needs a root go.mod with a module directive")
}

func directory(filename string) string { return path.Dir(filename) }

func graph(snapshot Snapshot, version string) (snapshotGraph, error) {
	module, err := modulePath(snapshot.Files["go.mod"])
	if err != nil {
		return snapshotGraph{}, err
	}
	g := snapshotGraph{module: module, nodes: map[string][]string{}, edges: map[string]*Edge{}, notes: append([]string{}, snapshot.Notes...)}
	nested := []string{}
	for name := range snapshot.Files {
		if name != "go.mod" && path.Base(name) == "go.mod" {
			nested = append(nested, path.Dir(name)+"/")
		}
	}
	parsed := map[string][]Evidence{}
	for name, data := range snapshot.Files {
		if !strings.HasSuffix(name, ".go") || strings.HasPrefix(name, "vendor/") {
			continue
		}
		skip := false
		for _, prefix := range nested {
			if strings.HasPrefix(name, prefix) {
				skip = true
				break
			}
		}
		for _, part := range strings.Split(name, "/") {
			if part == "testdata" || strings.HasPrefix(part, ".") || strings.HasPrefix(part, "_") {
				skip = true
			}
		}
		if skip {
			continue
		}
		set := token.NewFileSet()
		file, parseErr := parser.ParseFile(set, name, data, parser.ImportsOnly)
		if parseErr != nil {
			g.notes = append(g.notes, version+": skipped invalid Go imports in "+name)
			continue
		}
		dir := directory(name)
		g.nodes[dir] = append(g.nodes[dir], name)
		for _, imp := range file.Imports {
			value, quoteErr := strconv.Unquote(imp.Path.Value)
			if quoteErr != nil {
				continue
			}
			parsed[dir] = append(parsed[dir], Evidence{File: name, Line: set.PositionFor(imp.Pos(), false).Line,
				Import: value, Snapshot: version, Test: strings.HasSuffix(name, "_test.go")})
		}
	}
	for id, names := range g.nodes {
		sort.Strings(names)
		g.nodes[id] = names
	}
	for from, imports := range parsed {
		for _, evidence := range imports {
			to := ""
			if evidence.Import == module {
				to = "."
			} else if strings.HasPrefix(evidence.Import, module+"/") {
				to = strings.TrimPrefix(evidence.Import, module+"/")
			}
			if to == "" || to == from {
				continue
			}
			if _, exists := g.nodes[to]; !exists {
				g.notes = append(g.notes, version+": unresolved local import "+evidence.Import+" in "+evidence.File)
				continue
			}
			key := from + "\x00" + to
			if g.edges[key] == nil {
				g.edges[key] = &Edge{From: from, To: to, Evidence: []Evidence{}}
			}
			g.edges[key].Evidence = append(g.edges[key].Evidence, evidence)
		}
	}
	if len(nested) > 0 {
		g.notes = append(g.notes, version+": nested Go modules excluded")
	}
	sort.Strings(g.notes)
	return g, nil
}

// Build returns a conservative package-import review scope, not a runtime failure prediction.
func Build(before, after Snapshot, changes []Change, base, head, title string) (Report, error) {
	old, err := graph(before, "base")
	if err != nil {
		return Report{}, fmt.Errorf("base: %w", err)
	}
	now, err := graph(after, "head")
	if err != nil {
		return Report{}, fmt.Errorf("head: %w", err)
	}
	report := Report{Schema: "patch-orbit/v1", Title: title, Module: now.module, Base: base, Head: head,
		Nodes: []Node{}, Edges: []Edge{}, Changes: append([]Change{}, changes...), Notes: []string{}, TestCommand: "go test ./..."}
	ids := map[string]bool{}
	for id := range old.nodes {
		ids[id] = true
	}
	for id := range now.nodes {
		ids[id] = true
	}
	ordered := []string{}
	for id := range ids {
		ordered = append(ordered, id)
	}
	sort.Strings(ordered)
	changed := map[string]bool{}
	for i, change := range report.Changes {
		if change.Path == "go.mod" || change.Path == "go.sum" {
			report.ModuleChanged = true
		}
		if !strings.HasSuffix(change.Path, ".go") {
			continue
		}
		id := directory(change.Path)
		if ids[id] {
			changed[id] = true
			report.Changes[i].Package = id
			report.Changes[i].GoSource = true
		}
	}
	edges := map[string]*Edge{}
	for _, source := range []snapshotGraph{old, now} {
		for key, edge := range source.edges {
			if edges[key] == nil {
				edges[key] = &Edge{From: edge.From, To: edge.To, Evidence: []Evidence{}}
			}
			edges[key].Evidence = append(edges[key].Evidence, edge.Evidence...)
		}
	}
	keys := []string{}
	for key := range edges {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	reverse := map[string][]string{}
	for _, key := range keys {
		edge := edges[key]
		sort.Slice(edge.Evidence, func(i, j int) bool {
			a, b := edge.Evidence[i], edge.Evidence[j]
			if a.File != b.File {
				return a.File < b.File
			}
			if a.Line != b.Line {
				return a.Line < b.Line
			}
			return a.Snapshot < b.Snapshot
		})
		report.Edges = append(report.Edges, *edge)
		reverse[edge.To] = append(reverse[edge.To], edge.From)
	}
	distance := map[string]int{}
	chains := map[string][]string{}
	queue := []string{}
	for _, id := range ordered {
		if changed[id] || report.ModuleChanged {
			distance[id] = 0
			chains[id] = []string{id}
			queue = append(queue, id)
		}
	}
	for index := 0; index < len(queue); index++ {
		dependency := queue[index]
		for _, importer := range reverse[dependency] {
			if _, visited := distance[importer]; visited {
				continue
			}
			distance[importer] = distance[dependency] + 1
			chains[importer] = append([]string{importer}, chains[dependency]...)
			queue = append(queue, importer)
		}
	}
	for _, id := range ordered {
		files, present := now.nodes[id]
		if !present {
			files = old.nodes[id]
		}
		node := Node{ID: id, Files: files, Present: present, Changed: changed[id], Chain: []string{}}
		if value, exists := distance[id]; exists {
			node.Distance = &value
			node.Chain = chains[id]
			report.Summary.AffectedPackages++
			if !node.Changed && value == 1 {
				report.Summary.Direct++
			}
			if !node.Changed && value > 1 {
				report.Summary.Transitive++
			}
		} else {
			report.Summary.Unaffected++
		}
		if node.Changed {
			report.Summary.ChangedPackages++
		}
		report.Nodes = append(report.Nodes, node)
	}
	report.Summary.Packages = len(report.Nodes)
	report.Summary.ChangedFiles = len(changes)
	report.Notes = append(report.Notes, old.notes...)
	report.Notes = append(report.Notes, now.notes...)
	if report.ModuleChanged {
		report.Notes = append(report.Notes, "go.mod/go.sum changed: review every package")
	}
	if old.module != now.module {
		report.Notes = append(report.Notes, "module path changed between snapshots")
	}
	report.Notes = append(report.Notes, "package imports only; includes tests and all build-tag variants; excludes runtime calls, embed assets and external modules")
	return report, nil
}
