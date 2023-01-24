package commands

import (
	"tugboat/internal/cli/cmd/root"
	"tugboat/internal/cli/cmd/version"

	"github.com/spf13/cobra"
)

func NewCli() *cobra.Command {
	rootCmd := root.NewRootCommand()

	// add the commands to the root command
	addCommands(rootCmd)

	return rootCmd
}

// Adds all the commands to the given root command
func addCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		// version
		version.NewVersionCommand(),
	)
}
