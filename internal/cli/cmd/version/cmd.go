package version

import (
	"fmt"
	"tugboat/internal/cli"
	"tugboat/internal/pkg/flags"
	"tugboat/internal/version"

	"github.com/spf13/cobra"
)

func NewVersionCommand(globalFlags *flags.GlobalFlagGroup) *cobra.Command {
	versionFlags := flags.NewVersionFlagsGroup()

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Args:  cli.NoArgs,
		Long:  versionDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := flags.ToOptions(globalFlags, versionFlags)
			return runVersion(opts)
		},
	}

	flags.AddFlags(cmd, versionFlags)
	flags.Bind(cmd, versionFlags)

	return cmd
}

var versionDescription = `Show the Tugboat version information`

func runVersion(opts *flags.Options) error {
	if opts.Version.Short {
		fmt.Printf("%v\n", version.GetVersion())
	} else {
		fmt.Printf("tugboat version %s\n", version.GetFullVersionWithArch())
	}
	return nil
}
