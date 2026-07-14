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

func TestProgressFrameErasesLeftoverChars(t *testing.T) {
	// A shorter line must be padded so leftovers from a previous longer line
	// (e.g. "eta 1m15s" -> "eta 7s") do not linger on screen.
	frame, n := progressFrame("7s", len("1m15s"))
	if frame != "\r7s   " {
		t.Errorf("frame = %q, want %q", frame, "\r7s   ")
	}
	if n != 2 {
		t.Errorf("length = %d, want 2", n)
	}

	// A longer-or-equal line needs no padding.
	frame, n = progressFrame("1m15s", 2)
	if frame != "\r1m15s" {
		t.Errorf("frame = %q, want %q", frame, "\r1m15s")
	}
	if n != 5 {
		t.Errorf("length = %d, want 5", n)
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
