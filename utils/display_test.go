package utils

import (
	"strings"
	"testing"
)

func TestSanitizeDisplayTextRemovesZeroWidth(t *testing.T) {
	const zwsp = "\u200b"
	in := "Jhon Terry 😂" + zwsp + "😂"
	got := SanitizeDisplayText(in)
	if strings.Contains(got, zwsp) {
		t.Fatalf("zero-width space still present: %q", got)
	}
}
