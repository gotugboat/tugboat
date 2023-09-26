package reference

import (
	"fmt"
	"testing"
)

func TestNewUri_unofficial_namespace(t *testing.T) {
	testCases := []struct {
		name       string
		registry   string
		namespace  string
		image      string
		archOption ArchOption
		expected   string
	}{
		// registry, namespace
		{
			name:       "registry, duplicate namespace",
			registry:   "localhost.local",
			namespace:  "namespace",
			image:      "namespace/busybox",
			archOption: ArchOmit,
			expected:   "localhost.local/namespace/busybox:latest",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := &UriOptions{
				Registry:   tc.registry,
				Official:   false,
				ArchOption: tc.archOption,
			}
			uri, err := NewUri(fmt.Sprintf("%s/%s", tc.namespace, tc.image), opts)
			if err != nil {
				t.Errorf("An unexpected error occurred: %v", err)
			}

			actual := patchArch(uri.Remote())

			if tc.expected != actual {
				t.Errorf("expected value: %v; actual value: %v", tc.expected, actual)
			}
		})
	}
}

func TestNewUri_official(t *testing.T) {
	testCases := []struct {
		name       string
		registry   string
		image      string
		archOption ArchOption
		expected   string
	}{
		// registry, namespace
		{
			name:       "registry, namespace, arch append",
			registry:   "localhost.local",
			image:      "namespace/busybox",
			archOption: ArchAppend,
			expected:   "localhost.local/arch/busybox:latest",
		},
		{
			name:       "registry, namespace, arch prepend",
			registry:   "localhost.local",
			image:      "namespace/busybox",
			archOption: ArchPrepend,
			expected:   "localhost.local/arch/busybox:latest",
		},
		{
			name:       "registry, namespace, arch omitted",
			registry:   "localhost.local",
			image:      "namespace/busybox",
			archOption: ArchOmit,
			expected:   "localhost.local/arch/busybox:latest",
		},
		// no registry, namespace
		{
			name:       "no registry, namespace, arch append",
			registry:   "",
			image:      "namespace/busybox",
			archOption: ArchAppend,
			expected:   "docker.io/arch/busybox:latest",
		},
		{
			name:       "no registry, namespace, arch prepend",
			registry:   "",
			image:      "namespace/busybox",
			archOption: ArchPrepend,
			expected:   "docker.io/arch/busybox:latest",
		},
		{
			name:       "no registry, namespace, arch omitted",
			registry:   "",
			image:      "namespace/busybox",
			archOption: ArchOmit,
			expected:   "docker.io/arch/busybox:latest",
		},
		// registry, no namespace
		{
			name:       "registry, no namespace, arch append",
			registry:   "localhost.local",
			image:      "busybox",
			archOption: ArchAppend,
			expected:   "localhost.local/arch/busybox:latest",
		},
		{
			name:       "registry, no namespace, arch prepend",
			registry:   "localhost.local",
			image:      "busybox",
			archOption: ArchPrepend,
			expected:   "localhost.local/arch/busybox:latest",
		},
		{
			name:       "registry, no namespace, arch omit",
			registry:   "localhost.local",
			image:      "busybox",
			archOption: ArchOmit,
			expected:   "localhost.local/arch/busybox:latest",
		},
		// no registry, no namespace
		{
			name:       "no registry, no namespace, arch append",
			registry:   "",
			image:      "busybox",
			archOption: ArchAppend,
			expected:   "docker.io/arch/busybox:latest",
		},
		{
			name:       "no registry, no namespace, arch prepend",
			registry:   "",
			image:      "busybox",
			archOption: ArchPrepend,
			expected:   "docker.io/arch/busybox:latest",
		},
		{
			name:       "no registry, no namespace, arch omitted",
			registry:   "",
			image:      "busybox",
			archOption: ArchOmit,
			expected:   "docker.io/arch/busybox:latest",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := &UriOptions{
				Registry:   tc.registry,
				Official:   true,
				ArchOption: tc.archOption,
			}
			uri, err := NewUri(tc.image, opts)
			if err != nil {
				t.Errorf("An unexpected error occurred: %v", err)
			}

			actual := patchArch(uri.Remote())

			if tc.expected != actual {
				t.Errorf("expected value: %v; actual value: %v", tc.expected, actual)
			}
		})
	}
}

func TestNewUri_unofficial(t *testing.T) {
	testCases := []struct {
		name       string
		registry   string
		image      string
		archOption ArchOption
		expected   string
	}{
		// registry, namespace
		{
			name:       "registry, namespace, arch append",
			registry:   "localhost.local",
			image:      "namespace/busybox",
			archOption: ArchAppend,
			expected:   "localhost.local/namespace/busybox:latest-arch",
		},
		{
			name:       "registry, namespace, arch prepend",
			registry:   "localhost.local",
			image:      "namespace/busybox",
			archOption: ArchPrepend,
			expected:   "localhost.local/namespace/busybox:arch-latest",
		},
		{
			name:       "registry, namespace, arch omitted",
			registry:   "localhost.local",
			image:      "namespace/busybox",
			archOption: ArchOmit,
			expected:   "localhost.local/namespace/busybox:latest",
		},
		// no registry, namespace
		{
			name:       "no registry, namespace, arch append",
			registry:   "",
			image:      "namespace/busybox",
			archOption: ArchAppend,
			expected:   "docker.io/namespace/busybox:latest-arch",
		},
		{
			name:       "no registry, namespace, arch prepend",
			registry:   "",
			image:      "namespace/busybox",
			archOption: ArchPrepend,
			expected:   "docker.io/namespace/busybox:arch-latest",
		},
		{
			name:       "no registry, namespace, arch omitted",
			registry:   "",
			image:      "namespace/busybox",
			archOption: ArchOmit,
			expected:   "docker.io/namespace/busybox:latest",
		},
		// registry, no namespace
		{
			name:       "registry, no namespace, arch append",
			registry:   "localhost.local",
			image:      "busybox",
			archOption: ArchAppend,
			expected:   "localhost.local/busybox:latest-arch",
		},
		{
			name:       "registry, no namespace, arch prepend",
			registry:   "localhost.local",
			image:      "busybox",
			archOption: ArchPrepend,
			expected:   "localhost.local/busybox:arch-latest",
		},
		{
			name:       "registry, no namespace, arch omit",
			registry:   "localhost.local",
			image:      "busybox",
			archOption: ArchOmit,
			expected:   "localhost.local/busybox:latest",
		},
		// no registry, no namespace
		{
			name:       "no registry, no namespace, arch append",
			registry:   "",
			image:      "busybox",
			archOption: ArchAppend,
			expected:   "docker.io/library/busybox:latest-arch",
		},
		{
			name:       "no registry, no namespace, arch prepend",
			registry:   "",
			image:      "busybox",
			archOption: ArchPrepend,
			expected:   "docker.io/library/busybox:arch-latest",
		},
		{
			name:       "no registry, no namespace, arch omitted",
			registry:   "",
			image:      "busybox",
			archOption: ArchOmit,
			expected:   "docker.io/library/busybox:latest",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := &UriOptions{
				Registry:   tc.registry,
				Official:   false,
				ArchOption: tc.archOption,
			}
			uri, err := NewUri(tc.image, opts)
			if err != nil {
				t.Errorf("An unexpected error occurred: %v", err)
			}

			actual := patchArch(uri.Remote())

			if tc.expected != actual {
				t.Errorf("expected value: %v; actual value: %v", tc.expected, actual)
			}
		})
	}
}

func TestNewUriWithDigest_unofficial(t *testing.T) {
	testCases := []testCase{
		{
			name:       "no registry, no namespace, arch ignored",
			registry:   "",
			image:      "image@sha256:4abcc75d00d254c931cb5ce4a5c2ebb6aab90f19f70bf79d6734a0e3f8f2c72f",
			archOption: ArchPrepend,
			expected:   "docker.io/library/image@sha256:4abcc75d00d254c931cb5ce4a5c2ebb6aab90f19f70bf79d6734a0e3f8f2c72f",
		},
		{
			name:       "no registry, namespace, arch ignored",
			registry:   "",
			image:      "namespace/image@sha256:4abcc75d00d254c931cb5ce4a5c2ebb6aab90f19f70bf79d6734a0e3f8f2c72f",
			archOption: ArchPrepend,
			expected:   "docker.io/namespace/image@sha256:4abcc75d00d254c931cb5ce4a5c2ebb6aab90f19f70bf79d6734a0e3f8f2c72f",
		},
		{
			name:       "registry, namespace, arch ignored",
			registry:   "localhost.local",
			image:      "namespace/image@sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			archOption: ArchAppend,
			expected:   "localhost.local/namespace/image@sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		},
		{
			name:       "registry port, namespace, arch ignored",
			registry:   "localhost:5000",
			image:      "namespace/image@sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			archOption: ArchAppend,
			expected:   "localhost:5000/namespace/image@sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runTestCaseUnofficial(t, tc)
		})
	}
}

func TestNewUriWithDigest_official(t *testing.T) {
	testCases := []testCase{
		{
			name:       "no registry, no namespace, arch ignored",
			registry:   "",
			image:      "image@sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			archOption: ArchPrepend,
			expected:   "docker.io/arch/image@sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		},
		{
			name:       "no registry, namespace, arch ignored",
			registry:   "",
			image:      "namespace/image@sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			archOption: ArchPrepend,
			expected:   "docker.io/arch/image@sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		},
		{
			name:       "registry, namespace, arch ignored",
			registry:   "localhost.local",
			image:      "namespace/image@sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			archOption: ArchAppend,
			expected:   "localhost.local/arch/image@sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		},
		{
			name:       "registry port, namespace, arch ignored",
			registry:   "localhost:5000",
			image:      "namespace/image@sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			archOption: ArchAppend,
			expected:   "localhost:5000/arch/image@sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runTestCaseOfficial(t, tc)
		})
	}
}

func TestNewUriWithDefinedArch_official(t *testing.T) {
	testCases := []testCaseWithArch{
		// registry, namespace
		{
			name:       "arch amd64",
			registry:   "localhost.local",
			image:      "namespace/busybox",
			arch:       "amd64",
			archOption: ArchAppend,
			expected:   "localhost.local/amd64/busybox:latest",
		},
		{
			name:       "arch arm64",
			registry:   "localhost.local",
			image:      "namespace/busybox",
			arch:       "arm64",
			archOption: ArchAppend,
			expected:   "localhost.local/arm64/busybox:latest",
		},
		{
			name:       "arch ppc64le",
			registry:   "localhost.local",
			image:      "namespace/busybox",
			arch:       "ppc64le",
			archOption: ArchAppend,
			expected:   "localhost.local/ppc64le/busybox:latest",
		},
		{
			name:       "arch s390x",
			registry:   "localhost.local",
			image:      "namespace/busybox",
			arch:       "s390x",
			archOption: ArchAppend,
			expected:   "localhost.local/s390x/busybox:latest",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runTestCaseWithArchOfficial(t, tc)
		})
	}
}

func TestNewUriWithDefinedArch_unofficial(t *testing.T) {
	testCases := []testCaseWithArch{
		// registry, namespace
		{
			name:       "arch amd64",
			registry:   "localhost.local",
			image:      "namespace/busybox",
			arch:       "amd64",
			archOption: ArchAppend,
			expected:   "localhost.local/namespace/busybox:latest-amd64",
		},
		{
			name:       "arch arm64",
			registry:   "localhost.local",
			image:      "namespace/busybox",
			arch:       "arm64",
			archOption: ArchAppend,
			expected:   "localhost.local/namespace/busybox:latest-arm64",
		},
		{
			name:       "arch ppc64le",
			registry:   "localhost.local",
			image:      "namespace/busybox",
			arch:       "ppc64le",
			archOption: ArchAppend,
			expected:   "localhost.local/namespace/busybox:latest-ppc64le",
		},
		{
			name:       "arch s390x",
			registry:   "localhost.local",
			image:      "namespace/busybox",
			arch:       "s390x",
			archOption: ArchAppend,
			expected:   "localhost.local/namespace/busybox:latest-s390x",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runTestCaseWithArchUnofficial(t, tc)
		})
	}
}

func TestNewUriWithFullUri(t *testing.T) {
	testCases := []testCase{
		{
			name:       "Image with full uri string",
			registry:   "localhost.local",
			image:      "localhost.local/namespace/busybox:latest",
			archOption: ArchPrepend,
			expected:   "localhost.local/namespace/busybox:arch-latest",
		},
		{
			name:       "Image with full uri string",
			registry:   "localhost2.local",
			image:      "localhost.local/namespace/busybox:latest",
			archOption: ArchPrepend,
			expected:   "localhost.local/namespace/busybox:arch-latest",
		},
		{
			name:       "Image with full uri string",
			registry:   "localhost.local",
			image:      "localhost2.local/namespace/busybox:latest",
			archOption: ArchPrepend,
			expected:   "localhost2.local/namespace/busybox:arch-latest",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := &UriOptions{
				Registry:   tc.registry,
				Official:   false,
				ArchOption: tc.archOption,
			}
			uri, err := NewUri(tc.image, opts)
			if err != nil {
				t.Errorf("An unexpected error occurred: %v", err)
			}

			actual := patchArch(uri.Remote())

			if tc.expected != actual {
				t.Errorf("expected value: %v; actual value: %v", tc.expected, actual)
			}
		})
	}
}

func TestNewUriWithArchOptionAsString(t *testing.T) {
	expected := "docker.io/username/image:tag-arch"
	opts := &UriOptions{
		Official:   false,
		ArchOption: "append",
	}
	uri, _ := NewUri("username/image:tag", opts)
	actual := patchArch(uri.Remote())
	if expected != actual {
		t.Errorf("expected value: '%v'; actual value: '%v'", expected, actual)
	}

	expected = "docker.io/username/image:tag"
	opts = &UriOptions{
		Official:   false,
		ArchOption: "invalid",
	}
	uri, _ = NewUri("username/image:tag", opts)
	actual = patchArch(uri.Remote())
	if expected != actual {
		t.Errorf("expected value: '%v'; actual value: '%v'", expected, actual)
	}
}

type testCase struct {
	name       string
	registry   string
	image      string
	archOption ArchOption
	expected   string
}

func runTestCaseUnofficial(t *testing.T, tc testCase) {
	t.Helper()

	opts := &UriOptions{
		Registry:   tc.registry,
		Official:   false,
		ArchOption: tc.archOption,
	}
	uri, _ := NewUri(tc.image, opts)

	actual := patchArch(uri.Remote())

	if tc.expected != actual {
		t.Errorf("expected value: '%v'; actual value: '%v'", tc.expected, actual)
	}
}

func runTestCaseOfficial(t *testing.T, tc testCase) {
	t.Helper()

	opts := &UriOptions{
		Registry:   tc.registry,
		Official:   true,
		ArchOption: tc.archOption,
	}
	uri, _ := NewUri(tc.image, opts)

	actual := patchArch(uri.Remote())

	if tc.expected != actual {
		t.Errorf("expected value: '%v'; actual value: '%v'", tc.expected, actual)
	}
}

type testCaseWithArch struct {
	name       string
	registry   string
	image      string
	archOption ArchOption
	arch       string
	expected   string
}

func runTestCaseWithArchUnofficial(t *testing.T, tc testCaseWithArch) {
	t.Helper()

	opts := &UriOptions{
		Registry:   tc.registry,
		Official:   false,
		Arch:       tc.arch,
		ArchOption: tc.archOption,
	}
	uri, _ := NewUri(tc.image, opts)

	actual := uri.Remote()

	if tc.expected != actual {
		t.Errorf("expected value: '%v'; actual value: '%v'", tc.expected, actual)
	}
}

func runTestCaseWithArchOfficial(t *testing.T, tc testCaseWithArch) {
	t.Helper()

	opts := &UriOptions{
		Registry:   tc.registry,
		Official:   true,
		Arch:       tc.arch,
		ArchOption: tc.archOption,
	}
	uri, _ := NewUri(tc.image, opts)

	actual := uri.Remote()

	if tc.expected != actual {
		t.Errorf("expected value: '%v'; actual value: '%v'", tc.expected, actual)
	}
}
