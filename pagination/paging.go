// This package is for pagination
package pagination

import "math"

const (
	// Page size is the number of records in a page
	defaultPageSize = 20
	maxPageSize     = 100
	defaultPage     = 1
)

type PagingOptions struct {
	Page     int `json:"page" default:"1"`
	PageSize int `json:"pageSize" default:"20"`
}

// GetTotalPage calculates the total number of pages needed to display all records.
// It takes the total number of records and the desired page size as input.
// The output is the minimum number of pages required to display all records,
// given that each page displays at most 'pageSize' records.
func GetTotalPage(totalRecords int64, pageSize int) int {
	return int(math.Ceil(float64(totalRecords) / float64(pageSize)))
}

// GetPageAndPageSize validates the input page and pageSize and returns optimized values.
// If the page value is less than or equal to zero, it defaults to 'defaultPage'.
// If pageSize exceeds 'maxPageSize', it is set to 'maxPageSize'.
// This allows pageSize to be less than or equal to zero, accommodating APIs that use list API as a "ping" endpoint.
// Returns: Optimized 'page' and 'pageSize'.
func GetPageAndPageSize(page, pageSize int) (int, int) {
	if page <= 0 {
		page = defaultPage
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	return page, pageSize
}

// GetLimitOffset calculates the limit and offset for pagination based on the given page and pageSize.
// It returns the offset, which is the number of records to skip, and the limit, which is the number of records to return.
// The offset is calculated as (page - 1) * pageSize, and the limit is set to pageSize.
func GetLimitOffset(page, pageSize int) (int, int) {
	page, pageSize = GetPageAndPageSize(page, pageSize)
	return (page - 1) * int(pageSize), int(pageSize)
}
