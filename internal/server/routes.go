package server

import etaservice "github.com/nutanalabs/eta-service/internal/eta-service"

type Handlers struct {
	ETAHandler *etaservice.Handler

	// health handler
	// Other handlers
}


func (s *Server) InitRoutes(h Handlers) {
	
}