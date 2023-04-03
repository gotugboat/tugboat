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
	// docker flags
	DockerRegistryFlag = Flag{
		Name:       "docker-registry",
		ConfigName: "docker.registry",
		Value:      "docker.io",
		Usage:      "The docker registry to use",
		Persistent: true,
	}
	DockerNamespaceFlag = Flag{
		Name:       "docker-namespace",
		ConfigName: "docker.namespace",
		Value:      "",
		Usage:      "The namespace in the docker registry to use (DockerHub username if using DockerHub, any if using private registry)",
		Persistent: true,
	}
	DockerUsernameFlag = Flag{
		Name:       "docker-user",
		ConfigName: "docker.user",
		Value:      "",
		Usage:      "The username credential with access to the registry",
		Persistent: true,
	}
	DockerPasswordFlag = Flag{
		Name:       "docker-pass",
		ConfigName: "docker.pass",
		Value:      "",
		Usage:      "The password credential with access to the registry",
		Persistent: true,
	}
)

type DockerFlagGroup struct {
	IsExperimentalFlag *Flag
	RegistryFlag       *Flag
	NamespaceFlag      *Flag
	UsernameFlag       *Flag
	PasswordFlag       *Flag
}

type GlobalFlagGroup struct {
	ConfigFileFlag  *Flag
	DebugFlag       *Flag
	DryRunFlag      *Flag
	DockerFlagGroup *DockerFlagGroup
	OfficialFlag    *Flag
}

func NewGlobalFlagGroup() *GlobalFlagGroup {
	return &GlobalFlagGroup{
		ConfigFileFlag: &ConfigFileFlag,
		DebugFlag:      &DebugFlag,
		DryRunFlag:     &DryRunFlag,
		DockerFlagGroup: &DockerFlagGroup{
			RegistryFlag:  &DockerRegistryFlag,
			NamespaceFlag: &DockerNamespaceFlag,
			UsernameFlag:  &DockerUsernameFlag,
			PasswordFlag:  &DockerPasswordFlag,
		},
		OfficialFlag: &OfficialFlag,
	}
}

func (f *GlobalFlagGroup) Name() string {
	return "Global"
}

func (f *GlobalFlagGroup) Flags() []*Flag {
	return []*Flag{f.ConfigFileFlag, f.DebugFlag, f.DryRunFlag, f.OfficialFlag, f.DockerFlagGroup.RegistryFlag, f.DockerFlagGroup.NamespaceFlag, f.DockerFlagGroup.UsernameFlag, f.DockerFlagGroup.PasswordFlag}
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
		Docker: DockerOptions{
			Registry:  getString(f.DockerFlagGroup.RegistryFlag),
			Namespace: getString(f.DockerFlagGroup.NamespaceFlag),
			Username:  getString(f.DockerFlagGroup.UsernameFlag),
			Password:  getString(f.DockerFlagGroup.PasswordFlag),
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
