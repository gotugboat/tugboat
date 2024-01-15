package docker

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"tugboat/internal/clients/docker"
	"tugboat/internal/driver"
	"tugboat/internal/pkg/reference"
	"tugboat/internal/registry"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

// DockerDriver implements the Driver interface for building containers with Docker
type DockerDriver struct {
	Debug           bool
	DryRun          bool
	Official        bool
	ArchitectureTag string
	client          *client.Client
	registry        *registry.Registry
}

// NewDockerDriver creates a new instance of DockerDriver
func NewDockerDriver(opts driver.DriverOptions) (*DockerDriver, error) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return nil, err
	}

	return &DockerDriver{
		Debug:           opts.Debug,
		DryRun:          opts.DryRun,
		Official:        opts.Official,
		ArchitectureTag: opts.ArchitectureTag,
		client:          client,
		registry:        opts.Registry,
	}, nil
}

func (d *DockerDriver) BuildImage(ctx context.Context, opts driver.BuildOptions) (io.ReadCloser, error) {
	buildUris, err := driver.GenerateAllUris(d.registry.ServerAddress, d.registry.Namespace, opts.Tags, d.Official, reference.ArchOption(d.ArchitectureTag))
	if err != nil {
		return nil, err
	}

	return buildImage(ctx, buildUris, d.Official, d.DryRun, d.Debug, d.ArchitectureTag, opts)
}

func (d *DockerDriver) PullImage(ctx context.Context, image string) (io.ReadCloser, error) {
	uri, err := d.GetUri(image)
	if err != nil {
		return nil, err
	}

	log.Infof("Pulling %s", uri.Remote())

	if d.DryRun {
		return nil, nil
	}

	encodedRegistryAuth, err := encodeRegistryCredentials(d.registry)
	if err != nil {
		return nil, err
	}

	pullOpts := types.ImagePullOptions{
		RegistryAuth: encodedRegistryAuth,
	}

	response, err := d.client.ImagePull(ctx, image, pullOpts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (d *DockerDriver) PullImageWithArch(ctx context.Context, image string, architecture string) (io.ReadCloser, error) {
	uri, err := d.GetUriWithArch(image, architecture)
	if err != nil {
		return nil, err
	}

	response, err := d.pullImage(ctx, uri.Remote())
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (d *DockerDriver) pullImage(ctx context.Context, uri string) (io.ReadCloser, error) {
	log.Infof("Pulling %s", uri)

	if d.DryRun {
		return nil, nil
	}

	encodedRegistryAuth, err := encodeRegistryCredentials(d.registry)
	if err != nil {
		return nil, err
	}

	pullOpts := types.ImagePullOptions{
		RegistryAuth: encodedRegistryAuth,
	}

	response, err := d.client.ImagePull(ctx, uri, pullOpts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (d *DockerDriver) PushImage(ctx context.Context, image string) (io.ReadCloser, error) {
	uri, err := d.GetUri(image)
	if err != nil {
		return nil, err
	}

	log.Infof("Pushing %s", uri.Remote())

	if d.DryRun {
		return nil, nil
	}

	encodedRegistryAuth, err := encodeRegistryCredentials(d.registry)
	if err != nil {
		return nil, err
	}

	pushOpts := types.ImagePushOptions{
		RegistryAuth: encodedRegistryAuth,
	}

	response, err := d.client.ImagePush(ctx, uri.Remote(), pushOpts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (d *DockerDriver) PushImageWithArch(ctx context.Context, image string, architecture string) (io.ReadCloser, error) {
	uri, err := d.GetUriWithArch(image, architecture)
	if err != nil {
		return nil, err
	}

	response, err := d.pushImage(ctx, uri.Remote())
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (d *DockerDriver) pushImage(ctx context.Context, uri string) (io.ReadCloser, error) {
	log.Infof("Pushing %s", uri)

	if d.DryRun {
		return nil, nil
	}

	encodedRegistryAuth, err := encodeRegistryCredentials(d.registry)
	if err != nil {
		return nil, err
	}

	pushOpts := types.ImagePushOptions{
		RegistryAuth: encodedRegistryAuth,
	}

	response, err := d.client.ImagePush(ctx, uri, pushOpts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (d *DockerDriver) TagImage(ctx context.Context, sourceImage string, targetTag string) (string, error) {
	sourceUri, err := d.GetUri(sourceImage)
	if err != nil {
		return "", err
	}

	targetImage := fmt.Sprintf("%v:%v", sourceUri.ShortName(), targetTag)
	targetUri, err := d.GetUri(targetImage)
	if err != nil {
		return "", err
	}

	log.Infof("Tagging %v as %v", sourceUri.Remote(), targetUri.Remote())

	if d.DryRun {
		return targetUri.Remote(), nil
	}

	if err := d.client.ImageTag(ctx, sourceUri.Remote(), targetUri.Remote()); err != nil {
		return "", err
	}

	return targetUri.Remote(), nil
}

func (d *DockerDriver) TagImageWithArch(ctx context.Context, sourceImage string, targetTag string, architecture string) (string, error) {
	sourceUri, err := d.GetUriWithArch(sourceImage, architecture)
	if err != nil {
		return "", err
	}

	targetImage := fmt.Sprintf("%v:%v", sourceUri.ShortName(), targetTag)
	targetUri, err := d.GetUriWithArch(targetImage, architecture)
	if err != nil {
		return "", err
	}

	log.Infof("Tagging %v as %v", sourceUri.Remote(), targetUri.Remote())

	if d.DryRun {
		return targetUri.Remote(), nil
	}

	if err := d.client.ImageTag(ctx, sourceUri.Remote(), targetUri.Remote()); err != nil {
		return "", err
	}

	return targetUri.Remote(), nil
}

func (d *DockerDriver) CreateManifest(ctx context.Context, opts driver.ManifestCreateOptions) (io.ReadCloser, error) {
	// login to the registry
	if err := d.login(ctx); err != nil {
		return nil, err
	}

	// Generate the manifests for each desired tag
	for _, manifestTag := range opts.ManifestTags {
		// Generate the tagged uri to work with
		imageName := fmt.Sprintf("%s:%s", opts.ManifestList, manifestTag)
		manifestTagUri, err := reference.NewUri(fmt.Sprintf("%s/%s", d.registry.Namespace, imageName), &reference.UriOptions{
			Registry: d.registry.ServerAddress,
			Official: d.Official,
		})
		if err != nil {
			return nil, err
		}

		// Create the manifest
		if err := createManifest(
			ctx, manifestTagUri, opts.SupportedArchitectures, d.Official, d.ArchitectureTag, d.DryRun, d.Debug,
		); err != nil {
			return nil, err
		}

		// Annotate the manifest
		if err := annotateManifest(
			ctx, manifestTagUri, opts.SupportedArchitectures, d.Official, d.ArchitectureTag, d.DryRun, d.Debug,
		); err != nil {
			return nil, err
		}
	}

	// logout of the registry
	if err := d.logout(ctx); err != nil {
		return nil, err
	}

	return nil, nil
}

func (d *DockerDriver) PushManifest(ctx context.Context, manifestList string, opts driver.ManifestPushOptions) error {
	// login to the registry
	if err := d.login(ctx); err != nil {
		return err
	}

	// Generate the tagged uri to work with
	manifestUri, err := reference.NewUri(fmt.Sprintf("%s/%s", d.registry.Namespace, manifestList), &reference.UriOptions{
		Registry: d.registry.ServerAddress,
		Official: d.Official,
	})
	if err != nil {
		return err
	}

	// Push the manifest to the registry
	if err := pushManifest(ctx, manifestUri, d.DryRun, d.Debug, opts); err != nil {
		log.Errorf("pushing the manifest '%s' failed: %v", manifestUri.Remote(), err)
	}

	// logout of the registry
	if err := d.logout(ctx); err != nil {
		return err
	}

	return nil
}

func (d *DockerDriver) RemoveManifest(ctx context.Context, manifestLists []string) error {
	allReferences := []*reference.Reference{}

	for _, manifestList := range manifestLists {
		// Generate the tagged uri to work with
		manifestUri, err := reference.NewUri(fmt.Sprintf("%s/%s", d.registry.Namespace, manifestList), &reference.UriOptions{
			Registry: d.registry.ServerAddress,
			Official: d.Official,
		})
		if err != nil {
			return err
		}

		allReferences = append(allReferences, manifestUri)
	}

	if err := removeManifests(ctx, allReferences, d.DryRun, d.Debug); err != nil {
		return err
	}

	return nil
}

func (d *DockerDriver) GetUri(tag string) (*reference.Reference, error) {
	uri, err := driver.GenerateUri(d.registry.ServerAddress, d.registry.Namespace, tag, d.Official, reference.ArchOption(d.ArchitectureTag))
	if err != nil {
		return nil, err
	}

	return uri, nil
}

func (d *DockerDriver) GetUriWithArch(tag string, arch string) (*reference.Reference, error) {
	uri, err := driver.GenerateUriWithArch(d.registry.ServerAddress, d.registry.Namespace, tag, d.Official, reference.ArchOption(d.ArchitectureTag), arch)
	if err != nil {
		return nil, err
	}

	return uri, nil
}

func (d *DockerDriver) login(ctx context.Context) error {
	log.Infof("Logging into %v as %v", d.registry.ServerAddress, d.registry.User.Name)

	if d.DryRun {
		return nil
	}

	loginCmd := []string{"login", "--username", d.registry.User.Name, "--password", d.registry.User.Password, d.registry.ServerAddress}

	output, err := exec.Command("docker", loginCmd...).Output()
	if err != nil {
		log.Errorf("Docker login failed: %v", err)
		return ErrDockerLogin
	}
	log.Info(strings.TrimSuffix(string(output), "\n"))

	return nil
}

func (d *DockerDriver) logout(ctx context.Context) error {
	log.Infof("Logging out of %v", d.registry.ServerAddress)

	if d.DryRun {
		return nil
	}

	loginOutCmd := []string{"logout", d.registry.ServerAddress}

	output, err := exec.Command("docker", loginOutCmd...).Output()
	if err != nil {
		return ErrDockerLogout
	}
	log.Info(strings.TrimSuffix(string(output), "\n"))

	return nil
}
