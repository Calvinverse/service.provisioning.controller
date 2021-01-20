package registry

import (
	"context"
	"crypto/tls"
	"time"
)

// Registry describes the methods to interact with a service discovery system.
type Registry interface {
	Init(...Option) error
	Options() Options
	Register(...RegisterOption) error
	Deregister(...DeregisterOption) error
	GetService(string, ...GetOption) ([]*Service, error)
	ListServices(...ListOption) ([]*Service, error)
}

// Service describes a service as registered with a service discovery system.
type Service struct {
	Name      string            `json:"name"`
	Version   string            `json:"version"`
	Metadata  map[string]string `json:"metadata"`
	Endpoints []*Endpoint       `json:"endpoints"`
	Nodes     []*Node           `json:"nodes"`
}

// Node describes a compute node that executes one of the services registered with the service discovery system.
type Node struct {
	ID       string            `json:"id"`
	Address  string            `json:"address"`
	Metadata map[string]string `json:"metadata"`
}

// Endpoint describes an endpoint for a registered service.
type Endpoint struct {
	Name     string            `json:"name"`
	Request  *Value            `json:"request"`
	Response *Value            `json:"response"`
	Metadata map[string]string `json:"metadata"`
}

// Value describes the value of a configuration option.
type Value struct {
	Name   string   `json:"name"`
	Type   string   `json:"type"`
	Values []*Value `json:"values"`
}

// Options describes the configuration options for a service discovery system.
type Options struct {
	Addrs     []string
	Timeout   time.Duration
	Secure    bool
	TLSConfig *tls.Config
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// RegisterOptions describes the options for the registration of a service.
type RegisterOptions struct {
	TTL time.Duration
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// DeregisterOptions describes the options for the de-registration of a service.
type DeregisterOptions struct {
	Context context.Context
}

// GetOptions describes the options for getting information about a service.
type GetOptions struct {
	Context context.Context
}

// ListOptions describes the options for listing a number of services.
type ListOptions struct {
	Context context.Context
}

// Option describes the function that handles options.
type Option func(*Options)

// RegisterOption describes the function that handles registration options.
type RegisterOption func(*RegisterOptions)

// DeregisterOption describes the function that handles de-registration options.
type DeregisterOption func(*DeregisterOptions)

// GetOption describes the function that handles get options.
type GetOption func(*GetOptions)

// ListOption describes the function that handles list options.
type ListOption func(*ListOptions)
