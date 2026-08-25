package repository

import (
	"github.com/nutanalabs/eta-service/internal/constants"
	"github.com/nutanalabs/eta-service/internal/dataclients/mongo"
	"github.com/nutanalabs/eta-service/internal/types"
	"github.com/nutanalabs/rapido-mongo-go/mongo/results"
	"go.mongodb.org/mongo-driver/bson"
)

type Repository interface {
	GetETA(restaurant_id string, day_of_week string, time_slot string, lat, lng float64) (*types.GetETAResponse, error)
	// TODO : change time_slot to enum
	InsertETA(insertEtaRequest *types.InsertETARequest) error

	FetchRestaurantEstimates(restaurant_id string, day constants.Day, mealType constants.MealType) (*types.EtaRestaurantEstimates, error)
	FetchSublocalityEstimates(sublocality_id string, day constants.Day, mealType constants.MealType) (*types.EtaSublocalityEstimates, error)
	FetchRestaurantEstimatesByIDs(restaurantsID []string, day constants.Day, mealType constants.MealType) ([]types.EtaRestaurantEstimates, error)
	FetchSublocalityEstimatesByIDs(sublocalitiesID map[string]struct{}, day constants.Day, mealType constants.MealType) ([]types.EtaSublocalityEstimates, error)
}

type repositoryImpl struct {
	mongoRepository mongo.Repository
}

func NewRepository(mongoRepository mongo.Repository) Repository {
	return &repositoryImpl{
		mongoRepository: mongoRepository,
	}
}

func (r *repositoryImpl) GetETA(restaurant_id string, day_of_week string, time_slot string, lat, lng float64) (*types.GetETAResponse, error) {
	filter := bson.M{
		"restaurant_id": restaurant_id,
		"day_of_week": day_of_week,
		"time_slot": time_slot,
	}
	result := r.mongoRepository.FindOne(constants.ETA_RESTAURANT_ESTIMATES, filter, nil)


	var response types.GetETAResponse
	if err := result.Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (r *repositoryImpl) InsertETA(insertEtaRequest *types.InsertETARequest) error {
	_ , err := r.mongoRepository.InsertOne(constants.ETA_RESTAURANT_ESTIMATES, insertEtaRequest)

	return err
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

func (r *repositoryImpl) FetchSublocalityEstimatesByIDs(sublocalitiesId map[string]struct{}, day constants.Day, mealType constants.MealType) ([]types.EtaSublocalityEstimates, error) {
	ids := make([]string, 0, len(sublocalitiesId))
	for id := range sublocalitiesId {
		ids = append(ids, id)
	}

	filter := bson.M{
		"sublocalityId": bson.M{
			"$in": ids,
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