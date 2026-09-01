package config

type KafkaConfig struct {
	// Health Checks
	IsConsumerHealthCheckEnabled bool   `mapstructure:"isConsumerHealthCheckEnabled"`
	IsProducerHealthCheckEnabled bool   `mapstructure:"isProducerHealthCheckEnabled"`

	ConsumersEnabled             bool   `mapstructure:"consumersEnabled"`
	BootstrapServers             string `mapstructure:"bootstrapServers"`

	// Consumer
	GroupId          string `mapstructure:"groupId"`
	AutoOffsetReset  string `mapstructure:"autoOffsetReset"`
	EnableAutoCommit bool   `mapstructure:"enableAutoCommit"`
	SessionTimeout   int    `mapstructure:"sessionTimeout"`

	// Producer
	LingerIntervalInMs int `mapstructure:"lingerIntervalInMs"`
	BatchSizeInBytes   int `mapstructure:"batchSizeInBytes"`
	WorkerPoolSize     int `mapstructure:"workerPoolSize"`
}

func (c *Config) IsKafkaConsumerHealthCheckEnabled() bool {
	return c.Kafka.IsConsumerHealthCheckEnabled
}

func (c *Config) IsKafkaProducerHealthCheckEnabled() bool {
	return c.Kafka.IsProducerHealthCheckEnabled
}

func (c *Config) AreKafkaConsumersEnabled() bool {
	return c.Kafka.ConsumersEnabled
}

func (c *Config) GetKafkaBootstrapServers() string {
	return c.Kafka.BootstrapServers
}

func (c *Config) GetKafkaGroupId() string {
	return c.Kafka.GroupId
}

func (c *Config) GetKafkaAutoOffsetReset() string {
	return c.Kafka.AutoOffsetReset
}

func (c *Config) IsKafkaAutoCommitEnabled() bool {
	return c.Kafka.EnableAutoCommit
}

func (c *Config) GetKafkaSessionTimeout() int {
	return c.Kafka.SessionTimeout
}

func (c *Config) GetKafkaLingerIntervalInMs() int {
	return c.Kafka.LingerIntervalInMs
}

func (c *Config) GetKafkaBatchSizeInBytes() int {
	return c.Kafka.BatchSizeInBytes
}

func (c *Config) GetKafkaWorkerPoolSize() int {
	return c.Kafka.WorkerPoolSize
}

func (c *Config) GetKafkaConfig() KafkaConfig {
	var kafkaConfig = KafkaConfig{
		IsConsumerHealthCheckEnabled: c.Kafka.IsConsumerHealthCheckEnabled,
		IsProducerHealthCheckEnabled: c.Kafka.IsProducerHealthCheckEnabled,
		ConsumersEnabled: c.Kafka.ConsumersEnabled,
		BootstrapServers:         c.Kafka.BootstrapServers,

		GroupId:                  c.Kafka.GroupId,
		AutoOffsetReset:          c.Kafka.AutoOffsetReset,
		EnableAutoCommit:         c.Kafka.EnableAutoCommit,
		SessionTimeout:           c.Kafka.SessionTimeout,
		
		LingerIntervalInMs: c.Kafka.LingerIntervalInMs,
		WorkerPoolSize: c.Kafka.WorkerPoolSize,
		BatchSizeInBytes: c.Kafka.BatchSizeInBytes,
	}

	return kafkaConfig
}