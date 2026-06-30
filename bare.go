package main

import (
	"os"
	"path/filepath"
	"strings"
)

// findBareRepos returns the paths of directories within dir whose names end
// in ".git" — the convention for bare git repositories.
func findBareRepos(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var repos []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".git") {
			continue
		}
		repos = append(repos, filepath.Join(dir, entry.Name()))
	}
	return repos, nil
}
