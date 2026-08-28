package types

import (
	"github.com/nutanalabs/eta-service/internal/constants"
)

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}


// Fetch Eta Reques-Response
type FetchEtaRequestOptions struct {
	QosLevel constants.QosLevel `json:"qosLevel"`
}

type FetchEtaRequestEntity struct {
	RestaurantID       string   `json:"restaurantId"`
	RestaurantLocation Location `json:"restaurantLocation"`
}

type FetchEtaRequest struct {
	Surface      constants.Surface       `json:"surface"`
	DeliveryType constants.DeliveryType  `json:"deliveryType"`
	UserID       string                  `json:"userId"`
	UserLocation Location                `json:"userLocation"`
	Options      FetchEtaRequestOptions  `json:"options"`
	Entities     []FetchEtaRequestEntity `json:"entities"`
}

type FetchEtaResponse struct {
	RestaurantID string `json:"restaurantId"`
	EtaInSeconds uint   `json:"etaInSeconds"`
	DisplayMin   uint   `json:"displayMin,omitempty"`
	DisplayMax   uint   `json:"displayMax,omitempty"`
}

// Insert Eta Request 
type InsertRestaurantEstimateRequest struct {
	RestaurantID      string             `json:"restaurantId"`
	RatSeconds        float64            `json:"ratSeconds,omitempty"`
	RatSampleCount    int                `json:"ratSampleCount,omitempty"`
	KptSeconds        float64            `json:"kptSeconds,omitempty"`
	KptSampleCount    int                `json:"kptSampleCount,omitempty"`
	PickupSeconds     float64            `json:"pickupSeconds,omitempty"`
	PickupSampleCount int                `json:"pickupSampleCount,omitempty"`
	MealType          constants.MealType `json:"mealType"`
	Day               constants.Day      `json:"day"`
	CityID            string             `json:"cityId"`
	ZoneID            string             `json:"zoneId"`
	SublocalityID     string             `json:"sublocalityId"`
	UpdatedAt         float64            `json:"-"`
}

type InsertSublocalityEstimateRequest struct {
	SublocalityID  string             `json:"sublocalityId"`
	ZoneID         string             `json:"zoneId,omitempty"`
	CityID         string             `json:"cityId"`
	MealType       constants.MealType `json:"mealType"`
	Day            constants.Day      `json:"day"`
	CatSeconds     float64            `json:"catSeconds,omitempty"`
	CatSampleCount int                `json:"catSampleCount,omitempty"`
	FmSeconds      float64            `json:"fmSeconds,omitempty"`
	FmSampleCount  int                `json:"fmSampleCount"`
	UpdatedAt      float64            `json:"-"`
}

// RestaurantID is not bound from the request body; it is taken from the
// ":restaurantId" path param and set by the handler before validation.
type UpdateRestaurantEstimateRequest struct {
	RestaurantID string `json:"-"`

	Day      constants.Day      `json:"day" binding:"required"`
	MealType constants.MealType `json:"mealType" binding:"required"`

	RatSeconds        *float64 `json:"ratSeconds,omitempty"`
	RatSampleCount    *int     `json:"ratSampleCount,omitempty"`
	KptSeconds        *float64 `json:"kptSeconds,omitempty"`
	KptSampleCount    *int     `json:"kptSampleCount,omitempty"`
	PickupSeconds     *float64 `json:"pickupSeconds,omitempty"`
	PickupSampleCount *int     `json:"pickupSampleCount,omitempty"`
	CityID            *string  `json:"cityId,omitempty"`
	ZoneID            *string  `json:"zoneId,omitempty"`
	SublocalityID     *string  `json:"sublocalityId,omitempty"`

	UpdatedAt float64 `json:"-"`
}

// SublocalityID is not bound from the request body; it is taken from the
// ":sublocalityId" path param and set by the handler before validation.
type UpdateSublocalityEstimateRequest struct {
	SublocalityID string `json:"-"`

	Day      constants.Day      `json:"day" binding:"required"`
	MealType constants.MealType `json:"mealType" binding:"required"`

	ZoneID         *string  `json:"zoneId,omitempty"`
	CityID         *string  `json:"cityId,omitempty"`
	CatSeconds     *float64 `json:"catSeconds,omitempty"`
	CatSampleCount *int     `json:"catSampleCount,omitempty"`
	FmSeconds      *float64 `json:"fmSeconds,omitempty"`
	FmSampleCount  *int     `json:"fmSampleCount,omitempty"`

	UpdatedAt float64 `json:"-"`
}
