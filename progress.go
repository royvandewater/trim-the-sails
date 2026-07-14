package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

const barWidth = 30

// renderProgress formats a single-line progress bar for completed of total
// units of work, projecting an ETA for the remaining work from elapsed. Before
// any work finishes the ETA is unknown ("--"); once every unit is done it reads
// "done". A total of zero renders an empty string.
func renderProgress(completed, total int, elapsed time.Duration) string {
	if total == 0 {
		return ""
	}

	filled := completed * barWidth / total
	bar := strings.Repeat("#", filled) + strings.Repeat("-", barWidth-filled)

	eta := "--"
	switch {
	case completed >= total:
		eta = "done"
	case completed > 0:
		remaining := time.Duration(int64(elapsed) / int64(completed) * int64(total-completed))
		eta = remaining.Round(time.Second).String()
	}

	return fmt.Sprintf("[%s] %d/%d  eta %s", bar, completed, total, eta)
}

// progressBar renders a live progress bar as work completes. It only draws when
// its writer is a terminal, so piped or redirected output stays clean. advance
// is safe to call from multiple goroutines.
type progressBar struct {
	w     io.Writer
	total int
	start time.Time
	live  bool

	mu      sync.Mutex
	done    int
	lastLen int
}

// progressFrame builds an in-place terminal update: a carriage return, line,
// and enough trailing spaces to erase leftovers from a previously longer line
// (prevLen). It returns the frame to write and line's length.
func progressFrame(line string, prevLen int) (string, int) {
	pad := prevLen - len(line)
	if pad < 0 {
		pad = 0
	}
	return "\r" + line + strings.Repeat(" ", pad), len(line)
}

func newProgressBar(w io.Writer, total int) *progressBar {
	b := &progressBar{w: w, total: total, start: time.Now(), live: isTerminal(w)}
	b.draw()
	return b
}

func (b *progressBar) advance() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.done++
	b.render()
}

func (b *progressBar) draw() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.render()
}

// render redraws the bar in place. Callers must hold b.mu.
func (b *progressBar) render() {
	if !b.live {
		return
	}
	frame, n := progressFrame(renderProgress(b.done, b.total, time.Since(b.start)), b.lastLen)
	b.lastLen = n
	fmt.Fprint(b.w, frame)
}

// finish moves off the progress line so later output starts fresh.
func (b *progressBar) finish() {
	if !b.live {
		return
	}
	fmt.Fprintln(b.w)
}

// isTerminal reports whether w is a character device (a terminal), so we only
// draw the live bar when a human is watching.
func isTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}
