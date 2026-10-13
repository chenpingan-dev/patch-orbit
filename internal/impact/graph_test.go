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
