package utils

import (
	"fmt"
	"strings"
)

// AuthorWithFlair formats an author name and optional user flair for display.
func AuthorWithFlair(author, flair string) string {
	author = strings.TrimSpace(author)
	flair = strings.TrimSpace(flair)
	if author == "" {
		return flair
	}
	if flair == "" {
		return author
	}
	return fmt.Sprintf("%s [%s]", author, flair)
}
