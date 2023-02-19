package manifest

import (
	"tugboat/internal/cli"
	"tugboat/internal/pkg/flags"

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

	return nil
}

func getManifestTags(opts *flags.Options) ([]string, error) {
	// build the list of manifests
	var manifestTags []string
	manifestTags = append(manifestTags, opts.Manifest.Create.Tags...)
	if opts.Manifest.Create.Latest {
		manifestTags = append(manifestTags, "latest")
	}

	return manifestTags, nil
}
