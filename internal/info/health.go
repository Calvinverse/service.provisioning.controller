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
	instance HealthReporter
)

// HealthReporter defines a service that tracks the health of the application.
type HealthReporter interface {
	// Liveliness returns the status indicating if the application is healthy while processing requests.
	Liveliness() (HealthStatus, error)

	// Readiness returns the status indicating if the application is ready to process requests.
	Readiness() (HealthStatus, error)

	// Add a health check
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

// GetServiceWithDefaultSettings returns a health service with the default settings.
func GetServiceWithDefaultSettings() HealthReporter {
	once.Do(func() {
		instance = &healthReporter{
			instance: gosundheit.New(),
		}
	})

	return instance
}

// GetServiceWithHealthInstance returns a health service with the provided health instance. Note: for testing purposes only!
func GetServiceWithHealthInstance(healthInstance gosundheit.Health) HealthReporter {
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

func (h *healthReporter) Liveliness() (HealthStatus, error) {
	// Return the status of the different health checks
	return HealthStatus{}, nil
}

func (h *healthReporter) Readiness() (HealthStatus, error) {
	// If all health checks have been registered then we are good
	return HealthStatus{}, nil
}
