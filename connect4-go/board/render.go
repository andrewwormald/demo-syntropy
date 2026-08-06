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
	sb.WriteString("|\n")
	for row := 0; row < Rows; row++ {
		for col := 0; col < Cols; col++ {
			sb.WriteString("| ")
			sb.WriteByte(symbols[b.Cell(row, col)])
			sb.WriteString(" ")
		}
		sb.WriteString("|\n")
	}
	for col := 0; col < Cols; col++ {
		sb.WriteString("----")
	}
	sb.WriteString("-\n")
	return sb.String()
}
