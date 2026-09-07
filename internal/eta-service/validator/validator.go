package validator

import (
	"github.com/nutanalabs/eta-service/internal/config"
	"github.com/nutanalabs/eta-service/internal/constants"
	"github.com/nutanalabs/eta-service/internal/types"
)

type Validator interface {
	ValidateFetchEtaRequest(request *types.FetchEtaRequest) error
	ValidateInsertRestaurantComponentsRequest(request *types.InsertRestaurantComponentsRequest) error
	ValidateInsertSublocalityComponentsRequest(request *types.InsertSublocalityComponentsRequest) error

	ValidateUpdateRestaurantComponents(restaurantId string, request *types.UpdateRestaurantComponentsRequest) error
	ValidateUpdateSublocalityComponents(sublocalityId string, request *types.UpdateSublocalityComponentsRequest) error
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
	if request.UserLocation.Lat == 0 && request.UserLocation.Lng == 0 {
		return types.NewBadRequestError("invalid user location")
	}
	if request.Entities == nil || len(request.Entities) > v.config.Eta.MaxEntities {
		return types.NewBadRequestError("invalid entities")
	}
	if request.Surface == constants.Cart && request.OrderId == "" {
		return types.NewBadRequestError("invalid orderId")
	}
	return nil
}

func (v *validatorImpl) ValidateInsertRestaurantComponentsRequest(request *types.InsertRestaurantComponentsRequest) error {
	if request.RestaurantId == "" {
		return types.NewBadRequestError("invalid restaurantId")
	}
	if request.CityId == "" {
		return types.NewBadRequestError("invalid cityId")
	}
	if request.ZoneId == "" {
		return types.NewBadRequestError("invalid zoneId")
	}
	if request.SublocalityId == "" {
		return types.NewBadRequestError("invalid sublocalityId")
	}
	if len(request.DayType) == 0 {
		return types.NewBadRequestError("invalid dayType")
	}
	for _, day := range request.DayType {
		if !day.IsValid() {
			return types.NewBadRequestError("invalid dayType")
		}
	}
	if err := validateRestaurantMealComponents(request.Breakfast); err != nil {
		return err
	}
	if err := validateRestaurantMealComponents(request.Lunch); err != nil {
		return err
	}
	if err := validateRestaurantMealComponents(request.Snacks); err != nil {
		return err
	}
	if err := validateRestaurantMealComponents(request.Dinner); err != nil {
		return err
	}
	if err := validateRestaurantMealComponents(request.Latenight); err != nil {
		return err
	}
	return nil
}

func (v *validatorImpl) ValidateInsertSublocalityComponentsRequest(request *types.InsertSublocalityComponentsRequest) error {
	if request.SublocalityId == "" {
		return types.NewBadRequestError("invalid sublocalityId")
	}
	if request.ZoneId == "" {
		return types.NewBadRequestError("invalid zoneId")
	}
	if request.CityId == "" {
		return types.NewBadRequestError("invalid cityId")
	}
	if len(request.DayType) == 0 {
		return types.NewBadRequestError("invalid dayType")
	}
	for _, day := range request.DayType {
		if !day.IsValid() {
			return types.NewBadRequestError("invalid dayType")
		}
	}
	if err := validateSublocalityMealComponents(request.Breakfast); err != nil {
		return err
	}
	if err := validateSublocalityMealComponents(request.Lunch); err != nil {
		return err
	}
	if err := validateSublocalityMealComponents(request.Snacks); err != nil {
		return err
	}
	if err := validateSublocalityMealComponents(request.Dinner); err != nil {
		return err
	}
	if err := validateSublocalityMealComponents(request.Latenight); err != nil {
		return err
	}
	return nil
}

func (v *validatorImpl) ValidateUpdateRestaurantComponents(restaurantId string, request *types.UpdateRestaurantComponentsRequest) error {
	if restaurantId == "" {
		return types.NewBadRequestError("invalid restaurantId")
	}
	if !request.Day.IsValid() {
		return types.NewBadRequestError("invalid day")
	}
	if request.Days != nil {
		if len(*request.Days) == 0 {
			return types.NewBadRequestError("invalid days")
		}
		for _, day := range *request.Days {
			if !day.IsValid() {
				return types.NewBadRequestError("invalid days")
			}
		}
	}

	hasMealFieldUpdate := request.RatSeconds != nil || request.RatSampleCount != nil ||
		request.KptSeconds != nil || request.KptSampleCount != nil ||
		request.PickupSeconds != nil || request.PickupSampleCount != nil ||
		request.DelayDispatchSeconds != nil || request.DelayDispatchSampleCount != nil

	if request.MealType != "" && !request.MealType.IsValid() {
		return types.NewBadRequestError("invalid mealType")
	}
	if hasMealFieldUpdate && request.MealType == "" {
		return types.NewBadRequestError("mealType is required to update rat/kpt/pickup/delayDispatch fields")
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
	if request.DelayDispatchSeconds != nil && *request.DelayDispatchSeconds < 0 {
		return types.NewBadRequestError("invalid delayDispatchSeconds")
	}
	if request.DelayDispatchSampleCount != nil && *request.DelayDispatchSampleCount < 0 {
		return types.NewBadRequestError("invalid delayDispatchSampleCount")
	}
	if request.CityId != nil && *request.CityId == "" {
		return types.NewBadRequestError("invalid cityId")
	}
	if request.ZoneId != nil && *request.ZoneId == "" {
		return types.NewBadRequestError("invalid zoneId")
	}
	if request.SublocalityId != nil && *request.SublocalityId == "" {
		return types.NewBadRequestError("invalid sublocalityId")
	}
	return nil
}

func (v *validatorImpl) ValidateUpdateSublocalityComponents(sublocalityId string, request *types.UpdateSublocalityComponentsRequest) error {
	if sublocalityId == "" {
		return types.NewBadRequestError("invalid sublocalityId")
	}
	if !request.Day.IsValid() {
		return types.NewBadRequestError("invalid day")
	}
	if request.Days != nil {
		if len(*request.Days) == 0 {
			return types.NewBadRequestError("invalid days")
		}
		for _, day := range *request.Days {
			if !day.IsValid() {
				return types.NewBadRequestError("invalid days")
			}
		}
	}

	hasMealFieldUpdate := request.CatSeconds != nil || request.CatSampleCount != nil ||
		request.FmSeconds != nil || request.FmSampleCount != nil

	if request.MealType != "" && !request.MealType.IsValid() {
		return types.NewBadRequestError("invalid mealType")
	}
	if hasMealFieldUpdate && request.MealType == "" {
		return types.NewBadRequestError("mealType is required to update cat/fm fields")
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
	if request.ZoneId != nil && *request.ZoneId == "" {
		return types.NewBadRequestError("invalid zoneId")
	}
	if request.CityId != nil && *request.CityId == "" {
		return types.NewBadRequestError("invalid cityId")
	}
	return nil
}

// Helpers

func validateRestaurantMealComponents(estimate types.RestaurantMealComponents) error {
	if estimate.Rat.Seconds < 0 || estimate.Rat.SampleCount < 0 {
		return types.NewBadRequestError("invalid rat estimate")
	}
	if estimate.Kpt.Seconds < 0 || estimate.Kpt.SampleCount < 0 {
		return types.NewBadRequestError("invalid kpt estimate")
	}
	if estimate.Pickup.Seconds < 0 || estimate.Pickup.SampleCount < 0 {
		return types.NewBadRequestError("invalid pickup estimate")
	}
	if estimate.DelayDispatch.Seconds < 0 || estimate.DelayDispatch.SampleCount < 0 {
		return types.NewBadRequestError("invalid delayDispatch estimate")
	}
	return nil
}

func validateSublocalityMealComponents(estimate types.SublocalityMealComponents) error {
	if estimate.Cat.Seconds < 0 || estimate.Cat.SampleCount < 0 {
		return types.NewBadRequestError("invalid cat estimate")
	}
	if estimate.Fm.Seconds < 0 || estimate.Fm.SampleCount < 0 {
		return types.NewBadRequestError("invalid fm estimate")
	}
	return nil
}
