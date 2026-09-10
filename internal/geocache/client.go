package geocache

import (
	"time"

	"github.com/nutanalabs/geo-cache/builder"
	"github.com/nutanalabs/geo-cache/cache"
	geocfg "github.com/nutanalabs/geo-cache/config"
	logger "github.com/nutanalabs/rapido-logger-go"

	"github.com/nutanalabs/eta-service/internal/config"
)

func NewGeocache(cfg *config.Config) cache.GeoCache {
	gl := cfg.GeoLayer

	if gl.Host == "" || gl.Port <= 0 || len(gl.Cache.Layers) == 0 || gl.TimeoutInMs <= 0 {
		logger.Warn(logger.Format{
			Event:   "GEO_CACHE_INIT",
			Message: "geocache config is not provided or incomplete, skipping initialization",
		})
		return nil
	}

	b := builder.NewBuilder().WithGeoLayer(geocfg.GeoLayerConfig{
		Host:    gl.Host,
		Port:    gl.Port,
		Layers:  gl.Cache.Layers,
		Timeout: geocfg.TimeoutConfig{Default: gl.TimeoutInMs},
	})

	if gl.Cache.RefreshIntervalInSeconds > 0 {
		b = b.WithRefreshInterval(time.Duration(gl.Cache.RefreshIntervalInSeconds) * time.Second)
	}

	geoCache, err := b.Build()
	if err != nil {
		logger.Error(logger.Format{
			Event:   "GEO_CACHE_INIT",
			Message: "failed to build geo cache: " + err.Error(),
		})
		return nil
	}

	logger.Info(logger.Format{
		Event:   "GEO_CACHE_INIT",
		Message: "geo-cache initialized successfully",
	})

	return geoCache
}
