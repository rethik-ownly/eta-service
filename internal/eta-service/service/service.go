package service

import (
	"time"

	"github.com/nutanalabs/eta-service/internal/eta-service/repository"
	"github.com/nutanalabs/eta-service/internal/types"
	"github.com/nutanalabs/eta-service/internal/utils/common"
)

type Service interface {
	GetETA(restaurant_id string, lat, lon float64) (*types.GetETAResponse, error)
	InsertETA(insertEtaRequest *types.InsertETARequest) error
	
	FetchEta(request *types.FetchEtaRequest) ([]types.FetchEtaResponse, error)
}

type serviceImpl struct {
	repository repository.Repository
	commonUtils common.CommonUtils
}

func NewService(repository repository.Repository) Service {
	return &serviceImpl{
		repository: repository,
	}
}

func (s *serviceImpl) GetETA(restaurant_id string, lat, lon float64) (*types.GetETAResponse, error) {
	// Logic of time to calculate day_of_week and time_slot based on the request time
	return s.repository.GetETA(restaurant_id, "monday", "lunch", lat, lon)
}

func (s *serviceImpl) InsertETA(insertEtaRequest *types.InsertETARequest) error {
	return s.repository.InsertETA(insertEtaRequest)
}

func (s *serviceImpl) FetchEta(request *types.FetchEtaRequest) ([]types.FetchEtaResponse, error) {
	// Time , day , mealtype
	now := time.Now()
	day := s.commonUtils.GetDayFromTime(now)
	mealType := s.commonUtils.GetMealTypeFromTime(now)

	var response []types.FetchEtaResponse

	for _ , value := range request.Entities {
		restaurantEstimates, err := s.repository.FetchRestaurantEstimates(value.RestaurantID, day, mealType)

		if err != nil {
			return nil, err
		}

		sublocalityEstimates, err := s.repository.FetchSublocalityEstimates(restaurantEstimates.SublocalityID, day, mealType)

		if err != nil {
			return nil, err
		}

		// TODO : Last Mile
		lastMile := 0

		etaInSeconds := *restaurantEstimates.RatSeconds + max(*restaurantEstimates.KptSeconds, *sublocalityEstimates.CatSeconds + *sublocalityEstimates.FmSeconds + *restaurantEstimates.PickupSeconds) + float64(lastMile)

		response = append(response, types.FetchEtaResponse{
			RestaurantID: restaurantEstimates.RestaurantID,
			EtaInSeconds: uint(etaInSeconds),
			DisplayMin: 5,
			DisplayMax: 10,
		})
	}

	return response, nil
}