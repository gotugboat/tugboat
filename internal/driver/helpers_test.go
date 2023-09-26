package driver

import (
	"testing"
	"tugboat/internal/pkg/reference"
)

func TestGenerateUri(t *testing.T) {
	testCases := []struct {
		name       string
		registry   string
		namespace  string
		image      string
		archOption string
		expected   string
	}{
		{
			name:       "full uri string provided",
			registry:   "localhost.local",
			namespace:  "namespace",
			image:      "localhost.local/namespace/busybox:latest",
			archOption: "omit",
			expected:   "localhost.local/namespace/busybox:latest",
		},
	}

	isOfficial := false

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uri, err := GenerateUri(tc.registry, tc.namespace, tc.image, isOfficial, reference.ArchOption(tc.archOption))
			if err != nil {
				t.Errorf("An unexpected error occurred: %v", err)
			}

			actual := uri.Remote()

			if tc.expected != actual {
				t.Errorf("expected value: %v; actual value: %v", tc.expected, actual)
			}
		})
	}
}

func TestGenerateUriWithArch(t *testing.T) {
	testCases := []struct {
		name         string
		registry     string
		namespace    string
		image        string
		archOption   string
		architecture string
		expected     string
	}{
		{
			name:         "full uri string provided",
			registry:     "localhost.local",
			namespace:    "namespace",
			image:        "localhost.local/namespace/busybox:latest",
			archOption:   "omit",
			architecture: "arch",
			expected:     "localhost.local/namespace/busybox:latest",
		},
	}

	isOfficial := false

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uri, err := GenerateUriWithArch(tc.registry, tc.namespace, tc.image, isOfficial, reference.ArchOption(tc.archOption), tc.architecture)
			if err != nil {
				t.Errorf("An unexpected error occurred: %v", err)
			}

			actual := uri.Remote()

			if tc.expected != actual {
				t.Errorf("expected value: %v; actual value: %v", tc.expected, actual)
			}
		})
	}
}
