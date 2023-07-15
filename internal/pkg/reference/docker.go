package reference

import (
	_ "crypto/sha256"
	_ "crypto/sha512"
	"fmt"
	"strings"
	"tugboat/internal/pkg/reference/docker"
)

// Reference is an opaque object that include identifier such as a name, tag, repository, registry, etc...
type Reference struct {
	named docker.Named
	tag   string
}

type ArchOption string

var (
	ArchAppend  ArchOption = "append"
	ArchPrepend ArchOption = "prepend"
	ArchOmit    ArchOption = "omit"
)

type UriOptions struct {
	// The registry address for the image
	Registry string

	// Choose to mimic an official image naming format (i.e. registry/arch/image:tag)
	Official bool

	// Define an arch to use over the system architecture
	Arch string

	// Choose how the arch should be added to the image
	ArchOption ArchOption
}

// NewerUri returns a Reference from analyzing the given image and specified UriOptions.
func NewUri(image string, opts *UriOptions) (*Reference, error) {
	var uriString string

	arch := getArch()
	if opts.Arch != "" {
		arch = opts.Arch
	}

	// decide if this is already a uri

	uriString = generateUriString(image, opts.Registry, arch, opts.Official)

	ref, err := parse(uriString)
	if err != nil {
		return nil, err
	}

	// do not add an arch to the tag name when it contains a sha
	if strings.Contains(ref.tag, "@") {
		return ref, nil
	}

	if !opts.Official {
		// Prevent attaching an arch if the tag already contains the arch
		if strings.Contains(ref.Tag(), arch) {
			return ref, nil
		}

		switch opts.ArchOption {
		case ArchPrepend:
			ref.tag = fmt.Sprintf(":%s-%s", arch, ref.Tag())
		case ArchAppend:
			ref.tag = fmt.Sprintf(":%s-%s", ref.Tag(), arch)
		case ArchOmit:
			ref.tag = fmt.Sprintf(":%s", ref.Tag())
		}
	}
	return ref, nil
}

// Name returns the image's name. (ie: debian[:8.2])
func (r Reference) Name() string {
	return r.named.RemoteName() + r.tag
}

// ShortName returns the image's name (ie: debian)
func (r Reference) ShortName() string {
	return r.named.RemoteName()
}

// Tag returns the image's tag (or digest).
func (r Reference) Tag() string {
	if len(r.tag) > 1 {
		return r.tag[1:]
	}
	return ""
}

// Registry returns the image's registry. (ie: host[:port])
func (r Reference) Registry() string {
	return r.named.Hostname()
}

// Repository returns the image's repository. (ie: registry/name)
func (r Reference) Repository() string {
	return r.named.FullName()
}

// Remote returns the image's remote identifier. (ie: registry/name[:tag])
func (r Reference) Remote() string {
	return r.named.FullName() + r.tag
}

func clean(url string) string {
	s := url

	if strings.HasPrefix(url, "http://") {
		s = strings.Replace(url, "http://", "", 1)
	} else if strings.HasPrefix(url, "https://") {
		s = strings.Replace(url, "https://", "", 1)
	}

	return s
}

// parse returns a Reference from analyzing the given remote identifier.
func parse(remote string) (*Reference, error) {
	n, err := docker.ParseNamed(clean(remote))
	if err != nil {
		return nil, err
	}

	n = docker.WithDefaultTag(n)

	var t string
	switch x := n.(type) {
	case docker.Canonical:
		t = "@" + x.Digest().String()
	case docker.NamedTagged:
		t = ":" + x.Tag()
	}

	return &Reference{named: n, tag: t}, nil
}
