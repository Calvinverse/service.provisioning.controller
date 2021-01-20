package health

// Service defines a service that tracks the health of the application.
type Service interface {
	Readiness() (Status, error)

	Liveliness() (Status, error)
}

// Status stores the health status for the application.
type Status struct {
}

// CheckStatus stores the results of a health check.
type CheckStatus struct {
}

type healthReporter struct {
}
