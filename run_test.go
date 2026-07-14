package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestRunPrunesEachBareRepo(t *testing.T) {
	root := t.TempDir()
	names := []string{"alpha.git", "beta.git", "gamma.git"}
	var bares []string
	for _, name := range names {
		bares = append(bares, makeBareRepoWithDoomedBranch(t, root, name))
	}

	var out strings.Builder
	if err := run(root, nil, &out); err != nil {
		t.Fatalf("run: %v", err)
	}

	// The progress bar only draws to a terminal; a strings.Builder is not one,
	// so nothing should be written here.
	if out.String() != "" {
		t.Errorf("expected no output to a non-terminal writer, got: %q", out.String())
	}

	for _, bare := range bares {
		if contains(gitOut(t, bare, "branch", "--format=%(refname:short)"), "doomed") {
			t.Errorf("expected doomed branch pruned in %s", bare)
		}
	}
}

func TestRunHelpFlagPrintsUsageAndPrunesNothing(t *testing.T) {
	root := t.TempDir()
	bare := makeBareRepoWithDoomedBranch(t, root, "alpha.git")

	for _, flag := range []string{"--help", "-h"} {
		var out strings.Builder
		if err := run(root, []string{flag}, &out); err != nil {
			t.Fatalf("run %s: %v", flag, err)
		}
		if !strings.Contains(out.String(), "Usage") {
			t.Errorf("run %s: expected usage output, got: %q", flag, out.String())
		}
		if !contains(gitOut(t, bare, "branch", "--format=%(refname:short)"), "doomed") {
			t.Errorf("run %s: expected no pruning", flag)
		}
	}
}

func TestRunOnlyPrunesNamedRepos(t *testing.T) {
	root := t.TempDir()
	names := []string{"alpha.git", "beta.git", "gamma.git"}
	bares := make(map[string]string)
	for _, name := range names {
		bares[name] = makeBareRepoWithDoomedBranch(t, root, name)
	}

	var out strings.Builder
	if err := run(root, []string{"beta"}, &out); err != nil {
		t.Fatalf("run: %v", err)
	}

	if strings.Contains(out.String(), "alpha.git") || strings.Contains(out.String(), "gamma.git") {
		t.Errorf("expected only beta.git pruned, got: %q", out.String())
	}
	if contains(gitOut(t, bares["beta.git"], "branch", "--format=%(refname:short)"), "doomed") {
		t.Errorf("expected doomed branch pruned in beta.git")
	}
	if !contains(gitOut(t, bares["alpha.git"], "branch", "--format=%(refname:short)"), "doomed") {
		t.Errorf("expected doomed branch untouched in alpha.git")
	}
}

// makeBareRepoWithDoomedBranch builds an origin plus a bare repo tracking it,
// with a "doomed" branch whose upstream is deleted so a prune should remove it.
// It returns the bare repo's path.
func makeBareRepoWithDoomedBranch(t *testing.T, root, name string) string {
	t.Helper()
	origin := filepath.Join(root, strings.TrimSuffix(name, ".git")+".origin")
	bare := filepath.Join(root, name)

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

	return bare
}
