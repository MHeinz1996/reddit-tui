package common

import (
	"strings"

	"github.com/charmbracelet/glamour"
)

const DefaultMarkdownWidth = 100

func RenderMarkdown(source string) (string, error) {
	source = strings.TrimSpace(source)
	if source == "" {
		return "", nil
	}

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(DefaultMarkdownWidth),
	)
	if err != nil {
		return "", err
	}

	rendered, err := renderer.Render(source)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(rendered), nil
}
