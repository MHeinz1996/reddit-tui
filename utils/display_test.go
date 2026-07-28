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
