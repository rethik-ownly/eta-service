package utils

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/nutanalabs/eta-service/internal/types"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/nutanalabs/rapido-http-go/v3/errorhandler"
)

type resolvedHTTPError struct {
	statusCode int
	httpErr    *types.HTTPStatusError
}

func ResolveStatusCode(err error) int {
	statusCode, _ := ResolveHTTPStatusError(err, "")
	return statusCode
}

func ResolveHTTPStatusError(err error, fallbackInternalMessage string) (int , *types.HTTPStatusError) {
	if err == nil {
		return http.StatusInternalServerError, types.NewInternalServerError("unexpected nil error")
	}

	if resolved, ok := resolveFromHTTPStatusError(err, fallbackInternalMessage); ok {
		return resolved.statusCode, resolved.httpErr
	}

	if resolved, ok := resolveFromMongoError(err); ok {
		return resolved.statusCode, resolved.httpErr
	}

	if resolved, ok := resolveFromRapidoResponseError(err, fallbackInternalMessage); ok {
		return resolved.statusCode, resolved.httpErr
	}

	resolved := resolveDefault(err, fallbackInternalMessage)
	return resolved.statusCode, resolved.httpErr
}

func resolveFromHTTPStatusError(err error, fallbackInternalMessage string) (resolvedHTTPError, bool) {
	var typedErr *types.HTTPStatusError

	if !errors.As(err, &typedErr) {
		return resolvedHTTPError{}, false
	}

	statusCode, ok := parseStatusCode(typedErr.Code)
	if !ok {
		statusCode = http.StatusInternalServerError
	}
	displayMessage := typedErr.DisplayMessage

	// If internal server error then don't show the actual error to the end user instead show a fallback message
	if statusCode == http.StatusInternalServerError && fallbackInternalMessage != "" {
		displayMessage = fallbackInternalMessage
	}

	return resolvedHTTPError{
		statusCode: statusCode,
		httpErr: statusErrorFromCode(statusCode, displayMessage),
	}, true
}

func resolveFromMongoError(err error) (resolvedHTTPError, bool) {
	if errors.Is(err, mongo.ErrNoDocuments) {
		return resolvedHTTPError{
			statusCode: http.StatusNotFound,
			httpErr: statusErrorFromCode(http.StatusNotFound, err.Error()),
		}, true
	}

	if mongo.IsDuplicateKeyError(err) {
		return resolvedHTTPError{
			statusCode: http.StatusConflict,
			httpErr: statusErrorFromCode(http.StatusConflict, err.Error()),
		}, true
	}

	return resolvedHTTPError{}, false
}

// resolves rapido-http-go errors
func resolveFromRapidoResponseError(err error, fallbackInternalMessage string) (resolvedHTTPError, bool) {
	var responseErr *errorhandler.ResponseError
	if !errors.As(err, &responseErr) {
		return resolvedHTTPError{}, false
	}

	statusCode, ok := parseStatusCode(strconv.Itoa(responseErr.StatusCode))
	if !ok {
		statusCode = http.StatusInternalServerError
	}

	displayMessage := strings.TrimSpace(string(responseErr.Body))
	if displayMessage == "" {
		displayMessage = err.Error()
	}

	if statusCode == http.StatusInternalServerError && fallbackInternalMessage != "" {
		displayMessage = fallbackInternalMessage
	}

	return resolvedHTTPError{
		statusCode: statusCode,
		httpErr: statusErrorFromCode(statusCode, displayMessage),
	}, true
}

func parseStatusCode(code string) (int, bool) {
	statusCode, err := strconv.Atoi(strings.TrimSpace(code))
	if err != nil || statusCode < http.StatusBadRequest || statusCode > 599 {
		return 0, false
	}
	return statusCode, true
}


func resolveDefault(err error, fallbackInternalMessage string) resolvedHTTPError {
	displayMessage := err.Error()
	if fallbackInternalMessage != "" {
		displayMessage = fallbackInternalMessage
	}
	return resolvedHTTPError{
		statusCode: http.StatusInternalServerError,
		httpErr:    types.NewInternalServerError(displayMessage),
	}
}

func statusErrorFromCode(statusCode int, displayMessage string) *types.HTTPStatusError {
	switch statusCode {
	case http.StatusBadRequest:
		return types.NewBadRequestError(displayMessage)
	case http.StatusUnauthorized:
		return types.NewUnAuthorizedError(displayMessage)
	case http.StatusForbidden:
		return types.NewForbiddenError(displayMessage)
	case http.StatusNotFound:
		return types.NewNoResultsFoundRequestError(displayMessage)
	case http.StatusConflict:
		return types.NewConflictError(displayMessage)
	case http.StatusGone:
		return types.NewGoneError(displayMessage)
	case http.StatusTooManyRequests:
		return types.NewTooManyRequestsError(displayMessage)
	case http.StatusInternalServerError:
		return types.NewInternalServerError(displayMessage)
	default:
		message := strings.ToLower(http.StatusText(statusCode))
		if message == "" {
			message = "error"
		}
		return &types.HTTPStatusError{
			Message:        message,
			Code:           strconv.Itoa(statusCode),
			DisplayMessage: displayMessage,
		}
	}
}