package flags

var (
	TagTagsFlag = Flag{
		Name:       "tags",
		ConfigName: "",
		Value:      []string{},
		Usage:      "A list of tags to reference the image by",
	}
	TagPushFlag = Flag{
		Name:       "push",
		ConfigName: "tag.push",
		Value:      false,
		Usage:      "Push the tagged images to a container registry",
	}
)

type TagFlagGroup struct {
	TagTagsFlag *Flag
	TagPushFlag *Flag
}

func NewTagFlagsGroup() *TagFlagGroup {
	return &TagFlagGroup{
		TagTagsFlag: &TagTagsFlag,
		TagPushFlag: &TagPushFlag,
	}
}

func (f *TagFlagGroup) Name() string {
	return "Tag"
}

func (f *TagFlagGroup) Flags() []*Flag {
	return []*Flag{f.TagTagsFlag, f.TagPushFlag}
}

func (f *TagFlagGroup) ToOptions() TagOptions {
	opts := TagOptions{
		Tags: getStringSlice(f.TagTagsFlag),
		Push: getBool(f.TagPushFlag),
	}

	return opts
}
