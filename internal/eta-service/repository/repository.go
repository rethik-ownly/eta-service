package repository

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/nutanalabs/eta-service/internal/config"
	"github.com/nutanalabs/eta-service/internal/constants"
	"github.com/nutanalabs/eta-service/internal/dataclients/kafka"
	"github.com/nutanalabs/eta-service/internal/dataclients/mongo"
	"github.com/nutanalabs/eta-service/internal/types"
	logger "github.com/nutanalabs/rapido-logger-go"
	"github.com/nutanalabs/rapido-mongo-go/mongo/results"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//go:generate mockgen -source=./repository.go -destination=./repository_mock.go -package=repository
type Repository interface {
	FetchRestaurantComponents(restaurantId string, day constants.Day, mealType constants.MealType) (*types.RestaurantComponents, error)
	FetchSublocalityComponents(sublocalityId string, day constants.Day, mealType constants.MealType) (*types.SublocalityComponents, error)
	FetchRestaurantComponentsByIDs(restaurantsId []string, day constants.Day, mealType constants.MealType, batchSize int) ([]types.RestaurantComponents, error)
	FetchSublocalityComponentsByIDs(sublocalitiesId []string, day constants.Day, mealType constants.MealType, batchSize int) ([]types.SublocalityComponents, error)

	RestaurantDayOverlapExists(restaurantId string, days []constants.Day) (bool, error)
	SublocalityDayOverlapExists(sublocalityId string, days []constants.Day) (bool, error)

	InsertRestaurantComponents(request *types.InsertRestaurantComponentsRequest) error
	InsertSublocalityComponents(request *types.InsertSublocalityComponentsRequest) error

	UpdateRestaurantComponents(restaurantId string, request *types.UpdateRestaurantComponentsRequest) error
	UpdateSublocalityComponents(sublocalityId string, request *types.UpdateSublocalityComponentsRequest) error

	PublishFetchEtaEvent(eventType constants.EventType, request *types.FetchEtaRequest, response []types.FetchEtaResponse, err error)
}

type repositoryImpl struct {
	mongoRepository mongo.Repository
	kafkaRepository kafka.Repository
	config          *config.Config
}

func NewRepository(mongoRepository mongo.Repository, kafkaRepository kafka.Repository, config *config.Config) Repository {
	return &repositoryImpl{
		mongoRepository: mongoRepository,
		kafkaRepository: kafkaRepository,
		config:          config,
	}
}

func (r *repositoryImpl) FetchRestaurantComponents(restaurantId string, day constants.Day, mealType constants.MealType) (*types.RestaurantComponents, error) {
	filter := bson.M{
		"restaurantId": restaurantId,
		"dayType":      day,
	}
	opts := options.FindOne().SetProjection(restaurantMealProjection(mealType))

	result := r.mongoRepository.FindOne(constants.RESTAURANT_COMPONENTS, filter, opts)

	var components types.RestaurantComponents
	if err := result.Decode(&components); err != nil {
		return nil, err
	}
	return &components, nil
}

func (r *repositoryImpl) FetchSublocalityComponents(sublocalityId string, day constants.Day, mealType constants.MealType) (*types.SublocalityComponents, error) {
	filter := bson.M{
		"sublocalityId": sublocalityId,
		"dayType":       day,
	}
	opts := options.FindOne().SetProjection(sublocalityMealProjection(mealType))

	result := r.mongoRepository.FindOne(constants.SUBLOCALITY_COMPONENTS, filter, opts)

	var components types.SublocalityComponents
	if err := result.Decode(&components); err != nil {
		return nil, err
	}
	return &components, nil
}

func (r *repositoryImpl) FetchRestaurantComponentsByIDs(restaurantsId []string, day constants.Day, mealType constants.MealType, batchSize int) ([]types.RestaurantComponents, error) {
	if batchSize <= 0 {
		return nil, fmt.Errorf("invalid mongo queryBatchSize: %d", batchSize)
	}

	return fetchInBatches(restaurantsId, batchSize, func(batch []string) ([]types.RestaurantComponents, error) {
		return r.fetchRestaurantComponentsByIDsBatch(batch, day, mealType)
	})
}

func (r *repositoryImpl) fetchRestaurantComponentsByIDsBatch(restaurantsId []string, day constants.Day, mealType constants.MealType) ([]types.RestaurantComponents, error) {
	filter := bson.M{
		"restaurantId": bson.M{
			"$in": restaurantsId,
		},
		"dayType": day,
	}
	opts := options.Find().SetProjection(restaurantMealProjection(mealType))

	var restaurantComponents []types.RestaurantComponents
	queryResponse := results.QueryResponse{
		Data: &restaurantComponents,
	}

	err := r.mongoRepository.FindMany(constants.RESTAURANT_COMPONENTS, filter, &queryResponse, opts)

	if err != nil {
		return nil, err
	}

	if data, ok := queryResponse.Data.(*[]types.RestaurantComponents); ok && data != nil {
		return *data, nil
	}

	return []types.RestaurantComponents{}, nil
}

func (r *repositoryImpl) FetchSublocalityComponentsByIDs(sublocalitiesId []string, day constants.Day, mealType constants.MealType, batchSize int) ([]types.SublocalityComponents, error) {
	if batchSize <= 0 {
		return nil, fmt.Errorf("invalid mongo queryBatchSize: %d", batchSize)
	}

	return fetchInBatches(sublocalitiesId, batchSize, func(batch []string) ([]types.SublocalityComponents, error) {
		return r.fetchSublocalityComponentsByIDsBatch(batch, day, mealType)
	})
}

func (r *repositoryImpl) fetchSublocalityComponentsByIDsBatch(sublocalitiesId []string, day constants.Day, mealType constants.MealType) ([]types.SublocalityComponents, error) {
	filter := bson.M{
		"sublocalityId": bson.M{
			"$in": sublocalitiesId,
		},
		"dayType": day,
	}
	opts := options.Find().SetProjection(sublocalityMealProjection(mealType))

	var sublocalityComponents []types.SublocalityComponents
	queryResponse := results.QueryResponse{
		Data: &sublocalityComponents,
	}

	err := r.mongoRepository.FindMany(constants.SUBLOCALITY_COMPONENTS, filter, &queryResponse, opts)

	if err != nil {
		return nil, err
	}

	if data, ok := queryResponse.Data.(*[]types.SublocalityComponents); ok && data != nil {
		return *data, nil
	}

	return []types.SublocalityComponents{}, nil
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

	count, err := r.mongoRepository.CountDocuments(constants.RESTAURANT_COMPONENTS, filter, nil)
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
		"dayType": bson.M{
			"$in": days,
		},
	}

	count, err := r.mongoRepository.CountDocuments(constants.SUBLOCALITY_COMPONENTS, filter, nil)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *repositoryImpl) InsertRestaurantComponents(request *types.InsertRestaurantComponentsRequest) error {
	components := restaurantInsertRequestToComponents(request)
	_, err := r.mongoRepository.InsertOne(constants.RESTAURANT_COMPONENTS, components)
	if err != nil {
		return err
	}
	return nil
}

func (r *repositoryImpl) InsertSublocalityComponents(request *types.InsertSublocalityComponentsRequest) error {
	components := sublocalityInsertRequestToComponents(request)
	_, err := r.mongoRepository.InsertOne(constants.SUBLOCALITY_COMPONENTS, components)
	if err != nil {
		return err
	}
	return nil
}

func (r *repositoryImpl) UpdateRestaurantComponents(restaurantId string, request *types.UpdateRestaurantComponentsRequest) error {
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
		if request.DelayDispatchSeconds != nil {
			set[prefix+".delayDispatch.seconds"] = *request.DelayDispatchSeconds
		}
		if request.DelayDispatchSampleCount != nil {
			set[prefix+".delayDispatch.sampleCount"] = *request.DelayDispatchSampleCount
		}
	}

	_, err := r.mongoRepository.UpdateOne(constants.RESTAURANT_COMPONENTS, filter, bson.M{"$set": set}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (r *repositoryImpl) UpdateSublocalityComponents(sublocalityId string, request *types.UpdateSublocalityComponentsRequest) error {
	filter := bson.M{
		"sublocalityId": sublocalityId,
		"dayType":       request.Day,
	}

	set := bson.M{
		"updatedAt": request.UpdatedAt,
	}
	if request.Days != nil {
		set["dayType"] = *request.Days
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

	_, err := r.mongoRepository.UpdateOne(constants.SUBLOCALITY_COMPONENTS, filter, bson.M{"$set": set}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (r *repositoryImpl) PublishFetchEtaEvent(eventType constants.EventType, request *types.FetchEtaRequest, response []types.FetchEtaResponse, fetchErr error) {
	if !r.config.GetKafkaShadowEventsEnabled() {
		return
	}

	event := fetchEtaRequestToAnalyticsEvent(eventType, request, response, fetchErr)

	go func() {
		payload, marshalErr := json.Marshal(event)
		if marshalErr != nil {
			logger.Error(logger.Format{
				Event:   "MARSHAL_FETCH_ETA_ANALYTICS_EVENT",
				Message: fmt.Sprintf("error marshaling event: %v", marshalErr),
			})
			return
		}
		if sendErr := r.kafkaRepository.SendMessage(constants.FETCH_ETA_EVENTS_TOPIC, payload); sendErr != nil {
			logger.Error(logger.Format{
				Event:   "PUBLISH_FETCH_ETA_ANALYTICS_EVENT",
				Message: fmt.Sprintf("error publishing event: %v", sendErr),
			})
		}
	}()
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

// restaurantInsertRequestToComponents maps an insert request to the
// persisted/domain components type.
func restaurantInsertRequestToComponents(request *types.InsertRestaurantComponentsRequest) *types.RestaurantComponents {
	return &types.RestaurantComponents{
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

// sublocalityInsertRequestToComponents maps an insert request to the
// persisted/domain components type.
func sublocalityInsertRequestToComponents(request *types.InsertSublocalityComponentsRequest) *types.SublocalityComponents {
	return &types.SublocalityComponents{
		SublocalityId: request.SublocalityId,
		DayType:       request.DayType,
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

func fetchEtaRequestToAnalyticsEvent(eventType constants.EventType, request *types.FetchEtaRequest, response []types.FetchEtaResponse, fetchErr error) types.FetchEtaAnalyticsEvent {
	event := types.FetchEtaAnalyticsEvent{
		EventId:   uuid.New().String(),
		CreatedAt: time.Now().UnixMilli(),
		EventType: eventType,
		Status:    "SUCCESS",
		Response:  response,
	}

	if request != nil {
		event.RequestId = request.RequestId
		event.OrderId = request.OrderId
		event.Surface = request.Surface
		event.DeliveryType = request.DeliveryType
		event.Request = *request
	}

	if fetchErr != nil {
		event.Status = "ERROR"
		event.ErrorMessage = fetchErr.Error()
	}

	return event
}
