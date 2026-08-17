package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/nutanalabs/eta-service/internal/config"
	logger "github.com/nutanalabs/rapido-logger-go"
	"github.com/nutanalabs/rapido-mongo-go/mongo"
	"github.com/nutanalabs/rapido-mongo-go/mongo/options"
)

type Client interface {
	CheckHealth() error
	GetDatabaseClient() *mongo.Database
	Disconnect() error
}

type clientImpl struct {
	config *config.Config
	client *mongo.Client
}

func NewMongoClient(config *config.Config) Client {
	logger.Info(logger.Format{
		Event:   "NEW_MONGO_CLIENT",
		Message: "initializing Mongo Connection",
	})


	clientOptions := options.Client().
		SetAppName(config.GetAppName()).
		SetMaxConnecting(config.GetMaxConnecting()).
		SetMaxPoolSize(config.GetMaxPoolSize()).
		SetMinPoolSize(config.GetMinPoolSize()).
		SetRetryReads(config.IsRetryReadsEnabledForMongo()).
		SetTimeout(config.GetMongoReadTimeout()).
		SetConnectTimeout(config.GetMongoConnectionTimeout()).
		SetURI(config.GetMongoURI())
	
	if !config.IsMonitorMongoDriverEnabled() {
		clientOptions.DisableConnectionPoolMetrics().DisableCommandMetrics()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Error(logger.Format{
			Event:   "NEW_MONGO_CLIENT",
			Message: fmt.Sprintf("Error creating supply mongo client: %v", err),
		})
		panic(err)
	}

	return &clientImpl{
		client: client,
		config: config,
	}
}

func (c *clientImpl) CheckHealth() error {
	if err := c.client.CheckHealth(c.config.GetMongoReadTimeout()); err != nil {
		logger.Error(logger.Format{
			Event:   "MONGO_DB_HEALTH_CHECK",
			Message: fmt.Sprintf("mongoDB health check failed with error: %v", err),
		})
		return err
	}
	return nil
}

func (c *clientImpl) GetDatabaseClient() *mongo.Database {
	return c.client.Database(c.config.GetDatabase())
}

func (c *clientImpl) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), c.config.GetMongoReadTimeout())
	defer cancel()
	return c.client.Disconnect(ctx)
}



