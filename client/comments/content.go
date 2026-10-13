package comments

import (
	"reddittui/client/common"
	"reddittui/utils"
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
	text = utils.SanitizeDisplayText(text)
	return postTextTrimRegex.ReplaceAllString(strings.TrimSpace(text), "\n\n")
}
