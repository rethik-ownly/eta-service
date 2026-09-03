package types

import (
	"github.com/nutanalabs/eta-service/internal/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TimeSample struct {
	Seconds     float64 `json:"seconds,omitempty" bson:"seconds,omitempty"`
	SampleCount int     `json:"sampleCount,omitempty" bson:"sampleCount,omitempty"`
}

type RestaurantMealEstimate struct {
	Rat           TimeSample `json:"rat" bson:"rat"`
	Kpt           TimeSample `json:"kpt" bson:"kpt"`
	Pickup        TimeSample `json:"pickup" bson:"pickup"`
	DelayDispatch TimeSample `json:"delayDispatch" bson:"delayDispatch"`
}


type SublocalityMealEstimate struct {
	Cat TimeSample `json:"cat" bson:"cat"`
	Fm  TimeSample `json:"fm" bson:"fm"`
}

type EtaRestaurantEstimates struct {
	Id            primitive.ObjectID     `json:"-" bson:"_id,omitempty"`
	RestaurantId  string                 `json:"restaurantId" bson:"restaurantId"`
	DayType       []constants.Day        `json:"dayType" bson:"dayType"`
	CityId        string                 `json:"cityId" bson:"cityId"`
	ZoneId        string                 `json:"zoneId" bson:"zoneId"`
	SublocalityId string                 `json:"sublocalityId" bson:"sublocalityId"`
	Breakfast     RestaurantMealEstimate `json:"breakfast" bson:"breakfast,omitempty"`
	Lunch         RestaurantMealEstimate `json:"lunch" bson:"lunch,omitempty"`
	Snacks        RestaurantMealEstimate `json:"snacks" bson:"snacks,omitempty"`
	Dinner        RestaurantMealEstimate `json:"dinner" bson:"dinner,omitempty"`
	Latenight     RestaurantMealEstimate `json:"latenight" bson:"latenight,omitempty"`
	UpdatedAt     float64                `json:"updatedAt,omitempty" bson:"updatedAt"`
}

// MealSection returns the nested rat/kpt/pickup/delayDispatch section matching mealType.
func (e *EtaRestaurantEstimates) MealSection(mealType constants.MealType) RestaurantMealEstimate {
	switch mealType {
	case constants.Breakfast:
		return e.Breakfast
	case constants.Lunch:
		return e.Lunch
	case constants.Snack:
		return e.Snacks
	case constants.Dinner:
		return e.Dinner
	case constants.LateNight:
		return e.Latenight
	default:
		return RestaurantMealEstimate{}
	}
}

type EtaSublocalityEstimates struct {
	Id            primitive.ObjectID      `json:"-" bson:"_id,omitempty"`
	SublocalityId string                  `json:"sublocalityId" bson:"sublocalityId"`
	DayType       []constants.Day         `json:"dayType" bson:"dayType"`
	ZoneId        string                  `json:"zoneId,omitempty" bson:"zoneId,omitempty"`
	CityId        string                  `json:"cityId" bson:"cityId"`
	Breakfast     SublocalityMealEstimate `json:"breakfast" bson:"breakfast,omitempty"`
	Lunch         SublocalityMealEstimate `json:"lunch" bson:"lunch,omitempty"`
	Snacks        SublocalityMealEstimate `json:"snacks" bson:"snacks,omitempty"`
	Dinner        SublocalityMealEstimate `json:"dinner" bson:"dinner,omitempty"`
	Latenight     SublocalityMealEstimate `json:"latenight" bson:"latenight,omitempty"`
	UpdatedAt     float64                 `json:"updatedAt" bson:"updatedAt"`
}

// MealSection returns the nested cat/fm section matching mealType.
func (e *EtaSublocalityEstimates) MealSection(mealType constants.MealType) SublocalityMealEstimate {
	switch mealType {
	case constants.Breakfast:
		return e.Breakfast
	case constants.Lunch:
		return e.Lunch
	case constants.Snack:
		return e.Snacks
	case constants.Dinner:
		return e.Dinner
	case constants.LateNight:
		return e.Latenight
	default:
		return SublocalityMealEstimate{}
	}
}