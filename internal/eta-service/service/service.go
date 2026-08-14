package service

import "github.com/nutanalabs/eta-service/internal/eta-service/repository"

type Service interface {
	GetETA(lat, lon float64) (float64, error)
}

type serviceImpl struct {
	repository repository.Repository
}

func NewService(repository repository.Repository) Service {
	return &serviceImpl{
		repository: repository,
	}
}

func (s *serviceImpl) GetETA(lat, lon float64) (float64, error) {
	return s.repository.GetETA(lat, lon)
}