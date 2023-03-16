package manifest

import (
	"testing"
	"tugboat/internal/pkg/flags"

	"github.com/spf13/pflag"
)

func TestNewManifestCommand(t *testing.T) {
	globalFlags := flags.NewGlobalFlagGroup()
	cmd := NewManifestCommand(globalFlags)

	// validate the description strings
	expected := "Manage image manifests"
	if expected != cmd.Long {
		t.Errorf("expected %v, got %v", expected, cmd.Long)
	}

	expected = "Manage image manifests"
	if expected != cmd.Short {
		t.Errorf("expected %v, got %v", expected, cmd.Long)
	}

	// validate the number of commands attached to this command
	commands := cmd.Commands()
	expectedCommands := 1
	actualCommands := len(commands)
	if actualCommands != expectedCommands {
		t.Errorf("expected commands %v, got %v", expectedCommands, actualCommands)
	}

	// validate what flags are attached to this command
	if ok := cmd.HasLocalFlags(); ok {
		t.Error("expected no flags, but there are flags")
	}

	// validate the number of flags
	expectedFlagCount := 0
	actualFlagCount := 0
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		actualFlagCount++
	})

	if actualFlagCount != expectedFlagCount {
		t.Errorf("expected %v flags, got %v", expectedFlagCount, actualFlagCount)
	}
}
