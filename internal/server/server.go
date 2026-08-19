package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nutanalabs/eta-service/internal/config"
	logger "github.com/nutanalabs/rapido-logger-go"
)

type Server struct {
	config *config.Config
	engine *gin.Engine
	routerGroups RouterGroups
}

type RouterGroups struct {
	rootRouter *gin.Engine
}

func NewServer(config *config.Config) *Server {
	engine := gin.New()
	engine.Use(gin.Recovery())

	return &Server{
		config: config,
		engine: engine,
		routerGroups: RouterGroups{
			rootRouter: engine,
		},
	}
}

func (s *Server) Run(h Handlers) {
	s.InitRoutes(h)

	server := &http.Server {
		Addr: s.config.GetListenAddress(),
		Handler: s.engine,
	}

	go listenAndServe(server)

	waitForShutDown(server)
}

func listenAndServe(server *http.Server) {
	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		panic(err)
	}
}

func waitForShutDown(server *http.Server) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		logger.Error(logger.Format{Message: fmt.Sprintf("Server forced to shutdown: %v", err)})
	}

	logger.Info(logger.Format{Message: "server shutdown complete"})
}



