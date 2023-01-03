// This package is for pagination
package pagination

import "math"

const (
	// Page size is the number of records in a page
	defaultPageSize = 20
	maxPageSize     = 100
	defaultPage     = 1
)

// GetTotalPage calculates number of pages from number of records and page size.
func GetTotalPage(totalRecords int64, pageSize int32) int32 {
	return int32(math.Ceil(float64(totalRecords) / float64(pageSize)))
}

// GetPageAndPageSize validates and returns page size and limit
// `pageSize` is allowed to be less than or equal to zero in case some API using list API is a "ping" endpoint
func GetPageAndPageSize(page, pageSize int32) (int32, int32) {
	if page <= 0 {
		page = defaultPage
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	return page, pageSize
}
