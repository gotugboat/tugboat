package image

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"tugboat/internal/pkg/reference"
	"tugboat/internal/pkg/slices"

	log "github.com/sirupsen/logrus"
)

var ErrDockerLogin = errors.New("docker login error")
var ErrDockerLogout = errors.New("docker logout error")
var ErrCommandFailure = errors.New("command execution failed")

type DockerLoginOptions struct {
	ServerAddress string
	Username      string
	Password      string
	DryRun        bool
}

type PushManifestOptions struct {
	Purge  bool
	DryRun bool
	Debug  bool
}

type RmManifestOptions struct {
	DryRun bool
	Debug  bool
}

func dockerLogin(ctx context.Context, opts *DockerLoginOptions) error {
	log.Infof("Logging into %v as %v", opts.ServerAddress, opts.Username)

	if opts.DryRun {
		return nil
	}

	loginCmd := []string{"login", "--username", opts.Username, "--password", opts.Password, opts.ServerAddress}

	output, err := exec.Command("docker", loginCmd...).Output()
	if err != nil {
		log.Errorf("Docker login failed: %v", err)
		return ErrDockerLogin
	}
	log.Info(string(output))

	return nil
}

func dockerLogout(ctx context.Context, registry string, isDryRun bool) error {
	log.Infof("Logging out of %v", registry)

	if isDryRun {
		return nil
	}

	loginOutCmd := []string{"logout", registry}

	output, err := exec.Command("docker", loginOutCmd...).Output()
	if err != nil {
		return ErrDockerLogout
	}
	log.Info(string(output))

	return nil
}

func createManifest(ctx context.Context, reference *reference.Reference, opts ManifestCreateOptions) error {
	log.Infof("Creating Manifest for %v", reference.Remote())

	arguments, err := getCreateArgs(reference, opts)
	if err != nil {
		return err
	}
	manifestCreateCmd := getCommand(arguments)

	if err := executeCommand(manifestCreateCmd, opts.DryRun, opts.Debug); err != nil {
		return err
	}

	return nil
}

func annotateManifest(ctx context.Context, reference *reference.Reference, opts ManifestCreateOptions) error {
	log.Infof("Annotating Manifest for %v", reference.Remote())

	annotateCommands, err := getAnnotateCommands(reference, opts)
	if err != nil {
		return err
	}

	for _, cmd := range annotateCommands {
		if err := executeCommand(cmd, opts.DryRun, opts.Debug); err != nil {
			return err
		}
	}

	return nil
}

func pushManifest(ctx context.Context, reference *reference.Reference, opts PushManifestOptions) error {
	log.Infof("Pushing Manifest for %v", reference.Remote())

	arguments, err := getPushArgs(reference, opts)
	if err != nil {
		return err
	}
	pushCmd := getCommand(arguments)

	if err := executeCommand(pushCmd, opts.DryRun, opts.Debug); err != nil {
		return err
	}

	return nil
}

func removeManifests(ctx context.Context, references []*reference.Reference, opts RmManifestOptions) error {
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

	if err := executeCommand(pushCmd, opts.DryRun, opts.Debug); err != nil {
		return err
	}
	return nil
}

// Helper functions
type Command struct {
	Command string
	Args    []string
}

func getCommand(args []string) *Command {
	return &Command{
		Command: "docker",
		Args:    args,
	}
}

func executeCommand(cmd *Command, isDryRun bool, isDebug bool) error {

	if err := validateCommand(cmd); err != nil {
		return err
	}

	if isDebug {
		log.Debugf("Running command: %v %v", cmd.Command, cmd.Args)
	}

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

func getCreateArgs(ref *reference.Reference, opts ManifestCreateOptions) ([]string, error) {
	args := []string{"manifest", "create", ref.Remote()}

	for _, arch := range opts.SupportedArchitectures {
		// Generate the arch uri for the image
		uri, err := reference.NewUri(ref.Name(), &reference.UriOptions{
			Registry:   ref.Registry(),
			Official:   opts.Official,
			Arch:       arch,
			ArchOption: toArchOption(opts.ArchOption),
		})
		if err != nil {
			return nil, err
		}

		args = append(args, uri.Remote())
	}

	return args, nil

}

func getPushArgs(reference *reference.Reference, opts PushManifestOptions) ([]string, error) {
	args := []string{"manifest", "push"}

	if opts.Purge {
		args = append(args, "--purge")
	}

	args = append(args, reference.Remote())

	return args, nil

}

func getAnnotateCommands(ref *reference.Reference, opts ManifestCreateOptions) ([]*Command, error) {
	var annotateCommands []*Command
	baseArgs := []string{"manifest", "annotate", ref.Remote()}

	for _, arch := range opts.SupportedArchitectures {
		args := []string{}
		// Generate the arch uri for the image
		uri, err := reference.NewUri(ref.Name(), &reference.UriOptions{
			Registry:   ref.Registry(),
			Official:   opts.Official,
			Arch:       arch,
			ArchOption: toArchOption(opts.ArchOption),
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

func getRmArgs(references []*reference.Reference) ([]string, error) {
	args := []string{"manifest", "rm"}

	for _, ref := range references {
		args = append(args, ref.Remote())
	}

	return args, nil
}

// validateCommand ensures that the command only contains expected commands so unexpected programs are not run
func validateCommand(cmd *Command) error {
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
