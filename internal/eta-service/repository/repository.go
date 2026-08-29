package repository

import (
	"fmt"
	"math"

	"github.com/nutanalabs/eta-service/internal/constants"
	"github.com/nutanalabs/eta-service/internal/dataclients/mongo"
	"github.com/nutanalabs/eta-service/internal/types"
	"github.com/nutanalabs/rapido-mongo-go/mongo/results"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	FetchRestaurantEstimates(restaurantId string, day constants.Day, mealType constants.MealType) (*types.EtaRestaurantEstimates, error)
	FetchSublocalityEstimates(sublocalityId string, day constants.Day, mealType constants.MealType) (*types.EtaSublocalityEstimates, error)
	FetchRestaurantEstimatesByIDs(restaurantsId []string, day constants.Day, mealType constants.MealType, batchSize int) ([]types.EtaRestaurantEstimates, error)
	FetchSublocalityEstimatesByIDs(sublocalitiesId []string, day constants.Day, mealType constants.MealType, batchSize int) ([]types.EtaSublocalityEstimates, error)

	RestaurantDayOverlapExists(restaurantId string, days []constants.Day) (bool, error)
	SublocalityDayOverlapExists(sublocalityId string, days []constants.Day) (bool, error)

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

func (r *repositoryImpl) FetchRestaurantEstimates(restaurantId string, day constants.Day, mealType constants.MealType) (*types.EtaRestaurantEstimates, error) {
	filter := bson.M{
		"restaurantId": restaurantId,
		"dayType":      day,
	}
	opts := options.FindOne().SetProjection(restaurantMealProjection(mealType))

	result := r.mongoRepository.FindOne(constants.ETA_RESTAURANT_ESTIMATES, filter, opts)

	var estimates types.EtaRestaurantEstimates
	if err := result.Decode(&estimates); err != nil {
		return nil, err
	}
	return &estimates, nil
}

func (r *repositoryImpl) FetchSublocalityEstimates(sublocalityId string, day constants.Day, mealType constants.MealType) (*types.EtaSublocalityEstimates, error) {
	filter := bson.M{
		"sublocalityId": sublocalityId,
		"day":           day,
	}
	opts := options.FindOne().SetProjection(sublocalityMealProjection(mealType))

	result := r.mongoRepository.FindOne(constants.ETA_SUBLOCALITY_ESTIMATES, filter, opts)

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
		"dayType": day,
	}
	opts := options.Find().SetProjection(restaurantMealProjection(mealType))

	var restaurantsEstimates []types.EtaRestaurantEstimates
	queryResponse := results.QueryResponse{
		Data: &restaurantsEstimates,
	}

	err := r.mongoRepository.FindMany(constants.ETA_RESTAURANT_ESTIMATES, filter, &queryResponse, opts)

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
		"day": day,
	}
	opts := options.Find().SetProjection(sublocalityMealProjection(mealType))

	var sublocalitiesEstimates []types.EtaSublocalityEstimates
	queryResponse := results.QueryResponse{
		Data: &sublocalitiesEstimates,
	}

	err := r.mongoRepository.FindMany(constants.ETA_SUBLOCALITY_ESTIMATES, filter, &queryResponse, opts)

	if err != nil {
		return nil, err
	}

	if data, ok := queryResponse.Data.(*[]types.EtaSublocalityEstimates); ok && data != nil {
		return *data, nil
	}

	return []types.EtaSublocalityEstimates{}, nil
}

// RestaurantDayOverlapExists reports whether any existing document for
// restaurantId already covers one or more of the given days. Used by Insert
// to keep each day covered by at most one document per restaurant.
func (r *repositoryImpl) RestaurantDayOverlapExists(restaurantId string, days []constants.Day) (bool, error) {
	filter := bson.M{
		"restaurantId": restaurantId,
		"dayType": bson.M{
			"$in": days,
		},
	}

	count, err := r.mongoRepository.CountDocuments(constants.ETA_RESTAURANT_ESTIMATES, filter, nil)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// SublocalityDayOverlapExists reports whether any existing document for
// sublocalityId already covers one or more of the given days. Used by
// Insert to keep each day covered by at most one document per sublocality.
func (r *repositoryImpl) SublocalityDayOverlapExists(sublocalityId string, days []constants.Day) (bool, error) {
	filter := bson.M{
		"sublocalityId": sublocalityId,
		"day": bson.M{
			"$in": days,
		},
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
		"dayType":      request.Day,
	}

	set := bson.M{
		"updatedAt": request.UpdatedAt,
	}
	if request.Days != nil {
		set["dayType"] = *request.Days
	}
	if request.CityId != nil {
		set["cityId"] = *request.CityId
	}
	if request.ZoneId != nil {
		set["zoneId"] = *request.ZoneId
	}
	if request.SublocalityId != nil {
		set["sublocalityId"] = *request.SublocalityId
	}

	if request.MealType != "" {
		prefix := request.MealType.BsonKey()
		if request.RatSeconds != nil {
			set[prefix+".rat.seconds"] = *request.RatSeconds
		}
		if request.RatSampleCount != nil {
			set[prefix+".rat.sampleCount"] = *request.RatSampleCount
		}
		if request.KptSeconds != nil {
			set[prefix+".kpt.seconds"] = *request.KptSeconds
		}
		if request.KptSampleCount != nil {
			set[prefix+".kpt.sampleCount"] = *request.KptSampleCount
		}
		if request.PickupSeconds != nil {
			set[prefix+".pickup.seconds"] = *request.PickupSeconds
		}
		if request.PickupSampleCount != nil {
			set[prefix+".pickup.sampleCount"] = *request.PickupSampleCount
		}
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
	}

	set := bson.M{
		"updatedAt": request.UpdatedAt,
	}
	if request.Days != nil {
		set["day"] = *request.Days
	}
	if request.ZoneId != nil {
		set["zoneId"] = *request.ZoneId
	}
	if request.CityId != nil {
		set["cityId"] = *request.CityId
	}

	if request.MealType != "" {
		prefix := request.MealType.BsonKey()
		if request.CatSeconds != nil {
			set[prefix+".cat.seconds"] = *request.CatSeconds
		}
		if request.CatSampleCount != nil {
			set[prefix+".cat.sampleCount"] = *request.CatSampleCount
		}
		if request.FmSeconds != nil {
			set[prefix+".fm.seconds"] = *request.FmSeconds
		}
		if request.FmSampleCount != nil {
			set[prefix+".fm.sampleCount"] = *request.FmSampleCount
		}
	}

	_, err := r.mongoRepository.UpdateOne(constants.ETA_SUBLOCALITY_ESTIMATES, filter, bson.M{"$set": set}, nil)
	if err != nil {
		return err
	}
	return nil
}

// Helpers

// restaurantMealProjection limits the fetched fields to only what's needed
// to resolve an ETA for the given mealType: the join key (sublocalityId)
// and the single relevant nested meal section, instead of pulling all 5
// meal-type sections over the wire.
func restaurantMealProjection(mealType constants.MealType) bson.M {
	projection := bson.M{
		"restaurantId":  1,
		"sublocalityId": 1,
	}
	if key := mealType.BsonKey(); key != "" {
		projection[key] = 1
	}
	return projection
}

// sublocalityMealProjection limits the fetched fields to only what's needed
// to resolve an ETA for the given mealType: the join key (sublocalityId)
// and the single relevant nested meal section.
func sublocalityMealProjection(mealType constants.MealType) bson.M {
	projection := bson.M{
		"sublocalityId": 1,
	}
	if key := mealType.BsonKey(); key != "" {
		projection[key] = 1
	}
	return projection
}

// restaurantInsertRequestToEstimates maps an insert request to the
// persisted/domain estimates type.
func restaurantInsertRequestToEstimates(request *types.InsertRestaurantEstimateRequest) *types.EtaRestaurantEstimates {
	return &types.EtaRestaurantEstimates{
		RestaurantId:  request.RestaurantId,
		DayType:       request.DayType,
		CityId:        request.CityId,
		ZoneId:        request.ZoneId,
		SublocalityId: request.SublocalityId,
		Breakfast:     request.Breakfast,
		Lunch:         request.Lunch,
		Snacks:        request.Snacks,
		Dinner:        request.Dinner,
		Latenight:     request.Latenight,
		UpdatedAt:     request.UpdatedAt,
	}
}

// sublocalityInsertRequestToEstimates maps an insert request to the
// persisted/domain estimates type.
func sublocalityInsertRequestToEstimates(request *types.InsertSublocalityEstimateRequest) *types.EtaSublocalityEstimates {
	return &types.EtaSublocalityEstimates{
		SublocalityId: request.SublocalityId,
		Day:           request.Day,
		ZoneId:        request.ZoneId,
		CityId:        request.CityId,
		Breakfast:     request.Breakfast,
		Lunch:         request.Lunch,
		Snacks:        request.Snacks,
		Dinner:        request.Dinner,
		Latenight:     request.Latenight,
		UpdatedAt:     request.UpdatedAt,
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
