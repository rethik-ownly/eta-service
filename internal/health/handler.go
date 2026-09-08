package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nutanalabs/eta-service/internal/constants"
	"github.com/nutanalabs/eta-service/internal/dataclients"
	"github.com/nutanalabs/eta-service/internal/types"
)

type Handler struct {
	dataclients dataclients.DataClients
}

type SanityResponse struct {
	Client string `json:"client,omitempty"`
	Status string `json:"status"`
}

func NewHandler(dataclients dataclients.DataClients) *Handler {
	return &Handler{
		dataclients: dataclients,
	}
}

func (h *Handler) CheckSanity(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, SanityResponse{Status: constants.UP})
}

func (h *Handler) CheckHealth(ctx *gin.Context) {
	response, err := h.dataclients.CheckHealth()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.HealthCheckResponse{
			Checks: response,
			Status: constants.DOWN,
		})
		return
	}

	if response == nil {
		response = []types.HealthCheck{}
	}

	ctx.JSON(http.StatusOK, types.HealthCheckResponse{Checks: response, Status: constants.UP})
}