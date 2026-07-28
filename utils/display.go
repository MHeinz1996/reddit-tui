package utils

import "strings"

// SanitizeDisplayText removes invisible Unicode that breaks terminal width
// calculations (e.g. zero-width spaces Reddit inserts between emoji).
//
// Only truly zero-width formatting characters are stripped. U+FE0F
// (variation selector-16) and U+200D (zero width joiner) are deliberately
// preserved: they are load-bearing. VS16 promotes a text-presentation glyph
// such as U+26A0 to its double-width emoji form, and ZWJ fuses a multi-emoji
// sequence into a single glyph. Removing either makes Go measure the text as
// narrower (or, for ZWJ, far wider) than the terminal actually draws it, so
// lines survive truncation and then soft-wrap, which desynchronizes Bubble
// Tea's line accounting and duplicates content on screen.
func SanitizeDisplayText(s string) string {
	return strings.Map(func(r rune) rune {
		switch r {
		case '\u200b', '\u200c', '\ufeff':
			return -1
		default:
			return r
		}
	}, s)
}
