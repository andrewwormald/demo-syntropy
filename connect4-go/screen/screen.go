// Package screen writes successive text frames to a terminal, redrawing
// only the lines that changed between frames instead of reprinting
// everything.
package screen

import (
	"bytes"
	"fmt"
	"io"

	"connect4-go/framediff"
)

// Writer redraws only the changed lines of a frame on each call to Render,
// using ANSI cursor-movement escape sequences. It assumes the terminal
// cursor sits at column 0 of the frame's first line before the first call
// to Render, and leaves the cursor at column 0 of the line just below the
// last rendered line afterwards, so callers can write further output (such
// as an end-of-game announcement) immediately after.
type Writer struct {
	w    io.Writer
	prev []string
}

// New returns a Writer that writes frames to w.
func New(w io.Writer) *Writer {
	return &Writer{w: w}
}

// Render writes the lines of next that differ from the previously rendered
// frame, leaving unchanged lines untouched on screen. The first call to
// Render draws every non-empty line of next, since there is no previous
// frame to compare against. If next is identical to the previously
// rendered frame, Render writes nothing and leaves the cursor untouched.
func (s *Writer) Render(next []string) error {
	updates := framediff.Diff(s.prev, next)
	prevLen := len(s.prev)
	s.prev = next
	if len(updates) == 0 {
		return nil
	}

	var buf bytes.Buffer
	if prevLen > 0 {
		fmt.Fprintf(&buf, "\r\x1b[%dA", prevLen)
	}

	row := 0
	for _, u := range updates {
		if delta := u.Line - row; delta > 0 {
			fmt.Fprintf(&buf, "\x1b[%dB", delta)
		}
		buf.WriteString("\r\x1b[K")
		buf.WriteString(u.Text)
		row = u.Line
	}

	if delta := len(next) - row; delta > 0 {
		fmt.Fprintf(&buf, "\x1b[%dB", delta)
	}
	buf.WriteString("\r")

	_, err := s.w.Write(buf.Bytes())
	return err
}
