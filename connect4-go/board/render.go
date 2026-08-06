package board

import "strings"

// symbols maps each Cell value to the character used to render it.
var symbols = map[Cell]byte{
	Empty:   '.',
	Player1: 'X',
	Player2: 'O',
}

// Render returns an ASCII representation of the board's current state, with
// a selector row above the grid marking cursor as the column the current
// player will drop into, one row per grid row, and a border beneath the
// bottom row.
func Render(b *Board, cursor int) string {
	return strings.Join(RenderLines(b, cursor), "\n") + "\n"
}

// RenderLines returns the same ASCII representation as Render, but as a
// slice of lines rather than a single joined string, so callers can compare
// individual lines across renders.
func RenderLines(b *Board, cursor int) []string {
	lines := make([]string, 0, Rows+2)

	var sb strings.Builder
	for col := 0; col < Cols; col++ {
		sb.WriteString("| ")
		if col == cursor {
			sb.WriteString("v")
		} else {
			sb.WriteString(" ")
		}
		sb.WriteString(" ")
	}
	sb.WriteString("|")
	lines = append(lines, sb.String())

	for row := 0; row < Rows; row++ {
		sb.Reset()
		for col := 0; col < Cols; col++ {
			sb.WriteString("| ")
			sb.WriteByte(symbols[b.Cell(row, col)])
			sb.WriteString(" ")
		}
		sb.WriteString("|")
		lines = append(lines, sb.String())
	}

	sb.Reset()
	for col := 0; col < Cols; col++ {
		sb.WriteString("----")
	}
	sb.WriteString("-")
	lines = append(lines, sb.String())

	return lines
}
