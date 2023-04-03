package tmpl

import (
	"testing"
	"tugboat/internal/pkg/flags"
)

var opts = &flags.Options{
	Global: flags.GlobalOptions{
		Git: flags.Git{
			Branch:      "branch",
			Tag:         "tag",
			Commit:      "shortCommit",
			ShortCommit: "shortCommit",
			FullCommit:  "fullCommit",
		},
	},
	Image: flags.ImageOptions{
		Name:    "image",
		Version: "version",
	},
}

func TestCompileString(t *testing.T) {
	testCases := []struct {
		name     string
		template string
		expected string
	}{
		{
			name:     "docker image w/ version tag",
			template: `{{.ImageName}}:{{.Version}}`,
			expected: "image:version",
		},
		{
			name:     "docker image w/ all git commit info",
			template: `{{.ImageName}}:{{.Commit}}-{{.ShortCommit}}-{{.FullCommit}}`,
			expected: "image:shortCommit-shortCommit-fullCommit",
		},
		{
			name:     "docker image w/ all git branch and tag",
			template: `{{.ImageName}}:{{.Tag}}-{{.Branch}}`,
			expected: "image:tag-branch",
		},
		{
			name:     "replace",
			template: `{{ replace "image" "image" "another" }}`,
			expected: "another",
		},
		{
			name:     "tolower",
			template: `{{ tolower "LOWER" }}`,
			expected: "lower",
		},
		{
			name:     "toupper",
			template: `{{ toupper "upper" }}`,
			expected: "UPPER",
		},
		{
			name:     "trim",
			template: `{{ printf "  trimmed  " | trim }}`,
			expected: "trimmed",
		},
		{
			name:     "trimprefix",
			template: `{{ trimprefix "trimmed" "tr" }}`,
			expected: "immed",
		},
		{
			name:     "trimsuffix",
			template: `{{ trimsuffix "trimmed" "med" }}`,
			expected: "trim",
		},
		{
			name:     "title",
			template: `{{ title "title case" }}`,
			expected: "Title Case",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			compiledString, _ := CompileString(tc.template, opts)

			if tc.expected != compiledString {
				t.Errorf("expected '%v', got '%v'", tc.expected, compiledString)
			}
		})
	}
}

func TestCompileString_time(t *testing.T) {
	testCases := []struct {
		name     string
		template string
	}{
		{
			name:     "time YYYY-MM-DD",
			template: `{{ time "2006-01-02" }}`,
		},
		{
			name:     "time MM/DD/YYYY",
			template: `{{ time "01/02/2006" }}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			compiledString, _ := CompileString(tc.template, opts)

			if len(compiledString) == 0 {
				t.Errorf("expected a time, got '%v'", compiledString)
			}
		})
	}
}

func TestCompileStringSlice(t *testing.T) {
	testCases := []struct {
		name     string
		template []string
		expected []string
	}{
		{
			name:     "slice",
			template: []string{"{{.ImageName}}", "{{.Version}}", "{{.Branch}}"},
			expected: []string{"image", "version", "branch"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			compiledString, _ := CompileStringSlice(tc.template, opts)

			for i, item := range tc.expected {
				if item != compiledString[i] {
					t.Errorf("expected '%v', got '%v'", tc.expected, compiledString)
				}
			}
		})
	}
}
