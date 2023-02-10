package flags

var (
	BuildArgsFlag = Flag{
		Name:       "build-args",
		ConfigName: "build.args",
		Value:      []string{},
		Usage:      "Set build-time variables in a comma separated string (i.e. --build-args FOO=bar,BAR=foo)",
	}
	ContextFlag = Flag{
		Name:       "context",
		ConfigName: "build.context",
		Value:      ".",
		Usage:      "Docker image build path to use",
	}
	FileFlag = Flag{
		Name:       "file",
		ConfigName: "build.file",
		Shorthand:  "f",
		Value:      "Dockerfile",
		Usage:      "Name of the Dockerfile (Default is 'PATH/Dockerfile')",
	}
	TagsFlag = Flag{
		Name:       "tag",
		ConfigName: "build.tags",
		Shorthand:  "t",
		Value:      []string{},
		Usage:      "Name of the image and optionally a tag in the 'name:tag' format",
	}
	PushFlag = Flag{
		Name:       "push",
		ConfigName: "build.push",
		Value:      false,
		Usage:      "Push the image to a container registry after building",
	}
	PullFlag = Flag{
		Name:       "pull",
		ConfigName: "build.pull",
		Value:      false,
		Usage:      "Always attempt to pull a newer version of the image",
	}
	NoCacheFlag = Flag{
		Name:       "no-cache",
		ConfigName: "build.no-cache",
		Value:      false,
		Usage:      "Do not use cache when building the image",
	}
)

type BuildFlagGroup struct {
	BuildArgsFlag *Flag
	ContextFlag   *Flag
	FileFlag      *Flag
	TagsFlag      *Flag
	PushFlag      *Flag
	PullFlag      *Flag
	NoCacheFlag   *Flag
}

func NewBuildFlagsGroup() *BuildFlagGroup {
	return &BuildFlagGroup{
		BuildArgsFlag: &BuildArgsFlag,
		ContextFlag:   &ContextFlag,
		FileFlag:      &FileFlag,
		TagsFlag:      &TagsFlag,
		PushFlag:      &PushFlag,
		PullFlag:      &PullFlag,
		NoCacheFlag:   &NoCacheFlag,
	}
}

func (f *BuildFlagGroup) Name() string {
	return "Build"
}

func (f *BuildFlagGroup) Flags() []*Flag {
	return []*Flag{f.BuildArgsFlag, f.ContextFlag, f.FileFlag, f.TagsFlag, f.PushFlag, f.PullFlag, f.NoCacheFlag}
}

func (f *BuildFlagGroup) ToOptions() BuildOptions {
	buildOpts := BuildOptions{
		BuildArgs: getStringSlice(f.BuildArgsFlag),
		Context:   getString(f.ContextFlag),
		File:      getString(f.FileFlag),
		Tags:      getStringSlice(f.TagsFlag),
		Push:      getBool(f.PushFlag),
		Pull:      getBool(f.PullFlag),
		NoCache:   getBool(f.NoCacheFlag),
	}

	return buildOpts
}
