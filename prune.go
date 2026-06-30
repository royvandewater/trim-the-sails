package main

import "strings"

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
