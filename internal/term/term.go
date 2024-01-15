package term

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"tugboat/internal/pkg/slices"
	"tugboat/internal/types"

	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/moby/term"
)

type Command struct {
	Command string
	Args    []string
}

// Display client responses to the terminal
func DisplayResponse(r io.Reader) error {
	if r == nil {
		return nil
	}
	termFd, isTerm := term.GetFdInfo(os.Stderr)
	if err := jsonmessage.DisplayJSONMessagesStream(r, os.Stderr, termFd, isTerm, nil); err != nil {
		return err
	}
	return nil
}

func DisplayOutput(r io.Reader) error {
	if r == nil {
		return nil
	}
	if err := displayStream(r, os.Stdout); err != nil {
		return err
	}
	return nil
}

func StreamCommand(cmd *Command, isDryRun bool, isDebug bool) (io.ReadCloser, error) {
	if err := validateCommand(cmd); err != nil {
		return nil, err
	}

	command := exec.Command(cmd.Command, cmd.Args...)

	stdout, err := command.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("error creating stdout pipe: %v", err)
	}

	stderr, err := command.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("error creating stderr pipe: %v", err)
	}

	// Create a multiReadCloser to combine stdout and stderr
	mrc := &types.MultiReadCloser{
		Readers: []io.ReadCloser{stdout, stderr},
	}

	// Start the command
	if err := command.Start(); err != nil {
		mrc.Close() // Close the pipes in case of an error
		return nil, fmt.Errorf("error starting command: %v", err)
	}

	return mrc, nil
}

func displayStream(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Fprintf(w, "%s\n", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// validateCommand ensures that the command only contains expected commands so unexpected programs are not run
func validateCommand(cmd *Command) error {
	allowedCommands := []string{"docker"}
	disallowedCharacters := []string{";", "|", "&"}

	if !slices.Contains(allowedCommands, cmd.Command) {
		return fmt.Errorf("disallowed command: %s", cmd.Command)
	}

	input := strings.Join(cmd.Args, " ")
	for _, char := range disallowedCharacters {
		if strings.Contains(input, char) {
			return fmt.Errorf("disallowed character: %s", char)
		}
	}

	return nil
}
