package config

type HTTPTransportConfig struct {
	DisableMetrics            bool
	MaxIdleConnections        int
	MaxConnectionsPerHost     int
	MaxIdleConnectionsPerHost int
}

func (c *Config) GetMaxIdleConnections() int {
	return c.HTTPTransport.MaxIdleConnections
}

func (c *Config) GetMaxConnectionsPerHost() int {
	return c.HTTPTransport.MaxConnectionsPerHost
}

func (c *Config) GetMaxIdleConnectionsPerHost() int {
	return c.HTTPTransport.MaxIdleConnectionsPerHost
}

func (c *Config) GetDisableMetricsConfig() bool {
	return c.HTTPTransport.DisableMetrics
}