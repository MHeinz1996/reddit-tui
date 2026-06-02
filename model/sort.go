package model

type Sort string

const (
	SortHot           Sort = "hot"
	SortNew           Sort = "new"
	SortTop           Sort = "top"
	SortRising        Sort = "rising"
	SortControversial Sort = "controversial"
)

func (s Sort) Label() string {
	switch s {
	case SortHot:
		return "Hot"
	case SortNew:
		return "New"
	case SortTop:
		return "Top"
	case SortRising:
		return "Rising"
	case SortControversial:
		return "Controversial"
	default:
		return "Hot"
	}
}

func (s Sort) PathSegment() string {
	if s == "" || s == SortHot {
		return ""
	}
	return string(s)
}

func SortForKey(key string) (Sort, bool) {
	switch key {
	case "1":
		return SortHot, true
	case "2":
		return SortNew, true
	case "3":
		return SortRising, true
	case "4":
		return SortTop, true
	case "5":
		return SortControversial, true
	default:
		return "", false
	}
}
