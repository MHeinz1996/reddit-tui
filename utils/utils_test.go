package utils

import (
	"testing"
	"unicode/utf8"

	"github.com/charmbracelet/x/ansi"
)

func TestNormalizeSubreddit(t *testing.T) {
	tests := []struct {
		subreddit string
		want      string
	}{
		{"neovim", "r/neovim"},
		{"r/neovim", "r/neovim"},
	}

	for _, tt := range tests {
		got := NormalizeSubreddit(tt.subreddit)
		if got != tt.want {
			t.Errorf("got %s, want %s", got, tt.want)
		}
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		s     string
		width int
		want  string
	}{
		{"abc", 3, "abc"},
		{"abcd", 4, "abcd"},
		{"abcde", 4, "a..."},
		{"abcdef", 5, "ab..."},
		{"abcdefg", 6, "abc..."},

		// Multibyte text must be cut on rune boundaries, never mid-sequence.
		{"r/Ünicöde", 6, "r/Ü..."},
		{"r/café_lovers", 8, "r/caf..."},
		// Wide (2-cell) runes: the budget is columns, not runes or bytes.
		{"r/日本語のサブレディット", 10, "r/日本..."},
	}

	for _, tt := range tests {
		got := TruncateString(tt.s, tt.width)
		if got != tt.want {
			t.Errorf("got %s, want %s with input %s", got, tt.want, tt.s)
		}
	}
}

// Slicing by byte offset used to split multibyte runes, producing invalid UTF-8
// that the terminal cannot measure, which corrupted the surrounding redraw.
func TestTruncateStringKeepsValidUtf8(t *testing.T) {
	inputs := []string{
		"r/日本語のサブレディット",
		"r/Ünicöde_subreddit_name",
		"r/café_lovers_united",
		"r/🚀🚀🚀🚀🚀🚀🚀🚀",
		"r/emoji⚠️warning⚠️signs",
	}

	for _, in := range inputs {
		for w := 1; w <= 20; w++ {
			got := TruncateString(in, w)
			if !utf8.ValidString(got) {
				t.Errorf("TruncateString(%q, %d) = %q, which is not valid UTF-8", in, w, got)
			}
			if w > 3 && ansi.StringWidth(got) > w {
				t.Errorf("TruncateString(%q, %d) = %q, width %d exceeds budget %d",
					in, w, got, ansi.StringWidth(got), w)
			}
		}
	}
}

func TestClamp(t *testing.T) {
	tests := []struct {
		min  int
		max  int
		val  int
		want int
	}{
		{0, 2, 0, 0},
		{0, 2, 1, 1},
		{0, 2, 2, 2},
		{0, 2, 3, 2},
		{0, 2, -1, 0},
		{0, 10, 5, 5},
		{0, 10, -5, 0},
		{0, 10, 15, 10},
	}

	for _, tt := range tests {
		got := Clamp(tt.min, tt.max, tt.val)
		if got != tt.want {
			t.Errorf("got %d, want %d with input: min %d, max %d, val %d", got, tt.want, tt.min, tt.max, tt.val)
		}
	}
}

func TestGetSingularPlural(t *testing.T) {
	tests := []struct {
		s        string
		singular string
		plural   string
		want     string
	}{
		{"0", "banana", "bananas", "0 bananas"},
		{"1", "banana", "bananas", "1 banana"},
		{"2", "banana", "bananas", "2 bananas"},
		{"3", "banana", "bananas", "3 bananas"},
	}

	for _, tt := range tests {
		got := GetSingularPlural(tt.s, tt.singular, tt.plural)
		if got != tt.want {
			t.Errorf("got %s, want %s with input: s %s, singular %s, plural %s", got, tt.want, tt.s, tt.singular, tt.plural)
		}
	}
}
