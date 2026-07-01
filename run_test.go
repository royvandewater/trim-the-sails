package main

import (
	"fmt"
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

func TestRunReportsReposInDirectoryOrder(t *testing.T) {
	root := t.TempDir()
	names := []string{"alpha.git", "beta.git", "gamma.git"}
	var bares []string
	for _, name := range names {
		bares = append(bares, makeBareRepoWithDoomedBranch(t, root, name))
	}

	var out strings.Builder
	if err := run(root, &out); err != nil {
		t.Fatalf("run: %v", err)
	}

	var want strings.Builder
	for _, bare := range bares {
		fmt.Fprintln(&want, "pruning", bare)
	}
	if out.String() != want.String() {
		t.Errorf("output not in directory order:\ngot:  %q\nwant: %q", out.String(), want.String())
	}

	for _, bare := range bares {
		if contains(gitOut(t, bare, "branch", "--format=%(refname:short)"), "doomed") {
			t.Errorf("expected doomed branch pruned in %s", bare)
		}
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
