package registry

import (
	consul "github.com/hashicorp/consul/api"
)

// NewConsulDiscoveryClient creates an object which connects to the Consul service discovery agent
func NewConsulDiscoveryClient() (Registry, error) {

	// Want to do this async? Because this could be delayed if Consul isn't
	// running. We don't necessarily want to wait for it, but we should mark
	// ourselves as not-healthy if we're not registered with consul(?)
	c, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		return nil, err
	}

	client := &consulDiscoveryClient{
		consulAgent: c.Agent(),
	}

	return client, nil
}

type consulDiscoveryClient struct {
	consulAgent *consul.Agent
}

func (cd consulDiscoveryClient) Register() error {}

func (cd consulDiscoveryClient) Deregister() error {}

// health check
