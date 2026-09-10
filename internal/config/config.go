package config

import (
	"fmt"

	logger "github.com/nutanalabs/rapido-logger-go"
	"github.com/spf13/viper"
)

var AppConfig Config

type Config struct {
	AppName string
	Mongo  MongoConfig  `json:"mongo"`
	Redis  RedisConfig  `mapstructure:"redis" json:"redis"`
	Server ServerConfig `json:"server"`
	Log    LogConfig    `json:"log"`
	HTTPTransport HTTPTransportConfig `json:"http_transport"`
	ExternalServices ExternalServicesConfig `mapstructure:"externalServices"`
	Eta                 EtaConfig           `json:"eta" mapstructure:"eta"`
	EtaDefaultEstimates EtaDefaultEstimates `mapstructure:"etaDefaultEstimates"`
	Kafka               KafkaConfig         `json:"kafka" mapstructure:"kafka"`
	GeoLayer 			GeoLayerConfig 		`mapstructure:"geoLayer"`
}

func InitDefaultConfig() *Config {
	return InitConfig("application")
}

func InitConfig(configFile string) *Config {
	viper.AutomaticEnv()
	viper.SetConfigName(configFile)
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("rapido")
	viper.AddConfigPath("config")
	viper.AddConfigPath("../config/")
	viper.AddConfigPath("../../config/")
	viper.AddConfigPath("../../../config/")

	if err := viper.ReadInConfig(); err != nil {
		logger.Error(logger.Format{Message: fmt.Sprintf("Cannot read the config File: %s", err)})
		panic(err)
	}

	logger.Info(logger.Format{Message: fmt.Sprintf("Using config file: '%s'", viper.ConfigFileUsed())})

	if err := viper.Unmarshal(&AppConfig); err != nil {
		logger.Error(logger.Format{Message: fmt.Sprintf("Cannot unmarshal the config File: %s", err)})
		panic(err)
	}

	return &AppConfig
}



// Getters
func GetConfig() *Config {
	return &AppConfig
}

func (c *Config) GetAppName() string {
	return c.AppName
}