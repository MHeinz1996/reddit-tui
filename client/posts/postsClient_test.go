package posts

import (
	"reddittui/model"
	"testing"
)

func TestBuildPostsUrl(t *testing.T) {
	client := RedditPostsClient{BaseUrl: "https://old.reddit.com"}

	tests := []struct {
		name      string
		subreddit string
		after     string
		sort      model.Sort
		want      string
	}{
		{
			name: "home hot",
			sort: model.SortHot,
			want: "https://old.reddit.com",
		},
		{
			name: "home new",
			sort: model.SortNew,
			want: "https://old.reddit.com/new",
		},
		{
			name:  "home top with pagination",
			sort:  model.SortTop,
			after: "t3_abc",
			want:  "https://old.reddit.com/top?after=t3_abc",
		},
		{
			name:      "subreddit hot",
			subreddit: "dogs",
			sort:      model.SortHot,
			want:      "https://old.reddit.com/r/dogs",
		},
		{
			name:      "subreddit rising",
			subreddit: "dogs",
			sort:      model.SortRising,
			want:      "https://old.reddit.com/r/dogs/rising",
		},
		{
			name:      "subreddit controversial with pagination",
			subreddit: "dogs",
			sort:      model.SortControversial,
			after:     "t3_xyz",
			want:      "https://old.reddit.com/r/dogs/controversial?after=t3_xyz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := client.BuildPostsUrl(tt.subreddit, tt.after, tt.sort)
			if got != tt.want {
				t.Fatalf("BuildPostsUrl() = %q, want %q", got, tt.want)
			}
		})
	}
}
