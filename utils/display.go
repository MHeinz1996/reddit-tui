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
//
// Emoji skin tone modifiers (U+1F3FB-U+1F3FF) are also stripped, so emoji
// render in their default tone. Terminals whose font cannot fuse a modifier
// into its base emoji draw it as a separate box, taking two more cells than
// Go measures and shifting the rest of the line. The base emoji is kept, and
// ZWJ sequences stay valid once their modifiers are gone.
func SanitizeDisplayText(s string) string {
	return strings.Map(func(r rune) rune {
		if isSkinToneModifier(r) {
			return -1
		}

		switch r {
		case '\u200b', '\u200c', '\ufeff':
			return -1
		default:
			return r
		}
	}, s)
}

// isSkinToneModifier reports whether r is one of the five Fitzpatrick emoji
// modifiers.
func isSkinToneModifier(r rune) bool {
	return r >= '\U0001f3fb' && r <= '\U0001f3ff'
}
