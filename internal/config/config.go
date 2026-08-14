package config

type Config struct {
	Mongo MongoConfig `json:"mongo"`
	Server ServerConfig `json:"server"`
	Log LogConfig `json:"log"`
}