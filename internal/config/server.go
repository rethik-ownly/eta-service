package config

import "strconv"

type ServerConfig struct {
	Host string `json:"host"`
	Port int `json:"port"`
}

func (c *Config) GetListenAddress() string {
	return ":" + strconv.Itoa(c.Server.Port)
}