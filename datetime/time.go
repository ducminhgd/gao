package time

import "time"

// ToString converts a `time.Time` object to a RFC3399 format string
func ToString(v time.Time) string {
	return v.Format(time.RFC3339)
}
