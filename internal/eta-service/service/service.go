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
	GetETA(restaurant_id string, lat, lng float64) (*types.GetETAResponse, error)
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

func (s *serviceImpl) GetETA(restaurant_id string, lat, lng float64) (*types.GetETAResponse, error) {
	// Logic of time to calculate day_of_week and time_slot based on the request time
	return s.repository.GetETA(restaurant_id, "monday", "lunch", lat, lng)
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

	// Calculating distance matrix
	
	var sources []types.Location
	for _ , value := range request.Entities {
		sources = append(sources, value.RestaurantLocation)
	}

	// TODO : If-else based on surface
	distanceMatrixRequest := routingengine.DistanceMatrixRequest{
		Sources: sources,
		Destinations: []types.Location{request.UserLocation},
		Vehicle: constants.VEHICLE_TWO_WHEELER,
		RoutingPreference: constants.ROUTING_PREFERENCE_TRAFFIC_AWARE,
	}

	var distanceMatrixResponse *routingengine.DistanceMatrixResponse
	var err error

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

	for index , value := range request.Entities {
		restaurantEstimates, err := s.repository.FetchRestaurantEstimates(value.RestaurantID, day, mealType)

		if err != nil {
			return nil, err
		}

		sublocalityEstimates, err := s.repository.FetchSublocalityEstimates(restaurantEstimates.SublocalityID, day, mealType)

		if err != nil {
			return nil, err
		}

		// TODO : Last Mile
		lastMile := distanceMatrixResponse.Data[index][0].Duration.Value

		etaInSeconds := *restaurantEstimates.RatSeconds + max(*restaurantEstimates.KptSeconds, *sublocalityEstimates.CatSeconds + *sublocalityEstimates.FmSeconds + *restaurantEstimates.PickupSeconds) + float64(lastMile)

		response = append(response, types.FetchEtaResponse{
			RestaurantID: restaurantEstimates.RestaurantID,
			EtaInSeconds: uint(etaInSeconds),
			//TODO: displayMin, displayMax
		})
	}

	return response, nil
}