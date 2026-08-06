package screen

import (
	"bytes"
	"testing"
)

func TestRender(t *testing.T) {
	tests := []struct {
		name   string
		frames [][]string
		want   []string
	}{
		{
			name:   "first frame draws every non-empty line",
			frames: [][]string{{"a", "b", "c"}},
			want:   []string{"\r\x1b[Ka\x1b[1B\r\x1b[Kb\x1b[1B\r\x1b[Kc\x1b[1B\r"},
		},
		{
			name:   "identical frame writes nothing",
			frames: [][]string{{"a", "b", "c"}, {"a", "b", "c"}},
			want:   []string{"\r\x1b[Ka\x1b[1B\r\x1b[Kb\x1b[1B\r\x1b[Kc\x1b[1B\r", ""},
		},
		{
			name:   "single changed line",
			frames: [][]string{{"a", "b", "c"}, {"a", "X", "c"}},
			want: []string{
				"\r\x1b[Ka\x1b[1B\r\x1b[Kb\x1b[1B\r\x1b[Kc\x1b[1B\r",
				"\r\x1b[3A\x1b[1B\r\x1b[KX\x1b[2B\r",
			},
		},
		{
			name:   "shrinking frame leaves cursor after the new last line",
			frames: [][]string{{"a", "b", "c"}, {"a", "X"}},
			want: []string{
				"\r\x1b[Ka\x1b[1B\r\x1b[Kb\x1b[1B\r\x1b[Kc\x1b[1B\r",
				"\r\x1b[3A\x1b[1B\r\x1b[KX\x1b[1B\r",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			s := New(&buf)
			for i, frame := range tt.frames {
				buf.Reset()
				if err := s.Render(frame); err != nil {
					t.Fatalf("Render(%q) error = %v", frame, err)
				}
				if got := buf.String(); got != tt.want[i] {
					t.Errorf("Render(%q) wrote %q, want %q", frame, got, tt.want[i])
				}
			}
		})
	}
}
