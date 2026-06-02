package utils

import "testing"

func TestAuthorWithFlair(t *testing.T) {
	if got := AuthorWithFlair("alice", "Mod"); got != "alice [Mod]" {
		t.Fatalf("AuthorWithFlair() = %q", got)
	}
	if got := AuthorWithFlair("alice", ""); got != "alice" {
		t.Fatalf("AuthorWithFlair() = %q", got)
	}
}
