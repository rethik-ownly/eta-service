package types

// This file contains all the required types for any client of the form : `A library/object your service uses to talk to another system` 

type HealthCheck struct {
	Client string `json:"client"`
	Status string `json:"status"`
}

type HealthCheckResponse struct {
	Checks []HealthCheck `json:"checks"`
	Status string `json:"status"`
}

// HTTP response returned when this service calls another service
type HTTPResponse struct {
	StatusCode int
	Body       []byte
}

