package redis

import (
	"errors"
	"fmt"

	"github.com/nutanalabs/eta-service/internal/config"
	common "github.com/nutanalabs/eta-service/internal/utils/common"
	logger "github.com/nutanalabs/rapido-logger-go"
	redislib "github.com/nutanalabs/rapido-redis-go/redis"
	redisliboptions "github.com/nutanalabs/rapido-redis-go/redis/clientoptions"
	"github.com/redis/go-redis/v9"
)

type Repository interface {
	HSet(key string, value interface{}) error 
	HGet(key string, field string) (string, error)
	HDel(key string, field string) error
}

type repositoryImpl struct {
	config *config.Config
	redis  redislib.Client
	commonUtils common.CommonUtils
}

func NewRedisRepository(config *config.Config, redisClient Client) Repository {
	return &repositoryImpl{
		config: config,
		redis: redisClient.GetClient(),
	}
}

func (r *repositoryImpl) HSet(key string, value interface{}) error {
	queryTimeoutInMs := r.config.Redis.QueryTimeoutInMs
	hSetQueryOptions := redisliboptions.HSetQueryOptions{
		Key: &key,
		Value: value,
		TimeoutInMS: &queryTimeoutInMs,
	}
	_, err := r.redis.HSet(&hSetQueryOptions)
	if err != nil {
		logger.Error(logger.Format{
			Event: "REDIS_HSET",
			Message: fmt.Sprintf("failed to set redis, for key %s, value %s, with error - %v",
				key, r.commonUtils.ToJSON(value), err),
		})
		return err
	}
	logger.Debug(logger.Format{
		Event:   "REDIS_SET",
		Message: fmt.Sprintf("set redis for key %s, with value - %s", key, r.commonUtils.ToJSON(value)),
	})
	return nil
}

func (r *repositoryImpl) HGet(key string, field string) (string, error) {
	queryTimeoutInMs := r.config.Redis.QueryTimeoutInMs
	hGetQueryOptions := redisliboptions.HGetQueryOptions{
		Key: &key,
		Field: &field,
		TimeoutInMS: &queryTimeoutInMs,
	}
	value, err := r.redis.HGet(&hGetQueryOptions)
	if err != nil {
		logger.Error(logger.Format{
			Event: "REDIS_HGET",
			Message: fmt.Sprintf("failed to get redis, for key %s with error - %v",
				key, err),
		})
		if err == redis.Nil {
			return "", errors.New("no key found")
		}
		return "", err
	}
	return value, nil
}

func (r *repositoryImpl) HDel(key string, field string) error {
	queryTimeoutInMs := r.config.Redis.QueryTimeoutInMs
	hDeleteQueryOptions := redisliboptions.HDelQueryOptions{
		Key: &key,
		Field: &field,
		TimeoutInMS: &queryTimeoutInMs,
	}
	err := r.redis.HDel(&hDeleteQueryOptions)
	if err != nil {
		logger.Error(logger.Format{
			Event: "REDIS_SET",
			Message: fmt.Sprintf("failed to delete redis, for key %s, with error - %v",
				key, err),
		})
	}
	return err
}


