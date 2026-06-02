package common

import "strings"

// UserFlairFromTagline extracts the author's user flair from an old.reddit tagline node.
func UserFlairFromTagline(tagline HtmlNode) string {
	for child := range tagline.ChildNodes() {
		node := HtmlNode{Node: child}
		if !isUserFlairSpan(node) {
			continue
		}

		if text := strings.TrimSpace(node.TextContent()); text != "" {
			return text
		}

		if title := strings.TrimSpace(node.GetAttr("title")); title != "" {
			return title
		}
	}

	return ""
}

func isUserFlairSpan(n HtmlNode) bool {
	if !n.TagEquals("span") {
		return false
	}

	for _, class := range n.Classes() {
		switch class {
		case "flair", "flairrichtext":
			return true
		default:
			if strings.HasPrefix(class, "flair-") {
				return true
			}
		}
	}

	return false
}
