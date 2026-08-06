//! Connect 4 grid state and drop logic.

use std::error::Error;
use std::fmt;

pub const ROWS: usize = 6;
pub const COLS: usize = 7;

/// The contents of a single grid position.
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum Cell {
    Empty,
    Player1,
    Player2,
}

/// An error returned by [`Board::drop_piece`].
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum DropError {
    InvalidColumn,
    ColumnFull,
}

impl fmt::Display for DropError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            DropError::InvalidColumn => write!(f, "board: column out of range"),
            DropError::ColumnFull => write!(f, "board: column is full"),
        }
    }
}

impl Error for DropError {}

/// The Connect 4 grid state and whose turn it is.
pub struct Board {
    grid: [[Cell; COLS]; ROWS],
    current: Cell,
}

/// The four axes along which four-in-a-row can occur: horizontal, vertical,
/// and both diagonals. Each is checked in both senses by scanning outward
/// from every occupied cell.
const WIN_DIRECTIONS: [(isize, isize); 4] = [(0, 1), (1, 0), (1, 1), (1, -1)];

impl Board {
    /// Returns an empty board with Player1 to move first.
    pub fn new() -> Board {
        Board {
            grid: [[Cell::Empty; COLS]; ROWS],
            current: Cell::Player1,
        }
    }

    /// Returns the player whose turn it is to drop next.
    pub fn current_player(&self) -> Cell {
        self.current
    }

    /// Returns the contents of the grid at (row, col).
    pub fn cell(&self, row: usize, col: usize) -> Cell {
        self.grid[row][col]
    }

    /// Places the current player's piece into the given column, settling it
    /// to the lowest empty row. Returns the row the piece landed in and
    /// advances the current player. Returns `DropError::InvalidColumn` if
    /// `col` is out of range, and `DropError::ColumnFull` if the column has
    /// no empty rows.
    pub fn drop_piece(&mut self, col: i32) -> Result<usize, DropError> {
        if col < 0 || col as usize >= COLS {
            return Err(DropError::InvalidColumn);
        }
        let col = col as usize;

        for row in (0..ROWS).rev() {
            if self.grid[row][col] == Cell::Empty {
                self.grid[row][col] = self.current;
                self.current = match self.current {
                    Cell::Player1 => Cell::Player2,
                    _ => Cell::Player1,
                };
                return Ok(row);
            }
        }

        Err(DropError::ColumnFull)
    }

    /// Returns the player with four connected pieces in a row, column, or
    /// diagonal, or `Cell::Empty` if no player has won yet.
    pub fn winner(&self) -> Cell {
        for row in 0..ROWS {
            for col in 0..COLS {
                let player = self.grid[row][col];
                if player == Cell::Empty {
                    continue;
                }

                for &(dr, dc) in &WIN_DIRECTIONS {
                    let mut count = 1;
                    let mut r = row as isize + dr;
                    let mut c = col as isize + dc;
                    while count < 4
                        && r >= 0
                        && r < ROWS as isize
                        && c >= 0
                        && c < COLS as isize
                        && self.grid[r as usize][c as usize] == player
                    {
                        count += 1;
                        r += dr;
                        c += dc;
                    }
                    if count == 4 {
                        return player;
                    }
                }
            }
        }
        Cell::Empty
    }

    /// Reports whether every column is filled, meaning no more moves are
    /// possible.
    pub fn full(&self) -> bool {
        (0..COLS).all(|col| self.grid[0][col] != Cell::Empty)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn new_board_is_empty_with_player1_first() {
        let b = Board::new();

        assert_eq!(b.current_player(), Cell::Player1);
        for row in 0..ROWS {
            for col in 0..COLS {
                assert_eq!(b.cell(row, col), Cell::Empty);
            }
        }
    }

    #[test]
    fn drop_into_empty_column() {
        let mut b = Board::new();

        let row = b.drop_piece(3).expect("drop(3) should succeed");
        assert_eq!(row, ROWS - 1);
        assert_eq!(b.cell(row, 3), Cell::Player1);
        assert_eq!(b.current_player(), Cell::Player2);
    }

    #[test]
    fn drop_into_partially_filled_column() {
        let mut b = Board::new();

        let want_players = [Cell::Player1, Cell::Player2, Cell::Player1];
        for (i, want) in want_players.iter().enumerate() {
            let row = b.drop_piece(0).expect("drop(0) should succeed");
            assert_eq!(row, ROWS - 1 - i);
            assert_eq!(b.cell(row, 0), *want);
        }
    }

    #[test]
    fn drop_into_full_column() {
        let mut b = Board::new();

        for _ in 0..ROWS {
            b.drop_piece(2).expect("filling drop(2) should succeed");
        }

        let err = b.drop_piece(2).expect_err("drop(2) on full column should fail");
        assert_eq!(err, DropError::ColumnFull);
    }

    #[test]
    fn drop_invalid_column() {
        for &col in &[-1, COLS as i32, 100] {
            let mut b = Board::new();
            let err = b.drop_piece(col).expect_err("drop should fail for out-of-range column");
            assert_eq!(err, DropError::InvalidColumn);
        }
    }

    #[test]
    fn winner_none_on_empty_board() {
        let b = Board::new();
        assert_eq!(b.winner(), Cell::Empty);
    }

    #[test]
    fn winner_horizontal() {
        let mut b = Board::new();

        // Player1 drops in columns 0-3, Player2 drops elsewhere in between so
        // the row stays Player1's.
        for col in 0..4 {
            b.drop_piece(col).expect("drop should succeed");
            if col < 3 {
                b.drop_piece(col).expect("drop should succeed");
            }
        }

        assert_eq!(b.winner(), Cell::Player1);
    }

    #[test]
    fn winner_vertical() {
        let mut b = Board::new();

        for i in 0..4 {
            b.drop_piece(0).expect("drop should succeed");
            if i < 3 {
                b.drop_piece(1).expect("drop should succeed");
            }
        }

        assert_eq!(b.winner(), Cell::Player1);
    }

    #[test]
    fn winner_diagonal_up() {
        let mut b = Board::new();

        // Drop order chosen so Player1's turn lands on (5,0), (4,1), (3,2),
        // and (2,3), a rising diagonal. Column 4 takes one throwaway filler
        // drop to realign whose turn it is before the column 2 and column 3
        // stacks.
        let drops = [0, 1, 1, 2, 2, 4, 2, 3, 3, 3, 3];
        for col in drops {
            b.drop_piece(col).expect("drop should succeed");
        }

        let want_diagonal = [(5, 0), (4, 1), (3, 2), (2, 3)];
        for (row, col) in want_diagonal {
            assert_eq!(b.cell(row, col), Cell::Player1);
        }

        assert_eq!(b.winner(), Cell::Player1);
    }

    #[test]
    fn winner_diagonal_down() {
        let mut b = Board::new();

        // Drop order chosen so Player1's turn lands on (5,3), (4,2), (3,1),
        // and (2,0), a falling diagonal. Column 4 takes one throwaway filler
        // drop to realign whose turn it is before the column 1 and column 0
        // stacks.
        let drops = [3, 2, 2, 1, 1, 4, 1, 0, 0, 0, 0];
        for col in drops {
            b.drop_piece(col).expect("drop should succeed");
        }

        let want_diagonal = [(5, 3), (4, 2), (3, 1), (2, 0)];
        for (row, col) in want_diagonal {
            assert_eq!(b.cell(row, col), Cell::Player1);
        }

        assert_eq!(b.winner(), Cell::Player1);
    }

    #[test]
    fn full() {
        let mut b = Board::new();
        assert!(!b.full());

        for col in 0..COLS as i32 {
            for _ in 0..ROWS {
                b.drop_piece(col).expect("drop should succeed");
            }
        }

        assert!(b.full());
    }

    #[test]
    fn drop_does_not_advance_player_on_error() {
        let mut b = Board::new();

        b.drop_piece(-1).expect_err("drop(-1) should fail");
        assert_eq!(b.current_player(), Cell::Player1);

        for _ in 0..ROWS {
            b.drop_piece(1).expect("filling drop(1) should succeed");
        }
        let before = b.current_player();
        let err = b.drop_piece(1).expect_err("drop(1) on full column should fail");
        assert_eq!(err, DropError::ColumnFull);
        assert_eq!(b.current_player(), before);
    }
}
