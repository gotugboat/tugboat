package flags

var (
	ManifestCreateForFlag = Flag{
		Name:       "for",
		ConfigName: "manifest.create.for",
		Value:      []string{},
		Usage:      "A list of tags to create a manifest for",
	}
	ManifestCreateLatestFlag = Flag{
		Name:       "latest",
		ConfigName: "manifest.create.latest",
		Value:      false,
		Usage:      "Create a manifest for the latest tag",
	}
	ManifestCreatePushFlag = Flag{
		Name:       "push",
		ConfigName: "manifest.create.push",
		Value:      false,
		Usage:      "Push the tagged images to an image registry",
	}
)

type ManifestCreateFlagGroup struct {
	ManifestCreateForFlag    *Flag
	ManifestCreateLatestFlag *Flag
	ManifestCreatePushFlag   *Flag
}

func NewManifestCreateFlagGroup() *ManifestCreateFlagGroup {
	return &ManifestCreateFlagGroup{
		ManifestCreateForFlag:    &ManifestCreateForFlag,
		ManifestCreateLatestFlag: &ManifestCreateLatestFlag,
		ManifestCreatePushFlag:   &ManifestCreatePushFlag,
	}
}

func (f *ManifestCreateFlagGroup) Name() string {
	return "ManifestCreate"
}

func (f *ManifestCreateFlagGroup) Flags() []*Flag {
	return []*Flag{f.ManifestCreateForFlag, f.ManifestCreateLatestFlag, f.ManifestCreatePushFlag}
}

func (f *ManifestCreateFlagGroup) ToOptions() ManifestCreateOptions {
	opts := ManifestCreateOptions{
		Tags:   getStringSlice(f.ManifestCreateForFlag),
		Latest: getBool(f.ManifestCreateLatestFlag),
		Push:   getBool(f.ManifestCreatePushFlag),
	}

	return opts
}
