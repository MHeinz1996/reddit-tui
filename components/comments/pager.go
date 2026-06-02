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
	viewport        viewport.Model
	postText        string
	postUrl         string
	comments        []model.Comment
	keyMap          viewportKeyMap
	help            help.Model
	collapsed       map[int]bool
	selectedIndex   int
	commentStartLine map[int]int
	commentEndLine   map[int]int
	w, h            int
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
	return lipgloss.JoinVertical(lipgloss.Left, viewportView, helpView)
}

func (c *CommentsViewport) SetSize(w, h int) {
	c.w = w - viewportStyle.GetHorizontalFrameSize()
	c.h = h

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

	c.viewport.Width = c.w
	c.viewport.Height = c.h - helpHeight - 1
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

	postSection := parts[0]
	lineOffset := lipgloss.Height(postSection)
	currentLine := lineOffset

	for i := range c.comments {
		if c.isHidden(i) {
			continue
		}

		c.commentStartLine[i] = currentLine
		commentView := c.formatComment(c.comments[i], i)
		if len(commentView) == 0 {
			continue
		}

		parts = append(parts, commentView)
		commentHeight := lipgloss.Height(commentView)
		currentLine += commentHeight + 2
		c.commentEndLine[i] = currentLine - 3
	}

	return strings.Join(parts, "\n\n")
}

func (c *CommentsViewport) SetViewportContent() {
	content := c.GetViewportView()
	c.viewport.SetContent(content)
}

func (c *CommentsViewport) formatComment(comment model.Comment, i int) string {
	paddingW := comment.Depth * 2
	containerStyle := lipgloss.NewStyle().PaddingLeft(paddingW).Width(c.w - paddingW)

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
	authorView := authorStyle.Render(comment.Author)
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
