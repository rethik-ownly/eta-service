package redis

import (
	"fmt"
	"time"

	"github.com/nutanalabs/eta-service/internal/config"
	logger "github.com/nutanalabs/rapido-logger-go"
	redislib "github.com/nutanalabs/rapido-redis-go/redis"
	redisliboptions "github.com/nutanalabs/rapido-redis-go/redis/clientoptions"
)

type Client interface {
	CheckHealth() error
	GetClient() redislib.Client
	Close() error
}

type clientImpl struct {
	config *config.Config
	client redislib.Client
}

func NewRedisClient(config *config.Config) Client {
	client, err := newRedisClient(config)
	if err != nil {
		panic(fmt.Sprintf("error connecting to redis client : %v", err))
	}
	return &clientImpl{
		config: config,
		client: client,
	}
}

func (c *clientImpl) CheckHealth() error {
	options := &redisliboptions.HealthQueryOptions{
		TimeoutInMS: &c.config.Redis.ConnectionTimeoutInMs,
	}

	return c.client.CheckHealth(options)
}

func (c *clientImpl) GetClient() redislib.Client {
	return c.client
}

func (c *clientImpl) Close() error {
	if err := c.client.Disconnect(); err != nil {
		logger.Error(logger.Format{
			Event:   "CLOSE_REDIS_CONNECTION",
			Message: fmt.Sprintf("error closing Redis connection: %v", err),
		})
		return err
	}

	logger.Info(logger.Format{
		Event:   "CLOSE_REDIS_CONNECTION",
		Message: "Redis connection closed successfully",
	})
	return nil
}

func newRedisClient(config *config.Config) (redislib.Client, error) {
	logger.Info(logger.Format{
		Event:   "NEW_REDIS_CLIENT",
		Message: "initializing redis client connection",
	})
	return redislib.Connect(getRedisClientOptions(config)) 
}

func getRedisClientOptions(config *config.Config) *redisliboptions.ClientOptions {
	metricsOptions := &redisliboptions.MetricOptions{}
	metricsOptions.SetEnabled(config.Redis.IsMetricsEnabled)

	clientOptions := &redisliboptions.ClientOptions{}

	clientOptions.
		SetApplicationName(config.GetAppName()).
		SetAddress(config.GetAddresses()).
		SetPassword(config.Redis.Password).
		SetConnectionTimeout(time.Duration(config.Redis.ConnectionTimeoutInMs) * time.Millisecond).
		SetMetrics(metricsOptions)
	return clientOptions
}