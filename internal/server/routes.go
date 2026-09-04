package server

import (
	"github.com/gin-gonic/gin"
	etaservice "github.com/nutanalabs/eta-service/internal/eta-service"
	"github.com/nutanalabs/eta-service/internal/health"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Handlers struct {
	ETAHandler *etaservice.Handler
	HealthHandler *health.Handler

	// health handler
	// Other handlers
}




func (s *Server) InitRoutes(h Handlers) {
	router := s.routerGroups.rootRouter

	router.GET("/sanity", h.HealthHandler.CheckSanity)
	router.GET("/health", h.HealthHandler.CheckHealth)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	h.ETAHandler.InitRoutes(router)
}