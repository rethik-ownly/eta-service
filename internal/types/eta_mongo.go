package types

import (
	"github.com/nutanalabs/eta-service/internal/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TimeSample struct {
	Seconds     float64 `json:"seconds,omitempty" bson:"seconds,omitempty"`
	SampleCount int     `json:"sampleCount,omitempty" bson:"sampleCount,omitempty"`
}

type RestaurantMealComponents struct {
	Rat           TimeSample `json:"rat" bson:"rat"`
	Kpt           TimeSample `json:"kpt" bson:"kpt"`
	Pickup        TimeSample `json:"pickup" bson:"pickup"`
	DelayDispatch TimeSample `json:"delayDispatch" bson:"delayDispatch"`
}

type SublocalityMealComponents struct {
	Cat TimeSample `json:"cat" bson:"cat"`
	Fm  TimeSample `json:"fm" bson:"fm"`
}

type RestaurantComponents struct {
	Id            primitive.ObjectID       `json:"-" bson:"_id,omitempty"`
	RestaurantId  string                   `json:"restaurantId" bson:"restaurantId"`
	DayType       []constants.Day          `json:"dayType" bson:"dayType"`
	CityId        string                   `json:"cityId" bson:"cityId"`
	ZoneId        string                   `json:"zoneId" bson:"zoneId"`
	SublocalityId string                   `json:"sublocalityId" bson:"sublocalityId"`
	Breakfast     RestaurantMealComponents `json:"breakfast" bson:"breakfast,omitempty"`
	Lunch         RestaurantMealComponents `json:"lunch" bson:"lunch,omitempty"`
	Snacks        RestaurantMealComponents `json:"snacks" bson:"snacks,omitempty"`
	Dinner        RestaurantMealComponents `json:"dinner" bson:"dinner,omitempty"`
	Latenight     RestaurantMealComponents `json:"latenight" bson:"latenight,omitempty"`
	UpdatedAt     float64                  `json:"updatedAt,omitempty" bson:"updatedAt"`
}

// MealSection returns the nested rat/kpt/pickup/delayDispatch section matching mealType.
func (c *RestaurantComponents) MealSection(mealType constants.MealType) RestaurantMealComponents {
	switch mealType {
	case constants.Breakfast:
		return c.Breakfast
	case constants.Lunch:
		return c.Lunch
	case constants.Snack:
		return c.Snacks
	case constants.Dinner:
		return c.Dinner
	case constants.LateNight:
		return c.Latenight
	default:
		return RestaurantMealComponents{}
	}
}

type SublocalityComponents struct {
	Id            primitive.ObjectID        `json:"-" bson:"_id,omitempty"`
	SublocalityId string                    `json:"sublocalityId" bson:"sublocalityId"`
	DayType       []constants.Day           `json:"dayType" bson:"dayType"`
	ZoneId        string                    `json:"zoneId,omitempty" bson:"zoneId,omitempty"`
	CityId        string                    `json:"cityId" bson:"cityId"`
	Breakfast     SublocalityMealComponents `json:"breakfast" bson:"breakfast,omitempty"`
	Lunch         SublocalityMealComponents `json:"lunch" bson:"lunch,omitempty"`
	Snacks        SublocalityMealComponents `json:"snacks" bson:"snacks,omitempty"`
	Dinner        SublocalityMealComponents `json:"dinner" bson:"dinner,omitempty"`
	Latenight     SublocalityMealComponents `json:"latenight" bson:"latenight,omitempty"`
	UpdatedAt     float64                   `json:"updatedAt" bson:"updatedAt"`
}

// MealSection returns the nested cat/fm section matching mealType.
func (c *SublocalityComponents) MealSection(mealType constants.MealType) SublocalityMealComponents {
	switch mealType {
	case constants.Breakfast:
		return c.Breakfast
	case constants.Lunch:
		return c.Lunch
	case constants.Snack:
		return c.Snacks
	case constants.Dinner:
		return c.Dinner
	case constants.LateNight:
		return c.Latenight
	default:
		return SublocalityMealComponents{}
	}
}
