package image

import (
	"context"
	"tugboat/internal/term"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

func push(ctx context.Context, client *client.Client, registry Registry, image string, isDryRun bool) error {
	log.Infof("Pushing %s", image)

	if isDryRun {
		return nil
	}

	encodedRegistryAuth, err := encodeRegistryCredentials(registry)
	if err != nil {
		return err
	}

	pushOpts := types.ImagePushOptions{
		RegistryAuth: encodedRegistryAuth,
	}

	response, err := client.ImagePush(ctx, image, pushOpts)
	if err != nil {
		return err
	}
	defer response.Close()

	if err := term.DisplayResponse(response); err != nil {
		return err
	}

	return nil
}
