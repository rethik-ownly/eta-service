package etaservice

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nutanalabs/eta-service/internal/eta-service/service"
	"github.com/nutanalabs/eta-service/internal/types"
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
	restaurant_id := ctx.Query("restaurant_id")
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
	etaResponse , err := h.service.GetETA(restaurant_id, latFloat, lonFloat)
	
	if err != nil {
		// error response : TODO 
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get ETA"})
		return
	}

	ctx.JSON(http.StatusOK, etaResponse)
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
	
}