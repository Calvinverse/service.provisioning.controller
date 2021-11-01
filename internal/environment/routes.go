package environment

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/calvinverse/service.provisioning.controller/internal/config"
	"github.com/calvinverse/service.provisioning.controller/internal/db/repository"
	"github.com/calvinverse/service.provisioning.controller/internal/router"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/golang/gddo/httputil/header"
)

type malformedRequest struct {
	status int
	msg    string
}

func (mr *malformedRequest) Error() string {
	return mr.msg
}

type invalidEnvironmentField struct {
	field string
}

func (ie *invalidEnvironmentField) Error() string {
	return fmt.Sprintf("The environment data is missing information for field %s", ie.field)
}

// Environment is used to store information about an environment as used by the REST API
// This is a different type than the repository.Environment so that the data that is send
// over the REST API isn't directly linked to the data that is stored in the database. That
// way we can independently evolve both parts.
type Environment struct {
	// ID provides the unique identity of the environment.
	ID string `json:"_key"`

	// Name is the human readable name of the environment, e.g. 'production'.
	Name string `json:"name"`

	// Description is the human readable description of the environment.
	Description string `json:"description"`

	// CreatedOn is the date and time the environment was created. May be set to nil
	// if the environment has not been created yet.
	CreatedOn *time.Time `json:"created_on"`

	// DestroyedON is the date and time the environment was destroyed. May be set to
	// nil if the environment has not been destroyed yet.
	DestroyedOn *time.Time `json:"destroyed_on"`

	// DestructionPlannedOn is the date and time the environment is planned to be
	// destroyed. If the environment is not planned to be destroyed then it will
	// be set to nil.
	DestructionPlannedOn *time.Time `json:"destruction_planned_on"`
}

func EnvironmentFromStorage(environment *repository.Environment) (*Environment, error) {
	result := &Environment{
		ID:                   environment.ID,
		Name:                 environment.Name,
		Description:          environment.Description,
		CreatedOn:            environment.CreatedOn,
		DestroyedOn:          environment.DestroyedOn,
		DestructionPlannedOn: environment.DestructionPlannedOn,
	}

	return result, nil
}

func EnvironmentToStorage(environment *Environment) (*repository.Environment, error) {
	if environment.ID == "" {
		return nil, &invalidEnvironmentField{field: "ID"}
	}

	if environment.Name == "" {
		return nil, &invalidEnvironmentField{field: "Name"}
	}

	result := &repository.Environment{
		ID:                   environment.ID,
		Name:                 environment.Name,
		Description:          environment.Description,
		CreatedOn:            environment.CreatedOn,
		DestroyedOn:          environment.DestroyedOn,
		DestructionPlannedOn: environment.DestructionPlannedOn,
	}

	return result, nil
}

type EnvironmentRequest struct {
	Callback     string        `json:"callback"`
	Environments []Environment `json:"environments"`

	// Cursor?
}

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

// CreateEnvironment godoc
// @Summary Creates a new environment.
// @Description Creates a new environment based on the provided information.
// @Tags environment
// @Accept  json
// @Accept  xml
// @Produce  json
// @Produce  xml
// @Param id body environment.Environment true "Environment ID"
// @Success 201 {object} environment.Environment
// @Failure 404 {object} int
// @Failure 409 {object} int
// @Failure 500 {object} int
// @Router /v1/environment [put]
func (h *environmentRouter) create(w http.ResponseWriter, r *http.Request) {
	environment := &Environment{}
	err := h.decodeRequestBody(w, r, environment)
	if err != nil {
		if err1, ok := err.(*malformedRequest); ok {
			h.responseBody(w, r, err1.status, err1.msg)
			return
		} else {
			h.responseBody(w, r, http.StatusInternalServerError, "Failed to read the request body.")
			return
		}
	}

	storedEnv, err := EnvironmentToStorage(environment)
	if err != nil {
		h.responseBody(w, r, http.StatusBadRequest, "The environment data was missing required fields.")
		return
	}

	err = repository.StoreEnvironment(r.Context(), h.storage, storedEnv)
	if err != nil {
		log.
			WithError(err).
			Error("Failed to store the environment data")

		if _, ok := err.(*repository.DuplicateEnvironmentError); ok {
			h.responseBody(w, r, http.StatusConflict, "An environment with the given ID already exists.")
			return
		} else {
			h.responseBody(w, r, http.StatusInternalServerError, "Failed to store the environment data.")
			return
		}
	}

	h.responseBody(w, r, http.StatusCreated, environment)
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
	id := chi.URLParam(r, "id")

	environment, err := repository.FetchEnvironmentByID(r.Context(), h.storage, id)
	if err != nil {
		log.
			WithError(err).
			Error("Failed to store the environment data")

		if _, ok := err.(*repository.UnknownEnvironmentError); ok {
			h.responseBody(w, r, http.StatusNotFound, "An environment with the given ID could not be found.")
			return
		} else {
			h.responseBody(w, r, http.StatusInternalServerError, "Failed to retrieve the environment data.")
			return
		}
	}

	apiEnv, err := EnvironmentFromStorage(environment)
	if err != nil {
		h.responseBody(w, r, http.StatusInternalServerError, "Failed to retrieve the environment data.")
		return
	}

	h.responseBody(w, r, http.StatusOK, apiEnv)
}

// ListEnvironmentIDs godoc
// @Summary Provide the list of known environment IDs
// @Description Returns a list of known environment IDs.
// @Tags environment
// @Accept  json
// @Produce  json
// @Success 200 {array} environment.Environment
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

func (h *environmentRouter) decodeRequestBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	// If the Content-Type header is present, check that it has the value
	// application/json. Note that we are using the gddo/httputil/header
	// package to parse and extract the value here, so the check works
	// even if the client includes additional charset or boundary
	// information in the header.
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value == "application/json" {
			return h.decodeRequestBodyJson(w, r, dst)
		}

		if value == "application/xml" {
			return h.decodeRequestBodyXml(w, r, dst)
		}
	}

	msg := "Content-Type header is not application/json"
	return &malformedRequest{status: http.StatusUnsupportedMediaType, msg: msg}
}

func (h *environmentRouter) decodeRequestBodyJson(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	// Use http.MaxBytesReader to enforce a maximum read of 1MB from the
	// response body. A request body larger than that will now result in
	// Decode() returning a "http: request body too large" error.
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	// Setup the decoder and call the DisallowUnknownFields() method on it.
	// This will cause Decode() to return a "json: unknown field ..." error
	// if it encounters any extra unexpected fields in the JSON. Strictly
	// speaking, it returns an error for "keys which do not match any
	// non-ignored, exported fields in the destination".
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		// Catch any syntax errors in the JSON and send an error message
		// which interpolates the location of the problem to make it
		// easier for the client to fix.
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		// In some circumstances Decode() may also return an
		// io.ErrUnexpectedEOF error for syntax errors in the JSON. There
		// is an open issue regarding this at
		// https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		// Catch any type errors, like trying to assign a string in the
		// JSON request body to a int field in our Person struct. We can
		// interpolate the relevant field name and position into the error
		// message to make it easier for the client to fix.
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		// Catch the error caused by extra unexpected fields in the request
		// body. We extract the field name from the error message and
		// interpolate it in our custom error message. There is an open
		// issue at https://github.com/golang/go/issues/29035 regarding
		// turning this into a sentinel error.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		// An io.EOF error is returned by Decode() if the request body is
		// empty.
		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		// Catch the error caused by the request body being too large. Again
		// there is an open issue regarding turning this into a sentinel
		// error at https://github.com/golang/go/issues/30715.
		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &malformedRequest{status: http.StatusRequestEntityTooLarge, msg: msg}

		// Otherwise default to logging the error and sending a 500 Internal
		// Server Error response.
		default:
			return err
		}
	}

	// Call decode again, using a pointer to an empty anonymous struct as
	// the destination. If the request body only contained a single JSON
	// object this will return an io.EOF error. So if we get anything else,
	// we know that there is additional data in the request body.
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		return &malformedRequest{status: http.StatusBadRequest, msg: msg}
	}

	return nil
}

func (h *environmentRouter) decodeRequestBodyXml(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	return nil
}

func (h *environmentRouter) responseBody(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	mediatype, _, err := mime.ParseMediaType(r.Header.Get("Accept"))
	if err != nil {
		log.WithError(err).Error("Invalid 'Accept' header.")

		http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	switch mediatype {
	case "application/xml":
		render.Status(r, status)
		render.XML(w, r, data)
		return
	case "application/json":
		render.Status(r, status)
		render.JSON(w, r, data)
		return
	default:
		log.WithFields(log.Fields{
			"media_type": mediatype,
		}).Error("Invalid media type. Expected either 'application/json' or 'application/xml'.")

		http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}
}
