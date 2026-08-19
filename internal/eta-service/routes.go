package etaservice

import "github.com/gin-gonic/gin"

func (h *Handler) InitRoutes(router *gin.Engine) {
	router.GET("api/v1/eta", h.GetETA)
	router.POST("api/v1/eta", h.InsertETA)
}