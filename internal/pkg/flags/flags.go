package flags

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type Flag struct {
	// Name is for CLI flag and environment variable.
	// If this field is empty, it will be available only in config file.
	Name string

	// ConfigName is a key in config file. It is also used as the key for viper.
	ConfigName string

	// Shorthand is a shorthand letter.
	Shorthand string

	// Value is the default value. It must be filled to determine the flag type.
	Value interface{}

	// Usage explains how to use the flag.
	Usage string

	// Persistent represents if the flag is persistent
	Persistent bool

	// Deprecated represents if the flag is deprecated
	Deprecated bool
}

type FlagGroup interface {
	Name() string
	Flags() []*Flag
}

func AddFlags(cmd *cobra.Command, f ...FlagGroup) {
	for _, group := range f {
		for _, flag := range group.Flags() {
			addFlag(cmd, flag)
		}
	}
}

func Bind(cmd *cobra.Command, f FlagGroup) error {
	for _, flag := range f.Flags() {
		if err := bind(cmd, flag); err != nil {
			return errors.Errorf("flag groups: %v", err)
		}
	}
	return nil
}

func ToOptions(globalFlags *GlobalFlagGroup, f ...FlagGroup) *Options {
	opts := &Options{
		Global: globalFlags.ToOptions(),
	}

	for _, flagGroup := range f {
		switch v := flagGroup.(type) {
		case *BuildFlagGroup:
			opts.Build = v.ToOptions()
		case *VersionFlagGroup:
			opts.Version = v.ToOptions()
		}
	}

	return opts
}
