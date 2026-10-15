// Build a report from two real commits in an isolated synthetic fixture repository.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chenpingan-dev/patch-orbit/internal/impact"
	"github.com/chenpingan-dev/patch-orbit/internal/report"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func build() error {
	// A fresh temporary repository prevents accidental changes to an existing checkout.
	parent, err := filepath.Abs("build")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(parent, 0755); err != nil {
		return err
	}
	root, err := os.MkdirTemp(parent, "synthetic-")
	if err != nil {
		return err
	}
	if filepath.Dir(root) != parent || !strings.HasPrefix(filepath.Base(root), "synthetic-") {
		return fmt.Errorf("invalid fixture directory")
	}
	defer os.RemoveAll(root)
	run := func(args ...string) error {
		cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
		data, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("git: %w: %s", err, data)
		}
		return nil
	}
	write := func(name, text string) error {
		p := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(p), 0700); err != nil {
			return err
		}
		return os.WriteFile(p, []byte(text), 0600)
	}
	for _, args := range [][]string{{"init", "-b", "main"}, {"config", "user.name", "Synthetic demo"}, {"config", "user.email", "demo@example.invalid"}, {"config", "commit.gpgsign", "false"}} {
		if err := run(args...); err != nil {
			return err
		}
	}
	if err := write("go.mod", "module example.local/notebox\n\ngo 1.23\n"); err != nil {
		return err
	}
	graph := map[string][]string{
		"internal/model": {}, "internal/store": {"internal/model"}, "internal/export": {"internal/model"},
		"internal/search": {"internal/store"}, "internal/service": {"internal/store", "internal/export"},
		"internal/api": {"internal/service", "internal/search"}, "cmd/notebox": {"internal/api"},
		"cmd/export": {"internal/export"}, "internal/metrics": {}, "internal/config": {},
		"internal/health": {"internal/metrics"}, "cmd/check": {"internal/health", "internal/config"},
	}
	for name, imports := range graph {
		var text strings.Builder
		text.WriteString("package " + filepath.Base(name) + "\n\n")
		for _, imp := range imports {
			text.WriteString("import _ \"example.local/notebox/" + imp + "\"\n")
		}
		text.WriteString("\nconst Version = 1\n")
		if err := write(name+"/package.go", text.String()); err != nil {
			return err
		}
	}
	if err := run("add", "."); err != nil {
		return err
	}
	if err := run("commit", "-m", "Synthetic fixture base"); err != nil {
		return err
	}
	if err := write("internal/model/package.go", "package model\n\nconst Version = 2\n"); err != nil {
		return err
	}
	// Keep one import at a new line, remove one, and add one in the same package.
	service := "package service\n\n// Synthetic search integration\nimport (\n\t_ \"example.local/notebox/internal/export\"\n\t_ \"example.local/notebox/internal/search\"\n)\n\nconst Version = 2\n"
	if err := write("internal/service/package.go", service); err != nil {
		return err
	}
	if err := write("README.md", "Synthetic notebook fixture used only to build the demo\n"); err != nil {
		return err
	}
	if err := run("add", "."); err != nil {
		return err
	}
	if err := run("commit", "-m", "Synthetic fixture search integration"); err != nil {
		return err
	}
	value, err := impact.Analyze(context.Background(), root, "HEAD~1", "HEAD", "Notebook search integration")
	if err != nil {
		return err
	}
	value.Synthetic = true
	value.TitleI18n = map[string]string{"zh": "笔记检索接入", "en": "Notebook search integration"}
	html, err := report.HTML(value)
	if err != nil {
		return err
	}
	if err := os.MkdirAll("build", 0755); err != nil {
		return err
	}
	if err := os.WriteFile("build/index.html", html, 0644); err != nil {
		return err
	}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile("build/report.json", append(data, '\n'), 0644); err != nil {
		return err
	}
	fmt.Printf("Built synthetic report: %d changed packages, %d packages to review\n", value.Summary.ChangedPackages, value.Summary.AffectedPackages)
	return nil
}
func main() {
	if err := build(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
