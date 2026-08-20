package common

import (
	"net/http"
	"net/url"
	"testing"
)

func TestIsLoginRedirect(t *testing.T) {
	tests := []struct {
		name    string
		finalUrl string
		want    bool
	}{
		{
			name:     "login redirect with reason param",
			finalUrl: "https://old.reddit.com/login/?reason=lor2&dest=https%3A%2F%2Fold.reddit.com%2Fr%2Fpolitics",
			want:     true,
		},
		{
			name:     "login path without trailing slash",
			finalUrl: "https://old.reddit.com/login",
			want:     true,
		},
		{
			name:     "normal subreddit page",
			finalUrl: "https://old.reddit.com/r/politics",
			want:     false,
		},
		{
			name:     "home page",
			finalUrl: "https://old.reddit.com",
			want:     false,
		},
		{
			name:     "subreddit that merely starts with login",
			finalUrl: "https://old.reddit.com/r/loginhelp",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := url.Parse(tt.finalUrl)
			if err != nil {
				t.Fatalf("could not parse test url: %v", err)
			}

			res := &http.Response{
				StatusCode: http.StatusOK,
				Request:    &http.Request{URL: parsed},
			}

			if got := IsLoginRedirect(res); got != tt.want {
				t.Errorf("IsLoginRedirect(%q) = %v, want %v", tt.finalUrl, got, tt.want)
			}
		})
	}
}

func TestIsLoginRedirectHandlesMissingFields(t *testing.T) {
	if IsLoginRedirect(nil) {
		t.Error("IsLoginRedirect(nil) = true, want false")
	}

	if IsLoginRedirect(&http.Response{}) {
		t.Error("IsLoginRedirect(response with no request) = true, want false")
	}

	if IsLoginRedirect(&http.Response{Request: &http.Request{}}) {
		t.Error("IsLoginRedirect(request with no url) = true, want false")
	}
}
