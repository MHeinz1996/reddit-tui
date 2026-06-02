package common

import (
	"strings"
	"testing"
)

func TestRenderMarkdownList(t *testing.T) {
	raw := "Shopping List:\n\n* Milk\n* Eggs\n* Chicken"
	got, err := RenderMarkdown(raw)
	if err != nil {
		t.Fatalf("RenderMarkdown() error = %v", err)
	}

	for _, item := range []string{"Shopping List:", "Milk", "Eggs", "Chicken"} {
		if !strings.Contains(got, item) {
			t.Fatalf("RenderMarkdown() = %q, expected to contain %q", got, item)
		}
	}
}

func TestRenderMarkdownEmpty(t *testing.T) {
	got, err := RenderMarkdown("")
	if err != nil {
		t.Fatalf("RenderMarkdown() error = %v", err)
	}
	if got != "" {
		t.Fatalf("RenderMarkdown() = %q, want empty string", got)
	}
}
