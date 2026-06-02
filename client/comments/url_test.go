package comments

import (
	"reddittui/model"
	"testing"
)

func TestBuildCommentsUrl(t *testing.T) {
	baseURL := "https://old.reddit.com/r/dogs/comments/abc123/title/"

	tests := []struct {
		name  string
		sort  model.CommentSort
		want  string
	}{
		{
			name: "best omits sort param",
			sort: model.CommentSortBest,
			want: "https://old.reddit.com/r/dogs/comments/abc123/title/",
		},
		{
			name: "new",
			sort: model.CommentSortNew,
			want: "https://old.reddit.com/r/dogs/comments/abc123/title/?sort=new",
		},
		{
			name: "top",
			sort: model.CommentSortTop,
			want: "https://old.reddit.com/r/dogs/comments/abc123/title/?sort=top",
		},
		{
			name: "controversial",
			sort: model.CommentSortControversial,
			want: "https://old.reddit.com/r/dogs/comments/abc123/title/?sort=controversial",
		},
		{
			name: "old",
			sort: model.CommentSortOld,
			want: "https://old.reddit.com/r/dogs/comments/abc123/title/?sort=old",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildCommentsUrl(baseURL, tt.sort)
			if err != nil {
				t.Fatalf("BuildCommentsUrl() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("BuildCommentsUrl() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestStripCommentSortParam(t *testing.T) {
	rawURL := "https://old.reddit.com/r/dogs/comments/abc123/title/?sort=top&limit=500"
	got, err := StripCommentSortParam(rawURL)
	if err != nil {
		t.Fatalf("StripCommentSortParam() error = %v", err)
	}

	want := "https://old.reddit.com/r/dogs/comments/abc123/title/?limit=500"
	if got != want {
		t.Fatalf("StripCommentSortParam() = %q, want %q", got, want)
	}
}
