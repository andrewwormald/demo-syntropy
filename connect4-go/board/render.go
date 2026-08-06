package board

import "strings"

// symbols maps each Cell value to the character used to render it.
var symbols = map[Cell]byte{
	Empty:   '.',
	Player1: 'X',
	Player2: 'O',
}

// Render returns an ASCII representation of the board's current state, with
// one row per grid row followed by a border beneath the bottom row.
func Render(b *Board) string {
	var sb strings.Builder
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
