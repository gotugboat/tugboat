package manifest

import (
	"fmt"
	"tugboat/internal/cli"
	"tugboat/internal/pkg/flags"

	"github.com/spf13/cobra"
)

func NewManifestCommand(globalFlags *flags.GlobalFlagGroup) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "manifest COMMAND",
		Short: "Manage image manifests",
		Long:  manifestDescription,
		Args:  cli.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.UsageString())
		},
	}

	cmd.AddCommand(
		newCreateCommand(globalFlags),
	)

	return cmd
}

var manifestDescription = `Manage image manifests`
