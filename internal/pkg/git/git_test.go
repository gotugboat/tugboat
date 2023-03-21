package git

import (
	"context"
	"testing"
)

func TestIsRepo(t *testing.T) {
	testCases := []struct {
		name     string
		expected bool
	}{
		{
			name:     "is repository",
			expected: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			actual := IsRepo(ctx)

			if tc.expected != actual {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}
