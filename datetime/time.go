package datetime

import "time"

// ToString formats a time value as a string using RFC3339 format.
//
// It takes a time.Time parameter `v` and returns a string.
func ToString(v time.Time) string {
	return v.Format(time.RFC3339)
}
