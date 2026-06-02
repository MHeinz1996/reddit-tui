package comments

import (
	"reddittui/model"
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

func TestToggleCollapseNoChildrenIsNoop(t *testing.T) {
	c := NewCommentsViewport()
	c.comments = []model.Comment{{Author: "a", Text: "leaf", Depth: 0}}
	c.selectedIndex = 0

	c.toggleCollapseSelected()
	if len(c.collapsed) != 0 {
		t.Fatal("expected collapse to be ignored for comment without replies")
	}
}
