package types

import (
	"errors"
	"net/http"
	"strconv"
)

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

func NewConflictError(displayMessage string) *HTTPStatusError {
	return &HTTPStatusError{
		Message: "conflict",
		DisplayMessage: displayMessage,
		Code: strconv.Itoa(http.StatusConflict),
	}
}

// StatusCode returns the numeric HTTP status code for this error, defaulting
// to 500 if Code is missing or unparsable.
func (e *HTTPStatusError) StatusCode() int {
	if code, err := strconv.Atoi(e.Code); err == nil {
		return code
	}
	return http.StatusInternalServerError
}

// ToHTTPStatusError preserves the intended status code (e.g. conflict, bad
// request) when the given error is already a *HTTPStatusError, falling back
// to a generic internal server error otherwise.
func ToHTTPStatusError(err error) *HTTPStatusError {
	var statusErr *HTTPStatusError
	if errors.As(err, &statusErr) {
		return statusErr
	}
	return NewInternalServerError(err.Error())
}
