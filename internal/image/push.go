package image

import (
	"context"
	"tugboat/internal/driver"
	"tugboat/internal/term"
)

func Push(ctx context.Context, d driver.ImagePusher, image string) error {
	output, err := d.PushImage(ctx, image)
	if err != nil {
		return err
	}

	if output != nil {
		defer output.Close()

		if err := term.DisplayResponse(output); err != nil {
			return err
		}
	}

	return nil
}
