package service

import (
	"github.com/spf13/cobra"

	"github.com/calvinverse/service.provisioning.controller/internal/cmd"
	"github.com/calvinverse/service.provisioning.controller/internal/config"
	"github.com/calvinverse/service.provisioning.controller/internal/doc"
	"github.com/calvinverse/service.provisioning.controller/internal/health"
	"github.com/calvinverse/service.provisioning.controller/internal/router"
	"github.com/calvinverse/service.provisioning.controller/internal/web"
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

func (r *resolver) resolveAPIRouters() []router.APIRouter {
	docRouter := doc.NewDocumentationRouter(r.cfg)
	healthRouter := health.NewHealthAPIRouter()
	return []router.APIRouter{
		docRouter,
		healthRouter,
	}
}

func (r *resolver) ResolveCommands() []*cobra.Command {
	routerBuilder := r.resolveRouterBuilder()
	ServeCommandBuilder := cmd.NewServeCommandBuilder(r.cfg, routerBuilder)
	BootstrapCommandBuilder := cmd.NewBootstrapCommandBuilder(r.cfg, routerBuilder)

	if r.commands == nil {
		r.commands = []*cobra.Command{
			ServeCommandBuilder.New(),
			BootstrapCommandBuilder.New(),
		}
	}

	return r.commands
}

func (r *resolver) resolveRouterBuilder() router.Builder {
	apiRouters := r.resolveAPIRouters()
	webRouter := r.resolveWebRouter()
	return router.NewRouterBuilder(apiRouters, webRouter)
}

func (r *resolver) resolveWebRouter() router.WebRouter {
	return web.NewWebRouter(r.cfg)
}
