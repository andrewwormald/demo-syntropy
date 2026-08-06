//! ASCII rendering of a `Board`.

use crate::board::{Board, Cell, COLS, ROWS};

fn symbol(cell: Cell) -> char {
    match cell {
        Cell::Empty => '.',
        Cell::Player1 => 'X',
        Cell::Player2 => 'O',
    }
}

/// Returns an ASCII representation of the board's current state, with a
/// selector row above the grid marking `cursor` as the column the current
/// player will drop into, one row per grid row, and a border beneath the
/// bottom row.
pub fn render(b: &Board, cursor: usize) -> String {
    let mut out = String::new();
    for col in 0..COLS {
        out.push_str("| ");
        out.push(if col == cursor { 'v' } else { ' ' });
        out.push(' ');
    }
    out.push_str("|\n");

    for row in 0..ROWS {
        for col in 0..COLS {
            out.push_str("| ");
            out.push(symbol(b.cell(row, col)));
            out.push(' ');
        }
        out.push_str("|\n");
    }

    for _ in 0..COLS {
        out.push_str("----");
    }
    out.push_str("-\n");

    out
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn render_empty_board() {
        let b = Board::new();

        let got = render(&b, 0);

        let want_rows = ROWS + 2; // selector row + board rows + border row
        assert_eq!(got.matches('\n').count(), want_rows);
        assert!(!got.contains('X'));
        assert!(!got.contains('O'));
    }

    #[test]
    fn render_shows_dropped_pieces() {
        let mut b = Board::new();

        b.drop(3).expect("drop(3) should succeed");
        b.drop(3).expect("drop(3) should succeed");

        let got = render(&b, 0);

        let lines: Vec<&str> = got.trim_end_matches('\n').split('\n').collect();
        const SELECTOR_ROWS: usize = 1;
        let bottom_row = lines[SELECTOR_ROWS + ROWS - 1];
        let second_from_bottom_row = lines[SELECTOR_ROWS + ROWS - 2];

        assert!(
            bottom_row.contains('X'),
            "bottom row {bottom_row:?} does not contain Player1 symbol X"
        );
        assert!(
            second_from_bottom_row.contains('O'),
            "second from bottom row {second_from_bottom_row:?} does not contain Player2 symbol O"
        );
    }

    #[test]
    fn render_selector_marks_cursor_column() {
        let b = Board::new();

        for cursor in [0, 3, COLS - 1] {
            let got = render(&b, cursor);

            let lines: Vec<&str> = got.trim_end_matches('\n').split('\n').collect();
            let selector_row = lines[0];
            let cells: Vec<&str> = selector_row.trim_matches('|').split('|').collect();
            assert_eq!(
                cells.len(),
                COLS,
                "render(_, {cursor}) selector row {selector_row:?} has {} cells, want {COLS}",
                cells.len()
            );
            for (col, cell) in cells.iter().enumerate() {
                let marked = cell.contains('v');
                if col == cursor {
                    assert!(
                        marked,
                        "render(_, {cursor}) selector row {selector_row:?} does not mark column {col}"
                    );
                } else {
                    assert!(
                        !marked,
                        "render(_, {cursor}) selector row {selector_row:?} unexpectedly marks column {col}"
                    );
                }
            }
        }
    }
}
