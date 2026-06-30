package main

import "testing"

func TestParseGoneBranches(t *testing.T) {
	output := "" +
		"  feature-a abc1234 [origin/feature-a: gone] old work\n" +
		"* main      def5678 [origin/main] current\n" +
		"  feature-b aaa1111 [origin/feature-b: ahead 1] local only\n" +
		"+ feature-c bbb2222 [origin/feature-c: gone] in a worktree\n" +
		"  local-only ccc3333 no upstream at all\n"

	got := parseGoneBranches(output)

	want := []string{"feature-a", "feature-c"}
	if !equalStrings(got, want) {
		t.Fatalf("parseGoneBranches() = %v, want %v", got, want)
	}
}
