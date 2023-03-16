package slices

import "testing"

func TestContains(t *testing.T) {
	testCases := []struct {
		name     string
		slice    []string
		search   string
		expected bool
	}{
		{
			name:     "empty slice",
			slice:    []string{},
			search:   "foo",
			expected: false,
		},
		{
			name:     "without value",
			slice:    []string{"foo"},
			search:   "bar",
			expected: false,
		},
		{
			name:     "in first value",
			slice:    []string{"foo"},
			search:   "foo",
			expected: true,
		},
		{
			name:     "in second value",
			slice:    []string{"foo", "bar"},
			search:   "bar",
			expected: true,
		},
		{
			name:     "nil",
			slice:    nil,
			search:   "bar",
			expected: false,
		},
		{
			name:     "empty search value",
			slice:    []string{"foo", "bar"},
			search:   "",
			expected: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if value := Contains(tc.slice, tc.search); value != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, value)
			}
		})
	}
}
