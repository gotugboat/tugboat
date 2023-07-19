package tag

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

func NewTagCommand(globalFlags *flags.GlobalFlagGroup) *cobra.Command {
	tagFlags := flags.NewTagFlagsGroup()
	imageFlags := flags.NewImageFlagsGroup()

	cmd := &cobra.Command{
		Use:   "tag SOURCE_IMAGE",
		Short: "Create a tag that refers to another image",
		Long:  tagDescription,
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := flags.ToOptions(globalFlags, tagFlags, imageFlags)
			return runTag(opts, args)
		},
	}

	flags.AddFlags(cmd, tagFlags, imageFlags)
	flags.Bind(cmd, tagFlags)
	flags.Bind(cmd, imageFlags)

	return cmd
}

var tagDescription = `Create a tag that refers to another image`

func runTag(opts *flags.Options, args []string) error {
	log.Debugf("Tag Options: %+v", opts)
	log.Debugf("Tag Args: %+v", args)

	ctx := context.Background()

	compiledSourceImage, err := tmpl.CompileString(args[0], opts)
	if err != nil {
		return err
	}

	compiledTags, err := tmpl.CompileStringSlice(opts.Tag.Tags, opts)
	if err != nil {
		return err
	}

	log.Debugf("compiledSourceImage: %s", compiledSourceImage)
	log.Debugf("compiledTags: %s", compiledTags)

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

	tagOptions := image.TagOptions{
		SourceImage:            compiledSourceImage,
		Tags:                   compiledTags,
		Push:                   opts.Tag.Push,
		SupportedArchitectures: opts.Image.SupportedArchitectures,
	}

	if err := image.Tag(ctx, d, tagOptions); err != nil {
		return err
	}
	return nil
}
