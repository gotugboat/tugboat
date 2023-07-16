package image

import (
	"encoding/base64"
	"encoding/json"
	"tugboat/internal/pkg/reference"

	"github.com/docker/docker/api/types"
)

// Returns encoded registry credentials
func encodeRegistryCredentials(registry Registry) (string, error) {
	authConfig := types.AuthConfig{
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

// Return a given string as a reference.ArchOption
func toArchOption(value string) reference.ArchOption {
	switch value {
	case "prepend":
		return reference.ArchPrepend
	case "append":
		return reference.ArchAppend
	default:
		return reference.ArchOmit
	}
}
