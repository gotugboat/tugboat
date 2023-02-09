package commands

import (
	"tugboat/internal/cli/cmd/root"
	"tugboat/internal/cli/cmd/version"
	"tugboat/internal/pkg/flags"

	"github.com/spf13/cobra"
)

func NewCli() *cobra.Command {
	globalFlags := flags.NewGlobalFlagGroup()
	rootCmd := root.NewRootCommand(globalFlags)

	// add the commands to the root command
	addCommands(rootCmd, globalFlags)

	return rootCmd
}

// Adds all the commands to the given root command
func addCommands(cmd *cobra.Command, globalFlags *flags.GlobalFlagGroup) {
	cmd.AddCommand(
		// version
		version.NewVersionCommand(globalFlags),
	)
}
