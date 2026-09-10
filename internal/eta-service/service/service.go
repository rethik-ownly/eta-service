package service

import (
	"fmt"
	"time"

	"github.com/nutanalabs/eta-service/internal/config"
	"github.com/nutanalabs/eta-service/internal/constants"
	"github.com/nutanalabs/eta-service/internal/eta-service/repository"
	"github.com/nutanalabs/eta-service/internal/geocache"
	routingengine "github.com/nutanalabs/eta-service/internal/serviceclients/routing-engine"
	"github.com/nutanalabs/eta-service/internal/types"
	common "github.com/nutanalabs/eta-service/internal/utils/common"
	logger "github.com/nutanalabs/rapido-logger-go"
)

//go:generate mockgen -source=./service.go -destination=./service_mock.go -package=service
type Service interface {
	FetchEta(request *types.FetchEtaRequest) ([]types.FetchEtaResponse, error)
	InsertRestaurantComponents(request *types.InsertRestaurantComponentsRequest) error
	InsertSublocalityComponents(request *types.InsertSublocalityComponentsRequest) error

	UpdateRestaurantComponents(restaurantId string, request *types.UpdateRestaurantComponentsRequest) error
	UpdateSublocalityComponents(sublocalityId string, request *types.UpdateSublocalityComponentsRequest) error
}

type serviceImpl struct {
	repository    repository.Repository
	commonUtils   common.CommonUtils
	routingClient routingengine.RoutingEngineClient
	geoRepository geocache.GeoRepository
	config        *config.Config
}

func NewService(repository repository.Repository,
	commonUtils common.CommonUtils,
	routingClient routingengine.RoutingEngineClient,
	geoRepository geocache.GeoRepository,
	config *config.Config) Service {
	return &serviceImpl{
		repository:    repository,
		commonUtils:   commonUtils,
		routingClient: routingClient,
		geoRepository: geoRepository,
		config:        config,
	}
}

func (s *serviceImpl) FetchEta(request *types.FetchEtaRequest) ([]types.FetchEtaResponse, error) {
	now := time.Now()
	day := s.commonUtils.GetDayFromTime(now)
	mealType := s.commonUtils.GetMealTypeFromTime(now)
	batchSize := s.config.Mongo.QueryBatchSize

	// Get cityID and zoneID from user location using geo-cache
	var cityID, zoneID string
	if s.geoRepository != nil {
		cID, zID, err := s.geoRepository.FetchCityAndZone(request.UserLocation.Lat, request.UserLocation.Lng)
		if err != nil {
			logger.Warn(logger.Format{
				Message: "Failed to get cityID and zoneID from geo-cache",
				Data: map[string]string{
					"lat": fmt.Sprintf("User Latitude: %f", request.UserLocation.Lat),
					"lng":   fmt.Sprintf("User Longitude: %f", request.UserLocation.Lng),
					"error": err.Error(),
				},
			})
		} else {
			cityID = cID
			zoneID = zID
			logger.Debug(logger.Format{
				Message: "Resolved user location to cityID and zoneID",
			})
		}
	}

	sources, restaurantIDs := extractRestaurantLocationsAndIDs(request.Entities)

	componentsByRestaurant, restaurantComponents, restaurantFailed := s.loadRestaurantComponents(restaurantIDs, cityID, zoneID, day, mealType, batchSize)
	componentsBySublocality, sublocalityFailed := s.loadSublocalityComponents(restaurantComponents, cityID, zoneID, day, mealType, batchSize, restaurantFailed)

	distanceMatrix, routingFailed := s.loadRoutingDistanceMatrix(request, sources)

	response := make([]types.FetchEtaResponse, 0, len(request.Entities))
	for i, entity := range request.Entities {
		rest, restFallback := s.resolveRestaurantMealComponents(entity.RestaurantID, mealType, componentsByRestaurant, restaurantFailed)

		sublocalityID := restaurantSublocalityID(entity.RestaurantID, componentsByRestaurant, restaurantFailed)
		sub, subFallback := s.resolveSublocalityMealComponents(sublocalityID, mealType, componentsBySublocality, sublocalityFailed)

		lastMile := s.calculateLastMileDurationSeconds(request.UserLocation, entity.RestaurantLocation, distanceMatrix, routingFailed, i)
		kitchenOrDispatch := max(rest.Kpt.Seconds, sub.Cat.Seconds+sub.Fm.Seconds+rest.Pickup.Seconds+rest.DelayDispatch.Seconds)
		eta := rest.Rat.Seconds + kitchenOrDispatch + lastMile

		response = append(response, types.FetchEtaResponse{
			RestaurantID: entity.RestaurantID,
			EtaInSeconds: uint(eta),
			Source:       resolveEtaSource(restFallback || subFallback || routingFailed),
		})
	}

	s.repository.PublishFetchEtaEvent(constants.EventTypeNewEta, request, response, nil)
	return response, nil
}

func (s *serviceImpl) InsertRestaurantComponents(request *types.InsertRestaurantComponentsRequest) error {
	overlaps, err := s.repository.RestaurantDayOverlapExists(request.RestaurantId, request.DayType)
	if err != nil {
		return fmt.Errorf("checking existing restaurant components failed: %w", err)
	}
	if overlaps {
		return types.NewConflictError(fmt.Sprintf("restaurant components already exist for one or more of restaurantId=%s dayType=%v", request.RestaurantId, request.DayType))
	}

	request.UpdatedAt = float64(time.Now().Unix())
	return s.repository.InsertRestaurantComponents(request)
}

func (s *serviceImpl) InsertSublocalityComponents(request *types.InsertSublocalityComponentsRequest) error {
	overlaps, err := s.repository.SublocalityDayOverlapExists(request.SublocalityId, request.DayType)
	if err != nil {
		return fmt.Errorf("checking existing sublocality components failed: %w", err)
	}
	if overlaps {
		return types.NewConflictError(fmt.Sprintf("sublocality components already exist for one or more of sublocalityId=%s dayType=%v", request.SublocalityId, request.DayType))
	}

	request.UpdatedAt = float64(time.Now().Unix())
	return s.repository.InsertSublocalityComponents(request)
}

func (s *serviceImpl) UpdateRestaurantComponents(restaurantId string, request *types.UpdateRestaurantComponentsRequest) error {
	request.UpdatedAt = float64(time.Now().Unix())
	return s.repository.UpdateRestaurantComponents(restaurantId, request)
}

func (s *serviceImpl) UpdateSublocalityComponents(sublocalityId string, request *types.UpdateSublocalityComponentsRequest) error {
	request.UpdatedAt = float64(time.Now().Unix())
	return s.repository.UpdateSublocalityComponents(sublocalityId, request)
}

// FetchEta helpers

func extractRestaurantLocationsAndIDs(entities []types.FetchEtaRequestEntity) ([]types.Location, []string) {
	sources := make([]types.Location, len(entities))
	ids := make([]string, len(entities))
	for i, e := range entities {
		sources[i] = e.RestaurantLocation
		ids[i] = e.RestaurantID
	}
	return sources, ids
}

func (s *serviceImpl) loadRestaurantComponents(
	restaurantIDs []string,
	cityID string,
	zoneID string,
	day constants.Day,
	mealType constants.MealType,
	batchSize int,
) (map[string]types.RestaurantComponents, []types.RestaurantComponents, bool) {
	components, err := s.repository.FetchRestaurantComponentsByIDs(restaurantIDs, cityID, zoneID, day, mealType, batchSize)
	if err != nil {
		return nil, nil, true
	}

	byID := make(map[string]types.RestaurantComponents, len(components))
	for _, c := range components {
		byID[c.RestaurantId] = c
	}
	return byID, components, false
}

func (s *serviceImpl) loadSublocalityComponents(
	restaurantComponents []types.RestaurantComponents,
	cityID string,
	zoneID string,
	day constants.Day,
	mealType constants.MealType,
	batchSize int,
	restaurantFailed bool,
) (map[string]types.SublocalityComponents, bool) {
	if restaurantFailed {
		return nil, true
	}

	sublocalityIDs := uniqueSublocalityIDs(restaurantComponents)
	components, err := s.repository.FetchSublocalityComponentsByIDs(sublocalityIDs, cityID, zoneID, day, mealType, batchSize)
	if err != nil {
		return nil, true
	}

	byID := make(map[string]types.SublocalityComponents, len(components))
	for _, c := range components {
		byID[c.SublocalityId] = c
	}
	return byID, false
}

func (s *serviceImpl) loadRoutingDistanceMatrix(
	request *types.FetchEtaRequest,
	sources []types.Location,
) (*routingengine.DistanceMatrixResponse, bool) {
	req := routingengine.DistanceMatrixRequest{
		Sources:           sources,
		Destinations:      []types.Location{request.UserLocation},
		Vehicle:           constants.VEHICLE_TWO_WHEELER,
		RoutingPreference: constants.ROUTING_PREFERENCE_TRAFFIC_AWARE,
	}

	resp, err := s.callDistanceMatrix(request.Options.QosLevel, &req)
	return resp, err != nil
}

func restaurantSublocalityID(
	restaurantID string,
	componentsByRestaurant map[string]types.RestaurantComponents,
	restaurantFailed bool,
) string {
	if restaurantFailed {
		return ""
	}
	return componentsByRestaurant[restaurantID].SublocalityId
}

func (s *serviceImpl) calculateLastMileDurationSeconds(
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

func (s *serviceImpl) defaultRestaurantMealComponents() types.RestaurantMealComponents {
	d := s.config.EtaDefaultEstimates
	return types.RestaurantMealComponents{
		Rat:           types.TimeSample{Seconds: float64(d.RestaurantAcceptanceTime)},
		Kpt:           types.TimeSample{Seconds: float64(d.KitchenPreparationTime)},
		Pickup:        types.TimeSample{Seconds: float64(d.PickupTime)},
		DelayDispatch: types.TimeSample{Seconds: float64(d.DelayDispatchTime)},
	}
}

func (s *serviceImpl) defaultSublocalityMealComponents() types.SublocalityMealComponents {
	d := s.config.EtaDefaultEstimates
	return types.SublocalityMealComponents{
		Cat: types.TimeSample{Seconds: float64(d.CaptainAssignmentTime)},
		Fm:  types.TimeSample{Seconds: float64(d.FirstMileTime)},
	}
}

func hasRestaurantMealSamples(meal types.RestaurantMealComponents) bool {
	return meal.Rat.SampleCount > 0 ||
		meal.Kpt.SampleCount > 0 ||
		meal.Pickup.SampleCount > 0 ||
		meal.DelayDispatch.SampleCount > 0
}

func hasSublocalityMealSamples(meal types.SublocalityMealComponents) bool {
	return meal.Cat.SampleCount > 0 || meal.Fm.SampleCount > 0
}

func (s *serviceImpl) resolveRestaurantMealComponents(
	restaurantID string,
	mealType constants.MealType,
	componentsByRestaurant map[string]types.RestaurantComponents,
	restaurantMongoFailed bool,
) (types.RestaurantMealComponents, bool) {
	defaults := s.defaultRestaurantMealComponents()

	if restaurantMongoFailed {
		return defaults, true
	}

	restaurantComponents, ok := componentsByRestaurant[restaurantID]
	if !ok {
		return defaults, true
	}

	meal := restaurantComponents.MealSection(mealType)
	if !hasRestaurantMealSamples(meal) {
		return defaults, true
	}

	return meal, false
}

func (s *serviceImpl) resolveSublocalityMealComponents(
	sublocalityID string,
	mealType constants.MealType,
	componentsBySublocality map[string]types.SublocalityComponents,
	sublocalityMongoFailed bool,
) (types.SublocalityMealComponents, bool) {
	defaults := s.defaultSublocalityMealComponents()

	if sublocalityMongoFailed {
		return defaults, true
	}

	if sublocalityID == "" {
		return defaults, true
	}

	sublocalityComponents, ok := componentsBySublocality[sublocalityID]
	if !ok {
		return defaults, true
	}

	meal := sublocalityComponents.MealSection(mealType)
	if !hasSublocalityMealSamples(meal) {
		return defaults, true
	}

	return meal, false
}

func resolveEtaSource(usedFallback bool) string {
	if usedFallback {
		return constants.EtaSourceFallback
	}
	return constants.EtaSourceHistoric
}

func uniqueSublocalityIDs(restaurantComponents []types.RestaurantComponents) []string {
	seen := make(map[string]struct{}, len(restaurantComponents))
	for _, value := range restaurantComponents {
		seen[value.SublocalityId] = struct{}{}
	}

	result := make([]string, 0, len(seen))
	for id := range seen {
		result = append(result, id)
	}
	return result
}

func (s *serviceImpl) callDistanceMatrix(qosLevel constants.QosLevel, distanceMatrixRequest *routingengine.DistanceMatrixRequest) (*routingengine.DistanceMatrixResponse, error) {
	var distanceMatrixResponse *routingengine.DistanceMatrixResponse
	var err error
	if qosLevel != "" {
		distanceMatrixResponse, err = s.routingClient.GetDistanceMatrixWithQoS(distanceMatrixRequest, qosLevel)
	} else {
		distanceMatrixResponse, err = s.routingClient.GetDistanceMatrix(distanceMatrixRequest)
	}

	if err != nil {
		logger.Warn(logger.Format{
			Event:   "CALCULATE_ETA_DISTANCE",
			Message: "Distance matrix API failed, falling back to Haversine",
			Data: map[string]string{
				"error": err.Error(),
			},
		})
		return nil, fmt.Errorf("distance matrix API call failed: %w", err)
	}

	if len(distanceMatrixResponse.Data) == 0 || len(distanceMatrixResponse.Data) != len(distanceMatrixRequest.Sources) {
		logger.Warn(logger.Format{
			Event:   "CALCULATE_ETA_DISTANCE",
			Message: "invalid response from distance matrix API, falling back to Haversine",
		})
		return nil, fmt.Errorf("invalid response from distance matrix API: expected %d sources, got %d", len(distanceMatrixRequest.Sources), len(distanceMatrixResponse.Data))
	}
	return distanceMatrixResponse, nil
}
