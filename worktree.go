package main

import "strings"

// worktree is a linked working tree of a bare repo, identified by its path on
// disk and the short name of the branch checked out in it.
type worktree struct {
	Path   string
	Branch string
}

// pruneMergedWorktrees removes every linked worktree of repo whose branch's
// pull request has been merged, as reported by isMerged. Removing the worktree
// first frees its branch so the later gone-branch prune can delete it.
func pruneMergedWorktrees(repo string, isMerged func(branch string) (bool, error)) error {
	out, err := runGit(repo, "worktree", "list", "--porcelain")
	if err != nil {
		return err
	}

	for _, wt := range parseWorktrees(out) {
		merged, err := isMerged(wt.Branch)
		if err != nil {
			return err
		}
		if !merged {
			continue
		}
		if _, err := runGit(repo, "worktree", "remove", "--force", wt.Path); err != nil {
			return err
		}
	}
	return nil
}

// parseWorktrees reads `git worktree list --porcelain` output and returns the
// worktrees that have a branch checked out. The bare repo entry and any
// detached-HEAD worktrees have no branch, so they are skipped.
func parseWorktrees(output string) []worktree {
	var worktrees []worktree
	var path string
	for _, line := range strings.Split(output, "\n") {
		switch {
		case strings.HasPrefix(line, "worktree "):
			path = strings.TrimPrefix(line, "worktree ")
		case strings.HasPrefix(line, "branch "):
			ref := strings.TrimPrefix(line, "branch ")
			worktrees = append(worktrees, worktree{
				Path:   path,
				Branch: strings.TrimPrefix(ref, "refs/heads/"),
			})
		}
	}
	return worktrees
}
