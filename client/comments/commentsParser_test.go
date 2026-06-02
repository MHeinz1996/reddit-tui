package comments

import (
	"reddittui/client/common"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func renderMdHTML(htmlSnippet string) string {
	doc, err := html.Parse(strings.NewReader("<html><body>" + htmlSnippet + "</body></html>"))
	if err != nil {
		panic(err)
	}

	mdNode, ok := common.HtmlNode{Node: doc}.FindDescendant("div", "md")
	if !ok {
		panic("md node not found")
	}

	return strings.TrimSpace(renderHtmlNode(mdNode))
}

func TestRenderHtmlNodeUnorderedList(t *testing.T) {
	got := renderMdHTML(`<div class="md"><ul><li><p>First item</p></li><li><p>Second item</p></li></ul></div>`)
	want := "- First item\n- Second item"

	if got != want {
		t.Fatalf("renderHtmlNode() = %q, want %q", got, want)
	}
}

func TestRenderHtmlNodeOrderedList(t *testing.T) {
	got := renderMdHTML(`<div class="md"><ol><li><p>Step one</p></li><li><p>Step two</p></li></ol></div>`)
	want := "1. Step one\n2. Step two"

	if got != want {
		t.Fatalf("renderHtmlNode() = %q, want %q", got, want)
	}
}

func TestRenderHtmlNodeListWithPlainTextItem(t *testing.T) {
	got := renderMdHTML(`<div class="md"><ul><li>Plain item</li></ul></div>`)
	want := "- Plain item"

	if got != want {
		t.Fatalf("renderHtmlNode() = %q, want %q", got, want)
	}
}

func TestRenderHtmlNodeMixedContent(t *testing.T) {
	got := renderMdHTML(`<div class="md"><p>Intro</p><ul><li><p>Bullet</p></li></ul></div>`)
	want := "Intro\n- Bullet"

	if got != want {
		t.Fatalf("renderHtmlNode() = %q, want %q", got, want)
	}
}

func TestRenderHtmlNodeSkipsEmptyParagraphs(t *testing.T) {
	got := renderMdHTML(`<div class="md"><p>Shopping List:</p><p></p><p> </p><p></p></div>`)
	want := "Shopping List:"

	if got != want {
		t.Fatalf("renderHtmlNode() = %q, want %q", got, want)
	}
}

func TestRenderHtmlNodeSkipsEmptyListItems(t *testing.T) {
	got := renderMdHTML(`<div class="md"><p>Shopping List:</p><ul><li><p></p></li><li><p></p></li></ul></div>`)
	want := "Shopping List:"

	if got != want {
		t.Fatalf("renderHtmlNode() = %q, want %q", got, want)
	}
}

func TestRenderBodyContentPrefersMarkdown(t *testing.T) {
	htmlSnippet := `
<div class="entry">
  <form class="usertext">
    <div class="usertext-edit md-container" style="display:none">
      <textarea>Shopping List:

* Milk
* Eggs
* Chicken</textarea>
    </div>
    <div class="usertext-body may-blank-within md-container">
      <div class="md">
        <p>Shopping List:</p>
        <ul><li><p></p></li><li><p></p></li><li><p></p></li></ul>
      </div>
    </div>
  </form>
</div>`

	doc, err := html.Parse(strings.NewReader("<html><body>" + htmlSnippet + "</body></html>"))
	if err != nil {
		t.Fatalf("html.Parse() error = %v", err)
	}

	entryNode, ok := common.HtmlNode{Node: doc}.FindDescendant("div", "entry")
	if !ok {
		t.Fatal("entry node not found")
	}

	mdNode, ok := entryNode.FindDescendant("div", "md")
	if !ok {
		t.Fatal("md node not found")
	}

	got := renderBodyContent(entryNode, mdNode)
	for _, item := range []string{"Shopping List:", "Milk", "Eggs", "Chicken"} {
		if !strings.Contains(got, item) {
			t.Fatalf("renderBodyContent() = %q, expected to contain %q", got, item)
		}
	}
}

func TestRenderBodyContentFallsBackToHTML(t *testing.T) {
	htmlSnippet := `<div class="entry"><div class="md"><p>Hello world</p></div></div>`

	doc, err := html.Parse(strings.NewReader("<html><body>" + htmlSnippet + "</body></html>"))
	if err != nil {
		t.Fatalf("html.Parse() error = %v", err)
	}

	entryNode, ok := common.HtmlNode{Node: doc}.FindDescendant("div", "entry")
	if !ok {
		t.Fatal("entry node not found")
	}

	mdNode, ok := entryNode.FindDescendant("div", "md")
	if !ok {
		t.Fatal("md node not found")
	}

	got := renderBodyContent(entryNode, mdNode)
	if got != "Hello world" {
		t.Fatalf("renderBodyContent() = %q, want %q", got, "Hello world")
	}
}
