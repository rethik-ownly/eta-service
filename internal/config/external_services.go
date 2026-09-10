package config

import "github.com/nutanalabs/eta-service/internal/constants"

type ExternalServicesConfig struct {
	RoutingEngine RoutingEngineConfig `mapstructure:"routingEngine"`
}

type RoutingEngineConfig struct {
	Host              string                  `mapstructure:"host"`
	Port              int                     `mapstructure:"port"`
	DistanceMatrixAPI DistanceMatrixAPIConfig `mapstructure:"distanceMatrixAPI"`
}

type DistanceMatrixAPIConfig struct {
	Path        string `mapstructure:"path"`
	TimeoutInMs int    `mapstructure:"timeoutInMs"`
	QosLevel    constants.QosLevel `mapstructure:"qosLevel"`
}