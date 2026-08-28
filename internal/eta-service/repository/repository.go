package repository

import (
	"fmt"
	"math"

	"github.com/nutanalabs/eta-service/internal/constants"
	"github.com/nutanalabs/eta-service/internal/dataclients/mongo"
	"github.com/nutanalabs/eta-service/internal/types"
	"github.com/nutanalabs/rapido-mongo-go/mongo/results"
	"go.mongodb.org/mongo-driver/bson"
)

type Repository interface {
	FetchRestaurantEstimates(restaurant_id string, day constants.Day, mealType constants.MealType) (*types.EtaRestaurantEstimates, error)
	FetchSublocalityEstimates(sublocality_id string, day constants.Day, mealType constants.MealType) (*types.EtaSublocalityEstimates, error)
	FetchRestaurantEstimatesByIDs(restaurantsID []string, day constants.Day, mealType constants.MealType, batchSize int) ([]types.EtaRestaurantEstimates, error)
	FetchSublocalityEstimatesByIDs(sublocalitiesID []string, day constants.Day, mealType constants.MealType, batchSize int) ([]types.EtaSublocalityEstimates, error)

	RestaurantEstimateExists(restaurantId string, day constants.Day, mealType constants.MealType) (bool, error)
	SublocalityEstimateExists(sublocalityId string, day constants.Day, mealType constants.MealType) (bool, error)

	InsertRestaurantEstimates(request *types.InsertRestaurantEstimateRequest) error
	InsertSublocalityEstimates(request *types.InsertSublocalityEstimateRequest) error

	UpdateRestaurantEstimates(restaurantId string, request *types.UpdateRestaurantEstimateRequest) error
	UpdateSublocalityEstimates(sublocalityId string, request *types.UpdateSublocalityEstimateRequest) error
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
		"restaurantId": restaurant_id,
		"day":          day,
		"mealType":     mealType,
	}

	result := r.mongoRepository.FindOne(constants.ETA_RESTAURANT_ESTIMATES, filter, nil)

	var estimates types.EtaRestaurantEstimates
	if err := result.Decode(&estimates); err != nil {
		return nil, err
	}
	return &estimates, nil
}

func (r *repositoryImpl) FetchSublocalityEstimates(sublocality_id string, day constants.Day, mealType constants.MealType) (*types.EtaSublocalityEstimates, error) {
	filter := bson.M{
		"sublocalityId": sublocality_id,
		"day":           day,
		"mealType":      mealType,
	}

	result := r.mongoRepository.FindOne(constants.ETA_SUBLOCALITY_ESTIMATES, filter, nil)

	var estimates types.EtaSublocalityEstimates
	if err := result.Decode(&estimates); err != nil {
		return nil, err
	}
	return &estimates, nil
}

func (r *repositoryImpl) FetchRestaurantEstimatesByIDs(restaurantsId []string, day constants.Day, mealType constants.MealType, batchSize int) ([]types.EtaRestaurantEstimates, error) {
	if batchSize <= 0 {
		return nil, fmt.Errorf("invalid mongo queryBatchSize: %d", batchSize)
	}

	return fetchInBatches(restaurantsId, batchSize, func(batch []string) ([]types.EtaRestaurantEstimates, error) {
		return r.fetchRestaurantEstimatesByIDsBatch(batch, day, mealType)
	})
}

func (r *repositoryImpl) fetchRestaurantEstimatesByIDsBatch(restaurantsId []string, day constants.Day, mealType constants.MealType) ([]types.EtaRestaurantEstimates, error) {
	filter := bson.M{
		"restaurantId": bson.M{
			"$in": restaurantsId,
		},
		"day":      day,
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

func (r *repositoryImpl) FetchSublocalityEstimatesByIDs(sublocalitiesId []string, day constants.Day, mealType constants.MealType, batchSize int) ([]types.EtaSublocalityEstimates, error) {
	if batchSize <= 0 {
		return nil, fmt.Errorf("invalid mongo queryBatchSize: %d", batchSize)
	}

	return fetchInBatches(sublocalitiesId, batchSize, func(batch []string) ([]types.EtaSublocalityEstimates, error) {
		return r.fetchSublocalityEstimatesByIDsBatch(batch, day, mealType)
	})
}

func (r *repositoryImpl) fetchSublocalityEstimatesByIDsBatch(sublocalitiesId []string, day constants.Day, mealType constants.MealType) ([]types.EtaSublocalityEstimates, error) {
	filter := bson.M{
		"sublocalityId": bson.M{
			"$in": sublocalitiesId,
		},
		"day":      day,
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

func (r *repositoryImpl) RestaurantEstimateExists(restaurantId string, day constants.Day, mealType constants.MealType) (bool, error) {
	filter := bson.M{
		"restaurantId": restaurantId,
		"day":          day,
		"mealType":     mealType,
	}

	count, err := r.mongoRepository.CountDocuments(constants.ETA_RESTAURANT_ESTIMATES, filter, nil)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *repositoryImpl) SublocalityEstimateExists(sublocalityId string, day constants.Day, mealType constants.MealType) (bool, error) {
	filter := bson.M{
		"sublocalityId": sublocalityId,
		"day":           day,
		"mealType":      mealType,
	}

	count, err := r.mongoRepository.CountDocuments(constants.ETA_SUBLOCALITY_ESTIMATES, filter, nil)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *repositoryImpl) InsertRestaurantEstimates(request *types.InsertRestaurantEstimateRequest) error {
	estimates := restaurantInsertRequestToEstimates(request)
	_, err := r.mongoRepository.InsertOne(constants.ETA_RESTAURANT_ESTIMATES, estimates)
	if err != nil {
		return err
	}
	return nil
}

func (r *repositoryImpl) InsertSublocalityEstimates(request *types.InsertSublocalityEstimateRequest) error {
	estimates := sublocalityInsertRequestToEstimates(request)
	_, err := r.mongoRepository.InsertOne(constants.ETA_SUBLOCALITY_ESTIMATES, estimates)
	if err != nil {
		return err
	}
	return nil
}

func (r *repositoryImpl) UpdateRestaurantEstimates(restaurantId string, request *types.UpdateRestaurantEstimateRequest) error {
	filter := bson.M{
		"restaurantId": restaurantId,
		"day":          request.Day,
		"mealType":     request.MealType,
	}

	set := bson.M{
		"updatedAt": request.UpdatedAt,
	}
	if request.RatSeconds != nil {
		set["ratSeconds"] = *request.RatSeconds
	}
	if request.RatSampleCount != nil {
		set["ratSampleCount"] = *request.RatSampleCount
	}
	if request.KptSeconds != nil {
		set["kptSeconds"] = *request.KptSeconds
	}
	if request.KptSampleCount != nil {
		set["kptSampleCount"] = *request.KptSampleCount
	}
	if request.PickupSeconds != nil {
		set["pickupSeconds"] = *request.PickupSeconds
	}
	if request.PickupSampleCount != nil {
		set["pickupSampleCount"] = *request.PickupSampleCount
	}
	if request.CityID != nil {
		set["cityId"] = *request.CityID
	}
	if request.ZoneID != nil {
		set["zoneId"] = *request.ZoneID
	}
	if request.SublocalityID != nil {
		set["sublocalityId"] = *request.SublocalityID
	}

	_, err := r.mongoRepository.UpdateOne(constants.ETA_RESTAURANT_ESTIMATES, filter, bson.M{"$set": set}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (r *repositoryImpl) UpdateSublocalityEstimates(sublocalityId string, request *types.UpdateSublocalityEstimateRequest) error {
	filter := bson.M{
		"sublocalityId": sublocalityId,
		"day":           request.Day,
		"mealType":      request.MealType,
	}

	set := bson.M{
		"updatedAt": request.UpdatedAt,
	}
	if request.ZoneID != nil {
		set["zoneId"] = *request.ZoneID
	}
	if request.CityID != nil {
		set["cityId"] = *request.CityID
	}
	if request.CatSeconds != nil {
		set["catSeconds"] = *request.CatSeconds
	}
	if request.CatSampleCount != nil {
		set["catSampleCount"] = *request.CatSampleCount
	}
	if request.FmSeconds != nil {
		set["fmSeconds"] = *request.FmSeconds
	}
	if request.FmSampleCount != nil {
		set["fmSampleCount"] = *request.FmSampleCount
	}

	_, err := r.mongoRepository.UpdateOne(constants.ETA_SUBLOCALITY_ESTIMATES, filter, bson.M{"$set": set}, nil)
	if err != nil {
		return err
	}
	return nil
}

// Helpers

// restaurantInsertRequestToEstimates maps an insert request to the
// persisted/domain estimates type.
func restaurantInsertRequestToEstimates(request *types.InsertRestaurantEstimateRequest) *types.EtaRestaurantEstimates {
	return &types.EtaRestaurantEstimates{
		RestaurantID:      request.RestaurantID,
		RatSeconds:        request.RatSeconds,
		RatSampleCount:    request.RatSampleCount,
		KptSeconds:        request.KptSeconds,
		KptSampleCount:    request.KptSampleCount,
		PickupSeconds:     request.PickupSeconds,
		PickupSampleCount: request.PickupSampleCount,
		MealType:          request.MealType,
		Day:               request.Day,
		CityID:            request.CityID,
		ZoneID:            request.ZoneID,
		SublocalityID:     request.SublocalityID,
		UpdatedAt:         request.UpdatedAt,
	}
}

// sublocalityInsertRequestToEstimates maps an insert request to the
// persisted/domain estimates type.
func sublocalityInsertRequestToEstimates(request *types.InsertSublocalityEstimateRequest) *types.EtaSublocalityEstimates {
	return &types.EtaSublocalityEstimates{
		SublocalityID:  request.SublocalityID,
		ZoneID:         request.ZoneID,
		CityID:         request.CityID,
		MealType:       request.MealType,
		Day:            request.Day,
		CatSeconds:     request.CatSeconds,
		CatSampleCount: request.CatSampleCount,
		FmSeconds:      request.FmSeconds,
		FmSampleCount:  request.FmSampleCount,
		UpdatedAt:      request.UpdatedAt,
	}
}

func fetchInBatches[T any](
	ids []string,
	batchSize int,
	fetchFn func(batch []string) ([]T, error),
) ([]T, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	numberOfBatches := int(math.Ceil(float64(len(ids)) / float64(batchSize)))

	result := make([]T, 0, len(ids))

	for i := 0; i < numberOfBatches; i++ {
		start := i * batchSize

		end := start + batchSize
		if end > len(ids) {
			end = len(ids)
		}

		batchResult, err := fetchFn(ids[start:end])
		if err != nil {
			return nil, err
		}

		result = append(result, batchResult...)
	}

	return result, nil
}
