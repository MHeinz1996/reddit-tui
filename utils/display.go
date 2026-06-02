package utils

import "strings"

// SanitizeDisplayText removes invisible Unicode that breaks terminal width
// calculations (e.g. zero-width spaces Reddit inserts between emoji).
func SanitizeDisplayText(s string) string {
	return strings.Map(func(r rune) rune {
		switch r {
		case '\u200b', '\u200c', '\u200d', '\ufeff', '\ufe0f':
			return -1
		default:
			return r
		}
	}, s)
}
