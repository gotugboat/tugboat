package driver

import (
	"fmt"
	"strings"
	"tugboat/internal/pkg/reference"

	"github.com/pkg/errors"
)

func GenerateUri(registry string, namespace string, tag string, official bool, archOption reference.ArchOption) (*reference.Reference, error) {
	image := fmt.Sprintf("%s/%s", namespace, tag)
	s := strings.Split(tag, "/")
	if len(s) == 3 {
		image = tag
	}
	uri, err := reference.NewUri(image, &reference.UriOptions{
		Registry:   registry,
		Official:   official,
		ArchOption: archOption,
	})
	if err != nil {
		return nil, errors.Errorf("%v", err)
	}
	return uri, nil
}

func GenerateAllUris(registry string, namespace string, tags []string, official bool, archOption reference.ArchOption) ([]*reference.Reference, error) {
	buildTags := []*reference.Reference{}

	for _, tag := range tags {
		taggedUri, err := GenerateUri(registry, namespace, tag, official, archOption)
		if err != nil {
			return nil, errors.Errorf("%v", err)
		}
		buildTags = append(buildTags, taggedUri)
	}

	return buildTags, nil
}
