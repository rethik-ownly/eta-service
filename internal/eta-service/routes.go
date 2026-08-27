package etaservice

import "github.com/gin-gonic/gin"

func (h *Handler) InitRoutes(router *gin.Engine) {
	// The below request is a read operation but used POST because number of params required is high and body is suitable for it.
	router.POST("/api/v1/eta/estimate", h.FetchEta)
	router.POST("/api/v1/eta/promise", h.FetchEta)

	// Insert
	router.POST("/api/v1/eta/restaurant", h.InsertRestaurantEstimates)
	router.POST("/api/v1/eta/sublocality", h.InsertSublocalityEstimates)
}