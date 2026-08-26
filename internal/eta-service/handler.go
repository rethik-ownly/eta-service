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


func (h *Handler) InsertETA(ctx *gin.Context) {
	var insertEtaRequest types.InsertETARequest

	if err := ctx.ShouldBindJSON(&insertEtaRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error" : "Invalid Body"})
		return
	}

	err := h.service.InsertETA(&insertEtaRequest)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error" : "Failed to insert ETA"})
		return
	}

	ctx.JSON(http.StatusCreated, "Eta inserted")
}

func (h *Handler) FetchEta(ctx *gin.Context) {
	fmt.Println("In controller")
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

	ctx.JSON(http.StatusOK, resp)
}