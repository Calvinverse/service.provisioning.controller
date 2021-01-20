package health

import (
	"sync"

	gosundheit "github.com/AppsFlyer/go-sundheit"
)

const (
	Failed  string = "failed"
	Success string = "success"
)

var (
	once     sync.Once
	instance Service
)

// Service defines a service that tracks the health of the application.
type Service interface {
	Liveliness() (Status, error)

	Readiness() (Status, error)
}

// Status stores the health status for the application.
type Status struct {
	// Checks returns the collection of checks that were executed.
	Checks []CheckResult

	// IsHealthy returns the health status for the application.
	IsHealthy bool
}

// CheckResult stores the results of a health check.
type CheckResult struct {
	// Name returns the name of the check.
	Name string

	// IsSuccess returns the status of the check.
	IsSuccess bool
}

// GetServiceWithDefaultSettings returns a health service with the default settings.
func GetServiceWithDefaultSettings() Service {
	once.Do(func() {
		instance = &healthReporter{
			instance: gosundheit.New(),
		}
	})

	return instance
}

// GetServiceWithHealthInstance returns a health service with the provided health instance. Note: for testing purposes only!
func GetServiceWithHealthInstance(healthInstance gosundheit.Health) Service {
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

func (h *healthReporter) Liveliness() (Status, error) {
	// Return the status of the different health checks
	return Status{}, nil
}

func (h *healthReporter) Readiness() (Status, error) {
	// If all health checks have been registered then we are good
	return Status{}, nil
}
