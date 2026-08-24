package types

import (
	"github.com/nutanalabs/eta-service/internal/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
)




type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// Request types

type GetETAResponse struct {
	EtaID  				 string `json:"etaId" bson:"etaId"`
	EtaInSeconds 		 uint `json:"eta_in_seconds" bson:"eta_in_seconds"`
	CreatedAt            float64  `json:"createdAt" bson:"createdAt"`
	UpdatedAt            float64  `json:"updatedAt" bson:"updatedAt"`
}

type FetchEtaRequestOptions struct {
	QosLevel constants.QosLevel `json:"qos_level"`
}

type FetchEtaRequestEntity struct {
	RestaurantID string `json:"restaurant_id"`
	RestaurantLocation Location `json:"restaurant_location"`
}

type FetchEtaRequest struct {
	Surface      constants.Surface         `json:"surface" bson:"surface"`
	DeliveryType constants.DeliveryType	   `json:"delivery_type" bson:"delivery_type"`
	UserID       string                    `json:"user_id" bson:"user_id"`
	UserLocation Location                  `json:"user_location" bson:"user_location"`
	Options      FetchEtaRequestOptions    `json:"options" bson:"options"`
	Entities     []FetchEtaRequestEntity   `json:"entities" bson:"entities"`
}

// Response types

type InsertETARequest struct {
	RestaurantID string `json:"restaurant_id" bson:"restaurant_id"`
	DayOfWeek string `json:"day_of_week" bson:"day_of_week"`
	TimeSlot string `json:"time_slot" bson:"time_slot"`
	Lat float64 `json:"lat" bson:"lat"`
	Lon float64 `json:"lon" bson:"lon"`
	Eta uint `json:"eta" bson:"eta"`
	CreatedAt            float64  `json:"createdAt" bson:"createdAt"`
	UpdatedAt            float64  `json:"updatedAt" bson:"updatedAt"`
}

type FetchEtaResponse struct {
	RestaurantID string `json:"restaurant_id"`
	EtaInSeconds uint `json:"eta_in_seconds"`
	DisplayMin uint	`json:"display_min"`
	DisplayMax uint	`json:"display_max"`
}

//
type EtaRestaurantEstimates struct {
	ID 				  primitive.ObjectID 	`bson:"_id,omitempty"`
	RestaurantID 	  string  				`bson:"restaurant_id"`
    RatSeconds        *float64  			`bson:"rat_seconds,omitempty"`
    RatSampleCount    *int      			`bson:"rat_sample_count,omitempty"`
    KptSeconds        *float64				`bson:"kpt_seconds,omitempty"`
    KptSampleCount    *int					`bson:"kpt_sample_count,omitempty"`
    PickupSeconds     *float64				`bson:"pickup_seconds,omitempty"`
    PickupSampleCount *int					`bson:"pickup_sample_count,omitempty"`
	MealType 		  constants.MealType 	`bson:"meal_type"`
	Day 			  constants.Day			`bson:"day"`
	CityID 			  string 				`bson:"city_id"`
	ZoneID 			  string 				`bson:"zone_id"`
	SublocalityID 	  string 				`bson:"sublocality_id"`
}

type EtaSublocalityEstimates struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	SublocalityID string             `bson:"sublocality_id"`
	ZoneID        string             `bson:"zone_id,omitempty"`
    CityID        string             `bson:"city_id"`
	MealType     constants.MealType `bson:"meal_type"`
	Day 	 constants.Day       `bson:"day"`
	CatSeconds 	   *float64 		`bson:"cat_seconds,omitempty"`
	CatSampleCount *int 			`bson:"cat_sample_count,omitempty"`
	FmSeconds *float64 `bson:"fm_seconds,omitempty"`
	FmSampleCount *int `bson:"fm_sample_count"`
	UpdatedAt  float64		`bson:"updatedAt"`
}