package impact

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCommittedSnapshotsIgnoreWorktreeAndHandleSpaces(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git unavailable")
	}
	root := t.TempDir()
	gitRun := func(args ...string) string {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git: %v %s", err, out)
		}
		return string(out)
	}
	write := func(name, text string) {
		t.Helper()
		p := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(p), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(text), 0600); err != nil {
			t.Fatal(err)
		}
	}
	gitRun("init", "-b", "main")
	gitRun("config", "user.name", "Synthetic test")
	gitRun("config", "user.email", "test@example.invalid")
	gitRun("config", "commit.gpgsign", "false")
	write("go.mod", "module example.test/app\n\ngo 1.23\n")
	write("pkg space/item.go", "package item\nconst Value=1\n")
	write("api/api.go", "package api\nimport _ \"example.test/app/pkg space\"\n")
	gitRun("add", ".")
	gitRun("commit", "-m", "Synthetic base")
	write("pkg space/item.go", "package item\nconst Value=2\n")
	gitRun("add", ".")
	gitRun("commit", "-m", "Synthetic head")
	write("pkg space/item.go", "this uncommitted content must not be parsed")
	before := gitRun("status", "--porcelain")
	r, err := Analyze(context.Background(), root, "HEAD~1", "HEAD", "test")
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Changes) != 1 || r.Changes[0].Path != "pkg space/item.go" || r.Summary.AffectedPackages != 2 {
		t.Fatal(r)
	}
	if after := gitRun("status", "--porcelain"); after != before {
		t.Fatal("worktree changed")
	}
	if _, err := Analyze(context.Background(), root, "--help", "HEAD", "test"); err == nil {
		t.Fatal("option-like ref accepted")
	}
}
