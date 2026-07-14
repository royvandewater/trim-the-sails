package main

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// maxParallelism caps how many repos are pruned concurrently. Pruning is
// dominated by network-bound git fetches, so a handful at once is plenty.
const maxParallelism = 8

func main() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "trim-the-sails:", err)
		os.Exit(1)
	}
	if err := run(dir, os.Args[1:], os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "trim-the-sails:", err)
		os.Exit(1)
	}
}

// run finds the bare repos in dir and prunes their branches concurrently,
// capped at maxParallelism. When names is non-empty, only repos matching those
// names are pruned; otherwise every bare repo is. Progress is reported to w as
// a live bar with an ETA when w is a terminal, and stays silent otherwise so
// piped output remains clean.
func run(dir string, names []string, w io.Writer) error {
	for _, name := range names {
		if name == "--help" || name == "-h" {
			writeUsage(w)
			return nil
		}
	}

	repos, err := findBareRepos(dir)
	if err != nil {
		return err
	}
	repos = filterRepos(repos, names)

	errs := make([]error, len(repos))
	bar := newProgressBar(w, len(repos))

	var wg sync.WaitGroup
	sem := make(chan struct{}, maxParallelism)
	for i, repo := range repos {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			errs[i] = pruneRepo(repo)
			bar.advance()
		}()
	}
	wg.Wait()
	bar.finish()

	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

// writeUsage prints how to invoke the tool, including the optional repo-name
// arguments that narrow which bare repos get pruned.
func writeUsage(w io.Writer) {
	fmt.Fprint(w, `Usage: trim-the-sails [repo...]

Prunes branches of the bare git repos (directories ending in .git) in the
current directory. With no arguments, every bare repo is pruned. Given one or
more repo names (with or without the .git suffix), only those repos are pruned.

Flags:
  -h, --help   Show this help and exit.
`)
}
