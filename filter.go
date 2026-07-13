package main

import (
	"path/filepath"
	"strings"
)

// filterRepos narrows repos to those whose basename matches one of names, with
// or without the ".git" suffix. When names is empty, every repo is returned
// unchanged.
func filterRepos(repos []string, names []string) []string {
	if len(names) == 0 {
		return repos
	}

	wanted := make(map[string]bool, len(names))
	for _, name := range names {
		wanted[strings.TrimSuffix(name, ".git")] = true
	}

	var filtered []string
	for _, repo := range repos {
		base := strings.TrimSuffix(filepath.Base(repo), ".git")
		if wanted[base] {
			filtered = append(filtered, repo)
		}
	}
	return filtered
}
