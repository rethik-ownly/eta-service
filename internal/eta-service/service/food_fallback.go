package service

import (
	"github.com/nutanalabs/eta-service/internal/constants"
	"github.com/nutanalabs/eta-service/internal/types"
)

func (s *serviceImpl) defaultRestaurantMeal() types.RestaurantMealEstimate {
	d := s.config.EtaDefaultEstimates
	return types.RestaurantMealEstimate{
		Rat:           types.TimeSample{Seconds: float64(d.RestaurantAcceptanceTime)},
		Kpt:           types.TimeSample{Seconds: float64(d.KitchenPreparationTime)},
		Pickup:        types.TimeSample{Seconds: float64(d.PickupTime)},
		DelayDispatch: types.TimeSample{Seconds: float64(d.DelayDispatchTime)},
	}
}

func (s *serviceImpl) defaultSublocalityMeal() types.SublocalityMealEstimate {
	d := s.config.EtaDefaultEstimates
	return types.SublocalityMealEstimate{
		Cat: types.TimeSample{Seconds: float64(d.CaptainAssignmentTime)},
		Fm:  types.TimeSample{Seconds: float64(d.FirstMileTime)},
	}
}

func isRestaurantMealPopulated(meal types.RestaurantMealEstimate) bool {
	return meal.Rat.SampleCount > 0 ||
		meal.Kpt.SampleCount > 0 ||
		meal.Pickup.SampleCount > 0 ||
		meal.DelayDispatch.SampleCount > 0
}

func isSublocalityMealPopulated(meal types.SublocalityMealEstimate) bool {
	return meal.Cat.SampleCount > 0 || meal.Fm.SampleCount > 0
}

func (s *serviceImpl) resolveRestaurantFood(
	restaurantID string,
	mealType constants.MealType,
	estimatesByRestaurant map[string]types.EtaRestaurantEstimates,
	restaurantMongoFailed bool,
) (types.RestaurantMealEstimate, bool) {
	defaults := s.defaultRestaurantMeal()

	if restaurantMongoFailed {
		return defaults, true
	}

	restaurantEstimates, ok := estimatesByRestaurant[restaurantID]
	if !ok {
		return defaults, true
	}

	meal := restaurantEstimates.MealSection(mealType)
	if !isRestaurantMealPopulated(meal) {
		return defaults, true
	}

	return meal, false
}

func (s *serviceImpl) resolveSublocalityFood(
	sublocalityID string,
	mealType constants.MealType,
	estimatesBySublocality map[string]types.EtaSublocalityEstimates,
	sublocalityMongoFailed bool,
) (types.SublocalityMealEstimate, bool) {
	defaults := s.defaultSublocalityMeal()

	if sublocalityMongoFailed {
		return defaults, true
	}

	if sublocalityID == "" {
		return defaults, true
	}

	sublocalityEstimates, ok := estimatesBySublocality[sublocalityID]
	if !ok {
		return defaults, true
	}

	meal := sublocalityEstimates.MealSection(mealType)
	if !isSublocalityMealPopulated(meal) {
		return defaults, true
	}

	return meal, false
}

func etaSourceFromFallback(usedFallback bool) string {
	if usedFallback {
		return constants.EtaSourceFallback
	}
	return constants.EtaSourceHistoric
}
