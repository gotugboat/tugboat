package commands

import (
	"testing"
)

func TestNewCli(t *testing.T) {
	cli := NewCli()

	// validate the number of commands attached to the cli
	commands := cli.Commands()
	expectedNumCommands := 2 // the default completion and help commands are not counted
	actualNumCommands := len(commands)
	if actualNumCommands != expectedNumCommands {
		t.Errorf("expected commands %v, got %v", expectedNumCommands, actualNumCommands)
	}

	// validate what commands are attached to the cli (the default completion and help commands are not counted)
	expectedCommands := []string{
		"version",
		"build",
	}
	for _, command := range commands {
		if !contains(expectedCommands, command.Name()) {
			t.Errorf("%v is not an expected command", command.Name())
		}
	}
}

func contains(slice []string, element string) bool {
	for _, value := range slice {
		if value == element {
			return true
		}
	}
	return false
}
