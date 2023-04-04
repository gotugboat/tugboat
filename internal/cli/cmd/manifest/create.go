package manifest

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

func newCreateCommand(globalFlags *flags.GlobalFlagGroup) *cobra.Command {
	manifestCreateFlags := flags.NewManifestCreateFlagGroup()
	imageFlags := flags.NewImageFlagsGroup()

	cmd := &cobra.Command{
		Use:   "create IMAGE",
		Short: "Create a local annotated manifest list for pushing to a registry",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := flags.ToOptions(globalFlags, manifestCreateFlags, imageFlags)
			return createManifest(opts, args)
		},
	}

	flags.AddFlags(cmd, manifestCreateFlags, imageFlags)
	flags.Bind(cmd, manifestCreateFlags)
	flags.Bind(cmd, imageFlags)

	return cmd
}

func createManifest(opts *flags.Options, args []string) error {
	log.Debugf("Manifest Create Options: %+v", opts)
	log.Debugf("Manifest Create Args: %+v", args)

	ctx := context.Background()
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return err
	}

	compiledManifestList, err := tmpl.CompileString(args[0], opts)
	if err != nil {
		return err
	}

	manifestTags, err := getManifestTags(opts)
	if err != nil {
		return err
	}

	manifestCreateOpts := image.ManifestCreateOptions{
		ManifestList:           compiledManifestList,
		ManifestTags:           manifestTags,
		Push:                   opts.Manifest.Create.Push,
		SupportedArchitectures: opts.Image.SupportedArchitectures,
		Registry: image.NewRegistry(
			opts.Global.Docker.Registry,
			opts.Global.Docker.Namespace,
			opts.Global.Docker.Username,
			opts.Global.Docker.Password,
		),
		Official:   opts.Global.Official,
		DryRun:     opts.Global.DryRun,
		Debug:      opts.Global.Debug,
		ArchOption: flags.DefaultArchOption,
	}

	if err := image.ManifestCreate(ctx, client, manifestCreateOpts); err != nil {
		return err
	}
	return nil
}

func getManifestTags(opts *flags.Options) ([]string, error) {
	// Compile the tags
	compiledTags, err := tmpl.CompileStringSlice(opts.Manifest.Create.Tags, opts)
	if err != nil {
		return nil, err
	}

	// build the list of manifests
	var manifestTags []string
	manifestTags = append(manifestTags, compiledTags...)
	if opts.Manifest.Create.Latest {
		manifestTags = append(manifestTags, "latest")
	}

	return manifestTags, nil
}
