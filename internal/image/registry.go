package image

type Registry struct {
	ServerAddress string
	Namespace     string
	User          RegistryUser
}

type RegistryUser struct {
	Name     string
	Password string
}

type RegistryAuth struct {
	Username      string
	Password      string
	Namespace     string
	ServerAddress string
}

func NewRegistry(serverAddress, namespace, username, password string) Registry {
	return Registry{
		ServerAddress: serverAddress,
		Namespace:     namespace,
		User: RegistryUser{
			Name:     username,
			Password: password,
		},
	}
}
