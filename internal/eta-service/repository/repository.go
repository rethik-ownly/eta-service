package repository

import (
	"github.com/nutanalabs/eta-service/internal/constants"
	"github.com/nutanalabs/eta-service/internal/dataclients/mongo"
	"github.com/nutanalabs/eta-service/internal/types"
	"github.com/nutanalabs/rapido-mongo-go/mongo/results"
	"go.mongodb.org/mongo-driver/bson"
)

type Repository interface {

	FetchRestaurantEstimates(restaurant_id string, day constants.Day, mealType constants.MealType) (*types.EtaRestaurantEstimates, error)
	FetchSublocalityEstimates(sublocality_id string, day constants.Day, mealType constants.MealType) (*types.EtaSublocalityEstimates, error)
	FetchRestaurantEstimatesByIDs(restaurantsID []string, day constants.Day, mealType constants.MealType) ([]types.EtaRestaurantEstimates, error)
	FetchSublocalityEstimatesByIDs(sublocalitiesID []string, day constants.Day, mealType constants.MealType) ([]types.EtaSublocalityEstimates, error)

	InsertRestaurantEstimates(request *types.InsertRestaurantEstimateRequest) error
	InsertSublocalityEstimates(request *types.InsertSublocalityEstimateRequest) error 
}

type repositoryImpl struct {
	mongoRepository mongo.Repository
}

func NewRepository(mongoRepository mongo.Repository) Repository {
	return &repositoryImpl{
		mongoRepository: mongoRepository,
	}
}

func (r *repositoryImpl) FetchRestaurantEstimates(restaurant_id string, day constants.Day, mealType constants.MealType) (*types.EtaRestaurantEstimates, error) {
	filter := bson.M{
		"restaurantId" : restaurant_id,
		"day" : day, 
		"mealType": mealType,
	}

	result := r.mongoRepository.FindOne(constants.ETA_RESTAURANT_ESTIMATES, filter, nil)

	var response types.EtaRestaurantEstimates
	if err := result.Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (r *repositoryImpl) FetchSublocalityEstimates(sublocality_id string, day constants.Day, mealType constants.MealType) (*types.EtaSublocalityEstimates, error) {
	filter := bson.M{
		"sublocalityId": sublocality_id,
		"day": day,
		"mealType": mealType,
	}

	result := r.mongoRepository.FindOne(constants.ETA_SUBLOCALITY_ESTIMATES, filter, nil)

	var response types.EtaSublocalityEstimates
	if err := result.Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (r *repositoryImpl) FetchRestaurantEstimatesByIDs(restaurantsId []string, day constants.Day, mealType constants.MealType) ([]types.EtaRestaurantEstimates, error) {
	filter := bson.M{
		"restaurantId": bson.M{
			"$in": restaurantsId,
		},
		"day": day,
		"mealType": mealType,
	}

	var restaurantsEstimates []types.EtaRestaurantEstimates
	queryResponse := results.QueryResponse{
		Data: &restaurantsEstimates,
	}

	err := r.mongoRepository.FindMany(constants.ETA_RESTAURANT_ESTIMATES, filter, &queryResponse, nil)

	if err != nil {
		return nil, err
	}

	if data, ok := queryResponse.Data.(*[]types.EtaRestaurantEstimates); ok && data != nil {
		return *data, nil
	}

	return []types.EtaRestaurantEstimates{}, nil
}

func (r *repositoryImpl) FetchSublocalityEstimatesByIDs(sublocalitiesId []string, day constants.Day, mealType constants.MealType) ([]types.EtaSublocalityEstimates, error) {
	filter := bson.M{
		"sublocalityId": bson.M{
			"$in": sublocalitiesId,
		},
		"day": day,
		"mealType": mealType,
	}

	var sublocalitiesEstimates []types.EtaSublocalityEstimates
	queryResponse := results.QueryResponse{
		Data: &sublocalitiesEstimates,
	}

	err := r.mongoRepository.FindMany(constants.ETA_SUBLOCALITY_ESTIMATES, filter, &queryResponse, nil)

	if err != nil {
		return nil, err
	}

	if data, ok := queryResponse.Data.(*[]types.EtaSublocalityEstimates); ok && data != nil {
		return *data, nil
	}


	return []types.EtaSublocalityEstimates{}, nil
}

func (r *repositoryImpl) InsertRestaurantEstimates(request *types.InsertRestaurantEstimateRequest) error {
	_, err := r.mongoRepository.InsertOne(constants.ETA_RESTAURANT_ESTIMATES, request)
	if err != nil {
		return err
	}
	return nil
}

func (r *repositoryImpl) InsertSublocalityEstimates(request *types.InsertSublocalityEstimateRequest) error {
	_, err := r.mongoRepository.InsertOne(constants.ETA_SUBLOCALITY_ESTIMATES, request)
	if err != nil {
		return err
	}
	return nil
}