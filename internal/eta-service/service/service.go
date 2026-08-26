package service

import (
	"fmt"
	"time"

	"github.com/nutanalabs/eta-service/internal/constants"
	"github.com/nutanalabs/eta-service/internal/eta-service/repository"
	routingengine "github.com/nutanalabs/eta-service/internal/serviceclients/routing-engine"
	"github.com/nutanalabs/eta-service/internal/types"
	common "github.com/nutanalabs/eta-service/internal/utils/common"
)

type Service interface {
	InsertETA(insertEtaRequest *types.InsertETARequest) error
	
	FetchEta(request *types.FetchEtaRequest) ([]types.FetchEtaResponse, error)
}

type serviceImpl struct {
	repository repository.Repository
	commonUtils common.CommonUtils
	routingClient routingengine.RoutingEngineClient
}

func NewService(repository repository.Repository, commonUtils common.CommonUtils, routingClient routingengine.RoutingEngineClient) Service {
	return &serviceImpl{
		repository: repository,
		commonUtils: commonUtils,
		routingClient: routingClient,
	}
}

func (s *serviceImpl) InsertETA(insertEtaRequest *types.InsertETARequest) error {
	return s.repository.InsertETA(insertEtaRequest)
}

func (s *serviceImpl) FetchEta(request *types.FetchEtaRequest) ([]types.FetchEtaResponse, error) {
	fmt.Println("In service")
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

	// Fetching restaurant_Estimates By ID's
	restaurantsEstimates, err := s.repository.FetchRestaurantEstimatesByIDs(restaurantsID, day, mealType)
	if err != nil {
		return nil, fmt.Errorf("Fetching restaurant estimates failed : %w", err)
	}
	estimatesByRestaurant := make(map[string]types.EtaRestaurantEstimates, len(restaurantsEstimates))
	for _, e := range restaurantsEstimates {
		estimatesByRestaurant[e.RestaurantID] = e
	}

	// Fetching sublocality_estimates by ID's
	uniqueSublocalitesID := make(map[string]struct{})
	for _ , value := range restaurantsEstimates {
		if _ , exists := uniqueSublocalitesID[value.SublocalityID]; exists {
			continue
		}
		uniqueSublocalitesID[value.SublocalityID] = struct{}{}
	}

	sublocalitiesEstimates, err := s.repository.FetchSublocalityEstimatesByIDs(uniqueSublocalitesID, day, mealType)
	if err != nil {
		return nil, fmt.Errorf("Fetching sublocalities estimates failed : %w", err)
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

	var distanceMatrixResponse *routingengine.DistanceMatrixResponse

	qosLevel := request.Options.QosLevel
	if(qosLevel != "") {
		distanceMatrixResponse, err =  s.routingClient.GetDistanceMatrixWithQoS(&distanceMatrixRequest, qosLevel)
	}else {
		distanceMatrixResponse, err = s.routingClient.GetDistanceMatrix(&distanceMatrixRequest)
	}

	if err != nil {
		return nil, fmt.Errorf("distance matrix API call failed: %w", err)
	}
	
	var response []types.FetchEtaResponse
	for index, value := range request.Entities {
		restaurantEstimates, ok := estimatesByRestaurant[value.RestaurantID]
		if !ok {
			return nil, fmt.Errorf("restaurant estimates not found for %s", value.RestaurantID)
		}

		sublocalityEstimates, ok := estimatesBySublocality[restaurantEstimates.SublocalityID]
		if !ok {
			return nil, fmt.Errorf("sublocality estimates not found for %s", restaurantEstimates.SublocalityID)
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