package geocache

import (
	"fmt"

	"github.com/nutanalabs/geo-cache/cache"
	logger "github.com/nutanalabs/rapido-logger-go"
)

const (
	cityLayerType        = "city"
	zoneLayerType        = "zone"
	subLocalityLayerType = "sub_locality"
)

type GeoRepository interface {
	FetchCity(lat, lng float64) (string, error)
	FetchZone(lat, lng float64) (string, error)
	FetchSubLocality(lat, lng float64) (string, error)
	FetchLocationDetails(lat, lng float64) LocationDetails
	FetchCityAndZone(lat, lng float64) (cityID, zoneID string, err error)
}

type geoRepositoryImpl struct {
	geoCache cache.GeoCache
}

func NewGeoRepository(geoCache cache.GeoCache) GeoRepository {
	return &geoRepositoryImpl{
		geoCache: geoCache,
	}
}

// FetchCityAndZone returns both cityID and zoneID for the given location
// This is optimized for the common use case where both are needed
func (r *geoRepositoryImpl) FetchCityAndZone(lat, lng float64) (string, string, error) {
	if r.geoCache == nil {
		return "", "", fmt.Errorf("geo-cache is not initialized")
	}

	// Fetch zone first (most specific layer)
	zoneLayer, err := r.fetchFirstLayer(zoneLayerType, lat, lng)
	if err != nil {
		logger.Error(logger.Format{
			Event:   "FETCH_ZONE",
			Message: fmt.Sprintf("failed to get zone for lat: %f, lng: %f, error: %v", lat, lng, err),
		})
		return "", "", fmt.Errorf("failed to resolve zone")
	}

	if zoneLayer.ID == "" {
		return "", "", fmt.Errorf("no zone found for the given location")
	}

	zoneID := zoneLayer.ID

	// Try to get cityID from zone properties first
	var cityID string
	if zoneLayer.Properties != nil && zoneLayer.Properties.CityID != "" {
		cityID = zoneLayer.Properties.CityID
	} else {
		// Fallback: fetch city layer
		cityLayer, err := r.fetchFirstLayer(cityLayerType, lat, lng)
		if err != nil {
			logger.Warn(logger.Format{
				Event:   "FETCH_CITY",
				Message: fmt.Sprintf("failed to get city for lat: %f, lng: %f, using zone without city", lat, lng),
			})
			// Return zoneID even if city lookup fails
			return "", zoneID, nil
		}
		cityID = cityLayer.ID
		if cityLayer.Properties != nil && cityLayer.Properties.CityID != "" {
			cityID = cityLayer.Properties.CityID
		}
	}

	return cityID, zoneID, nil
}

// FetchLocationDetails returns comprehensive location information
func (r *geoRepositoryImpl) FetchLocationDetails(lat, lng float64) LocationDetails {
	cityLayer, _ := r.fetchFirstLayer(cityLayerType, lat, lng)
	zoneLayer, _ := r.fetchFirstLayer(zoneLayerType, lat, lng)
	localityLayer, _ := r.fetchFirstLayer(subLocalityLayerType, lat, lng)

	cityID := cityLayer.ID
	if cityLayer.Properties != nil && cityLayer.Properties.CityID != "" {
		cityID = cityLayer.Properties.CityID
	}
	if cityID == "" {
		cityID = cityLayer.ID
	}

	return LocationDetails{
		CityID:       cityID,
		CityName:     cityLayer.Name,
		ZoneID:       zoneLayer.ID,
		ZoneName:     zoneLayer.Name,
		LocalityID:   localityLayer.ID,
		LocalityName: localityLayer.Name,
	}
}

// FetchCity returns the cityID for the given location
func (r *geoRepositoryImpl) FetchCity(lat, lng float64) (string, error) {
	layer, err := r.fetchFirstLayer(cityLayerType, lat, lng)
	if err != nil {
		logger.Error(logger.Format{
			Event:   "FETCH_CITY",
			Message: fmt.Sprintf("failed to get city layers for lat: %f, lng: %f, error: %v", lat, lng, err),
		})
		return "", fmt.Errorf("failed to resolve city ID")
	}

	if layer.ID == "" {
		logger.Warn(logger.Format{
			Event:   "FETCH_CITY",
			Message: fmt.Sprintf("no city layers found for lat: %f, lng: %f", lat, lng),
		})
		return "", fmt.Errorf("no city found for the given location")
	}

	cityID := layer.ID
	if layer.Properties != nil && layer.Properties.CityID != "" {
		cityID = layer.Properties.CityID
	}
	if cityID == "" {
		logger.Error(logger.Format{
			Event:   "FETCH_CITY",
			Message: fmt.Sprintf("resolved city layer has no valid ID for lat: %f, lng: %f", lat, lng),
		})
		return "", fmt.Errorf("resolved city layer has no valid ID")
	}

	return cityID, nil
}

// FetchZone returns the zoneID for the given location
func (r *geoRepositoryImpl) FetchZone(lat, lng float64) (string, error) {
	layer, err := r.fetchFirstLayer(zoneLayerType, lat, lng)
	if err != nil {
		logger.Error(logger.Format{
			Event:   "FETCH_ZONE",
			Message: fmt.Sprintf("failed to get zone layers for lat: %f, lng: %f, error: %v", lat, lng, err),
		})
		return "", fmt.Errorf("failed to resolve zone ID")
	}

	if layer.ID == "" {
		return "", fmt.Errorf("no zone found for the given location")
	}

	return layer.ID, nil
}

// FetchSubLocality returns the sublocalityID for the given location
func (r *geoRepositoryImpl) FetchSubLocality(lat, lng float64) (string, error) {
	layer, err := r.fetchFirstLayer(subLocalityLayerType, lat, lng)
	if err != nil {
		logger.Error(logger.Format{
			Event:   "FETCH_SUB_LOCALITY",
			Message: fmt.Sprintf("failed to get sublocality layers for lat: %f, lng: %f, error: %v", lat, lng, err),
		})
		return "", fmt.Errorf("failed to resolve sublocality")
	}

	if layer.ID == "" {
		return "", fmt.Errorf("no sublocality found for the given location")
	}

	return layer.ID, nil
}

// fetchFirstLayer is a helper to get the first layer of a given type at a location
func (r *geoRepositoryImpl) fetchFirstLayer(layerType string, lat, lng float64) (layerInfo, error) {
	if r.geoCache == nil {
		return layerInfo{}, fmt.Errorf("geo-cache is not initialized")
	}

	layers, err := r.geoCache.GetGeoLayersByLatLng(layerType, lat, lng)
	if err != nil {
		return layerInfo{}, err
	}
	if len(layers) == 0 || layers[0] == nil {
		return layerInfo{}, nil
	}

	layer := layers[0]
	info := layerInfo{
		ID:   layer.ID,
		Name: layer.Name,
	}
	if layer.Properties != nil {
		info.Properties = &layerProperties{CityID: layer.Properties.CityID}
	}
	return info, nil
}

