package types

// APIResponse is the standard JSON envelope returned by all eta-service
// HTTP endpoints, for both success and error cases.
type APIResponse struct {
	Success bool             `json:"success"`
	Data    interface{}      `json:"data,omitempty"`
	Error   *HTTPStatusError `json:"error,omitempty"`
}
