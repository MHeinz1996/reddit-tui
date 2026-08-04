package comments

import (
	"fmt"
	"reddittui/model"
	"strings"
	"testing"
)

func testComments() []model.Comment {
	return []model.Comment{
		{Author: "a", Text: "root", Depth: 0},
		{Author: "b", Text: "reply", Depth: 1},
		{Author: "c", Text: "nested", Depth: 2},
		{Author: "d", Text: "root2", Depth: 0},
	}
}

// Line bookkeeping must match the rendered content exactly. When the tracked
// numbers run ahead, ensureSelectedVisible scrolls to rows that do not exist,
// which shifts the page up until the selected comment is off screen.
func TestCommentLineNumbersMatchRenderedContent(t *testing.T) {
	for _, wrapping := range []bool{false, true} {
		c := NewCommentsViewport()
		c.w = 40
		c.h = 20
		c.postText = "post line one\npost line two"

		for i := range 8 {
			comment := model.Comment{
				Author:    fmt.Sprintf("user%d", i),
				Text:      fmt.Sprintf("body%d", i),
				Timestamp: "1h",
				Points:    "5 points",
			}
			if wrapping {
				// Nesting plus long text forces multi-line comments, so a
				// per-comment height error cannot cancel itself out.
				comment.Depth = i % 4
				comment.Text = strings.Repeat(comment.Text+" ", 6)
			}
			c.comments = append(c.comments, comment)
		}

		lines := strings.Split(c.GetViewportView(), "\n")

		for i := range c.comments {
			author := fmt.Sprintf("user%d ", i)
			renderedStart := -1
			for n, line := range lines {
				if strings.Contains(line, author) {
					renderedStart = n
					break
				}
			}

			if got := c.commentStartLine[i]; got != renderedStart {
				t.Errorf("wrapping=%v comment %d: start line %d, rendered at %d", wrapping, i, got, renderedStart)
			}
			if got := c.commentEndLine[i]; got > len(lines)-1 {
				t.Errorf("wrapping=%v comment %d: end line %d past content of %d lines", wrapping, i, got, len(lines))
			}
			if got := c.commentEndLine[i]; got < c.commentStartLine[i] {
				t.Errorf("wrapping=%v comment %d: end line %d before start line %d", wrapping, i, got, c.commentStartLine[i])
			}
		}
	}
}

func TestIsHiddenWithPerCommentCollapse(t *testing.T) {
	c := NewCommentsViewport()
	c.comments = testComments()
	c.collapsed[0] = true

	if !c.isHidden(1) {
		t.Fatal("expected reply to be hidden when parent is collapsed")
	}
	if !c.isHidden(2) {
		t.Fatal("expected nested reply to be hidden when ancestor is collapsed")
	}
	if c.isHidden(3) {
		t.Fatal("expected unrelated root comment to remain visible")
	}
}

func TestVisibleIndices(t *testing.T) {
	c := NewCommentsViewport()
	c.comments = testComments()
	c.collapsed[0] = true

	visible := c.visibleIndices()
	if len(visible) != 2 || visible[0] != 0 || visible[1] != 3 {
		t.Fatalf("expected indices [0 3] visible, got %v", visible)
	}
}

func TestToggleCollapseSelected(t *testing.T) {
	c := NewCommentsViewport()
	c.comments = testComments()
	c.selectedIndex = 0

	c.toggleCollapseSelected()
	if !c.collapsed[0] {
		t.Fatal("expected selected comment to be collapsed")
	}
	if len(c.visibleIndices()) != 2 {
		t.Fatalf("expected two visible comments after collapse, got %d", len(c.visibleIndices()))
	}

	c.toggleCollapseSelected()
	if c.collapsed[0] {
		t.Fatal("expected selected comment to be expanded")
	}
}

func TestToggleCollapseLeafComment(t *testing.T) {
	c := NewCommentsViewport()
	c.w = 80
	c.comments = []model.Comment{{Author: "a", Text: "leaf", Depth: 0}}
	c.selectedIndex = 0

	c.toggleCollapseSelected()
	if !c.collapsed[0] {
		t.Fatal("expected leaf comment to be collapsible")
	}

	view := c.formatComment(c.comments[0], 0)
	if !strings.Contains(view, "comment hidden") {
		t.Fatalf("expected hidden placeholder, got %q", view)
	}
	if strings.Contains(view, "leaf") {
		t.Fatalf("expected comment body to be hidden, got %q", view)
	}
}

func TestCollapsedTopLevelWithRepliesShowsPlaceholder(t *testing.T) {
	c := NewCommentsViewport()
	c.w = 80
	c.comments = testComments()
	c.collapsed[0] = true

	view := c.formatComment(c.comments[0], 0)
	if !strings.Contains(view, "comment hidden") {
		t.Fatalf("expected hidden placeholder for collapsed top-level comment, got %q", view)
	}
	if strings.Contains(view, "root") {
		t.Fatalf("expected top-level comment body to be hidden, got %q", view)
	}
}

func TestCollapsedNestedWithRepliesShowsThreadHint(t *testing.T) {
	c := NewCommentsViewport()
	c.w = 80
	c.comments = testComments()
	c.collapsed[1] = true

	view := c.formatComment(c.comments[1], 1)
	if !strings.Contains(view, "reply") {
		t.Fatalf("expected nested comment body to remain visible, got %q", view)
	}
	if !strings.Contains(view, "1 comment hidden") {
		t.Fatalf("expected thread collapse hint, got %q", view)
	}
}
