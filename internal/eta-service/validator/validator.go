package validator

import (
	"github.com/nutanalabs/eta-service/internal/config"
	"github.com/nutanalabs/eta-service/internal/types"
)

type Validator interface {
	ValidateFetchEtaRequest(request *types.FetchEtaRequest) error
	ValidateInsertRestaurantEstimateRequest(request *types.InsertRestaurantEstimateRequest) error
	ValidateInsertSublocalityEstimateRequest(request *types.InsertSublocalityEstimateRequest) error
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