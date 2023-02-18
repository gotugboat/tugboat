package flags

var (
	ImageNameFlag = Flag{
		Name:       "",
		ConfigName: "image.name",
		Value:      "",
		Usage:      "Set the image name",
	}
	ImageArchitecturesFlag = Flag{
		Name:       "architectures",
		Shorthand:  "a",
		ConfigName: "image.supported-architectures",
		Value:      []string{},
		Usage:      "Define the supported image architectures",
	}
	ImageVersionFlag = Flag{
		Name:       "",
		ConfigName: "image.version",
		Value:      "",
		Usage:      "Define the version of your application",
	}
)

type ImageFlagGroup struct {
	ImageNameFlag          *Flag
	ImageArchitecturesFlag *Flag
	ImageVersionFlag       *Flag
}

func NewImageFlagsGroup() *ImageFlagGroup {
	return &ImageFlagGroup{
		ImageNameFlag:          &ImageNameFlag,
		ImageArchitecturesFlag: &ImageArchitecturesFlag,
		ImageVersionFlag:       &ImageVersionFlag,
	}
}

func (f *ImageFlagGroup) Name() string {
	return "Image"
}

func (f *ImageFlagGroup) Flags() []*Flag {
	return []*Flag{f.ImageNameFlag, f.ImageArchitecturesFlag, f.ImageVersionFlag}
}

func (f *ImageFlagGroup) ToOptions() ImageOptions {
	sanitizedVersion := getString(f.ImageVersionFlag)

	opts := ImageOptions{
		Name:                   getString(f.ImageNameFlag),
		SupportedArchitectures: getStringSlice(f.ImageArchitecturesFlag),
		Version:                sanitizedVersion,
	}

	return opts
}
