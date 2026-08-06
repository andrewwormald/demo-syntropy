package input

import "testing"

func TestDecode(t *testing.T) {
	tests := []struct {
		name         string
		buf          []byte
		wantKey      Key
		wantConsumed int
	}{
		{"left arrow", []byte{0x1b, '[', 'D'}, Left, 3},
		{"right arrow", []byte{0x1b, '[', 'C'}, Right, 3},
		{"enter carriage return", []byte{'\r'}, Enter, 1},
		{"enter line feed", []byte{'\n'}, Enter, 1},
		{"unrecognized letter", []byte{'a'}, Unknown, 1},
		{"lone escape byte", []byte{0x1b}, Unknown, 1},
		{"escape without bracket", []byte{0x1b, 'x', 'x'}, Unknown, 1},
		{"escape bracket unknown final byte", []byte{0x1b, '[', 'Z'}, Unknown, 1},
		{"empty buffer", []byte{}, Unknown, 0},
		{"left arrow followed by more input", []byte{0x1b, '[', 'D', 'a'}, Left, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, gotConsumed := Decode(tt.buf)
			if gotKey != tt.wantKey {
				t.Errorf("Decode(%v) key = %v, want %v", tt.buf, gotKey, tt.wantKey)
			}
			if gotConsumed != tt.wantConsumed {
				t.Errorf("Decode(%v) consumed = %d, want %d", tt.buf, gotConsumed, tt.wantConsumed)
			}
		})
	}
}
