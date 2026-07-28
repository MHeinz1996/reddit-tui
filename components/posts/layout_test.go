package posts

import (
	"fmt"
	"strings"
	"testing"

	"reddittui/components/styles"
	"reddittui/model"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// newTestPage builds a PostsPage without a RedditClient. Layout is independent
// of the client, and constructing one would touch the filesystem.
func newTestPage(w, h int) PostsPage {
	items := list.New(nil, NewPostsDelegate(), 0, 0)
	items.SetShowTitle(false)
	items.SetShowStatusBar(false)
	items.KeyMap.NextPage.SetEnabled(false)
	items.KeyMap.PrevPage.SetEnabled(false)
	items.SetFilteringEnabled(false)
	items.AdditionalShortHelpKeys = postsKeys.ShortHelp
	items.AdditionalFullHelpKeys = postsKeys.FullHelp

	page := PostsPage{
		list:           items,
		header:         NewPostsHeader(),
		sort:           model.SortHot,
		containerStyle: styles.GlobalStyle,
		focus:          true,
	}
	page.SetSize(w, h)
	return page
}

func testPosts(description string, titles ...string) model.Posts {
	var ps []model.Post
	for _, title := range titles {
		ps = append(ps, model.Post{
			PostTitle: title, TotalLikes: "1.2k", Subreddit: "r/golang",
			TotalComments: "42", FriendlyDate: "1 hour ago", Author: "bob",
		})
	}

	return model.Posts{Posts: ps, Description: description, Subreddit: "golang"}
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

func TestPostsPageFitsTerminal(t *testing.T) {
	sizes := []struct{ w, h int }{
		{60, 20}, {80, 24}, {80, 30}, {100, 40}, {120, 50}, {200, 60},
	}

	for _, size := range sizes {
		page := newTestPage(size.w, size.h)
		page.updatePosts(testPosts("A subreddit for gophers", "first post", "second post", "third post"))
		assertFits(t, fmt.Sprintf("%dx%d", size.w, size.h), page.View(), size.w, size.h)
	}
}

// Header height depends on how far the description wraps, which depends on the
// width. Measuring it before applying the new width left the page too tall
// after a resize.
func TestPostsPageFitsAfterResize(t *testing.T) {
	longDescription := strings.Repeat("a long subreddit description that wraps repeatedly ", 6)

	page := newTestPage(200, 30)
	page.updatePosts(testPosts(longDescription, "first post", "second post"))

	for _, size := range []struct{ w, h int }{{80, 30}, {60, 20}, {120, 45}, {70, 24}, {40, 18}} {
		page.SetSize(size.w, size.h)
		assertFits(t, fmt.Sprintf("resized to %dx%d", size.w, size.h), page.View(), size.w, size.h)
	}
}

// Emoji in post titles are the trigger that made this visible: a title whose
// width is misjudged survives truncation and then wraps in the terminal.
// Lengths are swept because the corruption only appears at the exact length
// that lands on the truncation boundary.
func TestPostsPageWithEmojiTitlesFits(t *testing.T) {
	const w, h = 80, 30
	prefixes := []string{
		"⚠️ ", "❤️ ", "ℹ️ ", "☠️ ", "➡️ ",
		"\U0001f468‍\U0001f469‍\U0001f467 ",
		"\U0001f1fa\U0001f1f8 ",
		"\U0001f44d\U0001f3fd ",
		"日本語 ",
	}

	for _, prefix := range prefixes {
		for n := 1; n <= 90; n++ {
			title := prefix + strings.Repeat("a", n)
			page := newTestPage(w, h)
			page.updatePosts(testPosts("desc", title, "another post"))
			assertFits(t, fmt.Sprintf("title %q len %d", prefix, n), page.View(), w, h)
		}
	}
}

// The selected row is rendered with a left border and different padding than an
// unselected one, so both states need checking at the truncation boundary.
func TestPostsPageSelectedRowFits(t *testing.T) {
	const w, h = 80, 30

	for n := 1; n <= 90; n++ {
		title := "⚠️ " + strings.Repeat("b", n)
		page := newTestPage(w, h)
		page.updatePosts(testPosts("desc", "plain first post", title))

		assertFits(t, fmt.Sprintf("unselected len %d", n), page.View(), w, h)

		page.list.Select(1)
		assertFits(t, fmt.Sprintf("selected len %d", n), page.View(), w, h)
	}
}
