package board

import (
	"strings"
	"testing"
)

func TestRenderEmptyBoard(t *testing.T) {
	b := New()

	got := Render(b, 0)

	wantRows := Rows + 2 // selector row + board rows + border row
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

	got := Render(b, 0)

	lines := strings.Split(strings.TrimRight(got, "\n"), "\n")
	const selectorRows = 1
	bottomRow := lines[selectorRows+Rows-1]
	secondFromBottomRow := lines[selectorRows+Rows-2]

	if !strings.Contains(bottomRow, "X") {
		t.Errorf("bottom row %q does not contain Player1 symbol X", bottomRow)
	}
	if !strings.Contains(secondFromBottomRow, "O") {
		t.Errorf("second from bottom row %q does not contain Player2 symbol O", secondFromBottomRow)
	}
}

func TestRenderLinesMatchesRender(t *testing.T) {
	b := New()
	if _, err := b.Drop(3); err != nil {
		t.Fatalf("Drop(3) returned error %v, want nil", err)
	}

	lines := RenderLines(b, 2)
	got := strings.Join(lines, "\n") + "\n"

	want := Render(b, 2)
	if got != want {
		t.Errorf("RenderLines() joined = %q, want Render() = %q", got, want)
	}
}

func TestRenderLinesCount(t *testing.T) {
	b := New()

	lines := RenderLines(b, 0)

	wantLines := Rows + 2 // selector row + board rows + border row
	if len(lines) != wantLines {
		t.Errorf("RenderLines() returned %d lines, want %d", len(lines), wantLines)
	}
	for i, line := range lines {
		if strings.Contains(line, "\n") {
			t.Errorf("RenderLines()[%d] = %q contains a newline", i, line)
		}
	}
}

func TestRenderLinesShowsDroppedPieces(t *testing.T) {
	b := New()
	if _, err := b.Drop(3); err != nil {
		t.Fatalf("Drop(3) returned error %v, want nil", err)
	}
	if _, err := b.Drop(3); err != nil {
		t.Fatalf("Drop(3) returned error %v, want nil", err)
	}

	lines := RenderLines(b, 0)

	const selectorRows = 1
	bottomRow := lines[selectorRows+Rows-1]
	secondFromBottomRow := lines[selectorRows+Rows-2]

	if !strings.Contains(bottomRow, "X") {
		t.Errorf("bottom row %q does not contain Player1 symbol X", bottomRow)
	}
	if !strings.Contains(secondFromBottomRow, "O") {
		t.Errorf("second from bottom row %q does not contain Player2 symbol O", secondFromBottomRow)
	}
}

func TestRenderSelectorMarksCursorColumn(t *testing.T) {
	b := New()

	for _, cursor := range []int{0, 3, Cols - 1} {
		got := Render(b, cursor)

		lines := strings.Split(strings.TrimRight(got, "\n"), "\n")
		selectorRow := lines[0]
		cells := strings.Split(strings.Trim(selectorRow, "|"), "|")
		if len(cells) != Cols {
			t.Fatalf("Render(_, %d) selector row %q has %d cells, want %d", cursor, selectorRow, len(cells), Cols)
		}
		for col, cell := range cells {
			marked := strings.Contains(cell, "v")
			if col == cursor && !marked {
				t.Errorf("Render(_, %d) selector row %q does not mark column %d", cursor, selectorRow, col)
			}
			if col != cursor && marked {
				t.Errorf("Render(_, %d) selector row %q unexpectedly marks column %d", cursor, selectorRow, col)
			}
		}
	}
}
