package repository

import (
	"github.com/nutanalabs/eta-service/internal/constants"
	"github.com/nutanalabs/eta-service/internal/dataclients/mongo"
	"github.com/nutanalabs/eta-service/internal/types"
	"go.mongodb.org/mongo-driver/bson"
)

type Repository interface {
	GetETA(restaurant_id string, day_of_week string, time_slot string, lat, lon float64) (*types.GetETAResponse, error)
	// TODO : change time_slot to enum
	InsertETA(insertEtaRequest *types.InsertETARequest) error
}

type repositoryImpl struct {
	mongoRepository mongo.Repository
}

func NewRepository(mongoRepository mongo.Repository) Repository {
	return &repositoryImpl{
		mongoRepository: mongoRepository,
	}
}

func (r *repositoryImpl) GetETA(restaurant_id string, day_of_week string, time_slot string, lat, lon float64) (*types.GetETAResponse, error) {
	filter := bson.M{
		"restaurant_id": restaurant_id,
		"day_of_week": day_of_week,
		"time_slot": time_slot,
	}
	result := r.mongoRepository.FindOne(constants.ETA_COLLECTION_NAME, filter, nil)


	var response types.GetETAResponse
	if err := result.Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (r *repositoryImpl) InsertETA(insertEtaRequest *types.InsertETARequest) error {
	_ , err := r.mongoRepository.InsertOne(constants.ETA_COLLECTION_NAME, insertEtaRequest)

	return err
}