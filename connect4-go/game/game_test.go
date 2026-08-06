package game

import (
	"bytes"
	"strings"
	"testing"

	"connect4-go/board"
	"connect4-go/input"
)

func TestHandleKeyMovesCursorClamped(t *testing.T) {
	g := New()

	g.HandleKey(input.Left)
	if g.Cursor != 0 {
		t.Errorf("Cursor after Left at 0 = %d, want 0", g.Cursor)
	}

	g.HandleKey(input.Right)
	g.HandleKey(input.Right)
	if g.Cursor != 2 {
		t.Errorf("Cursor after two Rights = %d, want 2", g.Cursor)
	}

	for i := 0; i < board.Cols; i++ {
		g.HandleKey(input.Right)
	}
	if g.Cursor != board.Cols-1 {
		t.Errorf("Cursor after overshooting Right = %d, want %d", g.Cursor, board.Cols-1)
	}
}

func TestHandleKeyEnterDropsPiece(t *testing.T) {
	g := New()
	g.Cursor = 3

	dropped := g.HandleKey(input.Enter)
	if !dropped {
		t.Fatalf("HandleKey(Enter) = false, want true")
	}
	if got := g.Board.Cell(board.Rows-1, 3); got != board.Player1 {
		t.Errorf("Cell(bottom, 3) = %v, want Player1", got)
	}
}

func TestHandleKeyEnterOnFullColumnDoesNotDrop(t *testing.T) {
	g := New()
	for i := 0; i < board.Rows; i++ {
		g.HandleKey(input.Enter)
	}

	dropped := g.HandleKey(input.Enter)
	if dropped {
		t.Errorf("HandleKey(Enter) on full column = true, want false")
	}
}

func TestFull(t *testing.T) {
	g := New()
	if g.Full() {
		t.Fatalf("Full() on empty board = true, want false")
	}

	for col := 0; col < board.Cols; col++ {
		g.Cursor = col
		for row := 0; row < board.Rows; row++ {
			g.HandleKey(input.Enter)
		}
	}

	if !g.Full() {
		t.Errorf("Full() on filled board = false, want true")
	}
}

func TestRunPlaysMovesFromInput(t *testing.T) {
	// Right, Right, Enter drops Player1 into column 2.
	r := bytes.NewReader([]byte{0x1b, '[', 'C', 0x1b, '[', 'C', '\r'})
	var w bytes.Buffer

	if err := Run(r, &w); err != nil {
		t.Fatalf("Run() returned error %v, want nil", err)
	}

	renders := strings.Split(strings.TrimSpace(w.String()), "\n\n")
	if len(renders) == 0 {
		t.Fatalf("Run() wrote no output")
	}
	last := renders[len(renders)-1]
	if !strings.Contains(last, "X") {
		t.Errorf("Run() final render = %q, want it to contain a dropped piece", last)
	}
}

func TestRunStopsOnEOFWithoutFullBoard(t *testing.T) {
	r := bytes.NewReader([]byte{'\r'})
	var w bytes.Buffer

	if err := Run(r, &w); err != nil {
		t.Fatalf("Run() returned error %v, want nil", err)
	}

	if w.Len() == 0 {
		t.Errorf("Run() wrote no output")
	}
}

func TestRunStopsWhenBoardFills(t *testing.T) {
	var seq []byte
	for col := 0; col < board.Cols; col++ {
		for row := 0; row < board.Rows; row++ {
			seq = append(seq, '\r')
		}
		if col < board.Cols-1 {
			seq = append(seq, 0x1b, '[', 'C')
		}
	}

	r := bytes.NewReader(seq)
	var w bytes.Buffer

	if err := Run(r, &w); err != nil {
		t.Fatalf("Run() returned error %v, want nil", err)
	}

	if !strings.Contains(w.String(), "O") {
		t.Errorf("Run() output missing Player2 pieces once board filled")
	}
}
