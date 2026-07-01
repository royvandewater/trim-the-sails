package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestParseWorktreesReturnsBranchedCheckouts(t *testing.T) {
	output := "" +
		"worktree /repos/repo.git\n" +
		"bare\n" +
		"\n" +
		"worktree /repos/repo.git/worktrees/feature-a\n" +
		"HEAD abc1234\n" +
		"branch refs/heads/feature-a\n" +
		"\n" +
		"worktree /repos/repo.git/worktrees/detached\n" +
		"HEAD def5678\n" +
		"detached\n" +
		"\n"

	got := parseWorktrees(output)

	want := []worktree{
		{Path: "/repos/repo.git/worktrees/feature-a", Branch: "feature-a"},
	}
	if len(got) != len(want) {
		t.Fatalf("parseWorktrees() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("parseWorktrees() = %v, want %v", got, want)
		}
	}
}

func TestPruneMergedWorktreesRemovesMergedAndKeepsRest(t *testing.T) {
	root := t.TempDir()
	origin := filepath.Join(root, "origin")
	bare := filepath.Join(root, "repo.git")

	git(t, root, "init", "-q", origin)
	git(t, origin, "config", "user.email", "test@example.com")
	git(t, origin, "config", "user.name", "Test")
	git(t, origin, "commit", "-q", "--allow-empty", "-m", "init")
	git(t, origin, "branch", "merged")
	git(t, origin, "branch", "open")

	git(t, root, "init", "-q", "--bare", bare)
	git(t, bare, "remote", "add", "origin", origin)
	git(t, bare, "fetch", "-q", "origin")

	mergedWT := filepath.Join(root, "merged")
	openWT := filepath.Join(root, "open")
	git(t, bare, "worktree", "add", "-q", mergedWT, "merged")
	git(t, bare, "worktree", "add", "-q", openWT, "open")

	// Fake merge check: only the "merged" branch's PR is merged.
	isMerged := func(branch string) (bool, error) {
		return branch == "merged", nil
	}

	if err := pruneMergedWorktrees(bare, isMerged); err != nil {
		t.Fatalf("pruneMergedWorktrees: %v", err)
	}

	paths := gitOut(t, bare, "worktree", "list", "--porcelain")
	if containsSubstring(paths, mergedWT) {
		t.Errorf("expected merged worktree removed, got: %v", paths)
	}
	if !containsSubstring(paths, openWT) {
		t.Errorf("expected open worktree to survive, got: %v", paths)
	}
}

func containsSubstring(lines []string, substr string) bool {
	for _, l := range lines {
		if strings.Contains(l, substr) {
			return true
		}
	}
	return false
}
