package etaservice

import "github.com/gin-gonic/gin"

func (h *Handler) InitRoutes(router *gin.Engine) {
	// The below request is a read operation but used POST because number of params required is high and body is suitable for it.
	router.POST("/api/v1/eta/estimate", h.FetchEta)
	router.POST("/api/v1/eta/promise", h.FetchEta)

	// Insert
	router.POST("/api/v1/eta/restaurant", h.InsertRestaurantComponents)
	router.POST("/api/v1/eta/sublocality", h.InsertSublocalityComponents)

	// Update
	router.PATCH("/api/v1/eta/restaurant/:restaurantId", h.UpdateRestaurantComponents)
	router.PATCH("/api/v1/eta/sublocality/:sublocalityId", h.UpdateSublocalityComponents)

	// Delete
	// router.DELETE("/api/v1/eta/restaurant/:restaurantId", h.DeleteRestaurantComponents)
	// router.DELETE("/api/v1/eta/sublocality/:sublocalityId", h.DeleteSublocalityComponents)
}