package flags

var (
	ShortFlag = Flag{
		Name:       "short",
		ConfigName: "options.version.short",
		Value:      false,
		Usage:      "Output a shortened version format",
	}
)

type VersionFlagGroup struct {
	ShortFlag *Flag
}

func NewVersionFlagsGroup() *VersionFlagGroup {
	return &VersionFlagGroup{
		ShortFlag: &ShortFlag,
	}
}

func (f *VersionFlagGroup) Name() string {
	return "Version"
}

func (f *VersionFlagGroup) Flags() []*Flag {
	return []*Flag{f.ShortFlag}
}

func (f *VersionFlagGroup) ToOptions() VersionOptions {
	opts := VersionOptions{
		Short: getBool(f.ShortFlag),
	}

	return opts
}
