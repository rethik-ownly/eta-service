package service

import (
	"github.com/nutanalabs/eta-service/internal/eta-service/repository"
	"github.com/nutanalabs/eta-service/internal/types"
)

type Service interface {
	GetETA(restaurant_id string, lat, lon float64) (*types.GetETAResponse, error)
	InsertETA(insertEtaRequest *types.InsertETARequest) error
}

type serviceImpl struct {
	repository repository.Repository
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