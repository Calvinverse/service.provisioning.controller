package health

import (
	"fmt"
	"mime"
	"net/http"
	"time"

	"github.com/calvinverse/service.provisioning.controller/internal/info"
	"github.com/calvinverse/service.provisioning.controller/internal/router"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	log "github.com/sirupsen/logrus"
)

// InfoResponse stores the response to an info request
type InfoResponse struct {
	BuildTime string `json:"buildtime"`
	Revision  string `json:"revision"`
	Version   string `json:"version"`
}

// LivelinessDetailedResponse stores detailed information about the liveliness of the application, indicating if the application is healthy
type LivelinessDetailedResponse struct {
	// Status of all the health checks
	Checks []CheckResult `json:"checks"` // <-- this is wrong. We're using an internal type externally

	foobar()

	// Global status
	Status string `json:"status"`

	// Time the liveliness response was created at
	Time string `json:"time"`
}

// LivelinessSummaryResponse stores condensed information about the liveliness of the application, indicating if the application is healthy
type LivelinessSummaryResponse struct {
	// Status of all health checks
	Checks map[string]string `json:"checks"`

	// Global status
	Status string `json:"status"`

	// Time the liveliness response was created at
	Time string `json:"time"`
}

// PingResponse stores the response to a Ping request
type PingResponse struct {
	Response string `json:"response"`
}

// ReadinessResponse stores information about the readiness of the application, indicating whether the application is ready to serve responses.
type ReadinessResponse struct {
	Time string `json:"time"`
}

// StartedResponse stores information about the starting state of the application, indicating whether the application has started successfully.
type StartedResponse struct {
	Time string `json:"time"`
}

// NewSelfAPIRouter returns an APIRouter instance for the health routes.
func NewSelfAPIRouter() router.APIRouter {
	return &selfRouter{
		healthService: GetServiceWithDefaultSettings(),
	}
}

// selfRouter defines an APIRouter that routes the 'self' metadata routes.
type selfRouter struct {
	healthService Service
}

func (h *selfRouter) Prefix() string {
	return "self"
}

// Routes creates the routes for the health package
func (h *selfRouter) Routes(prefix string, r chi.Router) {
	r.Get(fmt.Sprintf("%s/info", prefix), h.info)
	r.Get(fmt.Sprintf("%s/liveliness", prefix), h.liveliness)
	r.Get(fmt.Sprintf("%s/ping", prefix), h.ping)
	r.Get(fmt.Sprintf("%s/readiness", prefix), h.readiness)
	r.Get(fmt.Sprintf("%s/started", prefix), h.started)
}

func (h *selfRouter) Version() int8 {
	return 1
}

// Info godoc
// @Summary Respond to an info request
// @Description Respond to an info request with information about the application.
// @Tags health
// @Accept json
// @Accept xml
// @Produce json
// @Produce xml
// @Success 200 {object} health.InfoResponse
// @Failure 415 {string} string "Unsupported media type"
// @Router /v1/self/info [get]
func (h *selfRouter) info(w http.ResponseWriter, r *http.Request) {
	response := InfoResponse{
		BuildTime: info.BuildTime(),
		Revision:  info.Revision(),
		Version:   info.Version(),
	}

	h.responseBody(w, r, http.StatusOK, response)
}

// Liveliness godoc
// @Summary Respond to an liveliness request
// @Description Respond to an liveliness request with information about the status of the latest health checks.
// @Tags health
// @Accept json
// @Accept xml
// @Produce json
// @Produce xml
// @Param type query string false "options are summary or detailed" Enums(summary, detailed)
// @Success 200 {object} health.LivelinessDetailedResponse
// @Failure 415 {string} string "Unsupported media type"
// @Router /v1/self/liveliness [get]
func (h *selfRouter) liveliness(w http.ResponseWriter, r *http.Request) {
	healthStatus, err := h.healthService.Liveliness()
	if err != nil {
		t := time.Now()
		response := &LivelinessSummaryResponse{
			Checks: make(map[string]string),
			Status: Failed,
			Time:   t.Format("Mon Jan _2 15:04:05 2006"),
		}

		h.responseBody(w, r, http.StatusInternalServerError, response)
		return
	}

	responseType := r.URL.Query().Get("type")
	switch responseType {
	case "detailed":
		h.livelinessDetailedResponse(w, r, &healthStatus)
	case "summary":
		fallthrough
	default:
		h.livelinessSummaryResponse(w, r, &healthStatus)
	}
}

func (h *selfRouter) livelinessDetailedResponse(w http.ResponseWriter, r *http.Request, status *Status) {
	t := time.Now()

	statusText := Success
	responseCode := http.StatusOK
	if !status.IsHealthy {
		statusText = Failed
		responseCode = http.StatusInternalServerError
	}

	checkResults := status.Checks

	response := &LivelinessDetailedResponse{
		Checks: checkResults,
		Status: statusText,
		Time:   t.Format("Mon Jan _2 15:04:05 2006"),
	}

	h.responseBody(w, r, responseCode, response)
}

func (h *selfRouter) livelinessSummaryResponse(w http.ResponseWriter, r *http.Request, status *Status) {
	t := time.Now()

	statusText := Success
	responseCode := http.StatusOK
	if !status.IsHealthy {
		statusText = Failed
		responseCode = http.StatusInternalServerError
	}

	var checkResults map[string]string
	checkResults = make(map[string]string)

	for _, check := range status.Checks {
		checkResults[check.Name] = check.IsSuccess
	}

	response := &LivelinessSummaryResponse{
		Checks: checkResults,
		Status: statusText,
		Time:   t.Format("Mon Jan _2 15:04:05 2006"),
	}

	h.responseBody(w, r, responseCode, response)
}

// Ping godoc
// @Summary Respond to a ping request
// @Description Respond to a ping request with a pong response.
// @Tags health
// @Accept json
// @Accept xml
// @Produce json
// @Produce xml
// @Success 200 {object} health.PingResponse
// @Failure 415 {string} string "Unsupported media type"
// @Router /v1/self/ping [get]
func (h *selfRouter) ping(w http.ResponseWriter, r *http.Request) {
	t := time.Now()

	response := PingResponse{
		Response: fmt.Sprint("Pong - ", t.Format("Mon Jan _2 15:04:05 2006")),
	}

	h.responseBody(w, r, http.StatusOK, response)
}

func (h *selfRouter) readiness(w http.ResponseWriter, r *http.Request) {
}

func (h *selfRouter) started(w http.ResponseWriter, r *http.Request) {
}

func (h *selfRouter) responseBody(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	mediatype, _, err := mime.ParseMediaType(r.Header.Get("Accept"))
	if err != nil {
		log.Error(
			fmt.Sprintf(
				"Invalid 'Accept' header. Error was %v",
				err))

		http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	switch mediatype {
	case "application/xml":
		render.Status(r, status)
		render.XML(w, r, data)
	case "application/json":
		render.Status(r, status)
		render.JSON(w, r, data)
	default:
		log.Error(
			fmt.Sprintf(
				"Invalid media type. Expected either 'application/json' or 'application/xml', got %s.",
				mediatype))

		http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}
}

func (h *healthRouter) Prefix() string {
	return "self"
}

// Routes creates the routes for the health package
func (h *healthRouter) Routes(prefix string, r chi.Router) {
	r.Get(fmt.Sprintf("%s/info", prefix), h.info)
	r.Get(fmt.Sprintf("%s/liveliness", prefix), h.liveliness)
	r.Get(fmt.Sprintf("%s/ping", prefix), h.ping)
	r.Get(fmt.Sprintf("%s/readiness", prefix), h.readiness)
	r.Get(fmt.Sprintf("%s/started", prefix), h.started)
}

func (h *healthRouter) Version() int8 {
	return 1
}
