package build

import (
	"tugboat/internal/cli"
	"tugboat/internal/pkg/flags"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewBuildCommand(globalFlags *flags.GlobalFlagGroup) *cobra.Command {
	buildFlags := flags.NewBuildFlagsGroup()

	cmd := &cobra.Command{
		Use:   "build",
		Short: "Build a container",
		Long:  buildDescription,
		Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := flags.ToOptions(globalFlags, buildFlags)
			return runBuild(opts)
		},
	}

	flags.AddFlags(cmd, buildFlags)
	flags.Bind(cmd, buildFlags)

	return cmd
}

var buildDescription = `Build an image from a Dockerfile`

func runBuild(opts *flags.Options) error {
	log.Info("Running build")
	return nil
}
