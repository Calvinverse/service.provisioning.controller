package service

import (
	"github.com/spf13/cobra"

	"github.com/calvinverse/service.provisioning.controller/internal/cmd"
	"github.com/calvinverse/service.provisioning.controller/internal/config"
)

// Resolver defines the interface for Inversion-of-Control objects.
type Resolver interface {
	ResolveCommands() []*cobra.Command
}

// NewResolver returns a new Resolver instance
func NewResolver(config config.Configuration) Resolver {
	return &resolver{
		cfg: config,
	}
}

type resolver struct {
	cfg config.Configuration

	commands []*cobra.Command
}

func (r *resolver) ResolveCommands() []*cobra.Command {
	ServeCommandBuilder := cmd.NewServeCommandBuilder(r.cfg)

	if r.commands == nil {
		r.commands = []*cobra.Command{
			ServeCommandBuilder.New(),
		}
	}

	return r.commands
}
