package greet_test

import (
	"testing"

	"github.com/masonbesmer/dark-factory-lab/internal/greet"
)

func TestGreet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want string
	}{
		{"named", "Ada", "Hello, Ada!"},
		{"empty falls back to world", "", "Hello, world!"},
		{"spaces fall back to world", "   ", "Hello, world!"},
		{"tabs and newlines fall back to world", "\t\n ", "Hello, world!"},
		{"surrounding whitespace is trimmed", "  Grace  ", "Hello, Grace!"},
		{"inner whitespace is kept", "Ada Lovelace", "Hello, Ada Lovelace!"},
		{"unicode", "Zoë", "Hello, Zoë!"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := greet.Greet(tc.in); got != tc.want {
				t.Errorf("Greet(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
