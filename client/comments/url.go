package comments

import (
	"net/url"
	"reddittui/model"
)

func StripCommentSortParam(rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	query := parsed.Query()
	query.Del("sort")
	parsed.RawQuery = query.Encode()

	return parsed.String(), nil
}

func BuildCommentsUrl(baseURL string, sort model.CommentSort) (string, error) {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	query := parsed.Query()
	query.Del("sort")
	if sortValue := sort.QueryValue(); sortValue != "" {
		query.Set("sort", sortValue)
	}
	parsed.RawQuery = query.Encode()

	return parsed.String(), nil
}
