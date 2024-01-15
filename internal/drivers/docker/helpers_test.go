package docker

import (
	"fmt"
	"strings"
	"testing"
	"tugboat/internal/driver"
	"tugboat/internal/pkg/reference"
	"tugboat/internal/registry"
	"tugboat/internal/term"
)

var (
	image       = "image"
	manifestTag = "tag"
)

var basicCreateOpts = driver.ManifestCreateOptions{
	ManifestList:           image,
	ManifestTags:           []string{manifestTag},
	SupportedArchitectures: []string{"arm64"},
}

var basicPushOpts = driver.ManifestPushOptions{
	Purge: true,
}

func Test_getCommand(t *testing.T) {
	cmd := getCommand([]string{})

	if cmd.Command != "docker" {
		t.Errorf("expected command: docker, got %v", cmd.Command)
	}
}

func Test_validateCommand(t *testing.T) {
	testCases := []struct {
		name        string
		command     string
		args        []string
		expectedErr bool
	}{
		{
			name:        "valid command",
			command:     "docker",
			args:        []string{"images"},
			expectedErr: false,
		},
		{
			name:        "invalid char ;",
			command:     "docker",
			args:        []string{"images", ";", "rm", "-rf", "./non-existent-dir"},
			expectedErr: true,
		},
		{
			name:        "invalid char |",
			command:     "docker",
			args:        []string{"images", "|", "something"},
			expectedErr: true,
		},
		{
			name:        "invalid char &",
			command:     "docker",
			args:        []string{"images", "&&", "something"},
			expectedErr: true,
		},
		{
			name:        "invalid command",
			command:     "rm",
			args:        []string{"-rf", "./non-existent-dir"},
			expectedErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := &term.Command{
				Command: tc.command,
				Args:    tc.args,
			}
			err := validateCommand(cmd)

			if tc.expectedErr && err == nil {
				t.Errorf("expected command to be valid, but was invalid")
			}

			if !tc.expectedErr && err != nil {
				t.Errorf("expected command to be invalid, but was valid")
			}

		})
	}
}

func Test_getCreateArgs(t *testing.T) {
	registry, err := registry.NewRegistry("docker.io", "namespace", "username", "password")
	if err != nil {
		t.Errorf("create registry failed: %v", err)
	}

	imageName := fmt.Sprintf("%s:%s", image, manifestTag)
	ref, _ := reference.NewUri(fmt.Sprintf("%s/%s", registry.Namespace, imageName), &reference.UriOptions{
		Registry: registry.ServerAddress,
		Official: false,
	})

	expectedArgs := "manifest create docker.io/namespace/image:tag docker.io/namespace/image:arm64-tag"
	args, _ := getCreateArgs(ref, basicCreateOpts.SupportedArchitectures, false, "prepend")
	actualArgs := strings.Join(args, " ")

	if actualArgs != expectedArgs {
		t.Errorf("expected commands %v, got %v", expectedArgs, actualArgs)
	}
}

func Test_getPushArgs(t *testing.T) {
	registry, err := registry.NewRegistry("docker.io", "namespace", "username", "password")
	if err != nil {
		t.Errorf("create registry failed: %v", err)
	}

	imageName := fmt.Sprintf("%s:%s", image, manifestTag)
	ref, _ := reference.NewUri(fmt.Sprintf("%s/%s", registry.Namespace, imageName), &reference.UriOptions{
		Registry: registry.ServerAddress,
		Official: false,
	})

	expectedArgs := "manifest push --purge docker.io/namespace/image:tag"
	args, _ := getPushArgs(ref, basicPushOpts)
	actualArgs := strings.Join(args, " ")

	if actualArgs != expectedArgs {
		t.Errorf("expected commands %v, got %v", expectedArgs, actualArgs)
	}
}

func Test_getAnnotateCommands(t *testing.T) {
	registry, err := registry.NewRegistry("docker.io", "namespace", "username", "password")
	if err != nil {
		t.Errorf("create registry failed: %v", err)
	}

	imageName := fmt.Sprintf("%s:%s", image, manifestTag)
	ref, _ := reference.NewUri(fmt.Sprintf("%s/%s", registry.Namespace, imageName), &reference.UriOptions{
		Registry: registry.ServerAddress,
		Official: false,
	})

	expectedArgs := "manifest annotate docker.io/namespace/image:tag docker.io/namespace/image:arm64-tag --arch arm64"
	annotateCmds, _ := getAnnotateCommands(ref, basicCreateOpts.SupportedArchitectures, false, "prepend")

	actualArgs := strings.Join(annotateCmds[0].Args, " ")
	if actualArgs != expectedArgs {
		t.Errorf("expected commands %v, got %v", expectedArgs, actualArgs)
	}

}

func Test_getRmArgs(t *testing.T) {
	registry, err := registry.NewRegistry("docker.io", "namespace", "username", "password")
	if err != nil {
		t.Errorf("create registry failed: %v", err)
	}
	imageName := fmt.Sprintf("%s:%s", image, manifestTag)
	ref, _ := reference.NewUri(fmt.Sprintf("%s/%s", registry.Namespace, imageName), &reference.UriOptions{
		Registry: registry.ServerAddress,
		Official: false,
	})

	expectedArgs := "manifest rm docker.io/namespace/image:tag"
	rmArgs, _ := getRmArgs([]*reference.Reference{ref})

	actualArgs := strings.Join(rmArgs, " ")
	if expectedArgs != actualArgs {
		t.Errorf("expected args %v, got %v", expectedArgs, actualArgs)
	}
}
