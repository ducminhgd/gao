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
		pageSize          int32
		expectedTotalPage int32
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
		page  int32
		limit int32
	}
	testCases := []struct {
		name     string
		page     int32
		limit    int32
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
