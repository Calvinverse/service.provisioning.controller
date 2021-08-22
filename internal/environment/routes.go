package environment

import (
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/calvinverse/service.provisioning.controller/internal/config"
	"github.com/calvinverse/service.provisioning.controller/internal/db/repository"
	"github.com/calvinverse/service.provisioning.controller/internal/router"
	"github.com/go-chi/chi"
)

// NewEnvironmentAPIRouter returns an APIRouter instance for the environment routes.
func NewEnvironmentAPIRouter(config config.Configuration, storage *repository.Storage) router.APIRouter {
	return &environmentRouter{
		cfg:     config,
		storage: storage,
	}
}

// https://www.google.com/search?client=firefox-b-d&q=golang+chi+get+query+params
// https://github.com/pressly/imgry/blob/master/server/server.go
// https://github.com/pressly/imgry/blob/master/server/middleware.go
// https://github.com/pressly/imgry/blob/master/server/handlers.go
// https://github.com/pressly/imgry/blob/bbb40ff8100ff84b8290005ebe080b7b07939372/server/middleware.go

type environmentRouter struct {
	cfg     config.Configuration
	storage *repository.Storage
}

// Environment is used to store information about an environment as used by the REST API
type Environment struct {
}

type EnvironmentRequest struct {
	Callback    string      `json:"callback"`
	Environment Environment `json:"environment"`
}

// CreateEnvironment godoc
// @Summary Creates a new environment.
// @Description Creates a new environment based on the provided information.
// @Tags environment
// @Accept  json
// @Produce  json
// @Param id body environment.Environment true "Environment ID"
// @Success 201 {object} environment.Environment
// @Failure 404 {object} int
// @Failure 500 {object} int
// @Router /v1/environment [put]
func (h *environmentRouter) create(w http.ResponseWriter, r *http.Request) {
	//render.Status()
	// Receive a json file and queue the command

	currentTime := time.Now()
	environment := &repository.Environment{
		ID:                   "test-ID",
		Name:                 "test-name",
		Description:          "A test environment",
		CreatedOn:            &currentTime,
		DestroyedOn:          &time.Time{},
		DestructionPlannedOn: &time.Time{},
	}

	err := environment.Store(r.Context(), h.storage)
	if err != nil {
		log.
			WithError(err).
			Error("Failed to store the environment data")

		if _, ok := err.(*repository.DuplicateEnvironmentError); ok {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(`{"message": "An environment with the given ID already exists."}`))
			return
		} else {

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "Failed to store the environment data"}`))
			return
		}
	}

	json.NewEncoder(w).Encode(environment)
}

// DeleteEnvironment godoc
// @Summary Deletes an environment.
// @Description Deletes the environment with the given id.
// @Tags environment
// @Accept  json
// @Produce  json
// @Param id path string true "Environment ID"
// @Success 202 {object} environment.Environment
// @Failure 404 {object} int
// @Failure 500 {object} int
// @Router /v1/environment/{id} [delete]
func (h *environmentRouter) delete(w http.ResponseWriter, r *http.Request) {

}

// ShowEnvironment godoc
// @Summary Provide information about an environment.
// @Description Returns information about the environment with the given id.
// @Tags environment
// @Accept  json
// @Produce  json
// @Param id path string true "Environment ID"
// @Success 200 {object} environment.Environment
// @Failure 404 {object} int
// @Failure 500 {object} int
// @Router /v1/environment/{id} [get]
func (h *environmentRouter) get(w http.ResponseWriter, r *http.Request) {

}

// ListEnvironmentIDs godoc
// @Summary Provide the list of known environment IDs
// @Description Returns a list of known environment IDs.
// @Tags environment
// @Accept  json
// @Produce  json
// @Success 200 {array} string
// @Failure 404 {object} int
// @Failure 500 {object} int
// @Router /v1/environment/ [get]
func (h *environmentRouter) list(w http.ResponseWriter, r *http.Request) {

}

func (h *environmentRouter) update(w http.ResponseWriter, r *http.Request) {

}

func (h *environmentRouter) Prefix() string {
	return "environment"
}

// Routes creates the routes for the health package
func (h *environmentRouter) Routes(prefix string, r chi.Router) {
	r.Route(prefix, func(r chi.Router) {
		r.Get("/", h.list)
		r.Get("/{id}", h.get)
		r.Delete("/{id}", h.delete)
		r.Put("/", h.create)
	})
}

func (h *environmentRouter) Version() int8 {
	return 1
}
