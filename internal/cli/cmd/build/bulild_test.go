package build

import (
	"testing"
	"tugboat/internal/pkg/flags"

	"github.com/spf13/pflag"
)

func TestBuildCommand(t *testing.T) {
	globalFlags := flags.NewGlobalFlagGroup()
	cmd := NewBuildCommand(globalFlags)

	// validate the description strings
	expected := "Build an image from a Dockerfile"
	if expected != cmd.Long {
		t.Errorf("expected %v, got %v", expected, cmd.Long)
	}

	expected = "Build a container"
	if expected != cmd.Short {
		t.Errorf("expected %v, got %v", expected, cmd.Long)
	}

	// validate the number of commands attached to this command
	commands := cmd.Commands()
	expectedCommands := 0
	actualCommands := len(commands)
	if actualCommands != expectedCommands {
		t.Errorf("expected commands %v, got %v", expectedCommands, actualCommands)
	}

	// validate what flags are attached to this command
	if ok := cmd.HasLocalFlags(); !ok {
		t.Error("expected to see flags, but there are none")
	}

	// validate the number of flags
	expectedFlagCount := 7
	actualFlagCount := 0
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		actualFlagCount++
	})

	if actualFlagCount != expectedFlagCount {
		t.Errorf("expected %v flags, got %v", expectedFlagCount, actualFlagCount)
	}

	// validate each flag
	if _, err := cmd.Flags().GetStringSlice("build-args"); err != nil {
		t.Error(err)
	}

	if _, err := cmd.Flags().GetString("context"); err != nil {
		t.Error(err)
	}

	if _, err := cmd.Flags().GetString("file"); err != nil {
		t.Error(err)
	}

	if _, err := cmd.Flags().GetStringSlice("tag"); err != nil {
		t.Error(err)
	}

	if _, err := cmd.Flags().GetBool("push"); err != nil {
		t.Error(err)
	}

	if _, err := cmd.Flags().GetBool("pull"); err != nil {
		t.Error(err)
	}

	if _, err := cmd.Flags().GetBool("no-cache"); err != nil {
		t.Error(err)
	}
}
