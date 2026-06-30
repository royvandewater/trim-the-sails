package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestRunPrunesEachBareRepoAndReports(t *testing.T) {
	root := t.TempDir()
	origin := filepath.Join(root, "origin")
	bare := filepath.Join(root, "repo.git")

	git(t, root, "init", "-q", origin)
	git(t, origin, "config", "user.email", "test@example.com")
	git(t, origin, "config", "user.name", "Test")
	git(t, origin, "commit", "-q", "--allow-empty", "-m", "init")
	git(t, origin, "branch", "doomed")

	git(t, root, "init", "-q", "--bare", bare)
	git(t, bare, "remote", "add", "origin", origin)
	git(t, bare, "fetch", "-q", "origin")
	git(t, bare, "branch", "doomed", "origin/doomed")
	git(t, bare, "branch", "--set-upstream-to=origin/doomed", "doomed")
	git(t, origin, "branch", "-D", "doomed")

	var out strings.Builder
	if err := run(root, &out); err != nil {
		t.Fatalf("run: %v", err)
	}

	if !strings.Contains(out.String(), "repo.git") {
		t.Errorf("expected output to mention repo.git, got: %q", out.String())
	}
	if contains(gitOut(t, bare, "branch", "--format=%(refname:short)"), "doomed") {
		t.Errorf("expected run to prune doomed branch")
	}
}
