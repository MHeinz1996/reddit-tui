package posts

import (
	"strings"
	"testing"

	"reddittui/client/common"

	"golang.org/x/net/html"
)

func TestSanitizePostRemovesZeroWidthFromTitle(t *testing.T) {
	const raw = `<div class="thing">
		<a class="title">Jhon Terry 😂&#8203;&#8203;🤣</a>
		<div class="likes">967</div>
		<a class="comments" href="/r/chelseafc/comments/abc/x/">42 comments</a>
		<a class="author">user</a>
		<time class="live-timestamp">2 hours ago</time>
	</div>`

	doc, err := html.Parse(strings.NewReader(raw))
	if err != nil {
		t.Fatal(err)
	}

	post := sanitizePost(OldRedditPostsParser{}.parsePost(common.HtmlNode{Node: doc}))
	if strings.Contains(post.PostTitle, "\u200b") {
		t.Fatalf("title still has zero-width space: %q", post.PostTitle)
	}
}
