package common

import (
	"net/http"
	"strings"
)

const loginPathSegment = "/login"

// IsLoginRedirect reports whether a response is really a login page.
//
// Reddit answers logged-out requests with a 302 to /login/?reason=lor2 rather
// than a 4xx. Go's http.Client follows that redirect transparently, so the
// response arrives as a genuine 200 carrying login HTML. Without this check the
// parsers just find zero posts and report "not found", which is what made an
// expired session look like a missing subreddit.
//
// res.Request is the *final* request in the redirect chain, so its URL is where
// we actually landed rather than where we aimed.
func IsLoginRedirect(res *http.Response) bool {
	if res == nil || res.Request == nil || res.Request.URL == nil {
		return false
	}

	path := res.Request.URL.EscapedPath()
	return path == loginPathSegment || strings.HasPrefix(path, loginPathSegment+"/")
}

// AuthErrorForStatus maps authentication-related status codes to sentinel
// errors, returning nil for statuses that are not auth-related.
//
// Reddit answers a *rejected* session cookie with 403 "Blocked" rather than a
// login redirect, so 403 has to be treated as a probable auth failure — an
// expired cookie is otherwise indistinguishable from a generic fetch error.
func AuthErrorForStatus(statusCode int) error {
	switch statusCode {
	case http.StatusUnauthorized:
		return ErrNotAuthenticated
	case http.StatusForbidden:
		return ErrForbidden
	default:
		return nil
	}
}
