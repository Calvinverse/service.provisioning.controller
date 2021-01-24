package info

import (
	"sync"
	"time"

	gosundheit "github.com/AppsFlyer/go-sundheit"
)

const (
	// Failed indicates the health or a health check is failing.
	Failed string = "failed"

	// Success indicates the health or a health check is successful.
	Success string = "success"
)

var (
	once     sync.Once
	instance StatusReporter
)

// StatusReporter defines a service that tracks the health of the application.
type StatusReporter interface {
	// Liveliness returns the status indicating if the application is healthy while processing requests.
	Liveliness() (*HealthStatus, error)

	// Readiness returns the status indicating if the application is ready to process requests.
	Readiness() (*HealthStatus, error)

	// Started returns the information about the start of the application.
	Started() (*StartedStatus, error)
}

type HealthCenter interface {
	// Add a health check
	// Add a readiness check
	// Add a started check
}

// HealthStatus stores the health status for the application.
type HealthStatus struct {
	// Checks returns the collection of checks that were executed.
	Checks []HealthCheckResult

	// IsHealthy returns the health status for the application.
	IsHealthy bool
}

// HealthCheckResult stores the results of a health check.
type HealthCheckResult struct {
	// IsSuccess returns the status of the check.
	IsSuccess bool

	// Name returns the name of the check.
	Name string

	// The last time the check result was updated.
	Timestamp time.Time
}

// StartedStatus stores the application start information.
type StartedStatus struct {
	// The time the application was started.
	Timestamp time.Time
}

// GetServiceWithDefaultSettings returns a health service with the default settings.
func GetServiceWithDefaultSettings() StatusReporter {
	once.Do(func() {
		instance = &healthReporter{
			instance: gosundheit.New(),
		}
	})

	return instance
}

// GetServiceWithHealthInstance returns a health service with the provided health instance. Note: for testing purposes only!
func GetServiceWithHealthInstance(healthInstance gosundheit.Health) StatusReporter {
	once.Do(func() {
		instance = &healthReporter{
			instance: healthInstance,
		}
	})

	return instance
}

type healthReporter struct {
	instance gosundheit.Health
}

func (h *healthReporter) Liveliness() (*HealthStatus, error) {
	checkResults, healthy := h.instance.Results()

	var checks []HealthCheckResult
	checks = make([]HealthCheckResult, 0, len(checkResults))
	for name, check := range checkResults {
		checkResult := HealthCheckResult{
			IsSuccess: check.IsHealthy(),
			Name:      name,
			Timestamp: check.Timestamp,
		}

		checks = append(checks, checkResult)
	}

	// Return the status of the different health checks
	result := &HealthStatus{
		Checks:    checks,
		IsHealthy: healthy,
	}
	return result, nil
}

func (h *healthReporter) Readiness() (*HealthStatus, error) {
	// If all health checks have been registered then we are good
	return &HealthStatus{}, nil
}

func (h *healthReporter) Started() (*StartedStatus, error) {
	return &StartedStatus{}, nil
}
