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

type InsertETARequest struct {
	RestaurantID string  `json:"restaurantId" bson:"restaurantId"`
	DayOfWeek    string  `json:"dayOfWeek" bson:"dayOfWeek"`
	TimeSlot     string  `json:"timeSlot" bson:"timeSlot"`
	Lat          float64 `json:"lat" bson:"lat"`
	Lng        float64 `json:"lng" bson:"lng"`
	Eta          uint    `json:"eta" bson:"eta"`
	CreatedAt    float64 `json:"createdAt" bson:"createdAt"`
	UpdatedAt    float64 `json:"updatedAt" bson:"updatedAt"`
}

type FetchEtaResponse struct {
	RestaurantID string `json:"restaurantId"`
	EtaInSeconds uint   `json:"etaInSeconds"`
	DisplayMin   uint   `json:"displayMin,omitempty"`
	DisplayMax   uint   `json:"displayMax,omitempty"`
}

type EtaRestaurantEstimates struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	RestaurantID      string             `bson:"restaurantId"`
	RatSeconds        float64           `bson:"ratSeconds,omitempty"`
	RatSampleCount    int               `bson:"ratSampleCount,omitempty"`
	KptSeconds        float64           `bson:"kptSeconds,omitempty"`
	KptSampleCount    int               `bson:"kptSampleCount,omitempty"`
	PickupSeconds     float64           `bson:"pickupSeconds,omitempty"`
	PickupSampleCount int               `bson:"pickupSampleCount,omitempty"`
	MealType          constants.MealType `bson:"mealType"`
	Day               constants.Day      `bson:"day"`
	CityID            string             `bson:"cityId"`
	ZoneID            string             `bson:"zoneId"`
	SublocalityID     string             `bson:"sublocalityId"`
}

type EtaSublocalityEstimates struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	SublocalityID  string             `bson:"sublocalityId"`
	ZoneID         string             `bson:"zoneId,omitempty"`
	CityID         string             `bson:"cityId"`
	MealType       constants.MealType `bson:"mealType"`
	Day            constants.Day      `bson:"day"`
	CatSeconds     float64           `bson:"catSeconds,omitempty"`
	CatSampleCount int               `bson:"catSampleCount,omitempty"`
	FmSeconds      float64           `bson:"fmSeconds,omitempty"`
	FmSampleCount  int               `bson:"fmSampleCount"`
	UpdatedAt      float64            `bson:"updatedAt"`
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