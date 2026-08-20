package types

type ErrorResponse struct {
	Error HTTPStatusError `json:"error,omitempty"`
}

type HTTPStatusError struct {
	Message        string `json:"message"`
	Code           string `json:"code"`
	DisplayMessage string `json:"displayMessage"`
}

func (e *HTTPStatusError) Error() string {
	return e.Message
}
