package health

import (
	"bytes"
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

func TestInfoyWithAcceptHeaderSetToJson(t *testing.T) {
	request := setupRequest("/info", "application/json", make(map[string]string))

	w := httptest.NewRecorder()

	router := setupHttpRouter(&selfRouter{})
	router.ServeHTTP(w, request)

	validateInfoWithAcceptHeader(t, w, decodeJSONFromResponseBody)
}

func TestInfoyWithAcceptHeaderSetToXml(t *testing.T) {
	request := setupRequest("/info", "application/xml", make(map[string]string))

	w := httptest.NewRecorder()

	router := setupHttpRouter(&selfRouter{})
	router.ServeHTTP(w, request)

	validateInfoWithAcceptHeader(t, w, decodeXMLFromResponseBody)
}

func TestInfoWithoutHeader(t *testing.T) {
	request := setupRequest("/info", "", make(map[string]string))

	w := httptest.NewRecorder()

	router := setupHttpRouter(&selfRouter{})
	router.ServeHTTP(w, request)

	validateWithoutAcceptHeader(t, w, decodeJSONFromResponseBody)
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
	request := setupRequest("/liveliness", "", make(map[string]string))

	w := httptest.NewRecorder()

	healthService := &mockHealthService{
		error: &mockError{},
	}
	instance := &selfRouter{
		healthService: healthService,
	}

	router := setupHttpRouter(instance)
	router.ServeHTTP(w, request)

	validateWithoutAcceptHeader(t, w, decodeJSONFromResponseBody)
}

func TestLivelinessDetailedWithAcceptHeaderSetToJson(t *testing.T) {
	queryParameters := map[string]string{
		"type": "detailed",
	}

	request := setupRequest("/liveliness", "application/json", queryParameters)

	w := httptest.NewRecorder()

	numberOfChecks := 2

	healthService := createHealthServiceWithChecks(numberOfChecks)
	instance := &selfRouter{
		healthService: healthService,
	}

	router := setupHttpRouter(instance)
	router.ServeHTTP(w, request)

	actualResult := LivelinessDetailedResponse{}
	json.NewDecoder(w.Body).Decode(&actualResult)

	if status := w.Code; status != http.StatusOK {
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			http.StatusOK)
	}

	validateLivelinessDetailedResponse(t, numberOfChecks, actualResult)
}

func TestLivelinessDetailedWithAcceptHeaderSetToXml(t *testing.T) {
	queryParameters := map[string]string{
		"type": "detailed",
	}

	request := setupRequest("/liveliness", "application/xml", queryParameters)

	w := httptest.NewRecorder()

	numberOfChecks := 2

	healthService := createHealthServiceWithChecks(numberOfChecks)
	instance := &selfRouter{
		healthService: healthService,
	}

	router := setupHttpRouter(instance)
	router.ServeHTTP(w, request)

	actualResult := LivelinessDetailedResponse{}
	xml.NewDecoder(w.Body).Decode(&actualResult)

	if status := w.Code; status != http.StatusOK {
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			http.StatusOK)
	}

	validateLivelinessDetailedResponse(t, numberOfChecks, actualResult)
}

func TestLivelinessSummaryWithAcceptHeaderSetToJson(t *testing.T) {
	queryParameters := map[string]string{
		"type": "summary",
	}

	request := setupRequest("/liveliness", "application/json", queryParameters)

	w := httptest.NewRecorder()

	numberOfChecks := 2

	healthService := createHealthServiceWithChecks(numberOfChecks)
	instance := &selfRouter{
		healthService: healthService,
	}

	router := setupHttpRouter(instance)
	router.ServeHTTP(w, request)

	actualResult := LivelinessSummaryResponse{}
	json.NewDecoder(w.Body).Decode(&actualResult)

	if status := w.Code; status != http.StatusOK {
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			http.StatusOK)
	}

	validateLivelinessSummaryResponse(t, numberOfChecks, actualResult)
}

func TestLivelinessSummaryWithAcceptHeaderSetToXml(t *testing.T) {
	queryParameters := map[string]string{
		"type": "summary",
	}

	request := setupRequest("/liveliness", "application/xml", queryParameters)

	w := httptest.NewRecorder()

	numberOfChecks := 2

	healthService := createHealthServiceWithChecks(numberOfChecks)
	instance := &selfRouter{
		healthService: healthService,
	}

	router := setupHttpRouter(instance)
	router.ServeHTTP(w, request)

	actualResult := LivelinessSummaryResponse{}
	xml.NewDecoder(w.Body).Decode(&actualResult)

	if status := w.Code; status != http.StatusOK {
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			http.StatusOK)
	}

	validateLivelinessSummaryResponse(t, numberOfChecks, actualResult)
}

//
// ping
//

func TestPingyWithAcceptHeaderSetToJson(t *testing.T) {
	request := setupRequest("/ping", "application/json", make(map[string]string))

	w := httptest.NewRecorder()

	router := setupHttpRouter(&selfRouter{})
	router.ServeHTTP(w, request)

	validatePingWithAcceptHeader(t, w, decodeJSONFromResponseBody)
}

func TestPingyWithAcceptHeaderSetToXml(t *testing.T) {
	request := setupRequest("/ping", "application/xml", make(map[string]string))

	w := httptest.NewRecorder()

	router := setupHttpRouter(&selfRouter{})
	router.ServeHTTP(w, request)

	validatePingWithAcceptHeader(t, w, decodeXMLFromResponseBody)
}

func TestPingWithoutHeader(t *testing.T) {
	request := setupRequest("/ping", "", make(map[string]string))

	w := httptest.NewRecorder()

	router := setupHttpRouter(&selfRouter{})
	router.ServeHTTP(w, request)

	validateWithoutAcceptHeader(t, w, decodeJSONFromResponseBody)
}

// readiness - json
// readiness - xml
// readiness - no-accept

// started - json
// started - xml
// started - no-accept

//
// Helper functions
//

type decodeResponseBody func(buffer *bytes.Buffer, v interface{}) error

func decodeJSONFromResponseBody(buffer *bytes.Buffer, v interface{}) error {
	return json.NewDecoder(buffer).Decode(v)
}

func decodeXMLFromResponseBody(buffer *bytes.Buffer, v interface{}) error {
	return xml.NewDecoder(buffer).Decode(v)
}

type validateResponse func(t *testing.T, w *httptest.ResponseRecorder, decode decodeResponseBody)

//
// Setup functions
//

func createHealthServiceWithChecks(numberOfChecks int) *mockHealthService {
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
	return healthService
}

func setupHttpRouter(instance *selfRouter) *chi.Mux {
	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		instance.Routes("", r)
	})

	return router
}

func setupRequest(path string, acceptHeader string, queryItems map[string]string) *http.Request {
	request, _ := http.NewRequest("GET", path, nil)

	if acceptHeader != "" {
		request.Header.Set("Accept", acceptHeader)
	}

	for k, v := range queryItems {
		q := request.URL.Query()
		q.Add(k, v)
		request.URL.RawQuery = q.Encode()
	}

	return request
}

//
// Validation functions
//

func validateWithoutAcceptHeader(t *testing.T, w *httptest.ResponseRecorder, decode decodeResponseBody) {
	if status := w.Code; status != http.StatusUnsupportedMediaType {
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			http.StatusUnsupportedMediaType)
	}
}

func validateInfoWithAcceptHeader(t *testing.T, w *httptest.ResponseRecorder, decode decodeResponseBody) {
	actualResult := InfoResponse{}
	decode(w.Body, &actualResult)

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

func validatePingWithAcceptHeader(t *testing.T, w *httptest.ResponseRecorder, decode decodeResponseBody) {
	actualResult := PingResponse{}
	decode(w.Body, &actualResult)

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

func validateLivelinessDetailedResponse(t *testing.T, expectedNumberOfChecks int, result LivelinessDetailedResponse) {
	if result.Status != Success {
		t.Errorf("Handler returned unexpected status: got %s wanted %s", result.Status, Success)
	}

	if len(result.Checks) != expectedNumberOfChecks {
		t.Errorf("Handler returned unexpected number of checks: got %d wanted %d", len(result.Checks), expectedNumberOfChecks)
	}

	for i, c := range result.Checks {
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

func validateLivelinessSummaryResponse(t *testing.T, expectedNumberOfChecks int, result LivelinessSummaryResponse) {
	if result.Status != Success {
		t.Errorf("Handler returned unexpected status: got %s wanted %s", result.Status, Success)
	}

	if len(result.Checks) != expectedNumberOfChecks {
		t.Errorf("Handler returned unexpected number of checks: got %d wanted %d", len(result.Checks), expectedNumberOfChecks)
	}

	for i, k := range result.Checks {

		expectedName := strconv.Itoa(i)
		if k.Name != expectedName {
			t.Errorf("Check has an unexpected name: got %s wanted %s", k.Name, expectedName)
		}

		if k.Status != Success {
			t.Errorf("Check had an unexpected status. Expected Success got %s", k.Status)
		}
	}
}
