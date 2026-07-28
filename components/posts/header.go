package posts

import (
	"reddittui/components/colors"
	"reddittui/utils"

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

type PostsHeader struct {
	DescriptionStyle lipgloss.Style
	Title            string
	Description      string
	W                int
	maxHeight        int
}

func NewPostsHeader() PostsHeader {
	return PostsHeader{DescriptionStyle: defaultDescriptionStyle}
}

func (h *PostsHeader) SetSize(width, height int) {
	h.W = width - headerContainerStyle.GetHorizontalFrameSize()
	h.DescriptionStyle = h.DescriptionStyle.Width(h.W)
}

// SetMaxHeight caps how many rows the header may occupy, trimming the wrapped
// description when the terminal is too short to show all of it. A value of 0
// leaves the header unconstrained.
func (h *PostsHeader) SetMaxHeight(maxHeight int) {
	h.maxHeight = maxHeight
}

func (h PostsHeader) View() string {
	titleView := titleStyle.Render(utils.TruncateString(h.Title, h.W))
	descriptionView := h.DescriptionStyle.Render(h.Description)

	joinedView := lipgloss.JoinVertical(lipgloss.Left, titleView, descriptionView)

	containerStyle := headerContainerStyle
	if h.W > 0 {
		// The styled title bar can outgrow a narrow terminal, so cap the whole
		// block. Without this the container wraps it and the page grows
		// unexpectedly tall.
		containerStyle = containerStyle.MaxWidth(h.W)
	}
	if h.maxHeight > 0 {
		containerStyle = containerStyle.MaxHeight(h.maxHeight)
	}

	return containerStyle.Render(joinedView)
}

func (h *PostsHeader) SetContent(title, desc string) {
	h.Title = utils.NormalizeSubreddit(title)
	h.Description = desc
}
