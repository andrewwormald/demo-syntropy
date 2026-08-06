package framediff

import (
	"reflect"
	"testing"
)

func TestDiff(t *testing.T) {
	tests := []struct {
		name string
		prev []string
		next []string
		want []LineUpdate
	}{
		{
			name: "identical frames",
			prev: []string{"a", "b", "c"},
			next: []string{"a", "b", "c"},
			want: nil,
		},
		{
			name: "single changed line",
			prev: []string{"a", "b", "c"},
			next: []string{"a", "X", "c"},
			want: []LineUpdate{{Line: 1, Text: "X"}},
		},
		{
			name: "multiple changed lines",
			prev: []string{"a", "b", "c", "d"},
			next: []string{"X", "b", "Y", "Z"},
			want: []LineUpdate{
				{Line: 0, Text: "X"},
				{Line: 2, Text: "Y"},
				{Line: 3, Text: "Z"},
			},
		},
		{
			name: "no previous frame",
			prev: nil,
			next: []string{"a", "b", "c"},
			want: []LineUpdate{
				{Line: 0, Text: "a"},
				{Line: 1, Text: "b"},
				{Line: 2, Text: "c"},
			},
		},
		{
			name: "no previous frame with empty lines",
			prev: nil,
			next: []string{"a", "", "c"},
			want: []LineUpdate{
				{Line: 0, Text: "a"},
				{Line: 2, Text: "c"},
			},
		},
		{
			name: "next shorter than prev",
			prev: []string{"a", "b", "c"},
			next: []string{"a", "X"},
			want: []LineUpdate{{Line: 1, Text: "X"}},
		},
		{
			name: "both empty",
			prev: nil,
			next: nil,
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Diff(tt.prev, tt.next)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Diff(%q, %q) = %v, want %v", tt.prev, tt.next, got, tt.want)
			}
		})
	}
}
