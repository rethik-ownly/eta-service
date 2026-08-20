package utils

import (
	"errors"
	"net/http"

	"github.com/nutanalabs/eta-service/internal/types"
	"github.com/nutanalabs/rapido-http-go/v3/errorhandler"
)

type HTTPUtils interface {
	ParseErrorResponse(err error) types.HTTPResponse
	FetchErrorStatusCode(response *errorhandler.ResponseError) int
	FetchErrorBody(response *errorhandler.ResponseError) []byte
	BuildErrorResponse(err *types.HTTPStatusError) types.ErrorResponse
}

type httpUtilsImpl struct{}

func NewHttpUtils() HTTPUtils {
	return &httpUtilsImpl{}
}

func (hu *httpUtilsImpl) ParseErrorResponse(err error) types.HTTPResponse {
	var errorReponse errorhandler.ResponseError

	if !errors.As(err, &errorReponse) {
		return types.HTTPResponse{
			StatusCode: http.StatusInternalServerError, 
			Body: []byte(err.Error()),
		}
	}

	return types.HTTPResponse{
		StatusCode: hu.FetchErrorStatusCode(&errorReponse),
		Body: hu.FetchErrorBody(&errorReponse),
	}

}

func (hu *httpUtilsImpl) FetchErrorStatusCode(response *errorhandler.ResponseError) int {
	if response == nil {
		return http.StatusInternalServerError
	}
	return response.StatusCode
}

func (hu *httpUtilsImpl) FetchErrorBody(response *errorhandler.ResponseError) []byte {
	if response == nil {
		return nil
	}
	if response.Body == nil {
		return []byte(response.Error())
	}
	return response.Body
}

func (hu *httpUtilsImpl) BuildErrorResponse(err *types.HTTPStatusError) types.ErrorResponse {
	return types.ErrorResponse{
		Error: *err,
	}
}

