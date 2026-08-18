package etaservice

import (
	"github.com/google/wire"
	"github.com/nutanalabs/eta-service/internal/eta-service/repository"
	"github.com/nutanalabs/eta-service/internal/eta-service/service"
)

var Wireset = wire.NewSet(
	NewHandler,
	service.NewService,
	repository.NewRepository,
)