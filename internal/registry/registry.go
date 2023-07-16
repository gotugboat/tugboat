package registry

import "errors"

type Registry struct {
	ServerAddress string
	Namespace     string
	User          *RegistryUser
}

type RegistryUser struct {
	Name     string
	Password string
}

func NewRegistry(serverAddress, namespace, username, password string) (*Registry, error) {
	if serverAddress == "" || username == "" || password == "" {
		return nil, errors.New("invalid input parameters")
	}

	return &Registry{
		ServerAddress: serverAddress,
		Namespace:     namespace,
		User: &RegistryUser{
			Name:     username,
			Password: password,
		},
	}, nil
}
