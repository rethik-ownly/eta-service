package config

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type RedisConfig struct {
	Hosts string
	Password string
	Addresses []string
	ConnectionTimeoutInMs int
	QueryTimeoutInMs int
}

type RedisHosts struct {
	Host string `json:"host"`
	Port int `json:"port"`
}

func (r *RedisConfig) ParseHosts() error {
	var redisParsedConfig []RedisHosts
	if err := json.Unmarshal([]byte(r.Hosts), &redisParsedConfig); err != nil {
		fmt.Println("error occurred parsing redis config: ", err)
		return err
	}

	r.Addresses = make([]string, 0)
	for _, redis := range redisParsedConfig {
		redisHost := fmt.Sprintf("%s:%d", redis.Host, redis.Port)
		r.Addresses = append(r.Addresses, redisHost)
	}
	return nil
}

func (c *Config) GetHosts() string {
	return c.Redis.Hosts
} 

func (c *Config) GetRedisConnectionTimeout() time.Duration {
	duration, err := time.ParseDuration(strconv.Itoa(c.Redis.ConnectionTimeoutInMs) + "ms")

	if err != nil {
		panic(err)
	}

	return duration
}

func (c *Config) GetRedisQueryTimeout() time.Duration {
	duration, err := time.ParseDuration(strconv.Itoa(c.Redis.QueryTimeoutInMs) + "ms")

	if err != nil {
		panic(err)
	}

	return duration
}

func (c *Config) GetAddresses() []string {
	return c.Redis.Addresses
}