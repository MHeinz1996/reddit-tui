package model

import (
	"fmt"
	"reddittui/utils"
	"strings"
	"time"
)

type Post struct {
	PostTitle     string    `json:"title"`
	Author        string    `json:"author"`
	AuthorFlair   string    `json:"authorFlair"`
	Subreddit     string    `json:"subreddit"`
	FriendlyDate  string    `json:"friendlyDate"`
	Expiry        time.Time `json:"expiry"`
	PostUrl       string    `json:"postUrl"`
	CommentsUrl   string    `json:"commentsUrl"`
	TotalComments string    `json:"totalComments"`
	TotalLikes    string    `json:"totalLikes"`
}

type Posts struct {
	Description string
	Subreddit   string
	IsHome      bool
	Sort        Sort
	Posts       []Post
	After       string
	Expiry      time.Time
}

func (p Post) AuthorLabel() string {
	return utils.AuthorWithFlair(p.Author, p.AuthorFlair)
}

func (p Post) Title() string {
	return fmt.Sprintf(" %s  %s", p.TotalLikes, p.PostTitle)
}

func (p Post) Description() string {
	var sb strings.Builder
	if strings.TrimSpace(p.Subreddit) != "" {
		sb.WriteString(p.Subreddit)
		sb.WriteString("  ")
	}

	if strings.TrimSpace(p.TotalComments) == "" {
		fmt.Fprintf(&sb, "%d comments  ", 0)
	} else {
		fmt.Fprintf(&sb, "%s comments  ", p.TotalComments)
	}

	fmt.Fprintf(&sb, "submitted %s by %s", p.FriendlyDate, p.AuthorLabel())
	return sb.String()
}

func (p Post) FilterValue() string {
	return p.PostTitle
}
