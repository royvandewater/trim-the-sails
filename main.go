package main

import (
	"bytes"
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
	if err := run(dir, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "trim-the-sails:", err)
		os.Exit(1)
	}
}

// run finds every bare repo in dir and prunes its branches concurrently,
// capped at maxParallelism. Each repo's output is buffered and flushed to w in
// directory order so the report stays deterministic regardless of scheduling.
func run(dir string, w io.Writer) error {
	repos, err := findBareRepos(dir)
	if err != nil {
		return err
	}

	outputs := make([]bytes.Buffer, len(repos))
	errs := make([]error, len(repos))

	var wg sync.WaitGroup
	sem := make(chan struct{}, maxParallelism)
	for i, repo := range repos {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			fmt.Fprintln(&outputs[i], "pruning", repo)
			errs[i] = pruneRepo(repo)
		}()
	}
	wg.Wait()

	for i := range repos {
		outputs[i].WriteTo(w)
	}
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
