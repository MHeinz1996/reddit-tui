package comments

import (
	"fmt"
	"reddittui/model"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type CommentsViewport struct {
	viewport         viewport.Model
	postText         string
	postUrl          string
	comments         []model.Comment
	keyMap           viewportKeyMap
	help             help.Model
	collapsed        map[int]bool
	selectedIndex    int
	commentStartLine map[int]int
	commentEndLine   map[int]int
	w, h             int
}

func NewCommentsViewport() CommentsViewport {
	return CommentsViewport{
		viewport:         viewport.New(0, 0),
		keyMap:           commentsKeys,
		help:             help.New(),
		collapsed:        make(map[int]bool),
		selectedIndex:    -1,
		commentStartLine: make(map[int]int),
		commentEndLine:   make(map[int]int),
	}
}

func (c CommentsViewport) Update(msg tea.Msg) (CommentsViewport, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, c.keyMap.CursorUp):
			c.moveSelection(-1)
			return c, nil
		case key.Matches(msg, c.keyMap.CursorDown):
			c.moveSelection(1)
			return c, nil
		case key.Matches(msg, c.keyMap.GoToStart):
			c.selectFirst()
			return c, nil
		case key.Matches(msg, c.keyMap.GoToEnd):
			c.selectLast()
			return c, nil
		case key.Matches(msg, c.keyMap.CollapseComments):
			c.toggleCollapseSelected()
			return c, nil
		case key.Matches(msg, c.keyMap.ShowFullHelp),
			key.Matches(msg, c.keyMap.CloseFullHelp):
			c.help.ShowAll = !c.help.ShowAll
			// Full help occupies several rows instead of one, so the viewport
			// has to give that space back or the page outgrows the terminal.
			c.ResizeComponents()
			c.ensureSelectedVisible()
			return c, nil
		}
	}

	var cmd tea.Cmd
	c.viewport, cmd = c.viewport.Update(msg)
	return c, cmd
}

func (c CommentsViewport) View() string {
	viewportView := viewportStyle.Render(c.viewport.View())
	helpView := c.help.View(c.keyMap)

	// The full help block has a fixed natural height that cannot shrink. On a
	// short terminal it must be clipped rather than allowed to push the page
	// past the bottom of the screen.
	if c.h > 0 {
		helpView = lipgloss.NewStyle().MaxHeight(max(1, c.h-lipgloss.Height(viewportView))).Render(helpView)
	}

	return lipgloss.JoinVertical(lipgloss.Left, viewportView, helpView)
}

func (c *CommentsViewport) SetSize(w, h int) {
	c.w = w - viewportStyle.GetHorizontalFrameSize()
	c.h = h

	// The help view is joined outside viewportStyle, so it gets the full width.
	// Constraining it is what keeps it from overflowing and soft-wrapping.
	c.help.Width = w

	c.ResizeComponents()
	c.SetViewportContent()
}

func (c *CommentsViewport) SetContent(comments model.Comments) {
	c.postText = comments.PostText
	c.postUrl = comments.PostUrl
	c.comments = comments.Comments

	c.collapsed = make(map[int]bool)
	c.selectedIndex = -1
	if len(c.comments) > 0 {
		c.selectedIndex = 0
	}

	c.viewport.SetYOffset(0)
	c.ResizeComponents()
	c.SetViewportContent()
}

func (c *CommentsViewport) ResizeComponents() {
	helpHeight := lipgloss.Height(c.help.View(c.keyMap))

	c.viewport.Width = max(1, c.w)
	// viewportStyle contributes a bottom margin, hence the extra row. Clamp at
	// zero so a cramped terminal cannot produce a negative height, which the
	// viewport would render as content overflowing its bounds.
	c.viewport.Height = max(0, c.h-helpHeight-viewportStyle.GetVerticalFrameSize())
}

func (c *CommentsViewport) GetViewportView() string {
	var parts []string

	if len(c.postText) > 0 {
		parts = append(parts, c.postText)
	} else {
		parts = append(parts, c.postUrl)
	}

	c.commentStartLine = make(map[int]int)
	c.commentEndLine = make(map[int]int)

	// Parts are joined with a blank line between them, so each part occupies its
	// own height plus exactly one separator row. Counting more than that makes
	// the tracked line numbers run ahead of the rendered content, and
	// ensureSelectedVisible then scrolls past the selected comment.
	currentLine := lipgloss.Height(parts[0]) + 1

	for i := range c.comments {
		if c.isHidden(i) {
			continue
		}

		commentView := c.formatComment(c.comments[i], i)
		if len(commentView) == 0 {
			continue
		}

		parts = append(parts, commentView)
		commentHeight := lipgloss.Height(commentView)
		c.commentStartLine[i] = currentLine
		c.commentEndLine[i] = currentLine + commentHeight - 1
		currentLine += commentHeight + 1
	}

	return strings.Join(parts, "\n\n")
}

func (c *CommentsViewport) SetViewportContent() {
	content := c.GetViewportView()
	c.viewport.SetContent(content)
}

// minCommentWidth is the narrowest content column a deeply nested comment may
// be squeezed into. Without a floor, indentation eventually exceeds the
// available width and lipgloss receives a zero or negative Width, at which
// point it stops wrapping altogether and emits lines wider than the terminal.
const minCommentWidth = 20

func (c *CommentsViewport) formatComment(comment model.Comment, i int) string {
	paddingW := comment.Depth * 2
	if maxPadding := c.w - minCommentWidth; paddingW > maxPadding {
		paddingW = max(0, maxPadding)
	}

	// Derive the content column from the (already clamped) padding so the two
	// always sum to at most c.w, even when c.w itself is below the floor.
	contentW := max(1, c.w-paddingW)
	containerStyle := lipgloss.NewStyle().PaddingLeft(paddingW).Width(contentW)

	if c.collapsed[i] && c.showsHiddenPlaceholder(i) {
		hiddenMsg := collapsedStyle.Render("(comment hidden)")
		if i == c.selectedIndex {
			hiddenMsg = selectedCommentAuthorStyle.Render("(comment hidden)")
		}
		return containerStyle.Render(hiddenMsg)
	}

	var (
		authorAndDateView string
		pointsView        string
	)

	authorStyle := commentAuthorStyle
	if i == c.selectedIndex {
		authorStyle = selectedCommentAuthorStyle
	}
	authorView := authorStyle.Render(comment.AuthorLabel())
	dateView := commentDateStyle.Render(comment.Timestamp)
	pointsView = renderPoints(comment.Points)
	authorAndDateView = fmt.Sprintf("%s • %s • %s", authorView, dateView, pointsView)

	if c.collapsed[i] {
		children := c.countDescendants(i)
		if children == 1 {
			collapsedHintView := collapsedStyle.Render("(1 comment hidden)")
			authorAndDateView = fmt.Sprintf("%s  %s", authorAndDateView, collapsedHintView)
		} else if children > 1 {
			collapsedView := collapsedStyle.Render(fmt.Sprintf("(%d comments hidden)", children))
			authorAndDateView = fmt.Sprintf("%s  %s", authorAndDateView, collapsedView)
		}
	}

	joined := lipgloss.JoinVertical(lipgloss.Left, authorAndDateView, comment.Text)
	return containerStyle.Render(joined)
}

func renderPoints(pointsString string) string {
	parts := strings.Fields(pointsString)
	if len(parts) != 2 {
		return defaultPointsStyle.Render(pointsString)
	}

	if strings.Contains(parts[0], "-") {
		return negativePointsStyle.Render(pointsString)
	} else if strings.Contains(parts[0], "k") {
		return popularPointsStyle.Render(pointsString)
	}

	points, err := strconv.Atoi(parts[0])
	if err != nil {
		return defaultPointsStyle.Render(pointsString)
	} else if points >= 1000 {
		return popularPointsStyle.Render(pointsString)
	}

	return defaultPointsStyle.Render(pointsString)
}

func (c *CommentsViewport) moveSelection(delta int) {
	visible := c.visibleIndices()
	if len(visible) == 0 {
		c.selectedIndex = -1
		return
	}

	currentPos := -1
	for i, idx := range visible {
		if idx == c.selectedIndex {
			currentPos = i
			break
		}
	}

	if currentPos < 0 {
		c.selectedIndex = visible[0]
	} else {
		newPos := currentPos + delta
		if newPos < 0 {
			newPos = 0
		} else if newPos >= len(visible) {
			newPos = len(visible) - 1
		}
		c.selectedIndex = visible[newPos]
	}

	c.SetViewportContent()
	c.ensureSelectedVisible()
}

func (c *CommentsViewport) selectFirst() {
	visible := c.visibleIndices()
	if len(visible) == 0 {
		c.selectedIndex = -1
		return
	}

	c.selectedIndex = visible[0]
	c.SetViewportContent()
	c.viewport.GotoTop()
}

func (c *CommentsViewport) selectLast() {
	visible := c.visibleIndices()
	if len(visible) == 0 {
		c.selectedIndex = -1
		return
	}

	c.selectedIndex = visible[len(visible)-1]
	c.SetViewportContent()
	c.ensureSelectedVisible()
}

func (c *CommentsViewport) toggleCollapseSelected() {
	if c.selectedIndex < 0 {
		return
	}

	if c.collapsed[c.selectedIndex] {
		delete(c.collapsed, c.selectedIndex)
	} else {
		c.collapsed[c.selectedIndex] = true
	}

	c.SetViewportContent()
	c.ensureSelectedVisible()
}

func (c *CommentsViewport) showsHiddenPlaceholder(i int) bool {
	if !c.collapsed[i] {
		return false
	}

	return c.comments[i].Depth == 0 || !c.hasChildren(i)
}

func (c *CommentsViewport) ensureSelectedVisible() {
	if c.selectedIndex < 0 {
		return
	}

	start, ok := c.commentStartLine[c.selectedIndex]
	if !ok {
		return
	}

	end := c.commentEndLine[c.selectedIndex]
	if end < start {
		end = start
	}

	viewTop := c.viewport.YOffset
	viewBottom := viewTop + c.viewport.Height - 1

	if start < viewTop {
		c.viewport.SetYOffset(start)
	} else if end > viewBottom {
		newOffset := end - c.viewport.Height + 1
		if newOffset < 0 {
			newOffset = 0
		}
		c.viewport.SetYOffset(newOffset)
	}
}

func (c *CommentsViewport) isHidden(i int) bool {
	if i < 0 || i >= len(c.comments) {
		return true
	}

	depth := c.comments[i].Depth
	for j := i - 1; j >= 0; j-- {
		if c.comments[j].Depth < depth {
			if c.collapsed[j] {
				return true
			}
			depth = c.comments[j].Depth
		}
	}

	return false
}

func (c *CommentsViewport) visibleIndices() []int {
	var indices []int
	for i := range c.comments {
		if !c.isHidden(i) {
			indices = append(indices, i)
		}
	}
	return indices
}

func (c *CommentsViewport) hasChildren(i int) bool {
	return i+1 < len(c.comments) && c.comments[i+1].Depth > c.comments[i].Depth
}

func (c *CommentsViewport) countDescendants(i int) int {
	if i >= len(c.comments) {
		return 0
	}

	parentDepth := c.comments[i].Depth
	count := 0
	for j := i + 1; j < len(c.comments); j++ {
		if c.comments[j].Depth <= parentDepth {
			break
		}
		count++
	}

	return count
}
