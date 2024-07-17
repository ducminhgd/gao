package db

import "testing"

func TestOrderOpts_RawSQL(t *testing.T) {
	testCases := []struct {
		name   string
		orders []string
		sep    string
		expect string
	}{
		{
			name:   "empty order opts",
			orders: []string{},
			sep:    ".",
			expect: "",
		},
		{
			name:   "single order opt with ASC",
			orders: []string{"foo.asc"},
			sep:    ".",
			expect: "ORDER BY foo ASC",
		},
		{
			name:   "single order opt with DESC",
			orders: []string{"foo.desc"},
			sep:    ".",
			expect: "ORDER BY foo DESC",
		},
		{
			name:   "multiple order opts",
			orders: []string{"foo.asc", "bar.desc", "baz"},
			sep:    ".",
			expect: "ORDER BY foo ASC, bar DESC, baz",
		},
		{
			name:   "custom separator",
			orders: []string{"foo:asc", "bar:desc", "baz"},
			sep:    ":",
			expect: "ORDER BY foo ASC, bar DESC, baz",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := RawOrderBySQL(tc.orders, tc.sep)
			if output != tc.expect {
				t.Errorf("RawSQL(%v, %v): expected %v, got %v", tc.orders, tc.sep, tc.expect, output)
			}
		})
	}
}
