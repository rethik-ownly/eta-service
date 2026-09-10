package etaservice

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nutanalabs/eta-service/internal/constants"
	"github.com/nutanalabs/eta-service/internal/eta-service/service"
	"github.com/nutanalabs/eta-service/internal/eta-service/validator"
	"github.com/nutanalabs/eta-service/internal/metrics"
	"github.com/nutanalabs/eta-service/internal/types"
	httpUtils "github.com/nutanalabs/eta-service/internal/utils/http"
	logger "github.com/nutanalabs/rapido-logger-go"
)

type Handler struct {
	service   service.Service
	httpUtils httpUtils.HTTPUtils
	validator validator.Validator
	metrics   metrics.Metrics
}

func NewHandler(service service.Service, httpUtils httpUtils.HTTPUtils, validator validator.Validator, metrics metrics.Metrics) *Handler {
	return &Handler{
		service:   service,
		httpUtils: httpUtils,
		validator: validator,
		metrics: metrics,
	}
}

func (h *Handler) FetchEta(ctx *gin.Context) {
	startTime := time.Now()
	method := ctx.Request.Method
	route := ctx.FullPath()

	var fetchEtaRequest types.FetchEtaRequest
	if err := ctx.ShouldBindJSON(&fetchEtaRequest); err != nil {
		logger.Error(logger.Format{
			Event:   "BIND_FETCH_ETA_REQUEST",
			Message: fmt.Sprintf("error binding request with err: %v", err),
			Data: map[string]string{
				"error":  err.Error(),
				"method": method,
				"route":  route,
			},
		})
		h.recordMetrics(startTime, method, route, fetchEtaRequest.Surface, http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, h.httpUtils.BuildErrorResponse(types.NewBadRequestError(err.Error())))
		return
	}

	requestId := ctx.GetHeader("x-request-id")

	fetchEtaRequest.OrderId = ctx.GetHeader("x-order-id")
	fetchEtaRequest.RequestId = requestId

	if err := h.validator.ValidateFetchEtaRequest(&fetchEtaRequest); err != nil {
		logger.Error(logger.Format{
			Event:   "VALIDATE_FETCH_ETA_REQUEST",
			Message: fmt.Sprintf("error validating request with err: %v", err),
			Data: map[string]string{
				"error":  err.Error(),
				"method": method,
				"route":  route,
			},
		})
		h.recordMetrics(startTime, method, route, fetchEtaRequest.Surface,http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, h.httpUtils.BuildErrorResponse(types.NewBadRequestError(err.Error())))
		return
	}

	resp, err := h.service.FetchEta(&fetchEtaRequest)

	if err != nil {
		logger.Error(logger.Format{
			Event:   "FETCH_ETA_SERVICE_ERROR",
			Message: "failed to fetch eta",
			Data: map[string]string{
				"error":  err.Error(),
				"method": method,
				"route":  route,
			},
		})
		statusCode, httpErr := httpUtils.ResolveHTTPStatusError(err, "something went wrong. Try again")
		h.recordMetrics(startTime, method, route, fetchEtaRequest.Surface,statusCode)
		ctx.JSON(statusCode, h.httpUtils.BuildErrorResponse(httpErr))
		return
	}

	h.recordMetrics(startTime, method, route, fetchEtaRequest.Surface,http.StatusOK)
	ctx.JSON(http.StatusOK, h.httpUtils.BuildSuccessResponse(resp))
}

func (h *Handler) InsertRestaurantComponents(ctx *gin.Context) {
	method := ctx.Request.Method
	route := ctx.FullPath()
	var restaurantEstimates types.InsertRestaurantComponentsRequest

	if err := ctx.ShouldBindJSON(&restaurantEstimates); err != nil {
		logger.Error(logger.Format{
			Event:   "BIND_INSERT_RESTAURANT_ESTIMATE_REQUEST",
			Message: fmt.Sprintf("error binding request with err: %v", err),
			Data: map[string]string{
				"error":  err.Error(),
				"method": method,
				"route":  route,
			},
		})
		ctx.JSON(http.StatusBadRequest, h.httpUtils.BuildErrorResponse(types.NewBadRequestError(err.Error())))
		return
	}

	if err := h.validator.ValidateInsertRestaurantComponentsRequest(&restaurantEstimates); err != nil {
		logger.Error(logger.Format{
			Event:   "VALIDATE_INSERT_RESTAURANT_ESTIMATE_REQUEST",
			Message: fmt.Sprintf("error validating request with err: %v", err),
			Data: map[string]string{
				"error":  err.Error(),
				"method": method,
				"route":  route,
			},
		})
		ctx.JSON(http.StatusBadRequest, h.httpUtils.BuildErrorResponse(types.NewBadRequestError(err.Error())))
		return
	}

	err := h.service.InsertRestaurantComponents(&restaurantEstimates)

	if err != nil {
		logger.Error(logger.Format{
			Event:   "INSERT_RESTAURANT_ESTIMATE_SERVICE_ERROR",
			Message: "failed to insert restaurant estimates",
			Data: map[string]string{
				"error":  err.Error(),
				"method": method,
				"route":  route,
			},
		})
		statusCode, httpErr := httpUtils.ResolveHTTPStatusError(err, "something went wrong. Try again")
		ctx.JSON(statusCode, h.httpUtils.BuildErrorResponse(httpErr))
		return
	}

	ctx.JSON(http.StatusCreated, h.httpUtils.BuildSuccessResponse(gin.H{
		"message": "restaurant estimates inserted",
	}))
}

func (h *Handler) InsertSublocalityComponents(ctx *gin.Context) {
	method := ctx.Request.Method
	route := ctx.FullPath()
	var sublocalityEstimates types.InsertSublocalityComponentsRequest

	if err := ctx.ShouldBindJSON(&sublocalityEstimates); err != nil {
		logger.Error(logger.Format{
			Event:   "BIND_INSERT_SUBLOCALITY_ESTIMATE_REQUEST",
			Message: fmt.Sprintf("error binding request with err: %v", err),
			Data: map[string]string{
				"error":  err.Error(),
				"method": method,
				"route":  route,
			},
		})
		ctx.JSON(http.StatusBadRequest, h.httpUtils.BuildErrorResponse(types.NewBadRequestError(err.Error())))
		return
	}

	if err := h.validator.ValidateInsertSublocalityComponentsRequest(&sublocalityEstimates); err != nil {
		logger.Error(logger.Format{
			Event:   "VALIDATE_INSERT_SUBLOCALITY_ESTIMATE_REQUEST",
			Message: fmt.Sprintf("error validating request with err: %v", err),
			Data: map[string]string{
				"error":  err.Error(),
				"method": method,
				"route":  route,
			},
		})
		ctx.JSON(http.StatusBadRequest, h.httpUtils.BuildErrorResponse(types.NewBadRequestError(err.Error())))
		return
	}

	err := h.service.InsertSublocalityComponents(&sublocalityEstimates)

	if err != nil {
		logger.Error(logger.Format{
			Event:   "INSERT_SUBLOCALITY_ESTIMATE_SERVICE_ERROR",
			Message: "failed to insert sublocality estimates",
			Data: map[string]string{
				"error":  err.Error(),
				"method": method,
				"route":  route,
			},
		})
		statusCode, httpErr := httpUtils.ResolveHTTPStatusError(err, "something went wrong. Try again")
		ctx.JSON(statusCode, h.httpUtils.BuildErrorResponse(httpErr))
		return
	}

	ctx.JSON(http.StatusCreated, h.httpUtils.BuildSuccessResponse(gin.H{
		"message": "sublocality estimates inserted",
	}))
}

func (h *Handler) UpdateRestaurantComponents(ctx *gin.Context) {
	method := ctx.Request.Method
	route := ctx.FullPath()
	restaurantId := ctx.Param("restaurantId")

	var restaurantEstimates types.UpdateRestaurantComponentsRequest
	if err := ctx.ShouldBindJSON(&restaurantEstimates); err != nil {
		logger.Error(logger.Format{
			Event:   "BIND_UPDATE_RESTAURANT_ESTIMATE_REQUEST",
			Message: fmt.Sprintf("error binding request with err: %v", err),
			Data: map[string]string{
				"error":  err.Error(),
				"method": method,
				"route":  route,
			},
		})
		ctx.JSON(http.StatusBadRequest, h.httpUtils.BuildErrorResponse(types.NewBadRequestError(err.Error())))
		return
	}

	if err := h.validator.ValidateUpdateRestaurantComponents(restaurantId, &restaurantEstimates); err != nil {
		logger.Error(logger.Format{
			Event:   "VALIDATE_UPDATE_RESTAURANT_ESTIMATE_REQUEST",
			Message: fmt.Sprintf("error validating request with err: %v", err),
			Data: map[string]string{
				"error":  err.Error(),
				"method": method,
				"route":  route,
			},
		})
		ctx.JSON(http.StatusBadRequest, h.httpUtils.BuildErrorResponse(types.NewBadRequestError(err.Error())))
		return
	}

	err := h.service.UpdateRestaurantComponents(restaurantId, &restaurantEstimates)

	if err != nil {
		logger.Error(logger.Format{
			Event:   "UPDATE_RESTAURANT_ESTIMATE_SERVICE_ERROR",
			Message: "failed to update restaurant estimates",
			Data: map[string]string{
				"error":  err.Error(),
				"method": method,
				"route":  route,
			},
		})
		statusCode, httpErr := httpUtils.ResolveHTTPStatusError(err, "something went wrong. Try again")
		ctx.JSON(statusCode, h.httpUtils.BuildErrorResponse(httpErr))
		return
	}

	ctx.JSON(http.StatusOK, h.httpUtils.BuildSuccessResponse(gin.H{
		"message": "restaurant estimates updated",
	}))
}

func (h *Handler) UpdateSublocalityComponents(ctx *gin.Context) {
	method := ctx.Request.Method
	route := ctx.FullPath()
	sublocalityId := ctx.Param("sublocalityId")

	var sublocalityEstimates types.UpdateSublocalityComponentsRequest
	if err := ctx.ShouldBindJSON(&sublocalityEstimates); err != nil {
		logger.Error(logger.Format{
			Event:   "BIND_UPDATE_SUBLOCALITY_ESTIMATE_REQUEST",
			Message: fmt.Sprintf("error binding request with err: %v", err),
			Data: map[string]string{
				"error":  err.Error(),
				"method": method,
				"route":  route,
			},
		})
		ctx.JSON(http.StatusBadRequest, h.httpUtils.BuildErrorResponse(types.NewBadRequestError(err.Error())))
		return
	}

	if err := h.validator.ValidateUpdateSublocalityComponents(sublocalityId, &sublocalityEstimates); err != nil {
		logger.Error(logger.Format{
			Event:   "VALIDATE_UPDATE_SUBLOCALITY_ESTIMATE_REQUEST",
			Message: fmt.Sprintf("error validating request with err: %v", err),
			Data: map[string]string{
				"error":  err.Error(),
				"method": method,
				"route":  route,
			},
		})
		ctx.JSON(http.StatusBadRequest, h.httpUtils.BuildErrorResponse(types.NewBadRequestError(err.Error())))
		return
	}

	err := h.service.UpdateSublocalityComponents(sublocalityId, &sublocalityEstimates)

	if err != nil {
		logger.Error(logger.Format{
			Event:   "UPDATE_SUBLOCALITY_ESTIMATE_SERVICE_ERROR",
			Message: "failed to update sublocality estimates",
			Data: map[string]string{
				"error":  err.Error(),
				"method": method,
				"route":  route,
			},
		})
		statusCode, httpErr := httpUtils.ResolveHTTPStatusError(err, "something went wrong. Try again")
		ctx.JSON(statusCode, h.httpUtils.BuildErrorResponse(httpErr))
		return
	}

	ctx.JSON(http.StatusOK, h.httpUtils.BuildSuccessResponse(gin.H{
		"message": "sublocality estimates updated",
	}))
}

func (h *Handler) recordMetrics(startTime time.Time, method, route string, surface constants.Surface, status int) {
	durationMs := time.Since(startTime).Milliseconds()
	h.metrics.RecordHTTPServerRequest(method, route, surface, status, durationMs)
}