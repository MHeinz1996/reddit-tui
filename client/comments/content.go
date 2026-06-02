package comments

import (
	"reddittui/client/common"
	"strings"
)

func renderBodyContent(contentRoot, mdNode common.HtmlNode) string {
	if raw := findMarkdownSource(contentRoot); raw != "" {
		if rendered, err := common.RenderMarkdown(raw); err == nil && strings.TrimSpace(rendered) != "" {
			return normalizeContent(rendered)
		}
	}

	return normalizeContent(renderHtmlNode(mdNode))
}

func findMarkdownSource(contentRoot common.HtmlNode) string {
	if textarea, ok := contentRoot.FindDescendant("textarea"); ok {
		return strings.TrimSpace(textarea.TextContent())
	}

	return ""
}

func normalizeContent(text string) string {
	return postTextTrimRegex.ReplaceAllString(strings.TrimSpace(text), "\n\n")
}
