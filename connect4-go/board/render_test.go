package board

import (
	"strings"
	"testing"
)

func TestRenderEmptyBoard(t *testing.T) {
	b := New()

	got := Render(b)

	wantRows := Rows + 1 // board rows + border row
	if gotRows := strings.Count(got, "\n"); gotRows != wantRows {
		t.Errorf("Render() produced %d lines, want %d", gotRows, wantRows)
	}
	if strings.Contains(got, "X") || strings.Contains(got, "O") {
		t.Errorf("Render() of empty board contains a piece symbol: %q", got)
	}
}

func TestRenderShowsDroppedPieces(t *testing.T) {
	b := New()

	if _, err := b.Drop(3); err != nil {
		t.Fatalf("Drop(3) returned error %v, want nil", err)
	}
	if _, err := b.Drop(3); err != nil {
		t.Fatalf("Drop(3) returned error %v, want nil", err)
	}

	got := Render(b)

	lines := strings.Split(strings.TrimRight(got, "\n"), "\n")
	bottomRow := lines[Rows-1]
	secondFromBottomRow := lines[Rows-2]

	if !strings.Contains(bottomRow, "X") {
		t.Errorf("bottom row %q does not contain Player1 symbol X", bottomRow)
	}
	if !strings.Contains(secondFromBottomRow, "O") {
		t.Errorf("second from bottom row %q does not contain Player2 symbol O", secondFromBottomRow)
	}
}
