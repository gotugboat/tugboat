package term

import (
	"io"
	"os"

	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/moby/term"
)

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
