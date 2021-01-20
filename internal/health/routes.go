package health

import (
	"fmt"
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

type LivelinessDetailedResponse struct {
	// Global status
	Status string `json:"status"`

	// Time the liveliness response was created at
	Time string `json:"time"`

	Checks []HealthCheckStatus `json:"checks"`
}

// Liveliness stores the response to a liveliness request
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

type ReadinessResponse struct {
	Time string `json:"time"`
}

type StartedResponse struct {
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

	// return the correct document type based on the application type
	render.JSON(w, r, response)
}

func (h *healthRouter) readiness(w http.ResponseWriter, r *http.Request) {

}

func (h *healthRouter) started(w http.ResponseWriter, r *http.Request) {

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
