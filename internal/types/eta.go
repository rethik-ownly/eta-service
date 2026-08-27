package types

import (
	"github.com/nutanalabs/eta-service/internal/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type FetchEtaRequestOptions struct {
	QosLevel constants.QosLevel `json:"qosLevel"`
}

type FetchEtaRequestEntity struct {
	RestaurantID       string   `json:"restaurantId"`
	RestaurantLocation Location `json:"restaurantLocation"`
}

type FetchEtaRequest struct {
	Surface      constants.Surface      `json:"surface" bson:"surface"`
	DeliveryType constants.DeliveryType `json:"deliveryType" bson:"deliveryType"`
	UserID       string                 `json:"userId" bson:"userId"`
	UserLocation Location               `json:"userLocation" bson:"userLocation"`
	Options      FetchEtaRequestOptions `json:"options" bson:"options"`
	Entities     []FetchEtaRequestEntity `json:"entities" bson:"entities"`
}

type FetchEtaResponse struct {
	RestaurantID string `json:"restaurantId"`
	EtaInSeconds uint   `json:"etaInSeconds"`
	DisplayMin   uint   `json:"displayMin,omitempty"`
	DisplayMax   uint   `json:"displayMax,omitempty"`
}

type EtaRestaurantEstimates struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	RestaurantID      string             `bson:"restaurantId" json:"restaurantId"`
	RatSeconds        float64            `bson:"ratSeconds,omitempty" json:"ratSeconds,omitempty"`
	RatSampleCount    int                `bson:"ratSampleCount,omitempty" json:"ratSampleCount,omitempty"`
	KptSeconds        float64            `bson:"kptSeconds,omitempty" json:"kptSeconds,omitempty"`
	KptSampleCount    int                `bson:"kptSampleCount,omitempty" json:"kptSampleCount,omitempty"`
	PickupSeconds     float64            `bson:"pickupSeconds,omitempty" json:"pickupSeconds,omitempty"`
	PickupSampleCount int                `bson:"pickupSampleCount,omitempty" json:"pickupSampleCount,omitempty"`
	MealType          constants.MealType `bson:"mealType" json:"mealType"`
	Day               constants.Day      `bson:"day" json:"day"`
	CityID            string             `bson:"cityId" json:"cityId"`
	ZoneID            string             `bson:"zoneId" json:"zoneId"`
	SublocalityID     string             `bson:"sublocalityId" json:"sublocalityId"`
	UpdatedAt         float64            `bson:"updatedAt" json:"updatedAt,omitempty"`
}

type EtaSublocalityEstimates struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	SublocalityID  string             `bson:"sublocalityId" json:"sublocalityId"`
	ZoneID         string             `bson:"zoneId,omitempty" json:"zoneId,omitempty"`
	CityID         string             `bson:"cityId" json:"cityId"`
	MealType       constants.MealType `bson:"mealType" json:"mealType"`
	Day            constants.Day      `bson:"day" json:"day"`
	CatSeconds     float64            `bson:"catSeconds,omitempty" json:"catSeconds,omitempty"`
	CatSampleCount int                `bson:"catSampleCount,omitempty" json:"catSampleCount,omitempty"`
	FmSeconds      float64            `bson:"fmSeconds,omitempty" json:"fmSeconds,omitempty"`
	FmSampleCount  int                `bson:"fmSampleCount" json:"fmSampleCount"`
	UpdatedAt      float64            `bson:"updatedAt" json:"updatedAt"`
}

type EtaPlatformDefaults struct {
	ID 				primitive.ObjectID  `bson:"_id,omitempty"`
	CityID 			string 				`bson:"cityId"`
	RatSeconds      float64             `bson:"ratSeconds,omitempty"`
	KptSeconds      float64             `bson:"kptSeconds,omitempty"`
	PickupSeconds   float64             `bson:"pickupSeconds,omitempty"`
	CatSeconds      float64           	`bson:"catSeconds,omitempty"`
	FmSeconds       float64           	`bson:"fmSeconds,omitempty"`
	UpdatedAt       float64            	`bson:"updatedAt"`
}

type InsertRestaurantEstimateRequest struct {
	RestaurantID      string             `bson:"restaurantId" json:"restaurantId"`
	RatSeconds        float64            `bson:"ratSeconds,omitempty" json:"ratSeconds,omitempty"`
	RatSampleCount    int                `bson:"ratSampleCount,omitempty" json:"ratSampleCount,omitempty"`
	KptSeconds        float64            `bson:"kptSeconds,omitempty" json:"kptSeconds,omitempty"`
	KptSampleCount    int                `bson:"kptSampleCount,omitempty" json:"kptSampleCount,omitempty"`
	PickupSeconds     float64            `bson:"pickupSeconds,omitempty" json:"pickupSeconds,omitempty"`
	PickupSampleCount int                `bson:"pickupSampleCount,omitempty" json:"pickupSampleCount,omitempty"`
	MealType          constants.MealType `bson:"mealType" json:"mealType"`
	Day               constants.Day      `bson:"day" json:"day"`
	CityID            string             `bson:"cityId" json:"cityId"`
	ZoneID            string             `bson:"zoneId" json:"zoneId"`
	SublocalityID     string             `bson:"sublocalityId" json:"sublocalityId"`
	UpdatedAt         float64            `bson:"updatedAt" json:"-"`
}

type InsertSublocalityEstimateRequest struct {
	SublocalityID  string             `bson:"sublocalityId" json:"sublocalityId"`
	ZoneID         string             `bson:"zoneId,omitempty" json:"zoneId,omitempty"`
	CityID         string             `bson:"cityId" json:"cityId"`
	MealType       constants.MealType `bson:"mealType" json:"mealType"`
	Day            constants.Day      `bson:"day" json:"day"`
	CatSeconds     float64            `bson:"catSeconds,omitempty" json:"catSeconds,omitempty"`
	CatSampleCount int                `bson:"catSampleCount,omitempty" json:"catSampleCount,omitempty"`
	FmSeconds      float64            `bson:"fmSeconds,omitempty" json:"fmSeconds,omitempty"`
	FmSampleCount  int                `bson:"fmSampleCount" json:"fmSampleCount"`
	UpdatedAt      float64            `bson:"updatedAt" json:"-"`
}