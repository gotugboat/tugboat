package image

import (
	"context"
	"fmt"
	"tugboat/internal/pkg/docker"

	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

type TagOptions struct {
	// Should be in 'image[:tag]' format
	SourceImage            string
	Tags                   []string
	Push                   bool
	SupportedArchitectures []string

	Registry Registry
	Official bool
	DryRun   bool
	Debug    bool

	ArchOption string
}

func ImageTag(ctx context.Context, client *client.Client, opts TagOptions) error {
	if len(opts.Tags) == 0 {
		return ErrNoProvidedTags
	}

	if len(opts.SupportedArchitectures) == 0 {
		return ErrNoSupportedArchitectures
	}

	for _, arch := range opts.SupportedArchitectures {
		// Generate the uri for the source image
		sourceUri, err := docker.NewUri(fmt.Sprintf("%s/%s", opts.Registry.Namespace, opts.SourceImage), &docker.UriOptions{
			Registry:   opts.Registry.ServerAddress,
			Official:   opts.Official,
			Arch:       arch,
			ArchOption: toArchOption(opts.ArchOption),
		})
		if err != nil {
			return err
		}

		// Pull the image for each architecture to tag
		if err := pull(ctx, client, opts.Registry, sourceUri.Remote(), opts.DryRun); err != nil {
			return err
		}

		for _, targetTag := range opts.Tags {
			// Generate the uri for the target tag
			targetUri, err := docker.NewUri(fmt.Sprintf("%v:%v", sourceUri.ShortName(), targetTag), &docker.UriOptions{
				Registry:   opts.Registry.ServerAddress,
				Official:   opts.Official,
				Arch:       arch,
				ArchOption: toArchOption(opts.ArchOption),
			})
			if err != nil {
				return err
			}

			// Tag the image for each additional reference tag
			if err := tag(ctx, client, sourceUri.Remote(), targetUri.Remote(), opts.DryRun); err != nil {
				return err
			}

			if opts.Push {
				// Push the tagged image
				if err := push(ctx, client, opts.Registry, targetUri.Remote(), opts.DryRun); err != nil {
					log.Error(err)
				}
			}
		}
	}

	return nil
}

func tag(ctx context.Context, client *client.Client, sourceImage string, targetImage string, isDryRun bool) error {
	log.Infof("Tagging %v as %v", sourceImage, targetImage)

	if isDryRun {
		return nil
	}

	if err := client.ImageTag(ctx, sourceImage, targetImage); err != nil {
		return err
	}
	return nil
}
