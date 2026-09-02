package service

import (
	"testing"
)

func TestArticleListCacheKey(t *testing.T) {
	catID := uint(42)
	tests := []struct {
		name   string
		filter PublishedArticleFilter
		want   string
	}{
		{
			name:   "defaults",
			filter: PublishedArticleFilter{Page: 1},
			want:   "articles:list:1::::0",
		},
		{
			name:   "with keyword",
			filter: PublishedArticleFilter{Page: 2, Keyword: "go"},
			want:   "articles:list:2:go:::0",
		},
		{
			name:   "with category",
			filter: PublishedArticleFilter{Page: 1, CategoryID: &catID},
			want:   "articles:list:1::42::0",
		},
		{
			name:   "with tag",
			filter: PublishedArticleFilter{Page: 1, Tag: "docker"},
			want:   "articles:list:1:::docker:0",
		},
		{
			name:   "with author",
			filter: PublishedArticleFilter{Page: 1, AuthorID: 7},
			want:   "articles:list:1::::7",
		},
		{
			name: "all fields",
			filter: PublishedArticleFilter{
				Page:       3,
				Keyword:    "rust",
				CategoryID: &catID,
				Tag:        "web",
				AuthorID:   5,
			},
			want: "articles:list:3:rust:42:web:5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := articleListCacheKey(tt.filter)
			if got != tt.want {
				t.Errorf("articleListCacheKey() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNormalizePagination(t *testing.T) {
	tests := []struct {
		name       string
		page       int
		pageSize   int
		wantPage   int
		wantSize   int
	}{
		{"zero values", 0, 0, 1, 10},
		{"negative page", -1, -5, 1, 10},
		{"valid values", 3, 20, 3, 20},
		{"page size too large", 1, 100, 1, 30},
		{"page size exactly 30", 1, 30, 1, 30},
		{"page size 1", 1, 1, 1, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizePagination(tt.page, tt.pageSize)
			if got.Page != tt.wantPage || got.PageSize != tt.wantSize {
				t.Errorf("normalizePagination(%d, %d) = {Page:%d, PageSize:%d}, want {Page:%d, PageSize:%d}",
					tt.page, tt.pageSize, got.Page, got.PageSize, tt.wantPage, tt.wantSize)
			}
		})
	}
}

func TestParsePagination(t *testing.T) {
	tests := []struct {
		name       string
		pageRaw    string
		sizeRaw    string
		wantPage   int
		wantSize   int
	}{
		{"valid", "2", "20", 2, 20},
		{"empty strings", "", "", 0, 0},
		{"non-numeric", "abc", "xyz", 0, 0},
		{"mixed", "5", "", 5, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPage, gotSize := ParsePagination(tt.pageRaw, tt.sizeRaw)
			if gotPage != tt.wantPage || gotSize != tt.wantSize {
				t.Errorf("ParsePagination(%q, %q) = (%d, %d), want (%d, %d)",
					tt.pageRaw, tt.sizeRaw, gotPage, gotSize, tt.wantPage, tt.wantSize)
			}
		})
	}
}
