package flags

import (
	"context"
	"tugboat/internal/pkg/git"
	"tugboat/internal/version"
)

var (
	ConfigFileFlag = Flag{
		Name:       "config",
		ConfigName: "config",
		Shorthand:  "c",
		Value:      "",
		Usage:      "Custom path to a configuration file (optional)",
		Persistent: true,
	}
	DryRunFlag = Flag{
		Name:       "dry-run",
		ConfigName: "options.dry-run",
		Value:      false,
		Usage:      "Output what will happen, do not execute",
		Persistent: true,
	}
	DebugFlag = Flag{
		Name:       "debug",
		ConfigName: "options.debug",
		Value:      false,
		Usage:      "Run in debug mode",
		Persistent: true,
		Deprecated: true, // hide since it's not really intended to be used
	}
	OfficialFlag = Flag{
		Name:       "official",
		ConfigName: "publish.official",
		Value:      false,
		Usage:      "Mimic the official docker publish method for images in private registries",
		Persistent: true,
	}
	// registry flags
	RegistryUrlFlag = Flag{
		Name:       "registry",
		ConfigName: "registry.url",
		Value:      "docker.io",
		Usage:      "The registry to use",
		Persistent: true,
	}
	RegistryNamespaceFlag = Flag{
		Name:       "registry-namespace",
		ConfigName: "registry.namespace",
		Value:      "",
		Usage:      "The namespace in the registry to use (DockerHub username if using DockerHub, any if using private registry)",
		Persistent: true,
	}
	RegistryUsernameFlag = Flag{
		Name:       "registry-user",
		ConfigName: "registry.user",
		Value:      "",
		Usage:      "The username credential with access to the registry",
		Persistent: true,
	}
	RegistryPasswordFlag = Flag{
		Name:       "registry-password",
		ConfigName: "registry.password",
		Value:      "",
		Usage:      "The password credential with access to the registry",
		Persistent: true,
	}
	DriverNameFlag = Flag{
		Name:       "driver",
		ConfigName: "driver.name",
		Value:      "auto",
		Usage:      "The driver to use to manage containers",
		Persistent: true,
	}
)

type DriverFlagGroup struct {
	NameFlag *Flag
}

type RegistryFlagGroup struct {
	RegistryUrlFlag *Flag
	NamespaceFlag   *Flag
	UsernameFlag    *Flag
	PasswordFlag    *Flag
}

type GlobalFlagGroup struct {
	ConfigFileFlag    *Flag
	DebugFlag         *Flag
	DryRunFlag        *Flag
	DriverFlagGroup   *DriverFlagGroup
	RegistryFlagGroup *RegistryFlagGroup
	OfficialFlag      *Flag
}

func NewGlobalFlagGroup() *GlobalFlagGroup {
	return &GlobalFlagGroup{
		ConfigFileFlag: &ConfigFileFlag,
		DebugFlag:      &DebugFlag,
		DryRunFlag:     &DryRunFlag,
		DriverFlagGroup: &DriverFlagGroup{
			NameFlag: &DriverNameFlag,
		},
		RegistryFlagGroup: &RegistryFlagGroup{
			RegistryUrlFlag: &RegistryUrlFlag,
			NamespaceFlag:   &RegistryNamespaceFlag,
			UsernameFlag:    &RegistryUsernameFlag,
			PasswordFlag:    &RegistryPasswordFlag,
		},
		OfficialFlag: &OfficialFlag,
	}
}

func (f *GlobalFlagGroup) Name() string {
	return "Global"
}

func (f *GlobalFlagGroup) Flags() []*Flag {
	return []*Flag{f.ConfigFileFlag, f.DebugFlag, f.DryRunFlag, f.OfficialFlag, f.DriverFlagGroup.NameFlag, f.RegistryFlagGroup.RegistryUrlFlag, f.RegistryFlagGroup.NamespaceFlag, f.RegistryFlagGroup.UsernameFlag, f.RegistryFlagGroup.PasswordFlag}
}

func (f *GlobalFlagGroup) ToOptions() GlobalOptions {
	ctx := context.TODO()
	gitFullCommit, _ := git.Clean(git.Run(ctx, "rev-parse HEAD"))
	gitShortCommit, _ := git.Clean(git.Run(ctx, "log -1 --pretty=%h"))
	gitBranch, _ := git.Clean(git.Run(ctx, "rev-parse --abbrev-ref HEAD"))
	gitTag, _ := git.Clean(git.Run(ctx, "describe --tags"))

	opts := GlobalOptions{
		ConfigFile: getString(f.ConfigFileFlag),
		Debug:      getBool(f.DebugFlag),
		DryRun:     getBool(f.DryRunFlag),
		Driver: DriverOptions{
			Name: getString(f.DriverFlagGroup.NameFlag),
		},
		Registry: RegistryOptions{
			Url:       getString(f.RegistryFlagGroup.RegistryUrlFlag),
			Namespace: getString(f.RegistryFlagGroup.NamespaceFlag),
			Username:  getString(f.RegistryFlagGroup.UsernameFlag),
			Password:  getString(f.RegistryFlagGroup.PasswordFlag),
		},
		Official: getBool(f.OfficialFlag),
		Git: Git{
			Branch:      gitBranch,
			Commit:      gitFullCommit,
			FullCommit:  gitFullCommit,
			ShortCommit: gitShortCommit,
			Tag:         gitTag,
		},
		Version: Version{
			App: version.GetVersion(),
		},
	}

	return opts
}
