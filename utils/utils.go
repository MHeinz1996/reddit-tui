package utils

import (
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

func NormalizeSubreddit(subreddit string) string {
	if subreddit == "reddit.com" {
		return subreddit
	}

	if len(subreddit) >= 2 && subreddit[:2] == "r/" {
		return subreddit
	}

	return fmt.Sprintf("r/%s", subreddit)
}

// TruncateString shortens s so it occupies at most w terminal cells, appending
// an ellipsis when it has to cut. Width is measured in display cells rather
// than bytes so multibyte text (CJK, accents, emoji) is never sliced mid-rune,
// which would emit invalid UTF-8 and corrupt the terminal's width accounting.
func TruncateString(s string, w int) string {
	if w <= 0 {
		return s
	} else if ansi.StringWidth(s) <= w || w <= 3 {
		return s
	}

	return ansi.Truncate(s, w, "...")
}

func Clamp(min, max, val int) int {
	if val < min {
		return min
	} else if val > max {
		return max
	}

	return val
}

func GetSingularPlural(s, singular, plural string) string {
	if s == "1" {
		return fmt.Sprintf("%s %s", s, singular)
	}

	return fmt.Sprintf("%s %s", s, plural)
}
