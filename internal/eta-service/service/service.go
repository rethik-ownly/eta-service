package service

import (
	"fmt"
	"time"

	logger "github.com/nutanalabs/rapido-logger-go"
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

	restaurantsEstimates, restaurantMongoErr := s.repository.FetchRestaurantEstimatesByIDs(restaurantsID, day, mealType, batchSize)
	restaurantMongoFailed := restaurantMongoErr != nil

	estimatesByRestaurant := make(map[string]types.EtaRestaurantEstimates)
	if !restaurantMongoFailed {
		for _, e := range restaurantsEstimates {
			estimatesByRestaurant[e.RestaurantId] = e
		}
	}

	var sublocalitiesEstimates []types.EtaSublocalityEstimates
	sublocalityMongoFailed := false
	if !restaurantMongoFailed {
		uniqueSublocalitesID := getUniqueSublocalitiesId(restaurantsEstimates)
		var sublocalityMongoErr error
		sublocalitiesEstimates, sublocalityMongoErr = s.repository.FetchSublocalityEstimatesByIDs(uniqueSublocalitesID, day, mealType, batchSize)
		sublocalityMongoFailed = sublocalityMongoErr != nil
	} else {
		sublocalityMongoFailed = true
	}

	estimatesBySublocality := make(map[string]types.EtaSublocalityEstimates)
	if !sublocalityMongoFailed {
		for _, e := range sublocalitiesEstimates {
			estimatesBySublocality[e.SublocalityId] = e
		}
	}

	distanceMatrixRequest := routingengine.DistanceMatrixRequest{
		Sources:           sources,
		Destinations:      []types.Location{request.UserLocation},
		Vehicle:           constants.VEHICLE_TWO_WHEELER,
		RoutingPreference: constants.ROUTING_PREFERENCE_TRAFFIC_AWARE,
	}

	distanceMatrixResponse, routingClientErr := s.getDistanceMatrix(request.Options.QosLevel, &distanceMatrixRequest)
	routingClientFailed := routingClientErr != nil

	var response []types.FetchEtaResponse
	for index, value := range request.Entities {
		restSection, restUsedFallback := s.resolveRestaurantFood(value.RestaurantID, mealType, estimatesByRestaurant, restaurantMongoFailed)

		sublocalityID := ""
		if !restaurantMongoFailed {
			if restaurantEstimates, ok := estimatesByRestaurant[value.RestaurantID]; ok {
				sublocalityID = restaurantEstimates.SublocalityId
			}
		}

		subSection, subUsedFallback := s.resolveSublocalityFood(sublocalityID, mealType, estimatesBySublocality, sublocalityMongoFailed)

		var lastMile float64 
		if routingClientFailed {
			haversineDistance := s.commonUtils.GetHaversineDistance(request.UserLocation.Lat, request.UserLocation.Lng, value.RestaurantLocation.Lat, value.RestaurantLocation.Lng)
			lastMile = haversineDistance / constants.HF_SPEED
		}else {
			lastMile = distanceMatrixResponse.Data[index][0].Duration.Value
		}
		
		etaInSeconds := restSection.Rat.Seconds + max(restSection.Kpt.Seconds, subSection.Cat.Seconds+subSection.Fm.Seconds+restSection.Pickup.Seconds+restSection.DelayDispatch.Seconds) + lastMile

		response = append(response, types.FetchEtaResponse{
			RestaurantID: value.RestaurantID,
			EtaInSeconds: uint(etaInSeconds),
			Source:       etaSourceFromFallback(restUsedFallback || subUsedFallback || routingClientFailed),
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
