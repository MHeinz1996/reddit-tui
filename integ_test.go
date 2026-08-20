package main

import (
	"bytes"
	"os"
	"reddittui/components"
	"reddittui/config"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
)

const (
	testTimeout    = 20 * time.Second
	testDomain     = "old.reddit.com"
	testServerType = "old"
)

func TestStartup(t *testing.T) {
	t.Logf("Testing startup...")
	configuration := getTestConfig()

	tui := components.NewRedditTui(configuration, "", "")
	tm := teatest.NewTestModel(t, tui, teatest.WithInitialTermSize(300, 100))

	t.Logf("\tVerify the loading screen shows on startup...")
	WaitFor(t, tm, "loading reddit.com...")

	t.Logf("\tVerify the home page loads...")
	WaitFor(t, tm, "The front page of the internet")
}

func TestSwitchSubreddit(t *testing.T) {
	t.Logf("Testing switching subreddit...")
	configuration := getTestConfig()

	tui := components.NewRedditTui(configuration, "", "")
	tm := teatest.NewTestModel(t, tui, teatest.WithInitialTermSize(300, 100))

	t.Logf("\tVerify home page loads...")
	WaitFor(t, tm, "The front page of the internet")

	t.Logf("\tVerify subreddit selection modal shows...")
	WaitForWithInputs(t, tm, "s", "Choose a subreddit:")

	t.Logf("\tVerify dogs subreddit loads...")
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("dogs"),
	})
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	WaitFor(t, tm, "r/dogs")
}

func TestReturnToHomePage(t *testing.T) {
	t.Logf("Testing returning to the home page after switching subreddits...")
	configuration := getTestConfig()

	tui := components.NewRedditTui(configuration, "", "")
	tm := teatest.NewTestModel(t, tui, teatest.WithInitialTermSize(300, 100))

	t.Logf("\tVerify home page loads...")
	WaitFor(t, tm, "The front page of the internet")

	t.Logf("\tVerify subreddit selection modal shows...")
	WaitForWithInputs(t, tm, "s", "Choose a subreddit:")

	t.Logf("\tVerify dogs subreddit loads...")
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("dogs"),
	})
	tm.Send(tea.KeyMsg{
		Type: tea.KeyEnter,
	})
	WaitFor(t, tm, "r/dogs ")

	t.Logf("\tVerify home page loads")
	WaitForWithInputs(t, tm, "H", "The front page of the internet")
}

func TestShowComments(t *testing.T) {
	t.Logf("Testing show post comments...")
	configuration := getTestConfig()

	tui := components.NewRedditTui(configuration, "", "")
	tm := teatest.NewTestModel(t, tui, teatest.WithInitialTermSize(300, 100))

	t.Logf("\tVerify home page loads...")
	WaitFor(t, tm, "The front page of the internet")

	t.Logf("\tVerify subreddit selection modal shows...")
	WaitForWithInputs(t, tm, "s", "Choose a subreddit:")

	t.Logf("\tVerify dogs subreddit loads...")
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("dogs"),
	})
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	WaitFor(t, tm, "r/dogs", "/r/dogs")

	t.Logf("\tVerify comments header loads...")
	time.Sleep(time.Second)
	WaitForWithInputs(t, tm, "l", "submitted", "ago by", "point", "comment")
}

func TestLoadInitialPostFromId(t *testing.T) {
	t.Logf("Testing loading initial post...")
	configuration := getTestConfig()

	postId := "1jgxswb"
	tui := components.NewRedditTui(configuration, "", postId)
	tm := teatest.NewTestModel(t, tui, teatest.WithInitialTermSize(300, 100))

	t.Logf("\tVerify comments header loads...")
	time.Sleep(time.Second)
	WaitForWithInputs(t, tm, "l", "submitted", "ago by", "point", "comment")
}

func TestLoadInitialPostFromUrl(t *testing.T) {
	t.Logf("Testing loading initial post...")
	configuration := getTestConfig()

	postUrl := "https://old.reddit.com/r/dogs/comments/1jh0yne/dog_becoming_cuddlier_as_a_senior/"
	tui := components.NewRedditTui(configuration, "", postUrl)
	tm := teatest.NewTestModel(t, tui, teatest.WithInitialTermSize(300, 100))

	t.Logf("\tVerify comments header loads...")
	time.Sleep(time.Second)
	WaitForWithInputs(t, tm, "l", "submitted", "ago by", "point", "comment")
}

func TestLoadInitialSubredditAndCanGoBack(t *testing.T) {
	t.Logf("Testing loading subreddit...")
	configuration := getTestConfig()

	tui := components.NewRedditTui(configuration, "dogs", "")
	tm := teatest.NewTestModel(t, tui, teatest.WithInitialTermSize(300, 100))

	t.Logf("\tVerify dog subreddit loads...")
	WaitFor(t, tm, "r/dogs", "/r/dogs")
	time.Sleep(time.Second)

	// Go back
	t.Logf("\tVerify the home page loads...")
	WaitForWithInputs(t, tm, "h", "The front page of the internet")
}

// With no session at all, reddit answers with a 302 to its login page, which
// http.Client follows into a 200 full of login HTML. That used to parse to zero
// posts and get reported as a missing subreddit.
func TestStartupWithNoSessionShowsActionableError(t *testing.T) {
	t.Logf("Testing startup with no session cookie...")
	configuration := getTestConfig()
	configuration.Auth.SessionCookie = ""

	tui := components.NewRedditTui(configuration, "", "")
	tm := teatest.NewTestModel(t, tui, teatest.WithInitialTermSize(300, 100))

	t.Logf("\tVerify the session error names the setting to update...")
	WaitFor(t, tm, "expired", "sessionCookie", "reddit_session")
}

// A rejected session must produce an actionable message naming the setting to
// fix, rather than the generic "check the logfile" text or a bogus "subreddit
// not found". Reddit answers an invalid cookie with 403, so this is
// deterministic regardless of whether a real session is configured locally.
func TestStartupWithInvalidSessionShowsActionableError(t *testing.T) {
	t.Logf("Testing startup with an invalid session cookie...")
	configuration := getTestConfig()
	configuration.Auth.SessionCookie = "definitely-not-a-valid-session"

	tui := components.NewRedditTui(configuration, "", "")
	tm := teatest.NewTestModel(t, tui, teatest.WithInitialTermSize(300, 100))

	t.Logf("\tVerify the session error names the setting to update...")
	// Matching single tokens, since the modal wraps longer text.
	WaitFor(t, tm, "expired", "sessionCookie", "reddit_session")
}

// The session messages are long, and the modal applies both a wrapping Width
// and a truncating MaxWidth. Guard against the actionable part being cut off in
// an ordinary 80-column terminal.
func TestSessionErrorSurvivesNarrowTerminal(t *testing.T) {
	t.Logf("Testing session error in an 80x24 terminal...")
	configuration := getTestConfig()
	configuration.Auth.SessionCookie = ""

	tui := components.NewRedditTui(configuration, "", "")
	tm := teatest.NewTestModel(t, tui, teatest.WithInitialTermSize(80, 24))

	t.Logf("\tVerify the setting name is still visible...")
	WaitFor(t, tm, "sessionCookie")
}

func WaitFor(t *testing.T, tm *teatest.TestModel, messages ...string) {
	WaitForWithInputs(t, tm, "", messages...)
}

func WaitForWithInputs(t *testing.T, tm *teatest.TestModel, inputs string, messages ...string) {
	if len(inputs) > 0 {
		tm.Send(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune(inputs),
		})
	}

	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		for _, message := range messages {
			if !bytes.Contains(bts, []byte(message)) {
				return false
			}
		}

		return true
	}, teatest.WithCheckInterval(time.Millisecond*50), teatest.WithDuration(testTimeout))
}

func getTestConfig() config.Config {
	configuration := config.NewConfig()
	configuration.Core.BypassCache = true

	domain := os.Getenv("TEST_DOMAIN")
	serverType := os.Getenv("TEST_SERVER_TYPE")
	if len(domain) > 0 && len(serverType) > 0 {
		configuration.Server.Domain = domain
		configuration.Server.Type = serverType
	}

	// old.reddit.com redirects logged-out requests to a login page, so the
	// tests below that expect real content need a session. Export
	// REDDITTUI_SESSION_COOKIE to run them.
	configuration.Auth.SessionCookie = os.Getenv("REDDITTUI_SESSION_COOKIE")

	return configuration
}
