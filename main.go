package main

import (
	"fmt"
	"io"
	"os"
)

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

// run finds every bare repo in dir and prunes its branches, reporting each
// one to w.
func run(dir string, w io.Writer) error {
	repos, err := findBareRepos(dir)
	if err != nil {
		return err
	}

	for _, repo := range repos {
		fmt.Fprintln(w, "pruning", repo)
		if err := pruneRepo(repo); err != nil {
			return err
		}
	}
	return nil
}
