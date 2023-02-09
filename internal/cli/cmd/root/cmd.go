package root

import (
	"tugboat/internal/pkg/flags"
	"tugboat/internal/version"

	"github.com/spf13/cobra"
)

func NewRootCommand(globalFlags *flags.GlobalFlagGroup) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "tugboat",
		Version:       version.GetFullVersionWithArch(),
		Short:         "Build multi-arch images",
		Long:          rootDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Assign global flags
	flags.AddFlags(cmd, globalFlags)
	flags.Bind(cmd, globalFlags)

	return cmd
}

var rootDescription = `A tool to build and publish multi-architecture container images`
