package service

import (
	"github.com/nutanalabs/eta-service/internal/constants"
	routingengine "github.com/nutanalabs/eta-service/internal/serviceclients/routing-engine"
	"github.com/nutanalabs/eta-service/internal/types"
)

func extractSourcesAndIDs(entities []types.FetchEtaRequestEntity) ([]types.Location, []string) {
	sources := make([]types.Location, len(entities))
	ids := make([]string, len(entities))
	for i, e := range entities {
		sources[i] = e.RestaurantLocation
		ids[i] = e.RestaurantID
	}
	return sources, ids
}

func (s *serviceImpl) fetchRestaurantEstimates(
	restaurantIDs []string,
	day constants.Day,
	mealType constants.MealType,
	batchSize int,
) (map[string]types.EtaRestaurantEstimates, []types.EtaRestaurantEstimates, bool) {
	estimates, err := s.repository.FetchRestaurantEstimatesByIDs(restaurantIDs, day, mealType, batchSize)
	if err != nil {
		return nil, nil, true
	}

	byID := make(map[string]types.EtaRestaurantEstimates, len(estimates))
	for _, e := range estimates {
		byID[e.RestaurantId] = e
	}
	return byID, estimates, false
}

func (s *serviceImpl) fetchSublocalityEstimates(
	restaurantEstimates []types.EtaRestaurantEstimates,
	day constants.Day,
	mealType constants.MealType,
	batchSize int,
	restaurantFailed bool,
) (map[string]types.EtaSublocalityEstimates, bool) {
	if restaurantFailed {
		return nil, true
	}

	sublocalityIDs := getUniqueSublocalitiesId(restaurantEstimates)
	estimates, err := s.repository.FetchSublocalityEstimatesByIDs(sublocalityIDs, day, mealType, batchSize)
	if err != nil {
		return nil, true
	}

	byID := make(map[string]types.EtaSublocalityEstimates, len(estimates))
	for _, e := range estimates {
		byID[e.SublocalityId] = e
	}
	return byID, false
}

func (s *serviceImpl) fetchDistanceMatrix(
	request *types.FetchEtaRequest,
	sources []types.Location,
) (*routingengine.DistanceMatrixResponse, bool) {
	req := routingengine.DistanceMatrixRequest{
		Sources:           sources,
		Destinations:      []types.Location{request.UserLocation},
		Vehicle:           constants.VEHICLE_TWO_WHEELER,
		RoutingPreference: constants.ROUTING_PREFERENCE_TRAFFIC_AWARE,
	}

	resp, err := s.getDistanceMatrix(request.Options.QosLevel, &req)
	return resp, err != nil
}

func sublocalityIDFor(
	restaurantID string,
	estimatesByRestaurant map[string]types.EtaRestaurantEstimates,
	restaurantFailed bool,
) string {
	if restaurantFailed {
		return ""
	}
	return estimatesByRestaurant[restaurantID].SublocalityId
}

func (s *serviceImpl) lastMileSeconds(
	userLocation, restaurantLocation types.Location,
	distanceMatrix *routingengine.DistanceMatrixResponse,
	routingFailed bool,
	index int,
) float64 {
	if routingFailed {
		distance := s.commonUtils.GetHaversineDistance(
			userLocation.Lat, userLocation.Lng,
			restaurantLocation.Lat, restaurantLocation.Lng,
		)
		return distance / constants.HF_SPEED
	}
	return distanceMatrix.Data[index][0].Duration.Value
}
