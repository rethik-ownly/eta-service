package types

import (
	"github.com/nutanalabs/eta-service/internal/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EtaRestaurantEstimates is both the domain type used by the
// handler/service layers and the Mongo persistence shape for the
// eta_restaurant_estimates collection.
type EtaRestaurantEstimates struct {
	ID                primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	RestaurantID      string             `json:"restaurantId" bson:"restaurantId"`
	RatSeconds        float64            `json:"ratSeconds,omitempty" bson:"ratSeconds,omitempty"`
	RatSampleCount    int                `json:"ratSampleCount,omitempty" bson:"ratSampleCount,omitempty"`
	KptSeconds        float64            `json:"kptSeconds,omitempty" bson:"kptSeconds,omitempty"`
	KptSampleCount    int                `json:"kptSampleCount,omitempty" bson:"kptSampleCount,omitempty"`
	PickupSeconds     float64            `json:"pickupSeconds,omitempty" bson:"pickupSeconds,omitempty"`
	PickupSampleCount int                `json:"pickupSampleCount,omitempty" bson:"pickupSampleCount,omitempty"`
	MealType          constants.MealType `json:"mealType" bson:"mealType"`
	Day               constants.Day      `json:"day" bson:"day"`
	CityID            string             `json:"cityId" bson:"cityId"`
	ZoneID            string             `json:"zoneId" bson:"zoneId"`
	SublocalityID     string             `json:"sublocalityId" bson:"sublocalityId"`
	UpdatedAt         float64            `json:"updatedAt,omitempty" bson:"updatedAt"`
}

// EtaSublocalityEstimates is both the domain type used by the
// handler/service layers and the Mongo persistence shape for the
// eta_sublocality_estimates collection.
type EtaSublocalityEstimates struct {
	ID             primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	SublocalityID  string             `json:"sublocalityId" bson:"sublocalityId"`
	ZoneID         string             `json:"zoneId,omitempty" bson:"zoneId,omitempty"`
	CityID         string             `json:"cityId" bson:"cityId"`
	MealType       constants.MealType `json:"mealType" bson:"mealType"`
	Day            constants.Day      `json:"day" bson:"day"`
	CatSeconds     float64            `json:"catSeconds,omitempty" bson:"catSeconds,omitempty"`
	CatSampleCount int                `json:"catSampleCount,omitempty" bson:"catSampleCount,omitempty"`
	FmSeconds      float64            `json:"fmSeconds,omitempty" bson:"fmSeconds,omitempty"`
	FmSampleCount  int                `json:"fmSampleCount" bson:"fmSampleCount"`
	UpdatedAt      float64            `json:"updatedAt" bson:"updatedAt"`
}

// EtaPlatformDefaults is both the domain type used by the handler/service
// layers and the Mongo persistence shape for the eta_platform_defaults
// collection.
type EtaPlatformDefaults struct {
	ID            primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	CityID        string             `json:"cityId" bson:"cityId"`
	RatSeconds    float64            `json:"ratSeconds,omitempty" bson:"ratSeconds,omitempty"`
	KptSeconds    float64            `json:"kptSeconds,omitempty" bson:"kptSeconds,omitempty"`
	PickupSeconds float64            `json:"pickupSeconds,omitempty" bson:"pickupSeconds,omitempty"`
	CatSeconds    float64            `json:"catSeconds,omitempty" bson:"catSeconds,omitempty"`
	FmSeconds     float64            `json:"fmSeconds,omitempty" bson:"fmSeconds,omitempty"`
	UpdatedAt     float64            `json:"updatedAt" bson:"updatedAt"`
}
