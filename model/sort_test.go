package model

import "testing"

func TestSortForKey(t *testing.T) {
	tests := []struct {
		key  string
		want Sort
	}{
		{"1", SortHot},
		{"2", SortNew},
		{"3", SortRising},
		{"4", SortTop},
		{"5", SortControversial},
	}

	for _, tt := range tests {
		got, ok := SortForKey(tt.key)
		if !ok {
			t.Fatalf("SortForKey(%q) not found", tt.key)
		}
		if got != tt.want {
			t.Fatalf("SortForKey(%q) = %q, want %q", tt.key, got, tt.want)
		}
	}

	if _, ok := SortForKey("6"); ok {
		t.Fatal("expected SortForKey(6) to be false")
	}
}

func TestSortPathSegment(t *testing.T) {
	if SortHot.PathSegment() != "" {
		t.Fatalf("expected empty path for hot, got %q", SortHot.PathSegment())
	}

	if SortNew.PathSegment() != "new" {
		t.Fatalf("expected new path segment, got %q", SortNew.PathSegment())
	}
}
