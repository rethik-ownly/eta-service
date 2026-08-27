package validator

import (
	"github.com/nutanalabs/eta-service/internal/config"
	"github.com/nutanalabs/eta-service/internal/types"
)

type Validator interface {
	ValidateFetchEtaRequest(request *types.FetchEtaRequest) error
	ValidateInsertRestaurantEstimateRequest(request *types.InsertRestaurantEstimateRequest) error
	ValidateInsertSublocalityEstimateRequest(request *types.InsertSublocalityEstimateRequest) error

	ValidateUpdateRestaurantEstimates(restaurantId string, request *types.UpdateRestaurantEstimateRequest) error
	ValidateUpdateSublocalityEstimates(sublocalityId string, request *types.UpdateSublocalityEstimateRequest) error
}

type validatorImpl struct {
	config *config.Config
}

func NewValidator(config *config.Config) Validator {
	return &validatorImpl{
		config: config,
	}
} 

func (v *validatorImpl) ValidateFetchEtaRequest(request *types.FetchEtaRequest) error {
	if !request.Surface.IsValid(){
		return types.NewBadRequestError("invalid surface")
	}
	if !request.DeliveryType.IsValid(){
		return types.NewBadRequestError("invalid delivery type")
	}
	if request.Options.QosLevel != "" && !request.Options.QosLevel.IsValid() {
		return types.NewBadRequestError("invalid qosLevel")
	}
	if request.UserID == "" {
		return types.NewBadRequestError("invalid userId")
	}
	if request.Entities == nil || len(request.Entities) > v.config.Eta.MaxEntities {
		return types.NewBadRequestError("invalid entities")
	}
	return nil
}

func (v *validatorImpl) ValidateInsertRestaurantEstimateRequest(request *types.InsertRestaurantEstimateRequest) error {
	if request.RestaurantID == "" {
		return types.NewBadRequestError("invalid restaurantId")
	}
	if request.CityID == "" {
		return types.NewBadRequestError("invalid cityId")
	}
	if request.ZoneID == "" {
		return types.NewBadRequestError("invalid zoneId")
	}
	if request.SublocalityID == "" {
		return types.NewBadRequestError("invalid sublocalityId")
	}
	if !request.MealType.IsValid() {
		return types.NewBadRequestError("invalid mealType")
	}
	if !request.Day.IsValid() {
		return types.NewBadRequestError("invalid day")
	}
	if request.RatSeconds < 0 || request.KptSeconds < 0 || request.PickupSeconds < 0 {
		return types.NewBadRequestError("invalid estimate seconds")
	}
	if request.RatSampleCount < 0 || request.KptSampleCount < 0 || request.PickupSampleCount < 0 {
		return types.NewBadRequestError("invalid sample count")
	}
	return nil
}

func (v *validatorImpl) ValidateInsertSublocalityEstimateRequest(request *types.InsertSublocalityEstimateRequest) error {
	if request.SublocalityID == "" {
		return types.NewBadRequestError("invalid sublocalityId")
	}
	if request.ZoneID == "" {
		return types.NewBadRequestError("invalid zoneId")
	}
	if request.CityID == "" {
		return types.NewBadRequestError("invalid cityId")
	}
	if !request.MealType.IsValid() {
		return types.NewBadRequestError("invalid mealType")
	}
	if !request.Day.IsValid() {
		return types.NewBadRequestError("invalid day")
	}
	if request.CatSeconds < 0 || request.FmSeconds < 0 {
		return types.NewBadRequestError("invalid estimate seconds")
	}
	if request.CatSampleCount < 0 || request.FmSampleCount < 0 {
		return types.NewBadRequestError("invalid sample count")
	}
	return nil
}

func (v *validatorImpl) ValidateUpdateRestaurantEstimates(restaurantId string, request *types.UpdateRestaurantEstimateRequest) error {
	if restaurantId == "" {
		return types.NewBadRequestError("invalid restaurantId")
	}
	if !request.MealType.IsValid() {
		return types.NewBadRequestError("invalid mealType")
	}
	if !request.Day.IsValid() {
		return types.NewBadRequestError("invalid day")
	}
	if request.RatSeconds != nil && *request.RatSeconds < 0 {
		return types.NewBadRequestError("invalid ratSeconds")
	}
	if request.KptSeconds != nil && *request.KptSeconds < 0 {
		return types.NewBadRequestError("invalid kptSeconds")
	}
	if request.PickupSeconds != nil && *request.PickupSeconds < 0 {
		return types.NewBadRequestError("invalid pickupSeconds")
	}
	if request.RatSampleCount != nil && *request.RatSampleCount < 0 {
		return types.NewBadRequestError("invalid ratSampleCount")
	}
	if request.KptSampleCount != nil && *request.KptSampleCount < 0 {
		return types.NewBadRequestError("invalid kptSampleCount")
	}
	if request.PickupSampleCount != nil && *request.PickupSampleCount < 0 {
		return types.NewBadRequestError("invalid pickupSampleCount")
	}
	if request.CityID != nil && *request.CityID == "" {
		return types.NewBadRequestError("invalid cityId")
	}
	if request.ZoneID != nil && *request.ZoneID == "" {
		return types.NewBadRequestError("invalid zoneId")
	}
	if request.SublocalityID != nil && *request.SublocalityID == "" {
		return types.NewBadRequestError("invalid sublocalityId")
	}
	return nil
}

func (v *validatorImpl) ValidateUpdateSublocalityEstimates(sublocalityId string, request *types.UpdateSublocalityEstimateRequest) error {
	if sublocalityId == "" {
		return types.NewBadRequestError("invalid sublocalityId")
	}
	if !request.MealType.IsValid() {
		return types.NewBadRequestError("invalid mealType")
	}
	if !request.Day.IsValid() {
		return types.NewBadRequestError("invalid day")
	}
	if request.CatSeconds != nil && *request.CatSeconds < 0 {
		return types.NewBadRequestError("invalid catSeconds")
	}
	if request.FmSeconds != nil && *request.FmSeconds < 0 {
		return types.NewBadRequestError("invalid fmSeconds")
	}
	if request.CatSampleCount != nil && *request.CatSampleCount < 0 {
		return types.NewBadRequestError("invalid catSampleCount")
	}
	if request.FmSampleCount != nil && *request.FmSampleCount < 0 {
		return types.NewBadRequestError("invalid fmSampleCount")
	}
	if request.ZoneID != nil && *request.ZoneID == "" {
		return types.NewBadRequestError("invalid zoneId")
	}
	if request.CityID != nil && *request.CityID == "" {
		return types.NewBadRequestError("invalid cityId")
	}
	return nil
}