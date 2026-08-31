package etaservice

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nutanalabs/eta-service/internal/eta-service/service"
	"github.com/nutanalabs/eta-service/internal/eta-service/validator"
	"github.com/nutanalabs/eta-service/internal/types"
	httpUtils "github.com/nutanalabs/eta-service/internal/utils/http"
	logger "github.com/nutanalabs/rapido-logger-go"
)

type Handler struct {
	service service.Service
	httpUtils httpUtils.HTTPUtils
	validator validator.Validator
}

func NewHandler(service service.Service, httpUtils httpUtils.HTTPUtils, validator validator.Validator) *Handler {
	return &Handler{
		service: service,
		httpUtils: httpUtils,
		validator: validator,
	}
}

func (h *Handler) FetchEta(ctx *gin.Context) {
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
		ctx.JSON(http.StatusBadRequest, h.httpUtils.BuildErrorResponse(types.NewBadRequestError(err.Error())))
		return
	}

	requestId := ctx.GetHeader("x-request-id")
	if(requestId == "") {
		logger.Error(logger.Format{
			RequestID: requestId,
			Event:   "VALIDATE_REQUEST_ID",
			Message: "missing mandatory request id",
			Data: map[string]string{
				"method": method,
				"route":  route,
			},
		})
		ctx.JSON(http.StatusBadRequest, h.httpUtils.BuildErrorResponse(types.NewBadRequestError("missing mandatory request id")))
		return
	}

	fetchEtaRequest.OrderId = ctx.GetHeader("x-order-id")

	// Do validation 
	if err := h.validator.ValidateFetchEtaRequest(&fetchEtaRequest) ; err != nil {
		logger.Error(logger.Format{
			Event: "VALIDATE_FETCH_ETA_REQUEST",
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
		// TOdo : Resolve Error
		ctx.JSON(http.StatusInternalServerError, h.httpUtils.BuildErrorResponse(types.NewInternalServerError(err.Error())))
		return
	}

	ctx.JSON(http.StatusOK, h.httpUtils.BuildSuccessResponse(resp))
}

func (h *Handler) InsertRestaurantEstimates(ctx *gin.Context) {
	method := ctx.Request.Method
	route := ctx.FullPath()
	var restaurantEstimates types.InsertRestaurantEstimateRequest
	
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

	if err := h.validator.ValidateInsertRestaurantEstimateRequest(&restaurantEstimates) ; err != nil {
		logger.Error(logger.Format{
			Event: "VALIDATE_INSERT_RESTAURANT_ESTIMATE_REQUEST",
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

	err := h.service.InsertRestaurantEstimates(&restaurantEstimates)

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
		statusErr := types.ToHTTPStatusError(err)
		ctx.JSON(statusErr.StatusCode(), h.httpUtils.BuildErrorResponse(statusErr))
		return
	}

	ctx.JSON(http.StatusCreated, h.httpUtils.BuildSuccessResponse(gin.H{
		"message": "restaurant estimates inserted",
	}))
}

func (h *Handler) InsertSublocalityEstimates(ctx *gin.Context) {
	method := ctx.Request.Method
	route := ctx.FullPath()
	var sublocalityEstimates types.InsertSublocalityEstimateRequest
	
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

	if err := h.validator.ValidateInsertSublocalityEstimateRequest(&sublocalityEstimates) ; err != nil {
		logger.Error(logger.Format{
			Event: "VALIDATE_INSERT_SUBLOCALITY_ESTIMATE_REQUEST",
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

	err := h.service.InsertSublocalityEstimates(&sublocalityEstimates)

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
		statusErr := types.ToHTTPStatusError(err)
		ctx.JSON(statusErr.StatusCode(), h.httpUtils.BuildErrorResponse(statusErr))
		return
	}

	ctx.JSON(http.StatusCreated, h.httpUtils.BuildSuccessResponse(gin.H{
		"message": "sublocality estimates inserted",
	}))
}

func (h *Handler) UpdateRestaurantEstimates(ctx *gin.Context) {
	method := ctx.Request.Method
	route := ctx.FullPath()
	restaurantId := ctx.Param("restaurantId")

	var restaurantEstimates types.UpdateRestaurantEstimateRequest
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

	if err := h.validator.ValidateUpdateRestaurantEstimates(restaurantId, &restaurantEstimates); err != nil {
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

	err := h.service.UpdateRestaurantEstimates(restaurantId, &restaurantEstimates)

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
		// TOdo : Resolve Error
		ctx.JSON(http.StatusInternalServerError, h.httpUtils.BuildErrorResponse(types.NewInternalServerError(err.Error())))
		return
	}

	ctx.JSON(http.StatusOK, h.httpUtils.BuildSuccessResponse(gin.H{
		"message": "restaurant estimates updated",
	}))
}

func (h *Handler) UpdateSublocalityEstimates(ctx *gin.Context) {
	method := ctx.Request.Method
	route := ctx.FullPath()
	sublocalityId := ctx.Param("sublocalityId")

	var sublocalityEstimates types.UpdateSublocalityEstimateRequest
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

	if err := h.validator.ValidateUpdateSublocalityEstimates(sublocalityId, &sublocalityEstimates); err != nil {
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

	err := h.service.UpdateSublocalityEstimates(sublocalityId, &sublocalityEstimates)

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
		// TOdo : Resolve Error
		ctx.JSON(http.StatusInternalServerError, h.httpUtils.BuildErrorResponse(types.NewInternalServerError(err.Error())))
		return
	}

	ctx.JSON(http.StatusOK, h.httpUtils.BuildSuccessResponse(gin.H{
		"message": "sublocality estimates updated",
	}))
}