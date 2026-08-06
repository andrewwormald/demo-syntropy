// Package board implements the Connect 4 grid state and drop logic.
package board

import "errors"

const (
	Rows = 6
	Cols = 7
)

// Cell represents the contents of a single grid position.
type Cell int

const (
	Empty Cell = iota
	Player1
	Player2
)

var (
	ErrInvalidColumn = errors.New("board: column out of range")
	ErrColumnFull    = errors.New("board: column is full")
)

// Board holds the Connect 4 grid state and whose turn it is.
type Board struct {
	grid    [Rows][Cols]Cell
	current Cell
}

// New returns an empty board with Player1 to move first.
func New() *Board {
	return &Board{current: Player1}
}

// CurrentPlayer returns the player whose turn it is to drop next.
func (b *Board) CurrentPlayer() Cell {
	return b.current
}

// Cell returns the contents of the grid at (row, col).
func (b *Board) Cell(row, col int) Cell {
	return b.grid[row][col]
}

// Drop places the current player's piece into the given column, settling it
// to the lowest empty row. It returns the row the piece landed in and
// advances the current player. ErrInvalidColumn is returned if col is out of
// range, and ErrColumnFull if the column has no empty rows.
func (b *Board) Drop(col int) (int, error) {
	if col < 0 || col >= Cols {
		return -1, ErrInvalidColumn
	}

	for row := Rows - 1; row >= 0; row-- {
		if b.grid[row][col] == Empty {
			b.grid[row][col] = b.current
			if b.current == Player1 {
				b.current = Player2
			} else {
				b.current = Player1
			}
			return row, nil
		}
	}

	return -1, ErrColumnFull
}
