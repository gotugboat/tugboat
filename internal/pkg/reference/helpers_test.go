package reference

import (
	"runtime"
	"strings"
	"testing"
)

func Test_trimUri(t *testing.T) {
	testCases := []struct {
		name     string
		uri      string
		expected string
	}{
		{
			name:     "leading slash",
			uri:      "/namespace/image:tag",
			expected: "namespace/image:tag",
		},
		{
			name:     "double leading slash",
			uri:      "//image:tag",
			expected: "image:tag",
		},
		{
			name:     "triple leading slash",
			uri:      "///image:tag",
			expected: "image:tag",
		},
		{
			name:     "eight leading slash",
			uri:      "////////image:tag",
			expected: "image:tag",
		},
		{
			name:     "mixed with leading slash",
			uri:      "/registry///image:tag",
			expected: "registry/image:tag",
		},
		{
			name:     "mixed with mixed leading slash",
			uri:      "///registry///////image:tag",
			expected: "registry/image:tag",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := trimUri(tc.uri)

			if tc.expected != actual {
				t.Errorf("expected: %v; actual: %v", tc.expected, actual)
			}
		})
	}
}

func Test_generateUriString_unofficial(t *testing.T) {
	testCases := []struct {
		name     string
		image    string
		registry string
		expected string
	}{
		{
			name:     "full name",
			image:    "namespace/image:tag",
			registry: "registry",
			expected: "registry/namespace/image:tag",
		},
		{
			name:     "missing registry",
			image:    "/namespace/image:tag",
			registry: "",
			expected: "namespace/image:tag",
		},
		{
			name:     "missing registry and namespace",
			image:    "//image:tag",
			registry: "",
			expected: "image:tag",
		},
		{
			name:     "registry port and namespace",
			image:    "namespace/image:tag",
			registry: "localhost:5000",
			expected: "localhost:5000/namespace/image:tag",
		},
		{
			name:     "improper format, no registry",
			image:    "/////namespace///////image:tag",
			registry: "",
			expected: "namespace/image:tag",
		},
		{
			name:     "improper format with registry",
			image:    "/////namespace///////image:tag",
			registry: "localhost.local",
			expected: "localhost.local/namespace/image:tag",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := generateUriString(tc.image, tc.registry, "arch", false)

			if tc.expected != actual {
				t.Errorf("expected value: %v; actual value: %v", tc.expected, actual)
			}
		})
	}
}

func Test_generateUriString_official(t *testing.T) {
	testCases := []struct {
		name     string
		image    string
		registry string
		expected string
	}{
		{
			name:     "full name",
			image:    "namespace/image:tag",
			registry: "registry",
			expected: "registry/arch/image:tag",
		},
		{
			name:     "missing registry",
			image:    "/namespace/image:tag",
			registry: "",
			expected: "arch/image:tag",
		},
		{
			name:     "missing registry and namespace",
			image:    "//image:tag",
			registry: "",
			expected: "arch/image:tag",
		},
		{
			name:     "registry port and namespace",
			image:    "namespace/image:tag",
			registry: "localhost:5000",
			expected: "localhost:5000/arch/image:tag",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := generateUriString(tc.image, tc.registry, "arch", true)
			actual = patchArch(actual)

			if tc.expected != actual {
				t.Errorf("expected value: %v; actual value: %v", tc.expected, actual)
			}
		})
	}
}

// patchArch replaces the GOARCH when the unit test is running so tests can pass on any arch
func patchArch(uri string) string {
	arch := runtime.GOARCH

	if strings.Contains(uri, arch) {
		uri = strings.ReplaceAll(uri, arch, "arch")
	}
	return uri
}
