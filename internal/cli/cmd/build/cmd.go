package build

import (
	"context"
	"tugboat/internal/cli"
	"tugboat/internal/driver"
	"tugboat/internal/drivers"
	"tugboat/internal/image"
	"tugboat/internal/pkg/flags"
	"tugboat/internal/pkg/tmpl"
	"tugboat/internal/registry"

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

	// Compile the tags using the template
	compiledTags, err := tmpl.CompileStringSlice(opts.Build.Tags, opts)
	if err != nil {
		return err
	}

	log.Debugf("compiledTags: %v", compiledTags)

	registry, err := registry.NewRegistry(
		opts.Global.Registry.Url,
		opts.Global.Registry.Namespace,
		opts.Global.Registry.Username,
		opts.Global.Registry.Password,
	)
	if err != nil {
		return err
	}

	driverOpts := driver.DriverOptions{
		Registry:        registry,
		DryRun:          opts.Global.DryRun,
		Debug:           opts.Global.Debug,
		ArchitectureTag: flags.DefaultArchOption,
	}
	d, err := drivers.NewDriver(opts.Global.Driver.Name, driverOpts)
	if err != nil {
		return err
	}

	buildOpts := image.BuildOptions{
		Dockerfile: opts.Build.File,
		Context:    opts.Build.Context,
		Tags:       compiledTags,
		BuildArgs:  opts.Build.BuildArgs,
		Pull:       opts.Build.Pull,
		NoCache:    opts.Build.NoCache,
		Push:       opts.Build.Push,
	}
	if err := image.Build(ctx, d, buildOpts); err != nil {
		return err
	}

	return nil
}
