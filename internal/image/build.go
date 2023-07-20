package image

import (
	"context"
	"tugboat/internal/driver"
	"tugboat/internal/term"
)

type BuildOptions struct {
	Context    string
	Dockerfile string
	Tags       []string
	BuildArgs  []string
	Push       bool
	Pull       bool
	NoCache    bool
}

func Build(ctx context.Context, d driver.ImageBuilderPusher, opts BuildOptions) error {
	buildOpts := driver.BuildOptions{
		Context:    opts.Context,
		Dockerfile: opts.Dockerfile,
		Tags:       opts.Tags,
		BuildArgs:  opts.BuildArgs,
		Rm:         true,
		Pull:       opts.Pull,
		NoCache:    opts.NoCache,
		Push:       opts.Push,
	}
	output, err := d.BuildImage(ctx, buildOpts)
	if err != nil {
		return err
	}

	if output != nil {
		defer output.Close()

		if err := term.DisplayResponse(output); err != nil {
			return err
		}
	}

	if opts.Push {
		for _, tag := range opts.Tags {
			pushOutput, err := d.PushImage(ctx, tag)
			if err != nil {
				return err
			}

			if pushOutput != nil {
				defer pushOutput.Close()

				if err := term.DisplayResponse(pushOutput); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
