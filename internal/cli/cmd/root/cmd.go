package root

import (
	"tugboat/internal/version"

	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "tugboat",
		Version:       version.GetFullVersionWithArch(),
		Short:         "Build multi-arch images",
		Long:          rootDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	return cmd
}

var rootDescription = `A tool to build and publish multi-architecture container images`
