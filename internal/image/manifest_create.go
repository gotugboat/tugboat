package image

import (
	"context"
	"fmt"
	"tugboat/internal/pkg/reference"

	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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

	Registry Registry
	Official bool
	DryRun   bool
	Debug    bool

	ArchOption string
}

func ManifestCreate(ctx context.Context, client *client.Client, opts ManifestCreateOptions) error {
	if len(opts.ManifestTags) == 0 {
		return errors.Wrap(ErrNoProvidedTags, "Create manifest failed")
	}

	if len(opts.SupportedArchitectures) == 0 {
		return errors.Wrap(ErrNoSupportedArchitectures, "Create manifest failed")
	}

	// Generate the uri for the manifest list
	manifestListUri, err := reference.NewUri(fmt.Sprintf("%s/%s", opts.Registry.Namespace, opts.ManifestList), &reference.UriOptions{
		Registry: opts.Registry.ServerAddress,
		Official: opts.Official,
	})
	if err != nil {
		return err
	}

	// Pre-pull all the image for each architecture
	for _, arch := range opts.SupportedArchitectures {
		for _, manifestTag := range opts.ManifestTags {
			// Generate the tagged uri to pull
			imageName := fmt.Sprintf("%s:%s", manifestListUri.ShortName(), manifestTag)
			manifestTagUri, err := reference.NewUri(fmt.Sprintf("%s/%s", opts.Registry.Namespace, imageName), &reference.UriOptions{
				Registry:   opts.Registry.ServerAddress,
				Official:   opts.Official,
				Arch:       arch,
				ArchOption: toArchOption(opts.ArchOption),
			})
			if err != nil {
				return err
			}

			// Pull the image
			if err := pull(ctx, client, opts.Registry, manifestTagUri.Remote(), opts.DryRun); err != nil {
				return err
			}
		}
	}

	// login to the registry (required to create a manifest)
	loginOpts := &DockerLoginOptions{
		ServerAddress: opts.Registry.ServerAddress,
		Username:      opts.Registry.User.Name,
		Password:      opts.Registry.User.Password,
		DryRun:        opts.DryRun,
	}
	if err := dockerLogin(ctx, loginOpts); err != nil {
		return err
	}

	// Generate the manifests for each desired tag
	manifestsToPush := []*reference.Reference{}
	for _, manifestTag := range opts.ManifestTags {
		// Generate the tagged uri to work with
		imageName := fmt.Sprintf("%s:%s", manifestListUri.ShortName(), manifestTag)
		manifestTagUri, err := reference.NewUri(fmt.Sprintf("%s/%s", opts.Registry.Namespace, imageName), &reference.UriOptions{
			Registry: opts.Registry.ServerAddress,
			Official: opts.Official,
		})
		if err != nil {
			return err
		}

		// Create the manifest
		if err := createManifest(ctx, manifestTagUri, opts); err != nil {
			return err
		}

		// Annotate the manifest
		if err := annotateManifest(ctx, manifestTagUri, opts); err != nil {
			return err
		}

		// Push the manifest
		if opts.Push {
			manifestsToPush = append(manifestsToPush, manifestTagUri)
			log.Debugf("%s manifest staged to push", manifestTagUri.Remote())
		}
	}

	if opts.Push {
		// push all the manifests to the registry
		for _, manifest := range manifestsToPush {
			pushOpts := PushManifestOptions{
				Purge:  true,
				DryRun: opts.DryRun,
				Debug:  opts.Debug,
			}
			if err := pushManifest(ctx, manifest, pushOpts); err != nil {
				log.Errorf("pushing the manifest '%s' failed: %v", manifest.Remote(), err)
			}
		}
	}

	// logout of the registry
	if err := dockerLogout(ctx, opts.Registry.ServerAddress, opts.DryRun); err != nil {
		return err
	}

	return nil
}
