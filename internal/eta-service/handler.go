package etaservice

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nutanalabs/eta-service/internal/eta-service/service"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetETA(ctx *gin.Context) {
	lat, lon := ctx.Query("lat"), ctx.Query("lon")

	// Validate latitude and longitude
	latFloat, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude"})
		return
	}

	lonFloat, err := strconv.ParseFloat(lon, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
		return
	}
	
	// Call servie to get ETA
	eta, err := h.service.GetETA(latFloat, lonFloat)
	
	// TODO : DTOtypes for request and response structure.
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get ETA"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"eta": eta})
}