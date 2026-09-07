package service

import (
	"fmt"
	"time"

	"github.com/nutanalabs/eta-service/internal/config"
	"github.com/nutanalabs/eta-service/internal/constants"
	"github.com/nutanalabs/eta-service/internal/eta-service/repository"
	routingengine "github.com/nutanalabs/eta-service/internal/serviceclients/routing-engine"
	"github.com/nutanalabs/eta-service/internal/types"
	common "github.com/nutanalabs/eta-service/internal/utils/common"
	logger "github.com/nutanalabs/rapido-logger-go"
)

type Service interface {
	FetchEta(request *types.FetchEtaRequest) ([]types.FetchEtaResponse, error)
	InsertRestaurantEstimates(request *types.InsertRestaurantEstimateRequest) error
	InsertSublocalityEstimates(request *types.InsertSublocalityEstimateRequest) error

	UpdateRestaurantEstimates(restaurantId string, request *types.UpdateRestaurantEstimateRequest) error
	UpdateSublocalityEstimates(sublocalityId string, request *types.UpdateSublocalityEstimateRequest) error
}

type serviceImpl struct {
	repository    repository.Repository
	commonUtils   common.CommonUtils
	routingClient routingengine.RoutingEngineClient
	config        *config.Config
}

func NewService(repository repository.Repository,
	commonUtils common.CommonUtils,
	routingClient routingengine.RoutingEngineClient,
	config *config.Config) Service {
	return &serviceImpl{
		repository:    repository,
		commonUtils:   commonUtils,
		routingClient: routingClient,
		config:        config,
	}
}

func (s *serviceImpl) FetchEta(request *types.FetchEtaRequest) ([]types.FetchEtaResponse, error) {
	now := time.Now()
	day := s.commonUtils.GetDayFromTime(now)
	mealType := s.commonUtils.GetMealTypeFromTime(now)
	batchSize := s.config.Mongo.QueryBatchSize

	sources, restaurantIDs := extractRestaurantLocationsAndIDs(request.Entities)

	estimatesByRestaurant, restaurantEstimates, restaurantFailed := s.loadRestaurantEstimates(restaurantIDs, day, mealType, batchSize)
	estimatesBySublocality, sublocalityFailed := s.loadSublocalityEstimates(restaurantEstimates, day, mealType, batchSize, restaurantFailed)

	distanceMatrix, routingFailed := s.loadRoutingDistanceMatrix(request, sources)

	response := make([]types.FetchEtaResponse, 0, len(request.Entities))
	for i, entity := range request.Entities {
		rest, restFallback := s.resolveRestaurantMealEstimate(entity.RestaurantID, mealType, estimatesByRestaurant, restaurantFailed)

		sublocalityID := restaurantSublocalityID(entity.RestaurantID, estimatesByRestaurant, restaurantFailed)
		sub, subFallback := s.resolveSublocalityMealEstimate(sublocalityID, mealType, estimatesBySublocality, sublocalityFailed)

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

func (s *serviceImpl) InsertRestaurantEstimates(request *types.InsertRestaurantEstimateRequest) error {
	overlaps, err := s.repository.RestaurantDayOverlapExists(request.RestaurantId, request.DayType)
	if err != nil {
		return fmt.Errorf("checking existing restaurant estimate failed: %w", err)
	}
	if overlaps {
		return types.NewConflictError(fmt.Sprintf("restaurant estimate already exists for one or more of restaurantId=%s dayType=%v", request.RestaurantId, request.DayType))
	}

	request.UpdatedAt = float64(time.Now().Unix())
	return s.repository.InsertRestaurantEstimates(request)
}

func (s *serviceImpl) InsertSublocalityEstimates(request *types.InsertSublocalityEstimateRequest) error {
	overlaps, err := s.repository.SublocalityDayOverlapExists(request.SublocalityId, request.DayType)
	if err != nil {
		return fmt.Errorf("checking existing sublocality estimate failed: %w", err)
	}
	if overlaps {
		return types.NewConflictError(fmt.Sprintf("sublocality estimate already exists for one or more of sublocalityId=%s dayType=%v", request.SublocalityId, request.DayType))
	}

	request.UpdatedAt = float64(time.Now().Unix())
	return s.repository.InsertSublocalityEstimates(request)
}

func (s *serviceImpl) UpdateRestaurantEstimates(restaurantId string, request *types.UpdateRestaurantEstimateRequest) error {
	request.UpdatedAt = float64(time.Now().Unix())
	return s.repository.UpdateRestaurantEstimates(restaurantId, request)
}

func (s *serviceImpl) UpdateSublocalityEstimates(sublocalityId string, request *types.UpdateSublocalityEstimateRequest) error {
	request.UpdatedAt = float64(time.Now().Unix())
	return s.repository.UpdateSublocalityEstimates(sublocalityId, request)
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

func (s *serviceImpl) loadRestaurantEstimates(
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

func (s *serviceImpl) loadSublocalityEstimates(
	restaurantEstimates []types.EtaRestaurantEstimates,
	day constants.Day,
	mealType constants.MealType,
	batchSize int,
	restaurantFailed bool,
) (map[string]types.EtaSublocalityEstimates, bool) {
	if restaurantFailed {
		return nil, true
	}

	sublocalityIDs := uniqueSublocalityIDs(restaurantEstimates)
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
	estimatesByRestaurant map[string]types.EtaRestaurantEstimates,
	restaurantFailed bool,
) string {
	if restaurantFailed {
		return ""
	}
	return estimatesByRestaurant[restaurantID].SublocalityId
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

func (s *serviceImpl) defaultRestaurantMealEstimate() types.RestaurantMealEstimate {
	d := s.config.EtaDefaultEstimates
	return types.RestaurantMealEstimate{
		Rat:           types.TimeSample{Seconds: float64(d.RestaurantAcceptanceTime)},
		Kpt:           types.TimeSample{Seconds: float64(d.KitchenPreparationTime)},
		Pickup:        types.TimeSample{Seconds: float64(d.PickupTime)},
		DelayDispatch: types.TimeSample{Seconds: float64(d.DelayDispatchTime)},
	}
}

func (s *serviceImpl) defaultSublocalityMealEstimate() types.SublocalityMealEstimate {
	d := s.config.EtaDefaultEstimates
	return types.SublocalityMealEstimate{
		Cat: types.TimeSample{Seconds: float64(d.CaptainAssignmentTime)},
		Fm:  types.TimeSample{Seconds: float64(d.FirstMileTime)},
	}
}

func hasRestaurantMealSamples(meal types.RestaurantMealEstimate) bool {
	return meal.Rat.SampleCount > 0 ||
		meal.Kpt.SampleCount > 0 ||
		meal.Pickup.SampleCount > 0 ||
		meal.DelayDispatch.SampleCount > 0
}

func hasSublocalityMealSamples(meal types.SublocalityMealEstimate) bool {
	return meal.Cat.SampleCount > 0 || meal.Fm.SampleCount > 0
}

func (s *serviceImpl) resolveRestaurantMealEstimate(
	restaurantID string,
	mealType constants.MealType,
	estimatesByRestaurant map[string]types.EtaRestaurantEstimates,
	restaurantMongoFailed bool,
) (types.RestaurantMealEstimate, bool) {
	defaults := s.defaultRestaurantMealEstimate()

	if restaurantMongoFailed {
		return defaults, true
	}

	restaurantEstimates, ok := estimatesByRestaurant[restaurantID]
	if !ok {
		return defaults, true
	}

	meal := restaurantEstimates.MealSection(mealType)
	if !hasRestaurantMealSamples(meal) {
		return defaults, true
	}

	return meal, false
}

func (s *serviceImpl) resolveSublocalityMealEstimate(
	sublocalityID string,
	mealType constants.MealType,
	estimatesBySublocality map[string]types.EtaSublocalityEstimates,
	sublocalityMongoFailed bool,
) (types.SublocalityMealEstimate, bool) {
	defaults := s.defaultSublocalityMealEstimate()

	if sublocalityMongoFailed {
		return defaults, true
	}

	if sublocalityID == "" {
		return defaults, true
	}

	sublocalityEstimates, ok := estimatesBySublocality[sublocalityID]
	if !ok {
		return defaults, true
	}

	meal := sublocalityEstimates.MealSection(mealType)
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

func uniqueSublocalityIDs(restaurantsEstimates []types.EtaRestaurantEstimates) []string {
	seen := make(map[string]struct{}, len(restaurantsEstimates))
	for _, value := range restaurantsEstimates {
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
