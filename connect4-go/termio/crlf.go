// Package termio provides small io.Writer helpers for talking to a
// terminal that has been put into raw mode.
package termio

import "io"

// CRLFWriter wraps an io.Writer and translates each "\n" not already
// preceded by "\r" into "\r\n" before writing. Raw mode (as set by
// term.MakeRaw) disables the tty driver's output post-processing, so a
// lone "\n" no longer returns the cursor to column 0; writes made while
// raw mode is active need to supply the "\r" themselves.
type CRLFWriter struct {
	w        io.Writer
	lastByte byte
}

// NewCRLFWriter returns a CRLFWriter wrapping w.
func NewCRLFWriter(w io.Writer) *CRLFWriter {
	return &CRLFWriter{w: w}
}

// Write implements io.Writer, translating "\n" to "\r\n" (unless already
// preceded by "\r", including across separate Write calls) and writing
// the result to the underlying writer. It returns the number of bytes of
// p consumed; on success this is always len(p), even though more bytes
// may have been written to the underlying writer.
func (c *CRLFWriter) Write(p []byte) (int, error) {
	start := 0
	prev := c.lastByte
	for i, b := range p {
		if b == '\n' && prev != '\r' {
			if _, err := c.w.Write(p[start:i]); err != nil {
				return start, err
			}
			if _, err := c.w.Write([]byte("\r\n")); err != nil {
				return start, err
			}
			start = i + 1
		}
		prev = b
	}
	if len(p) > 0 {
		c.lastByte = p[len(p)-1]
	}
	if start < len(p) {
		if _, err := c.w.Write(p[start:]); err != nil {
			return start, err
		}
	}
	return len(p), nil
}
