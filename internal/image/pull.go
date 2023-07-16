package image

import (
	"context"
	"tugboat/internal/registry"
	"tugboat/internal/term"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

func pull(ctx context.Context, client *client.Client, registry *registry.Registry, image string, isDryRun bool) error {
	log.Infof("Pulling %s", image)

	if isDryRun {
		return nil
	}

	encodedRegistryAuth, err := encodeRegistryCredentials(registry)
	if err != nil {
		return err
	}

	pullOpts := types.ImagePullOptions{
		RegistryAuth: encodedRegistryAuth,
	}

	response, err := client.ImagePull(ctx, image, pullOpts)
	if err != nil {
		return err
	}
	defer response.Close()

	if err := term.DisplayResponse(response); err != nil {
		return err
	}

	return nil
}
