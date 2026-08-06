// Package game wires the board, input, and render packages into an
// interactive Connect 4 game loop.
package game

import (
	"io"

	"connect4-go/board"
	"connect4-go/input"
)

// Game holds the in-progress board state and the column currently selected
// by the player for their next drop.
type Game struct {
	Board  *board.Board
	Cursor int
}

// New returns a Game with an empty board and the cursor on the leftmost
// column.
func New() *Game {
	return &Game{Board: board.New()}
}

// HandleKey applies a decoded key event to the game state. Left and Right
// move the cursor, clamped to the board's columns. Enter drops the current
// player's piece into the selected column; a full column is silently
// ignored, matching how a player would just try another column. It reports
// whether a piece was dropped.
func (g *Game) HandleKey(key input.Key) bool {
	switch key {
	case input.Left:
		if g.Cursor > 0 {
			g.Cursor--
		}
	case input.Right:
		if g.Cursor < board.Cols-1 {
			g.Cursor++
		}
	case input.Enter:
		if _, err := g.Board.Drop(g.Cursor); err == nil {
			return true
		}
	}
	return false
}

// Full reports whether every column is filled, ending the game in a draw.
func (g *Game) Full() bool {
	for col := 0; col < board.Cols; col++ {
		if g.Board.Cell(0, col) == board.Empty {
			return false
		}
	}
	return true
}

// Run decodes raw terminal byte sequences from r and drives the game
// against them, writing the rendered board to w after every recognized key.
// It returns when the board fills up or r is exhausted.
func Run(r io.Reader, w io.Writer) error {
	g := New()

	if _, err := io.WriteString(w, board.Render(g.Board)); err != nil {
		return err
	}

	var pending []byte
	chunk := make([]byte, 64)
	for {
		n, err := r.Read(chunk)
		if n > 0 {
			pending = append(pending, chunk[:n]...)
		}

		for len(pending) > 0 {
			key, consumed := input.Decode(pending)
			if consumed == 0 {
				break
			}
			pending = pending[consumed:]

			if g.HandleKey(key) {
				if _, werr := io.WriteString(w, board.Render(g.Board)); werr != nil {
					return werr
				}
				if g.Full() {
					return nil
				}
			}
		}

		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}
