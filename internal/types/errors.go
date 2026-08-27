package types

import (
	"net/http"
	"strconv"
)

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

func NewInternalServerError(displayMessage string) *HTTPStatusError {
	return &HTTPStatusError{
		Message: "internal server error",
		DisplayMessage: displayMessage,
		Code: strconv.Itoa(http.StatusInternalServerError),
	}
}

func NewBadRequestError(displayMessage string) *HTTPStatusError {
	return &HTTPStatusError{
		Message: "bad request",
		DisplayMessage: displayMessage,
		Code: strconv.Itoa(http.StatusBadRequest),
	}
}
