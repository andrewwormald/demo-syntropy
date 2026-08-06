// Package input decodes raw terminal byte sequences (as read from a
// terminal put into raw mode) into game key events, independent of any
// game logic.
package input

// Key identifies a decoded key press.
type Key int

const (
	// Unknown is returned for byte sequences that don't map to a
	// recognized key.
	Unknown Key = iota
	Left
	Right
	Enter
)

// Decode inspects a raw byte sequence read from a terminal in raw mode and
// returns the Key it represents along with the number of bytes consumed
// from buf. Left and Right correspond to the ANSI escape sequences emitted
// by the arrow keys (ESC [ D and ESC [ C respectively); Enter corresponds
// to a carriage return or line feed. If buf does not start with a
// recognized sequence, Decode returns Unknown and consumes a single byte
// so callers can skip past it.
func Decode(buf []byte) (Key, int) {
	if len(buf) == 0 {
		return Unknown, 0
	}

	switch buf[0] {
	case '\r', '\n':
		return Enter, 1
	case 0x1b: // ESC
		if len(buf) >= 3 && buf[1] == '[' {
			switch buf[2] {
			case 'C':
				return Right, 3
			case 'D':
				return Left, 3
			}
		}
		return Unknown, 1
	default:
		return Unknown, 1
	}
}
