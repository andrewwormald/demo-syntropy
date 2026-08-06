package game

import (
	"bytes"
	"fmt"
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

// TestRunEmitsDiffBasedFrames drives Run through a scripted key sequence
// and asserts the literal byte stream it writes: a full initial draw of
// every board line, then for each subsequent key only the cursor-move,
// clear-line, and content bytes for the lines that actually changed,
// exactly as screen.Writer documents. This ties frame, framediff, screen,
// and the game loop together end to end, rather than only checking that
// the final render contains an expected substring.
func TestRunEmitsDiffBasedFrames(t *testing.T) {
	// Right, Right, Enter: moves the cursor to column 2 and drops
	// Player1's piece there.
	r := bytes.NewReader([]byte{0x1b, '[', 'C', 0x1b, '[', 'C', '\r'})
	var w bytes.Buffer

	if err := Run(r, &w); err != nil {
		t.Fatalf("Run() returned error %v, want nil", err)
	}

	initial := board.RenderLines(board.New(), 0)
	afterRight1 := board.RenderLines(board.New(), 1)
	afterRight2 := board.RenderLines(board.New(), 2)
	dropped := board.New()
	if _, err := dropped.Drop(2); err != nil {
		t.Fatalf("Drop(2) returned error %v, want nil", err)
	}
	afterEnter := board.RenderLines(dropped, 2)

	numLines := len(initial)

	var want bytes.Buffer
	// Initial full draw: every line is "changed" since there is no
	// previous frame, so each is preceded by a cursor-down-and-clear and
	// followed by a move to the next line.
	for _, line := range initial {
		want.WriteString("\r\x1b[K")
		want.WriteString(line)
		want.WriteString("\x1b[1B")
	}
	want.WriteString("\r")

	// Right: only the header line (index 0) changed, so the writer moves
	// back to the top, rewrites just that line, then moves back down to
	// below the last line.
	fmt.Fprintf(&want, "\r\x1b[%dA", numLines)
	want.WriteString("\r\x1b[K")
	want.WriteString(afterRight1[0])
	fmt.Fprintf(&want, "\x1b[%dB\r", numLines)

	// Right again: same shape, header line only.
	fmt.Fprintf(&want, "\r\x1b[%dA", numLines)
	want.WriteString("\r\x1b[K")
	want.WriteString(afterRight2[0])
	fmt.Fprintf(&want, "\x1b[%dB\r", numLines)

	// Enter: only the row the piece landed in (index 6: the header line
	// plus the bottom grid row) changed, so the writer skips down past the
	// unchanged lines above it before clearing and rewriting it.
	droppedLine := 6
	fmt.Fprintf(&want, "\r\x1b[%dA", numLines)
	fmt.Fprintf(&want, "\x1b[%dB", droppedLine)
	want.WriteString("\r\x1b[K")
	want.WriteString(afterEnter[droppedLine])
	fmt.Fprintf(&want, "\x1b[%dB\r", numLines-droppedLine)

	if got := w.String(); got != want.String() {
		t.Errorf("Run() wrote:\n%q\nwant:\n%q", got, want.String())
	}
}

func TestRunStopsOnQuit(t *testing.T) {
	// Right, then Ctrl-C: the cursor move should be rendered but the
	// game should stop before any further input is read.
	r := bytes.NewReader([]byte{0x1b, '[', 'C', 0x03, 0x1b, '[', 'C'})
	var w bytes.Buffer

	if err := Run(r, &w); err != nil {
		t.Fatalf("Run() returned error %v, want nil", err)
	}

	if w.Len() == 0 {
		t.Errorf("Run() wrote no output")
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

func TestRunAnnouncesWinner(t *testing.T) {
	// Alternating drops into columns 0 and 1 land Player1's pieces at every
	// row of column 0 (Player2 always drops into column 1 in between),
	// connecting four for Player1 vertically.
	enter := byte('\r')
	right := []byte{0x1b, '[', 'C'}
	left := []byte{0x1b, '[', 'D'}
	var keys []byte
	for i := 0; i < 4; i++ {
		keys = append(keys, enter)
		if i < 3 {
			keys = append(keys, right...)
			keys = append(keys, enter)
			keys = append(keys, left...)
		}
	}

	r := bytes.NewReader(keys)
	var w bytes.Buffer

	if err := Run(r, &w); err != nil {
		t.Fatalf("Run() returned error %v, want nil", err)
	}

	if !strings.Contains(w.String(), "Player 1 wins!") {
		t.Errorf("Run() output = %q, want it to contain a win announcement", w.String())
	}
}

func TestAnnouncementDraw(t *testing.T) {
	// Column order that fills the entire board without ever connecting
	// four, found by randomized search over valid drop sequences. Applied
	// directly through the board rather than via Run's raw key decoding,
	// since driving 42 drops across all 7 columns would need far more
	// arrow-key bytes than fit in Run's read buffer.
	drawCols := []int{
		5, 3, 6, 6, 4, 4, 5, 6, 2, 4, 0, 3, 2, 1, 5, 1, 6, 5, 6, 5, 4, 1,
		2, 6, 0, 4, 4, 0, 1, 2, 5, 3, 1, 2, 3, 2, 3, 0, 3, 1, 0, 0,
	}

	b := board.New()
	for _, col := range drawCols {
		if _, err := b.Drop(col); err != nil {
			t.Fatalf("Drop(%d) returned error %v, want nil", col, err)
		}
	}

	if got, want := announcement(b), "Draw!\n"; got != want {
		t.Errorf("announcement() = %q, want %q", got, want)
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
