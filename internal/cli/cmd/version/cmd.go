package version

import (
	"fmt"
	"tugboat/internal/version"

	"github.com/spf13/cobra"
)

type versionOpts struct {
	short bool
}

func NewVersionCommand() *cobra.Command {
	var opts versionOpts

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Args:  cobra.NoArgs,
		Long:  versionDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVersion(opts)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opts.short, "short", "", false, "Output a shortened version format")

	return cmd
}

var versionDescription = `Show the Tugboat version information`

func runVersion(opts versionOpts) error {
	if opts.short {
		fmt.Printf("%v\n", version.GetVersion())
	} else {
		fmt.Println(version.GetFullVersionWithArch())
	}
	return nil
}
