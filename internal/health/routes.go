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

// CheckInformation stores information about the status of a health check.
type CheckInformation struct {
	// Name returns the name of the health check.
	Name string `json:"name"`

	// Status returns the status of the health check, either success or failure.
	Status string `json:"status"`

	// Timestamp returns the time the healtcheck was executed.
	Timestamp string `json:"timestamp"`
}

// InfoResponse stores the response to an info request
type InfoResponse struct {
	// BuildTime stores the date and time the application was built.
	BuildTime string `json:"buildtime"`

	// Revision stores the GIT SHA of the commit on which the application build was based.
	Revision string `json:"revision"`

	// Version stores the version number of the application.
	Version string `json:"version"`
}

// LivelinessDetailedResponse stores detailed information about the liveliness of the application, indicating if the application is healthy
type LivelinessDetailedResponse struct {
	// Status of all the health checks
	Checks []CheckInformation `json:"checks"`

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
		healthStatus = Status{
			Checks:    make([]CheckResult, 0, 0),
			IsHealthy: false,
		}

		h.livelinessSummaryResponse(w, r, &healthStatus)
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

	statusText := statusToText(status.IsHealthy)
	responseCode := statusToResponseCode(status.IsHealthy)

	var checkResults []CheckInformation
	checkResults = make([]CheckInformation, len(status.Checks))
	for _, check := range status.Checks {

		result := CheckInformation{
			Name:      check.Name,
			Status:    statusToText(check.IsSuccess),
			Timestamp: check.Timestamp.Format(time.RFC3339),
		}
		checkResults = append(checkResults, result)
	}

	response := &LivelinessDetailedResponse{
		Checks: checkResults,
		Status: statusText,
		Time:   t.Format(time.RFC3339),
	}

	h.responseBody(w, r, responseCode, response)
}

func (h *selfRouter) livelinessSummaryResponse(w http.ResponseWriter, r *http.Request, status *Status) {
	t := time.Now()

	statusText := statusToText(status.IsHealthy)
	responseCode := statusToResponseCode(status.IsHealthy)

	var checkResults map[string]string
	checkResults = make(map[string]string)

	for _, check := range status.Checks {
		checkResults[check.Name] = statusToText(check.IsSuccess)
	}

	response := &LivelinessSummaryResponse{
		Checks: checkResults,
		Status: statusText,
		Time:   t.Format(time.RFC3339),
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
		Response: fmt.Sprint("Pong - ", t.Format(time.RFC3339)),
	}

	h.responseBody(w, r, http.StatusOK, response)
}

// Readiness godoc
// @Summary Respond to an readiness request
// @Description Respond to an readiness request with information about ability of the application to start serving requests.
// @Tags health
// @Accept json
// @Accept xml
// @Produce json
// @Produce xml
// @Success 200 {object} health.ReadinessResponse
// @Router /v1/self/readiness [get]
func (h *selfRouter) readiness(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNotImplemented)
}

// Started godoc
// @Summary Respond to an started request
// @Description Respond to an started request with information indicating if the application has started successfully.
// @Tags health
// @Accept json
// @Accept xml
// @Produce json
// @Produce xml
// @Success 200 {object} health.StartedResponse
// @Router /v1/self/started [get]
func (h *selfRouter) started(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNotImplemented)
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

func statusToResponseCode(status bool) int {
	responseCode := http.StatusOK
	if !status {
		responseCode = http.StatusInternalServerError
	}
	return responseCode
}

func statusToText(status bool) string {
	statusText := Success
	if !status {
		statusText = Failed
	}

	return statusText
}
