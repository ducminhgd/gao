package datetime

import (
	"testing"
	"time"
)

func TestToString(t *testing.T) {
	t.Parallel()

	// test case for ToString function
	testCases := []struct {
		name   string
		time   time.Time
		expect string
	}{
		{
			name:   "standard time",
			time:   time.Date(2021, 11, 10, 23, 0, 0, 0, time.UTC),
			expect: "2021-11-10T23:00:00Z",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := ToString(tc.time)
			if output != tc.expect {
				t.Errorf("ToString(%v): expected %v, got %v", tc.time, tc.expect, output)
			}
		})
	}
}
