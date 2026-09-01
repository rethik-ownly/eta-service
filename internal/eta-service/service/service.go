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
	// Time , day , mealtype
	now := time.Now()
	day := s.commonUtils.GetDayFromTime(now)
	mealType := s.commonUtils.GetMealTypeFromTime(now)

	var sources []types.Location
	var restaurantsID []string
	for _, value := range request.Entities {
		sources = append(sources, value.RestaurantLocation)
		restaurantsID = append(restaurantsID, value.RestaurantID)
	}

	batchSize := s.config.Mongo.QueryBatchSize

	// Fetching restaurant_Estimates By ID's
	restaurantsEstimates, err := s.repository.FetchRestaurantEstimatesByIDs(restaurantsID, day, mealType, batchSize)

	if err != nil {
		err = fmt.Errorf("Fetching restaurant estimates failed : %w", err)
		s.repository.PublishFetchEtaEvent(constants.EventTypeNewEta, request, nil, err)
		return nil, err
	}

	estimatesByRestaurant := make(map[string]types.EtaRestaurantEstimates, len(restaurantsEstimates))
	for _, e := range restaurantsEstimates {
		estimatesByRestaurant[e.RestaurantId] = e
	}

	// Fetching sublocality_estimates by ID's
	uniqueSublocalitesID := getUniqueSublocalitiesId(restaurantsEstimates)

	sublocalitiesEstimates, err := s.repository.FetchSublocalityEstimatesByIDs(uniqueSublocalitesID, day, mealType, batchSize)

	if err != nil {
		err = fmt.Errorf("Fetching sublocality estimates failed : %w", err)
		s.repository.PublishFetchEtaEvent(constants.EventTypeNewEta, request, nil, err)
		return nil, err
	}

	estimatesBySublocality := make(map[string]types.EtaSublocalityEstimates, len(sublocalitiesEstimates))
	for _, e := range sublocalitiesEstimates {
		estimatesBySublocality[e.SublocalityId] = e
	}

	// TODO : If-else based on surface

	// Calculating distance matrix
	distanceMatrixRequest := routingengine.DistanceMatrixRequest{
		Sources:           sources,
		Destinations:      []types.Location{request.UserLocation},
		Vehicle:           constants.VEHICLE_TWO_WHEELER,
		RoutingPreference: constants.ROUTING_PREFERENCE_TRAFFIC_AWARE,
	}

	distanceMatrixResponse, err := s.getDistanceMatrix(request.Options.QosLevel, &distanceMatrixRequest)

	if err != nil {
		s.repository.PublishFetchEtaEvent(constants.EventTypeNewEta, request, nil, err)
		return nil, err
	}

	var response []types.FetchEtaResponse
	for index, value := range request.Entities {
		restaurantEstimates, ok := estimatesByRestaurant[value.RestaurantID]
		if !ok {
			continue
		}

		sublocalityEstimates, ok := estimatesBySublocality[restaurantEstimates.SublocalityId]
		if !ok {
			continue
		}

		restSection := restaurantEstimates.MealSection(mealType)
		subSection := sublocalityEstimates.MealSection(mealType)

		lastMile := distanceMatrixResponse.Data[index][0].Duration.Value

		etaInSeconds := restSection.Rat.Seconds + max(restSection.Kpt.Seconds, subSection.Cat.Seconds+subSection.Fm.Seconds+restSection.Pickup.Seconds) + lastMile

		response = append(response, types.FetchEtaResponse{
			RestaurantID: restaurantEstimates.RestaurantId,
			EtaInSeconds: uint(etaInSeconds),
			//TODO: displayMin, displayMax ( what to do if 2 min ? )
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
	overlaps, err := s.repository.SublocalityDayOverlapExists(request.SublocalityId, request.Day)
	if err != nil {
		return fmt.Errorf("checking existing sublocality estimate failed: %w", err)
	}
	if overlaps {
		return types.NewConflictError(fmt.Sprintf("sublocality estimate already exists for one or more of sublocalityId=%s day=%v", request.SublocalityId, request.Day))
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
		return nil, fmt.Errorf("distance matrix API call failed: %w", err)
	}

	if len(distanceMatrixResponse.Data) == 0 || len(distanceMatrixResponse.Data) != len(distanceMatrixRequest.Sources) {
		return nil, fmt.Errorf("invalid response from distance matrix API: expected %d sources, got %d", len(distanceMatrixRequest.Sources), len(distanceMatrixResponse.Data))
	}
	return distanceMatrixResponse, nil
}
