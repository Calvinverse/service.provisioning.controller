package doc

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/calvinverse/service.provisioning.controller/internal/config"
	"github.com/calvinverse/service.provisioning.controller/internal/router"
	"github.com/go-chi/chi"

	httpSwagger "github.com/swaggo/http-swagger"

	log "github.com/sirupsen/logrus"
)

// NewDocumentationRouter returns an APIRouter instance for the documentation routes.
func NewDocumentationRouter(config config.Configuration) router.APIRouter {
	return &docRouter{
		cfg: config,
	}
}

type docRouter struct {
	cfg config.Configuration
}

func (d *docRouter) Prefix() string {
	return "doc"
}

// GetOpenApiJson godoc
// @Summary Returns the OpenAPI document for the current service
// @Description Returns the OpenAPI document for the current service
// @Tags doc
// @Accept  json
// @Produce  json
// @Success 200 {object} environment.Environment
// @Router /v1/doc [get]
func (d *docRouter) Routes(prefix string, r chi.Router) {
	filesDir := ""
	if d.cfg.IsSet("doc.path") {
		filesDir = d.cfg.GetString("doc.path")
	} else {
		ex, err := os.Executable()
		if err != nil {
			log.Fatal(
				fmt.Sprintf(
					"Failed to locate the directory of the executable. Error was: %v",
					err))
		}

		workDir := filepath.Dir(ex)
		filesDir = filepath.Join(workDir, "api")
	}

	log.Debug(
		fmt.Sprintf(
			"Using doc directory %s",
			filesDir))

	r.Get("/swagger", httpSwagger.WrapHandler)
	r.Get(prefix, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filesDir+"/swagger.json")
	})
}

func (d *docRouter) Version() int8 {
	return 1
}
