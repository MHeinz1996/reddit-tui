package comments

import (
	"fmt"
	"strings"
	"testing"

	"reddittui/components/styles"
	"reddittui/model"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// newTestPage builds a CommentsPage without a RedditClient. Layout is
// independent of the client, and constructing one would touch the filesystem.
func newTestPage(w, h int) CommentsPage {
	page := CommentsPage{
		header:         NewCommentsHeader(),
		pager:          NewCommentsViewport(),
		containerStyle: styles.GlobalStyle,
		sort:           model.CommentSortBest,
		focus:          true,
	}
	page.SetSize(w, h)
	return page
}

func nestedComments(maxDepth int, text string) model.Comments {
	var cs []model.Comment
	for d := 0; d <= maxDepth; d++ {
		cs = append(cs, model.Comment{
			Author: "someuser", Points: "12 points", Timestamp: "1 hour ago",
			Depth: d, Text: text,
		})
	}

	return model.Comments{
		PostTitle: "A post title", Subreddit: "golang", PostAuthor: "bob",
		PostTimestamp: "1 hour ago", PostPoints: "5", Comments: cs,
	}
}

// assertFits checks the invariant the terminal actually cares about: the view
// must never be wider or taller than the space it was given. A view that
// exceeds either bound gets soft-wrapped or scrolled by the terminal, which
// desynchronizes Bubble Tea's line accounting and duplicates content on screen.
func assertFits(t *testing.T, label, view string, w, h int) {
	t.Helper()

	if got := lipgloss.Height(view); got > h {
		t.Errorf("%s: view is %d rows tall, exceeds terminal height %d", label, got, h)
	}

	for i, line := range strings.Split(view, "\n") {
		if got := ansi.StringWidth(line); got > w {
			t.Errorf("%s: line %d is %d cells wide, exceeds terminal width %d: %q",
				label, i, got, w, ansi.Strip(line))
			return
		}
	}
}

func TestCommentsPageFitsTerminal(t *testing.T) {
	sizes := []struct{ w, h int }{
		{60, 20}, {80, 24}, {80, 30}, {100, 40}, {120, 50}, {200, 60},
	}

	for _, size := range sizes {
		page := newTestPage(size.w, size.h)
		page.updateComments(nestedComments(3, "a short reply"))
		assertFits(t, fmt.Sprintf("%dx%d", size.w, size.h), page.View(), size.w, size.h)
	}
}

// The help line is rendered outside the viewport's frame and used to be left
// unconstrained, so at typical widths it wrapped and pushed the page one row
// past the bottom of the terminal on every single comments view.
func TestCommentsPageHelpLineFitsWidth(t *testing.T) {
	for _, w := range []int{60, 80, 90, 96, 97, 100, 120} {
		page := newTestPage(w, 30)
		page.updateComments(nestedComments(1, "a reply"))

		helpView := page.pager.help.View(page.pager.keyMap)
		if got := ansi.StringWidth(helpView); got > w {
			t.Errorf("termW=%d: help line is %d cells wide: %q", w, got, ansi.Strip(helpView))
		}
		if got := lipgloss.Height(helpView); got != 1 {
			t.Errorf("termW=%d: help line wrapped to %d rows", w, got)
		}
	}
}

// Toggling full help ("?") expands the help block considerably; the page must
// still fit inside the terminal.
func TestCommentsPageFitsWithFullHelp(t *testing.T) {
	for _, size := range []struct{ w, h int }{{80, 30}, {100, 40}} {
		page := newTestPage(size.w, size.h)
		page.updateComments(nestedComments(2, "a reply"))

		page, _ = page.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")})
		assertFits(t, fmt.Sprintf("full help %dx%d", size.w, size.h), page.View(), size.w, size.h)
	}
}

// Indentation grows with comment depth. Once it consumed the whole content
// column, lipgloss received a zero or negative width and stopped wrapping,
// emitting lines far wider than the terminal.
func TestDeeplyNestedCommentsStayInBounds(t *testing.T) {
	const w, h = 80, 30
	longText := "this is a reply with a decent amount of text in it that has to wrap somewhere sensible"

	for _, depth := range []int{0, 10, 20, 30, 38, 50, 80} {
		page := newTestPage(w, h)
		page.updateComments(nestedComments(depth, longText))
		assertFits(t, fmt.Sprintf("depth=%d", depth), page.View(), w, h)
	}
}

// Header height depends on how far the description wraps, which depends on the
// width. Measuring it before applying the new width left the page too tall
// after a resize.
func TestCommentsPageFitsAfterResize(t *testing.T) {
	longTitle := strings.Repeat("an extremely long reddit post title that keeps going ", 4)

	page := newTestPage(200, 30)
	page.updateComments(nestedComments(2, "a reply"))

	for _, size := range []struct{ w, h int }{{80, 30}, {60, 20}, {120, 45}, {70, 24}} {
		page.SetSize(size.w, size.h)
		assertFits(t, fmt.Sprintf("resized to %dx%d", size.w, size.h), page.View(), size.w, size.h)
	}

	// Same again, with a title long enough to wrap several times.
	page = newTestPage(200, 30)
	page.updateComments(model.Comments{
		PostTitle: longTitle, Subreddit: "golang", PostAuthor: "bob",
		PostTimestamp: "1 hour ago", PostPoints: "5",
		Comments: []model.Comment{{Author: "x", Text: "hi", Points: "1 point", Timestamp: "now"}},
	})
	for _, size := range []struct{ w, h int }{{80, 30}, {56, 20}, {36, 18}} {
		page.SetSize(size.w, size.h)
		assertFits(t, fmt.Sprintf("long title at %dx%d", size.w, size.h), page.View(), size.w, size.h)
	}
}

// Emoji in post titles and flair are the trigger that made this visible: their
// width must survive the trip through sanitizing and truncation.
func TestCommentsPageWithEmojiContentFits(t *testing.T) {
	const w, h = 80, 30
	texts := []string{
		"⚠️ " + strings.Repeat("warning about this ", 6),
		strings.Repeat("\U0001f468‍\U0001f469‍\U0001f467 ", 12),
		strings.Repeat("❤️", 40),
		strings.Repeat("\U0001f1fa\U0001f1f8", 40),
		strings.Repeat("日本語のテキスト", 12),
	}

	for i, text := range texts {
		page := newTestPage(w, h)
		page.updateComments(nestedComments(3, text))
		assertFits(t, fmt.Sprintf("emoji case %d", i), page.View(), w, h)
	}
}
