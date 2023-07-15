package reference

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"
)

// Returns a string that is formatted for official or unofficial images.
func generateUriString(image string, registry string, arch string, isOfficial bool) string {
	var uriString string

	// clean the provided uri
	uriString = trimUri(image)

	// find the namespace
	var namespace string
	if strings.Contains(uriString, "/") {
		// This is a shortname definition of the image name, we need to remove the namespace from imageName
		s := strings.Split(uriString, "/")
		if len(s) == 2 {
			namespace = s[0]
			image = s[1]
		} else if len(s) == 3 {
			registry = s[0]
			namespace = s[1]
			image = s[2]
		}
	}

	// build image so it can be compatible with Docker Hub deployment i.e. registry/namespace/image:arch-tag
	uriString = trimUri(fmt.Sprintf("%s/%s/%s", registry, namespace, image))

	// build image so that it mimics the official build format i.e. registry/arch/image:tag
	if isOfficial {
		uriString = trimUri(fmt.Sprintf("%s/%s/%s", registry, arch, image))
	}

	return uriString
}

// strip off all leading '/' characters, change `//` to `/`
func trimUri(uri string) string {
	re := regexp.MustCompile(`/+`)
	result := re.ReplaceAllString(uri, "/")
	result = strings.TrimLeft(result, "/")
	return result
}

// Returns the platform architecture
func getArch() string {
	return runtime.GOARCH
}
