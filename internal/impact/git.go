package impact

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const MaxBlob = 1024 * 1024
const MaxSnapshot = 32 * 1024 * 1024
const MaxFiles = 3000

var hashPattern = regexp.MustCompile(`^[0-9a-f]{40,64}$`)

func git(ctx context.Context, repo string, args ...string) ([]byte, error) {
	full := append([]string{"--no-pager", "-C", repo}, args...)
	command := exec.CommandContext(ctx, "git", full...)
	command.Env = append(os.Environ(), "GIT_NO_LAZY_FETCH=1", "GIT_TERMINAL_PROMPT=0")
	var stderr bytes.Buffer
	command.Stderr = &stderr
	output, err := command.Output()
	if err != nil {
		return nil, fmt.Errorf("git %s: %s", args[0], strings.TrimSpace(stderr.String()))
	}
	return output, nil
}

func resolve(ctx context.Context, repo, ref string) (string, error) {
	if ref == "" || strings.HasPrefix(ref, "-") || strings.ContainsAny(ref, "\x00\r\n") {
		return "", fmt.Errorf("invalid commit reference")
	}
	data, err := git(ctx, repo, "rev-parse", "--verify", "--end-of-options", ref+"^{commit}")
	if err != nil {
		return "", err
	}
	sha := strings.TrimSpace(string(data))
	if !hashPattern.MatchString(sha) {
		return "", fmt.Errorf("unexpected commit hash")
	}
	return sha, nil
}

type blob struct {
	name, sha string
	size      int
}

func readSnapshot(ctx context.Context, repo, sha string) (Snapshot, error) {
	data, err := git(ctx, repo, "ls-tree", "-r", "-z", "--long", sha)
	if err != nil {
		return Snapshot{}, err
	}
	snapshot := Snapshot{Files: map[string][]byte{}, Notes: []string{}}
	blobs := []blob{}
	total := 0
	for _, entry := range bytes.Split(data, []byte{0}) {
		parts := bytes.SplitN(entry, []byte{'\t'}, 2)
		if len(parts) != 2 {
			continue
		}
		fields := strings.Fields(string(parts[0]))
		name := string(parts[1])
		if len(fields) != 4 || fields[1] != "blob" || fields[0] == "120000" {
			continue
		}
		if !strings.HasSuffix(name, ".go") && name != "go.mod" && !strings.HasSuffix(name, "/go.mod") {
			continue
		}
		size, parseErr := strconv.Atoi(fields[3])
		if parseErr != nil || size < 0 {
			return Snapshot{}, fmt.Errorf("invalid blob size")
		}
		if size > MaxBlob {
			snapshot.Notes = append(snapshot.Notes, "skipped oversized blob: "+name)
			continue
		}
		if !hashPattern.MatchString(fields[2]) {
			return Snapshot{}, fmt.Errorf("invalid blob hash")
		}
		total += size
		if len(blobs) >= MaxFiles || total > MaxSnapshot {
			return Snapshot{}, fmt.Errorf("snapshot exceeds 3000 source files or 32 MiB")
		}
		blobs = append(blobs, blob{name: name, sha: fields[2], size: size})
	}
	if len(blobs) == 0 {
		return snapshot, nil
	}
	var requests strings.Builder
	for _, blob := range blobs {
		requests.WriteString(blob.sha + "\n")
	}
	command := exec.CommandContext(ctx, "git", "--no-pager", "-C", repo, "cat-file", "--batch")
	command.Env = append(os.Environ(), "GIT_NO_LAZY_FETCH=1", "GIT_TERMINAL_PROMPT=0")
	command.Stdin = strings.NewReader(requests.String())
	var stderr bytes.Buffer
	command.Stderr = &stderr
	stdout, err := command.StdoutPipe()
	if err != nil {
		return Snapshot{}, err
	}
	if err = command.Start(); err != nil {
		return Snapshot{}, err
	}
	reader := bufio.NewReader(stdout)
	readErr := func() error {
		for _, blob := range blobs {
			header, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			fields := strings.Fields(header)
			if len(fields) != 3 || fields[0] != blob.sha || fields[1] != "blob" || fields[2] != strconv.Itoa(blob.size) {
				return fmt.Errorf("unexpected git blob header")
			}
			content := make([]byte, blob.size)
			if _, err = io.ReadFull(reader, content); err != nil {
				return err
			}
			separator, err := reader.ReadByte()
			if err != nil || separator != '\n' {
				return fmt.Errorf("invalid git blob boundary")
			}
			snapshot.Files[blob.name] = content
		}
		return nil
	}()
	if readErr != nil {
		_ = command.Process.Kill()
	}
	waitErr := command.Wait()
	if readErr != nil {
		return Snapshot{}, readErr
	}
	if waitErr != nil {
		return Snapshot{}, fmt.Errorf("git cat-file: %s", strings.TrimSpace(stderr.String()))
	}
	return snapshot, nil
}

func Analyze(ctx context.Context, repo, baseRef, headRef, title string) (Report, error) {
	base, err := resolve(ctx, repo, baseRef)
	if err != nil {
		return Report{}, err
	}
	head, err := resolve(ctx, repo, headRef)
	if err != nil {
		return Report{}, err
	}
	before, err := readSnapshot(ctx, repo, base)
	if err != nil {
		return Report{}, err
	}
	after, err := readSnapshot(ctx, repo, head)
	if err != nil {
		return Report{}, err
	}
	data, err := git(ctx, repo, "diff", "--no-ext-diff", "--no-textconv", "--no-renames", "--name-status", "-z", base, head, "--")
	if err != nil {
		return Report{}, err
	}
	parts := bytes.Split(data, []byte{0})
	changes := []Change{}
	for i := 0; i+1 < len(parts); i += 2 {
		if len(parts[i]) == 0 {
			continue
		}
		changes = append(changes, Change{Status: string(parts[i]), Path: string(parts[i+1])})
	}
	return Build(before, after, changes, base, head, title)
}
