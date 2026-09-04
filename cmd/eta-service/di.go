//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/nutanalabs/eta-service/internal/config"
	"github.com/nutanalabs/eta-service/internal/dataclients"
	etaservice "github.com/nutanalabs/eta-service/internal/eta-service"
	"github.com/nutanalabs/eta-service/internal/health"
	"github.com/nutanalabs/eta-service/internal/httpclient"
	"github.com/nutanalabs/eta-service/internal/server"
	"github.com/nutanalabs/eta-service/internal/utils"
	routingengine "github.com/nutanalabs/eta-service/internal/serviceclients/routing-engine"
	"github.com/nutanalabs/eta-service/internal/metrics"
)

type ServerDependencies struct {
	config *config.Config
	server *server.Server
	handlers server.Handlers
	dataclients dataclients.DataClients
}

func InitDependencies() (ServerDependencies, error) {
	wire.Build(
		wire.Struct(new(ServerDependencies), "*"),
		wire.Struct(new(server.Handlers), "*"),
		config.GetConfig,
		server.Wireset,
		httpclient.Wireset,
		etaservice.Wireset,
		dataclients.Wireset,
		health.Wireset,
		utils.Wireset,
		routingengine.Wireset,
		metrics.Wireset,
	)

	return ServerDependencies{}, nil
}
