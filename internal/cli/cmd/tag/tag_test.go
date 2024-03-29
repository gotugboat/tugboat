package tag

import (
	"testing"
	"tugboat/internal/pkg/flags"

	"github.com/spf13/pflag"
)

func TestTagCommand(t *testing.T) {
	globalFlags := flags.NewGlobalFlagGroup()
	cmd := NewTagCommand(globalFlags)

	// validate the description strings
	expected := "Create a tag that refers to another image"
	if expected != cmd.Long {
		t.Errorf("expected %v, got %v", expected, cmd.Long)
	}

	expected = "Create a tag that refers to another image"
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
	expectedFlagCount := 3
	actualFlagCount := 0
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		actualFlagCount++
	})

	if actualFlagCount != expectedFlagCount {
		t.Errorf("expected %v flags, got %v", expectedFlagCount, actualFlagCount)
	}

	// validate each flag
	if _, err := cmd.Flags().GetStringSlice("architectures"); err != nil {
		t.Error(err)
	}

	if _, err := cmd.Flags().GetStringSlice("tags"); err != nil {
		t.Error(err)
	}

	if _, err := cmd.Flags().GetBool("push"); err != nil {
		t.Error(err)
	}
}
