package dataclients

import (
	"fmt"

	"github.com/nutanalabs/eta-service/internal/config"
	"github.com/nutanalabs/eta-service/internal/constants"
	"github.com/nutanalabs/eta-service/internal/dataclients/mongo"
	"github.com/nutanalabs/eta-service/internal/types"
	logger "github.com/nutanalabs/rapido-logger-go"
)

type DataClients interface {
	Stop()
	CheckHealth() ([]types.HealthCheck, error)
}

type dataClientsImpl struct {
	config *config.Config
	mongodbClient mongo.Client
}

func NewDataClients(config *config.Config, mongodbClient mongo.Client) DataClients {
	return &dataClientsImpl{
		config: config,
		mongodbClient: mongodbClient,
	}
}

func (dc *dataClientsImpl) Stop() {
	if err := dc.mongodbClient.Disconnect() ; err != nil {
		logger.Error(logger.Format{
			Event:   "CLOSE_MONGO_DB_CONNECTION",
			Message: fmt.Sprintf("error while closing mongoDB connection: %v", err),
		})
		return
	}

	logger.Info(logger.Format{
		Event:   "CLOSE_MONGO_DB_CONNECTION",
		Message: "mongoDB connection closed successfully",
	})

	// TODO : Stop Redis, kafkaConsumer
}

func (dc *dataClientsImpl) CheckHealth() ([]types.HealthCheck, error) {
	var healthCheckResponse []types.HealthCheck
	var healthCheckError error

	if dc.config.Mongo.IsHealthCheckEnabled {
		mongoDbHealthCheck := types.HealthCheck{Client: "MongoDB", Status: constants.UP}
		if err := dc.mongodbClient.CheckHealth(); err != nil {
			healthCheckError = err
			mongoDbHealthCheck.Status = constants.DOWN
		}
		healthCheckResponse = append(healthCheckResponse, mongoDbHealthCheck)
	}

	// TODO : CheckHealth Redis, KafkaConsumer 

	return healthCheckResponse, healthCheckError
}

