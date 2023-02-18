package root

import (
	"os"
	"tugboat/internal/config"
	"tugboat/internal/pkg/flags"
	"tugboat/internal/pkg/logging"
	"tugboat/internal/version"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewRootCommand(globalFlags *flags.GlobalFlagGroup) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tugboat",
		Version: version.GetFullVersionWithArch(),
		Short:   "Build multi-arch images",
		Long:    rootDescription,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Load the config
			configPath := viper.GetString(flags.ConfigFileFlag.ConfigName)
			if err := config.LoadConfig(configPath); err != nil {
				return err
			}

			globalOptions := globalFlags.ToOptions()

			logging.Initialize(os.Stderr, globalOptions.Debug)

			if globalOptions.DryRun {
				log.Warn("Dry run in progress, nothing will be executed")
			}

			return nil
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Assign global flags
	flags.AddFlags(cmd, globalFlags)
	flags.Bind(cmd, globalFlags)

	return cmd
}

var rootDescription = `A tool to build and publish multi-architecture container images`
