package messages

import (
	"errors"
	"reddittui/client/common"
	"reddittui/model"

	tea "github.com/charmbracelet/bubbletea"
)

// These are kept to a single line each; the error modal renders the text inline
// after an "Error:" prefix and relies on lipgloss to wrap.
const (
	// SessionErrorText covers an unambiguous auth failure: a login redirect or a 401.
	SessionErrorText = "Reddit session missing or expired. Set auth.sessionCookie in ~/.config/reddittui/reddittui.toml to your reddit_session browser cookie, then restart reddittui."
	// ForbiddenErrorText covers a 403, which reddit returns both for a rejected
	// session cookie and for private or banned content, so it names both causes.
	ForbiddenErrorText = "Reddit refused this request (403). Your reddit_session cookie may have expired, or this content may be private. Check auth.sessionCookie in ~/.config/reddittui/reddittui.toml, then restart reddittui."
)

// LoadError builds the modal message for a failed load, preferring an
// actionable message when the cause is a known one. fallbackMsg is used for
// causes we cannot explain to the user.
func LoadError(err error, fallbackMsg string) ShowErrorModalMsg {
	switch {
	case errors.Is(err, common.ErrNotAuthenticated):
		return ShowErrorModalMsg{ErrorMsg: SessionErrorText, Actionable: true}
	case errors.Is(err, common.ErrForbidden):
		return ShowErrorModalMsg{ErrorMsg: ForbiddenErrorText, Actionable: true}
	default:
		return ShowErrorModalMsg{ErrorMsg: fallbackMsg}
	}
}

type ErrorModalMsg struct {
	ErrorMsg string
	OnClose  tea.Cmd
	// Actionable marks a message that tells the user how to fix the problem, so
	// the initialization handler preserves it instead of replacing it with the
	// generic "check the logfile" text.
	Actionable bool
}

type ChangePostsSortMsg struct {
	Home bool
	Sort model.Sort
}

type ChangeCommentsSortMsg struct {
	Sort model.CommentSort
}

type (
	CleanCacheMsg      struct{}
	GoBackMsg          struct{}
	LoadCommentsMsg    string
	LoadHomeMsg        struct{}
	LoadMorePostsMsg   bool
	LoadSubredditMsg   string
	RefreshCommentsMsg struct{}
	RefreshPostsMsg    bool
	UpdateCommentsMsg  model.Comments
	UpdatePostsMsg     model.Posts
	AddMorePostsMsg    model.Posts
	LoadingCompleteMsg struct{}

	OpenModalMsg        struct{}
	ExitModalMsg        struct{}
	ShowSpinnerModalMsg string

	ShowErrorModalMsg ErrorModalMsg

	OpenUrlMsg string
)

func CleanCache() tea.Msg {
	return CleanCacheMsg{}
}

func GoBack() tea.Msg {
	return GoBackMsg{}
}

func LoadHome() tea.Msg {
	return LoadHomeMsg{}
}

func LoadMorePosts(home bool) tea.Cmd {
	return func() tea.Msg {
		return LoadMorePostsMsg(home)
	}
}

func ChangePostsSort(home bool, sort model.Sort) tea.Cmd {
	return func() tea.Msg {
		return ChangePostsSortMsg{Home: home, Sort: sort}
	}
}

func ChangeCommentsSort(sort model.CommentSort) tea.Cmd {
	return func() tea.Msg {
		return ChangeCommentsSortMsg{Sort: sort}
	}
}

func LoadSubreddit(subreddit string) tea.Cmd {
	return func() tea.Msg {
		return LoadSubredditMsg(subreddit)
	}
}

func LoadComments(url string) tea.Cmd {
	return func() tea.Msg {
		return LoadCommentsMsg(url)
	}
}

func RefreshPosts(home bool) tea.Cmd {
	return func() tea.Msg {
		return RefreshPostsMsg(home)
	}
}

func RefreshComments() tea.Cmd {
	return func() tea.Msg {
		return RefreshCommentsMsg{}
	}
}

func LoadingComplete() tea.Msg {
	return LoadingCompleteMsg{}
}

func OpenModal() tea.Msg {
	return OpenModalMsg{}
}

func ExitModal() tea.Msg {
	return ExitModalMsg{}
}

func ShowSpinnerModal(loadingMsg string) tea.Cmd {
	return func() tea.Msg {
		return ShowSpinnerModalMsg(loadingMsg)
	}
}

func ShowErrorModal(errorMsg string) tea.Cmd {
	return func() tea.Msg {
		return ShowErrorModalMsg{ErrorMsg: errorMsg}
	}
}

func ShowActionableErrorModal(errorMsg string) tea.Cmd {
	return func() tea.Msg {
		return ShowErrorModalMsg{ErrorMsg: errorMsg, Actionable: true}
	}
}

func ShowErrorModalWithCallback(errorMsg string, callback tea.Cmd) tea.Cmd {
	return func() tea.Msg {
		return ShowErrorModalMsg{ErrorMsg: errorMsg, OnClose: callback}
	}
}

func HideSpinnerModal() tea.Msg {
	return ExitModalMsg{}
}

func OpenUrl(url string) tea.Cmd {
	return func() tea.Msg {
		return OpenUrlMsg(url)
	}
}
