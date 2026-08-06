//! Wires the board, input, and renderer modules into an interactive
//! Connect 4 game loop.

use std::io::{self, Read, Write};

use crate::board::{Board, Cell, COLS};
use crate::input::{self, Key};
use crate::renderer::render;

/// Holds the in-progress board state and the column currently selected for
/// the next drop.
pub struct Game {
    board: Board,
    cursor: usize,
}

impl Game {
    /// Returns a Game with an empty board and the cursor on the leftmost
    /// column.
    pub fn new() -> Self {
        Game {
            board: Board::new(),
            cursor: 0,
        }
    }

    /// Applies a decoded key event to the game state. Left and Right move
    /// the cursor, clamped to the board's columns. Enter drops the current
    /// player's piece into the selected column; a full column is silently
    /// ignored, matching how a player would just try another column.
    /// Returns whether a piece was dropped.
    pub fn handle_key(&mut self, key: Key) -> bool {
        match key {
            Key::Left => {
                if self.cursor > 0 {
                    self.cursor -= 1;
                }
            }
            Key::Right => {
                if self.cursor < COLS - 1 {
                    self.cursor += 1;
                }
            }
            Key::Enter => {
                if self.board.drop(self.cursor).is_ok() {
                    return true;
                }
            }
            Key::Unknown | Key::Quit => {}
        }
        false
    }
}

impl Default for Game {
    fn default() -> Self {
        Self::new()
    }
}

/// Returns the end-of-game message for the current board state, or `None`
/// if the game is still in progress.
fn announcement(board: &Board) -> Option<&'static str> {
    match board.winner() {
        Cell::Player1 => return Some("Player 1 wins!\n"),
        Cell::Player2 => return Some("Player 2 wins!\n"),
        Cell::Empty => {}
    }
    if board.full() {
        return Some("Draw!\n");
    }
    None
}

/// Decodes raw terminal byte sequences from `r` and drives the game against
/// them, writing the rendered board to `w` after every recognized key. It
/// returns once a player wins, the board fills up, a Quit key is decoded, or
/// `r` is exhausted, writing an end-of-game announcement in the first two
/// cases.
pub fn run<R: Read, W: Write>(mut r: R, mut w: W) -> io::Result<()> {
    let mut g = Game::new();

    w.write_all(render(&g.board, g.cursor).as_bytes())?;

    let mut pending: Vec<u8> = Vec::new();
    let mut chunk = [0u8; 64];
    loop {
        let n = r.read(&mut chunk)?;
        if n > 0 {
            pending.extend_from_slice(&chunk[..n]);
        }

        loop {
            let (key, consumed) = input::decode(&pending);
            if consumed == 0 {
                break;
            }
            pending.drain(..consumed);

            if key == Key::Unknown {
                continue;
            }
            if key == Key::Quit {
                return Ok(());
            }

            let dropped = g.handle_key(key);
            w.write_all(render(&g.board, g.cursor).as_bytes())?;
            if dropped {
                if let Some(msg) = announcement(&g.board) {
                    w.write_all(msg.as_bytes())?;
                    return Ok(());
                }
            }
        }

        if n == 0 {
            return Ok(());
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::board::{Cell, ROWS};

    #[test]
    fn handle_key_moves_cursor_clamped() {
        let mut g = Game::new();

        g.handle_key(Key::Left);
        assert_eq!(g.cursor, 0);

        g.handle_key(Key::Right);
        g.handle_key(Key::Right);
        assert_eq!(g.cursor, 2);

        for _ in 0..COLS {
            g.handle_key(Key::Right);
        }
        assert_eq!(g.cursor, COLS - 1);
    }

    #[test]
    fn handle_key_enter_drops_piece() {
        let mut g = Game::new();
        g.cursor = 3;

        let dropped = g.handle_key(Key::Enter);
        assert!(dropped, "handle_key(Enter) = false, want true");
        assert_eq!(g.board.cell(ROWS - 1, 3), Cell::Player1);
    }

    #[test]
    fn handle_key_enter_on_full_column_does_not_drop() {
        let mut g = Game::new();
        for _ in 0..ROWS {
            g.handle_key(Key::Enter);
        }

        let dropped = g.handle_key(Key::Enter);
        assert!(!dropped, "handle_key(Enter) on full column = true, want false");
    }

    #[test]
    fn run_plays_moves_from_input() {
        // Right, Right, Enter drops Player1 into column 2.
        let input = [0x1b, b'[', b'C', 0x1b, b'[', b'C', b'\r'];
        let mut out = Vec::new();

        run(&input[..], &mut out).expect("run() should not error");

        let rendered = String::from_utf8(out).expect("output should be valid UTF-8");
        assert!(
            rendered.contains('X'),
            "run() output = {rendered:?}, want it to contain a dropped piece"
        );
    }

    #[test]
    fn run_stops_on_quit() {
        // Right, then Ctrl-C: the cursor move should be rendered but the
        // game should stop before any further input is read.
        let input = [0x1b, b'[', b'C', 0x03, 0x1b, b'[', b'C'];
        let mut out = Vec::new();

        run(&input[..], &mut out).expect("run() should not error");

        assert!(!out.is_empty(), "run() wrote no output");
    }

    #[test]
    fn run_stops_on_eof() {
        let input = [b'\r'];
        let mut out = Vec::new();

        run(&input[..], &mut out).expect("run() should not error");

        assert!(!out.is_empty(), "run() wrote no output");
    }

    #[test]
    fn run_announces_winner() {
        // Alternating drops into columns 0 and 1 land Player1's pieces at
        // every row of column 0 (Player2 always drops into column 1 in
        // between), connecting four for Player1 vertically.
        let enter = b'\r';
        let right = [0x1b, b'[', b'C'];
        let left = [0x1b, b'[', b'D'];
        let mut keys = Vec::new();
        for i in 0..4 {
            keys.push(enter);
            if i < 3 {
                keys.extend_from_slice(&right);
                keys.push(enter);
                keys.extend_from_slice(&left);
            }
        }

        let mut out = Vec::new();
        run(&keys[..], &mut out).expect("run() should not error");

        let rendered = String::from_utf8(out).expect("output should be valid UTF-8");
        assert!(
            rendered.contains("Player 1 wins!"),
            "run() output = {rendered:?}, want it to contain a win announcement"
        );
    }

    #[test]
    fn announcement_draw() {
        // Column order that fills the entire board without ever connecting
        // four, found by randomized search over valid drop sequences.
        // Applied directly through the board rather than via run's raw key
        // decoding, since driving 42 drops across all 7 columns would need
        // far more arrow-key bytes than fit in run's read buffer.
        let draw_cols = [
            5, 3, 6, 6, 4, 4, 5, 6, 2, 4, 0, 3, 2, 1, 5, 1, 6, 5, 6, 5, 4, 1, 2, 6, 0, 4, 4, 0, 1,
            2, 5, 3, 1, 2, 3, 2, 3, 0, 3, 1, 0, 0,
        ];

        let mut b = Board::new();
        for col in draw_cols {
            b.drop(col).unwrap_or_else(|e| panic!("drop({col}) returned {e:?}"));
        }

        assert_eq!(announcement(&b), Some("Draw!\n"));
    }
}
