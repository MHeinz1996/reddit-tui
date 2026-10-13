package utils

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

func TestSanitizeDisplayTextRemovesZeroWidth(t *testing.T) {
	const zwsp = "\u200b"
	in := "Jhon Terry 😂" + zwsp + "😂"
	got := SanitizeDisplayText(in)
	if strings.Contains(got, zwsp) {
		t.Fatalf("zero-width space still present: %q", got)
	}
}

// Sanitizing must never change how wide the text renders. Stripping VS16 or ZWJ
// made Go measure emoji differently than the terminal draws them, so lines
// passed truncation and then soft-wrapped, duplicating content on screen.
func TestSanitizeDisplayTextPreservesWidth(t *testing.T) {
	inputs := []string{
		"⚠️ PSA about this thing",
		"❤️ a post title",
		"ℹ️ info",
		"➡️ next",
		"☠️ danger",
		"\U0001f468‍\U0001f469‍\U0001f467 family",
		"\U0001f469‍❤️‍\U0001f468 couple",
		"\U0001f3f3️‍\U0001f308 pride",
		"\U0001f44d\U0001f3fd thumbs",
		"\U0001f1fa\U0001f1f8 flag",
	}

	for _, in := range inputs {
		got := SanitizeDisplayText(in)
		if want, have := ansi.StringWidth(in), ansi.StringWidth(got); want != have {
			t.Errorf("SanitizeDisplayText(%q) changed display width: %d -> %d (result %q)",
				in, want, have, got)
		}
	}
}

func TestSanitizeDisplayTextStripsSkinTones(t *testing.T) {
	const zwj = "‍"

	tests := []struct {
		name string
		in   string
		want string
	}{
		{"pointing up", "\U0001f446\U0001f3fd100% this", "\U0001f446100% this"},
		{"every tone", "\U0001f44d\U0001f3fb\U0001f44d\U0001f3fc\U0001f44d\U0001f3fd\U0001f44d\U0001f3fe\U0001f44d\U0001f3ff", "\U0001f44d\U0001f44d\U0001f44d\U0001f44d\U0001f44d"},
		{"zwj sequence", "\U0001f469\U0001f3fd" + zwj + "\U0001f680", "\U0001f469" + zwj + "\U0001f680"},
		{"untoned emoji untouched", "\U0001f446 plain", "\U0001f446 plain"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SanitizeDisplayText(tt.in); got != tt.want {
				t.Errorf("SanitizeDisplayText(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestSanitizeDisplayTextKeepsLoadBearingJoiners(t *testing.T) {
	const (
		vs16 = "️"
		zwj  = "‍"
	)

	if got := SanitizeDisplayText("⚠" + vs16); !strings.Contains(got, vs16) {
		t.Errorf("variation selector-16 was stripped: %q", got)
	}

	family := "\U0001f468" + zwj + "\U0001f469" + zwj + "\U0001f467"
	if got := SanitizeDisplayText(family); !strings.Contains(got, zwj) {
		t.Errorf("zero width joiner was stripped: %q", got)
	}
}
