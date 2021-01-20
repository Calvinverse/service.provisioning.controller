package health

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/calvinverse/service.provisioning.controller/internal/info"
	"github.com/go-chi/chi"
)

//
// Mocks
//

type mockHealthService struct {
	status Status
	error  error
}

func (h *mockHealthService) Liveliness() (Status, error) {
	return h.status, h.error
}

func (h *mockHealthService) Readiness() (Status, error) {
	return h.status, h.error
}

type mockError struct{}

func (e *mockError) Error() string {
	return "some text"
}

//
// Info
//

func TestInfoWithAcceptHeaderSetToJson(t *testing.T) {
	request, _ := http.NewRequest("GET", "/info", nil)
	request.Header.Set("Accept", "application/json")

	w := httptest.NewRecorder()

	instance := &selfRouter{}

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		instance.Routes("", r)
	})

	router.ServeHTTP(w, request)

	actualResult := InfoResponse{}
	json.NewDecoder(w.Body).Decode(&actualResult)

	if status := w.Code; status != http.StatusOK {
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			http.StatusOK)
	}

	if actualResult.BuildTime != info.BuildTime() {
		t.Errorf("Handler returned unexpected build time: got %s wanted %s", actualResult.BuildTime, info.BuildTime())
	}

	if actualResult.Revision != info.Revision() {
		t.Errorf("Handler returned unexpected revision: got %s wanted %s", actualResult.Revision, info.Revision())
	}

	if actualResult.Version != info.Version() {
		t.Errorf("Handler returned unexpected build time: got %s wanted %s", actualResult.Version, info.Version())
	}
}

func TestInfoWithAcceptHeaderSetToXml(t *testing.T) {
	request, _ := http.NewRequest("GET", "/info", nil)
	request.Header.Set("Accept", "application/xml")

	w := httptest.NewRecorder()

	instance := &selfRouter{}

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		instance.Routes("", r)
	})

	router.ServeHTTP(w, request)

	actualResult := InfoResponse{}
	xml.NewDecoder(w.Body).Decode(&actualResult)

	if status := w.Code; status != http.StatusOK {
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			http.StatusOK)
	}

	if actualResult.BuildTime != info.BuildTime() {
		t.Errorf("Handler returned unexpected build time: got %s wanted %s", actualResult.BuildTime, info.BuildTime())
	}

	if actualResult.Revision != info.Revision() {
		t.Errorf("Handler returned unexpected revision: got %s wanted %s", actualResult.Revision, info.Revision())
	}

	if actualResult.Version != info.Version() {
		t.Errorf("Handler returned unexpected build time: got %s wanted %s", actualResult.Version, info.Version())
	}
}

func TestInfoWithNoAccept(t *testing.T) {
	request, _ := http.NewRequest("GET", "/info", nil)

	w := httptest.NewRecorder()

	instance := &selfRouter{}

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		instance.Routes("", r)
	})

	router.ServeHTTP(w, request)

	if status := w.Code; status != http.StatusUnsupportedMediaType {
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			http.StatusUnsupportedMediaType)
	}
}

//
// liveliness
//

func TestLivelinessWithFailingHealthAndHeaderSetToJson(t *testing.T) {
	request, _ := http.NewRequest("GET", "/liveliness", nil)
	request.Header.Set("Accept", "application/json")

	w := httptest.NewRecorder()

	healthService := &mockHealthService{
		error: &mockError{},
	}
	instance := &selfRouter{
		healthService: healthService,
	}

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		instance.Routes("", r)
	})

	router.ServeHTTP(w, request)

	actualResult := LivelinessSummaryResponse{}
	json.NewDecoder(w.Body).Decode(&actualResult)

	if status := w.Code; status != http.StatusInternalServerError {
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			http.StatusInternalServerError)
	}

	if actualResult.Status != Failed {
		t.Errorf("Handler returned unexpected status: got %s wanted %s", actualResult.Status, Failed)
	}

	if len(actualResult.Checks) != 0 {
		t.Errorf("Handler returned unexpected number of checks: got %d wanted %d", len(actualResult.Checks), 0)
	}
}

func TestLivelinessWithNoAccept(t *testing.T) {
	request, _ := http.NewRequest("GET", "/liveliness", nil)

	w := httptest.NewRecorder()

	healthService := &mockHealthService{
		error: &mockError{},
	}
	instance := &selfRouter{
		healthService: healthService,
	}

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		instance.Routes("", r)
	})

	router.ServeHTTP(w, request)

	if status := w.Code; status != http.StatusUnsupportedMediaType {
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			http.StatusUnsupportedMediaType)
	}
}

func TestLivelinessDetailedWithAcceptHeaderSetToJson(t *testing.T) {
	request, _ := http.NewRequest("GET", "/liveliness", nil)
	request.Header.Set("Accept", "application/json")
	q := request.URL.Query()
	q.Add("type", "detailed")

	w := httptest.NewRecorder()

	numberOfChecks := 10

	var checks []CheckResult
	for i := 0; i < numberOfChecks; i++ {
		check := CheckResult{
			IsSuccess: true,
			Name:      strconv.Itoa(i),
			Timestamp: time.Date(2021, time.January, i, i, i, i, 0, time.Local),
		}
		checks = append(checks, check)
	}

	status := Status{
		Checks:    checks,
		IsHealthy: true,
	}

	healthService := &mockHealthService{
		status: status,
		error:  nil,
	}

	instance := &selfRouter{
		healthService: healthService,
	}

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		instance.Routes("", r)
	})

	router.ServeHTTP(w, request)

	actualResult := LivelinessDetailedResponse{}
	json.NewDecoder(w.Body).Decode(&actualResult)

	if status := w.Code; status != http.StatusOK {
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			http.StatusOK)
	}

	if actualResult.Status != Success {
		t.Errorf("Handler returned unexpected status: got %s wanted %s", actualResult.Status, Success)
	}

	if len(actualResult.Checks) != numberOfChecks {
		t.Errorf("Handler returned unexpected number of checks: got %d wanted %d", len(actualResult.Checks), numberOfChecks)
	}

	for i, c := range actualResult.Checks {
		if c.Name != strconv.Itoa(i) {
			t.Errorf("Check %d has an unexpected name: got %s wanted %s", i, c.Name, strconv.Itoa(i))
		}

		if c.Status != Success {
			t.Errorf("Check %d had an unexpected status. Expected Success got %s", i, c.Status)
		}

		parsedTime, err := time.Parse(time.RFC3339, c.Timestamp)
		if err != nil {
			t.Errorf("Check %d contained a timestamp that was not parsable. Got %s", i, c.Timestamp)
		}

		expectedTime := time.Date(2021, time.January, i, i, i, i, 0, time.Local)
		if parsedTime != expectedTime {
			t.Errorf("Check %d had an unexpected timestamp. Got %s wanted %s", i, c.Timestamp, expectedTime.Format(time.RFC3339))
		}
	}
}

func TestLivelinessDetailedWithAcceptHeaderSetToXml(t *testing.T) {
	request, _ := http.NewRequest("GET", "/liveliness", nil)
	request.Header.Set("Accept", "application/xml")
	q := request.URL.Query()
	q.Add("type", "detailed")

	w := httptest.NewRecorder()

	numberOfChecks := 10

	var checks []CheckResult
	for i := 0; i < numberOfChecks; i++ {
		check := CheckResult{
			IsSuccess: true,
			Name:      strconv.Itoa(i),
			Timestamp: time.Date(2021, time.January, i, i, i, i, 0, time.Local),
		}
		checks = append(checks, check)
	}

	status := Status{
		Checks:    checks,
		IsHealthy: true,
	}

	healthService := &mockHealthService{
		status: status,
		error:  nil,
	}

	instance := &selfRouter{
		healthService: healthService,
	}

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		instance.Routes("", r)
	})

	router.ServeHTTP(w, request)

	actualResult := LivelinessDetailedResponse{}
	xml.NewDecoder(w.Body).Decode(&actualResult)

	if status := w.Code; status != http.StatusOK {
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			http.StatusOK)
	}

	if actualResult.Status != Success {
		t.Errorf("Handler returned unexpected status: got %s wanted %s", actualResult.Status, Success)
	}

	if len(actualResult.Checks) != numberOfChecks {
		t.Errorf("Handler returned unexpected number of checks: got %d wanted %d", len(actualResult.Checks), numberOfChecks)
	}

	for i, c := range actualResult.Checks {
		if c.Name != strconv.Itoa(i) {
			t.Errorf("Check %d has an unexpected name: got %s wanted %s", i, c.Name, strconv.Itoa(i))
		}

		if c.Status != Success {
			t.Errorf("Check %d had an unexpected status. Expected Success got %s", i, c.Status)
		}

		parsedTime, err := time.Parse(time.RFC3339, c.Timestamp)
		if err != nil {
			t.Errorf("Check %d contained a timestamp that was not parsable. Got %s", i, c.Timestamp)
		}

		expectedTime := time.Date(2021, time.January, i, i, i, i, 0, time.Local)
		if parsedTime != expectedTime {
			t.Errorf("Check %d had an unexpected timestamp. Got %s wanted %s", i, c.Timestamp, expectedTime.Format(time.RFC3339))
		}
	}
}

func TestLivelinessSummaryWithAcceptHeaderSetToJson(t *testing.T) {
	request, _ := http.NewRequest("GET", "/liveliness", nil)
	request.Header.Set("Accept", "application/json")
	q := request.URL.Query()
	q.Add("type", "summary")

	w := httptest.NewRecorder()

	numberOfChecks := 10

	var checks []CheckResult
	for i := 0; i < numberOfChecks; i++ {
		check := CheckResult{
			IsSuccess: true,
			Name:      strconv.Itoa(i),
			Timestamp: time.Date(2021, time.January, i, i, i, i, 0, time.Local),
		}
		checks = append(checks, check)
	}

	status := Status{
		Checks:    checks,
		IsHealthy: true,
	}

	healthService := &mockHealthService{
		status: status,
		error:  nil,
	}

	instance := &selfRouter{
		healthService: healthService,
	}

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		instance.Routes("", r)
	})

	router.ServeHTTP(w, request)

	actualResult := LivelinessSummaryResponse{}
	json.NewDecoder(w.Body).Decode(&actualResult)

	if status := w.Code; status != http.StatusOK {
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			http.StatusOK)
	}

	if actualResult.Status != Success {
		t.Errorf("Handler returned unexpected status: got %s wanted %s", actualResult.Status, Success)
	}

	if len(actualResult.Checks) != numberOfChecks {
		t.Errorf("Handler returned unexpected number of checks: got %d wanted %d", len(actualResult.Checks), numberOfChecks)
	}

	for k, v := range actualResult.Checks {
		if k != "" {
			t.Errorf("Check has an unexpected name: got %s wanted %s", k, "")
		}

		if v != Success {
			t.Errorf("Check had an unexpected status. Expected Success got %s", v)
		}
	}
}

func TestLivelinessSummaryWithAcceptHeaderSetToXml(t *testing.T) {
	t.Fail()
}

//
// ping
//

func TestPingWithAcceptHeaderSetToJson(t *testing.T) {
	request, _ := http.NewRequest("GET", "/ping", nil)
	request.Header.Set("Accept", "application/json")

	w := httptest.NewRecorder()

	instance := &selfRouter{}

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		instance.Routes("", r)
	})

	router.ServeHTTP(w, request)

	actualResult := PingResponse{}
	json.NewDecoder(w.Body).Decode(&actualResult)

	if status := w.Code; status != http.StatusOK {
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			http.StatusOK)
	}

	if !strings.HasPrefix(actualResult.Response, "Pong - ") {
		t.Errorf("Handler returned unexpected response: got %s wanted 'Pong'", actualResult.Response)
	}
}

func TestPingWithAcceptHeaderSetToXml(t *testing.T) {
	request, _ := http.NewRequest("GET", "/ping", nil)
	request.Header.Set("Accept", "application/xml")

	w := httptest.NewRecorder()

	instance := &selfRouter{}

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		instance.Routes("", r)
	})

	router.ServeHTTP(w, request)

	actualResult := PingResponse{}
	xml.NewDecoder(w.Body).Decode(&actualResult)

	if status := w.Code; status != http.StatusOK {
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			http.StatusOK)
	}

	if !strings.HasPrefix(actualResult.Response, "Pong - ") {
		t.Errorf("Handler returned unexpected response: got %s wanted 'Pong'", actualResult.Response)
	}
}

func TestPingWithNoAccept(t *testing.T) {
	request, _ := http.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()

	instance := &selfRouter{}

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		instance.Routes("", r)
	})

	router.ServeHTTP(w, request)

	if status := w.Code; status != http.StatusUnsupportedMediaType {
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			http.StatusUnsupportedMediaType)
	}
}

// readiness - json
// readiness - xml
// readiness - no-accept

// started - json
// started - xml
// started - no-accept
