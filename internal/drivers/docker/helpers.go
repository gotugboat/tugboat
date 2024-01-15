package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"slices"
	"strings"
	"tugboat/internal/driver"
	"tugboat/internal/pkg/reference"
	"tugboat/internal/registry"
	"tugboat/internal/term"

	"github.com/docker/cli/cli/command/image/build"
	"github.com/docker/docker/api/types"
	registrytypes "github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/idtools"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var ErrDockerLogin = errors.New("docker login error")
var ErrDockerLogout = errors.New("docker logout error")
var ErrCommandFailure = errors.New("command execution failed")

func buildImage(ctx context.Context, references []*reference.Reference, isOfficial bool, isDryRun bool, isDebug bool, archOption string, opts driver.BuildOptions) (io.ReadCloser, error) {
	arguments, err := getBuildArgs(references, isOfficial, archOption, opts)
	if err != nil {
		return nil, err
	}
	cmd := getCommand(arguments)

	log.Infof("Building %s using %s/%s", references[0].Remote(), opts.Context, opts.Dockerfile)
	log.Debugf("command: %v", cmd)

	if isDryRun {
		return nil, nil
	}

	outputStream, err := term.StreamCommand(cmd, isDryRun, isDebug)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil, err
	}

	return outputStream, nil
}

func getBuildArgs(refs []*reference.Reference, isOfficial bool, archOption string, opts driver.BuildOptions) ([]string, error) {
	args := []string{"build"}

	for _, tag := range refs {
		args = append(args, "-t", tag.Remote())
	}

	args = append(args, "-f", fmt.Sprintf("%s/%s", opts.Context, opts.Dockerfile))

	for _, buildArg := range opts.BuildArgs {
		args = append(args, "--build-arg", buildArg)
	}

	if opts.NoCache {
		args = append(args, "--no-cache")
	}

	if opts.Pull {
		args = append(args, "--pull")
	}

	args = append(args, opts.Context)

	return args, nil
}

func createManifest(ctx context.Context, reference *reference.Reference, supportedArchitectures []string, isOfficial bool, archOption string, isDryRun bool, isDebug bool) error {
	log.Infof("Creating Manifest for %v", reference.Remote())

	arguments, err := getCreateArgs(reference, supportedArchitectures, isOfficial, archOption)
	if err != nil {
		return err
	}
	manifestCreateCmd := getCommand(arguments)

	if err := executeCommand(manifestCreateCmd, isDryRun, isDebug); err != nil {
		return err
	}

	return nil
}

func getCreateArgs(ref *reference.Reference, supportedArchitectures []string, isOfficial bool, archOption string) ([]string, error) {
	args := []string{"manifest", "create", ref.Remote()}

	for _, arch := range supportedArchitectures {
		// Generate the arch uri for the image
		uri, err := reference.NewUri(ref.Name(), &reference.UriOptions{
			Registry:   ref.Registry(),
			Official:   isOfficial,
			Arch:       arch,
			ArchOption: reference.ArchOption(archOption),
		})
		if err != nil {
			return nil, err
		}

		args = append(args, uri.Remote())
	}

	return args, nil
}

func annotateManifest(ctx context.Context, reference *reference.Reference, supportedArchitectures []string, isOfficial bool, archOption string, isDryRun bool, isDebug bool) error {
	log.Infof("Annotating Manifest for %v", reference.Remote())

	annotateCommands, err := getAnnotateCommands(reference, supportedArchitectures, isOfficial, archOption)
	if err != nil {
		return err
	}

	for _, cmd := range annotateCommands {
		if err := executeCommand(cmd, isDryRun, isDebug); err != nil {
			return err
		}
	}

	return nil
}

func getAnnotateCommands(ref *reference.Reference, supportedArchitectures []string, isOfficial bool, archOption string) ([]*term.Command, error) {
	var annotateCommands []*term.Command
	baseArgs := []string{"manifest", "annotate", ref.Remote()}

	for _, arch := range supportedArchitectures {
		args := []string{}
		// Generate the arch uri for the image
		uri, err := reference.NewUri(ref.Name(), &reference.UriOptions{
			Registry:   ref.Registry(),
			Official:   isOfficial,
			Arch:       arch,
			ArchOption: reference.ArchOption(archOption),
		})
		if err != nil {
			return nil, err
		}

		args = append(args, baseArgs...)
		args = append(args, uri.Remote(), "--arch", arch)

		cmd := getCommand(args)
		annotateCommands = append(annotateCommands, cmd)

	}

	return annotateCommands, nil
}

func pushManifest(ctx context.Context, reference *reference.Reference, isDryRun bool, isDebug bool, opts driver.ManifestPushOptions) error {
	log.Infof("Pushing Manifest for %v", reference.Remote())

	arguments, err := getPushArgs(reference, opts)
	if err != nil {
		return err
	}

	pushCmd := getCommand(arguments)

	if err := executeCommand(pushCmd, isDryRun, isDebug); err != nil {
		return err
	}

	return nil
}

func getPushArgs(reference *reference.Reference, opts driver.ManifestPushOptions) ([]string, error) {
	args := []string{"manifest", "push"}

	if opts.Purge {
		args = append(args, "--purge")
	}

	args = append(args, reference.Remote())

	return args, nil
}

func removeManifests(ctx context.Context, references []*reference.Reference, isDryRun bool, isDebug bool) error {
	allReferences := []string{}
	for _, ref := range references {
		allReferences = append(allReferences, ref.Remote())
	}

	log.Infof("Removing Manifests for %v", strings.Join(allReferences, " "))

	arguments, err := getRmArgs(references)
	if err != nil {
		return err
	}

	pushCmd := getCommand(arguments)

	if err := executeCommand(pushCmd, isDryRun, isDebug); err != nil {
		return err
	}

	return nil
}

func getRmArgs(references []*reference.Reference) ([]string, error) {
	args := []string{"manifest", "rm"}

	for _, ref := range references {
		args = append(args, ref.Remote())
	}

	return args, nil
}

// packageBuildContext creates a tarball from the context directory, excluding files marked in the dockerignore file
func packageBuildContext(context string, dockerfile string) (io.ReadCloser, error) {
	excludes, err := build.ReadDockerignore(context)
	if err != nil {
		return nil, err
	}

	if err := build.ValidateContextDirectory(context, excludes); err != nil {
		return nil, errors.Wrap(err, "error checking context")
	}

	excludes = build.TrimBuildFilesFromExcludes(excludes, dockerfile, false)
	buildContext, err := archive.TarWithOptions(context, &archive.TarOptions{
		ExcludePatterns: excludes,
		ChownOpts:       &idtools.Identity{UID: 0, GID: 0},
	})
	if err != nil {
		return nil, err
	}

	return buildContext, nil
}

func imageBuildOptions(buildUris []*reference.Reference, opts driver.BuildOptions) types.ImageBuildOptions {
	// Prepare the tags
	var buildTags []string
	for _, uri := range buildUris {
		buildTags = append(buildTags, uri.Remote())
	}

	// Prepare the build arguments
	buildArgs := make(map[string]*string)
	for _, pair := range opts.BuildArgs {
		// Split the pair on the equals sign
		kv := strings.Split(pair, "=")
		// Add the key-value pair to the map
		buildArgs[kv[0]] = &kv[1]
	}

	return types.ImageBuildOptions{
		Dockerfile: opts.Dockerfile,
		Tags:       buildTags,
		NoCache:    opts.NoCache,
		Remove:     opts.Rm,
		PullParent: opts.Pull,
		BuildArgs:  buildArgs,
	}
}

// Returns base64 encoded registry credentials
func encodeRegistryCredentials(registry *registry.Registry) (string, error) {
	authConfig := registrytypes.AuthConfig{
		Username:      registry.User.Name,
		Password:      registry.User.Password,
		ServerAddress: registry.ServerAddress,
	}
	authConfigAsBytes, err := json.Marshal(authConfig)
	if err != nil {
		return "", err
	}
	encodedAuthConfig := base64.URLEncoding.EncodeToString(authConfigAsBytes)
	return encodedAuthConfig, nil
}

// Helper functions
func getCommand(args []string) *term.Command {
	return &term.Command{
		Command: "docker",
		Args:    args,
	}
}

func executeCommand(cmd *term.Command, isDryRun bool, isDebug bool) error {

	if err := validateCommand(cmd); err != nil {
		return err
	}

	log.Debugf("Running command: %v %v", cmd.Command, cmd.Args)

	if isDryRun {
		return nil
	}

	output, err := exec.Command(cmd.Command, cmd.Args...).Output()
	if err != nil {
		if isDebug {
			log.Errorf("Failure: [%s %s] ", cmd.Command, cmd.Args)
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				if len(exitErr.Stderr) > 0 {
					log.Debugf("Exit code: %d\n****\n%s****\n", exitErr.ExitCode(), string(exitErr.Stderr))
				}
			}
			if len(output) > 0 {
				log.Debugf("STDOUT\n****\n%s****\n", string(output))
			}
		} else {
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				if len(exitErr.Stderr) > 0 {
					fmt.Printf("%v", string(exitErr.Stderr))
				}
			}
		}
		return ErrCommandFailure
	}
	return nil
}

// validateCommand ensures that the command only contains expected commands so unexpected programs are not run
func validateCommand(cmd *term.Command) error {
	allowedCommands := []string{"docker"}
	disallowedCharacters := []string{";", "|", "&"}

	if !slices.Contains(allowedCommands, cmd.Command) {
		return fmt.Errorf("disallowed command: %s", cmd.Command)
	}

	input := strings.Join(cmd.Args, " ")
	for _, char := range disallowedCharacters {
		if strings.Contains(input, char) {
			return fmt.Errorf("disallowed character: %s", char)
		}
	}

	return nil
}
