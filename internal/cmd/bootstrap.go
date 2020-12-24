package cmd

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"

	"github.com/calvinverse/service.provisioning.controller/internal/config"
	"github.com/calvinverse/service.provisioning.controller/internal/router"
)

// BootstrapCommandBuilder creates new Cobra Commands for the bootstrap capability.
type BootstrapCommandBuilder interface {
	New() *cobra.Command
}

// NewBootstrapCommandBuilder creates a new instance of the BootstrapCommandBuilder interface.
func NewBootstrapCommandBuilder(config config.Configuration, builder router.Builder) ServeCommandBuilder {
	return &bootstrapCommandBuilder{
		cfg:     config,
		builder: builder,
	}
}

type bootstrapCommandBuilder struct {
	cfg     config.Configuration
	builder router.Builder
}

func (s bootstrapCommandBuilder) New() *cobra.Command {
	return &cobra.Command{
		Use:   "bootstrap",
		Short: "Bootstraps the environment that the service is supposed to execute in",
		Long:  "Bootstraps the environment that the service is supposed to execute in. Assumes no other environment of this kind exists.",
		RunE:  s.executeServer,
	}
}

func (s bootstrapCommandBuilder) executeServer(cmd *cobra.Command, args []string) error {
	log.Printf("Pretending to bootstrap \n")

	return nil
}
