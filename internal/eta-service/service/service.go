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
}

type serviceImpl struct {
	repository repository.Repository
	commonUtils common.CommonUtils
	routingClient routingengine.RoutingEngineClient
	config *config.Config
}

func NewService(repository repository.Repository, 
	commonUtils common.CommonUtils, 
	routingClient routingengine.RoutingEngineClient,
	config *config.Config) Service {
	return &serviceImpl{
		repository: repository,
		commonUtils: commonUtils,
		routingClient: routingClient,
		config: config,
	}
}

func (s *serviceImpl) FetchEta(request *types.FetchEtaRequest) ([]types.FetchEtaResponse, error) {
	// Time , day , mealtype
	now := time.Now()
	day := s.commonUtils.GetDayFromTime(now)
	mealType := s.commonUtils.GetMealTypeFromTime(now)

	var sources []types.Location
	var restaurantsID []string
	for _ , value := range request.Entities {
		sources = append(sources, value.RestaurantLocation)
		restaurantsID = append(restaurantsID, value.RestaurantID)
	}

	batchSize := s.config.Mongo.QueryBatchSize

	// Fetching restaurant_Estimates By ID's
	restaurantsEstimates, err := s.repository.FetchRestaurantEstimatesByIDs(restaurantsID, day, mealType, batchSize)

	if err != nil {
		return nil, fmt.Errorf("Fetching restaurant estimates failed : %w", err)
	}

	estimatesByRestaurant := make(map[string]types.EtaRestaurantEstimates, len(restaurantsEstimates))
	for _, e := range restaurantsEstimates {
		estimatesByRestaurant[e.RestaurantID] = e
	}

	// Fetching sublocality_estimates by ID's
	uniqueSublocalitesID := getUniqueSublocalitiesId(restaurantsEstimates)

	sublocalitiesEstimates, err := s.repository.FetchSublocalityEstimatesByIDs(uniqueSublocalitesID, day, mealType, batchSize)

	if err != nil {
		return nil, fmt.Errorf("Fetching sublocality estimates failed : %w", err)
	}

	estimatesBySublocality := make(map[string]types.EtaSublocalityEstimates, len(sublocalitiesEstimates))
	for _, e := range sublocalitiesEstimates {
		estimatesBySublocality[e.SublocalityID] = e
	}

	// TODO : If-else based on surface

	// Calculating distance matrix
	distanceMatrixRequest := routingengine.DistanceMatrixRequest{
		Sources: sources,
		Destinations: []types.Location{request.UserLocation},
		Vehicle: constants.VEHICLE_TWO_WHEELER,
		RoutingPreference: constants.ROUTING_PREFERENCE_TRAFFIC_AWARE,
	}

	distanceMatrixResponse, err := s.getDistanceMatrix(request.Options.QosLevel, &distanceMatrixRequest)

	if err != nil {
		return nil, err
	}

	var response []types.FetchEtaResponse
	for index, value := range request.Entities {
		restaurantEstimates, ok := estimatesByRestaurant[value.RestaurantID]
		if !ok {
			continue
		}

		sublocalityEstimates, ok := estimatesBySublocality[restaurantEstimates.SublocalityID]
		if !ok {
			continue
		}

		lastMile := distanceMatrixResponse.Data[index][0].Duration.Value

		etaInSeconds := restaurantEstimates.RatSeconds + max(restaurantEstimates.KptSeconds, sublocalityEstimates.CatSeconds + sublocalityEstimates.FmSeconds + restaurantEstimates.PickupSeconds) + lastMile

		response = append(response, types.FetchEtaResponse{
			RestaurantID: restaurantEstimates.RestaurantID,
			EtaInSeconds: uint(etaInSeconds),
			//TODO: displayMin, displayMax
		})
	}

	return response, nil
}

func (s *serviceImpl) InsertRestaurantEstimates(request *types.InsertRestaurantEstimateRequest) error {
	request.UpdatedAt = float64(time.Now().Unix())
	return s.repository.InsertRestaurantEstimates(request)
}

func (s *serviceImpl) InsertSublocalityEstimates(request *types.InsertSublocalityEstimateRequest) error {
	request.UpdatedAt = float64(time.Now().Unix())
	return s.repository.InsertSublocalityEstimates(request)
}


// Helpers

func getUniqueSublocalitiesId(restaurantsEstimates []types.EtaRestaurantEstimates) []string {
	uniqueSublocalitiesId := make(map[string]struct{})

	for _ , value := range restaurantsEstimates {
		if _, exists := uniqueSublocalitiesId[value.SublocalityID]; exists {
			continue
		}
		uniqueSublocalitiesId[value.SublocalityID] = struct{}{}
	}

	var result []string
	
	for sublocaityId  := range uniqueSublocalitiesId {
		result = append(result, sublocaityId)
	}

	return result
}

func (s *serviceImpl) getDistanceMatrix(qosLevel constants.QosLevel, distanceMatrixRequest *routingengine.DistanceMatrixRequest) (*routingengine.DistanceMatrixResponse, error) {
	var distanceMatrixResponse *routingengine.DistanceMatrixResponse
	var err error
	if(qosLevel != "") {
		distanceMatrixResponse, err =  s.routingClient.GetDistanceMatrixWithQoS(distanceMatrixRequest, qosLevel)
	}else {
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