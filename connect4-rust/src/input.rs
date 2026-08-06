//! Decodes raw terminal byte sequences (as read from a terminal put into
//! raw mode) into game key events, independent of any game logic.

/// A decoded key press.
#[derive(Debug, PartialEq, Eq, Clone, Copy)]
pub enum Key {
    /// Returned for byte sequences that don't map to a recognized key.
    Unknown,
    Left,
    Right,
    Enter,
    Quit,
}

/// Inspects a raw byte sequence read from a terminal in raw mode and
/// returns the `Key` it represents along with the number of bytes consumed
/// from `buf`. Left and Right correspond to the ANSI escape sequences
/// emitted by the arrow keys (ESC [ D and ESC [ C respectively); Enter
/// corresponds to a carriage return or line feed; Quit corresponds to
/// Ctrl-C (which raw mode delivers as the literal byte 0x03 instead of
/// SIGINT) or the 'q'/'Q' keys. If `buf` does not start with a recognized
/// sequence, `decode` returns `Unknown` and consumes a single byte so
/// callers can skip past it.
pub fn decode(buf: &[u8]) -> (Key, usize) {
    if buf.is_empty() {
        return (Key::Unknown, 0);
    }

    match buf[0] {
        b'\r' | b'\n' => (Key::Enter, 1),
        0x03 | b'q' | b'Q' => (Key::Quit, 1),
        0x1b => {
            if buf.len() >= 3 && buf[1] == b'[' {
                match buf[2] {
                    b'C' => return (Key::Right, 3),
                    b'D' => return (Key::Left, 3),
                    _ => {}
                }
            }
            (Key::Unknown, 1)
        }
        _ => (Key::Unknown, 1),
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn decodes_left_arrow() {
        assert_eq!(decode(&[0x1b, b'[', b'D']), (Key::Left, 3));
    }

    #[test]
    fn decodes_right_arrow() {
        assert_eq!(decode(&[0x1b, b'[', b'C']), (Key::Right, 3));
    }

    #[test]
    fn decodes_enter_carriage_return() {
        assert_eq!(decode(&[b'\r']), (Key::Enter, 1));
    }

    #[test]
    fn decodes_enter_line_feed() {
        assert_eq!(decode(&[b'\n']), (Key::Enter, 1));
    }

    #[test]
    fn decodes_quit_ctrl_c() {
        assert_eq!(decode(&[0x03]), (Key::Quit, 1));
    }

    #[test]
    fn decodes_quit_lowercase_q() {
        assert_eq!(decode(&[b'q']), (Key::Quit, 1));
    }

    #[test]
    fn decodes_quit_uppercase_q() {
        assert_eq!(decode(&[b'Q']), (Key::Quit, 1));
    }

    #[test]
    fn unrecognized_letter_is_unknown() {
        assert_eq!(decode(&[b'a']), (Key::Unknown, 1));
    }

    #[test]
    fn lone_escape_byte_is_unknown() {
        assert_eq!(decode(&[0x1b]), (Key::Unknown, 1));
    }

    #[test]
    fn escape_without_bracket_is_unknown() {
        assert_eq!(decode(&[0x1b, b'x', b'x']), (Key::Unknown, 1));
    }

    #[test]
    fn escape_bracket_unknown_final_byte_is_unknown() {
        assert_eq!(decode(&[0x1b, b'[', b'Z']), (Key::Unknown, 1));
    }

    #[test]
    fn empty_buffer_is_unknown() {
        assert_eq!(decode(&[]), (Key::Unknown, 0));
    }

    #[test]
    fn left_arrow_followed_by_more_input() {
        assert_eq!(decode(&[0x1b, b'[', b'D', b'a']), (Key::Left, 3));
    }
}
