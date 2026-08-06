package main

import "testing"

func TestRenderBoard(t *testing.T) {
	got := renderBoard()

	wantRows := rows + 1 // board rows + border row
	gotRows := 0
	for _, c := range got {
		if c == '\n' {
			gotRows++
		}
	}
	if gotRows != wantRows {
		t.Errorf("renderBoard() produced %d lines, want %d", gotRows, wantRows)
	}
}
