package server

import (
	etaservice "github.com/nutanalabs/eta-service/internal/eta-service"
	"github.com/nutanalabs/eta-service/internal/health"
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

	h.ETAHandler.InitRoutes(router)
}