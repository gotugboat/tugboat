package image

import (
	"context"
	"tugboat/internal/driver"
	"tugboat/internal/term"

	log "github.com/sirupsen/logrus"
)

type TagOptions struct {
	// Should be in 'image[:tag]' format
	SourceImage            string
	Tags                   []string
	Push                   bool
	SupportedArchitectures []string
}

func Tag(ctx context.Context, d driver.ImageBuilderPusher, opts TagOptions) error {
	if len(opts.Tags) == 0 {
		return ErrNoProvidedTags
	}

	if len(opts.SupportedArchitectures) == 0 {
		return ErrNoSupportedArchitectures
	}

	for _, arch := range opts.SupportedArchitectures {
		pullOutput, err := d.PullImageWithArch(ctx, opts.SourceImage, arch)
		if err != nil {
			return err
		}

		if pullOutput != nil {
			defer pullOutput.Close()

			if err := term.DisplayResponse(pullOutput); err != nil {
				return err
			}
		}

		for _, targetTag := range opts.Tags {
			taggedUri, err := d.TagImage(ctx, opts.SourceImage, targetTag)
			if err != nil {
				return err
			}

			if opts.Push {
				pushOutput, err := d.PushImageWithArch(ctx, taggedUri, arch)
				if err != nil {
					log.Errorf("tag push: %v", err)
				}

				if pushOutput != nil {
					defer pushOutput.Close()

					if err := term.DisplayResponse(pushOutput); err != nil {
						return err
					}
				}

			}
		}
	}

	return nil
}
