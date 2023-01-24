package version

import (
	"testing"
)

func TestVersionCommand(t *testing.T) {
	cmd := NewVersionCommand()

	// validate the description string
	expected := "Show the Tugboat version information"
	if expected != cmd.Long {
		t.Errorf("expected %v, got %v", expected, cmd.Long)
	}

	expected = "Show version information"
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

	_, err := cmd.Flags().GetBool("short")
	if err != nil {
		t.Error(err)
	}
}
