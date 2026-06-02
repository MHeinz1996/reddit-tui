package model

import (
	"fmt"
	"reddittui/utils"
	"strings"
	"time"
)

type Comment struct {
	Author      string `json:"author"`
	AuthorFlair string `json:"authorFlair"`
	Text        string `json:"text"`
	Points      string `json:"points"`
	Timestamp   string `json:"timestamp"`
	Depth       int    `json:"depth"`
}

type Comments struct {
	PostTitle       string      `json:"title"`
	PostAuthor      string      `json:"author"`
	PostAuthorFlair string      `json:"postAuthorFlair"`
	Subreddit       string      `json:"subreddit"`
	PostPoints      string      `json:"points"`
	PostText        string      `json:"text"`
	PostUrl         string      `json:"url"`
	PostTimestamp   string      `json:"timestamp"`
	Sort            CommentSort `json:"sort"`
	Expiry          time.Time   `json:"expiry"`
	Comments        []Comment   `json:"comments"`
}

func (c Comment) Title() string {
	return formatDepth(c.Text, c.Depth)
}

func (c Comment) AuthorLabel() string {
	return utils.AuthorWithFlair(c.Author, c.AuthorFlair)
}

func (c Comment) Description() string {
	desc := fmt.Sprintf("%s  by %s  %s", c.Points, c.AuthorLabel(), c.Timestamp)
	return formatDepth(desc, c.Depth)
}

func (c Comments) PostAuthorLabel() string {
	return utils.AuthorWithFlair(c.PostAuthor, c.PostAuthorFlair)
}

func (c Comment) FilterValue() string {
	return c.Author
}

func formatDepth(s string, depth int) string {
	var results strings.Builder
	for range depth {
		results.WriteString("  ")
	}
	results.WriteString(s)

	return results.String()
}
