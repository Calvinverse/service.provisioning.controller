package health

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/calvinverse/service.provisioning.controller/internal/info"
	"github.com/go-chi/chi"
)

func TestInfoWithAcceptHeaderSetToJson(t *testing.T) {
	request, _ := http.NewRequest("GET", "/ping", nil)
	request.Header.Set("Accept", "application/json")

	w := httptest.NewRecorder()

	instance := &healthRouter{}

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
	request, _ := http.NewRequest("GET", "/ping", nil)
	request.Header.Set("Accept", "application/xml")

	w := httptest.NewRecorder()

	instance := &healthRouter{}

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
	request, _ := http.NewRequest("GET", "/ping", nil)

	w := httptest.NewRecorder()

	instance := &healthRouter{}

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

func TestPingWithAcceptHeaderSetToJson(t *testing.T) {
	request, _ := http.NewRequest("GET", "/ping", nil)
	request.Header.Set("Accept", "application/json")

	w := httptest.NewRecorder()

	instance := &healthRouter{}

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

	instance := &healthRouter{}

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

	instance := &healthRouter{}

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
