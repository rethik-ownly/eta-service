package config

import "fmt"

type MongoConfig struct {
	Hosts string
	User string
	Password string
	Database string // The particular DB inside a cluster.
	AuthSource string // The DB in which username, password are checked against
	AppName string

	// Pool Config
	MaxConnecting             int
	MaxPoolSize               int
	MinPoolSize               int

	// Timeout Config 
	WriteTimeoutInMs          int
	ReadTimeoutInMs           int
	MaxIdleTimeInMs           int
	ConnectTimeoutInMs        int64

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