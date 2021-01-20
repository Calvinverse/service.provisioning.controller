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
)

// InfoResponse stores the response to an info request
type InfoResponse struct {
	BuildTime string `json:"buildtime"`
	Revision  string `json:"revision"`
	Version   string `json:"version"`
}

// LivelinessDetailedResponse stores detailed information about the liveliness of the application, indicating if the application is healthy
type LivelinessDetailedResponse struct {
	// Global status
	Status string `json:"status"`

	// Time the liveliness response was created at
	Time string `json:"time"`

	Checks []CheckStatus `json:"checks"`
}

// LivelinessSummaryResponse stores condensed information about the liveliness of the application, indicating if the application is healthy
type LivelinessSummaryResponse struct {
	// Global status
	Status string `json:"status"`

	// Status of all health checks
	Checks []string `json:"checks"`
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

// NewHealthAPIRouter returns an APIRouter instance for the health routes.
func NewHealthAPIRouter() router.APIRouter {
	return &healthRouter{}
}

type healthRouter struct{}

// Info godoc
// @Summary Respond to an info request
// @Description Respond to an info request with information about the application.
// @Tags health
// @Accept  json
// @Produce  json
// @Success 200 {object} health.InfoResponse
// @Router /v1/self/info [get]
func (h *healthRouter) info(w http.ResponseWriter, r *http.Request) {
	response := InfoResponse{
		BuildTime: info.BuildTime(),
		Revision:  info.Revision(),
		Version:   info.Version(),
	}

	render.JSON(w, r, response)
}

func (h *healthRouter) liveliness(w http.ResponseWriter, r *http.Request) {
	// Render liveliness status - If the 'json=detailed' query string is present show the
	// extended JSON. Otherwise just return a short JSON response that only changes
	// when the health status changes.
	//
	// The short version is used for Consul etc. which prefer having unchanged responses
	// for unchanged conditions

	// r.URL.Query().Get("type") == "detailed"
}

// Ping godoc
// @Summary Respond to a ping request
// @Description Respond to a ping request with information about the application.
// @Tags health
// @Accept  json
// @Produce  json
// @Success 200 {object} health.PingResponse
// @Router /v1/self/ping [get]
func (h *healthRouter) ping(w http.ResponseWriter, r *http.Request) {
	t := time.Now()

	response := PingResponse{
		Response: fmt.Sprint("Pong - ", t.Format("Mon Jan _2 15:04:05 2006")),
	}

	h.responseBody(w, r, response)
}

func (h *healthRouter) readiness(w http.ResponseWriter, r *http.Request) {

}

func (h *healthRouter) started(w http.ResponseWriter, r *http.Request) {

}

func (h *healthRouter) responseBody(w http.ResponseWriter, r *http.Request, data interface{}) {
	mediatype, _, err := mime.ParseMediaType(r.Header.Get("Accept"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	switch mediatype {
	case "application/xml":
		render.XML(w, r, data)
	case "application/json":
		render.JSON(w, r, data)
	default:
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
