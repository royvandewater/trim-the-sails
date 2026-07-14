package main

import (
	"strings"
	"testing"
	"time"
)

func TestRenderProgress(t *testing.T) {
	tests := []struct {
		name      string
		completed int
		total     int
		elapsed   time.Duration
		want      string
	}{
		{
			name:      "no progress yet has unknown eta",
			completed: 0,
			total:     10,
			elapsed:   0,
			want:      "[------------------------------] 0/10  eta --",
		},
		{
			name:      "half done projects remaining time",
			completed: 5,
			total:     10,
			elapsed:   10 * time.Second,
			want:      "[###############---------------] 5/10  eta 10s",
		},
		{
			name:      "complete reports done",
			completed: 10,
			total:     10,
			elapsed:   20 * time.Second,
			want:      "[##############################] 10/10  eta done",
		},
		{
			name:      "zero total renders empty",
			completed: 0,
			total:     0,
			elapsed:   0,
			want:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := renderProgress(tt.completed, tt.total, tt.elapsed); got != tt.want {
				t.Errorf("renderProgress(%d, %d, %v) = %q, want %q", tt.completed, tt.total, tt.elapsed, got, tt.want)
			}
		})
	}
}

func TestProgressBarSilentWhenNotTerminal(t *testing.T) {
	var out strings.Builder
	bar := newProgressBar(&out, 3)
	bar.advance()
	bar.advance()
	bar.finish()

	if out.String() != "" {
		t.Errorf("expected no output to a non-terminal writer, got: %q", out.String())
	}
}
