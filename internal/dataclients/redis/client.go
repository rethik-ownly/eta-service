package redis

import (
	"context"
	"fmt"
	"time"
	"github.com/nutanalabs/eta-service/internal/config"
	logger "github.com/nutanalabs/rapido-logger-go"
	"github.com/redis/go-redis/v9"
)

type Client interface {
	CheckHealth() error
	GetClient() *redis.ClusterClient
	Close() error
}

type clientImpl struct {
	config *config.Config
	clusterClient *redis.ClusterClient
}

func NewRedisClient(config *config.Config) Client {
	err := config.Redis.ParseHosts()

	if err != nil {
		logger.Error(logger.Format{
			Event: "ERROR_PARSING_REDIS_HOSTS",
			Message: fmt.Sprintf("error parsing redis client : %v", err),
		})
	}

	clusterClient, err := newRedisClusterClient(config.Redis)
	if err != nil {
		logger.Error(logger.Format{
			Event:   "REDIS_CLIENT_CONNECTION_TIMEOUT",
			Message: fmt.Sprintf("error connecting to redis client : %v", err),
		})
	}

	return &clientImpl{
		config: config,
		clusterClient: clusterClient,
	}
}

func newRedisClusterClient(redisConfig config.RedisConfig) (*redis.ClusterClient, error) {
	logEventName := "redis.client.newRedisClusterClient"

	logger.Info(logger.Format{
		Event:  logEventName,
		Message: "Initializing redis connection",
	})
	clusterOptions := &redis.ClusterOptions{
		Addrs: redisConfig.Addresses,
		Password: redisConfig.Password,
		ReadTimeout: time.Duration(redisConfig.QueryTimeoutInMs),
		WriteTimeout: time.Duration(redisConfig.QueryTimeoutInMs),
	}

	rdb := redis.NewClusterClient(clusterOptions)

	pong, err := checkRedisClusterHealth(rdb, redisConfig.ConnectionTimeoutInMs)
	if err != nil {
		logger.Error(logger.Format{
			Event:   logEventName,
			Message: fmt.Sprintf("Error creating redis client: %v", err),
		})
		return rdb, err
	}

	logger.Info(logger.Format{
		Event:   logEventName,
		Message: fmt.Sprintf("Successfully connected to redis %v", pong),
	})

	return rdb, nil
} 

func (c *clientImpl) CheckHealth() error {
	_, err := checkRedisClusterHealth(c.clusterClient, c.config.Redis.ConnectionTimeoutInMs)
	if err != nil {
		return err
	}
	return nil
}

func checkRedisClusterHealth(redisClient *redis.ClusterClient, connectionTimeoutInMs int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(connectionTimeoutInMs)*time.Millisecond)
	defer cancel()

	pong, err := redisClient.Ping(ctx).Result()

	if err != nil {
		logger.Error(logger.Format{
			Event:   "REDIS_CLUSTER_HEALTH_CHECK",
			Message: fmt.Sprintf("error while redis health check %v", err),
		})
		return "", err
	}

	return pong, nil
}

func (c *clientImpl) GetClient() *redis.ClusterClient {
	return c.clusterClient
}

func (c *clientImpl) Close() error {
	redisClient := c.clusterClient
	err := redisClient.Close()
	if err != nil {
		logger.Error(logger.Format{
			Message: fmt.Sprintf(
				"error while closing redis client: redisClusterClientErr=>%v",
				err,
			),
		})
		return fmt.Errorf(
			"{redisClusterClientCloseErr}: %v",
			err,
		)
	}
	return nil
}