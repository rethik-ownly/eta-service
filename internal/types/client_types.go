package types

type HealthCheck struct {
	Client string `json:"client"`
	Status string `json:"status"`
}

type HealthCheckResponse struct {
	Checks []HealthCheck `json:"checks"`
	Status string `json:"status"`
}

type HTTPResponse struct {
	StatusCode int
	Body       []byte
}

