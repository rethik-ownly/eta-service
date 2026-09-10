package config

type GeoCacheConfig struct {
	Layers                   []string `mapstructure:"layers"`
	RefreshIntervalInSeconds int      `mapstructure:"refreshIntervalInSeconds"`
}

type GeoLayerConfig struct {
	Host        string         `mapstructure:"host"`
	Port        int            `mapstructure:"port"`
	TimeoutInMs int            `mapstructure:"timeoutInMs"`
	Cache       GeoCacheConfig `mapstructure:"cache"`
}

