package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindBareReposReturnsDirsEndingInDotGit(t *testing.T) {
	dir := t.TempDir()
	mustMkdir(t, filepath.Join(dir, "alpha.git"))
	mustMkdir(t, filepath.Join(dir, "beta.git"))
	mustMkdir(t, filepath.Join(dir, "notarepo"))
	mustWriteFile(t, filepath.Join(dir, "gamma.git")) // a file, not a dir

	got, err := findBareRepos(dir)
	if err != nil {
		t.Fatalf("findBareRepos: %v", err)
	}

	want := []string{
		filepath.Join(dir, "alpha.git"),
		filepath.Join(dir, "beta.git"),
	}
	if !equalStrings(got, want) {
		t.Fatalf("findBareRepos() = %v, want %v", got, want)
	}
}

func mustMkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.Mkdir(path, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
}

func mustWriteFile(t *testing.T, path string) {
	t.Helper()
	if err := os.WriteFile(path, nil, 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
