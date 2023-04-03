package image

import (
	"fmt"
	"strings"
	"testing"
	"tugboat/internal/pkg/docker"
	"tugboat/internal/pkg/flags"
)

var (
	image       = "image"
	manifestTag = "tag"
)

var basicCreateOpts = ManifestCreateOptions{
	ManifestList:           image,
	ManifestTags:           []string{manifestTag},
	Push:                   false,
	SupportedArchitectures: []string{"arm64"},
	Registry: NewRegistry(
		"docker.io",
		"namespace",
		"username",
		"password",
	),
	Official:   false,
	DryRun:     true,
	Debug:      false,
	ArchOption: flags.DefaultArchOption,
}

var basicPushOpts = PushManifestOptions{
	Purge:  true,
	DryRun: false,
	Debug:  false,
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
			cmd := &Command{
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
	imageName := fmt.Sprintf("%s:%s", image, manifestTag)
	ref, _ := docker.NewUri(fmt.Sprintf("%s/%s", basicCreateOpts.Registry.Namespace, imageName), &docker.UriOptions{
		Registry: basicCreateOpts.Registry.ServerAddress,
		Official: basicCreateOpts.Official,
	})

	expectedArgs := "manifest create docker.io/namespace/image:tag docker.io/namespace/image:arm64-tag"
	args, _ := getCreateArgs(ref, basicCreateOpts)
	actualArgs := strings.Join(args, " ")

	if actualArgs != expectedArgs {
		t.Errorf("expected commands %v, got %v", expectedArgs, actualArgs)
	}
}

func Test_getPushArgs(t *testing.T) {
	imageName := fmt.Sprintf("%s:%s", image, manifestTag)
	ref, _ := docker.NewUri(fmt.Sprintf("%s/%s", basicCreateOpts.Registry.Namespace, imageName), &docker.UriOptions{
		Registry: basicCreateOpts.Registry.ServerAddress,
		Official: basicCreateOpts.Official,
	})

	expectedArgs := "manifest push --purge docker.io/namespace/image:tag"
	args, _ := getPushArgs(ref, basicPushOpts)
	actualArgs := strings.Join(args, " ")

	if actualArgs != expectedArgs {
		t.Errorf("expected commands %v, got %v", expectedArgs, actualArgs)
	}
}

func Test_getAnnotateCommands(t *testing.T) {

	imageName := fmt.Sprintf("%s:%s", image, manifestTag)
	ref, _ := docker.NewUri(fmt.Sprintf("%s/%s", basicCreateOpts.Registry.Namespace, imageName), &docker.UriOptions{
		Registry: basicCreateOpts.Registry.ServerAddress,
		Official: basicCreateOpts.Official,
	})

	expectedArgs := "manifest annotate docker.io/namespace/image:tag docker.io/namespace/image:arm64-tag --arch arm64"
	annotateCmds, _ := getAnnotateCommands(ref, basicCreateOpts)
	// actualArgs := strings.Join(args, " ")

	actualArgs := strings.Join(annotateCmds[0].Args, " ")
	if actualArgs != expectedArgs {
		t.Errorf("expected commands %v, got %v", expectedArgs, actualArgs)
	}

}

func Test_getRmArgs(t *testing.T) {
	imageName := fmt.Sprintf("%s:%s", image, manifestTag)
	ref, _ := docker.NewUri(fmt.Sprintf("%s/%s", basicCreateOpts.Registry.Namespace, imageName), &docker.UriOptions{
		Registry: basicCreateOpts.Registry.ServerAddress,
		Official: basicCreateOpts.Official,
	})

	expectedArgs := "manifest rm docker.io/namespace/image:tag"
	rmArgs, _ := getRmArgs([]*docker.Reference{ref})

	actualArgs := strings.Join(rmArgs, " ")
	if expectedArgs != actualArgs {
		t.Errorf("expected args %v, got %v", expectedArgs, actualArgs)
	}
}
