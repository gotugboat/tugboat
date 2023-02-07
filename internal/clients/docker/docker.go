package docker

import (
	"github.com/docker/docker/client"
)

func NewClientWithOpts(opts ...client.Opt) (*client.Client, error) {
	client, err := client.NewClientWithOpts(opts...)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewClientFromEnv() (*client.Client, error) {
	client, err := NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	return client, nil
}
