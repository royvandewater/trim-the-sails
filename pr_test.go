package main

import "testing"

func TestParsePRMerged(t *testing.T) {
	cases := []struct {
		name   string
		output string
		want   bool
	}{
		{"no PR", "[]\n", false},
		{"merged PR", `[{"number":7}]` + "\n", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parsePRMerged(tc.output)
			if err != nil {
				t.Fatalf("parsePRMerged(%q): %v", tc.output, err)
			}
			if got != tc.want {
				t.Errorf("parsePRMerged(%q) = %v, want %v", tc.output, got, tc.want)
			}
		})
	}
}
