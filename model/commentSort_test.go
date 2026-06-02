package model

import "testing"

func TestCommentSortForKey(t *testing.T) {
	tests := []struct {
		key  string
		want CommentSort
	}{
		{"1", CommentSortBest},
		{"2", CommentSortNew},
		{"3", CommentSortTop},
		{"4", CommentSortControversial},
		{"5", CommentSortOld},
	}

	for _, tt := range tests {
		got, ok := CommentSortForKey(tt.key)
		if !ok {
			t.Fatalf("CommentSortForKey(%q) not found", tt.key)
		}
		if got != tt.want {
			t.Fatalf("CommentSortForKey(%q) = %q, want %q", tt.key, got, tt.want)
		}
	}
}

func TestCommentSortQueryValue(t *testing.T) {
	if CommentSortBest.QueryValue() != "" {
		t.Fatalf("expected empty query value for best, got %q", CommentSortBest.QueryValue())
	}

	if CommentSortNew.QueryValue() != "new" {
		t.Fatalf("expected new query value, got %q", CommentSortNew.QueryValue())
	}
}
