package dataclients

import (
	"github.com/google/wire"
	"github.com/nutanalabs/eta-service/internal/dataclients/mongo"
)

var Wireset = wire.NewSet(
	NewDataClients,
	mongo.NewMongoClient,
	mongo.NewMongoRepository,
	
)