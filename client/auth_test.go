package client

import (
	"net/url"
	"testing"
)

func TestParseSessionCookie(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantNames  []string
		wantValues []string
	}{
		{
			name:       "bare cookie value",
			input:      "123456,2026-08-20T00:00:00,abcdef0123456789",
			wantNames:  []string{sessionCookieName},
			wantValues: []string{"123456,2026-08-20T00:00:00,abcdef0123456789"},
		},
		{
			name:       "bare value with surrounding whitespace",
			input:      "  abc123  ",
			wantNames:  []string{sessionCookieName},
			wantValues: []string{"abc123"},
		},
		{
			name:       "full cookie header pasted from devtools",
			input:      "reddit_session=abc123; token_v2=xyz789",
			wantNames:  []string{"reddit_session", "token_v2"},
			wantValues: []string{"abc123", "xyz789"},
		},
		{
			name:       "single name=value pair",
			input:      "reddit_session=abc123",
			wantNames:  []string{"reddit_session"},
			wantValues: []string{"abc123"},
		},
		{
			name:       "empty input yields nothing",
			input:      "   ",
			wantNames:  nil,
			wantValues: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cookies := parseSessionCookie(tt.input)

			if len(cookies) != len(tt.wantNames) {
				t.Fatalf("got %d cookies, want %d", len(cookies), len(tt.wantNames))
			}

			for i, cookie := range cookies {
				if cookie.Name != tt.wantNames[i] {
					t.Errorf("cookie %d name = %q, want %q", i, cookie.Name, tt.wantNames[i])
				}

				if cookie.Value != tt.wantValues[i] {
					t.Errorf("cookie %d value = %q, want %q", i, cookie.Value, tt.wantValues[i])
				}

				if cookie.Domain != cookieDomain {
					t.Errorf("cookie %d domain = %q, want %q", i, cookie.Domain, cookieDomain)
				}
			}
		})
	}
}

// A bare reddit_session value can itself contain "=", so it must not be
// mistaken for a name=value pair.
func TestParseSessionCookieBareValueContainingEquals(t *testing.T) {
	cookies := parseSessionCookie("abc,def=,ghi")

	if len(cookies) != 1 {
		t.Fatalf("got %d cookies, want 1", len(cookies))
	}

	if cookies[0].Name != sessionCookieName {
		t.Errorf("name = %q, want %q", cookies[0].Name, sessionCookieName)
	}

	if cookies[0].Value != "abc,def=,ghi" {
		t.Errorf("value = %q, want the whole bare value", cookies[0].Value)
	}
}

func TestIsRedditHost(t *testing.T) {
	tests := []struct {
		host string
		want bool
	}{
		{"reddit.com", true},
		{"old.reddit.com", true},
		{"www.reddit.com", true},
		{"OLD.REDDIT.COM", true},
		{"old.reddit.com.", true},
		{"redlib.catsarch.com", false},
		{"notreddit.com", false},
		{"reddit.com.evil.example", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.host, func(t *testing.T) {
			if got := isRedditHost(tt.host); got != tt.want {
				t.Errorf("isRedditHost(%q) = %v, want %v", tt.host, got, tt.want)
			}
		})
	}
}

func TestNewSessionCookieJarAttachesCookieToRedditHost(t *testing.T) {
	baseUrl := "https://old.reddit.com"
	jar := newSessionCookieJar(baseUrl, "abc123")

	if jar == nil {
		t.Fatal("jar is nil")
	}

	parsed, err := url.Parse(baseUrl)
	if err != nil {
		t.Fatalf("could not parse base url: %v", err)
	}

	cookies := jar.Cookies(parsed)
	if len(cookies) != 1 {
		t.Fatalf("got %d cookies for reddit host, want 1", len(cookies))
	}

	if cookies[0].Name != sessionCookieName || cookies[0].Value != "abc123" {
		t.Errorf("got cookie %s=%s, want %s=abc123",
			cookies[0].Name, cookies[0].Value, sessionCookieName)
	}
}

// The session cookie grants full account access, so it must never be sent to a
// user-configured non-reddit server such as a redlib instance.
func TestNewSessionCookieJarWithholdsCookieFromNonRedditHost(t *testing.T) {
	baseUrl := "https://redlib.example.com"
	jar := newSessionCookieJar(baseUrl, "abc123")

	if jar == nil {
		t.Fatal("jar is nil")
	}

	parsed, err := url.Parse(baseUrl)
	if err != nil {
		t.Fatalf("could not parse base url: %v", err)
	}

	if cookies := jar.Cookies(parsed); len(cookies) != 0 {
		t.Errorf("got %d cookies for non-reddit host, want 0 (credential leak)", len(cookies))
	}
}

func TestNewSessionCookieJarWithoutCookieIsUsable(t *testing.T) {
	jar := newSessionCookieJar("https://old.reddit.com", "")
	if jar == nil {
		t.Fatal("jar is nil; an unauthenticated client still needs a usable jar")
	}

	parsed, _ := url.Parse("https://old.reddit.com")
	if cookies := jar.Cookies(parsed); len(cookies) != 0 {
		t.Errorf("got %d cookies with no cookie configured, want 0", len(cookies))
	}
}

func TestNewSessionCookieJarCoversSubdomainRedirect(t *testing.T) {
	jar := newSessionCookieJar("https://old.reddit.com", "abc123")

	// Reddit can bounce between old. and www. hosts; the cookie should follow.
	parsed, err := url.Parse("https://www.reddit.com/r/politics")
	if err != nil {
		t.Fatalf("could not parse url: %v", err)
	}

	if cookies := jar.Cookies(parsed); len(cookies) != 1 {
		t.Errorf("got %d cookies for www.reddit.com, want 1", len(cookies))
	}
}
