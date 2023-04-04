package build

import (
	"context"
	"tugboat/internal/cli"
	"tugboat/internal/clients/docker"
	"tugboat/internal/image"
	"tugboat/internal/pkg/flags"
	"tugboat/internal/pkg/tmpl"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewBuildCommand(globalFlags *flags.GlobalFlagGroup) *cobra.Command {
	buildFlags := flags.NewBuildFlagsGroup()
	imageFlags := flags.NewImageFlagsGroup()

	cmd := &cobra.Command{
		Use:   "build",
		Short: "Build a container",
		Long:  buildDescription,
		Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := flags.ToOptions(globalFlags, buildFlags, imageFlags)
			return runBuild(opts)
		},
	}

	flags.AddFlags(cmd, buildFlags)
	flags.Bind(cmd, buildFlags)

	return cmd
}

var buildDescription = `Build an image from a Dockerfile`

func runBuild(opts *flags.Options) error {
	log.Debugf("Build Options: %+v", opts)

	ctx := context.Background()
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return err
	}

	// Compile the tags using the template
	compiledTags, err := tmpl.CompileStringSlice(opts.Build.Tags, opts)
	if err != nil {
		return err
	}

	log.Debugf("compiledTags: %v", compiledTags)

	buildOpts := image.BuildOptions{
		Dockerfile: opts.Build.File,
		Context:    opts.Build.Context,
		Tags:       compiledTags,
		BuildArgs:  opts.Build.BuildArgs,
		Rm:         true,
		Pull:       opts.Build.Pull,
		NoCache:    opts.Build.NoCache,
		Push:       opts.Build.Push,
		DryRun:     opts.Global.DryRun,
		Debug:      opts.Global.Debug,
		Registry: image.NewRegistry(
			opts.Global.Docker.Registry,
			opts.Global.Docker.Namespace,
			opts.Global.Docker.Username,
			opts.Global.Docker.Password,
		),
		ArchOption: flags.DefaultArchOption,
	}

	if err := image.ImageBuild(ctx, client, buildOpts); err != nil {
		return err
	}
	return nil
}
