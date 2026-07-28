package posts

import (
	"fmt"
	"log/slog"
	"reddittui/client"
	"reddittui/client/common"
	"reddittui/components/messages"
	"reddittui/components/styles"
	"reddittui/model"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultHeaderTitle       = "reddit.com"
	defaultHeaderDescription = "The front page of the internet"
	postsErrorText           = "Could not load posts. Please try again in a few moments."
	subredditNotFoundText    = "Subreddit not found"
)

// minListHeight is how many rows the post list keeps for itself before the
// header is allowed to grow into them.
const minListHeight = 4

type PostsPage struct {
	Subreddit      string
	sort           model.Sort
	posts          model.Posts
	redditClient   client.RedditClient
	header         PostsHeader
	list           list.Model
	focus          bool
	Home           bool
	containerStyle lipgloss.Style
}

func NewPostsPage(redditClient client.RedditClient, home bool) PostsPage {
	items := list.New(nil, NewPostsDelegate(), 0, 0)
	items.SetShowTitle(false)
	items.SetShowStatusBar(false)
	items.KeyMap.NextPage.SetEnabled(false)
	items.KeyMap.PrevPage.SetEnabled(false)
	items.SetFilteringEnabled(false)
	items.AdditionalShortHelpKeys = postsKeys.ShortHelp
	items.AdditionalFullHelpKeys = postsKeys.FullHelp

	header := NewPostsHeader()
	if home {
		header.SetContent(defaultHeaderTitle, defaultHeaderDescription)
	}

	containerStyle := styles.GlobalStyle

	return PostsPage{
		list:           items,
		redditClient:   redditClient,
		header:         header,
		Home:           home,
		sort:           model.SortHot,
		containerStyle: containerStyle,
	}
}

func (p PostsPage) Init() tea.Cmd {
	return nil
}

func (p PostsPage) Update(msg tea.Msg) (PostsPage, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	if p.focus {
		p, cmd = p.handleFocusedMessages(msg)
		cmds = append(cmds, cmd)
	}

	p, cmd = p.handleGlobalMessages(msg)
	cmds = append(cmds, cmd)

	return p, tea.Batch(cmds...)
}

func (p PostsPage) handleGlobalMessages(msg tea.Msg) (PostsPage, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.LoadHomeMsg:
		if p.Home {
			return p, p.loadHome()
		}

	case messages.LoadSubredditMsg:
		if !p.Home {
			subreddit := string(msg)
			return p, p.loadSubreddit(subreddit)
		}

	case messages.LoadMorePostsMsg:
		isHome := bool(msg)
		if p.Home == isHome {
			return p, p.loadMorePosts()
		}

	case messages.ChangePostsSortMsg:
		if p.Home == msg.Home {
			if p.sort == msg.Sort {
				return p, messages.LoadingComplete
			}

			p.sort = msg.Sort
			return p, p.reloadPosts()
		}

	case messages.RefreshPostsMsg:
		isHome := bool(msg)
		if p.Home == isHome {
			return p, p.refreshPosts()
		}

	case messages.UpdatePostsMsg:
		posts := model.Posts(msg)
		if posts.IsHome == p.Home {
			p.updatePosts(posts)
			return p, messages.LoadingComplete
		}

	case messages.AddMorePostsMsg:
		posts := model.Posts(msg)
		if posts.IsHome == p.Home {
			p.addPosts(posts)
			return p, messages.LoadingComplete
		}
	}

	return p, nil
}

func (p PostsPage) handleFocusedMessages(msg tea.Msg) (PostsPage, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter", "right", "l":
			item := p.list.SelectedItem()
			post, ok := item.(model.Post)
			if !ok {
				return p, nil
			}

			return p, func() tea.Msg {
				return messages.LoadCommentsMsg(post.CommentsUrl)
			}

		case "q", "Q":
			// Ignore q keystrokes to list.Modal. since it will default to sending a Quit message
			// instead of showing the quit modal. Tui component will correctly handle quit mesages
			return p, nil

		case "L":
			return p, messages.LoadMorePosts(p.Home)

		case "r", "R":
			return p, messages.RefreshPosts(p.Home)

		case "1", "2", "3", "4", "5":
			sort, ok := model.SortForKey(keypress)
			if ok {
				return p, messages.ChangePostsSort(p.Home, sort)
			}

		case "H":
			return p, messages.LoadHome

		case "esc", "backspace", "left", "h":
			return p, messages.GoBack
		}
	}

	var cmd tea.Cmd
	p.list, cmd = p.list.Update(msg)
	return p, cmd
}

func (p PostsPage) View() string {
	if len(p.posts.Posts) == 0 {
		return p.containerStyle.Render("")
	}

	var (
		headerView = p.header.View()
		listView   = p.list.View()
		innerH     = p.containerStyle.GetHeight() - p.containerStyle.GetVerticalFrameSize()
	)

	// The list keeps a minimum height of its own, so on a very short terminal it
	// can render taller than the rows it was given. Clip it rather than let the
	// page overflow the screen.
	if innerH > 0 {
		available := max(1, innerH-lipgloss.Height(headerView))
		listView = lipgloss.NewStyle().MaxHeight(available).Render(listView)
	}

	joined := lipgloss.JoinVertical(lipgloss.Left, headerView, listView)
	return p.containerStyle.Render(joined)
}

func (p *PostsPage) SetSize(w, h int) {
	p.containerStyle = p.containerStyle.Width(w).Height(h)
	p.resizeComponents()
}

func (p *PostsPage) Focus() {
	p.focus = true
}

func (p *PostsPage) Blur() {
	p.focus = false
}

func (p *PostsPage) resizeComponents() {
	var (
		w         = p.containerStyle.GetWidth() - p.containerStyle.GetHorizontalFrameSize()
		h         = p.containerStyle.GetHeight() - p.containerStyle.GetVerticalFrameSize()
		listWidth = w - postsListStyle.GetHorizontalFrameSize()
	)

	// Size the header before measuring it: its height depends on how far the
	// description wraps, which depends on the width being set first. Measuring
	// beforehand uses the previous width and leaves the page too tall after a
	// resize, which overflows the terminal and corrupts the redraw.
	p.header.SetSize(w, h)

	// Reserve room for the list. On a short or narrow terminal the wrapped
	// description can otherwise claim every available row (and then some),
	// pushing the page past the bottom of the screen.
	headerHeight := min(lipgloss.Height(p.header.View()), max(0, h-minListHeight))
	p.header.SetMaxHeight(headerHeight)

	p.list.SetSize(listWidth, h-headerHeight)
}

func (p *PostsPage) loadHome() tea.Cmd {
	p.sort = model.SortHot

	return func() tea.Msg {
		posts, err := p.redditClient.GetHomePosts("", p.sort)
		if err != nil {
			slog.Error(postsErrorText, "error", err)
			return messages.ShowErrorModalMsg{ErrorMsg: postsErrorText}
		}

		return messages.UpdatePostsMsg(posts)
	}
}

func (p *PostsPage) reloadPosts() tea.Cmd {
	return func() tea.Msg {
		var (
			posts model.Posts
			err   error
		)

		if p.Home {
			posts, err = p.redditClient.GetHomePosts("", p.sort)
		} else {
			posts, err = p.redditClient.GetSubredditPosts(p.Subreddit, "", p.sort)
		}

		if err == common.ErrNotFound {
			slog.Error(subredditNotFoundText, "error", err, "subreddit", p.Subreddit)
			return messages.ShowErrorModalMsg{ErrorMsg: fmt.Sprintf("%s: %s", subredditNotFoundText, p.Subreddit)}
		} else if err != nil {
			slog.Error(postsErrorText, "error", err)
			return messages.ShowErrorModalMsg{ErrorMsg: postsErrorText}
		}

		return messages.UpdatePostsMsg(posts)
	}
}

func (p *PostsPage) refreshPosts() tea.Cmd {
	if p.Home {
		p.redditClient.InvalidatePostsCache("", p.sort)
	} else {
		p.redditClient.InvalidatePostsCache(p.Subreddit, p.sort)
	}

	return p.reloadPosts()
}

func (p *PostsPage) loadMorePosts() tea.Cmd {
	return func() tea.Msg {
		var (
			posts model.Posts
			err   error
		)

		if len(p.posts.After) == 0 {
			slog.Error(postsErrorText, "error", err)
			return messages.ShowErrorModalMsg{ErrorMsg: postsErrorText}
		}

		if p.posts.IsHome {
			posts, err = p.redditClient.GetHomePosts(p.posts.After, p.sort)
		} else {
			posts, err = p.redditClient.GetSubredditPosts(p.Subreddit, p.posts.After, p.sort)
		}

		if err != nil {
			slog.Error(postsErrorText, "error", err)
			return messages.ShowErrorModalMsg{ErrorMsg: postsErrorText}
		}

		return messages.AddMorePostsMsg(posts)
	}
}

func (p PostsPage) loadSubreddit(subreddit string) tea.Cmd {
	p.sort = model.SortHot

	return func() tea.Msg {
		posts, err := p.redditClient.GetSubredditPosts(subreddit, "", p.sort)
		if err == common.ErrNotFound {
			slog.Error(subredditNotFoundText, "error", err, "subreddit", subreddit)
			return messages.ShowErrorModalMsg{ErrorMsg: fmt.Sprintf("%s: %s", subredditNotFoundText, subreddit)}
		} else if err != nil {
			slog.Error(postsErrorText, "error", err)
			return messages.ShowErrorModalMsg{ErrorMsg: postsErrorText}
		}

		return messages.UpdatePostsMsg(posts)
	}
}

func (p *PostsPage) updatePosts(posts model.Posts) {
	p.posts = posts
	if posts.Sort != "" {
		p.sort = posts.Sort
	}

	if posts.IsHome {
		p.header.SetContent(defaultHeaderTitle, p.sortDescription(defaultHeaderDescription))
	} else {
		p.header.SetContent(posts.Subreddit, p.sortDescription(posts.Description))
		p.Subreddit = posts.Subreddit
	}

	p.list.ResetSelected()

	var listItems []list.Item
	for _, p := range posts.Posts {
		listItems = append(listItems, p)
	}
	p.list.SetItems(listItems)

	// Need to set size again when content loads so padding and margins are correct
	p.resizeComponents()
}

func (p *PostsPage) addPosts(posts model.Posts) {
	p.posts.Posts = append(p.posts.Posts, posts.Posts...)
	p.posts.Posts = dedupePosts(p.posts.Posts)
	p.posts.After = posts.After

	var listItems []list.Item
	for _, post := range p.posts.Posts {
		listItems = append(listItems, post)
	}
	p.list.SetItems(listItems)

	// Need to set size again when content loads so padding and margins are correct
	p.resizeComponents()
}

func dedupePosts(posts []model.Post) []model.Post {
	seen := make(map[string]struct{}, len(posts))
	out := make([]model.Post, 0, len(posts))

	for _, post := range posts {
		key := post.CommentsUrl
		if key == "" {
			key = post.PostTitle + "|" + post.Author + "|" + post.FriendlyDate
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, post)
	}

	return out
}

func (p PostsPage) sortDescription(description string) string {
	sortLabel := fmt.Sprintf("sorted by %s", p.sort.Label())
	if len(description) == 0 {
		return sortLabel
	}

	return fmt.Sprintf("%s • %s", description, sortLabel)
}
