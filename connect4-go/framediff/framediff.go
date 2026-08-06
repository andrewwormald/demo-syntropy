// Package framediff computes the minimal set of line changes between two
// rendered text frames so a caller can redraw only the lines that changed
// instead of reprinting the whole frame.
package framediff

// LineUpdate describes a single line that changed between two frames.
type LineUpdate struct {
	// Line is the zero-based index of the changed line.
	Line int
	// Text is the new content of the line.
	Text string
}

// Diff compares prev against next and returns a LineUpdate for every line
// whose content differs, in ascending line order. A nil or empty prev (for
// example, when there is no previous frame yet) is treated as if every line
// in next were previously empty, so every non-empty line in next is
// reported as changed. If next is shorter than prev, lines beyond
// len(next)-1 are not reported, since a caller only redraws lines that
// still exist.
func Diff(prev, next []string) []LineUpdate {
	var updates []LineUpdate
	for i, line := range next {
		var prevLine string
		if i < len(prev) {
			prevLine = prev[i]
		}
		if line != prevLine {
			updates = append(updates, LineUpdate{Line: i, Text: line})
		}
	}
	return updates
}
