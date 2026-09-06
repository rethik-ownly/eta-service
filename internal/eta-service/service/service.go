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

	sources, restaurantIDs := extractSourcesAndIDs(request.Entities)

	estimatesByRestaurant, restaurantEstimates, restaurantFailed := s.fetchRestaurantEstimates(restaurantIDs, day, mealType, batchSize)
	estimatesBySublocality, sublocalityFailed := s.fetchSublocalityEstimates(restaurantEstimates, day, mealType, batchSize, restaurantFailed)

	distanceMatrix, routingFailed := s.fetchDistanceMatrix(request, sources)

	response := make([]types.FetchEtaResponse, 0, len(request.Entities))
	for i, entity := range request.Entities {
		rest, restFallback := s.resolveRestaurantFood(entity.RestaurantID, mealType, estimatesByRestaurant, restaurantFailed)

		sublocalityID := sublocalityIDFor(entity.RestaurantID, estimatesByRestaurant, restaurantFailed)
		sub, subFallback := s.resolveSublocalityFood(sublocalityID, mealType, estimatesBySublocality, sublocalityFailed)

		lastMile := s.lastMileSeconds(request.UserLocation, entity.RestaurantLocation, distanceMatrix, routingFailed, i)
		kitchenOrDispatch := max(rest.Kpt.Seconds, sub.Cat.Seconds+sub.Fm.Seconds+rest.Pickup.Seconds+rest.DelayDispatch.Seconds)
		eta := rest.Rat.Seconds + kitchenOrDispatch + lastMile

		response = append(response, types.FetchEtaResponse{
			RestaurantID: entity.RestaurantID,
			EtaInSeconds: uint(eta),
			Source:       etaSourceFromFallback(restFallback || subFallback || routingFailed),
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

// Helpers

func getUniqueSublocalitiesId(restaurantsEstimates []types.EtaRestaurantEstimates) []string {
	uniqueSublocalitiesId := make(map[string]struct{})

	for _, value := range restaurantsEstimates {
		if _, exists := uniqueSublocalitiesId[value.SublocalityId]; exists {
			continue
		}
		uniqueSublocalitiesId[value.SublocalityId] = struct{}{}
	}

	var result []string

	for sublocaityId := range uniqueSublocalitiesId {
		result = append(result, sublocaityId)
	}

	return result
}

func (s *serviceImpl) getDistanceMatrix(qosLevel constants.QosLevel, distanceMatrixRequest *routingengine.DistanceMatrixRequest) (*routingengine.DistanceMatrixResponse, error) {
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
