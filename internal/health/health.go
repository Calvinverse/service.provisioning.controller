package health

type HealthService interface {
	Readiness() HealthStatus, error

	Liveliness() HealthStatus, error
}

type HealthStatus struct {

}

type healthReporter {

}