package flags

import (
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func addFlag(cmd *cobra.Command, flag *Flag) {
	if flag == nil || flag.Name == "" {
		return
	}

	var flags *pflag.FlagSet
	if flag.Persistent {
		flags = cmd.PersistentFlags()
	} else {
		flags = cmd.Flags()
	}

	// Evaluate what type of flag should being added
	switch v := flag.Value.(type) {
	case int:
		flags.IntP(flag.Name, flag.Shorthand, v, flag.Usage)
	case string:
		flags.StringP(flag.Name, flag.Shorthand, v, flag.Usage)
	case []string:
		flags.StringSliceP(flag.Name, flag.Shorthand, v, flag.Usage)
	case bool:
		flags.BoolP(flag.Name, flag.Shorthand, v, flag.Usage)
	case time.Duration:
		flags.DurationP(flag.Name, flag.Shorthand, v, flag.Usage)
	}

	if flag.Deprecated {
		flags.MarkHidden(flag.Name)
	}
}

func bind(cmd *cobra.Command, flag *Flag) error {
	if flag == nil {
		return nil
	} else if flag.Name == "" {
		// This flag is only available in the config file
		viper.SetDefault(flag.ConfigName, flag.Value)
		return nil
	}

	if flag.Persistent {
		if err := viper.BindPFlag(flag.ConfigName, cmd.PersistentFlags().Lookup(flag.Name)); err != nil {
			return err
		}
	} else {
		// fmt.Println("Config Name:", flag.ConfigName, " Flag Name:", cmd.Flags().Lookup(flag.Name))
		if err := viper.BindPFlag(flag.ConfigName, cmd.Flags().Lookup(flag.Name)); err != nil {
			return err
		}
	}

	// Bind the yaml configs to an env var (i.e log.format -> LOG_FORMAT)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	return nil
}

func getString(flag *Flag) string {
	if flag == nil {
		return ""
	}
	return viper.GetString(flag.ConfigName)
}

func getStringSlice(flag *Flag) []string {
	if flag == nil {
		return nil
	}
	v := viper.GetStringSlice(flag.ConfigName)

	// Separate env values containing a ','
	switch {
	case len(v) == 0: // no strings
		return nil
	case len(v) == 1 && strings.Contains(v[0], ","): // , separated string
		v = strings.Split(v[0], ",")
	}
	return v
}

func getBool(flag *Flag) bool {
	if flag == nil {
		return false
	}
	return viper.GetBool(flag.ConfigName)
}
