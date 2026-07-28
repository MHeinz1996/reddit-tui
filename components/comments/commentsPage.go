package comments

import (
	"log/slog"
	"reddittui/client"
	"reddittui/components/messages"
	"reddittui/components/styles"
	"reddittui/model"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var commentsErrorText = "Could not load comments. Please try again in a few moments."

// minPagerHeight is how many rows the pager keeps for itself before the header
// is allowed to grow into them.
const minPagerHeight = 4

type CommentsPage struct {
	redditClient   client.RedditClient
	header         CommentsHeader
	pager          CommentsViewport
	containerStyle lipgloss.Style
	commentsUrl    string
	sort           model.CommentSort
	postUrl        string
	focus          bool
}

func NewCommentsPage(redditClient client.RedditClient) CommentsPage {
	header := NewCommentsHeader()
	vp := NewCommentsViewport()

	return CommentsPage{
		redditClient:   redditClient,
		header:         header,
		pager:          vp,
		sort:           model.CommentSortBest,
		containerStyle: styles.GlobalStyle,
	}
}

func (c CommentsPage) Init() tea.Cmd {
	return nil
}

func (c CommentsPage) Update(msg tea.Msg) (CommentsPage, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if c.focus {
		c, cmd = c.handleFocusedMessages(msg)
		cmds = append(cmds, cmd)
	}

	c, cmd = c.handleGlobalMessages(msg)
	cmds = append(cmds, cmd)

	return c, tea.Batch(cmds...)
}

func (c CommentsPage) handleGlobalMessages(msg tea.Msg) (CommentsPage, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.LoadCommentsMsg:
		url := string(msg)
		baseURL, err := client.StripCommentsBaseUrl(url)
		if err != nil {
			slog.Error(commentsErrorText, "error", err)
			return c, messages.ShowErrorModal(commentsErrorText)
		}

		c.commentsUrl = baseURL
		c.sort = model.CommentSortBest
		return c, c.fetchComments()

	case messages.ChangeCommentsSortMsg:
		if c.sort == msg.Sort {
			return c, messages.LoadingComplete
		}

		c.sort = msg.Sort
		return c, c.fetchComments()

	case messages.RefreshCommentsMsg:
		return c, c.refreshComments()

	case messages.UpdateCommentsMsg:
		c.updateComments(model.Comments(msg))
		return c, messages.LoadingComplete
	}

	return c, nil
}

func (c CommentsPage) handleFocusedMessages(msg tea.Msg) (CommentsPage, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "H":
			return c, messages.LoadHome

		case "escape", "backspace", "left", "h":
			return c, messages.GoBack

		case "o", "O":
			return c, messages.OpenUrl(c.postUrl)

		case "r", "R":
			return c, messages.RefreshComments()

		case "1", "2", "3", "4", "5":
			sort, ok := model.CommentSortForKey(keypress)
			if ok {
				return c, messages.ChangeCommentsSort(sort)
			}
		}
	}

	var cmd tea.Cmd
	c.pager, cmd = c.pager.Update(msg)
	return c, cmd
}

func (c CommentsPage) View() string {
	var (
		headerView = c.header.View()
		pagerView  = c.pager.View()
		innerH     = c.containerStyle.GetHeight() - c.containerStyle.GetVerticalFrameSize()
	)

	// Final guard: whatever the components produce, the page must not exceed the
	// rows the terminal gave us, or the redraw desynchronizes and duplicates.
	if innerH > 0 {
		available := max(1, innerH-lipgloss.Height(headerView))
		pagerView = lipgloss.NewStyle().MaxHeight(available).Render(pagerView)
	}

	joined := lipgloss.JoinVertical(lipgloss.Center, headerView, pagerView)
	return c.containerStyle.Render(joined)
}

func (c *CommentsPage) SetSize(w, h int) {
	c.containerStyle = c.containerStyle.Width(w).Height(h)
	c.resizeComponents()
}

func (c *CommentsPage) Focus() {
	c.focus = true
}

func (c *CommentsPage) Blur() {
	c.focus = false
}

func (c *CommentsPage) resizeComponents() {
	var (
		w = c.containerStyle.GetWidth() - c.containerStyle.GetHorizontalFrameSize()
		h = c.containerStyle.GetHeight() - c.containerStyle.GetVerticalFrameSize()
	)

	// Size the header before measuring it: its height depends on how far the
	// description wraps, which depends on the width being set first. Measuring
	// beforehand uses the previous width and leaves the page too tall after a
	// resize, which overflows the terminal and corrupts the redraw.
	c.header.SetSize(w, h)

	// Reserve room for the pager. On a short or narrow terminal the wrapped
	// header can otherwise claim every available row (and then some), pushing
	// the page past the bottom of the screen.
	headerHeight := min(lipgloss.Height(c.header.View()), max(0, h-minPagerHeight))
	c.header.SetMaxHeight(headerHeight)

	c.pager.SetSize(w, h-headerHeight)
}

func (c CommentsPage) fetchComments() tea.Cmd {
	return func() tea.Msg {
		comments, err := c.redditClient.GetComments(c.commentsUrl, c.sort)
		if err != nil {
			slog.Error(commentsErrorText, "error", err)
			return messages.ShowErrorModalMsg{ErrorMsg: commentsErrorText}
		}

		return messages.UpdateCommentsMsg(comments)
	}
}

func (c CommentsPage) refreshComments() tea.Cmd {
	c.redditClient.InvalidateCommentsCache(c.commentsUrl, c.sort)
	return c.fetchComments()
}

func (c *CommentsPage) updateComments(comments model.Comments) {
	if comments.Sort != "" {
		c.sort = comments.Sort
	}

	c.header.SetContent(comments, c.sort)
	c.pager.SetContent(comments)
	c.postUrl = comments.PostUrl

	c.resizeComponents()
}
