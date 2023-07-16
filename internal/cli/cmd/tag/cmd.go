package tag

import (
	"context"
	"tugboat/internal/cli"
	"tugboat/internal/clients/docker"
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
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return err
	}

	compiledSourceImage, err := tmpl.CompileString(args[0], opts)
	if err != nil {
		return err
	}

	compiledTags, err := tmpl.CompileStringSlice(opts.Tag.Tags, opts)
	if err != nil {
		return err
	}

	// create the registry
	registry, err := registry.NewRegistry(
		opts.Global.Docker.Registry,
		opts.Global.Docker.Namespace,
		opts.Global.Docker.Username,
		opts.Global.Docker.Password,
	)
	if err != nil {
		return err
	}

	log.Debugf("compiledSourceImage: %s", compiledSourceImage)
	log.Debugf("compiledTags: %s", compiledTags)

	tagOptions := image.TagOptions{
		SourceImage:            compiledSourceImage,
		Tags:                   compiledTags,
		Push:                   opts.Tag.Push,
		SupportedArchitectures: opts.Image.SupportedArchitectures,
		Registry:               registry,
		Official:               opts.Global.Official,
		DryRun:                 opts.Global.DryRun,
		Debug:                  opts.Global.Debug,
		ArchOption:             flags.DefaultArchOption,
	}

	if err := image.ImageTag(ctx, client, tagOptions); err != nil {
		return err
	}
	return nil
}
