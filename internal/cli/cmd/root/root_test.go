package root

import (
	"testing"
	"tugboat/internal/pkg/flags"
)

func TestRootCommand(t *testing.T) {
	globalFlags := flags.NewGlobalFlagGroup()
	cmd := NewRootCommand(globalFlags)

	// validate the description strings
	expected := "A tool to build and publish multi-architecture container images"
	if expected != cmd.Long {
		t.Errorf("expected %v, got %v", expected, cmd.Long)
	}

	expected = "Build multi-arch images"
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

	if _, err := cmd.Flags().GetBool("dry-run"); err != nil {
		t.Error(err)
	}

	if _, err := cmd.Flags().GetBool("debug"); err != nil {
		t.Error(err)
	}

	// validate command settings
	if cmd.SilenceUsage != true {
		t.Error("SilenceUsage should be false")
	}

	if cmd.SilenceErrors != true {
		t.Error("SilenceErrors should be false")
	}
}
