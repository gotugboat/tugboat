package driver

import (
	"context"
	"io"
	"tugboat/internal/pkg/reference"
	"tugboat/internal/registry"
)

// ImageBuilder represents the functionality for building container images
type ImageBuilder interface {
	BuildImage(ctx context.Context, options BuildOptions) (io.ReadCloser, error)
	PullImage(ctx context.Context, image string) (io.ReadCloser, error)
	PullImageWithArch(ctx context.Context, image string, architecture string) (io.ReadCloser, error)
}

// ImagePusher represents the functionality for pushing container images
type ImagePusher interface {
	PushImage(ctx context.Context, image string) (io.ReadCloser, error)
	PushImageWithArch(ctx context.Context, image string, architecture string) (io.ReadCloser, error)
	TagImage(ctx context.Context, sourceImage string, targetTag string) (string, error)
	TagImageWithArch(ctx context.Context, sourceImage string, targetTag string, architecture string) (string, error)
}

// ImageBuilderPusher represents the functionality for building and pushing container images
type ImageBuilderPusher interface {
	ImageBuilder
	ImagePusher
}

// ManifestSet represents the functionality for working with image manifests
type ManifestSet interface {
	CreateManifest(ctx context.Context, options ManifestCreateOptions) (io.ReadCloser, error)
	PushManifest(ctx context.Context, manifestList string, options ManifestPushOptions) error
	RemoveManifest(ctx context.Context, manifestLists []string) error
}

// Driver represents the complete set of functionality provided by a container driver
type Driver interface {
	ImageBuilder
	ImagePusher
	ManifestSet
	GetUri(tag string) (*reference.Reference, error)
}

type DriverOptions struct {
	Debug           bool
	DryRun          bool
	Official        bool
	ArchitectureTag string
	Registry        *registry.Registry
}

type BuildOptions struct {
	Context    string
	Dockerfile string
	Tags       []string
	BuildArgs  []string
	Rm         bool
	Pull       bool
	NoCache    bool
	Push       bool
}

type PushOptions struct {
	RegistryURL string
}

type ManifestCreateOptions struct {
	ManifestList           string
	ManifestTags           []string
	SupportedArchitectures []string
}

type ManifestPushOptions struct {
	Purge bool
}
