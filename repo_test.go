package main

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestPruneRepoDeletesBranchesWithGoneUpstream(t *testing.T) {
	root := t.TempDir()
	origin := filepath.Join(root, "origin")
	bare := filepath.Join(root, "repo.git")

	// Origin with main + two feature branches.
	git(t, root, "init", "-q", origin)
	git(t, origin, "config", "user.email", "test@example.com")
	git(t, origin, "config", "user.name", "Test")
	git(t, origin, "commit", "-q", "--allow-empty", "-m", "init")
	git(t, origin, "branch", "keep")
	git(t, origin, "branch", "doomed")

	// Bare repo tracking origin, with local branches set to track it.
	git(t, root, "init", "-q", "--bare", bare)
	git(t, bare, "remote", "add", "origin", origin)
	git(t, bare, "fetch", "-q", "origin")
	git(t, bare, "branch", "keep", "origin/keep")
	git(t, bare, "branch", "doomed", "origin/doomed")
	git(t, bare, "branch", "--set-upstream-to=origin/keep", "keep")
	git(t, bare, "branch", "--set-upstream-to=origin/doomed", "doomed")

	// Delete doomed on origin so its upstream becomes gone after prune.
	git(t, origin, "branch", "-D", "doomed")

	if err := pruneRepo(bare); err != nil {
		t.Fatalf("pruneRepo: %v", err)
	}

	branches := gitOut(t, bare, "branch", "--format=%(refname:short)")
	if contains(branches, "doomed") {
		t.Errorf("expected doomed branch to be pruned, got: %v", branches)
	}
	if !contains(branches, "keep") {
		t.Errorf("expected keep branch to survive, got: %v", branches)
	}
}

func git(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

func gitOut(t *testing.T, dir string, args ...string) []string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	var lines []string
	for _, l := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if l != "" {
			lines = append(lines, l)
		}
	}
	return lines
}

func contains(haystack []string, needle string) bool {
	for _, h := range haystack {
		if h == needle {
			return true
		}
	}
	return false
}
