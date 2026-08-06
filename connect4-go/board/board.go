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
	Empty   Cell = 0
	Player1 Cell = 1
	Player2 Cell = 2
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

// winDirections are the four axes along which four-in-a-row can occur:
// horizontal, vertical, and both diagonals. Each is checked in both senses
// by scanning outward from every occupied cell.
var winDirections = [4][2]int{
	{0, 1},
	{1, 0},
	{1, 1},
	{1, -1},
}

// Winner returns the player with four connected pieces in a row, column, or
// diagonal, or Empty if no player has won yet.
func (b *Board) Winner() Cell {
	for row := 0; row < Rows; row++ {
		for col := 0; col < Cols; col++ {
			player := b.grid[row][col]
			if player == Empty {
				continue
			}

			for _, dir := range winDirections {
				count := 1
				r, c := row+dir[0], col+dir[1]
				for count < 4 && r >= 0 && r < Rows && c >= 0 && c < Cols && b.grid[r][c] == player {
					count++
					r += dir[0]
					c += dir[1]
				}
				if count == 4 {
					return player
				}
			}
		}
	}
	return Empty
}

// Full reports whether every column is filled, meaning no more moves are
// possible.
func (b *Board) Full() bool {
	for col := 0; col < Cols; col++ {
		if b.grid[0][col] == Empty {
			return false
		}
	}
	return true
}
