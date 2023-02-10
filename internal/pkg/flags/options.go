package flags

const DefaultArchOption = "prepend"

type Options struct {
	Global   GlobalOptions
	Build    BuildOptions
	Image    ImageOptions
	Manifest ManifestOptions
	Tag      TagOptions
	Version  VersionOptions
}

type GlobalOptions struct {
	ConfigFile string
	Docker     DockerOptions
	Debug      bool
	DryRun     bool
	Official   bool
	Git        Git
	Version    Version
}

type BuildOptions struct {
	BuildArgs []string
	Context   string
	File      string
	Tags      []string
	Push      bool
	Pull      bool
	NoCache   bool
}

type DockerOptions struct {
	Registry  string
	Namespace string
	Username  string
	Password  string
}

type ImageOptions struct {
	Name                   string
	SupportedArchitectures []string
	Version                string
}

type ManifestOptions struct {
	Create ManifestCreateOptions
	Push   ManifestPushOptions
}

type ManifestCreateOptions struct {
	Tags   []string
	Latest bool
	Push   bool
}

type ManifestPushOptions struct {
	Purge bool
}

type TagOptions struct {
	Tags []string
	Push bool
}

type VersionOptions struct {
	Short bool
}

type Git struct {
	Branch      string
	Commit      string
	FullCommit  string
	ShortCommit string
	Tag         string
}

type Version struct {
	App    string
	Config string
}
