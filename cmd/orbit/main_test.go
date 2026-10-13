package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestOutputMustNotOverwriteWithoutForce(t *testing.T) {
	p := filepath.Join(t.TempDir(), "report.html")
	if err := os.WriteFile(p, []byte("keep"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := writeOutput(p, []byte("replace"), false); err == nil {
		t.Fatal("overwrote output")
	}
	data, _ := os.ReadFile(p)
	if string(data) != "keep" {
		t.Fatal("content changed")
	}
	if err := writeOutput(p, []byte("replace"), true); err != nil {
		t.Fatal(err)
	}
	var out, errors bytes.Buffer
	if run([]string{"-html", p}, &out, &errors) != 2 {
		t.Fatal("existing output accepted")
	}
	if run([]string{"-help"}, &out, &errors) != 0 {
		t.Fatal("help failed")
	}
}
