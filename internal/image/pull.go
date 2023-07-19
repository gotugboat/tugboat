package image

import (
	"context"
	"tugboat/internal/driver"
	"tugboat/internal/term"
)

func Pull(ctx context.Context, d driver.ImageBuilder, image string) error {
	output, err := d.PullImage(ctx, image)
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
