package manifest

import (
	"context"
	"fmt"
	"tugboat/internal/driver"
	"tugboat/internal/term"

	"github.com/pkg/errors"
)

var (
	ErrNoProvidedTags           = errors.New("tags must be provided")
	ErrNoSupportedArchitectures = errors.New("there are no supported architectures define")
)

type ManifestCreateOptions struct {
	ManifestList           string
	ManifestTags           []string
	Push                   bool
	SupportedArchitectures []string
}

func Create(ctx context.Context, d driver.Driver, opts ManifestCreateOptions) error {
	if len(opts.ManifestTags) == 0 {
		return errors.Wrap(ErrNoProvidedTags, "Create manifest failed")
	}

	if len(opts.SupportedArchitectures) == 0 {
		return errors.Wrap(ErrNoSupportedArchitectures, "Create manifest failed")
	}

	// Pre-pull all the images for each architecture
	for _, arch := range opts.SupportedArchitectures {
		for _, manifestTag := range opts.ManifestTags {
			// Generate the uri for the manifest list
			manifestListUri := fmt.Sprintf("%s:%s", opts.ManifestList, manifestTag)
			output, err := d.PullImageWithArch(ctx, manifestListUri, arch)
			if err != nil {
				return err
			}

			if output != nil {
				defer output.Close()

				if err := term.DisplayResponse(output); err != nil {
					return err
				}
			}
		}
	}

	// Create the manifests
	output, err := d.CreateManifest(ctx, driver.ManifestCreateOptions{
		ManifestList:           opts.ManifestList,
		ManifestTags:           opts.ManifestTags,
		SupportedArchitectures: opts.SupportedArchitectures,
	})
	if err != nil {
		return err
	}

	if output != nil {
		defer output.Close()

		if err := term.DisplayResponse(output); err != nil {
			return err
		}
	}

	// push all the manifests to the registry, removing them from the local disk
	if opts.Push {
		for _, manifestTag := range opts.ManifestTags {
			manifestName := fmt.Sprintf("%s:%s", opts.ManifestList, manifestTag)
			if err := d.PushManifest(ctx, manifestName, driver.ManifestPushOptions{Purge: true}); err != nil {
				return err
			}
		}
	}

	return nil
}
