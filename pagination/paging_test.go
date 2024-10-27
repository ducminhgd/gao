// This package is for pagination
package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTotalPage(t *testing.T) {
	testCases := []struct {
		name              string
		totalRecords      int64
		pageSize          int
		expectedTotalPage int
	}{
		{
			name:              "totalRecords=100,pageSize=10",
			totalRecords:      100,
			pageSize:          10,
			expectedTotalPage: 10,
		},
		{
			name:              "totalRecords=105,pageSize=10",
			totalRecords:      105,
			pageSize:          10,
			expectedTotalPage: 11,
		},
		{
			name:              "totalRecords=1,pageSize=10",
			totalRecords:      1,
			pageSize:          10,
			expectedTotalPage: 1,
		},
		{
			name:              "totalRecords=0,pageSize=10",
			totalRecords:      1,
			pageSize:          10,
			expectedTotalPage: 1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			totalPage := GetTotalPage(tc.totalRecords, tc.pageSize)
			assert.Equal(t, tc.expectedTotalPage, totalPage)
		})
	}
}

func TestGetPageAndPageSize(t *testing.T) {
	type Expected struct {
		page  int
		limit int
	}
	testCases := []struct {
		name     string
		page     int
		limit    int
		expected Expected
	}{
		{
			name:  "page=0,limit=0",
			page:  0,
			limit: 0,
			expected: Expected{
				page:  1,
				limit: 0,
			},
		},
		{
			name:  "page=1,limit=1",
			page:  1,
			limit: 1,
			expected: Expected{
				page:  1,
				limit: 1,
			},
		},
		{
			name:  "page=1,limit=1000",
			page:  1,
			limit: 1_000,
			expected: Expected{
				page:  1,
				limit: 100,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			page, limit := GetPageAndPageSize(tc.page, tc.limit)
			assert.Equal(t, tc.expected.page, page)
			assert.Equal(t, tc.expected.limit, limit)
		})
	}
}

func TestGetLimitOffset(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		pageSize int
		expected struct {
			offset int
			limit  int
		}
	}{
		{
			name:     "first page",
			page:     1,
			pageSize: 10,
			expected: struct {
				offset int
				limit  int
			}{offset: 0, limit: 10},
		},
		{
			name:     "second page",
			page:     2,
			pageSize: 10,
			expected: struct {
				offset int
				limit  int
			}{offset: 10, limit: 10},
		},
		{
			name:     "last page",
			page:     3,
			pageSize: 10,
			expected: struct {
				offset int
				limit  int
			}{offset: 20, limit: 10},
		},
		{
			name:     "page size is 0",
			page:     1,
			pageSize: 0,
			expected: struct {
				offset int
				limit  int
			}{offset: 0, limit: 0},
		},
		{
			name:     "page is 0",
			page:     0,
			pageSize: 10,
			expected: struct {
				offset int
				limit  int
			}{offset: 0, limit: 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualOffset, actualLimit := GetLimitOffset(tt.page, tt.pageSize)
			assert.Equal(t, tt.expected.offset, actualOffset)
			assert.Equal(t, tt.expected.limit, actualLimit)
		})
	}
}
