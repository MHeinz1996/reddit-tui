package posts

import (
	"testing"

	"reddittui/model"
)

func TestDedupePosts(t *testing.T) {
	posts := []model.Post{
		{PostTitle: "Daily Discussion", CommentsUrl: "https://reddit.com/a"},
		{PostTitle: "Daily Discussion", CommentsUrl: "https://reddit.com/a"},
		{PostTitle: "Other post", CommentsUrl: "https://reddit.com/b"},
	}

	got := dedupePosts(posts)
	if len(got) != 2 {
		t.Fatalf("dedupePosts() len = %d, want 2", len(got))
	}
	if got[0].CommentsUrl != "https://reddit.com/a" || got[1].CommentsUrl != "https://reddit.com/b" {
		t.Fatalf("dedupePosts() = %+v", got)
	}
}
