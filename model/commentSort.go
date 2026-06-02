package model

type CommentSort string

const (
	CommentSortBest           CommentSort = "confidence"
	CommentSortNew            CommentSort = "new"
	CommentSortTop            CommentSort = "top"
	CommentSortControversial  CommentSort = "controversial"
	CommentSortOld            CommentSort = "old"
)

func (s CommentSort) Label() string {
	switch s {
	case CommentSortBest, "":
		return "Best"
	case CommentSortNew:
		return "New"
	case CommentSortTop:
		return "Top"
	case CommentSortControversial:
		return "Controversial"
	case CommentSortOld:
		return "Old"
	default:
		return "Best"
	}
}

func (s CommentSort) QueryValue() string {
	if s == "" || s == CommentSortBest {
		return ""
	}
	return string(s)
}

func CommentSortForKey(key string) (CommentSort, bool) {
	switch key {
	case "1":
		return CommentSortBest, true
	case "2":
		return CommentSortNew, true
	case "3":
		return CommentSortTop, true
	case "4":
		return CommentSortControversial, true
	case "5":
		return CommentSortOld, true
	default:
		return "", false
	}
}
