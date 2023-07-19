package image

import "errors"

var (
	ErrNoProvidedTags           = errors.New("tags must be provided")
	ErrNoSupportedArchitectures = errors.New("there are no supported architectures define")
)
