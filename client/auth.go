package client

import (
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

const (
	sessionCookieName = "reddit_session"
	redditDomain      = "reddit.com"
	// cookieDomain is deliberately the parent domain so the session survives a
	// redirect hop between old.reddit.com and www.reddit.com.
	cookieDomain = "." + redditDomain
)

// newSessionCookieJar builds a cookie jar seeded with the configured reddit
// session cookie.
//
// The cookie value is never logged: it grants full access to the account.
func newSessionCookieJar(baseUrl, sessionCookie string) http.CookieJar {
	jar, err := cookiejar.New(nil)
	if err != nil {
		slog.Warn("Could not create cookie jar, continuing unauthenticated", "error", err)
		return nil
	}

	if sessionCookie == "" {
		return jar
	}

	parsed, err := url.Parse(baseUrl)
	if err != nil {
		slog.Warn("Could not parse base url, session cookie not attached", "error", err)
		return jar
	}

	// Guard against leaking a reddit account credential to a third party. The
	// server domain is user-configurable, so it may point at a redlib instance.
	if !isRedditHost(parsed.Hostname()) {
		slog.Warn(
			"Session cookie configured but server is not a reddit.com host; cookie not sent",
			"host", parsed.Hostname())
		return jar
	}

	cookies := parseSessionCookie(sessionCookie)
	if len(cookies) == 0 {
		slog.Warn("Configured session cookie is empty after parsing, continuing unauthenticated")
		return jar
	}

	jar.SetCookies(parsed, cookies)
	slog.Debug("Attached reddit session cookie", "host", parsed.Hostname(), "count", len(cookies))
	return jar
}

func isRedditHost(host string) bool {
	host = strings.ToLower(strings.TrimSuffix(host, "."))
	return host == redditDomain || strings.HasSuffix(host, "."+redditDomain)
}

// parseSessionCookie accepts either a bare reddit_session value or a full
// "name=value; name2=value2" cookie header pasted out of browser devtools.
func parseSessionCookie(sessionCookie string) []*http.Cookie {
	sessionCookie = strings.TrimSpace(sessionCookie)
	if sessionCookie == "" {
		return nil
	}

	if looksLikeCookieHeader(sessionCookie) {
		if parsed, err := http.ParseCookie(sessionCookie); err == nil && len(parsed) > 0 {
			for _, cookie := range parsed {
				cookie.Domain = cookieDomain
				cookie.Path = "/"
			}
			return parsed
		}
		slog.Warn("Could not parse session cookie as a cookie header, treating it as a bare value")
	}

	return []*http.Cookie{{
		Name:   sessionCookieName,
		Value:  sessionCookie,
		Domain: cookieDomain,
		Path:   "/",
	}}
}

// looksLikeCookieHeader distinguishes "reddit_session=abc; other=def" from a
// bare cookie value. A bare value can itself contain "=", so require the text
// before the first "=" to be a plausible cookie name.
func looksLikeCookieHeader(s string) bool {
	name, _, found := strings.Cut(s, "=")
	if !found {
		return false
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return false
	}

	return strings.IndexFunc(name, func(r rune) bool {
		return r == ',' || r == ';' || r == ' ' || r == '\t'
	}) == -1
}
