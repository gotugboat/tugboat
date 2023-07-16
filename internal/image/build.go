package image

import (
	"context"
	"fmt"
	"io"
	"strings"
	"tugboat/internal/pkg/reference"
	"tugboat/internal/term"

	"github.com/docker/cli/cli/command/image/build"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/idtools"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type BuildOptions struct {
	Context    string
	Dockerfile string
	Tags       []string
	BuildArgs  []string
	Push       bool
	Rm         bool
	Pull       bool
	NoCache    bool

	Registry Registry
	Official bool
	DryRun   bool
	Debug    bool

	ArchOption string
}

func ImageBuild(ctx context.Context, client *client.Client, opts BuildOptions) error {
	var buildContext io.ReadCloser

	// Create a tarball from the context directory
	if buildContext == nil && !opts.DryRun {
		excludes, err := build.ReadDockerignore(opts.Context)
		if err != nil {
			return err
		}

		if err := build.ValidateContextDirectory(opts.Context, excludes); err != nil {
			return errors.Wrap(err, "error checking context")
		}

		excludes = build.TrimBuildFilesFromExcludes(excludes, opts.Dockerfile, false)
		buildContext, err = archive.TarWithOptions(opts.Context, &archive.TarOptions{
			ExcludePatterns: excludes,
			ChownOpts:       &idtools.Identity{UID: 0, GID: 0},
		})
		if err != nil {
			return err
		}
	}

	// Build the image using the context directory
	buildUris, err := generateAllUris(opts)
	if err != nil {
		return err
	}
	buildOpts := imageBuildOptions(buildUris, opts)

	log.Infof("Building %s using %s/%s", buildOpts.Tags[0], opts.Context, opts.Dockerfile)

	if !opts.DryRun {
		response, err := client.ImageBuild(ctx, buildContext, buildOpts)
		if err != nil {
			return errors.Wrap(err, "Image build error")
		}
		defer response.Body.Close()

		if err := term.DisplayResponse(response.Body); err != nil {
			return err
		}
	}

	// Push the image and any additional tags
	if opts.Push {
		for _, uri := range buildUris {
			if err := push(ctx, client, opts.Registry, uri.Remote(), opts.DryRun); err != nil {
				log.Error(err)
			}
		}
	}

	return nil
}

func generateAllUris(opts BuildOptions) ([]*reference.Reference, error) {
	buildTags := []*reference.Reference{}

	for _, tag := range opts.Tags {
		taggedUri, err := reference.NewUri(fmt.Sprintf("%s/%s", opts.Registry.Namespace, tag), &reference.UriOptions{
			Registry:   opts.Registry.ServerAddress,
			Official:   opts.Official,
			ArchOption: toArchOption(opts.ArchOption),
		})
		if err != nil {
			return nil, errors.Errorf("%v", err)
		}
		buildTags = append(buildTags, taggedUri)
	}

	return buildTags, nil
}

func imageBuildOptions(buildUris []*reference.Reference, opts BuildOptions) types.ImageBuildOptions {
	// Prepare the tags
	var buildTags []string
	for _, uri := range buildUris {
		buildTags = append(buildTags, uri.Remote())
	}

	// Prepare the build arguments
	buildArgs := make(map[string]*string)
	for _, pair := range opts.BuildArgs {
		// Split the pair on the equals sign
		kv := strings.Split(pair, "=")
		// Add the key-value pair to the map
		buildArgs[kv[0]] = &kv[1]
	}

	return types.ImageBuildOptions{
		Dockerfile: opts.Dockerfile,
		Tags:       buildTags,
		NoCache:    opts.NoCache,
		Remove:     opts.Rm,
		PullParent: opts.Pull,
		BuildArgs:  buildArgs,
	}
}
