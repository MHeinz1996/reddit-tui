package common

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestUserFlairFromTagline(t *testing.T) {
	const raw = `<p class="tagline">
		<a class="author">renome</a>
		<span class="flair flair-player" title="Celery">Celery</span>
		<span class="userattrs"></span>
	</p>`

	doc, err := html.Parse(strings.NewReader(raw))
	if err != nil {
		t.Fatal(err)
	}

	tagline, ok := HtmlNode{Node: doc}.FindDescendant("p", "tagline")
	if !ok {
		t.Fatal("tagline not found")
	}

	if got := UserFlairFromTagline(tagline); got != "Celery" {
		t.Fatalf("UserFlairFromTagline() = %q, want Celery", got)
	}
}

func TestUserFlairFromTaglineRichText(t *testing.T) {
	const raw = `<p class="tagline">
		<a class="author">MarinaGranovskaia</a>
		<span class="flairrichtext flaircolorlight flair" title="Palmer :C_Palmer:">
			<span>Palmer </span>
		</span>
	</p>`

	doc, err := html.Parse(strings.NewReader(raw))
	if err != nil {
		t.Fatal(err)
	}

	tagline, ok := HtmlNode{Node: doc}.FindDescendant("p", "tagline")
	if !ok {
		t.Fatal("tagline not found")
	}

	got := UserFlairFromTagline(tagline)
	if !strings.Contains(got, "Palmer") {
		t.Fatalf("UserFlairFromTagline() = %q, want text containing Palmer", got)
	}
}
