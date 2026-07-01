package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// pruneRepo fetches with --prune to drop stale remote-tracking refs, removes
// any worktree whose branch's PR has been merged, then deletes any local
// branches whose upstream is now gone.
func pruneRepo(repo string) error {
	if _, err := runGit(repo, "fetch", "--all", "--prune"); err != nil {
		return err
	}

	if err := pruneMergedWorktrees(repo, func(branch string) (bool, error) {
		return prMerged(repo, branch)
	}); err != nil {
		return err
	}

	out, err := runGit(repo, "branch", "-vv")
	if err != nil {
		return err
	}

	for _, branch := range parseGoneBranches(out) {
		if _, err := runGit(repo, "branch", "-D", branch); err != nil {
			return err
		}
	}
	return nil
}

func runGit(repo string, args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", repo}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s: %w\n%s", strings.Join(args, " "), err, out)
	}
	return string(out), nil
}

// parseGoneBranches reads `git branch -vv` output and returns the names of
// branches whose upstream has been deleted, marked by git as ": gone]".
func parseGoneBranches(output string) []string {
	var gone []string
	for _, line := range strings.Split(output, "\n") {
		if !strings.Contains(line, ": gone]") {
			continue
		}
		// Strip the leading marker ("* ", "+ ", or "  ") then take the
		// branch name, which is the first whitespace-delimited field.
		fields := strings.Fields(strings.TrimLeft(line, "*+ "))
		if len(fields) == 0 {
			continue
		}
		gone = append(gone, fields[0])
	}
	return gone
}
