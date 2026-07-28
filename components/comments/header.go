package comments

import (
	"fmt"
	"reddittui/components/colors"
	"reddittui/model"
	"reddittui/utils"
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

var (
	headerContainerStyle = lipgloss.NewStyle().MarginBottom(2)
	titleStyle           = lipgloss.NewStyle().
				MarginBottom(1).
				Padding(0, 2).
				Height(1).
				Background(colors.AdaptiveColors(colors.Blue, colors.Indigo)).
				Foreground(colors.AdaptiveColors(colors.White, colors.Sand))

	defaultDescriptionStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colors.AdaptiveColor(colors.Text))
)

type CommentsHeader struct {
	DescriptionStyle lipgloss.Style
	Title            string
	Description      string
	Author           string
	AuthorFlair      string
	Timestamp        string
	Points           string
	SortLabel        string
	TotalComments    int
	W                int
	maxHeight        int
}

func NewCommentsHeader() CommentsHeader {
	return CommentsHeader{DescriptionStyle: defaultDescriptionStyle}
}

func (h *CommentsHeader) SetSize(width, height int) {
	h.W = width - headerContainerStyle.GetHorizontalFrameSize()
	h.DescriptionStyle = h.DescriptionStyle.Width(h.W)
}

// SetMaxHeight caps how many rows the header may occupy, trimming the wrapped
// post title when the terminal is too short to show all of it. A value of 0
// leaves the header unconstrained.
func (h *CommentsHeader) SetMaxHeight(maxHeight int) {
	h.maxHeight = maxHeight
}

func (h CommentsHeader) View() string {
	titleView := titleStyle.Render(utils.TruncateString(h.Title, h.W))
	descriptionView := h.DescriptionStyle.Render(h.Description)

	authorView := postAuthorStyle.Render(utils.AuthorWithFlair(h.Author, h.AuthorFlair))
	timestampView := postTimestampStyle.Render(fmt.Sprintf("submitted %s by", h.Timestamp))
	authorTimestampView := fmt.Sprintf("%s %s", timestampView, authorView)

	postPointsView := postPointsStyle.Render(utils.GetSingularPlural(h.Points, "point", "points"))
	totalCommentsView := totalCommentsStyle.Render(utils.GetSingularPlural(strconv.Itoa(h.TotalComments), "comment", "comments"))
	sortView := postTimestampStyle.Render(fmt.Sprintf("sorted by %s", h.SortLabel))
	pointsAndCommentsView := fmt.Sprintf("%s • %s • %s", postPointsView, totalCommentsView, sortView)

	joinedView := lipgloss.JoinVertical(lipgloss.Left, titleView, descriptionView, authorTimestampView, pointsAndCommentsView)

	containerStyle := headerContainerStyle
	if h.W > 0 {
		// The metadata rows are assembled from independently styled fragments
		// and can outgrow a narrow terminal, so cap the whole block. Without
		// this the container wraps them and the page grows unexpectedly tall.
		containerStyle = containerStyle.MaxWidth(h.W)
	}
	if h.maxHeight > 0 {
		containerStyle = containerStyle.MaxHeight(h.maxHeight)
	}

	return containerStyle.Render(joinedView)
}

func (h *CommentsHeader) SetContent(comments model.Comments, sort model.CommentSort) {
	h.Title = utils.NormalizeSubreddit(comments.Subreddit)
	h.Description = comments.PostTitle
	h.Author = comments.PostAuthor
	h.AuthorFlair = comments.PostAuthorFlair
	h.TotalComments = len(comments.Comments)
	h.Timestamp = comments.PostTimestamp
	h.Points = comments.PostPoints
	h.SortLabel = sort.Label()
}
