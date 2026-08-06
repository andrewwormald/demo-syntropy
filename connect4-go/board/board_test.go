package board

import "testing"

func TestNew(t *testing.T) {
	b := New()

	if got := b.CurrentPlayer(); got != Player1 {
		t.Errorf("CurrentPlayer() = %v, want %v", got, Player1)
	}

	for row := 0; row < Rows; row++ {
		for col := 0; col < Cols; col++ {
			if got := b.Cell(row, col); got != Empty {
				t.Errorf("Cell(%d, %d) = %v, want Empty", row, col, got)
			}
		}
	}
}

func TestDropIntoEmptyColumn(t *testing.T) {
	b := New()

	row, err := b.Drop(3)
	if err != nil {
		t.Fatalf("Drop(3) returned error %v, want nil", err)
	}
	if want := Rows - 1; row != want {
		t.Errorf("Drop(3) row = %d, want %d", row, want)
	}
	if got := b.Cell(row, 3); got != Player1 {
		t.Errorf("Cell(%d, 3) = %v, want Player1", row, got)
	}
	if got := b.CurrentPlayer(); got != Player2 {
		t.Errorf("CurrentPlayer() after drop = %v, want Player2", got)
	}
}

func TestDropIntoPartiallyFilledColumn(t *testing.T) {
	b := New()

	wantPlayers := []Cell{Player1, Player2, Player1}
	for i, want := range wantPlayers {
		row, err := b.Drop(0)
		if err != nil {
			t.Fatalf("Drop(0) call %d returned error %v, want nil", i, err)
		}

		wantRow := Rows - 1 - i
		if row != wantRow {
			t.Fatalf("Drop(0) call %d row = %d, want %d", i, row, wantRow)
		}
		if got := b.Cell(row, 0); got != want {
			t.Errorf("Cell(%d, 0) call %d = %v, want %v", row, i, got, want)
		}
	}
}

func TestDropIntoFullColumn(t *testing.T) {
	b := New()

	for i := 0; i < Rows; i++ {
		if _, err := b.Drop(2); err != nil {
			t.Fatalf("Drop(2) fill call %d returned error %v, want nil", i, err)
		}
	}

	row, err := b.Drop(2)
	if err != ErrColumnFull {
		t.Errorf("Drop(2) on full column error = %v, want ErrColumnFull", err)
	}
	if row != -1 {
		t.Errorf("Drop(2) on full column row = %d, want -1", row)
	}
}

func TestDropInvalidColumn(t *testing.T) {
	tests := []struct {
		name string
		col  int
	}{
		{"negative", -1},
		{"equal to Cols", Cols},
		{"far out of range", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := New()

			row, err := b.Drop(tt.col)
			if err != ErrInvalidColumn {
				t.Errorf("Drop(%d) error = %v, want ErrInvalidColumn", tt.col, err)
			}
			if row != -1 {
				t.Errorf("Drop(%d) row = %d, want -1", tt.col, row)
			}
		})
	}
}

func TestWinnerNoneOnEmptyBoard(t *testing.T) {
	b := New()

	if got := b.Winner(); got != Empty {
		t.Errorf("Winner() on empty board = %v, want Empty", got)
	}
}

func TestWinnerHorizontal(t *testing.T) {
	b := New()

	// Player1 drops in columns 0-3, Player2 drops elsewhere in between so
	// the row stays Player1's.
	for col := 0; col < 4; col++ {
		if _, err := b.Drop(col); err != nil {
			t.Fatalf("Drop(%d) returned error %v, want nil", col, err)
		}
		if col < 3 {
			if _, err := b.Drop(col); err != nil {
				t.Fatalf("Drop(%d) returned error %v, want nil", col, err)
			}
		}
	}

	if got := b.Winner(); got != Player1 {
		t.Errorf("Winner() = %v, want Player1", got)
	}
}

func TestWinnerVertical(t *testing.T) {
	b := New()

	for i := 0; i < 4; i++ {
		if _, err := b.Drop(0); err != nil {
			t.Fatalf("Drop(0) call %d returned error %v, want nil", i, err)
		}
		if i < 3 {
			if _, err := b.Drop(1); err != nil {
				t.Fatalf("Drop(1) call %d returned error %v, want nil", i, err)
			}
		}
	}

	if got := b.Winner(); got != Player1 {
		t.Errorf("Winner() = %v, want Player1", got)
	}
}

func TestWinnerDiagonalUp(t *testing.T) {
	b := New()

	// Drop order chosen so Player1's turn lands on (5,0), (4,1), (3,2), and
	// (2,3), a rising diagonal. Column 4 takes one throwaway filler drop to
	// realign whose turn it is before the column 2 and column 3 stacks.
	drops := []int{0, 1, 1, 2, 2, 4, 2, 3, 3, 3, 3}
	for _, col := range drops {
		if _, err := b.Drop(col); err != nil {
			t.Fatalf("Drop(%d) returned error %v, want nil", col, err)
		}
	}

	wantDiagonal := [][2]int{{5, 0}, {4, 1}, {3, 2}, {2, 3}}
	for _, cell := range wantDiagonal {
		if got := b.Cell(cell[0], cell[1]); got != Player1 {
			t.Fatalf("Cell(%d, %d) = %v, want Player1", cell[0], cell[1], got)
		}
	}

	if got := b.Winner(); got != Player1 {
		t.Errorf("Winner() = %v, want Player1", got)
	}
}

func TestWinnerDiagonalDown(t *testing.T) {
	b := New()

	// Drop order chosen so Player1's turn lands on (5,3), (4,2), (3,1), and
	// (2,0), a falling diagonal. Column 4 takes one throwaway filler drop to
	// realign whose turn it is before the column 1 and column 0 stacks.
	drops := []int{3, 2, 2, 1, 1, 4, 1, 0, 0, 0, 0}
	for _, col := range drops {
		if _, err := b.Drop(col); err != nil {
			t.Fatalf("Drop(%d) returned error %v, want nil", col, err)
		}
	}

	wantDiagonal := [][2]int{{5, 3}, {4, 2}, {3, 1}, {2, 0}}
	for _, cell := range wantDiagonal {
		if got := b.Cell(cell[0], cell[1]); got != Player1 {
			t.Fatalf("Cell(%d, %d) = %v, want Player1", cell[0], cell[1], got)
		}
	}

	if got := b.Winner(); got != Player1 {
		t.Errorf("Winner() = %v, want Player1", got)
	}
}

func TestFull(t *testing.T) {
	b := New()
	if b.Full() {
		t.Fatalf("Full() on empty board = true, want false")
	}

	for col := 0; col < Cols; col++ {
		for row := 0; row < Rows; row++ {
			if _, err := b.Drop(col); err != nil {
				t.Fatalf("Drop(%d) returned error %v, want nil", col, err)
			}
		}
	}

	if !b.Full() {
		t.Errorf("Full() on filled board = false, want true")
	}
}

func TestDropDoesNotAdvancePlayerOnError(t *testing.T) {
	b := New()

	if _, err := b.Drop(-1); err == nil {
		t.Fatal("Drop(-1) returned nil error, want ErrInvalidColumn")
	}
	if got := b.CurrentPlayer(); got != Player1 {
		t.Errorf("CurrentPlayer() after invalid drop = %v, want unchanged Player1", got)
	}

	for i := 0; i < Rows; i++ {
		if _, err := b.Drop(1); err != nil {
			t.Fatalf("Drop(1) fill call %d returned error %v, want nil", i, err)
		}
	}
	before := b.CurrentPlayer()
	if _, err := b.Drop(1); err != ErrColumnFull {
		t.Fatalf("Drop(1) on full column error = %v, want ErrColumnFull", err)
	}
	if got := b.CurrentPlayer(); got != before {
		t.Errorf("CurrentPlayer() after full-column drop = %v, want unchanged %v", got, before)
	}
}
