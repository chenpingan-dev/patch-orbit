package impact

import (
	"reflect"
	"strings"
	"testing"
)

func snapshot(sources map[string]string) Snapshot {
	s := Snapshot{Files: map[string][]byte{"go.mod": []byte("module example.test/app\n\ngo 1.23\n")}}
	for name, source := range sources {
		s.Files[name] = []byte(source)
	}
	return s
}
func node(t *testing.T, r Report, id string) Node {
	t.Helper()
	for _, n := range r.Nodes {
		if n.ID == id {
			return n
		}
	}
	t.Fatalf("missing %s", id)
	return Node{}
}
func TestReverseImportsAndEvidence(t *testing.T) {
	s := snapshot(map[string]string{"model/model.go": "package model", "store/store.go": "package store\nimport _ \"example.test/app/model\"", "api/api.go": "package api\nimport _ \"example.test/app/store\"", "util/util.go": "package util"})
	r, err := Build(s, s, []Change{{Path: "model/model.go", Status: "M"}}, "before", "after", "test")
	if err != nil {
		t.Fatal(err)
	}
	if r.Summary.AffectedPackages != 3 || r.Summary.Transitive != 1 || r.Summary.Unaffected != 1 {
		t.Fatal(r.Summary)
	}
	n := node(t, r, "api")
	if !reflect.DeepEqual(n.Chain, []string{"api", "store", "model"}) {
		t.Fatal(n.Chain)
	}
	if r.Edges[0].Evidence[0].Line != 2 || len(r.Edges[0].Evidence) != 2 {
		t.Fatal(r.Edges)
	}
}
func TestDeletedPackageAndRemovedImportsRemainEvidence(t *testing.T) {
	before := snapshot(map[string]string{"old/old.go": "package old", "api/api.go": "package api\nimport _ \"example.test/app/old\""})
	after := snapshot(map[string]string{"api/api.go": "package api"})
	r, err := Build(before, after, []Change{{Path: "old/old.go", Status: "D"}}, "b", "h", "test")
	if err != nil {
		t.Fatal(err)
	}
	if node(t, r, "old").Present || *node(t, r, "api").Distance != 1 || r.Edges[0].Evidence[0].Snapshot != "base" {
		t.Fatal(r)
	}
}
func TestConservativeScopeAndExclusions(t *testing.T) {
	s := snapshot(map[string]string{"a/a.go": "package a", "b/b.go": "package b", "b/b_test.go": "package b\nimport _ \"example.test/app/a\"", "c/c_windows.go": "//go:build windows\npackage c\nimport _ \"example.test/app/a\"", "vendor/x/x.go": "package x", "testdata/x.go": "package x", "_ignored/x.go": "package x", "nested/go.mod": "module other.test/app", "nested/n.go": "package nested", "bad/bad.go": "not go"})
	r, err := Build(s, s, []Change{{Path: "a/a.go", Status: "M"}}, "b", "h", "test")
	if err != nil {
		t.Fatal(err)
	}
	if r.Summary.Packages != 3 || r.Summary.AffectedPackages != 3 {
		t.Fatal(r.Summary)
	}
	found := false
	for _, e := range r.Edges {
		for _, proof := range e.Evidence {
			if proof.Test {
				found = true
			}
		}
	}
	if !found {
		t.Fatal("test evidence missing")
	}
	if !strings.Contains(strings.Join(r.Notes, "\n"), "nested Go modules excluded") {
		t.Fatal(r.Notes)
	}
}
func TestModuleChangeAndImportCycle(t *testing.T) {
	s := snapshot(map[string]string{"a/a.go": "package a\nimport _ \"example.test/app/b\"", "b/b.go": "package b\nimport _ \"example.test/app/a\"", "c/c.go": "package c"})
	r, err := Build(s, s, []Change{{Path: "a/a.go", Status: "M"}}, "b", "h", "test")
	if err != nil {
		t.Fatal(err)
	}
	if r.Summary.AffectedPackages != 2 {
		t.Fatal(r.Summary)
	}
	r, err = Build(s, s, []Change{{Path: "go.sum", Status: "M"}}, "b", "h", "test")
	if err != nil {
		t.Fatal(err)
	}
	if !r.ModuleChanged || r.Summary.AffectedPackages != 3 {
		t.Fatal(r.Summary)
	}
}
func TestMissingModuleIsExplicitError(t *testing.T) {
	_, err := Build(Snapshot{}, Snapshot{}, nil, "b", "h", "test")
	if err == nil {
		t.Fatal("missing module accepted")
	}
}

func TestEvidenceRetainsSnapshotSpecificLinesAndFiles(t *testing.T) {
	before := snapshot(map[string]string{
		"old/package.go":         "package old",
		"new/package.go":         "package new",
		"stable/package.go":      "package stable",
		"api/package.go":         "package api\nimport _ \"example.test/app/service\"",
		"service/package.go":     "package service\n\nimport (\n _ \"example.test/app/old\"\n _ \"example.test/app/stable\"\n)\n",
		"service/legacy_test.go": "package service\nimport _ \"example.test/app/old\"\n",
	})
	after := snapshot(map[string]string{
		"old/package.go":              "package old",
		"new/package.go":              "package new",
		"stable/package.go":           "package stable",
		"api/package.go":              "package api\nimport _ \"example.test/app/service\"",
		"service/package.go":          "package service\n\n// Integration changed\n\nimport (\n _ \"example.test/app/stable\"\n _ \"example.test/app/new\"\n)\n",
		"service/replacement_test.go": "package service\nimport _ \"example.test/app/new\"\n",
	})
	report, err := Build(before, after, []Change{
		{Path: "service/package.go", Status: "M"},
		{Path: "service/legacy_test.go", Status: "D"},
		{Path: "service/replacement_test.go", Status: "A"},
	}, "before", "after", "test")
	if err != nil {
		t.Fatal(err)
	}
	want := map[string][]Evidence{
		"old": {
			{File: "service/legacy_test.go", Line: 2, Import: "example.test/app/old", Snapshot: "base", Test: true},
			{File: "service/package.go", Line: 4, Import: "example.test/app/old", Snapshot: "base"},
		},
		"new": {
			{File: "service/package.go", Line: 7, Import: "example.test/app/new", Snapshot: "head"},
			{File: "service/replacement_test.go", Line: 2, Import: "example.test/app/new", Snapshot: "head", Test: true},
		},
		"stable": {
			{File: "service/package.go", Line: 5, Import: "example.test/app/stable", Snapshot: "base"},
			{File: "service/package.go", Line: 6, Import: "example.test/app/stable", Snapshot: "head"},
		},
	}
	for _, edge := range report.Edges {
		if edge.From != "service" {
			continue
		}
		if !reflect.DeepEqual(edge.Evidence, want[edge.To]) {
			t.Errorf("%s evidence: got %#v, want %#v", edge.To, edge.Evidence, want[edge.To])
		}
		delete(want, edge.To)
	}
	if len(want) != 0 {
		t.Fatalf("missing edges: %v", want)
	}
	if !reflect.DeepEqual(node(t, report, "api").Chain, []string{"api", "service"}) || report.Summary.AffectedPackages != 2 {
		t.Fatal("import comparison changed the conservative review scope", report)
	}
}

func TestEvidenceUsesPhysicalSourceLines(t *testing.T) {
	sources := snapshot(map[string]string{
		"model/package.go": "package model",
		"api/package.go":   "package api\nimport (\n//line generated.go:900\n _ \"example.test/app/model\"\n)\n",
	})
	report, err := Build(sources, sources, []Change{{Path: "model/package.go", Status: "M"}}, "before", "after", "test")
	if err != nil {
		t.Fatal(err)
	}
	for _, proof := range report.Edges[0].Evidence {
		if proof.File != "api/package.go" || proof.Line != 4 {
			t.Errorf("import location must use the physical repository file: %#v", proof)
		}
	}
}
