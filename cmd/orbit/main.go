package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/chenpingan-dev/patch-orbit/internal/impact"
	"github.com/chenpingan-dev/patch-orbit/internal/report"
)

func run(args []string, out, stderr io.Writer) int {
	flags := flag.NewFlagSet("patch-orbit", flag.ContinueOnError)
	flags.SetOutput(stderr)
	repo := flags.String("repo", ".", "explicit Git repository with a root go.mod")
	base := flags.String("base", "HEAD~1", "base commit")
	head := flags.String("head", "HEAD", "head commit")
	title := flags.String("title", "Change impact", "report title")
	htmlPath := flags.String("html", "orbit.html", "standalone HTML output")
	jsonPath := flags.String("json", "", "optional JSON report output")
	force := flags.Bool("force", false, "allow replacing existing output files")
	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}
	if flags.NArg() != 0 {
		fmt.Fprintln(stderr, "patch-orbit: unexpected positional arguments")
		return 2
	}
	paths := []string{*htmlPath}
	if *jsonPath != "" {
		paths = append(paths, *jsonPath)
	}
	seen := map[string]bool{}
	for _, name := range paths {
		absolute, err := filepath.Abs(name)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
		if runtime.GOOS == "windows" {
			absolute = strings.ToLower(absolute)
		}
		if seen[absolute] {
			fmt.Fprintln(stderr, "patch-orbit: outputs must use different paths")
			return 2
		}
		seen[absolute] = true
		if !*force {
			if _, err := os.Stat(name); err == nil {
				fmt.Fprintln(stderr, "patch-orbit: output exists; choose a new path or use -force")
				return 2
			}
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	value, err := impact.Analyze(ctx, *repo, *base, *head, *title)
	if err != nil {
		fmt.Fprintln(stderr, "patch-orbit:", err)
		return 2
	}
	html, err := report.HTML(value)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	if *jsonPath != "" {
		data, err := json.MarshalIndent(value, "", "  ")
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
		if err = writeOutput(*jsonPath, append(data, '\n'), *force); err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
	}
	if err = writeOutput(*htmlPath, html, *force); err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	fmt.Fprintf(out, "Report: %s\n%d changed packages; %d packages to review\n", *htmlPath, value.Summary.ChangedPackages, value.Summary.AffectedPackages)
	return 0
}

func writeOutput(name string, data []byte, force bool) error {
	mode := os.O_WRONLY | os.O_CREATE | os.O_EXCL
	if force {
		mode = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	}
	file, err := os.OpenFile(name, mode, 0600)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	closeErr := file.Close()
	if err != nil {
		return err
	}
	return closeErr
}

func main() { os.Exit(run(os.Args[1:], os.Stdout, os.Stderr)) }
