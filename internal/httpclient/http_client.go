package httpclient

import (
	"github.com/nutanalabs/eta-service/internal/config"
	"github.com/nutanalabs/rapido-http-go/v3/httpclient"
)

func NewHTTPClient(config *config.Config) httpclient.Client {
	return httpclient.New(getHTTPClientSetting(config))
}

func getHTTPClientSetting(config *config.Config) httpclient.Setting {
	return httpclient.Setting{
		AppName: config.GetAppName(),
		MaxIdleConnections: config.GetMaxIdleConnections(),
		MaxIdleConnectionsPerHost: config.GetMaxIdleConnectionsPerHost(),
		DisableMetrics: config.GetDisableMetricsConfig(),
		CBSettings: nil,
	}
}