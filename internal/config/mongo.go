package config

import (
	"fmt"
	"strconv"
	"time"
)

type MongoConfig struct {
	Hosts string
	User string
	Password string
	Database string // The particular DB inside a cluster.
	AuthSource string // The DB in which username, password are checked against
	AppName string

	// Pool Config
	MaxConnecting             uint64
	MaxPoolSize               uint64
	MinPoolSize               uint64

	// Timeout Config 
	WriteTimeoutInMs          int
	ReadTimeoutInMs           int
	MaxIdleTimeInMs           int
	ConnectTimeoutInMs        int

	// Replicas
	ReplicaSet string
	ReadPreference string

	// Retry
	RetryWritesEnabled        bool
	RetryReadsEnabled         bool

	// Observability 
	IsHealthCheckEnabled      bool
	MonitorMongoDriverEnabled bool
}

// TODO : Do we need Hidden Hosts ?

// utility functions
func (c *Config) GetMongoURI() string {
	// Handle case where no authentication is required
	if c.Mongo.User == "" || c.Mongo.Password == "" {
		return fmt.Sprintf(
			"mongodb://%s/%s?readPreference=secondaryPreferred&replicaSet=%s&appName=%s",
			c.Mongo.Hosts,
			c.Mongo.Database,
			c.Mongo.ReplicaSet,
			c.Mongo.AppName,
		)
	}

	// Handle case with authentication
	authSourceParam := ""
	if c.Mongo.AuthSource != "" {
		authSourceParam = fmt.Sprintf("&authSource=%s", c.Mongo.AuthSource)
	}

	return fmt.Sprintf(
		"mongodb://%s:%s@%s/%s?readPreference=secondaryPreferred&replicaSet=%s&appName=%s%s",
		c.Mongo.User,
		c.Mongo.Password,
		c.Mongo.Hosts,
		c.Mongo.Database,
		c.Mongo.ReplicaSet,
		c.Mongo.AppName,
		authSourceParam,
	)
}

// Getters
func (c *Config) GetUser() string {
	return c.Mongo.User
}

func (c *Config) GetDatabase() string {
	return c.Mongo.Database
}

func (c *Config) IsRetryReadsEnabledForMongo() bool {
	return c.Mongo.RetryReadsEnabled
}

func (c *Config) IsRetryWritesEnabledForMongo() bool {
	return c.Mongo.RetryWritesEnabled
}

func (c *Config) IsMonitorMongoDriverEnabled() bool {
	return c.Mongo.MonitorMongoDriverEnabled
}

func (c *Config) GetMaxConnecting() uint64 {
	return c.Mongo.MaxConnecting
}

func (c *Config) GetMaxPoolSize() uint64 {
	return c.Mongo.MaxPoolSize
}

func (c *Config) GetMinPoolSize() uint64 {
	return c.Mongo.MinPoolSize
}

func (c *Config) GetMongoReadTimeout() time.Duration {
	duration, err := time.ParseDuration(strconv.Itoa(c.Mongo.ReadTimeoutInMs) + "ms")
	if err != nil {
		panic(err)
	}
	return duration
}

func (c *Config) GetMongoWriteTimeout() time.Duration {
	duration, err := time.ParseDuration(strconv.Itoa(c.Mongo.WriteTimeoutInMs) + "ms")
	if err != nil {
		panic(err)
	}
	return duration
}

func (c *Config) GetMongoConnectionTimeout() time.Duration {
	duration, err := time.ParseDuration(strconv.Itoa(c.Mongo.ConnectTimeoutInMs) + "ms")
	if err != nil {
		panic(err)
	}
	return duration
}