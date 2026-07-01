package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// prMerged reports whether branch's pull request has been merged, asking the
// GitHub CLI for merged PRs whose head is that branch.
func prMerged(repo, branch string) (bool, error) {
	out, err := runGitHub(repo, "pr", "list", "--head", branch, "--state", "merged", "--json", "number")
	if err != nil {
		return false, err
	}
	return parsePRMerged(out)
}

// parsePRMerged reports whether `gh pr list --json number` returned any PRs.
func parsePRMerged(output string) (bool, error) {
	var prs []struct {
		Number int `json:"number"`
	}
	if err := json.Unmarshal([]byte(output), &prs); err != nil {
		return false, fmt.Errorf("parse gh output: %w\n%s", err, output)
	}
	return len(prs) > 0, nil
}

func runGitHub(repo string, args ...string) (string, error) {
	cmd := exec.Command("gh", args...)
	cmd.Dir = repo
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("gh %s: %w\n%s", strings.Join(args, " "), err, out)
	}
	return string(out), nil
}
