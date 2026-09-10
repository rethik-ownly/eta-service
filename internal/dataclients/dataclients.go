package dataclients

import (
	"fmt"

	"github.com/nutanalabs/eta-service/internal/config"
	"github.com/nutanalabs/eta-service/internal/constants"
	"github.com/nutanalabs/eta-service/internal/dataclients/kafka"
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
	kafkaProducerClient kafka.ProducerClient
}

func NewDataClients(config *config.Config, mongodbClient mongo.Client, kafkaProducerClient kafka.ProducerClient) DataClients {
	return &dataClientsImpl{
		config:              config,
		mongodbClient:       mongodbClient,
		kafkaProducerClient: kafkaProducerClient,
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

	if dc.kafkaProducerClient != nil {
		dc.kafkaProducerClient.Close()
		logger.Info(logger.Format{
			Event:   "CLOSE_KAFKA_PRODUCER",
			Message: "Kafka producer closed successfully",
		})
	}

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

	if dc.config.Kafka.IsProducerHealthCheckEnabled {
		if dc.kafkaProducerClient != nil {
			kafkaProducerHealthCheck := types.HealthCheck{Client: "KafkaProducer", Status: constants.UP}
			if err := dc.kafkaProducerClient.CheckHealth(); err != nil {
				healthCheckError = err
				kafkaProducerHealthCheck.Status = constants.DOWN
				logger.Error(logger.Format{
					Event:   "KAFKA_PRODUCER_HEALTH_CHECK",
					Message: fmt.Sprintf("Kafka health check failed with error: %v", err),
				})
			}
			healthCheckResponse = append(healthCheckResponse, kafkaProducerHealthCheck)
		}
	}

	// TODO : CheckHealth Redis, KafkaConsumer 

	return healthCheckResponse, healthCheckError
}

